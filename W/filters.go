package W

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"time"
)

func PanicFilter(ctx *Context) {
	defer func() {
		if err := recover(); err != nil {
			//L.Panic(errors.New(`Internal Server Error`), ``, err)
			err2, ok := err.(error)
			if !ok {
				err2 = errors.New(fmt.Sprintf("%# v", err))
			}
			err_str := err2.Error()
			L.LOG.Errorf(err_str)
			str := L.StackTrace(2)
			L.LOG.Criticalf("StackTrace: %s", str)
			// TODO: continue this
			//ctx.Abort(500, ctx.Engine.DebugMode, err2)
			if !ctx.Engine.DebugMode {
				// ctx.Engine.SendDebugMail(ctx.RequestDebugStr() + "\n\n" + err_str)
			}
		}
	}()
	ctx.Next()(ctx)
}

func LogFilter(ctx *Context) {
	start := time.Now()
	ctx.Next()(ctx)
	L.Trace()
	var codeStr string
	code := ctx.Response.StatusCode()
	if code < 400 {
		codeStr = L.BgGreen(`%s`, color.BlueString(`%3d`, code))
	} else {
		codeStr = L.BgRed(`%3d`, code)
	}
	ones := `nano`
	elapsed := float64(time.Since(start).Nanoseconds())
	if elapsed > 1000000000.0 {
		elapsed /= 1000000000.0
		ones = `sec`
	} else if elapsed > 1000000.0 {
		elapsed /= 1000000.0
		ones = `mili`
	} else if elapsed > 1000.0 {
		elapsed /= 1000.0
		ones = `micro`
	}
	referrer := ctx.RequestCtx.Referer()
	url := ctx.RequestCtx.RequestURI()
	msg := fmt.Sprintf(`[%s] %7d %7.2f %5s | %4s %-40s | %-40s | %15s | %s || %s`,
		codeStr,
		ctx.Buffer.Len(),
		elapsed,
		ones,
		ctx.RequestCtx.Method(),
		url,
		referrer,
		ctx.Session.IpAddr,
		ctx.Session.String(),
		ctx.Posts().String(),
	)
	msg = S.Replace(msg, `%`, `%%`)
	L.LOG.Notice(msg)
}

func SessionFilter(ctx *Context) {
	ctx.Session = &Session{}
	ctx.Session.Load(ctx)
	ctx.Next()(ctx)
	ctx.Session.Save(ctx)
}
