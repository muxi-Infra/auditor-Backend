package router

import (
	"github.com/gin-gonic/gin"
	"muxi_auditor/api/request"
	"muxi_auditor/api/response"
	"muxi_auditor/pkg/ginx"
)

type ProjectController interface {
	GetProjectList(ctx *gin.Context) (response.Response, error)
	Create(ctx *gin.Context, req request.CreateProject) (response.Response, error)
	Detail(ctx *gin.Context, req request.GetProjectDetail) (response.Response, error)
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
	authGroup.GET("/detail", authMiddleware, ginx.WrapReq(c.Detail))
}
