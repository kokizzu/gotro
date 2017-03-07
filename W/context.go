package W

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/X"
	"github.com/valyala/fasthttp"
)

var strPost = []byte(`POST`)

type Context struct {
	*fasthttp.RequestCtx
	Session     *Session
	Title       string
	Engine      *Engine
	Buffer      bytes.Buffer
	Route       string
	Actions     []Action
	PostCache   *Posts
	NoLayout    bool
	ContentType string
}

// get url parameter as string
func (ctx *Context) ParamStr(key string) string {
	return X.ToS(ctx.RequestCtx.UserValue(key))
}

// get url parameter as integer
func (ctx *Context) ParamInt(key string) int64 {
	return X.ToI(ctx.RequestCtx.UserValue(key))
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
	ctx.Buffer.Write(buf)
}

// append buffer
func (ctx *Context) AppendBuffer(buff bytes.Buffer) {
	ctx.Buffer.Write(buff.Bytes())
}

// append string
func (ctx *Context) AppendString(txt string) {
	ctx.Buffer.WriteString(txt)
}

// append json
func (ctx *Context) AppendJson(any Ajax) {
	buf, err := json.Marshal(any.SX)
	L.IsError(err, `error converting to json`)
	ctx.Buffer.Write(buf)
}

func (ctx *Context) Render(path string, locals M.SX) {
	ctx.Engine.Template(path).Render(&ctx.Buffer, locals)
}

func (ctx *Context) PartialNoDebug(path string, locals M.SX) string {
	buff := bytes.Buffer{}
	ctx.Engine.Template(path).Render(&buff, locals)
	return buff.String()
}

func (ctx *Context) Finish() {
	ctx.SetContentType(ctx.ContentType)
	if ctx.NoLayout {
		fmt.Fprint(ctx, ctx.Buffer)
	} else {
		buff := bytes.Buffer{}
		ctx.Engine.Template(`layout`).Render(&buff, M.SX{
			`title`:         ctx.Title,
			`project_name`:  ctx.Engine.Name,
			`assets`:        ctx.Engine.Assets,
			`contents`:      ctx.Buffer.String(),
			`is_superadmin`: ctx.IsWebMaster(),
			`debug_mode`:    ctx.Engine.DebugMode,
		})
		fmt.Fprint(ctx, buff)
	}
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
	return Webmasters[ctx.Session.GetStr(`email`)] != ``
}

func (ctx *Context) Posts() *Posts {
	if ctx.PostCache == nil {
		p := &Posts{}
		p.FromContext(ctx)
		ctx.PostCache = p
	}
	return ctx.PostCache
}

// call next filter or action/handler
func (ctx *Context) Next() Action {
	if ctx.Actions == nil || len(ctx.Actions) == 0 {
		panic(`action-chain unavailable`)
	}
	action := ctx.Actions[0]
	ctx.Actions = ctx.Actions[1:]
	return action
}

func (ctx *Context) Error(code int, info string) {
	ctx.SetStatusCode(code)
	ctx.Render(`error`, M.SX{
		`requested_path`: ctx.Request.URI().String(),
		`error_code`:     code,
		`project_name`:   ctx.Engine.Name,
		`webmaster`:      ctx.Engine.WebMasterAnchor,
		`error_title`:    Errors[code],
		`error_detail`:   info,
	})
}
