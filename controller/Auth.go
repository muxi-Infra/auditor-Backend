package controller

import (
	"context"
	"muxi_auditor/pkg/ginx"
	"muxi_auditor/pkg/jwt"
	//"errors"
	"github.com/gin-gonic/gin"
	//api_errors "muxi_auditor/api/errors"
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
	Login(ctx context.Context, email string) (string, string, error)
	Register(ctx context.Context, email string, username string) (string, error)
	Logout(ctx *gin.Context) error
	UpdateMyInfo(ctx context.Context, req request.UpdateUserReq) error
}

func NewOAuthController(client *client.OAuthClient, service *service.AuthService) *AuthController {
	return &AuthController{
		client:  client,
		service: service,
	}
}

func (c *AuthController) Login(ctx *gin.Context, req request.LoginReq) (response.Response, error) {

	////随便写的逻辑,你需要进行更改
	accessToken, err := c.client.GetOAuth(req.Code)
	if err != nil {
		return response.Response{}, err
	}
	email, err := c.client.GetEmail(accessToken)
	if err != nil {
		return response.Response{}, err
	}
	role, token, err := c.service.Login(ctx, email)
	if err != nil {
		return response.Response{}, err
	}
	if role == "0" {
		return response.Response{
			Msg:  "",
			Code: 20001,
			Data: map[string]interface{}{
				"token":    "",
				"username": "",
				"role":     0,
				"email":    email,
			},
		}, nil
	}

	return response.Response{
		Msg:  "",
		Code: 200,
		Data: map[string]interface{}{
			"token":    token,
			"username": role,
			"role":     1,
		},
	}, nil
	//返回
	//return response.Response{
	//	Msg:  "",
	//	Code: 0,
	//	Data: nil,
	//}, nil
}
func (c *AuthController) Register(ctx *gin.Context, req request.RegisterReq) (response.Response, error) {

	token, err := c.service.Register(ctx, req.Email, req.Name)
	if err != nil {
		return response.Response{}, err
	}
	return response.Response{
		Msg:  "",
		Code: 200,
		Data: map[string]interface{}{
			"token":    token,
			"username": req.Name,
			"role":     1,
		},
	}, nil
}
func (c *AuthController) Logout(ctx *gin.Context) (response.Response, error) {
	_, err := ginx.GetClaims[jwt.UserClaims](ctx)
	if err != nil {
		return response.Response{
			Msg:  "",
			Code: 40001,
			Data: nil,
		}, err
	}
	err = c.service.Logout(ctx)
	if err != nil {
		return response.Response{
			Msg:  "",
			Code: 40001,
			Data: nil,
		}, err
	}
	return response.Response{
		Msg:  "成功登出",
		Code: 200,
		Data: nil,
	}, nil
}
func (c *AuthController) UpdateMyInfo(ctx *gin.Context, req request.UpdateUserReq) (response.Response, error) {
	err := c.service.UpdateMyInfo(ctx, req)
	if err != nil {
		return response.Response{}, err
	}
	return response.Response{
		Msg:  "更新用户信息成功",
		Code: 200,
		Data: nil,
	}, nil
}
