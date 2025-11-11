package response

type AuditByLLMResponse struct {
	ID     uint   `json:"id"`
	Status int    `json:"status"`
	Reason string `json:"reason"`
}
