package router

import (
	"github.com/gin-gonic/gin"
	"muxi_auditor/api/request"
	"muxi_auditor/api/response"
	"muxi_auditor/pkg/ginx"
	"muxi_auditor/pkg/jwt"
)

type UserController interface {
	GetUsers(g *gin.Context, req request.GetUserReq) (response.Response, error)
	UpdateUsers(g *gin.Context, req request.UpdateUserRoleReq, cla jwt.UserClaims) (response.Response, error)
}

func UserRoutes(
	s *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	c UserController,
) {
	//认证服务
	UserGroup := s.Group("/user")
	UserGroup.POST("/getUsers", authMiddleware, ginx.WrapReq(c.GetUsers))
	UserGroup.POST("/updateUser", authMiddleware, ginx.WrapClaimsAndReq(c.UpdateUsers))
}
