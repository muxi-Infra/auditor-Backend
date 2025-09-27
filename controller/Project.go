package controller

import (
	"context"
	"errors"
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
	GetProjectList(ctx context.Context, cla jwt.UserClaims) ([]model.ProjectList, error)
	Create(ctx context.Context, project request.CreateProject) (uint, string, error)
	Detail(ctx context.Context, id uint) (response.GetDetailResp, error)
	Delete(ctx context.Context, cla jwt.UserClaims, p uint) error
	Update(ctx context.Context, id uint, req request.UpdateProject) error
	GetUsers(ctx context.Context, id uint) ([]model.UserResponse, error)
	ReturnApiKey(apiKey string, hookUrl string) error
	GetAllTags(ctx context.Context, pid uint) ([]string, error)
	AddUsers(ctx context.Context, role int, uid uint, key string, req []request.AddUser) error
	DeleteUser(ctx context.Context, role int, uid uint, key string, ids []uint) error
	GiveProjectRole(ctx context.Context, userRole int, uid uint, key string, req []request.AddUser) ([]request.AddUser, error)
	SelectUser(ctx context.Context, query string, apiKey string) ([]model.User, error)
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
func (ctrl *ProjectController) GetProjectList(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error) {

	list, err := ctrl.service.GetProjectList(ctx, cla)
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
	id, _, err := ctrl.service.Create(ctx, req)
	if err != nil {
		return response.Response{
			Code: 400,
			Msg:  "创建项目失败",
			Data: err,
		}, nil
	}
	return response.Response{
		Code: 200,
		Msg:  "创建成功",
		Data: id,
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
		if err.Error() == "无权限" {
			return response.Response{
				Code: 400,
				Msg:  "无权限",
				Data: nil,
			}, nil
		}
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
func (ctrl *ProjectController) GetUsers(ctx *gin.Context) (response.Response, error) {
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
		}, errors.New("no power")
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

// AddUsers 添加项目成员
// @Summary 添加项目成员
// @Description 根据api_key和用户id，向项目中批量添加用户
// @Tags Project
// @Accept json
// @Produce json
// @Param api_key header string true "API 认证密钥(api_key)"
// @Param request body request.AddUsersReq true "添加用户请求体"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response "添加成功"
// @Failure 400 {object} response.Response "请求错误（缺少api_key或参数错误）"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/project/addUsers [post]
func (ctrl *ProjectController) AddUsers(ctx *gin.Context, req request.AddUsersReq, cla jwt.UserClaims) (response.Response, error) {
	uid := cla.Uid
	key := ctx.GetHeader("api_key")
	if key == "" {
		return response.Response{
			Code: 400,
			Msg:  "api_key is necessary",
		}, errors.New("api_key is necessary")
	}
	userRole := cla.UserRule
	err := ctrl.service.AddUsers(ctx, userRole, uid, key, req.AddUsers)
	if err != nil {
		return response.Response{
			Code: 400,
			Msg:  err.Error(),
		}, err
	}
	return response.Response{
		Code: 200,
		Msg:  "添加成功",
	}, nil

}

// DeleteUsers 删除项目成员
// @Summary 删除项目成员
// @Description 根据api_key和用户id，批量删除项目中的用户
// @Tags Project
// @Accept json
// @Produce json
// @Param api_key header string true "API 认证密钥(api_key)"
// @Param request body request.DeleteUsers true "删除用户请求体"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求错误（缺少api_key或参数错误）"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/project/deleteUsers [delete]
func (ctrl *ProjectController) DeleteUsers(ctx *gin.Context, req request.DeleteUsers, cla jwt.UserClaims) (response.Response, error) {
	uid := cla.Uid
	key := ctx.GetHeader("api_key")
	if key == "" {
		return response.Response{
			Code: 400,
			Msg:  "api_key is necessary",
		}, nil
	}
	userRole := cla.UserRule
	err := ctrl.service.DeleteUser(ctx, userRole, uid, key, req.Ids)
	if err != nil {
		return response.Response{
			Code: 400,
			Msg:  err.Error(),
		}, err
	}
	return response.Response{
		Code: 200,
		Msg:  "删除成功",
	}, nil
}

// GiveProjectRole 更新项目成员权限
// @Summary 更新项目成员权限
// @Description 根据 api_key 和用户信息，批量更新项目中的用户权
// @Tags Project
// @Accept json
// @Produce json
// @Param api_key header string true "API 认证密钥(api_key)"
// @Param request body request.AddUsersReq true "更新用户角色请求体（包含用户 ID 与新角色）"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求错误（缺少 api_key 或参数错误）"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/project/giveProjectRole [put]
func (ctrl *ProjectController) GiveProjectRole(ctx *gin.Context, req request.AddUsersReq, cla jwt.UserClaims) (response.Response, error) {
	uid := cla.Uid
	key := ctx.GetHeader("api_key")
	if key == "" {
		return response.Response{
			Code: 400,
			Msg:  "api_key is necessary",
		}, nil
	}
	userRole := cla.UserRule
	data, err := ctrl.service.GiveProjectRole(ctx, userRole, uid, key, req.AddUsers)
	if err != nil {
		return response.Response{
			Code: 400,
			Msg:  err.Error(),
			Data: data,
		}, err
	}
	return response.Response{
		Code: 200,
		Msg:  "更新成功",
		Data: data,
	}, nil
}

// SelectUser 搜索用户
// @Summary 根据用户名称搜索用户
// @Description 根据用户名称搜索用户
// @Tags Project
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param api_key header string true "API 认证密钥(api_key)"
// @Param the_query query string true "查询关键字"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求错误（参数错误/无query）"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/project/selectUser [get]
func (ctrl *ProjectController) SelectUser(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error) {
	if cla.UserRule == 0 {
		return response.Response{
			Msg:  "no power",
			Code: 403,
			Data: nil,
		}, errors.New("no power")
	}
	key := ctx.GetHeader("api_key")
	if key == "" {
		return response.Response{
			Code: 400,
			Msg:  "api_key is necessary",
		}, errors.New("api_key is necessary")
	}
	query := ctx.DefaultQuery("the_query", "")
	if query == "" {
		return response.Response{
			Msg:  "query is necessary",
			Code: 400,
			Data: nil,
		}, errors.New("query is necessary")
	}
	users, err := ctrl.service.SelectUser(ctx, query, key)
	if err != nil {
		return response.Response{
			Msg:  "数据搜索有误",
			Code: 400,
			Data: nil,
		}, err
	}
	var da []response.UserInfo
	for _, user := range users {
		da = append(da, response.UserInfo{
			Avatar: user.Avatar,
			Name:   user.Name,
			Id:     user.ID,
			Role:   user.UserRole,
			Email:  user.Email,
		})
	}
	return response.Response{
		Msg:  "",
		Code: 200,
		Data: da,
	}, nil
}
