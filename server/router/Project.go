package router

import (
	"github.com/gin-gonic/gin"

	"github.com/muxi-Infra/auditor-Backend/api/request"
	"github.com/muxi-Infra/auditor-Backend/api/response"
	"github.com/muxi-Infra/auditor-Backend/pkg/ginx"
	"github.com/muxi-Infra/auditor-Backend/pkg/jwt"
)

// ProjectController 项目方面接口(即外接应用)
type ProjectController interface {
	GetProjectList(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error)
	Create(ctx *gin.Context, req request.CreateProject) (response.Response, error)
	Detail(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error)
	Delete(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error)
	Update(ctx *gin.Context, req request.UpdateProject, cla jwt.UserClaims) (response.Response, error)
	GetUsers(g *gin.Context) (response.Response, error)
	GetAllTags(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error)
	AddUsers(ctx *gin.Context, req request.AddUsersReq, cla jwt.UserClaims) (response.Response, error)
	DeleteUsers(ctx *gin.Context, req request.DeleteUsers, cla jwt.UserClaims) (response.Response, error)
	GiveProjectRole(ctx *gin.Context, req request.AddUsersReq, cla jwt.UserClaims) (response.Response, error)
	SelectUser(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error)
	GetItemNums(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error)
}

func RegisterProjectRoutes(
	s *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	c ProjectController,
) {
	// 项目服务
	projectGroup := s.Group("/project")
	projectGroup.GET("/getProjectList", authMiddleware, ginx.WrapClaims(c.GetProjectList))
	projectGroup.POST("/create", authMiddleware, ginx.WrapReq(c.Create))
	projectGroup.DELETE("/:project_id", authMiddleware, ginx.WrapClaims(c.Delete))
	projectGroup.GET("/:project_id/detail", authMiddleware, ginx.WrapClaims(c.Detail))
	projectGroup.POST("/:project_id/update", authMiddleware, ginx.WrapClaimsAndReq(c.Update))
	projectGroup.GET("/:project_id/getUsers", authMiddleware, ginx.Wrap(c.GetUsers))
	projectGroup.GET("/:project_id/getAllTags", authMiddleware, ginx.WrapClaims(c.GetAllTags))
	projectGroup.POST("/addUsers", authMiddleware, ginx.WrapClaimsAndReq(c.AddUsers))
	projectGroup.DELETE("/deleteUsers", authMiddleware, ginx.WrapClaimsAndReq(c.DeleteUsers))
	projectGroup.PUT("giveProjectRole", authMiddleware, ginx.WrapClaimsAndReq(c.GiveProjectRole))
	projectGroup.GET("/selectUser", authMiddleware, ginx.WrapClaims(c.SelectUser))
	projectGroup.GET("/:project_id/getItemNums", authMiddleware, ginx.WrapClaims(c.GetItemNums))
}
