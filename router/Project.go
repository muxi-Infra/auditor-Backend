package router

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/ginx"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type ProjectController interface {
	GetProjectList(ctx *gin.Context) (response.Response, error)
	Create(ctx *gin.Context, req request.CreateProject) (response.Response, error)
	Detail(ctx *gin.Context) (response.Response, error)
	Delete(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error)
	Update(ctx *gin.Context, req request.UpdateProject, cla jwt.UserClaims) (response.Response, error)
	GetUsers(g *gin.Context, cla jwt.UserClaims) (response.Response, error)
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
	authGroup.GET("/:project_id/detail", authMiddleware, ginx.Wrap(c.Detail))
	authGroup.POST("/:project_id/update", authMiddleware, ginx.WrapClaimsAndReq(c.Update))
	authGroup.GET("/:project_id/getUsers", authMiddleware, ginx.WrapClaims(c.GetUsers))
}
