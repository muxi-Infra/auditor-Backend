package strategy

type AIAuditStrategy string

const (
	AuditByBackend AIAuditStrategy = "BACKEND"
	AuditByFront   AIAuditStrategy = "FRONT"
)
