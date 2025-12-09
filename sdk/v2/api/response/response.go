package response

// 我认为对于sdk而言，用户需要的是确定类型的返回结果，这样用户更易阅读，使用起来更加轻松固针对每个接口独立设计resp,并不使用通用resp.

type UploadItemResp struct {
	Basic  Basic `json:"basic"`
	ItemID int   `json:"data"`
	// 这里始终传errorx.SDKError
}

type UpdateItemResp struct {
	Basic  Basic `json:"basic"`
	ItemID int   `json:"data"`
}

type DeleteItemResp struct {
	Basic  Basic `json:"basic"`
	ItemID int   `json:"data"`
}

type GetItemsResp struct {
	Basic Basic         `json:"basic"`
	Items []ItemsStatus `json:"data"`
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

type Basic struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Errorx error  `json:"errorx"`
}
