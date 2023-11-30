package domain

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/X"
	"github.com/kokizzu/json5b/encoding/json5b"
	"github.com/kokizzu/lexid"
	"github.com/kpango/fastime"
	"github.com/rs/zerolog/log"

	"example2/conf"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file common.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type common.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type common.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type common.go
// go:generate msgp -tests=false -file common.go -o  common__MSG.GEN.go
//go:generate farify doublequote --file common.go

type RawFile struct {
	FileName string `json:"fileName" form:"fileName" query:"fileName" long:"fileName" msg:"fileName"`
	Mime     string `json:"mime" form:"mime" query:"mime" long:"mime" msg:"mime"`
	Size     int64  `json:"size" form:"size" query:"size" long:"size" msg:"size"`
	saveFunc func(string) error
	openFunc func() (multipart.File, error)
}

type rawFileReaderCloser struct {
	*bytes.Reader
}

func (r rawFileReaderCloser) Close() error {
	return nil
}

func NewLocalRawFileFromReader(fileName string, reader io.ReadCloser) *RawFile {
	buf := bytes.Buffer{}
	_, _ = io.Copy(&buf, reader)
	reader.Close()
	byt := rawFileReaderCloser{bytes.NewReader(buf.Bytes())}
	mime, err := mimetype.DetectReader(byt.Reader)
	byt.Reader.Seek(0, io.SeekStart)
	L.PanicIf(err, `NewLocalRawFile.mimetype.DetectReader`)
	return &RawFile{
		FileName: fileName,
		Mime:     mime.String(),
		Size:     byt.Size(),
		saveFunc: func(s string) error {
			fo, err := os.Create(s)
			if L.IsError(err, `RawFile.safeFunc.os.Create`) {
				return err
			}
			_, err = io.Copy(fo, byt)
			if L.IsError(err, `RawFile.safeFunc.io.Copy`) {
				return err
			}
			err = fo.Close()
			if L.IsError(err, `RawFile.safeFunc.fo.Close`) {
				return err
			}
			return nil
		},
		openFunc: func() (multipart.File, error) {
			return byt, nil
		},
	}
}

type RequestCommon struct {
	TracerContext context.Context `json:"-" form:"tracerContext" query:"tracerContext" long:"tracerContext" msg:"-"`
	RequestId     string          `json:"requestId,string" form:"requestId" query:"requestId" long:"requestId" msg:"requestId"`
	SessionToken  string          `json:"sessionToken" form:"sessionToken" query:"sessionToken" long:"sessionToken" msg:"sessionToken"`
	UserAgent     string          `json:"userAgent" form:"userAgent" query:"userAgent" long:"userAgent" msg:"userAgent"`
	IpAddress     string          `json:"ipAddress" form:"ipAddress" query:"ipAddress" long:"ipAddress" msg:"ipAddress"`
	OutputFormat  string          `json:"outputFormat,omitempty" form:"outputFormat" query:"outputFormat" long:"outputFormat" msg:"outputFormat"` // defaults to json
	RawFile       *RawFile        `json:"rawFile" form:"rawFile" query:"rawFile" long:"rawFile" msg:"rawFile"`
	Debug         bool            `json:"debug,omitempty" form:"debug" query:"debug" long:"debug" msg:"debug"`
	Header        string          `json:"header,omitempty" form:"header" query:"header" long:"header" msg:"header"`
	RawBody       string          `json:"rawBody,omitempty" form:"rawBody" query:"rawBody" long:"rawBody" msg:"rawBody"`
	Host          string          `json:"host" form:"host" query:"host" long:"host" msg:"host"`
	Action        string          `json:"action" form:"action" query:"action" long:"action" msg:"action"`
	Lat           float64         `json:"lat" form:"lat" query:"lat" long:"lat" msg:"lat"`
	Long          float64         `json:"long" form:"long" query:"long" long:"long" msg:"long"`
	SessionUser   Session         `json:"-" form:"-" query:"-" long:"-" msg:"-"`

	// in seconds
	now   int64     `json:"-" form:"now" query:"now" long:"now" msg:"-"`
	start time.Time `json:"-"` // for latency measurement
}

func NewLocalRequestCommon(sessionToken, userAgent string) RequestCommon {
	return RequestCommon{
		RequestId:    lexid.ID(),
		SessionToken: sessionToken,
		UserAgent:    userAgent,
		IpAddress:    `127.0.0.1`,
	}
}

