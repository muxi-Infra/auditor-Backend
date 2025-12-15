package dao

import (
	"context"
	"gorm.io/gorm"

	"github.com/muxi-Infra/auditor-Backend/repository/model"
)

type ProjectDAOInterface interface {
	GetProjectRole(c context.Context, pid uint) (string, error)
	CountItems(c context.Context, pid uint) (int64, error)
}

type ProjectDAO struct {
	DB *gorm.DB
}

func NewProjectDAO(db *gorm.DB) *ProjectDAO {
	return &ProjectDAO{
		DB: db,
	}
}

func (d *ProjectDAO) GetProjectRole(c context.Context, pid uint) (string, error) {
	var pro model.Project
	err := d.DB.WithContext(c).First(&pro, pid).Error
	if err != nil {
		return "", err
	}
	return pro.AuditRule, nil
}

func (d *ProjectDAO) CountItems(c context.Context, pid uint) (int64, error) {
	var count int64

	err := d.DB.WithContext(c).Model(&model.Item{}).Where("project_id = ?", pid).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, err
}
