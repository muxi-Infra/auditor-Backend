package dao

import (
	"context"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/model"
	"gorm.io/gorm"
)

type ProjectDAOInterface interface {
	GetProjectRole(c context.Context, pid uint) (string, error)
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
