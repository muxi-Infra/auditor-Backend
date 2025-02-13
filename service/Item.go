package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"muxi_auditor/api/request"
	"muxi_auditor/api/response"
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

func NewItemService(userDAO *dao.UserDAO, redisJwtHandler *jwt.RedisJWTHandler) *ItemService {
	return &ItemService{userDAO: userDAO, redisJwtHandler: redisJwtHandler}
}
func (s *ItemService) Select(ctx context.Context, req request.SelectReq) ([]model.Project, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	projects, err := s.userDAO.Select(ctx, req)
	if err != nil {
		return nil, err
	}

	return projects, nil
}
func (s *ItemService) Audit(ctx context.Context, req request.AuditReq, id uint) error {
	err := s.userDAO.AuditItem(ctx, req, id)
	if err != nil {
		return err
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
func (s *ItemService) Upload(ctx context.Context, req request.UploadReq) error {
	claims, err := apikey.ParseAPIKey(req.ApiKey)
	if err != nil {
		return err
	}
	layout := "2006-01-02 15:04:05"
	publicTime, err := time.Parse(layout, req.PublicTime)
	if err != nil {
		return errors.New("时间转换出错")
	}
	projectID := uint(claims["sub"].(float64))
	err = s.userDAO.Upload(ctx, req, projectID, publicTime)
	if err != nil {
		return err
	}
	reqBody := response.Response{
		Code: 200,
		Msg:  "请求成功",
		Data: Data{
			Id:     req.Id,
			Status: "未审核",
			Msg:    "操作成功",
		},
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	reqs, err := http.NewRequest("POST", req.HookUrl, bytes.NewBuffer(body))
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
