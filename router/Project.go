package router

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/ginx"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// ProjectController 项目方面接口(即外接应用)
type ProjectController interface {
	GetProjectList(ctx *gin.Context) (response.Response, error)
	Create(ctx *gin.Context, req request.CreateProject) (response.Response, error)
	Detail(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error)
	Delete(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error)
	Update(ctx *gin.Context, req request.UpdateProject, cla jwt.UserClaims) (response.Response, error)
	GetUsers(g *gin.Context, cla jwt.UserClaims) (response.Response, error)
	GetAllTags(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error)
	AddUsers(ctx *gin.Context, req request.AddUsersReq, cla jwt.UserClaims) (response.Response, error)
	DeleteUsers(ctx *gin.Context, req request.DeleteUsers, cla jwt.UserClaims) (response.Response, error)
	GiveProjectRole(ctx *gin.Context, req request.AddUsersReq, cla jwt.UserClaims) (response.Response, error)
}

func RegisterProjectRoutes(
	s *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	c ProjectController,
) {
	//认证服务
	authGroup := s.Group("/project")
	authGroup.GET("/getProjectList", authMiddleware, ginx.Wrap(c.GetProjectList))
	authGroup.POST("/create", authMiddleware, ginx.WrapReq(c.Create))
	authGroup.DELETE("/:project_id/delete", authMiddleware, ginx.WrapClaims(c.Delete))
	authGroup.GET("/:project_id/detail", authMiddleware, ginx.WrapClaims(c.Detail))
	authGroup.POST("/:project_id/update", authMiddleware, ginx.WrapClaimsAndReq(c.Update))
	authGroup.GET("/:project_id/getUsers", authMiddleware, ginx.WrapClaims(c.GetUsers))
	authGroup.GET("/:project_id/getAllTags", authMiddleware, ginx.WrapClaims(c.GetAllTags))
	authGroup.POST("/addUsers", authMiddleware, ginx.WrapClaimsAndReq(c.AddUsers))
	authGroup.DELETE("/deleteUsers", authMiddleware, ginx.WrapClaimsAndReq(c.DeleteUsers))
	authGroup.PUT("giveProjectRole", authMiddleware, ginx.WrapClaimsAndReq(c.GiveProjectRole))
}
