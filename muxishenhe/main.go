package main

import (
	"github.com/gin-gonic/gin"
	"muxishenhe/login"
	"muxishenhe/mytoken"
	"net/http"
)

func main() {

	r := gin.Default()
	r.POST("/login", func(c *gin.Context) {
		a := login.GetCode(c)
		b := login.GetToken(a)
		d := login.Getuserinfo(b.AccessToken)
		if d.Email == "" || d.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"Message": "登录失败",
			})
		}
		token, err := mytoken.GenerateJwt(b.AccessToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"Message": "生成JWT失败",
			})
		}
		c.Header("Authorization", "Bearer "+token)
		c.JSON(200, gin.H{
			"code":    "200",
			"Message": "登录成功",
		})
	})
	r.Run(":8080")
}
