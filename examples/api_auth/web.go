package main

import (
	"github.com/sirupsen/logrus"
	"github.com/xxxmailk/cera/http"
	"github.com/xxxmailk/cera/router"
)

func main() {
	r := router.New()
	r.GET("/auth/login", &Login{})
	h := http.NewHttpServe("127.0.0.1", "9999")
	h.SetLogger(logrus.New())
	h.SetRouter(r)
	//h.UseMiddleWare()
	h.Start()
}
