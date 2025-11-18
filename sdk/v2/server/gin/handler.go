package gin

type HandlerFunc func(c *Context) (any, error)
