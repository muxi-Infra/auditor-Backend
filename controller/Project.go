package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"muxi_auditor/api/request"
	"muxi_auditor/api/response"
	"muxi_auditor/pkg/ginx"
	"muxi_auditor/pkg/jwt"
	"muxi_auditor/repository/model"
	"muxi_auditor/service"
)

type ProjectController struct {
	service ProjectService
}
type ProjectService interface {
	GetProjectList(ctx context.Context, logo string) ([]model.ProjectList, error)
	Create(ctx context.Context, name string, logo string, audioRule string, ids []uint) error
	Detail(ctx context.Context, id uint) (response.GetDetailResp, error)
	Delete(ctx context.Context, cla jwt.UserClaims, req request.DeleteProject) error
}

func NewProjectController(service *service.ProjectService) *ProjectController {
	return &ProjectController{
		service: service,
	}
}

// GetProjectList @Summary 获取项目列表
// @Description 获取所有项目列表，根据 logo 过滤
// @Tags Project
// @Accept json
// @Produce json
// @Param X-Header-Param header string true "Logo过滤字段"
// @Success 200 {object} response.Response "成功返回项目列表"
// @Failure 400 {object} response.Response "获取项目列表失败"
// @Security ApiKeyAuth
// @Router /api/v1/project/getProjectList [get]
func (ctrl *ProjectController) GetProjectList(ctx *gin.Context) (response.Response, error) {
	logo := ctx.GetHeader("X-Header-Param")
	list, err := ctrl.service.GetProjectList(ctx, logo)
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

// Create @Summary 创建项目
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
	err = ctrl.service.Create(ctx, req.Name, req.Logo, req.AudioRule, req.UserIds)
	if err != nil {
		return response.Response{
			Code: 400,
			Msg:  "创建项目失败",
		}, nil
	}
	return response.Response{
		Code: 200,
		Msg:  "创建成功",
		Data: nil,
	}, nil
}

// Detail @Summary 获取项目详细信息
// @Description 根据项目 ID 获取项目的详细信息
// @Tags Project
// @Accept json
// @Produce json
// @Param projectId query uint true "项目ID"
// @Success 200 {object} response.Response{data=response.GetDetailResp} "获取项目详细信息成功"
// @Failure 400 {object} response.Response "获取项目详细信息失败"
// @Security ApiKeyAuth
// @Router /api/v1/project/detail [get]
func (ctrl *ProjectController) Detail(ctx *gin.Context, req request.GetProjectDetail) (response.Response, error) {
	_, err := ginx.GetClaims[jwt.UserClaims](ctx)
	if err != nil {
		return response.Response{
			Msg:  "",
			Code: 40001,
			Data: nil,
		}, err
	}
	detail, err := ctrl.service.Detail(ctx, req.ProjectId)
	if err != nil {
		return response.Response{
			Msg:  "",
			Code: 40001,
			Data: nil,
		}, err
	}

	return response.Response{
		Code: 200,
		Msg:  "获取成功",
		Data: detail,
	}, nil

}

// Delete
// @Summary 删除项目
// @Description 通过项目 ID 删除指定的项目
// @Tags Project
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param delete body request.DeleteProject true "删除项目请求参数"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "删除失败"
// @Security ApiKeyAuth
// @Router /api/v1/project/delete [delete]
func (ctrl *ProjectController) Delete(ctx *gin.Context, cla jwt.UserClaims, req request.DeleteProject) (response.Response, error) {
	err := ctrl.service.Delete(ctx, cla, req)
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
