package main

import (
	"fmt"
	"github.com/kokizzu/gotro/D/Pg"
	"github.com/kokizzu/gotro/D/Rd"
	"github.com/kokizzu/gotro/F"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/W"
	//"github.com/pkg/profile"
	"example-complete/sql"
	"os"
	"runtime"
	"strings"
	"time"
)

var VERSION string
var DEBUG_MODE = (VERSION == ``)
var LISTEN_ADDR = `:3001`
var ROOT_DIR string
var PROJECT_NAME string
var DOMAIN string

// filter the page that may or may may not be accessed
func AuthFilter(ctx *W.Context) {
	L.Trace()
	handled := false
	//Log.Describe(ctx.Session)
	//L.Describe(ctx.Session.RedisKey())
	if ctx.Session.GetInt(`user_id`) > 0 {
		// logged in
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
	sql.PROJECT_NAME = PROJECT_NAME
	sql.DOMAIN = DOMAIN
	//defer profile.Start(profile.CPUProfile).Stop()
	redis_conn := Rd.NewRedisSession(``, ``, 9, `session::`)
	W.InitSession(sql.PROJECT_NAME, 7*24*time.Hour, 1*24*time.Hour, *redis_conn, *redis_conn)
	W.Mailers = map[string]*W.SmtpConfig{
		``: {
			Name:     `Mailer Daemon`,
			Username: `TODO_CHANGE_THIS`,
			Password: `TODO_CHANGE_THIS`,
			Hostname: `smtp.gmail.com`,
			Port:     587,
		},
	}
	W.Assets = ASSETS
	W.Webmasters = sql.WEBMASTER_EMAILS
	W.Routes = HANDLERS
	W.Filters = []W.Action{AuthFilter}
	Pg.InitOfficeMail(`@` + sql.DOMAIN)
	//Pg.DEBUG = DEBUG_MODE
	// web engine
	server := W.NewEngine(DEBUG_MODE, false, sql.PROJECT_NAME+VERSION, ROOT_DIR)
	server.StartServer(LISTEN_ADDR)
}
