package router

import (
	"github.com/gin-gonic/gin"
	"muxi_auditor/api/response"
	"muxi_auditor/pkg/ginx"
)

type TubeController interface {
	GetQiToken(g *gin.Context) (response.Response, error)
}

func TubeRoutes(
	s *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	c TubeController,
) {
	tubeGroup := s.Group("/tube")
	tubeGroup.Use(authMiddleware)
	tubeGroup.GET("/GetQiToken", authMiddleware, ginx.Wrap(c.GetQiToken))
}
