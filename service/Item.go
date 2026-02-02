package service

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/muxi-Infra/auditor-Backend/api/request"
	"github.com/muxi-Infra/auditor-Backend/pkg/apikey"
	"github.com/muxi-Infra/auditor-Backend/pkg/jwt"
	"github.com/muxi-Infra/auditor-Backend/pkg/logger"
	"github.com/muxi-Infra/auditor-Backend/repository/dao"
	"github.com/muxi-Infra/auditor-Backend/repository/model"
)

type ItemService struct {
	userDAO         *dao.UserDAO
	redisJwtHandler *jwt.RedisJWTHandler
	logger          logger.Logger
}
type Data struct {
	Id     uint
	Status string
	Msg    string
}

var M = map[model.ItemStatus]string{
	model.Pending: "未审核",
	model.Pass:    "通过",
	model.Reject:  "不通过",
}

func NewItemService(userDAO *dao.UserDAO, redisJwtHandler *jwt.RedisJWTHandler, lo logger.Logger) *ItemService {
	return &ItemService{userDAO: userDAO, redisJwtHandler: redisJwtHandler, logger: lo}
}
func (s *ItemService) Select(ctx context.Context, req request.SelectReq) ([]model.Item, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	p := req.PageSize * (req.Page - 1)

	items, err := s.userDAO.Select(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(items) > p {
		if len(items) > p+req.PageSize {
			return items[p : p+req.PageSize], nil
		} else {
			return items[p:], nil
		}
	}
	return nil, nil
}
func (s *ItemService) Audit(ctx context.Context, req request.AuditReq, id uint) (request.WebHookData, model.Item, error) {

	err := s.userDAO.AuditItem(ctx, req.ItemId, req.Status, req.Reason, id)

	if err != nil {
		return request.WebHookData{}, model.Item{}, err
	}
	item, err := s.userDAO.SelectItemById(ctx, req.ItemId)
	if err != nil {
		return request.WebHookData{}, model.Item{}, err
	}
	reqBody := request.WebHookData{
		Id:     item.HookId,
		Status: M[item.Status],
		Msg:    req.Reason,
	}

	return reqBody, item, nil
}

func (s *ItemService) Hook(reqbody request.WebHookData, item model.Item) error {
	try := os.Getenv("HOOK_TRY_MAX")
	num, err := strconv.Atoi(try) // 将 string 转成 int
	if err != nil {
		return errors.New("回调次数环境变量需要为整数")
	}
	if num > 10 {
		return errors.New("too many hooks")
	}
	var req = request.HookPayload{
		Event: "audit result back",
		Data:  reqbody,
		Try:   num,
	}
	_, err = hookBack(item.HookUrl, req, "")
	if err != nil {
		s.logger.Error("hook back failed", logger.Error(err))
		return err
	}
	return nil
}

func (s *ItemService) RoleBack(item model.Item) {
	err := s.userDAO.RollBack(item.ID, 0, "")
	if err != nil {
		s.logger.Error("role back error", logger.Error(fmt.Errorf("回滚失败: item=%+v, 原因: %w", item, err)))
	}
}

func (s *ItemService) SearchHistory(ctx context.Context, id uint) ([]model.Item, error) {
	var items []model.Item
	err := s.userDAO.SearchHistory(ctx, &items, id)
	if err != nil {
		return []model.Item{}, err
	}
	return items, nil
}

func (s *ItemService) Upload(ctx context.Context, req request.UploadReq, key string) (uint, error) {
	claims, err := apikey.ParseAPIKey(key)
	if err != nil {
		return 0, err
	}
	unixTimestamp1 := int64(req.PublicTime)
	if unixTimestamp1 > 1e10 {
		unixTimestamp1 /= 1000
	}
	publicTime := time.Unix(unixTimestamp1, 0)

	projectID := uint(claims["sub"].(float64))
	id, err := s.userDAO.Upload(ctx, req, projectID, publicTime)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *ItemService) GetDetail(ctx context.Context, id uint) (model.Item, error) {
	item, err := s.userDAO.GetItemDetail(ctx, id)
	if err != nil {
		return model.Item{}, errors.New("获取条目失败")
	}
	return item, nil
}

// AuditMany 批量审核方法实现
func (s *ItemService) AuditMany(ctx context.Context, reqs []request.AuditReq, uid uint) []request.WebHookData {
	var (
		datas []request.WebHookData
		mu    sync.Mutex
	)

	g, ctx := errgroup.WithContext(ctx)

	for _, req := range reqs {
		re := req // 防止闭包引用错误

		g.Go(func() error {
			data, item, err := s.Audit(ctx, re, uid)

			// 把结果 append 到 datas（保护 datas 的并发写）
			if err != nil {
				mu.Lock()
				defer mu.Unlock()
				data.Id = re.ItemId

				data.Msg = err.Error()
				datas = append(datas, data)

			} else {
				go func() {
					err = s.Hook(data, item)
					if err != nil {
						s.RoleBack(item)
					}
				}()
			}
			return nil
		})
	}

	_ = g.Wait()
	return datas
}
