package model

type AuditResult struct {
	ID         uint
	Result     int     `json:"result"` // 仅允许 pass|review|reject
	Reason     string  `json:"reason"`
	Confidence float32 `json:"confidence"` //between 0 and 1
}
