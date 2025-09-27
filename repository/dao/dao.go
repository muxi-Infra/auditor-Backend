package dao

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/model"
	"gorm.io/gorm"
)

const Nothing = 0

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&model.User{}, &model.Project{}, &model.UserProject{}, &model.Item{}, &model.Comment{}, &model.History{})
}
