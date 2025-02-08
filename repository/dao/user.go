package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"muxi_auditor/api/response"
	"muxi_auditor/pkg/apikey"
	"muxi_auditor/repository/model"
	"strconv"
)

type UserDAO struct {
	DB *gorm.DB
}

// NewUserDAO 创建一个新的 UserDAO
func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db}
}

func (d *UserDAO) Create(ctx context.Context, user *model.User) error {
	if err := d.DB.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (d *UserDAO) Read(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := d.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// 预计用不上
func (d *UserDAO) Update(ctx context.Context, user *model.User) error {
	if err := d.DB.WithContext(ctx).Save(user).Error; err != nil {
		return err
	}
	return nil
}

// 预计用不上
func (d *UserDAO) Delete(ctx context.Context, id uint) error {
	if err := d.DB.WithContext(ctx).Delete(&model.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (d *UserDAO) List(ctx context.Context) ([]model.User, error) {
	var users []model.User
	if err := d.DB.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
func (d *UserDAO) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := d.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func (d *UserDAO) FindByProjectID(ctx context.Context, id int) ([]model.User, error) {
	var users []model.User
	err := d.DB.WithContext(ctx).Preload("Projects").Joins("JOIN user_projects ON user_projects.user_id = users.id").
		Where("user_projects.project_id = ? AND user_projects.role = ?", id, 1).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (d *UserDAO) FindByUserIDs(ctx context.Context, ids []uint) ([]model.User, error) {
	var users []model.User
	err := d.DB.WithContext(ctx).Where("id IN ?", ids).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (d *UserDAO) GetResponse(ctx context.Context, users []model.User) ([]model.UserResponse, error) {
	var userResponses []model.UserResponse
	for _, user := range users {
		var projectPermits []model.ProjectPermit
		for _, project := range user.Projects {
			var userProject model.UserProject
			d.DB.WithContext(ctx).Where("user_id = ? AND project_id = ?", user.ID, project.ID).First(&userProject)

			projectPermits = append(projectPermits, model.ProjectPermit{
				ProjectID:   project.ID,
				ProjectName: project.ProjectName,
				ProjectRole: userProject.Role,
			})
		}

		userResponses = append(userResponses, model.UserResponse{
			Name:          user.Name,
			UserID:        user.ID,
			Avatar:        user.Avatar,
			ProjectPermit: projectPermits,
			Role:          user.UserRole,
		})
	}

	return userResponses, nil
}
func (d *UserDAO) ChangeProjectRole(ctx context.Context, userId int, projectPermit []model.ProjectPermit, role int) error {
	var user model.User
	err := d.DB.WithContext(ctx).Preload("Projects").Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		return errors.New("未找到该用户")
	}
	user.UserRole = role
	var userProject model.UserProject
	for _, project := range projectPermit {
		userProject.Role = project.ProjectRole
		userProject.UserID = user.ID
		userProject.ProjectID = project.ProjectID
		d.DB.WithContext(ctx).Save(&userProject)
	}
	return nil
}
func (d *UserDAO) GetProjectList(ctx context.Context) ([]model.ProjectList, error) {
	var projects []model.Project
	if err := d.DB.WithContext(ctx).Find(&projects).Error; err != nil {
		return nil, errors.New("查询数据库错误")
	}
	var projectlist []model.ProjectList
	for _, project := range projects {
		projectlist = append(projectlist, model.ProjectList{
			ProjectId:   project.ID,
			ProjectName: project.ProjectName,
		})
	}
	return projectlist, nil
}
func (d *UserDAO) CreateProject(ctx context.Context, project *model.Project) error {
	if err := d.DB.WithContext(ctx).Create(project).Error; err != nil {
		return errors.New("创建项目失败")
	}
	key, err := apikey.GenerateAPIKey(strconv.Itoa(int(project.ID)))
	if err != nil {
		return errors.New("生成apikey失败")
	}
	project.Apikey = key
	if err := d.DB.WithContext(ctx).Save(project).Error; err != nil {
		return err
	}
	return nil

}
func (d *UserDAO) GetProjectDetails(ctx context.Context, id uint) (response.GetDetailResp, error) {
	var project model.Project
	err := d.DB.WithContext(ctx).Preload("Items").Preload("Users").First(&project, id).Error
	if err != nil {
		return response.GetDetailResp{}, err
	}
	countMap := map[int]int{
		0: 0,
		1: 0,
		2: 0,
	}
	for _, item := range project.Items {
		countMap[item.Status]++
	}
	var users []model.UserResponse
	for _, user := range project.Users {
		users = append(users, model.UserResponse{
			Name:   user.Name,
			UserID: user.ID,
			Avatar: user.Avatar,
		})
	}

	re := response.GetDetailResp{
		TotleNumber:   countMap[0] + countMap[1] + countMap[2],
		CurrentNumber: countMap[0],
		Apikey:        project.Apikey,
		AuditRule:     project.AudioRule,
		Members:       users,
	}
	return re, nil
}
