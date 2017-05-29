package main

import (
	"fmt"
	"github.com/kokizzu/gotro/D/Rd"
	"github.com/kokizzu/gotro/F"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/W"
	"os"
	"runtime"
	"strings"
	"time"
)

var VERSION string
var DEBUG_MODE = (VERSION == ``)
var LISTEN_ADDR = `:3001`
var ROOT_DIR string

var WEBMASTER_EMAILS = M.SS{
	`CHANGEME@gmail.com`: `CHANGEME Name`,
}

const PROJECT_NAME = `CHANGEME Example`

// filter the page that may or may may not be accessed
func AuthFilter(ctx *W.Context) {
	L.Trace()
	handled := false
	if ctx.Session.GetStr(`user_id`) > `` {
		// logged in
	} else {
		//ctx.Error(403, ``)
		// you can block the page for non-logged-in users here (handled=true)
	}
	ctx.Session.Touch()
	if !handled {
		cpu := L.PercentCPU()
		if cpu > 95.0 {
			W.Sessions.Inc(`throttle_counter`)
			if !ctx.IsAjax() {
				ctx.Error(503, `Server Overloaded`)
			} else {
				ctx.AppendString(`{"errors":["error 503: server is overloaded, please wait for a moment.."]}`)
			}
			fmt.Println(`Throttled: ` + F.ToS(cpu) + ` %`)
			return
		}
		ctx.Next()(ctx)
	}
}

// initialize loading time
func init() {
	var err error
	ROOT_DIR, err = os.Getwd() // filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		_, path, _, _ := runtime.Caller(0)
		slash_pos := strings.LastIndex(path, `/`) + 1
		ROOT_DIR = path[:slash_pos]
	} else {
		ROOT_DIR += `/`
	}
}

func main() {
	redis_conn := Rd.NewRedisSession(``, ``, 9, `session::`)
	global_conn := Rd.NewRedisSession(``, ``, 10, `session::`)
	W.InitSession(`Aaa`, 2*24*time.Hour, 1*24*time.Hour, *redis_conn, *global_conn)
	W.Mailers = map[string]*W.SmtpConfig{
		``: {
			Name:     `Mailer Daemon`,
			Username: `CHANGEME`,
			Password: `CHANGEME`,
			Hostname: `smtp.gmail.com`,
			Port:     587,
		},
	}
	W.Assets = ASSETS
	W.Webmasters = WEBMASTER_EMAILS
	W.Routes = ROUTERS
	W.Filters = []W.Action{AuthFilter}
	// web engine

	server := W.NewEngine(DEBUG_MODE, false, PROJECT_NAME+VERSION, ROOT_DIR)
	server.StartServer(LISTEN_ADDR)
}
