package controller

import (
	"context"
	"muxi_auditor/config"
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

type UserInfo struct {
	Token string `json:"token"`
	Name  string `json:"name"`
	Role  int    `json:"role"`
	Email string `json:"email"`
}
type AuthController struct {
	client  *client.OAuthClient
	service AuthService
	qi      *config.QiNiuYunConfig
}

type AuthService interface {
	Login(ctx context.Context, email string) (string, string, int, error)
	Register(ctx context.Context, email string, username string) (string, error)
	Logout(ctx *gin.Context) error
	UpdateMyInfo(ctx context.Context, req request.UpdateUserReq, id uint) error
	GetQiToken(ctx context.Context) (string, error)
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
	username, token, role, err := c.service.Login(ctx, email)
	if err != nil {
		return response.Response{}, err
	}
	if role == 0 {
		return response.Response{
			Msg:  "",
			Code: 20001,
			Data: UserInfo{
				Token: "",
				Name:  username,
				Role:  0,
				Email: email,
			},
		}, nil
	}

	return response.Response{
		Msg:  "",
		Code: 200,
		Data: UserInfo{
			Token: token,
			Name:  username,
			Role:  role,
			Email: email,
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
		Data: UserInfo{
			Token: token,
			Name:  req.Name,
			Role:  0,
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
	token, err := ginx.GetClaims[jwt.UserClaims](ctx)
	if err != nil {
		return response.Response{
			Msg:  "Invalid or expired token",
			Code: 40001,
			Data: nil,
		}, err
	}
	err = c.service.UpdateMyInfo(ctx, req, token.Uid)
	if err != nil {
		return response.Response{}, err
	}
	return response.Response{
		Msg:  "更新用户信息成功",
		Code: 200,
		Data: nil,
	}, nil
}
func (c *AuthController) GetQiToken(ctx *gin.Context) (response.Response, error) {
	token,err:=c.service.GetQiToken(ctx)
	if err != nil {
		return response.Response{
			Code: 400,
			Data: nil,
			Msg: "获取图床token失败",
		}, err
	}
	return response.Response{
		Code: 200,
		Data: token,
		Msg:"获取成功",
	},nil
}