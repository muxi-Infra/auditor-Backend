package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"

	"github.com/muxi-Infra/auditor-Backend/api/request"
	"github.com/muxi-Infra/auditor-Backend/api/response"
	"github.com/muxi-Infra/auditor-Backend/pkg/jwt"
	"github.com/muxi-Infra/auditor-Backend/repository/model"
	"github.com/muxi-Infra/auditor-Backend/service"
)

type UserController struct {
	service UserService
}
type UserService interface {
	UpdateUserRole(ctx context.Context, userId uint, projectPermit []model.ProjectPermit, role int) error
	GetMyInfo(ctx context.Context, id uint) (*model.User, error)
	UpdateMyInfo(ctx context.Context, req request.UpdateUserReq, id uint) error
	GetUsers(ctx context.Context, req request.GetUsers) ([]response.UserAllInfo, error)
	GetUserInfo(ctx context.Context, id uint) (*model.User, error)
	UpdateUsersRole(ctx context.Context, li []request.UserRole) error
	NoPermissionUsers(ctx context.Context) ([]response.UserInfo, error)
	GetRoleInProject(ctx context.Context, projectId uint, uid uint) (int, error)
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

// UpdateUsers 更新用户权限等信息
// @Summary 更新用户角色
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
func (c *UserController) UpdateUsers(ctx *gin.Context, req request.UpdateUserRoleReq, cla jwt.UserClaims) (response.Response, error) {

	if cla.UserRule != 2 {
		return response.Response{
			Msg:  "无权限",
			Code: 40001,
			Data: nil,
		}, nil
	}
	fmt.Println(1)
	err := c.service.UpdateUserRole(ctx, req.UserId, req.ProjectPermit, req.Role)
	if err != nil {
		return response.Response{}, err
	}
	return response.Response{
		Code: 200,
		Msg:  "修改成功",
		Data: nil,
	}, nil
}

// GetMyInfo 获取自己信息
// @Summary 获取自己信息
// @Description 获取用户名，邮箱，权限
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=response.UserInfo} "获取信息成功"
// @Failure 400 {object} response.Response{data=nil} "失败"
// @Security ApiKeyAuth
// @Router /api/v1/user/getMyInfo [get]
func (c *UserController) GetMyInfo(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error) {

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
			Id:     cla.Uid,
			Name:   user.Name,
			Role:   user.UserRole,
			Email:  user.Email,
			Avatar: user.Avatar,
		},
		Msg: "获取成功",
	}, nil
}

// UpdateMyInfo 更新自己信息
// @Summary 更新用户信息
// @Description 更新当前用户的信息，如邮箱、名称和头像
// @Tags User
// @Accept json
// @Produce json
// @Param update body request.UpdateUserReq true "更新用户信息请求体"
// @Success 200 {object} response.Response "成功更新用户信息"
// @Failure 400 {object} response.Response "Invalid or expired token"
// @Security ApiKeyAuth
// @Router /api/v1/user/updateMyInfo [post]
func (c *UserController) UpdateMyInfo(ctx *gin.Context, req request.UpdateUserReq, cla jwt.UserClaims) (response.Response, error) {

	err := c.service.UpdateMyInfo(ctx, req, cla.Uid)
	if err != nil {
		return response.Response{}, err
	}
	return response.Response{
		Msg:  "更新用户信息成功",
		Code: 200,
		Data: nil,
	}, nil
}

