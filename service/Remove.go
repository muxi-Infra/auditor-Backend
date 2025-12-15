package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/muxi-Infra/auditor-Backend/api/request"
	"github.com/muxi-Infra/auditor-Backend/pkg/apikey"
	"github.com/muxi-Infra/auditor-Backend/repository/dao"
	"github.com/muxi-Infra/auditor-Backend/repository/model"
)

type RemoveService struct {
	UDB dao.UserDAOInterface
	IDB dao.ItemDaoInterface
	CDB dao.CommentDaoInterface
}

func NewRemoveService(udb *dao.UserDAO, idb *dao.ItemDao, cdb *dao.CommentDao) *RemoveService {
	return &RemoveService{UDB: udb, IDB: idb, CDB: cdb}
}

// CheckPower 鉴权，看是否是已注册应用
func (service *RemoveService) CheckPower(c context.Context, ac string) (bool, uint, error) {
	claims, err := apikey.ParseAPIKey(ac)
	if err != nil {
		return false, 0, err
	}
	projectIdFloat := claims["sub"].(float64)
	projectIdUint := uint(projectIdFloat)
	return true, projectIdUint, nil
}

func (service *RemoveService) Upload(c context.Context, req request.UploadReq, projectId uint) (uint, error) {
	now := time.Now()
	id, err := service.UDB.Upload(c, req, projectId, now)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (service *RemoveService) Update(c context.Context, req request.RemoveUpdateReq, projectId uint) (uint, error) {
	//装换成time.Time类型
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t := time.Unix(int64(req.PublicTime), 0).In(loc)

	id, err := service.IDB.UpdateItem(c, req, projectId, t)
	if err != nil {
		return id, err
	}

	err = service.CDB.UpdateComments(c, id, &req.Content.LastComment, &req.Content.NextComment)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (service *RemoveService) Delete(c context.Context, hookId uint, projectId uint) error {
	err := service.UDB.DeleteItemByHookId(c, hookId, projectId)
	if err != nil {
		return err
	}
	return nil
}

func (service *RemoveService) Get(c context.Context, itemIds []uint, projectId uint) (*model.RemoveItemsStatus, error) {
	var items = model.RemoveItemsStatus{
		Items: make([]model.RemoveItemStatus, 0, 10),
	}

	if len(itemIds) == 1 && itemIds[0] == 0 {
		re, err := service.UDB.GetItems(c, projectId)
		if err != nil {
			return nil, err
		}

		for _, v := range re {
			items.Items = append(items.Items, model.RemoveItemStatus{
				Status: v.Status,
				HookId: v.HookId,
			})
		}
		return &items, nil
	}

	var lastErr []string //记录循环中的每一次错误

	for _, id := range itemIds {
		re, err := service.UDB.GetItemByHookId(c, id)
		if err != nil {
			lastErr = append(lastErr, fmt.Errorf("item:%d:%w", id, err).Error())
			continue
		}
		items.Items = append(items.Items, model.RemoveItemStatus{
			Status: re.Status,
			HookId: id,
		})
	}
	if len(lastErr) > 0 {
		return &items, errors.New(strings.Join(lastErr, ","))
	}
	return &items, nil
}
