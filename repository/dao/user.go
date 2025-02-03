package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"muxi_auditor/repository/model"
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
func (d *UserDAO) FindByID(ctx context.Context, id int) ([]model.User, error) {
	var users []model.User
	err := d.DB.WithContext(ctx).Preload("Projects").Joins("JOIN user_projects ON user_projects.user_id = users.id").
		Where("user_projects.project_id = ? AND user_projects.role = ?", id, 1).Find(&users).Error
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
			d.DB.Where("user_id = ? AND project_id = ?", user.ID, project.ID).First(&userProject)

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
func (d *UserDAO) ChangeProjectRole(ctx context.Context, userId int, projectPermit []model.ProjectPermit) error {
	var user model.User
	err := d.DB.WithContext(ctx).Preload("Projects").Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		return errors.New("未找到该用户")
	}
	var userProject model.UserProject
	for _, project := range projectPermit {
		userProject.Role = project.ProjectRole
		userProject.UserID = user.ID
		userProject.ProjectID = project.ProjectID
		d.DB.WithContext(ctx).Save(&userProject)
	}
	return nil
}
