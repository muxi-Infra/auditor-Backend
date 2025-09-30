package langchain

type AuditAI interface {
	Init()
	GetToolList() []string
	SendMessage() interface{}
}
