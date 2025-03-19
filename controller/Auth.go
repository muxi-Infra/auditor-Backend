package controller

import (
	"context"
	"fmt"
	"muxi_auditor/config"
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
	qi      *config.QiNiuYunConfig
}

type AuthService interface {
	Login(ctx context.Context, email string) (string, int, error)
	Logout(ctx *gin.Context) error
}

func NewOAuthController(client *client.OAuthClient, service *service.AuthService) *AuthController {
	return &AuthController{
		client:  client,
		service: service,
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 通过邮箱登录，返回用户的 Token
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body request.LoginReq true "登录请求体"
// @Success 200 {object} response.Response{data=string} "成功返回Token"
// @Success 20001 {object} response.Response{data=string} "审核中"
// @Failure 400 {object} response.Response{data=nil} "错误信息"
// @Router /api/v1/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context, req request.LoginReq) (response.Response, error) {

	////随便写的逻辑,你需要进行更改
	accessToken, err := c.client.GetOAuth(req.Code)
	if err != nil {
		return response.Response{
			Msg:  "获取accessToken失败",
			Data: err.Error(),
		}, err
	}
	email, err := c.client.GetEmail(accessToken)
	if err != nil || email == "" {
		return response.Response{
			Msg: "获取email失败",
		}, err
	}
	token, role, err := c.service.Login(ctx, email)

	if err != nil {
		return response.Response{
			Msg:  "login服务出错",
			Data: err.Error(),
		}, err
	}
	if role == 0 {
		return response.Response{
			Msg:  "审核中",
			Code: 20001,
			Data: "",
		}, nil
	}

	return response.Response{
		Msg:  "",
		Code: 200,
		Data: token,
	}, nil
	//返回
	//return response.Response{
	//	Msg:  "",
	//	Code: 0,
	//	Data: nil,
	//}, nil
}

// Logout 用户登出
// @Summary 用户登出
// @Description 清除用户 Token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "成功登出"
// @Failure 400 {object} response.Response "错误信息"
// @Security ApiKeyAuth
// @Router /api/v1/auth/logout [get]
func (c *AuthController) Logout(ctx *gin.Context) (response.Response, error) {
	err := c.service.Logout(ctx)
	fmt.Println(err)
	if err != nil {
		return response.Response{
			Msg:  "",
			Code: 400,
			Data: err.Error(),
		}, err
	}
	fmt.Println(2)
	return response.Response{
		Msg:  "成功登出",
		Code: 200,
		Data: nil,
	}, nil
}
