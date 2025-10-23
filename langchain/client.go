package langchain

type AuditAI interface {
	GetToolList() []string
	SendMessage() interface{}
}
