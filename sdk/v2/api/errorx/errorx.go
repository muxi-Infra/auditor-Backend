package errorx

import "fmt"

// 似乎把成功的code也放在error这里有些不太合适，但一时也没想好放哪
const (
	// 200 成功系列
	SuccessCode          = 20010
	SeverDataIllegalCode = 20011
	// 400系列
	RequestErrCode = 40010
	MarshalErrCode = 40019
)

var (
	DefaultErr = &SDKError{
		HttpCode: 0,
		Code:     0,
		Msg:      "",
		Category: "",
		Cause:    nil,
	}

	MarshalErr = func(e error) *SDKError {
		return New(0, MarshalErrCode, "marshal wrong", e)
	}

	SeverDataIllegal = func(e error, httpCode int) *SDKError {
		return New(httpCode, SeverDataIllegalCode, "server return data illegal", e)
	}
)

type SDKError struct {
	HttpCode int    `json:"httpCode"`
	Code     int    `json:"code"`
	Msg      string `json:"message"`
	Category string `json:"category"` //具体分类
	Cause    error  `json:"cause"`    // 具体错误原因
}

func (e *SDKError) Error() string {
	return fmt.Sprintf("type:%s [%d] %s : %v", e.Category, e.Code, e.Msg, e.Cause)
}

func New(httpCode int, code int, message string, cause error) *SDKError {
	return &SDKError{
		HttpCode: httpCode,
		Code:     code,
		Msg:      message,
		Cause:    cause,
	}
}
