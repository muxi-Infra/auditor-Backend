package controller

import (
	"context"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/ginx"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/jwt"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/model"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ProjectController struct {
	service ProjectService
}
type ProjectService interface {
	GetProjectList(ctx context.Context) ([]model.ProjectList, error)
	Create(ctx context.Context, name string, url string, logo string, audioRule string, ids []uint) (uint, error)
	Detail(ctx context.Context, id uint) (response.GetDetailResp, error)
	Delete(ctx context.Context, cla jwt.UserClaims, p uint) error
	Update(ctx context.Context, id uint, req request.UpdateProject) error
	GetUsers(ctx context.Context, id uint) ([]model.UserResponse, error)
	ReturnApiKey(apiKey string, hookUrl string) error
	GetAllTags(ctx context.Context, pid uint) ([]string, error)
}

func NewProjectController(service *service.ProjectService) *ProjectController {
	return &ProjectController{
		service: service,
	}
}

// GetProjectList 获取项目列表
// @Summary 获取项目列表
// @Description 获取所有项目列表，根据 logo 过滤
// @Tags Project
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "成功返回项目列表"
// @Failure 400 {object} response.Response "获取项目列表失败"
// @Security ApiKeyAuth
// @Router /api/v1/project/getProjectList [get]
func (ctrl *ProjectController) GetProjectList(ctx *gin.Context) (response.Response, error) {

	list, err := ctrl.service.GetProjectList(ctx)
	if err != nil {
		return response.Response{
			Code: 400,
			Msg:  "获取列表失败",
			Data: nil,
		}, err
	}
	return response.Response{
		Data: list,
		Code: 200,
		Msg:  "获取列表成功",
	}, nil
}

// Create 创建项目
// @Summary 创建项目
// @Description 根据请求体参数创建新的项目
// @Tags Project
// @Accept json
// @Produce json
// @Param createProject body request.CreateProject true "创建项目请求体"
// @Success 200 {object} response.Response "项目创建成功"
// @Failure 400 {object} response.Response "无权限或创建失败"
// @Security ApiKeyAuth
// @Router /api/v1/project/create [post]
func (ctrl *ProjectController) Create(ctx *gin.Context, req request.CreateProject) (response.Response, error) {
	token, err := ginx.GetClaims[jwt.UserClaims](ctx)
	if err != nil {
		return response.Response{
			Msg:  "",
			Code: 40001,
			Data: nil,
		}, err
	}
	if token.UserRule != 2 {
		return response.Response{
			Code: 400,
			Msg:  "无权限",
		}, nil
	}
	k, err := ctrl.service.Create(ctx, req.Name, req.HookUrl, req.Logo, req.AudioRule, req.UserIds)
	if err != nil {
		return response.Response{
			Code: 400,
			Msg:  "创建项目失败",
		}, nil
	}
	return response.Response{
		Code: 200,
		Msg:  "创建成功",
		Data: k,
	}, nil
}

// Detail 获取项目详细信息
// @Summary 获取项目详细信息
// @Description 根据项目 ID 获取项目的详细信息
// @Tags Project
// @Accept json
// @Produce json
// @Param project_id query uint true "项目ID"
// @Success 200 {object} response.Response{data=response.GetDetailResp} "获取项目详细信息成功"
// @Failure 400 {object} response.Response "获取项目详细信息失败"
// @Security ApiKeyAuth
// @Router /api/v1/project/{project_id}/detail [get]
func (ctrl *ProjectController) Detail(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error) {
	projectID := ctx.Param("project_id")
	if projectID == "" {
		return response.Response{
			Code: 400,
			Msg:  "需要project_id",
		}, nil
	}
	pid, err := strconv.ParseUint(projectID, 10, 64)
	if err != nil {
		return response.Response{
			Code: 400,
			Msg:  "获取project_id失败",
		}, err
	}
	p := uint(pid)

	detail, err := ctrl.service.Detail(ctx, p)
	if err != nil {
		return response.Response{
			Msg:  "",
			Code: 40001,
			Data: nil,
		}, err
	}
	if cla.UserRule != 2 {
		detail.Apikey = "*********"
	}
	return response.Response{
		Code: 200,
		Msg:  "获取成功",
		Data: detail,
	}, nil

}

