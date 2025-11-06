package errorx

import "github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/stringx"

type AIError struct {
	Domain  string // 错误码
	Message string // 错误消息
	err     error
}

func (e *AIError) Error() string {
	return stringx.Build("[", e.Domain, "]: ", e.Message, "---", "err:", e.err.Error())
}

// Unwrap 便于去error.Is
func (e *AIError) Unwrap() error {
	return e.err
}

func (e *AIError) Wrap(err error) *AIError {
	e.err = err
	return e
}

func NewAIError(err error, domain, msg string) *AIError {
	return &AIError{
		err:     err,
		Domain:  domain,
		Message: msg,
	}
}

// 预定义常用错误
var (
	ErrTaskTypeNotRegistered = &AIError{Domain: "TaskTypeNotRegistered", Message: "任务类型未注册"}
	ErrDataInvalid           = &AIError{Domain: "DataInvalid", Message: "返回数据格式错误"}
	ErrDecodeFailed          = &AIError{Domain: "DecodeFailed", Message: "数据解析失败"}
	ErrClientInitFailed      = &AIError{Domain: "ClientInitErr", Message: "ai客户端初始化出错"}
	ErrTextAuditErr          = &AIError{Domain: "阿里云内容审核", Message: "文本审核失败"}
)
