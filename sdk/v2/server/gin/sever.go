package gin

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/v2/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/v2/api/response"
	"github.com/gin-gonic/gin"
)

type Server struct {
	group *gin.RouterGroup
}

func NewGinRegistrar(g *gin.RouterGroup) *Server {
	return &Server{group: g}
}

func (gr *Server) WebHook(path string, chain *Chain, fn func(g *gin.Context, Req *request.HookPayload) (response.Resp, error)) {
	gr.group.POST(path, func(c *gin.Context) {
		wrapped := chain.Wrap(HandlerWrapper(fn))
		ctx := NewContext(c)
		data, err := wrapped(ctx)

		if err != nil {
			c.JSON(ctx.Writer.Status(), response.Resp{
				Code: 0,
				Msg:  err.Error(),
				Data: nil,
			})
		}

		c.JSON(ctx.Writer.Status(), response.Resp{
			Code: 200,
			Msg:  "success",
			Data: data,
		})
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