// Delete 删除项目
// @Summary 删除项目
// @Description 通过项目 ID 删除指定的项目
// @Tags Project
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param project_id path int true "项目ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "删除失败"
// @Security ApiKeyAuth
// @Router /api/v1/project/{project_id}/delete [delete]
func (ctrl *ProjectController) Delete(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error) {

	projectID := ctx.Param("project_id")
	if projectID == "" {
		return response.Response{
			Code: 400,
			Msg:  "需要project_id",
		}, nil
	}
	pid, err := strconv.ParseUint(projectID, 10, 64)
	if err != nil {
		return response.Response{
			Code: 400,
			Msg:  "获取project_id失败",
		}, err
	}
	p := uint(pid)
	err = ctrl.service.Delete(ctx, cla, p)
	if err != nil {
		return response.Response{
			Code: 400,
			Msg:  "",
			Data: nil,
		}, err
	}
	return response.Response{
		Code: 200,
		Msg:  "删除项目成功",
		Data: nil,
	}, nil
}

// Update 更新项目信息
// @Summary 更新项目
// @Description 更新项目信息，只有用户权限为 2（管理员）时才能操作
// @Tags Project
// @Accept json
// @Produce json
// @Param project_id path int true "项目ID"
// @Param request body request.UpdateProject true "更新项目信息"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求错误（参数错误/无权限等）"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/project/{project_id}/update [post]
func (ctrl *ProjectController) Update(ctx *gin.Context, req request.UpdateProject, cla jwt.UserClaims) (response.Response, error) {
	projectID := ctx.Param("project_id")
	if projectID == "" {
		return response.Response{
			Code: 400,
			Msg:  "需要project_id",
		}, nil
	}
	pid, err := strconv.ParseUint(projectID, 10, 64)
	if err != nil {
		return response.Response{
			Code: 400,
			Msg:  "获取project_id失败",
		}, err
	}
	p := uint(pid)
	uRole := cla.UserRule
	if uRole != 2 {
		return response.Response{
			Code: 400,
			Msg:  "无权限",
			Data: nil,
		}, nil
	}
	err = ctrl.service.Update(ctx, p, req)
	if err != nil {
		return response.Response{
			Msg:  "更新失败",
			Code: 400,
			Data: nil,
		}, err
	}
	return response.Response{
		Code: 200,
		Msg:  "更新成功",
	}, nil
}

// GetUsers 获取项目成员列表
// @Summary 获取项目成员
// @Description 根据 project_id 获取该项目的用户列表
// @Tags Project
// @Accept json
// @Produce json
// @Param project_id path int true "项目ID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求错误（参数错误/无 project_id）"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/project/{project_id}/getUsers [get]
func (ctrl *ProjectController) GetUsers(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error) {
	projectID := ctx.Param("project_id")
	if projectID == "" {
		return response.Response{
			Code: 400,
			Msg:  "需要project_id",
		}, nil
	}
	pid, err := strconv.ParseUint(projectID, 10, 64)
	if err != nil {
		return response.Response{
			Code: 400,
			Msg:  "获取project_id失败",
		}, err
	}
	p := uint(pid)
	userResponse, err := ctrl.service.GetUsers(ctx, p)
	if err != nil {
		return response.Response{}, err
	}

	return response.Response{
		Msg:  "获取成功",
		Code: 200,
		Data: userResponse,
	}, nil
}

// GetAllTags 获取tags
// @Summary 获取某个项目中所有的标签
// @Description 根据 project_id 获取该项目的所有标签
// @Tags Project
// @Accept json
// @Produce json
// @Param project_id path int true "项目ID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求错误（参数错误/无 project_id）"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/project/{project_id}/getAllTags [get]
func (ctrl *ProjectController) GetAllTags(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error) {
	role := cla.UserRule
	if role == 0 {
		return response.Response{
			Code: 403,
			Msg:  "无权限",
		}, nil
	}
	projectID := ctx.Param("project_id")
	if projectID == "" {
		return response.Response{
			Code: 400,
			Msg:  "需要project_id",
		}, nil
	}
	pid, err := strconv.ParseUint(projectID, 10, 64)
	if err != nil {
		return response.Response{
			Code: 400,
			Msg:  "获取project_id失败",
		}, err
	}
	p := uint(pid)
	re, err := ctrl.service.GetAllTags(ctx, p)
	if err != nil {
		return response.Response{
			Code: 400,
		}, err
	}
	return response.Response{
		Code: 200,
		Data: re,
		Msg:  "获取tags成功",
	}, nil
}