func (l *RequestCommon) ToFiberCtx(ctx *fiber.Ctx, out any, rc *ResponseCommon, in any) error {
	if rc.StatusCode != http.StatusOK {
		ctx.Status(rc.StatusCode)
	}
	if rc.Redirect != `` {
		_ = ctx.Redirect(rc.Redirect, rc.StatusCode)
	}
	rc.DecorateSession(ctx)
	switch l.OutputFormat {
	case ``, `json`, fiber.MIMEApplicationJSON:
		ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		if l.Debug {
			rc.Debug = in
		}
		if l.Action == UserAutoLoginLinkAction { // to prevent / became %2f, ? became %3f
			buffer := &bytes.Buffer{}
			encoder := json.NewEncoder(buffer)
			encoder.SetEscapeHTML(false)
			err := encoder.Encode(out)
			L.Print(`AutoLoginLink: ` + buffer.String()) // TODO: remove after debugging
			if L.IsError(err, `json.Encode: %#v`, out) {
				return err
			}
			ctx.Write(buffer.Bytes())
		} else {
			byt, err := json.Marshal(out)
			if L.IsError(err, `json.Marshal: %#v`, out) {
				spew.Dump(in)
				spew.Dump(out)
				return err
			}
			_, err = ctx.Write(byt)
			if L.IsError(err, `ctx.Write failed: `+string(byt)) {
				return err
			}
			// TODO: log size/bytes written
			if l.Debug || rc.HasError() {
				L.Describe(in)
				log.Print(string(byt))
			}
		}
	case `html`:
		// do nothing
	default:
		return errors.New(`ToFiberCtx unhandled format: ` + l.OutputFormat)
	}
	return nil
}

func (l *RequestCommon) UnixNow() int64 {
	if l.now == 0 {
		l.now = fastime.UnixNow()
	}
	return l.now
}

func (i *RequestCommon) TimeNow() time.Time {
	return time.Unix(i.UnixNow(), 0)
}

func (i *RequestCommon) Latency() float64 {
	return fastime.Since(i.start).Seconds()
}

func (l *RequestCommon) FromFiberCtx(ctx *fiber.Ctx, tracerCtx context.Context) {
	l.RequestId = lexid.ID()
	l.SessionToken = ctx.Cookies(conf.CookieName, l.SessionToken)
	l.UserAgent = utils.CopyString(string(ctx.Request().Header.UserAgent()))
	l.Host = ctx.Protocol() + `://` + ctx.Hostname()
	// from nginx reverse proxy
	l.IpAddress = ctx.IP()
	if l.IpAddress == `` {
		l.IpAddress = `0.0.0.0`
	}
	// "Accept":"*/*", "Connection":"close", "Content-Length":"0", "Host":"admin.hapstr.xyz", "User-Agent":"curl/7.81.0", "X-Forwarded-For":"182.253.163.10", "X-Forwarded-Proto":"https", "X-Real-Ip":"182.253.163.10"
	l.now = fastime.UnixNow()
	l.start = fastime.Now()
	file, err := ctx.FormFile(`rawFile`)
	if err == nil {
		l.RawFile = &RawFile{
			FileName: file.Filename,
			Size:     file.Size,
			saveFunc: func(to string) error {
				return ctx.SaveFile(file, to)
			},
			openFunc: file.Open,
		}
		for _, v := range file.Header {
			l.RawFile.Mime = v[0]
			break
		}
	}
	l.TracerContext = tracerCtx
	l.SessionUser = TryDecodeSession(*l)
}

func (l *RequestCommon) ToCli(file *os.File, out any, rc ResponseCommon) {
	var byt []byte
	var err error
	switch l.OutputFormat {
	case `json`, fiber.MIMEApplicationJSON:
		byt, err = json.MarshalIndent(out, ``, `  `)
		if L.IsError(err, `json.MarshalIndent: %#v`, out) {
			return
		}
		_, err = file.Write(byt)
		if L.IsError(err, `file.Write failed: `+string(byt)) {
			return
		}
	default: // empty format also goes here
		byt, err = json5b.MarshalIndent(out, ``, `  `)
	}
	if L.IsError(err, `marshal: %#v`, out) {
		return
	}
	_, err = file.Write(byt)
	if L.IsError(err, `file.Write failed: `+string(byt)) {
		return
	}

	_, _ = os.Stderr.WriteString(M.SX{
		`statusCode`:   rc.StatusCode,
		`sessionToken`: rc.SessionToken,
		`redirect`:     rc.Redirect,
		`error`:        rc.Error,
		`debug`:        rc.Debug,
		`latency`:      l.Latency(),
		`requestId`:    l.RequestId,
	}.ToJsonPretty())

}

