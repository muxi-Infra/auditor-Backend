package dao

import (
	"gorm.io/gorm"

	"github.com/muxi-Infra/auditor-Backend/repository/model"
)

const Nothing = 0

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&model.User{}, &model.Project{}, &model.UserProject{}, &model.Item{}, &model.Comment{}, &model.History{})
}
