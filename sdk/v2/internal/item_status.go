package internal

type ItemStatus struct {
	HookId int `json:"hook_id"`
	Status int `json:"status"`
}

type ItemsStatus struct {
	Items []ItemStatus `json:"items"`
}
