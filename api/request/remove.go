package request

type RemoveUpdateReq struct {
	HookUrl    string      `json:"hook_url"`
	Id         uint        `json:"id"` //hook_id
	Author     string      `json:"author"`
	PublicTime int64       `json:"public_time"`
	Tags       []string    `json:"tags"`
	Content    Contents    `json:"content"`
	Extra      interface{} `json:"extra"`
}
