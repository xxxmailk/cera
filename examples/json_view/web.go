package main

import (
	"github.com/sirupsen/logrus"
	"github.com/xxxmailk/cera/http"
	"github.com/xxxmailk/cera/middlewares/access"
	"github.com/xxxmailk/cera/middlewares/auth"
	"github.com/xxxmailk/cera/router"
)

func main() {
	r := router.New()
	r.GET("/auth/login", &Login{})
	r.POST("/auth/login", &Login{})
	r.ANY("/", &Paas{})
	h := http.NewHttpServe("127.0.0.1", "9999")
	logger := logrus.New()
	logger.SetLevel(logrus.TraceLevel)
	logger.Formatter = &logrus.TextFormatter{}
	logger.ReportCaller = true
	h.SetLogger(logger)
	au := auth.NewCeraAuth("root", "P@ssw0rd", "/login", "haha", 300, &auth.CeraAuthResult{}, logger, []string{"/aa"})
	h.UseMiddleWare(au)
	acc := access.NewAccessMiddleware(logger)
	h.AtLast(acc)
	h.SetRouter(r)
	h.Start()
}
