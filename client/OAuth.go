package client

import (
	"encoding/json"
	"errors"
	"fmt"
	conf "github.com/cqhasy/2025-Muxi-Team-auditor-Backend/config"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type OAuthClient struct {
	addr         string `yaml:"addr"`
	clientID     string `yaml:"client_id"`
	clientSecret string `yaml:"client_secret"`
}
type OAuthToken struct {
	AccessToken string `json:"access_token"`
}
type userEmail struct {
	Email string `json:"email"`
}
type Response1 struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data OAuthToken `json:"data"`
}
type Response2 struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"`
	Data userEmail `json:"data"`
}

func NewOAuthClient(config *conf.OAuthConfig) *OAuthClient {
	return &OAuthClient{
		addr:         config.Addr,
		clientID:     config.ClientID,
		clientSecret: config.ClientSecret,
	}
}

//这里放的是用来调用固定第三方的客户端结构体

// 待完善
func (c *OAuthClient) GetOAuth(code string) (string, error) {
	formData := url.Values{}
	formData.Set("client_secret", c.clientSecret)

	formData.Set("code", code)
	params := url.Values{}
	params.Set("client_id", c.clientID)
	params.Set("grant_type", "authorization_code")
	params.Set("response_type", "token")

	req, err := http.NewRequest("POST", c.addr+"/oauth/token"+"?"+params.Encode(), strings.NewReader(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("获取token的请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("获取token的响应失败：: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var re Response1
	err = json.Unmarshal(body, &re)
	if err != nil || re.Data.AccessToken == "" {
		return "", errors.New("未能获取accessToken")
	}

	return re.Data.AccessToken, nil
}
func (c *OAuthClient) GetEmail(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", c.addr+"/user", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("token", accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var userEmail Response2
	err = json.Unmarshal(body, &userEmail)
	if err != nil {
		return "", err
	}
	return userEmail.Data.Email, nil
}
