package dao

import (
	"context"
	"encoding/json"
	"gorm.io/gorm"
	"time"

	"github.com/muxi-Infra/auditor-Backend/api/request"
	"github.com/muxi-Infra/auditor-Backend/repository/model"
)

type ItemDaoInterface interface {
	AuditItem(id uint, status int, reason string) error
	GetOneItem() (*model.Item, error)
	FindItemByID(ctx context.Context, id uint) (*model.Item, error)
	UpdateItem(ctx context.Context, req request.RemoveUpdateReq, id uint, time time.Time) (uint, error)
}

type ItemDao struct {
	DB *gorm.DB
}

func NewItemDao(db *gorm.DB) *ItemDao {
	return &ItemDao{DB: db}
}

func (d *ItemDao) AuditItem(id uint, status int, reason string) error {
	return d.DB.Model(&model.Item{}).
		Where("id = ?", id).
		UpdateColumns(map[string]interface{}{
			"status": status,
			"reason": reason,
		}).Error
}

func (d *ItemDao) GetOneItem() (*model.Item, error) {
	var item model.Item
	err := d.DB.Where("status = ?", model.Pending).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (d *ItemDao) FindItemByID(ctx context.Context, id uint) (*model.Item, error) {
	var item model.Item
	err := d.DB.WithContext(ctx).Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (d *ItemDao) UpdateItem(ctx context.Context, req request.RemoveUpdateReq, id uint, t time.Time) (uint, error) {
	var it model.Item
	err := d.DB.WithContext(ctx).Model(&model.Item{}).Where("hook_id=?", req.Id).First(&it).Error
	if err != nil {
		return 0, err
	}

	updates := map[string]interface{}{}

	updates["status"] = model.Pending
	updates["project_id"] = id

	if req.Author != "" {
		updates["author"] = req.Author
	}
	if len(req.Tags) > 0 {
		updates["tags"] = req.Tags
	}
	if req.Content.Topic.Title != "" {
		updates["title"] = req.Content.Topic.Title
	}
	if req.Content.Topic.Content != "" {
		updates["content"] = req.Content.Topic.Content
	}
	if req.PublicTime != 0 {
		updates["public_time"] = t
	}
	if len(req.Content.Topic.Pictures) > 0 {
		b, err := json.Marshal(req.Content.Topic.Pictures)
		if err != nil {
			return 0, err
		}

		updates["pictures"] = b
	}
	if req.HookUrl != "" {
		updates["hook_url"] = req.HookUrl
	}

	err = d.DB.WithContext(ctx).
		Model(&model.Item{}).
		Where("id = ?", it.ID).
		Updates(updates).Error
	if err != nil {
		return 0, err
	}

	return it.ID, nil
}
