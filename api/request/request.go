package request

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/model"
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
	Avatar string `json:"avatar"`
}
type GetUserReq struct {
	Project_id uint `json:"project_id"`
}
type UpdateUserRoleReq struct {
	Role          int                   `json:"role"` //用户权限
	UserId        uint                  `json:"user_id"`
	ProjectPermit []model.ProjectPermit `json:"project_permit"` //允许的项目列表
}
type CreateProject struct {
	Name      string `json:"name"`
	Logo      string `json:"logo"`
	AudioRule string `json:"audio_rule"` //审核规则
	UserIds   []uint `json:"user_ids"`
	HookUrl   string `json:"hook_url"`
}
type GetProjectDetail struct {
	ProjectId uint `json:"project_id"`
}
type SelectReq struct {
	ProjectID uint     `json:"project_id"`
	RoundTime [][]int  `json:"round_time"` //日期
	Tags      []string `json:"tags"`       //标签
	Statuses  []int    `json:"statuses"`
	Auditors  []uint   `json:"auditors"`
	Page      int      `json:"page"`
	PageSize  int      `json:"page_size"`
	Query     string   `json:"query"` //查询字段
}
type AuditReq struct {
	Reason string `json:"reason"`
	Status int    `json:"status"` //0未审核，1通过，2未通过
	ItemId uint   `json:"item_id"`
}
type UploadReq struct {
	HookUrl    string            `json:"hook_url"`
	Id         uint              `json:"id"` //hook_id
	Auditor    uint              `json:"auditor"`
	Author     string            `json:"author"`
	PublicTime int64             `json:"public_time"`
	Tags       []string          `json:"tags"`
	Content    response.Contents `json:"content"`
	Extra      interface{}       `json:"extra"`
}
type DeleteProject struct {
	ProjectId uint `json:"project_id"`
}
type UpdateProject struct {
	Logo      string `json:"logo"`
	AudioRule string `json:"audio_rule"`
}

type GetUsers struct {
	Query    string `json:"query"`
	Page     int    `json:"page"`
	PageSize int    `json:"size"`
}
type HookPayload struct {
	Event string `json:"event"`
	Data  any    `json:"data,omitempty"`
	Try   int    `json:"try"` // 重试次数，最大不超过五
}
type ReturnSecret struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Message   string `json:"message"`
}
