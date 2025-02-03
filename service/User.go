package service

import (
	"context"
	"muxi_auditor/pkg/jwt"
	"muxi_auditor/repository/dao"
	"muxi_auditor/repository/model"
)

type UserService struct {
	userDAO         *dao.UserDAO
	redisJwtHandler *jwt.RedisJWTHandler
}

func NewUserService(userDAO *dao.UserDAO, redisJwtHandler *jwt.RedisJWTHandler) *UserService {
	return &UserService{userDAO: userDAO, redisJwtHandler: redisJwtHandler}
}

func (userService *UserService) GetUsers(ctx context.Context, id int) ([]model.UserResponse, error) {
	users, err := userService.userDAO.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	userResponse, err := userService.userDAO.GetResponse(ctx, users)
	if err != nil {
		return nil, err
	}
	return userResponse, nil
}
func (userService *UserService) UpdateUserRole(ctx context.Context, userId int, projectPermit []model.ProjectPermit) error {
	err := userService.userDAO.ChangeProjectRole(ctx, userId, projectPermit)
	if err != nil {
		return err
	}
	return nil
}
