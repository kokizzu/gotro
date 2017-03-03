package W

import (
	"bytes"
	"encoding/json"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/valyala/fasthttp"
)

var strPost = []byte(`POST`)

type Context struct {
	fasthttp.RequestCtx
	Session *Session
	Title   string
	Engine  *Engine
	Params  *Params
	Buffer  bytes.Buffer
}

// check if request method is POST
func (ctx *Context) IsAjax() bool {
	method := ctx.Request.Header.Method()
	return bytes.Compare(method, strPost) == 0
}

// get requested host
func (ctx *Context) Host() string {
	return string(ctx.RequestCtx.Host())
}

// append bytes
func (ctx *Context) AppendBytes(buf []byte) {
	L.Trace()
	ctx.Buffer.Write(buf)
}

// append buffer
func (ctx *Context) AppendBuffer(buff bytes.Buffer) {
	ctx.Buffer.Write(buff.Bytes())
}

// append string
func (ctx *Context) AppendString(txt string) {
	L.Trace()
	ctx.Buffer.WriteString(txt)
}

// append json
func (ctx *Context) AppendJson(any Ajax) {
	L.Trace()
	buf, err := json.Marshal(any.SX)
	L.IsError(err, `error converting to json`)
	ctx.Buffer.Write(buf)
}

func (ctx *Context) Render(path string, locals M.SX) {
	// TODO: continue this: use Z template engine
}

func (ctx *Context) PartialNoDebug(path string, locals M.SX) string {
	// TODO: continue this: return rendered partial as string
	return ``
}

// TODO: test this, make sure it returns only first segment
func (ctx *Context) FirstPath() string {
	uri := ctx.Request.RequestURI()
	first_slash := bytes.IndexByte(uri[1:], '/')
	if first_slash < 0 {
		return string(uri[1:])
	}
	return string(uri[1:first_slash])
}

func (ctx *Context) IsWebMaster() bool {
	// TODO: continue this
	return true
}

func (ctx *Context) Posts() *Posts {
	// TODO: fix this, should return correctly
	return &Posts{}
}

func (ctx *Context) Error(text string, any ...interface{}) {
	// TODO: continue this
}

func (ctx *Context) ServeError(code int64) {
	// TODO: copy from W.ServeError
}

func (ctx *Context) Next() Action {
	// TODO: continue this
	return func(ctx *Context) {

	}
}