func (l *RequestCommon) FromCli(action string, payload []byte, in any) bool {
	err := json5b.Unmarshal(payload, &in)
	if L.IsError(err, `json5b.Unmarshal`) {
		return false
	}
	err = json.Unmarshal(payload, &in)
	if L.IsError(err, `json.Unmarshal`) {
		return false
	}
	l.RequestId = lexid.ID()
	// TODO: read from args/stdin/config-file other than json
	// l.SessionToken =
	// _, err = flags.Parse(&l)
	// L.PanicIf(err, `flags.Parse`)
	l.UserAgent = `CLI` // TODO: add input format combination, eg. json-stdin
	l.IpAddress = `127.0.0.1`
	l.TracerContext = context.Background()
	l.now = fastime.UnixNow()
	l.start = fastime.Now()
	l.Action = action
	return true
}

func (l *RequestCommon) FirstSegment() string {
	if l.Action == `` {
		return ``
	}
	segments := S.Split(l.Action, `/`)
	if len(segments) > 0 {
		return segments[0]
	}
	return ``
}

type ResponseCommon struct {
	SessionToken string `json:"sessionToken" form:"sessionToken" query:"sessionToken" long:"sessionToken" msg:"sessionToken"`
	Error        string `json:"error" form:"error" query:"error" long:"error" msg:"error"`
	StatusCode   int    `json:"status" form:"statusCode" query:"statusCode" long:"statusCode" msg:"statusCode"`
	Debug        any    `json:"debug,omitempty" form:"debug" query:"debug" long:"debug" msg:"debug"`
	Redirect     string `json:"redirect,omitempty" form:"redirect" query:"redirect" long:"redirect" msg:"redirect"`

	// action trace
	traces []any  // if you need to add trace or log that can be queried
	actor  uint64 // currently logged-in user, automatically filled by .mustLogin or .mustAdmin
	refId  uint64 // refer of primary table id being actioned
}

func (o *ResponseCommon) HasError() bool {
	return o.StatusCode >= 400 || len(o.Error) > 0
}

func (o *ResponseCommon) SetRedirect(to string) {
	o.StatusCode = 303 // force GET
	o.Redirect = to
}

func (o *ResponseCommon) SetError(code int, errStr string) {
	o.StatusCode = code
	o.Error = errStr
}

func (o *ResponseCommon) SetErrorf(code int, errFmt string, args ...any) {
	o.StatusCode = code
	o.Error = fmt.Sprintf(errFmt, args...)
}

func (l *ResponseCommon) DecorateSession(ctx *fiber.Ctx) {
	if l.SessionToken != `` {
		if l.SessionToken == conf.CookieLogoutValue {
			ctx.ClearCookie(conf.CookieName)
			ctx.Cookie(&fiber.Cookie{
				Name:    conf.CookieName,
				Expires: time.Unix(0, 0),
			})
			return
		}
		ctx.Cookie(&fiber.Cookie{
			Name:  conf.CookieName,
			Value: l.SessionToken,
			// HTTPOnly: true,
			Expires: time.Now().AddDate(0, 0, conf.CookieDays),
		})
	}
}

func (o *ResponseCommon) AddTrace(act string) {
	o.traces = append(o.traces, act)
}

func (a *ResponseCommon) AddDbChangeLogs(x []A.X) { // array of [field, old, new]
	for _, v := range x {
		a.traces = append(a.traces, v)
	}
}

func (a *ResponseCommon) AddDbChange(field, old, new any) {
	a.traces = append(a.traces, []any{field, old, new})
}

func (o *ResponseCommon) Traces() string {
	if o.traces == nil {
		return ``
	}
	return X.ToJson(o.traces)
}

func (d *Domain) segmentsFromSession(s *Session) M.SB {
	s.IsSuperAdmin = d.Superadmins[s.Email]
	s.Segments = M.SB{}
	for _, role := range s.Roles {
		switch role {
		case TenantAdminSegment:
			s.Segments[TenantAdminSegment] = true
			s.Segments[ReportViewerSegment] = true
			s.Segments[EntryUserSegment] = true
			s.Segments[UserSegment] = true
			s.Segments[GuestSegment] = true
		case EntryUserSegment:
			s.Segments[EntryUserSegment] = true
			s.Segments[UserSegment] = true
			s.Segments[GuestSegment] = true
		case ReportViewerSegment:
			s.Segments[ReportViewerSegment] = true
			s.Segments[UserSegment] = true
			s.Segments[GuestSegment] = true
		case UserSegment:
			s.Segments[GuestSegment] = true
			s.Segments[UserSegment] = true
		case GuestSegment:
			s.Segments[GuestSegment] = true
		}
	}
	if s.IsSuperAdmin {
		s.Segments[SuperAdminSegment] = true
	}
	return s.Segments
}
