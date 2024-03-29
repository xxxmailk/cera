package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/valyala/fasthttp"

	"github.com/xxxmailk/cera/log"
	"github.com/xxxmailk/cera/middlewares"
	"time"
)

type CeraAuth struct {
	Username    string
	Password    string
	LoginUrl    string
	SecurityKey string
	ExpireTime  int64 // expire time: seconds
	IgnoreUrls  []string
	Resultor    CeraAuthResultor
	Log         log.SimpleLogger
	ctx         *fasthttp.RequestCtx
	middlewares.Middleware
}

type CeraAuthClaims struct {
	PoweredBy string
	jwt.StandardClaims
}

type CeraAuthResultor interface {
	JsonRs(token, iss, exp string) ([]byte, error)
}

type CeraAuthResult struct {
	Token    string `json:"Token"`
	IssuedAt string `json:"IssuedAt"`
	ExpireAt string `json:"ExpiresAt"`
}

func (c *CeraAuthResult) JsonRs(token, iss, exp string) ([]byte, error) {
	c.Token = token
	c.IssuedAt = iss
	c.ExpireAt = exp
	return json.Marshal(c)
}

func NewCeraAuth(
	username, password, loginUrl, securityKey string,
	expiredTime int64,
	resultStruct CeraAuthResultor,
	logger log.SimpleLogger,
	ignoreUrls []string) *CeraAuth {
	if loginUrl == "" {
		loginUrl = "/crea_auth/login"
	}
	c := new(CeraAuth)
	if resultStruct == nil {
		c.Resultor = new(CeraAuthResult)
	} else {
		c.Resultor = resultStruct
	}
	c.Username = username
	c.Password = password
	c.LoginUrl = loginUrl
	c.SecurityKey = securityKey
	c.ExpireTime = expiredTime
	c.IgnoreUrls = ignoreUrls
	c.Log = logger
	return c
}

type XAuthErr struct {
	Error string
}

// e.g. url: /crea_auth/login
// default: /crea_auth/login
func (a *CeraAuth) SetLoginUri(url string) {
	if url == "" {
		a.LoginUrl = "/crea_auth/login"
	}
	a.LoginUrl = url

}

func (a *CeraAuth) Handle(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	a.UnBreak()
	a.ctx = ctx
	if a.ignore() {
		a.Log.Debugf("auth ignored %s", a.ctx.URI().Path())
		return ctx
	}
	if a.isLoginUrl() {
		a.Log.Debugf("handle login url %s", a.ctx.URI().Path())
		if !a.headerAuth() && !a.paramAuth() {
			a.ctx.SetContentType("application/json")
			e, _ := json.Marshal(&XAuthErr{Error: "username or password not valid"})
			a.ctx.SetStatusCode(403)
			a.ctx.Write(e)
			a.Break()
			return ctx
		} else {
			a.login()
			return ctx
		}
	}
	if err := a.verifyToken(); err == nil {
		return ctx
	} else {
		a.Log.Debugf("login required %s method %s", a.ctx.URI().Path(), a.ctx.Method())
		e, _ := json.Marshal(&XAuthErr{Error: fmt.Sprintf("auth login required, %s", err)})
		a.ctx.SetContentType("application/json")
		a.ctx.SetStatusCode(fasthttp.StatusForbidden)
		a.ctx.Write(e)
		a.Break()
	}
	return ctx
}

func (a *CeraAuth) headerAuth() bool {
	var user, pass string
	user = string(a.ctx.Request.Header.Peek("X-Auth-Username"))
	pass = string(a.ctx.Request.Header.Peek("X-Auth-Password"))
	if user == "" {
		user = string(a.ctx.Request.Header.Peek("X-Auth-User"))
	}
	if pass == "" {
		pass = string(a.ctx.Request.Header.Peek("X-Auth-Key"))
	}
	if user == "" || pass == "" {
		return false
	}
	if a.Username == user && a.Password == pass {
		return true
	}
	return false
}

func (a *CeraAuth) paramAuth() bool {
	var user, pass string
	arg := a.ctx.PostArgs()
	user = string(arg.Peek("Username"))
	if user == "" {
		user = string(arg.Peek("username"))
	}
	pass = string(arg.Peek("Password"))
	if pass == "" {
		pass = string(arg.Peek("password"))
	}
	if user == "" || pass == "" {
		return false
	}
	if user == a.Username && pass == a.Password {
		return true
	}
	return false
}

func (a *CeraAuth) login() {
	var err error
	now := time.Now()
	cla := &CeraAuthClaims{
		PoweredBy: "Crea",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(time.Second * time.Duration(a.ExpireTime)).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    a.Username,
		},
	}
	ExpireAt := time.Unix(cla.StandardClaims.ExpiresAt, 0).Format(time.RFC3339)
	IssuedAt := time.Unix(cla.StandardClaims.IssuedAt, 0).Format(time.RFC3339)
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, cla)
	Token, err := tk.SignedString([]byte(a.SecurityKey))
	if err != nil {
		a.Log.Debugf("signed token failed %s", err)
	}
	js, err := a.Resultor.JsonRs(Token, IssuedAt, ExpireAt)
	if err != nil {
		a.Log.Errorf("marshal json result failed %s", err)
	}
	a.ctx.SetContentType("application/json")
	a.ctx.SetStatusCode(200)
	a.ctx.Write(js)
	a.Break()
}

func (a *CeraAuth) isLoginUrl() bool {
	if a.ctx.IsPost() {
		if bytes.EqualFold(a.ctx.Request.URI().Path(), []byte(a.LoginUrl)) {
			return true
		}
	}
	return false
}

func (a *CeraAuth) ignore() bool {
	for _, v := range a.IgnoreUrls {
		if bytes.EqualFold(a.ctx.Request.URI().Path(), []byte(v)) {
			return true
		}
	}
	return false
}

func (a *CeraAuth) verifyToken() error {
	tk := a.ctx.Request.Header.Peek("X-Auth-Token")
	_, err := a.verifyAction(string(tk))
	return err
}

func (a *CeraAuth) verifyAction(strToken string) (*CeraAuthClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &CeraAuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.SecurityKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*CeraAuthClaims)
	if !ok {
		return nil, errors.New("duplicated login")
	}
	if err := token.Claims.Valid(); err != nil {
		return nil, err
	}
	return claims, nil
}
