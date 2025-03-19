package router

import (
	"github.com/gin-gonic/gin"
	"muxi_auditor/api/request"
	"muxi_auditor/api/response"
	"muxi_auditor/pkg/ginx"
	"muxi_auditor/pkg/jwt"
)

type UserController interface {
	UpdateUsers(g *gin.Context, req request.UpdateUserRoleReq, cla jwt.UserClaims) (response.Response, error)
	GetMyInfo(g *gin.Context, cla jwt.UserClaims) (response.Response, error)
	UpdateMyInfo(g *gin.Context, req request.UpdateUserReq, cla jwt.UserClaims) (response.Response, error)
}

func UserRoutes(
	s *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	c UserController,
) {
	//认证服务
	UserGroup := s.Group("/user")

	UserGroup.POST("/updateUser", authMiddleware, ginx.WrapClaimsAndReq(c.UpdateUsers))
	UserGroup.GET("/getMyInfo", authMiddleware, ginx.WrapClaims(c.GetMyInfo))
	UserGroup.POST("/updateMyInfo", authMiddleware, ginx.WrapClaimsAndReq(c.UpdateMyInfo))
}
