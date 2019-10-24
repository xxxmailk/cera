package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

type Proxy struct {
	IP   string
	Port string
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received ruest %s %s %s\n", r.Method, r.Host, r.RemoteAddr)

	c := new(http.Client)

	uri := "http://" + net.JoinHostPort(p.IP, p.Port) + r.URL.String()
	fmt.Println("Forward to url:", r.URL.String())

	req, err := http.NewRequest(r.Method, uri, r.Body)
	if err != nil {
		fmt.Print("http.NewRequest ", err.Error())
		return
	}

	res, err := c.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for key, value := range res.Header {
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}

	w.WriteHeader(res.StatusCode)
	// 处理页面
	resource, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	result := bytes.Replace(resource, []byte("<a class=\"navbar-page-btn\" ng-click=\"ctrl.showSearch()\">"),
		[]byte("<a class=\"navbar-page-btn\" ng-click=\"ctrl.showSearch()\" style=\"display:hidden;\">"), -1)

	// 页面返回
	_, err = w.Write(result)
	if err != nil {
		fmt.Println(err.Error())
	}
	if err := res.Body.Close(); err != nil {
		fmt.Println("close request failed,", err.Error())
	}
}
