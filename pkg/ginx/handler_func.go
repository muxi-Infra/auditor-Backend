package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"muxi_auditor/pkg/jwt"

	"muxi_auditor/pkg/errorx"
	"muxi_auditor/pkg/logger"
	"net/http"
)

// 受制于泛型，这里只能使用包变量
var log logger.Logger

var vector *prometheus.CounterVec

func SetPrometheusCounter(opt prometheus.CounterOpts) {
	// 初始化 Prometheus 计数器向量
	vector = prometheus.NewCounterVec(opt, []string{"code"})
	prometheus.MustRegister(vector)
}

func SetLogger(l logger.Logger) {
	// 设置日志记录器
	log = l
}

// WrapClaimsAndReq 如果做成中间件来源出去，那么直接耦合 UserClaims 也是不好的。,用于处理有请求体而且要用户验证的请求
func WrapClaimsAndReq[Req any](fn func(*gin.Context, Req, jwt.UserClaims) (Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//解析请求
		var req Req
		bind(ctx, req)
		//调用函数
		res, err := fn(ctx, req, getClaims(ctx))
		//处理结果
		respond(ctx, res, err)
	}
}

// WrapReq 。用于处理有请求体的请求
func WrapReq[Req any](fn func(*gin.Context, Req) (Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//解析参数
		var req Req
		bind(ctx, req)
		// 调用业务逻辑函数
		res, err := fn(ctx, req)
		//处理返回结果
		respond(ctx, res, err)
	}
}

// Wrap 。用于处理没有请求体的请求
func Wrap(fn func(*gin.Context) (Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := fn(ctx)
		respond(ctx, res, err)
	}
}

// WrapClaims 用于处理有用户验证但是没有请求体的请求
func WrapClaims(fn func(*gin.Context, jwt.UserClaims) (Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := fn(ctx, getClaims(ctx))
		respond(ctx, res, err)
	}
}

// 解析参数通用函数
func bind(ctx *gin.Context, req any) {
	var err error
	// 根据请求方法选择合适的绑定方式
	if ctx.Request.Method == http.MethodGet {
		err = ctx.ShouldBindQuery(&req) // 处理GET请求的查询参数
	} else {
		err = ctx.ShouldBind(&req) // 处理POST、PUT等请求的请求体数据
	}

	if err != nil {
		log.Error("解析请求失败", logger.Error(err))
	}

}

// 获取Claims通用函数
func getClaims(ctx *gin.Context) (claims jwt.UserClaims) {

	rawVal, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		log.Error("无法获得 claims",
			logger.String("path", ctx.Request.URL.Path))
		return
	}
	// 注意，这里要求放进去 ctx 的不能是*UserClaims，这是常见的一个错误
	claims, ok = rawVal.(jwt.UserClaims)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		log.Error("无法获得 claims",
			logger.String("path", ctx.Request.URL.Path))
		return
	}
	return claims
}

// 返回响应通用函数
func respond(ctx *gin.Context, res Response, err error) {
	if err != nil {
		customError := errorx.ToCustomError(err)
		//处理错误
		res.Code = customError.Code
		res.Msg = customError.Msg
		res.Data = nil
		ctx.JSON(customError.HttpCode, res)

		return
	} else {
		//vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		ctx.JSON(http.StatusOK, res)
		return
	}

}
