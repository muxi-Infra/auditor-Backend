package gin

import (
	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
}

func NewContext(c *gin.Context) *Context {
	return &Context{Context: c}
}

func (ctx *Context) Type() string {
	return "gin-audit-ctx"
}

func (ctx *Context) GetContext() *gin.Context {
	return ctx.Context
}

func (ctx *Context) BindJson(obj any) error {
	return ctx.ShouldBindJSON(obj)
}
