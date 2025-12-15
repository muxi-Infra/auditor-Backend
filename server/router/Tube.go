package router

import (
	"github.com/gin-gonic/gin"

	"github.com/muxi-Infra/auditor-Backend/api/response"
	"github.com/muxi-Infra/auditor-Backend/pkg/ginx"
)

// TubeController 获取图床token的接口
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
