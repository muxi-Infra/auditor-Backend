package dao

import (
	"gorm.io/gorm"
	"muxi_auditor/repository/model"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&model.User{}, &model.Project{}, &model.UserProject{}, &model.Item{}, &model.Comment{})
}