// SelectUsers 获取或查询所有用户信息
// @Summary 获取或查询所有用户信息
// @Description 获取或查询所有用户信息包括权限等
// @Tags User
// @Accept json
// @Produce json
// @Param the_query query string false "查询关键字"
// @Param page query int false "页码 (默认: 1)"
// @Param pageSize query int false "每页数量 (默认: 10)"
// @Success 200 {object} response.Response "成功获取用户信息"
// @Failure 400 {object} response.Response "Invalid or expired token"
// @Security ApiKeyAuth
// @Router /api/v1/user/getUsers [get]
func (c *UserController) SelectUsers(ctx *gin.Context) (response.Response, error) {
	query1 := ctx.DefaultQuery("the_query", "")
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if err != nil {
		pageSize = 10
	}
	var req request.GetUsers
	req.Query = query1
	req.Page = page
	req.PageSize = pageSize
	re, err := c.service.GetUsers(ctx, req)
	if err != nil {
		return response.Response{
			Code: 400,
			Data: nil,
		}, err
	}
	return response.Response{
		Code: 200,
		Data: re,
		Msg:  "success",
	}, err
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 通过用户 ID 获取详细信息
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "用户 ID"
// @Success 200 {object} response.Response{data=response.UserInfo} "获取成功"
// @Failure 400 {object} response.Response "请求错误"
// @Router /api/v1/user/{id}/getUserInfo [get]
func (c *UserController) GetUserInfo(ctx *gin.Context) (response.Response, error) {
	ID := ctx.Param("id")
	if ID == "" {
		return response.Response{
			Code: 400,
			Msg:  "需要user_id",
		}, nil
	}
	pid, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		return response.Response{
			Code: 400,
			Msg:  "获取user_id失败",
		}, err
	}
	u := uint(pid)
	user, err := c.service.GetUserInfo(ctx, u)
	if err != nil {
		return response.Response{
			Code: 400,
			Data: err.Error(),
			Msg:  "获取信息失败",
		}, err
	}
	return response.Response{
		Code: 200,
		Data: response.UserInfo{
			Id:     user.ID,
			Name:   user.Name,
			Role:   user.UserRole,
			Email:  user.Email,
			Avatar: user.Avatar,
		},
		Msg: "获取成功",
	}, nil
}

// GetNoPermissionUsers 获取待审核成员
// @Summary 获取待审核成员
// @Description 获取已注册但待审核成员
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]response.UserInfo} "获取成功"
// @Failure 400 {object} response.Response "请求错误"
// @Router /api/v1/user/getNoPermissionUsers [get]
func (c *UserController) GetNoPermissionUsers(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error) {
	if cla.UserRule != 2 {
		return response.Response{
			Msg:  "no power",
			Code: 403,
			Data: nil,
		}, nil
	}
	r, err := c.service.NoPermissionUsers(ctx)
	if err != nil {
		return response.Response{
			Msg:  "获取失败",
			Code: 400,
			Data: nil,
		}, err
	}
	return response.Response{
		Msg:  "获取待审核成员成功",
		Code: 200,
		Data: r,
	}, nil
}

// ChangeUsersRole 批量更新权限
// @Summary 批量更新权限
// @Description 批量更新用户在审核平台的权限
// @Tags User
// @Accept json
// @Produce json
// @Param update body request.ChangeUserRolesReq true "更新用户在审核平台的权限请求体"
// @Success 200 {object} response.Response "成功更新用户信息"
// @Failure 400 {object} response.Response "Invalid or expired token"
// @Failure 403 {object} response.Response "no power"
// @Security ApiKeyAuth
// @Router /api/v1/user/changeRoles [post]
func (c *UserController) ChangeUsersRole(ctx *gin.Context, req request.ChangeUserRolesReq, cla jwt.UserClaims) (response.Response, error) {
	if cla.UserRule != 2 {
		return response.Response{
			Code: 403,
			Data: nil,
			Msg:  "no power",
		}, nil
	}
	err := c.service.UpdateUsersRole(ctx, req.List)
	if err != nil {
		return response.Response{
			Code: 400,
			Data: nil,
			Msg:  "更新用户权限失败",
		}, err
	}
	return response.Response{
		Code: 200,
		Data: nil,
		Msg:  "success",
	}, err
}

// GetProjectRole 获取projectRole
// @Summary 获取projectRole
// @Description 获取用户在某个project的projectRole
// @Tags User
// @Accept json
// @Produce json
// @Param project_id path int true "project ID"
// @Success 200 {object} response.Response{data=int} "获取成功"
// @Failure 400 {object} response.Response "请求错误"
// @Router /api/v1/user/getProjectRole/{project_id} [get]
func (c *UserController) GetProjectRole(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error) {
	pid := ctx.Param("project_id")
	if pid == "" {
		return response.Response{
			Msg:  "project_id is empty",
			Code: 400,
			Data: nil,
		}, errors.New("project_id is necessary")
	}
	projectID, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		return response.Response{
			Msg:  "project_id is invalid",
			Code: 400,
			Data: nil,
		}, errors.New("project_id is invalid")
	}
	r, err := c.service.GetRoleInProject(ctx, uint(projectID), cla.Uid)
	if err != nil {
		return response.Response{
			Msg:  "获取project_role失败",
			Code: 400,
			Data: nil,
		}, err
	}
	return response.Response{
		Msg:  "获取成功",
		Code: 200,
		Data: r,
	}, nil
}
