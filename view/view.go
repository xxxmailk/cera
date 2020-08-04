package view

import (
	"crypto/sha1"
	"encoding/binary"
	"github.com/valyala/fasthttp"
	"github.com/xxxmailk/cera/log"
	"golang.org/x/net/xsrftoken"
	"html/template"
	"math/rand"
	"strconv"
	"time"
)

const CERASALT = "Crea@2019=="

var (
	HTMLContentType = []byte("text/plain; charset=utf-8")
	JSONContentType = []byte("application/json; charset=utf-8")
)

func init() {
	rand.Seed(time.Now().Unix())
	sha1.New()
}

type MethodViewer interface {
	viewer
	Init()
	Get()
	Post()
	Head()
	Options()
	Put()
	Patch()
	Delete()
	Trace()
	Render()
	SetLogger(log.SimpleLogger)
}

type viewer interface {
	Before()
	After()
	GetCtx() *fasthttp.RequestCtx
	SetCtx(ctx *fasthttp.RequestCtx)
}

type CreaCookie struct {
	ActionId [16]byte // action id of session,generate by timestamp and random uint
	XsrfKey  [8]byte  // xsrf key
	XsrfUid  [8]byte  // xsrf uid
}

// parse cookie to struct
func ParseCookie(ck []byte) *CreaCookie {
	c := new(CreaCookie)
	copy(c.ActionId[:], ck[0:15])
	copy(c.XsrfKey[:], ck[16:23])
	copy(c.XsrfUid[:], ck[24:31])
	return c
}

// generate struct to byte slice
func (c *CreaCookie) ToByte() []byte {
	b := make([]byte, 32)
	copy(b[0:15], c.ActionId[:])
	copy(b[16:23], c.XsrfKey[:])
	copy(b[24:31], c.XsrfUid[:])
	return b
}

// generate new action id
func newActionId() [16]byte {
	buf := [16]byte{}
	r := make([]byte, 8)
	tm := make([]byte, 8)
	// covert int64 to byte
	binary.BigEndian.PutUint64(r, rand.Uint64())
	binary.BigEndian.PutUint64(tm, uint64(time.Now().Unix()))
	// add random int and timestamp to buffer
	copy(buf[0:7], r[:])
	copy(buf[8:15], tm[:])
	return buf
}

type View struct {
	Tpl    string                 // template name
	Data   map[string]interface{} // stored user values
	Ctx    *fasthttp.RequestCtx
	Cookie *fasthttp.Cookie
	Logger log.SimpleLogger
}

// combine this struct and rewrite those functions to reply http methods
func (r *View) Init() {
	r.Data = make(map[string]interface{})
}

func (r *View) Before() {}

func (r *View) Get() {
	if err := r.Html404(); err != nil {
		r.Logger.Errorf("handle 404 with GET method error %s", err)
	}
}

func (r *View) Head() {
	if err := r.Html404(); err != nil {
		r.Logger.Errorf("handle 404 with Head method error %s", err)
	}
}

func (r *View) Options() {
	if err := r.Html404(); err != nil {
		r.Logger.Errorf("handle 404 with Option method error %s", err)
	}
}

func (r *View) Post() {
	if err := r.Html404(); err != nil {
		r.Logger.Errorf("handle 404 with Post method error %s", err)
	}
}

func (r *View) Put() {
	if err := r.Html404(); err != nil {
		r.Logger.Errorf("handle 404 with Put method error %s", err)
	}
}

func (r *View) Patch() {
	if err := r.Html404(); err != nil {
		r.Logger.Errorf("handle 404 with Patch method error %s", err)
	}
}

func (r *View) Delete() {
	if err := r.Html404(); err != nil {
		r.Logger.Errorf("handle 404 with Delete method error %s", err)
	}
}

func (r *View) Trace() {
	if err := r.Html404(); err != nil {
		r.Logger.Errorf("handle 404 with Trace method error %s", err)
	}
}

func (r *View) SetLogger(l log.SimpleLogger) {
	r.Logger = l
}

func (r *View) After() {}

func (r *View) Render() {
	t := template.Must(template.ParseGlob("./template/*.htm"))
	r.Ctx.Response.Header.SetContentType("text/html; charset=utf-8")
	err := t.ExecuteTemplate(r.Ctx.Response.BodyWriter(), r.Tpl, r.Data)
	if err != nil {
		r.Logger.Errorf("render template failed, %s", err)
		return
	}
}

