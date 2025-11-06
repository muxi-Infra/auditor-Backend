package ali

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	gre "github.com/alibabacloud-go/green-20220302/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/errorx"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/model"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/logger"
)

const (
	ConnectTimeout = 3000
	ReadTimeout    = 6000
)

type AlClient struct {
	AccessKeyId     string
	AccessKeySecret string
	client          *gre.Client
	runtime         *util.RuntimeOptions
	log             logger.Logger
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
		AccessKeyId: accessKeyId, AccessKeySecret: accessKeySecret,
		client: ac, runtime: &util.RuntimeOptions{}}
	for _, o := range op {
		o(c)
	}
	return c
}

func (ac *AlClient) WrapLogger(logger logger.Logger) {
	ac.log = logger
}

func (ac *AlClient) SendMessage(content string, pics []string) (model.AuditResult, error) {
	// 为提高速度，这里并发审核图片和文本。
	var (
		textResult   *model.TextAuditResult
		imageResults []model.ImageAuditResult
		textErr      error
		wg           sync.WaitGroup
	)
	// 文本审核
	wg.Add(1)
	go func() {
		defer wg.Done()
		tr, err := ac.auditText(content)
		if err != nil {
			textErr = err
			return
		}
		textResult = parseTextResponse(tr)
	}()
	// 图片审核
	if len(pics) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ir := ac.auditImages(transformPics(pics))
			imageResults = parseImageResponse(ir)
		}()
	}
	wg.Wait()
	if textErr != nil {
		return model.AuditResult{}, errorx.ErrTextAuditErr.Wrap(textErr)
	}
	re := merge(textResult, imageResults)
	return re, nil
}

func (ac *AlClient) Transform(role string, contents response.Contents) (content string, pics []string) {
	return parseContent(contents)
}

// auditImage 批量审核图片；参数详细信息请见https://www.alibabacloud.com/help/zh/content-moderation/latest/image-review-enhanced-api
func (ac *AlClient) auditImages(pics []model.ImageParameters) []*gre.ImageModerationResponseBodyData {
	var wait sync.WaitGroup
	res := make([]*gre.ImageModerationResponseBodyData, 0, 5)
	for _, pic := range pics {
		wait.Add(1)
		go func() {
			defer wait.Done()
			re, err := ac.auditImage(pic)
			if err != nil {
				ac.log.Error(err.Error())
				return
			}
			res = append(res, re)
		}()
	}
	wait.Wait()
	return res
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
		return nil, fmt.Errorf("%s  body-code:%d", *body.Msg, *body.Code)
	}
	data := body.GetData()
	return data, nil
}
