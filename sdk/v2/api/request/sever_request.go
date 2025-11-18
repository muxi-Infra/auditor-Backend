package request

type HookPayload struct {
	Event string      `json:"event"`
	Data  WebHookData `json:"data,omitempty"`
	Try   int         `json:"try"` // 重试次数，最大不超过五
}

type WebHookData struct {
	Id     uint
	Status string
	Msg    string
}
