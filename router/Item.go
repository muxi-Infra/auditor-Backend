package router

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/ginx"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// ItemController 需要审核的条目方面接口
type ItemController interface {
	Select(g *gin.Context, req request.SelectReq) (response.Response, error)
	Audit(g *gin.Context, req request.AuditReq, cla jwt.UserClaims) (response.Response, error)
	SearchHistory(g *gin.Context, cla jwt.UserClaims) (response.Response, error)
	Upload(g *gin.Context, req request.UploadReq) (response.Response, error)
	Detail(g *gin.Context) (response.Response, error)
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
	ItemGroup.PUT("/upload", authMiddleware, ginx.WrapReq(c.Upload))
	ItemGroup.GET("/:item_id/detail", authMiddleware, ginx.Wrap(c.Detail))
}
