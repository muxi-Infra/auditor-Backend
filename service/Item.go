package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"muxi_auditor/api/request"
	"muxi_auditor/pkg/apikey"
	"muxi_auditor/pkg/jwt"
	"muxi_auditor/repository/dao"
	"muxi_auditor/repository/model"
	"net/http"
	"time"
)

type ItemService struct {
	userDAO         *dao.UserDAO
	redisJwtHandler *jwt.RedisJWTHandler
}
type Data struct {
	Id     int
	Status string
	Msg    string
}

var M = map[int]string{
	0: "未审核",
	1: "通过",
	2: "不通过",
}

func NewItemService(userDAO *dao.UserDAO, redisJwtHandler *jwt.RedisJWTHandler) *ItemService {
	return &ItemService{userDAO: userDAO, redisJwtHandler: redisJwtHandler}
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
func (s *ItemService) Audit(ctx context.Context, req request.AuditReq, id uint) (Data, model.Item, error) {

	err := s.userDAO.AuditItem(ctx, req.ItemId, req.Status, req.Reason, id)

	if err != nil {
		return Data{}, model.Item{}, err
	}
	item, err := s.userDAO.SelectItemById(ctx, req.ItemId)
	if err != nil {
		return Data{}, model.Item{}, err
	}
	reqBody := Data{
		Id:     item.HookId,
		Status: M[item.Status],
		Msg:    req.Reason,
	}

	return reqBody, item, nil
}
func (s *ItemService) Hook(reqbody Data, item model.Item) error {
	body, err := json.Marshal(reqbody)
	if err != nil {
		return err
	}

	reqs, err := http.NewRequest("POST", item.HookUrl, bytes.NewBuffer(body))

	if err != nil {

		return err
	}
	reqs.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(reqs)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {

		return errors.New("回调HookUrl失败")
	}
	return nil
}
func (s *ItemService) RoleBack(item model.Item) error {
	err := s.userDAO.RollBack(item.ID, 0, item.Reason)
	if err != nil {
		return errors.New("回滚失败")
	}
	return nil
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
func (s *ItemService) GetTags(ctx context.Context, id uint) ([]string, error) {
	items, err := s.userDAO.GetItems(ctx, id)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	tags := make([]string, 0)
	for _, item := range items {
		for _, tag := range item.Tags {
			m[tag] = tag
		}
	}
	for tag, _ := range m {
		tags = append(tags, tag)
	}
	return tags, nil
}
