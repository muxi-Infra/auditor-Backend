package gin

type MiddlewareFunc func(next HandlerFunc) HandlerFunc

type Chain struct {
	middlewares []MiddlewareFunc
}

func NewChain(middlewares ...MiddlewareFunc) *Chain {
	return &Chain{middlewares: middlewares}
}

func (c *Chain) Use(m ...MiddlewareFunc) {
	c.middlewares = append(c.middlewares, m...)
}

func (c *Chain) Wrap(h HandlerFunc) HandlerFunc {
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		h = c.middlewares[i](h)
	}
	return h
}
