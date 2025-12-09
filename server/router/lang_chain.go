package router

import (
	"github.com/gin-gonic/gin"

	"github.com/muxi-Infra/auditor-Backend/api/request"
	"github.com/muxi-Infra/auditor-Backend/api/response"
	"github.com/muxi-Infra/auditor-Backend/pkg/ginx"
	"github.com/muxi-Infra/auditor-Backend/pkg/jwt"
)

type LLMController interface {
	Audit(ctx *gin.Context, req request.AuditByLLMReq, cla jwt.UserClaims) (response.Response, error)
	Close()
}

// LLMRoutes 其他应用上传或修改item的接口
func LLMRoutes(
	s *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	c LLMController,
) {
	LLMGroup := s.Group("/llm")
	LLMGroup.POST("/audit", authMiddleware, ginx.WrapClaimsAndReq(c.Audit))
}
