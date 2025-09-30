package errorx

type AIError struct {
	Code    string // 错误码
	Message string // 错误消息
	err     error
}

func (e *AIError) Error() string {
	return "[" + e.Code + "] " + e.Message + ":" + e.err.Error()
}

// Unwrap 便于去error.Is
func (e *AIError) Unwrap() error {
	return e.err
}

// 预定义常用错误
var (
	ErrTaskTypeNotRegistered = &AIError{Code: "TaskTypeNotRegistered", Message: "任务类型未注册"}
	ErrDataInvalid           = &AIError{Code: "DataInvalid", Message: "返回数据格式错误"}
	ErrDecodeFailed          = &AIError{Code: "DecodeFailed", Message: "数据解析失败"}
)
