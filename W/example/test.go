package main

import (
	"fmt"
	"github.com/kokizzu/gotro/D/Rd"
	"github.com/kokizzu/gotro/F"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/W"
	"github.com/kokizzu/gotro/X"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"
)

var VERSION string
var DEBUG_MODE = (VERSION == ``)
var LISTEN_ADDR = `:9090`
var ROOT_DIR string

var ASSETS = [][2]string{
	//// http://api.jquery.com/ 1.11.1
	{`js`, `jquery`},
}

var ROUTERS = map[string]W.Action{
	``:                            LoginExample,
	`login_example`:               LoginExample,
	`logout_example`:              LogoutExample,
	`post_values_example`:         PostValuesExample,
	`named_params_example/:value`: NamedParamsExample,
	`query_string_example`:        QueryStringExample,
}

var WEBMASTER_EMAILS = M.SS{
	`test@test.com`: `Test`,
}

func LoginExample(ctx *W.Context) {
	user_id := ctx.Session.GetInt(`user_id`)
	if ctx.IsAjax() {
		ajax := W.Ajax{M.SX{`is_success`: true}}
		username := ctx.Posts().GetStr(`username`)
		password := ctx.Posts().GetStr(`password`)
		if username != password {
			ajax.Set(`is_success`, false)
			ajax.Error(`301 Wrong username or password; username is case sensitive`)
			ctx.AppendJson(ajax)
			return
		}
		id := rand.Intn(1000)
		ctx.Session.Login(M.SX{
			`username`: username,
			`user_id`:  id,
			`level`:    M.SX{},
		})
		ajax.Set(`logged`, id)
		ctx.AppendJson(ajax)
		return
	}
	ctx.Render(`login_example`, M.SX{
		`title`:   `Login Example`,
		`user_id`: user_id,
	})
}

func LogoutExample(ctx *W.Context) {
	ajax := W.Ajax{M.SX{`is_success`: true}}
	ctx.Session.Logout()
	ctx.AppendJson(ajax)
}

func PostValuesExample(ctx *W.Context) {
	if ctx.IsAjax() {
		ajax := W.Ajax{M.SX{`is_success`: true}}
		value := ctx.Posts().GetStr(`value`)
		ajax.Set(`value`, value)
		ctx.AppendJson(ajax)
		return
	}
	ctx.Render(`post_values_example`, M.SX{
		`title`: `Post Values Example`,
	})
}

func NamedParamsExample(ctx *W.Context) {
	ctx.Render(`named_params_example`, M.SX{
		`title`: `Named Params Example`,
		`value`: ctx.ParamStr(`value`),
	})
}

func QueryStringExample(ctx *W.Context) {
	params := ctx.QueryParams()
	data := M.SX{}
	params.VisitAll(func(key, value []byte) {
		data.Set(X.ToS(key), X.ToS(value))
	})
	ctx.Render(`query_string_example`, M.SX{
		`title`: `Query String Example`,
		`data`:  data,
	})
}

const PROJECT_NAME = `Gotro Example`
const DOMAIN = `localhost`

// filter the page that may or may may not be accessed
func AuthFilter(ctx *W.Context) {
	L.Trace()
	handled := false
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
	redis_conn := Rd.NewRedisSession(``, ``, 9, `session::`)
	W.InitSession(`Test`, 2*24*time.Hour, 1*24*time.Hour, *redis_conn)
	W.Mailers = map[string]*W.SmtpConfig{
		``: {
			Name:     `Mailer Daemon`,
			Username: `test.test`,
			Password: `123456`,
			Hostname: `smtp.gmail.com`,
			Port:     587,
		},
	}
	W.Assets = ASSETS
	W.Webmasters = WEBMASTER_EMAILS
	W.Routes = ROUTERS
	W.Filters = []W.Action{AuthFilter}
	runtime.GOMAXPROCS(int(L.NUM_CPU))
	// web engine
	server := W.NewEngine(DEBUG_MODE, false, PROJECT_NAME+VERSION, ROOT_DIR)
	server.StartServer(LISTEN_ADDR)
}
