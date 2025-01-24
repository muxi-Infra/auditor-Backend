package token

type Token struct {
	AccessToken    string `json:"access_token"`
	AccessExpired  int    `json:"access_expired"`
	RefreshToken   string `json:"refresh_token"`
	RefreshExpired int    `json:"refresh_expired"`
}
