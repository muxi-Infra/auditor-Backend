package middleware

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	api_errors "muxi_auditor/api/errors"
	"muxi_auditor/api/response"
	"muxi_auditor/config"
	"muxi_auditor/ioc"
	"muxi_auditor/pkg/errorx"
	"muxi_auditor/pkg/ginx"
	"muxi_auditor/pkg/jwt"
	"muxi_auditor/pkg/logger"
	"net/http"
	"time"
)

type CorsMiddleware struct {
	allowedOrigins []string
}

type AuthMiddleware struct {
	jwtHandler *jwt.RedisJWTHandler
}

type LoggerMiddleware struct {
	log        logger.Logger
	prometheus *ioc.Prometheus
}

func NewCorsMiddleware(conf *config.MiddlewareConf) *CorsMiddleware {
	return &CorsMiddleware{allowedOrigins: conf.AllowedOrigins}
}

func (cm *CorsMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 允许的请求头
		AllowHeaders: []string{"Content-ContentType", "Authorization", "Origin"},
		// 是否允许携带凭证（如 Cookies）
		AllowCredentials: true,
		// 解决跨域问题,这个地方允许所有请求跨域了,之后要改成允许前端的请求,比如localhost
		AllowOriginFunc: func(origin string) bool {
			//暂时允许所有跨域请求,根据需要进行调整
			return true
			//只允许在列表里面的origin可以跨域
			//if slices.Contains(cm.allowedOrigins, origin) {
			//	return true
			//} else {
			//	return false
			//}
		},

		// 预检请求的缓存时间
		MaxAge: 12 * time.Hour,
	})
}

func NewAuthMiddleware(jwtHandler *jwt.RedisJWTHandler) *AuthMiddleware {
	return &AuthMiddleware{jwtHandler: jwtHandler}
}

func (am *AuthMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求中提取并解析 Token
		userClaims, err := am.jwtHandler.ParseToken(ctx)
		if err != nil {
			ctx.Error(api_errors.UNAUTHORIED_ERROR(err))
			return
		}
		// 将解析后的用户信息存入上下文，供后续逻辑使用
		ginx.SetClaims[jwt.UserClaims](ctx, userClaims)

		// 继续处理请求
		ctx.Next()
	}
}

func NewLoggerMiddleware(
	log logger.Logger,
	prometheus *ioc.Prometheus,
) *LoggerMiddleware {
	return &LoggerMiddleware{
		log:        log,
		prometheus: prometheus,
	}
}

func (lm *LoggerMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.FullPath()

		// 记录活跃连接数
		lm.prometheus.ActiveConnections.WithLabelValues(path).Inc()
		defer func() {

			// 记录响应信息
			lm.prometheus.ActiveConnections.WithLabelValues(path).Dec()
			status := ctx.Writer.Status()
			lm.prometheus.RouterCounter.WithLabelValues(ctx.Request.Method, path, http.StatusText(status)).Inc()
			lm.prometheus.DurationTime.WithLabelValues(path, http.StatusText(status)).Observe(time.Since(start).Seconds())
		}()

		ctx.Next() // 执行后续逻辑

		// 处理返回值或错误
		res, httpCode := lm.handleResponse(ctx)
		if !ctx.IsAborted() { // 避免重复返回响应
			ctx.JSON(httpCode, res)
		}
	}
}

// 提取的日志逻辑：记录自定义错误日志
func (lm *LoggerMiddleware) logCustomError(customError *errorx.CustomError, ctx *gin.Context) {
	lm.log.Error("处理请求出错",
		logger.Error(customError),
		logger.String("timestamp", time.Now().Format(time.RFC3339)),
		logger.String("ip", ctx.ClientIP()),
		logger.String("path", ctx.Request.URL.Path),
		logger.String("method", ctx.Request.Method),
		logger.String("headers", fmt.Sprintf("%v", ctx.Request.Header)),
		logger.Int("httpCode", customError.HttpCode),
		logger.Int("code", customError.Code),
		logger.String("msg", customError.Msg),
		logger.String("category", customError.Category),
		logger.String("file", customError.File),
		logger.Int("line", customError.Line),
		logger.String("function", customError.Function),
	)
}

// 提取的日志逻辑：记录未知错误日志
func (lm *LoggerMiddleware) logUnexpectedError(err error, ctx *gin.Context) {
	lm.log.Error("意外错误类型",
		logger.Error(err),
		logger.String("timestamp", time.Now().Format(time.RFC3339)),
		logger.String("ip", ctx.ClientIP()),
		logger.String("path", ctx.Request.URL.Path),
		logger.String("method", ctx.Request.Method),
		logger.String("headers", fmt.Sprintf("%v", ctx.Request.Header)),
	)
}
func (lm *LoggerMiddleware) commonInfo(ctx *gin.Context) {
	lm.log.Info("意外错误类型",
		logger.String("timestamp", time.Now().Format(time.RFC3339)),
		logger.String("ip", ctx.ClientIP()),
		logger.String("path", ctx.Request.URL.Path),
		logger.String("method", ctx.Request.Method),
		logger.String("headers", fmt.Sprintf("%v", ctx.Request.Header)),
	)
}

// 处理响应逻辑
func (lm *LoggerMiddleware) handleResponse(ctx *gin.Context) (response.Response, int) {
	var res response.Response
	httpCode := http.StatusOK

	//有错误则进行错误处理
	if len(ctx.Errors) > 0 {
		err := ctx.Errors.Last().Err
		customError := errorx.ToCustomError(err)
		if customError == nil {
			lm.logUnexpectedError(err, ctx)
			return response.Response{Code: api_errors.ERROR_TYPE_ERROR_CODE, Msg: err.Error(), Data: nil}, http.StatusInternalServerError
		}
		lm.logCustomError(customError, ctx)
		return response.Response{Code: customError.Code, Msg: customError.Msg, Data: nil}, customError.HttpCode
	} else {

		//无错误则记录常规日志
		lm.commonInfo(ctx)
		res = ginx.GetResp[response.Response](ctx)
	}

	return res, httpCode
}
