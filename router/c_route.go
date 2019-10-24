package router

import (
	"github.com/valyala/fasthttp"
)

const (
	GET = iota
	POST
	HEAD
	PUT
	DELETE
	CONNECT
	OPTIONS
	TRACE
	PATCH
	COPY
)

type IRouterHandler interface {
	Handler(ctx *fasthttp.RequestCtx)
}

func NewRouter() *Router {
	r := new(Router)
	r.RedirectFixedPath = true
	r.RedirectTrailingSlash = true
	r.HandleMethodNotAllowed = true
	r.HandleOPTIONS = true
	return r
}
