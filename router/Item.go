package router

import (
	"github.com/gin-gonic/gin"
	"muxi_auditor/api/request"
	"muxi_auditor/api/response"
	"muxi_auditor/pkg/ginx"
	"muxi_auditor/pkg/jwt"
)

type ItemController interface {
	Select(g *gin.Context, cla jwt.UserClaims, req request.SelectReq) (response.Response, error)
	Audit(g *gin.Context, cla jwt.UserClaims, req request.AuditReq) (response.Response, error)
	SearchHistory(g *gin.Context, cla jwt.UserClaims) (response.Response, error)
	Upload(g *gin.Context, cla jwt.UserClaims, req request.UploadReq) (response.Response, error)
}

func ItemRoutes(
	s *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	c ItemController,
) {
	ItemGroup := s.Group("/item")
	ItemGroup.POST("/select", ginx.WrapClaimsAndReq(c.Select))
	ItemGroup.POST("/audit", ginx.WrapClaimsAndReq(c.Audit))
	ItemGroup.GET("/searchHistory", authMiddleware, ginx.WrapClaims(c.SearchHistory))
	ItemGroup.POST("/upload", authMiddleware, ginx.WrapClaimsAndReq(c.Upload))
}
