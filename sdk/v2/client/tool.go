package client

const (
	Pending = iota
	Pass
	Reject
)

func TransformStatusToString(status int) string {
	switch status {
	case Pending:
		return "Pending"
	case Pass:
		return "Pass"
	case Reject:
		return "Reject"
	default:
		return "UnknownStatus"
	}
}
