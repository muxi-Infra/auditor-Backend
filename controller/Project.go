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
}

func NewProjectController(service *service.ProjectService) *ProjectController {
	return &ProjectController{
		service: service,
	}
}
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
