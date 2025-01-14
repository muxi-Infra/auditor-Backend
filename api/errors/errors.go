package errors

import "muxi_auditor/pkg/errorx"

// OAUTH_ERROR
var (
	OAUTH_GETINFO_ERROR = func(err error) error {
		return errorx.New(500, 50001, "从通行证获取用户信息失败!", "OAuth", err)
	}
	LOGIN_ERROR = func(err error) error {
		return errorx.New(500, 50002, "系统发生内部错误,登陆失败!", "User", err)
	}
)

// Common 错误快捷方法,一般不推荐使用
func BadRequest(msg string) error {
	return errorx.New(400, 40001, msg, "Common", nil)
}

func Unauthorized(msg string) error {
	return errorx.New(401, 40001, msg, "Common", nil)
}

func InternalServerError(msg string, cause error) error {
	return errorx.New(500, 50001, msg, "Common", cause)
}

func NewError(httpCode int, code int, msg string, err error) error {
	return errorx.New(httpCode, code, msg, "common", err)
}
