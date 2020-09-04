package middlewares

import "github.com/valyala/fasthttp"

type MiddlewareInterface interface {
	Handle(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx
	Break()
	UnBreak()
	IsBreakHere() bool
}

type Middleware struct {
	b bool
}

func (m *Middleware) Handle(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	return ctx
}

func (m *Middleware) Break() {
	m.b = true
}

func (m *Middleware) UnBreak() {
	m.b = false
}

func (m *Middleware) IsBreakHere() bool {
	return m.b
}
