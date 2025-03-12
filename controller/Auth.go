package controller

import (
	"context"
	"muxi_auditor/config"
	"muxi_auditor/pkg/ginx"
	"muxi_auditor/pkg/jwt"
	"muxi_auditor/repository/model"

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
	UpdateMyInfo(ctx context.Context, req request.UpdateUserReq, id uint) error
	GetQiToken(ctx context.Context) (string, error)
	GetMyInfo(ctx context.Context, id uint) (*model.User, error)
}

func NewOAuthController(client *client.OAuthClient, service *service.AuthService) *AuthController {
	return &AuthController{
		client:  client,
		service: service,
	}
}

// Login @Summary 用户登录
// @Description 通过邮箱登录，返回用户的 Token
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body request.LoginReq true "登录请求体"
// @Success 200 {object} response.Response{data=string} "成功返回Token"
// @Success 20001 {object} response.Response{data=string} "审核中"
// @Failure 400 {object} response.Response{data=nil} "错误信息"
// @Router /api/v1/user/login [post]
func (c *AuthController) Login(ctx *gin.Context, req request.LoginReq) (response.Response, error) {

	////随便写的逻辑,你需要进行更改
	accessToken, err := c.client.GetOAuth(req.Code)
	if err != nil {
		return response.Response{}, err
	}
	email, err := c.client.GetEmail(accessToken)
	if err != nil || email == "" {
		return response.Response{}, err
	}
	token, role, err := c.service.Login(ctx, email)

	if err != nil {
		return response.Response{}, err
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

// Logout @Summary 用户登出
// @Description 清除用户 Token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "成功登出"
// @Failure 400 {object} response.Response "错误信息"
// @Security ApiKeyAuth
// @Router /api/v1/user/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) (response.Response, error) {
	_, err := ginx.GetClaims[jwt.UserClaims](ctx)
	if err != nil {
		return response.Response{
			Msg:  "",
			Code: 400,
			Data: nil,
		}, err
	}
	err = c.service.Logout(ctx)
	if err != nil {
		return response.Response{
			Msg:  "",
			Code: 400,
			Data: nil,
		}, err
	}
	return response.Response{
		Msg:  "成功登出",
		Code: 200,
		Data: nil,
	}, nil
}

// UpdateMyInfo @Summary 更新用户信息
// @Description 更新当前用户的信息，如邮箱、名称和头像
// @Tags Auth
// @Accept json
// @Produce json
// @Param update body request.UpdateUserReq true "更新用户信息请求体"
// @Success 200 {object} response.Response "成功更新用户信息"
// @Failure 400 {object} response.Response "Invalid or expired token"
// @Security ApiKeyAuth
// @Router /api/v1/user/updateMyInfo [post]
func (c *AuthController) UpdateMyInfo(ctx *gin.Context, req request.UpdateUserReq) (response.Response, error) {
	token, err := ginx.GetClaims[jwt.UserClaims](ctx)
	if err != nil {
		return response.Response{
			Msg:  "Invalid or expired token",
			Code: 400,
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

// GetQiToken @Summary 获取七牛云上传 Token
// @Description 获取用于上传文件的七牛云 Token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "返回上传 Token"
// @Failure 400 {object} response.Response "获取图床token失败"
// @Security ApiKeyAuth
// @Router /api/v1/user/GetQiToken [get]
func (c *AuthController) GetQiToken(ctx *gin.Context) (response.Response, error) {
	token, err := c.service.GetQiToken(ctx)
	if err != nil {
		return response.Response{
			Code: 400,
			Data: nil,
			Msg:  "获取图床token失败",
		}, err
	}
	return response.Response{
		Code: 200,
		Data: token,
		Msg:  "获取成功",
	}, nil
}

// GetMyInfo @Summary 获取自己信息
// @Description 获取用户名，邮箱，权限
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=response.UserInfo} "获取信息成功"
// @Failure 400 {object} response.Response{data=nil} "失败"
// @Security ApiKeyAuth
// @Router /api/v1/user/GetMyInfo [get]
func (c *AuthController) GetMyInfo(ctx *gin.Context, cla *jwt.UserClaims) (response.Response, error) {
	user, err := c.service.GetMyInfo(ctx, cla.Uid)
	if err != nil {
		return response.Response{

			Code: 400,
			Data: nil,
			Msg:  "获取失败",
		}, err
	}

	return response.Response{
		Code: 200,
		Data: response.UserInfo{
			Name:   user.Name,
			Role:   user.UserRole,
			Email:  user.Email,
			Avatar: user.Avatar,
		},
		Msg: "获取成功",
	}, nil
}
