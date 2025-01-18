package ginx

import (
	"errors"
	"github.com/gin-gonic/gin"
	api_errors "muxi_auditor/api/errors"
	"net/http"
)

// ctx表示上下文,req表示请求结构体,Resp表示响应结构体(这里全部填response.Response),UserClaims表示用户信息
func WrapClaimsAndReq[Req any, UserClaims any, Resp any](fn func(*gin.Context, Req, UserClaims) (Resp, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//解析请求
		var req Req
		err := bind(ctx, req)
		if err != nil {
			ctx.Set("err", err)
		}
		//调用函数
		res, err := fn(ctx, req, getClaims[UserClaims](ctx))
		if err != nil {
			ctx.Set("err", err)
		} else {
			ctx.Set("resp", res)
		}

	}
}

// WrapReq 。用于处理有请求体的请求
// ctx表示上下文,req表示请求结构体,Resp表示响应结构体(这里全部填response.Response)
func WrapReq[Req any, Resp any](fn func(*gin.Context, Req) (Resp, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//解析参数
		var req Req
		err := bind(ctx, req)
		if err != nil {
			ctx.Set("err", err)
		}

		// 调用业务逻辑函数
		res, err := fn(ctx, req)
		if err != nil {
			ctx.Set("err", err)
		} else {
			ctx.Set("resp", res)
		}
	}
}

// Wrap 。用于处理没有请求体的请求
// ctx表示上下文,Resp表示响应结构体(这里全部填response.Response)
func Wrap[Resp any](fn func(*gin.Context) (Resp, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := fn(ctx)
		if err != nil {
			ctx.Set("err", err)
		} else {
			ctx.Set("resp", res)
		}
	}
}

// WrapClaims 用于处理有用户验证但是没有请求体的请求
// ctx表示上下文,Resp表示响应结构体(这里全部填response.Response),UserClaims表示用户信息
func WrapClaims[UserClaims any, Resp any](fn func(*gin.Context, UserClaims) (Resp, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := fn(ctx, getClaims[UserClaims](ctx))
		if err != nil {
			ctx.Set("err", err)
		} else {
			ctx.Set("resp", res)
		}
	}
}

// 解析参数通用函数
func bind(ctx *gin.Context, req any) error {
	var err error
	// 根据请求方法选择合适的绑定方式
	if ctx.Request.Method == http.MethodGet {
		err = ctx.ShouldBindQuery(&req) // 处理GET请求的查询参数
	} else {
		err = ctx.ShouldBind(&req) // 处理POST、PUT等请求的请求体数据
	}

	if err != nil {
		return api_errors.BAD_ENTITY_ERROR(err)
	}

	return nil
}

// 获取Claims通用函数
func getClaims[UserClaims any](ctx *gin.Context) (claims UserClaims) {
	rawVal, ok := ctx.Get("user")
	if !ok {
		ctx.Set("err", api_errors.BAD_ENTITY_ERROR(errors.New("从上下文获取userClaims失败")))
		ctx.Abort()
		return
	}

	// 注意，这里要求放进去 ctx 的不能是*UserClaims，这是常见的一个错误
	claims, ok = rawVal.(UserClaims)
	if !ok {
		ctx.Set("err", api_errors.BAD_ENTITY_ERROR(errors.New("userClaims类型断言失败了")))
		ctx.Abort()
		return
	}
	return claims
}
