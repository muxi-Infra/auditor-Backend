package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"muxi_auditor/api/response"
	"muxi_auditor/repository/dao"
	"muxi_auditor/repository/model"
	"muxi_auditor/service"
	"testing"
)

// 假设你有一个 Test DB 环境，配置如下：
func setupDB() (*gorm.DB, error) {
	// 使用真实数据库连接
	dsn := "root:mysecretpassword@tcp(mysql:3306)/my_test_db?charset=utf8&parseTime=True&loc=Local"
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
func initDB(db *gorm.DB) {
	// 自动迁移：确保数据库表与模型结构一致
	if err := db.AutoMigrate(&model.User{}, &model.Project{}, &model.Item{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 插入用户数据
	users := []model.User{
		{Model: gorm.Model{ID: 1}, Name: "User 1", Avatar: "logo.png"},
		{Model: gorm.Model{ID: 2}, Name: "User 2", Avatar: "logo.png"},
		{Model: gorm.Model{ID: 3}, Name: "User 3", Avatar: "logo.png"},
	}
	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			log.Fatalf("Failed to insert user data: %v", err)
		}
	}

	// 插入项目数据
	projects := []model.Project{{
		Model:       gorm.Model{ID: 1},
		ProjectName: "Project A",
		Logo:        "logo.png",
		AudioRule:   "Audio Rule A",
		Users: []model.User{
			{Model: gorm.Model{ID: 1}, Name: "User 1", Avatar: "logo.png"},
			{Model: gorm.Model{ID: 2}, Name: "User 2", Avatar: "logo.png"},
			{Model: gorm.Model{ID: 3}, Name: "User 3", Avatar: "logo.png"},
		},
		Items: []model.Item{
			{Model: gorm.Model{ID: 1}, Author: "Item 1", Status: 1},
			{Model: gorm.Model{ID: 2}, Auditor: "Item 2", Status: 0},
		},
		Apikey: "apikey123",
	},
		{
			Model:       gorm.Model{ID: 2},
			ProjectName: "Project B",
			Logo:        "logo.png",
			AudioRule:   "Audio Rule B",
			Users: []model.User{
				{Model: gorm.Model{ID: 4}, Name: "User 4", Avatar: "logo.png"},
			},
			Items: []model.Item{
				{Model: gorm.Model{ID: 3}, Author: "Item 3", Status: 1},
			},
		},
	}

	// 插入项目数据
	for _, project := range projects {
		if err := db.Create(&project).Error; err != nil {
			log.Fatalf("Failed to insert project data: %v", err)
		}
	}

}
func TestJcCreate(t *testing.T) {
	db, err := setupDB()
	if err != nil {
		t.Error(err)
	}
	initDB(db)
	d := dao.NewUserDAO(db)
	projectService := service.NewProjectService(d, nil)
	err = projectService.Create(context.Background(), "Project C", "logo.png", "Audio Rule", []uint{1, 2, 3})
	assert.NoError(t, err, "Expected no error during project creation")

}
func TestJcGetProjectList(t *testing.T) {
	db, err := setupDB()
	if err != nil {
		t.Error(err)
	}
	initDB(db)
	d := dao.NewUserDAO(db)
	projectService := service.NewProjectService(d, nil)
	l, err := projectService.GetProjectList(context.Background(), "logo")
	assert.NoError(t, err, "Expected no error during project list")
	assert.Equal(t, []model.ProjectList{{ProjectId: 1, Name: "Project A"}, {ProjectId: 2, Name: "Project B"}}, l)

}
func TestJcDetail(t *testing.T) {
	db, err := setupDB()
	if err != nil {
		t.Error(err)
	}
	initDB(db)
	d := dao.NewUserDAO(db)
	projectService := service.NewProjectService(d, nil)
	dt, err := projectService.Detail(context.Background(), uint(1))
	assert.NoError(t, err, "Expected no error during project detail")
	assert.Equal(t, response.GetDetailResp{
		TotleNumber:   2,
		CurrentNumber: 1,
		Apikey:        "apikey123",
		AuditRule:     "Audio Rule A",
		Members: []model.UserResponse{
			{UserID: 1, Name: "User 1", Avatar: "logo.png"},
			{UserID: 2, Name: "User 2", Avatar: "logo.png"},
		},
	}, dt)
}
