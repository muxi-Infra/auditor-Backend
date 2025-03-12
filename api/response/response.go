package response

import (
	"muxi_auditor/repository/model"
	"time"
)

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
type SelectResp struct {
	ProjectId uint   `json:"project_id"`
	Items     []Item `json:"items"`
}

type Item struct {
	ItemId     uint      `json:"item_id"`
	Author     string    `json:"author"`
	Tags       []string  `json:"tags"`
	Status     int       `json:"status"`
	PublicTime time.Time `json:"public_time"`
	Auditor    string    `json:"auditor"`
	Content    Contents  `json:"content"`
}
type Contents struct {
	Topic       Topics  `json:"topic"`
	LastComment Comment `json:"last_comment"`
	NextComment Comment `json:"next_comment"`
}
type Topics struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Pictures []string `json:"pictures"`
}
type Comment struct {
	Content  string   `json:"content"`
	Pictures []string `json:"pictures"`
}
type UserInfo struct {
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
	Role   int    `json:"role"`
	Email  string `json:"email"`
}
