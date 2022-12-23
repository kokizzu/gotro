package domain

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/jessevdk/go-flags"
	"github.com/kokizzu/id64"
	"github.com/kpango/fastime"

	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/W2/example/conf"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file common.go
//go:generate replacer 'Id" form' 'Id,string" form' type common.go
//go:generate replacer 'json:"id"' 'json:"id,string"' type common.go
//go:generate replacer 'By" form' 'By,string" form' type common.go
// go:generate msgp -tests=false -file common.go -o  common__MSG.GEN.go
//go:generate farify doublequote --file common.go

type RequestCommon struct {
	TracerContext context.Context   `json:"-" form:"tracerContext" query:"tracerContext" long:"tracerContext" msg:"-"`
	RequestId     uint64            `json:"requestId,string" form:"requestId" query:"requestId" long:"requestId" msg:"requestId"`
	SessionToken  string            `json:"sessionToken" form:"sessionToken" query:"sessionToken" long:"sessionToken" msg:"sessionToken"`
	UserAgent     string            `json:"userAgent" form:"userAgent" query:"userAgent" long:"userAgent" msg:"userAgent"`
	IpAddress     string            `json:"ipAddress" form:"ipAddress" query:"ipAddress" long:"ipAddress" msg:"ipAddress"`
	OutputFormat  string            `json:"outputFormat,omitempty" form:"outputFormat" query:"outputFormat" long:"outputFormat" msg:"outputFormat"` // defaults to json
	Uploads       map[string]string `json:"uploads,omitempty" form:"uploads" query:"uploads" long:"uploads" msg:"uploads"`                          // form key and temporary file path
	Now           int64             `json:"-" form:"now" query:"now" long:"now" msg:"-"`
	Debug         bool              `json:"debug,omitempty" form:"debug" query:"debug" long:"debug" msg:"debug"`
	Header        string            `json:"header,omitempty" form:"header" query:"header" long:"header" msg:"header"`
	RawBody       string            `json:"rawBody,omitempty" form:"rawBody" query:"rawBody" long:"rawBody" msg:"rawBody"`
	Host          string            `json:"host" form:"host" query:"host" long:"host" msg:"host"`
}

func (l *RequestCommon) ToFiberCtx(ctx *fiber.Ctx, out any) error {
	defer l.deleteTempFiles()
	switch l.OutputFormat {
	case ``, `json`, fiber.MIMEApplicationJSON:
		ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		byt, err := json.Marshal(out)
		//L.Print(string(byt))
		if L.IsError(err, `json.Marshal: %#v`, out) {
			return err
		}
		_, err = ctx.Write(byt)
		if L.IsError(err, `ctx.Write failed: `+string(byt)) {
			return err
		}
		// TODO: log size/bytes written
	default:
		return errors.New(`ToFiberCtx unhandled format: ` + l.OutputFormat)
	}
	return nil
}

func (l *RequestCommon) UnixNow() int64 {
	if l.Now == 0 {
		l.Now = fastime.UnixNow()
	}
	return l.Now
}

func (i *RequestCommon) TimeNow() time.Time {
	return time.Unix(i.UnixNow(), 0)
}

func (l *RequestCommon) FromFiberCtx(ctx *fiber.Ctx, tracerCtx context.Context) {
	l.RequestId = id64.UID()
	l.SessionToken = ctx.Cookies(conf.CookieName, l.SessionToken)
	l.UserAgent = string(ctx.Request().Header.UserAgent())
	l.IpAddress = ctx.IP()
	l.Host = ctx.Protocol() + `://` + ctx.Hostname()
	l.Now = fastime.UnixNow()
	// TODO: check ctx.Accepts()
	//ctx.Request().URI().RequestURI()
	file, err := ctx.FormFile(`fileBinary`)
	if err == nil {
		l.Uploads = map[string]string{}
		target := `/tmp/` + id64.SID() + filepath.Ext(file.Filename)
		err = ctx.SaveFile(file, target)
		if !L.IsError(err, `failed saving upload to: `+target) {
			l.Uploads[file.Filename] = target
		}
	}
	//L.Describe(l)
	l.TracerContext = tracerCtx
}

