package service

import (
	"context"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
	"muxi_auditor/config"
	"muxi_auditor/pkg/jwt"
	"muxi_auditor/repository/dao"
	"time"
)

type TubeService struct {
	userDAO         *dao.UserDAO
	redisJwtHandler *jwt.RedisJWTHandler
	conf            *config.QiNiuYunConfig
}

func NewTubeService(userDAO *dao.UserDAO, redisJwtHandler *jwt.RedisJWTHandler, conf *config.QiNiuYunConfig) *TubeService {
	return &TubeService{userDAO: userDAO, redisJwtHandler: redisJwtHandler, conf: conf}
}

func (s *TubeService) GetQiToken(ctx context.Context) (string, error) {
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
