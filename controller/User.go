package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"muxi_auditor/api/request"
	"muxi_auditor/api/response"
	"muxi_auditor/pkg/jwt"
	"muxi_auditor/repository/model"
	"muxi_auditor/service"
)

type UserController struct {
	service UserService
}
type UserService interface {
	GetUsers(ctx context.Context, id uint) ([]model.UserResponse, error)
	UpdateUserRole(ctx context.Context, userId uint, projectPermit []model.ProjectPermit, role int) error
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

// GetUsers @Summary 获取用户列表
// @Description 根据项目 ID 获取用户列表，要求角色为 2 才有权限
// @Tags User
// @Accept json
// @Produce json
// @Param GetUserReq body request.GetUserReq true "获取用户请求体"
// @Success 200 {object} response.Response{data=[]model.UserResponse} "成功获取用户列表"
// @Failure 40001 {object} response.Response "无权限"
// @Failure 400 {object} response.Response "获取失败"
// @Security ApiKeyAuth
// @Router /api/v1/user/getUsers [post]
func (uc *UserController) GetUsers(ctx *gin.Context, req request.GetUserReq) (response.Response, error) {
	if req.Role != 2 {
		return response.Response{
			Msg:  "无权限",
			Code: 40001,
			Data: nil,
		}, nil

	}
	userResponse, err := uc.service.GetUsers(ctx, req.Project_id)
	if err != nil {
		return response.Response{}, err
	}

	return response.Response{
		Msg:  "获取成功",
		Code: 200,
		Data: userResponse,
	}, nil
}

// UpdateUsers @Summary 更新用户角色
// @Description 修改指定用户的角色，根据项目权限设置角色信息
// @Tags User
// @Accept json
// @Produce json
// @Param UpdateUserRoleReq body request.UpdateUserRoleReq true "更新用户角色请求体"
// @Success 200 {object} response.Response "成功更新用户角色"
// @Failure 40001 {object} response.Response "无权限"
// @Failure 400 {object} response.Response "修改失败"
// @Security ApiKeyAuth
// @Router /api/v1/user/updateUser [post]
func (uc *UserController) UpdateUsers(ctx *gin.Context, req request.UpdateUserRoleReq, cla jwt.UserClaims) (response.Response, error) {

	if cla.UserRule != 2 {
		return response.Response{
			Msg:  "无权限",
			Code: 40001,
			Data: nil,
		}, nil
	}
	err := uc.service.UpdateUserRole(ctx, req.UserId, req.ProjectPermit, req.Role)
	if err != nil {
		return response.Response{}, err
	}
	return response.Response{
		Code: 200,
		Msg:  "修改成功",
		Data: nil,
	}, nil
}
