package login

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"muxishenhe/token"
	"muxishenhe/user"
	"net/http"
)

type Code struct {
	Code    string `json:"code"`
	Expired string `json:"expired"`
}
type AccessCode struct {
	Code         string `json:"code"`
	ClientSecret string `json:"client_secret"`
}

// 从前端拿到code
func GetCode(c *gin.Context) AccessCode {
	var code Code
	var accesscode AccessCode
	if err := c.ShouldBind(&code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return AccessCode{
			Code:         "",
			ClientSecret: "",
		}
	}
	accesscode = AccessCode{
		Code:         code.Code,
		ClientSecret: "d8e4b18e-8c61-40b8-8e2c-c666e8b1164a",
	}
	return accesscode
}

// 从通行证拿到验证的token
func GetToken(accesscode AccessCode) token.Token {
	var accesstoken token.Token
	a := map[string][]string{
		"client_secret": {accesscode.ClientSecret},
		"code":          {accesscode.Code},
	}
	res, err := http.PostForm("http://pass.muxi-tech.xyz/oauth/token", a)
	if err != nil {
		log.Println(err)
		return accesstoken
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return accesstoken
	}
	err = json.Unmarshal(body, &accesstoken)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	return accesstoken

}
func Getuserinfo(accesstoken string) user.UserStruct {
	var user user.UserStruct
	res, err := http.NewRequest("GET", "http://pass.muxi-tech.xyz/auth/api/user", nil)
	res.Header.Set("token", accesstoken)
	c := http.Client{}
	resp, err := c.Do(res)
	if err != nil {
		log.Println(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println(err)
	}
	return user
}
