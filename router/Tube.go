package router

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/ginx"
	"github.com/gin-gonic/gin"
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
