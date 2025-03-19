package response

type LoginResp struct {
	Token string `json:"token"`
}
type Response struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
type GetDetailResp struct {
	TotleNumber   int    `json:"totle_number"`   //项目中item总数
	CurrentNumber int    `json:"current_number"` //未审核的数目
	Apikey        string `json:"api_key"`        //由project_id生成的key
	AuditRule     string `json:"audit_rule"`
}
type SelectResp struct {
	Items []Item `json:"items"`
}

type Item struct {
	ItemId     uint     `json:"item_id"`
	Author     string   `json:"author"`
	Tags       []string `json:"tags"`
	Status     int      `json:"status"`
	PublicTime int64    `json:"public_time"`
	Auditor    uint     `json:"auditor"`
	Content    Contents `json:"content"` //item具体内容，包含题目内容和评论
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
	Role   int    `json:"role"` //用户权限
	Email  string `json:"email"`
}
