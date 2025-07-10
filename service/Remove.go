package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/apikey"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/dao"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/model"
	"strings"
	"time"
)

type RemoveService struct {
	Db dao.UserDAOInterface
}

func NewRemoveService(db dao.UserDAOInterface) *RemoveService {
	return &RemoveService{Db: db}
}

// CheckPower 鉴权，看是否是已注册应用
func (service *RemoveService) CheckPower(c context.Context, ac string) (bool, uint, error) {
	claims, err := apikey.ParseAPIKey(ac)
	if err != nil {
		return false, 0, err
	}
	return true, claims["sub"].(uint), nil
}
func (service *RemoveService) Upload(c context.Context, req request.UploadReq, projectId uint) (uint, error) {
	now := time.Now()
	id, err := service.Db.Upload(c, req, projectId, now)
	if err != nil {
		return id, err
	}
	return id, nil
}
func (service *RemoveService) Update(c context.Context, req request.UploadReq, projectId uint) (uint, error) {
	//装换成time.Time类型
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t := time.Unix(int64(req.PublicTime), 0).In(loc)
	id, err := service.Db.UpdateItem(c, req, projectId, t)
	if err != nil {
		return id, err
	}
	return id, nil
}
func (service *RemoveService) Delete(c context.Context, hookId uint, projectId uint) error {

	err := service.Db.DeleteItemByHookId(c, hookId, projectId)
	if err != nil {
		return err
	}
	return nil
}
func (service *RemoveService) Get(c context.Context, itemIds []uint, projectId uint) ([]model.Item, error) {
	if len(itemIds) == 1 && itemIds[0] == 0 {
		re, err := service.Db.GetItems(c, projectId)
		if err != nil {
			return nil, err
		}
		return re, nil
	}
	var items []model.Item
	var lastErr []string //记录循环中的每一次错误
	for _, id := range itemIds {
		re, err := service.Db.GetItemByHookId(c, id)
		if err != nil {
			lastErr = append(lastErr, fmt.Errorf("item:%d:%w", id, err).Error())
			continue
		}
		items = append(items, re)
	}
	if len(lastErr) > 0 {
		return items, errors.New(strings.Join(lastErr, ","))
	}
	return items, nil
}
