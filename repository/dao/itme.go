package dao

import (
	"context"
	"gorm.io/gorm"
	"time"

	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/model"
)

const (
	Pending = iota
	pass
	reject
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
	err := d.DB.Where("status = ?", Pending).First(&item).Error
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
	err := d.DB.WithContext(ctx).Where("hook_id=?", req.Id).First(&it).Error
	if err != nil {
		return 0, err
	}
	it.Status = Pending
	it.ProjectId = id

	if req.Author != "" {
		it.Author = req.Author
	}
	if len(req.Tags) > 0 {
		it.Tags = req.Tags
	}
	if req.Content.Topic.Title != "" {
		it.Title = req.Content.Topic.Title
	}
	if req.Content.Topic.Content != "" {
		it.Content = req.Content.Topic.Content
	}
	if req.PublicTime != 0 {
		it.PublicTime = t
	}
	if len(req.Content.Topic.Pictures) > 0 {
		it.Pictures = req.Content.Topic.Pictures
	}
	if req.HookUrl != "" {
		it.HookUrl = req.HookUrl
	}

	err = d.DB.WithContext(ctx).Select("*").Updates(&it).Error
	if err != nil {
		return 0, err
	}

	return it.ID, nil
}
