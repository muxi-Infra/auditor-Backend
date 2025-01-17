package response

type LoginResp struct {
	Token string `json:"token"`
}
type Response struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
