package dao

import (
	"context"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/model"
	"gorm.io/gorm"
)

const (
	NotAudit = iota
	reject   = iota
	pass     = iota
)

type ItemDaoInterface interface {
	AuditItem(id uint, status int, reason string) error
	GetOneItem() (*model.Item, error)
	FindItemByID(ctx context.Context, id uint) (*model.Item, error)
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
	err := d.DB.Where("status = ?", NotAudit).First(&item).Error
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
