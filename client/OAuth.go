package client

import conf "muxi_auditor/config"

type OAuthClient struct {
	addr         string `yaml:"addr"`
	clientID     string `yaml:"client_id"`
	clientSecret string `yaml:"client_secret"`
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
	return "", nil
}
