package router

import (
	"github.com/gin-gonic/gin"
	"muxi_auditor/api/request"
	"muxi_auditor/api/response"
	"muxi_auditor/pkg/ginx"
	"muxi_auditor/pkg/jwt"
)

type ItemController interface {
	Select(g *gin.Context, req request.SelectReq) (response.Response, error)
	Audit(g *gin.Context, req request.AuditReq, cla jwt.UserClaims) (response.Response, error)
	SearchHistory(g *gin.Context, cla jwt.UserClaims) (response.Response, error)
	Upload(g *gin.Context, req request.UploadReq, cla jwt.UserClaims) (response.Response, error)
}

func ItemRoutes(
	s *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	c ItemController,
) {
	ItemGroup := s.Group("/item")
	ItemGroup.POST("/select", authMiddleware, ginx.WrapReq(c.Select))
	ItemGroup.POST("/audit", authMiddleware, ginx.WrapClaimsAndReq(c.Audit))
	ItemGroup.GET("/searchHistory", authMiddleware, ginx.WrapClaims(c.SearchHistory))
	ItemGroup.POST("/upload", authMiddleware, ginx.WrapClaimsAndReq(c.Upload))
}
