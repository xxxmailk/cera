package main

import "github.com/xxxmailk/cera/view"

type Login struct {
	view.View
}

func (l *Login) Get() {
	l.Data["a"] = "test"
	l.Tpl = "login"
}
