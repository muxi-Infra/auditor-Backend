package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"muxi_auditor/api/errors"
	"muxi_auditor/api/request"
	"muxi_auditor/pkg/jwt"
	"muxi_auditor/repository/dao"
	"muxi_auditor/repository/model"
)

func NewAuthService(userDAO *dao.UserDAO, redisJwtHandler *jwt.RedisJWTHandler) *AuthService {
	return &AuthService{userDAO: userDAO, redisJwtHandler: redisJwtHandler}
}

type AuthService struct {
	userDAO         *dao.UserDAO
	redisJwtHandler *jwt.RedisJWTHandler
}

func (s *AuthService) Login(ctx context.Context, email string) (string, string, error) {
	//随便写的逻辑,需要修改
	user, err := s.userDAO.FindByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "0", "", nil
	}
	token, err := s.redisJwtHandler.Jwt.SetJWTToken(user.ID, user.Name)
	if err != nil {
		return "", "", err
	}
	return user.Name, token, nil
	//执行注册的具体逻辑

}
func (s *AuthService) Register(ctx context.Context, email string, username string) (string, error) {
	user := model.User{
		Email: email,
		Name:  username,
	}
	err := s.userDAO.Create(ctx, &user)
	if err != nil {
		return "", errors.LOGIN_ERROR(err)
	}
	token, err := s.redisJwtHandler.Jwt.SetJWTToken(user.ID, user.Name)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) Logout(ctx *gin.Context) error {
	err := s.redisJwtHandler.ClearToken(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (s *AuthService) UpdateMyInfo(ctx context.Context, req request.UpdateUserReq) error {
	var user model.User
	user.Email = req.Email
	user.Name = req.Name
	err := s.userDAO.Update(ctx, &user)
	if err != nil {
		return err
	}
	return nil
}
