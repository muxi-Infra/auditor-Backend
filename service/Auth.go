package service

import (
	"context"
	merr "errors"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
	"muxi_auditor/api/errors"
	"muxi_auditor/api/request"
	"muxi_auditor/config"
	"muxi_auditor/pkg/jwt"
	"muxi_auditor/repository/dao"
	"muxi_auditor/repository/model"
	"time"
)

func NewAuthService(userDAO *dao.UserDAO, redisJwtHandler *jwt.RedisJWTHandler, conf *config.QiNiuYunConfig) *AuthService {
	return &AuthService{userDAO: userDAO, redisJwtHandler: redisJwtHandler, conf: conf}
}

type AuthService struct {
	userDAO         *dao.UserDAO
	redisJwtHandler *jwt.RedisJWTHandler
	conf            *config.QiNiuYunConfig
}

func (s *AuthService) Login(ctx context.Context, email string) (string, string, int, error) {
	//随便写的逻辑,需要修改
	user, err := s.userDAO.FindByEmail(ctx, email)
	if err != nil {
		return "", "", 0, err
	}
	if user == nil {
		return "", "", 0, nil
	}
	token, err := s.redisJwtHandler.Jwt.SetJWTToken(user.ID, user.Name, user.UserRole)
	if err != nil {
		return "", "", 0, err
	}
	return user.Name, token, user.UserRole, nil
	//执行注册的具体逻辑

}
func (s *AuthService) Register(ctx context.Context, email string, username string) (string, error) {
	user := model.User{
		Email:    email,
		Name:     username,
		UserRole: 0,
	}
	err := s.userDAO.Create(ctx, &user)
	if err != nil {
		return "", errors.LOGIN_ERROR(err)
	}
	token, err := s.redisJwtHandler.Jwt.SetJWTToken(user.ID, user.Name, user.UserRole)
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
func (s *AuthService) UpdateMyInfo(ctx context.Context, req request.UpdateUserReq, id uint) error {
	if req.Email == "" || req.Name == "" {
		return merr.New("email and name are required")
	}
	// Check if user exists
	existingUser, err := s.userDAO.Read(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return merr.New("user not found")
	}
	var user model.User
	user.Email = req.Email
	user.Name = req.Name
	user.Avatar = req.Avatar
	err = s.userDAO.Update(ctx, &user, id)
	if err != nil {
		return err
	}
	return nil
}
func (s *AuthService) GetQiToken(ctx context.Context) (string, error) {
	accesskey := s.conf.AccessKey
	secretkey := s.conf.SecretKey
	bucket := s.conf.Bucket
	mac := credentials.NewCredentials(accesskey, secretkey)
	putPolicy, err := uptoken.NewPutPolicy(bucket, time.Now().Add(1*time.Hour))
	if err != nil {
		return "", err
	}
	upToken, err := uptoken.NewSigner(putPolicy, mac).GetUpToken(context.Background())
	if err != nil {
		return upToken, err
	}
	return upToken, nil
}
