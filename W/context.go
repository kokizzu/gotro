package W

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/T"
	"github.com/kokizzu/gotro/X"
	"github.com/valyala/fasthttp"
	"mime/multipart"
	"net/http"
	"path/filepath"
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
	return `http` + S.If(ctx.RequestCtx.IsTLS(), `s`) + `://` + string(ctx.RequestCtx.Host())
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
func (ctx *Context) AppendJson(any M.SX) {
	buf, err := json.Marshal(any)
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
		fmt.Fprint(ctx, ctx.Buffer.String())
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
		fmt.Fprint(ctx, buff.String())
	}
}

// TODO: test this, make sure it returns only first segment
func (ctx *Context) FirstPath() string {
	uri := ctx.Request.RequestURI()
	first_slash := bytes.IndexByte(uri[1:], '/')
	if first_slash < 0 {
		return string(uri[1:])
	}
	return string(uri[1 : first_slash+1])
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

// get parsed ?a=b&c=d, this is different from Param*() func
func (ctx *Context) QueryParams() *QueryParams {
	return &QueryParams{ctx.RequestCtx.QueryArgs()}
}

// request URL
func (ctx *Context) RequestURL() string {
	return string(ctx.RequestURI())
}

// uploaded file
func (ctx *Context) UploadedFile(id string) (fileName, ext, contentType string, reader multipart.File) {
	fh, err := ctx.RequestCtx.FormFile(id)
	if L.IsError(err, `Parameter multipart.FileHeader missing:`+id) {
		ext = err.Error()
		return
	}
	buff := make([]byte, 512)
	reader, err = fh.Open()
	if L.IsError(err, `Opening multipart.FileHeader: `+id) {
		return
	}
	_, err = reader.Read(buff)
	if L.IsError(err, `Reading first 512-byte multipart.FileHeader.Reader: `+id) {
		ext = err.Error()
		reader.Close()
		return
	}
	contentType = http.DetectContentType(buff) // do not trust header.Get(`content-type`)[0]
	reader.Seek(0, 0)
	fileName = fh.Filename
	if fileName == `` {
		fileName = `tmp`
	}
	ext = filepath.Ext(fileName)
	ext = S.ToLower(ext)
	return
}

// debug info
func (ctx *Context) RequestDebugStr() string {
	return ctx.Session.IpAddr + S.WebBR +
		ctx.Session.UserAgent + S.WebBR +
		ctx.Title + S.WebBR +
		T.DateTimeStr() + S.WebBR +
		S.IfElse(ctx.IsAjax(), `POST`, `GET`) + ` ` + S.IfEmpty(string(ctx.Path()), `/`) + S.WebBR +
		`Session: ` + ctx.Session.String() + S.WebBR +
		ctx.Posts().String() + S.WebBR + S.WebBR
}

// non debug info
func (ctx *Context) RequestHtmlStr() string {
	return `Request Path: ` + ctx.RequestURL() + S.WebBR +
		`Server Time: ` + T.DateTimeStr() + S.WebBR +
		`Session: ` + ctx.Session.String() + S.WebBR
}
