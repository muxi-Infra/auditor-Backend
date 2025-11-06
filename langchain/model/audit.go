package model

const (
	Pending = iota
	Pass
	Reject
	PassBeforeHook
	RejectBeforeHook
)

type AuditContent[T1 any, T2 any] struct {
	Text T1
	Imgs T2
}

type AuditResult struct {
	ID         uint
	Result     int     `json:"result"` // 仅允许 pass|review|reject
	Reason     string  `json:"reason"`
	Confidence float32 `json:"confidence"` //between 0 and 100
}

type ImageAuditResult struct {
	ID         string
	Result     int      `json:"result"`
	Reason     []string `json:"reason"`
	Confidence float32  `json:"confidence"`
}

type TextAuditResult struct {
	ID         string
	Result     int      `json:"result"`
	Reason     []string `json:"reason"`
	Confidence float32  `json:"confidence"`
}
