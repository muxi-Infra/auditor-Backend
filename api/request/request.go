package request

import (
	"muxi_auditor/api/response"
	"muxi_auditor/repository/model"
)

type LoginReq struct {
	Code string `json:"code"`
}
type RegisterReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
type UpdateUserReq struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}
type GetUserReq struct {
	Role       int  `json:"role"`
	Project_id uint `json:"project_id"`
}
type UpdateUserRoleReq struct {
	Role          int                   `json:"role"`
	UserId        uint                  `json:"user_id"`
	ProjectPermit []model.ProjectPermit `json:"project_permit"`
}
type CreateProject struct {
	Name      string `json:"name"`
	Logo      string `json:"logo"`
	AudioRule string `json:"audio_rule"`
	UserIds   []uint `json:"user_ids"`
}
type GetProjectDetail struct {
	ProjectId uint `json:"project_id"`
}
type SelectReq struct {
	ProjectID uint `json:"project_id"`
	//RoundTime [][]int  `json:"round_time"`
	Tag      string `json:"tag"`
	Status   string `json:"status"`
	Auditor  string `json:"auditor"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Query    string `json:"query"`
}
type AuditReq struct {
	ProjectID uint   `json:"project_id"`
	Reason    string `json:"reason"`
	Status    int    `json:"status"`
	ItemId    uint   `json:"item_id"`
}
type UploadReq struct {
	ApiKey     string            `json:"api_key"`
	HookUrl    string            `json:"hook_url"`
	Id         int               `json:"id"`
	Auditor    string            `json:"auditor"`
	Author     string            `json:"author"`
	PublicTime string            `json:"public_time"`
	Tags       []string          `json:"tags"`
	Content    response.Contents `json:"content"`
	Extra      interface{}       `json:"extra"`
}
type DeleteProject struct {
	ProjectId uint `json:"project_id"`
}