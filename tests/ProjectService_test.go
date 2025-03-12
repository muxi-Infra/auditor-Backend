package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"muxi_auditor/api/response"
	"muxi_auditor/repository/model"
	"muxi_auditor/service"
	"muxi_auditor/tests/mocks"
	"testing"
)

func TestCreate(t *testing.T) {
	// 创建一个 mocks UserDAO
	mockUserDAO := new(mocks.UserDAOInterface)

	// 设置期望：mocks 对 FindByUserIDs 和 CreateProject 的调用
	mockUserDAO.On("FindByUserIDs", mock.Anything, []uint{1, 2, 3}).Return([]model.User{
		{Model: gorm.Model{ID: 1}, Name: "User 1"},
		{Model: gorm.Model{ID: 2}, Name: "User 2"},
		{Model: gorm.Model{ID: 3}, Name: "User 3"},
	}, nil)

	mockUserDAO.On("CreateProject", mock.Anything, mock.Anything).Return(nil).Once()
	// 创建 ProjectService 实例
	projectService := service.NewProjectService(mockUserDAO, nil)

	// 调用 Create 方法
	err := projectService.Create(context.Background(), "Project A", "logo.png", "Audio Rule", []uint{1, 2, 3})
	assert.NoError(t, err, "Expected no error during project creation")

	// 验证 CreateProject 方法是否被正确调用
	mockUserDAO.AssertExpectations(t)
}
func TestGetProjectList(t *testing.T) {
	mockUserDAO := new(mocks.UserDAOInterface)
	mockUserDAO.On("GetProjectList", mock.Anything, "logo").Return([]model.Project{
		{
			Model: gorm.Model{
				ID: 1, // 项目ID
			},
			ProjectName: "Project A",
			Logo:        "logo.png",
			AudioRule:   "Audio Rule A",
			Users: []model.User{
				{Model: gorm.Model{ID: 1}, Name: "User 1"},
				{Model: gorm.Model{ID: 2}, Name: "User 2"},
			},
			Items: []model.Item{
				{Model: gorm.Model{ID: 1}, Author: "Item 1"},
				{Model: gorm.Model{ID: 2}, Auditor: "Item 2"},
			},
			Apikey: "apikey123",
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			ProjectName: "Project B",
			Logo:        "logo2.png",
			AudioRule:   "Audio Rule B",
			Users: []model.User{
				{Model: gorm.Model{ID: 3}, Name: "User 3"},
			},
			Items: []model.Item{
				{Model: gorm.Model{ID: 3}, Author: "Item 3"},
			},
			Apikey: "apikey456",
		},
	}, nil)
	projectService := service.NewProjectService(mockUserDAO, nil)
	l, err := projectService.GetProjectList(context.Background(), "logo")
	assert.NoError(t, err, "Expected no error during project list")
	assert.Equal(t, []model.ProjectList{{ProjectId: 1, ProjectName: "Project A"}, {ProjectId: 2, ProjectName: "Project B"}}, l, "Expected project list to be equal")
	mockUserDAO.AssertExpectations(t)
}
func TestDetail(t *testing.T) {
	mockUserDAO := new(mocks.UserDAOInterface)

	// 使用 mock.Anything 来匹配所有类型
	mockUserDAO.On("GetProjectDetails", mock.Anything, mock.Anything).Return(model.Project{
		Model:       gorm.Model{ID: 1},
		ProjectName: "Project A",
		Logo:        "logo.png",
		AudioRule:   "Audio Rule A",
		Users: []model.User{
			{Model: gorm.Model{ID: 1}, Name: "User 1", Avatar: "logo.png"},
			{Model: gorm.Model{ID: 2}, Name: "User 2", Avatar: "logo.png"},
		},
		Items: []model.Item{
			{Model: gorm.Model{ID: 1}, Author: "Item 1", Status: 1},
			{Model: gorm.Model{ID: 2}, Auditor: "Item 2", Status: 0},
		},
		Apikey: "apikey123",
	}, nil)

	projectService := service.NewProjectService(mockUserDAO, nil)
	d, err := projectService.Detail(context.Background(), uint(1))
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
	}, d)

	mockUserDAO.AssertExpectations(t)
}
