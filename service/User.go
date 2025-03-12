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

func (userService *UserService) GetUsers(ctx context.Context, id uint) ([]model.UserResponse, error) {
	users, err := userService.userDAO.FindByProjectID(ctx, id)
	if err != nil {
		return nil, err
	}
	userResponse, err := userService.userDAO.GetResponse(ctx, users)
	if err != nil {
		return nil, err
	}
	return userResponse, nil
}
func (userService *UserService) UpdateUserRole(ctx context.Context, userId uint, projectPermit []model.ProjectPermit, role int) error {
	user, err := userService.userDAO.PPFUserByid(ctx, userId)
	if err != nil {
		return err
	}
	user.UserRole = role
	err = userService.userDAO.Update(ctx, &user, userId)
	if err != nil {
		return err
	}

	err = userService.userDAO.ChangeProjectRole(ctx, user, projectPermit)
	if err != nil {
		return err
	}
	return nil
}
