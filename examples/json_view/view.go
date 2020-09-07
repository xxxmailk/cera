package main

import "github.com/xxxmailk/cera/view"

type Login struct {
	view.ApiView
}

func (l *Login) Get() {
	l.Data["a"] = "test"
}
func (l *Login) Post() {
	post := l.Ctx.Request.PostArgs()
	get := l.Ctx.Request.URI().QueryArgs()
	l.Data["a"] = "test"
	l.Data["port"] = post
	l.Data["get"] = get
}

type Paas struct {
	view.ApiView
}

func (p *Paas) Get() {
	p.Ctx.Response.Header.SetStatusCode(600)
	p.Data["test"] = "aaa"
}
