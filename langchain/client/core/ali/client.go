package ali

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/errorx"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	gre "github.com/alibabacloud-go/green-20220302/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/model"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/logger"
)

const (
	Domain            = "阿里云内容审核"
	ConnectTimeout    = 3000
	ReadTimeout       = 6000
	MaxWorks          = 50
	OnceMaxGoroutines = 10
)

type token struct{}

type AlClient struct {
	AccessKeyId     string
	AccessKeySecret string
	client          *gre.Client
	runtime         *util.RuntimeOptions
	log             logger.Logger
	workers         chan token // 并发限制器，保证一次最多同时执行50个任务
}

type AlClientOpt func(*AlClient)

func NewAlClient(accessKeyId, accessKeySecret, region, endpoint string, op ...AlClientOpt) *AlClient {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		ConnectTimeout:  tea.Int(ConnectTimeout),
		ReadTimeout:     tea.Int(ReadTimeout),
		RegionId:        tea.String(region),
		Endpoint:        tea.String(endpoint),
	}

	ac, err := gre.NewClient(config)
	if err != nil {
		panic(err)
	}

	c := &AlClient{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		client:          ac,
		runtime:         &util.RuntimeOptions{},
		workers:         make(chan token, MaxWorks),
	}

	for i := 0; i < MaxWorks; i++ {
		c.workers <- token{}
	}

	for _, o := range op {
		o(c)
	}
	return c
}

func (ac *AlClient) WrapLogger(logger logger.Logger) {
	ac.log = logger
}

func (ac *AlClient) SendMessage(content string, pics []string) (model.AuditResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	select {
	case tk := <-ac.workers:
		{
			defer func() { ac.workers <- tk }()

			var (
				textResult   *model.TextAuditResult
				imageResults []model.ImageAuditResult
			)
			g := new(errgroup.Group)
			g.SetLimit(OnceMaxGoroutines)

			if content != "" {
				g.Go(func() error {
					tr, err := ac.auditText(content)
					if err != nil {
						return err
					}
					textResult = parseTextResponse(tr)
					return nil
				})
			}

			var res []*gre.ImageModerationResponseBodyData
			if len(pics) > 0 {
				tar := transformPics(pics)

				res = make([]*gre.ImageModerationResponseBodyData, len(tar), len(tar))

				for k, pic := range tar {
					idx := k
					p := pic
					g.Go(func() error {
						re, err := ac.auditImage(p)
						if err != nil {
							return err
						}

						if re != nil {
							res[idx] = re
						}
						return nil
					})
				}
			}
			if err := g.Wait(); err != nil {
				return model.AuditResult{}, err
			}

			if len(res) > 0 {
				imageResults = parseImageResponse(res)
			}
			re := merge(textResult, imageResults)
			return re, nil
		}
	case <-ctx.Done():
		return model.AuditResult{}, errors.New("阿里审核系统繁忙，等待超时")
	}

}

func (ac *AlClient) Transform(role string, contents response.Contents) (content string, pics []string) {
	return parseContent(contents)
}

func (ac *AlClient) auditText(content string) (*gre.TextModerationPlusResponseBodyData, error) {
	serviceParameters, _ := json.Marshal(
		map[string]interface{}{
			"content": content,
		},
	)
	request := gre.TextModerationPlusRequest{
		Service:           tea.String("comment_detection_pro"),
		ServiceParameters: tea.String(string(serviceParameters)),
	}

	result, err := ac.client.TextModerationPlusWithOptions(&request, ac.runtime)
	if err != nil {
		return nil, err
	}

	if *result.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response not success. status:%d", *result.StatusCode)
	}
	body := result.Body
	if *body.Code != http.StatusOK {
		return nil, fmt.Errorf("%s  body-code:%d", *body.Message, *body.Code)
	}

	data := body.Data
	return data, nil
}

// auditImage 批量审核图片；参数详细信息请见https://www.alibabacloud.com/help/zh/content-moderation/latest/image-review-enhanced-api
func (ac *AlClient) auditImage(pic model.ImageParameters) (*gre.ImageModerationResponseBodyData, error) {
	serviceParameters, err := json.Marshal(
		pic,
	)

	if err != nil {
		return nil, err
	}

	imageModerationRequest := &gre.ImageModerationRequest{
		Service:           tea.String("baselineCheck"), // 基线检测
		ServiceParameters: tea.String(string(serviceParameters)),
	}

	result, err := ac.client.ImageModerationWithOptions(imageModerationRequest, ac.runtime)
	if err != nil {
		return nil, err
	}

	if *result.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response not success. status:%d", *result.StatusCode)
	}

	body := result.Body

	if *body.Code != http.StatusOK {
		return nil, errorx.ErrUnSupportImage.
			SetError(fmt.Errorf("%s  body-code:%d", *body.Msg, *body.Code)).
			SetMessage("该图片可能不支持：" + pic.ImageUrl)
	}

	data := body.GetData()
	return data, nil
}
