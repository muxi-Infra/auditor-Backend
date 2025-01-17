package controller

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	api_errors "muxi_auditor/api/errors"
	"muxi_auditor/api/request"
	"muxi_auditor/api/response"
	"muxi_auditor/client"
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

func (c *AuthController) Login(ctx *gin.Context, req request.LoginReq) (response.Response, error) {

	////随便写的逻辑,你需要进行更改
	//username, err := c.client.GetOAuth(req.Code)
	//if err != nil {
	//	return response.Response{}, err
	//}
	//
	//err = c.service.Login(ctx, username)
	//if err != nil {
	//	return response.Response{}, err
	//}
	return response.Response{}, api_errors.LOGIN_ERROR(errors.New("登陆失败测试"))
	//返回
	//return response.Response{
	//	Msg:  "",
	//	Code: 0,
	//	Data: nil,
	//}, nil
}
