package response

import "muxi_auditor/repository/model"

type LoginResp struct {
	Token string `json:"token"`
}
type Response struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
type GetDetailResp struct {
	TotleNumber   int                  `json:"totle_number"`
	CurrentNumber int                  `json:"current_number"`
	Apikey        string               `json:"api_key"`
	AuditRule     string               `json:"audit_rule"`
	Members       []model.UserResponse `json:"members"`
}
