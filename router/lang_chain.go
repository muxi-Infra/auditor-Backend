package router

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/ginx"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type LLMController interface {
	Audit(ctx *gin.Context, req request.AuditByLLMReq, cla jwt.UserClaims) (response.Response, error)
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
