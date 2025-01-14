package service

import (
	"context"
	"gorm.io/gorm"
	"muxi_auditor/api/errors"
	"muxi_auditor/repository/dao"
	"muxi_auditor/repository/model"
)

func NewAuthService(userDAO *dao.UserDAO) *AuthService {
	return &AuthService{userDAO: userDAO}
}

type AuthService struct {
	userDAO *dao.UserDAO
}

func (s *AuthService) Login(ctx context.Context, username string) error {
	//随便写的逻辑,需要修改
	err := s.userDAO.Create(ctx, &model.User{
		Name:  username,
		Model: gorm.Model{},
	})
	if err != nil {
		return errors.LOGIN_ERROR(err)
	}
	//执行注册的具体逻辑
	return nil
}
