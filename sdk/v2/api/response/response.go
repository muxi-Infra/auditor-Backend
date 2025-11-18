package response

// 我认为对于sdk而言，用户需要的是确定类型的返回结果，这样用户更易阅读，使用起来更加轻松固针对每个接口独立设计resp,并不使用通用resp.

type UploadItemResp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	ItemID int    `json:"data"`
	Errorx error  `json:"errorx"` // 这里始终传errorx.SDKError
}

type UpdateItemResp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	ItemID int    `json:"data"`
	Errorx error  `json:"errorx"`
}

type DeleteItemResp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	ItemID int    `json:"data"`
	Errorx error  `json:"errorx"`
}

type GetItemsResp struct {
	Code   int           `json:"code"`
	Msg    string        `json:"msg"`
	Items  []ItemsStatus `json:"data"`
	Errorx error         `json:"errorx"`
}

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type ItemsStatus struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
}
