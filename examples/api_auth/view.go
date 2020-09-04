package main

import "github.com/xxxmailk/cera/view"

type Login struct {
	view.ApiView
}

func (l *Login) Get() {
	l.Data["a"] = "test"
}

type Index struct {
	view.ApiView
}

func (i *Index) Get() {
	i.Data["hello"] = "world"
}
