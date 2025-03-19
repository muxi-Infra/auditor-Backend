package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"math/rand"
	"muxi_auditor/pkg/jwt"
	"muxi_auditor/repository/dao"
	"muxi_auditor/repository/model"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
func NewAuthService(userDAO *dao.UserDAO, redisJwtHandler *jwt.RedisJWTHandler) *AuthService {
	return &AuthService{userDAO: userDAO, redisJwtHandler: redisJwtHandler}
}

type AuthService struct {
	userDAO         *dao.UserDAO
	redisJwtHandler *jwt.RedisJWTHandler
}

func (s *AuthService) Login(ctx context.Context, email string) (string, int, error) {

	user, err := s.userDAO.FindByEmail(ctx, email)
	if err != nil {
		return "", 0, err
	}
	if user == nil {
		s.userDAO.Create(ctx, &model.User{
			Email:    email,
			Name:     RandomString(5),
			UserRole: 0,
		})
		return "", 0, nil
	}
	token, err := s.redisJwtHandler.Jwt.SetJWTToken(user.ID, user.Email, user.UserRole)
	if err != nil {
		return "", 0, err
	}
	//err = s.redisJwtHandler.CheckLogin(ctx, email)
	//if err != nil {
	//	return "", 0, err
	//}
	//err = s.redisJwtHandler.Login(ctx, email)
	//if err != nil {
	//	return "", 0, err
	//}
	return token, user.UserRole, nil

}

func (s *AuthService) Logout(ctx *gin.Context) error {
	err := s.redisJwtHandler.ClearToken(ctx)

	if err != nil {
		return err
	}

	return nil
}

//func (s *AuthService) UpdateMyInfo(ctx context.Context, req request.UpdateUserReq, id uint) error {
//	if req.Name == "" && req.Avatar == "" {
//		return merr.New(" name or avatar are required")
//	}
//	// Check if user exists
//	existingUser, err := s.userDAO.Read(ctx, id)
//	if err != nil {
//		return err
//	}
//	if existingUser == nil {
//		return merr.New("user not found")
//	}
//	var user model.User
//	if req.Name != "" {
//		user.Name = req.Name
//	}
//
//	if req.Avatar != "" {
//		user.Avatar = req.Avatar
//	}
//	err = s.userDAO.Update(ctx, &user, id)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//func (s *AuthService) GetQiToken(ctx context.Context) (string, error) {
//	accesskey := s.conf.AccessKey
//	secretkey := s.conf.SecretKey
//	bucket := s.conf.Bucket
//	mac := credentials.NewCredentials(accesskey, secretkey)
//	putPolicy, err := uptoken.NewPutPolicy(bucket, time.Now().Add(1*time.Hour))
//	if err != nil {
//		return "", err
//	}
//	upToken, err := uptoken.NewSigner(putPolicy, mac).GetUpToken(context.Background())
//	if err != nil {
//		return upToken, err
//	}
//	return upToken, nil
//}
//func (s *AuthService) GetMyInfo(ctx context.Context, id uint) (*model.User, error) {
//	user, err := s.userDAO.Read(ctx, id)
//	if err != nil {
//		return nil, err
//	}
//	if user == nil {
//		return nil, merr.New("user not found")
//	}
//	return user, nil
//}
