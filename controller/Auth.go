package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"muxi_auditor/api/request"
	"muxi_auditor/client"
	"muxi_auditor/pkg/ginx"
	"muxi_auditor/service"
)

type AuthController struct {
	client  *client.OAuthClient
	service AuthService
}

type AuthService interface {
	Login(ctx context.Context, username string) error
}

func NewOAuthController(client *client.OAuthClient, service *service.AuthService) *AuthController {
	return &AuthController{
		client:  client,
		service: service,
	}
}

func (c *AuthController) Login(ctx *gin.Context, req request.LoginReq) (ginx.Response, error) {
	//随便写的逻辑,你需要进行更改
	username, err := c.client.GetOAuth(req.Code)
	if err != nil {
		return ginx.Response{}, err
	}

	err = c.service.Login(ctx, username)
	if err != nil {
		return ginx.Response{}, err
	}

	//返回
	return ginx.Response{
		Msg:  "",
		Code: 0,
		Data: nil,
	}, err
}