// 获取参数，通过标准get url方式传值 e.g. http://xxx.com/?id=1
func (r *View) GetArgString(key string) string {
	return string(r.Ctx.Request.URI().QueryArgs().Peek(key))
}

// 获取参数，通过标准get url方式传值 e.g. http://xxx.com/?id=1
func (r *View) GetArgBytes(key string) []byte {
	return r.Ctx.Request.URI().QueryArgs().Peek(key)
}

// 获取参数，通过标准get url方式传值 e.g. http://xxx.com/?id=1
func (r *View) GetArgInt(key string) (int, error) {
	s := r.GetArgString(key)
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (r *View) GetCtx() *fasthttp.RequestCtx {
	return r.Ctx
}

func (r *View) SetCtx(ctx *fasthttp.RequestCtx) {
	r.Ctx = ctx
}

func Switcher(v MethodViewer) {
	// running before method priority
	ctx := v.GetCtx()
	v.Before()
	method := ctx.Method()
	switch string(method) {
	case fasthttp.MethodGet:
		v.Get()
		v.Render()
	case fasthttp.MethodPost:
		v.Post()
		v.Render()
	case fasthttp.MethodHead:
		v.Head()
		v.Render()
	case fasthttp.MethodOptions:
		v.Options()
		v.Render()
	case fasthttp.MethodPut:
		v.Put()
		v.Render()
	case fasthttp.MethodPatch:
		v.Patch()
		v.Render()
	case fasthttp.MethodDelete:
		v.Delete()
		v.Render()
	case fasthttp.MethodTrace:
		v.Trace()
		v.Render()
	default:
		HtmlUnknownMethod(ctx)
	}
	v.After()
}

func (r *View) GetPostArgs(key string) string {

	return string(r.Ctx.PostArgs().Peek(key))
}

func (r *View) setXsrfToken() {
	buf := make([]byte, 8)
	bufU := make([]byte, 8)
	random := rand.Uint64()
	randomU := rand.Uint64()
	binary.BigEndian.PutUint64(buf, random)
	binary.BigEndian.PutUint64(bufU, randomU)
	xsrf := xsrftoken.Generate(string(buf), strconv.FormatUint(rand.Uint64(), 16), strconv.FormatInt(time.Now().Unix(), 16))

	c := new(fasthttp.Cookie)
	c.SetKey("CreaCookie")
	// todo : set session key id
	//c.SetValue()
	r.Data["XSRF"] = xsrf
	r.Ctx.Response.Header.SetCookie(c)
}

// xsrf token check, if xsrf token is not valid, response permission denied to client
//func XsrfValidate(ctx *fasthttp.RequestCtx) {
//	// xsrf token
//	xsrf := string(ctx.Request.PostArgs("XSRF_TOKEN"))
//	// xsrf uid
//	xsu := string(r.Ctx.Request.Header.Cookie("xsu"))
//	// xsrf key
//	xsk := string(r.Ctx.Request.Header.Cookie("xsk"))
//	// xsrf action
//	xsa := string(r.Ctx.Request.Header.Cookie("xsa"))
//	if !xsrftoken.Valid(xsrf, xsk, xsu, string(xsa)) {
//		html := `
//<html>
//<head>
//<title>Permission denied!xsrf</title>
//</head>
//<body style="background:#000;text-align:center;">
//<span style="font-size:5em;color:#fff;"><b>403, sorry, xsrf token check failed! :) </b></span>
//</body>
//</html>
//`
//		r.Ctx.SetStatusCode(403)
//		if _, err := r.Ctx.Write([]byte(html)); err != nil {
//			log.Print(err)
//		}
//		r.Ctx.Done()
//	}
//}

func HtmlUnknownMethod(ctx *fasthttp.RequestCtx) error {
	html := `
<html>
<head>
<title>Page not found</title>
</head>
<body style="background:#000;text-align:center;">
<span style="font-size:5em;color:#fff;"><b>500, sorry! unknown http method :) </b></span>
</body>
</html>
`
	ctx.SetStatusCode(500)
	if _, err := ctx.Write([]byte(html)); err != nil {
		return err
	}
	return nil
}

func (r *View) Html404() error {
	html := `
<body style="background:#000;text-align:center;">
<span style="font-size:5em;color:#fff;"><b>404 Sorry, Page not found! :) </b></span>
</body>
`
	r.Ctx.SetStatusCode(404)
	r.Ctx.SetBodyString(html)
	if _, err := r.Ctx.Write([]byte(html)); err != nil {
		return err
	}
	return nil
}
