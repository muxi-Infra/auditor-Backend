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
