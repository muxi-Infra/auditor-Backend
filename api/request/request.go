package request

import "muxi_auditor/repository/model"

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
	Role       int `json:"role"`
	Project_id int `json:"project_id"`
}
type UpdateUserRoleReq struct {
	Role          int                   `json:"role"`
	UserId        int                   `json:"user_id"`
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
