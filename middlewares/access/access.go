package access

import (
	"github.com/valyala/fasthttp"
	"github.com/xxxmailk/cera/log"
	"github.com/xxxmailk/cera/middlewares"
	"strings"
)

type Access struct {
	Log log.SimpleLogger
	middlewares.Middleware
}

func NewAccessMiddleware(logger log.SimpleLogger) *Access {
	return &Access{Log: logger}
}

func (a *Access) Handle(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	a.Log.Infof("access url %s method %s from %s agent %s request body %d bytes %v bytes sent",
		ctx.Path(),
		ctx.Request.Header.Method(),
		getRealIP(ctx),
		ctx.UserAgent(),
		ctx.Request.Header.ContentLength(),
		ctx.Response.Header.ContentLength())
	return ctx
}

func getRealIP(ctx *fasthttp.RequestCtx) string {
	clientIP := string(ctx.Request.Header.Peek("X-Forwarded-For"))
	if clientIP == "" {
		return ctx.Conn().RemoteAddr().String()
	}
	if index := strings.IndexByte(clientIP, ','); index >= 0 {
		clientIP = clientIP[0:index]
		//获取最开始的一个 即 1.1.1.1
	}
	clientIP = strings.TrimSpace(clientIP)
	if len(clientIP) > 0 {
		return clientIP
	}
	clientIP = strings.TrimSpace(string(ctx.Request.Header.Peek("X-Real-Ip")))
	if len(clientIP) > 0 {
		return clientIP
	}
	return ""
}
