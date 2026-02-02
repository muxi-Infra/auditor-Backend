package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/muxi-Infra/auditor-Backend/sdk/v2/api/errorx"
	"github.com/muxi-Infra/auditor-Backend/sdk/v2/api/request"
	"github.com/muxi-Infra/auditor-Backend/sdk/v2/api/response"
)

type Server struct {
	group *gin.RouterGroup
}

func NewGinRegistrar(g *gin.RouterGroup) *Server {
	return &Server{group: g}
}

func (gr *Server) WebHook(path string, chain *Chain, fn func(g *gin.Context, Req *request.HookPayload) (response.Resp, error)) {
	gr.group.POST(path, func(c *gin.Context) {
		ctx := NewContext(c)
		handler := chain.Wrap(HandlerWrapper(fn))

		res, err := handler(ctx)

		if err != nil {
			// 只记录 error，不写 JSON
			if len(c.Errors) == 0 {
				c.Error(err)
			}
			return
		}

		// 只 Set，不写 JSON
		c.Set("sdk.result", res)
	})
}

func HandlerWrapper[Req any, Resp any](fn func(*gin.Context, Req) (Resp, error)) HandlerFunc {
	return func(c *Context) (any, error) {
		var req Req
		err := c.BindJson(&req)
		if err != nil {
			return nil, err
		}

		resp, err := fn(c.GetContext(), req)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// SDKResponseMiddleware SDK 兜底响应中间件（一定要在用户中间件之后）
func SDKResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// 先判断用户自定义的插件是否已经写入
		if c.Writer.Written() {
			return
		}
		// todo:这里的error系统需要系统的设计一下，暂且先这样
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			c.JSON(400, response.Resp{
				Code: errorx.CustomErrCode,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		val, ok := c.Get("sdk.result")
		if !ok {
			return
		}

		res, ok := val.(response.Resp)
		if !ok {
			c.JSON(500, response.Resp{
				Code: errorx.SDKResponseErrCode,
				Msg:  "invalid response type",
				Data: nil,
			})
			return
		}
		c.JSON(200, res)
	}
}