func (l *RequestCommon) ToCli(file *os.File, out any) {
	defer l.deleteTempFiles()
	switch l.OutputFormat {
	case ``, `json`, fiber.MIMEApplicationJSON:
		byt, err := json.Marshal(out)
		if L.IsError(err, `json.Marshal: %#v`, out) {
			return
		}
		_, err = file.Write(byt)
		if L.IsError(err, `file.Write failed: `+string(byt)) {
			return
		}
		// TODO: log size/bytes written
	default:
		L.Print(`ToCli unhandled format: ` + l.OutputFormat)
	}
}

func (l *RequestCommon) FromCli(file *os.File, ctx context.Context) {
	l.RequestId = id64.UID()
	// TODO: read from args/stdin/config-file?
	// l.SessionToken =
	_, err := flags.Parse(&l)
	L.PanicIf(err, `flags.Parse`)
	l.UserAgent = `CLI` // TODO: add input format combination, eg. json-stdin
	l.IpAddress = `127.0.0.1`
	l.TracerContext = ctx
	l.Now = fastime.UnixNow()
}

func (l *RequestCommon) deleteTempFiles() {
	for _, tmpPath := range l.Uploads {
		// TODO: delete temporary uploads
		_ = tmpPath
	}
}

type ResponseCommon struct {
	SessionToken string `json:"sessionToken" form:"sessionToken" query:"sessionToken" long:"sessionToken" msg:"sessionToken"`
	Error        string `json:"error" form:"error" query:"error" long:"error" msg:"error"`
	StatusCode   int    `json:"status" form:"statusCode" query:"statusCode" long:"statusCode" msg:"statusCode"`
	Debug        any    `json:"debug,omitempty" form:"debug" query:"debug" long:"debug" msg:"debug"`
	Redirect     string `json:"redirect,omitempty" form:"redirect" query:"redirect" long:"redirect" msg:"redirect"`
}

func (o *ResponseCommon) HasError() bool {
	return o.StatusCode >= 400 || len(o.Error) > 0
}

func (o *ResponseCommon) SetRedirect(to string) {
	o.StatusCode = 302
	o.Redirect = to
}

func (o *ResponseCommon) SetError(code int, errStr string) {
	o.StatusCode = code
	o.Error = errStr
}

func (l *ResponseCommon) ToFiberCtx(ctx *fiber.Ctx, inRc *RequestCommon, in any) {
	if l.SessionToken != `` {
		if l.SessionToken == conf.CookieLogoutValue {
			ctx.ClearCookie(conf.CookieName)
		} else {
			ctx.Cookie(&fiber.Cookie{
				Name:    conf.CookieName,
				Value:   l.SessionToken,
				Expires: time.Now().AddDate(0, 0, conf.CookieDays),
			})
		}
	}
	if l.StatusCode > 0 {
		res := ctx.Response()
		res.SetStatusCode(l.StatusCode)
		if l.Redirect != `` {
			err := ctx.Redirect(l.Redirect)
			L.IsError(err, `ResponseCommon.ToFiberCtx.Redirect failed: `+l.Redirect)
		}
	}
	if inRc.Debug {
		inRc.TracerContext = nil
		L.Describe(in)
		l.Debug = in
	}
}

func (l *ResponseCommon) ToCli(file *os.File) {
	// TODO: write to stdout/config-file?
}

// requestCommon struct generator

func ExpiredRC() RequestCommon {
	return RequestCommon{SessionToken: conf.AdminTestExpiredSession}
}
func AdminRC() RequestCommon {
	return RequestCommon{SessionToken: conf.AdminTestSessionToken}
}
func EmptyRC() RequestCommon {
	return RequestCommon{SessionToken: conf.AdminTestSessionToken}
}
func NewRC(sessionToken string) RequestCommon {
	return RequestCommon{SessionToken: sessionToken}
}
