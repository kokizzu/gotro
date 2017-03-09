package main

import (
	"fmt"
	"github.com/kokizzu/gotro/D/Rd"
	"github.com/kokizzu/gotro/F"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/W"
	"github.com/kokizzu/gotro/X"
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

var HANDLERS = map[string]W.Action{
	``:            Home,
	`hello/:name`: Hello,
	`data`:        Data,
}

var WEBMASTER_EMAILS = M.SS{
	`test@test.com`: `Test`,
}

func Home(ctx *W.Context) {
	if ctx.IsAjax() {
		name := ctx.Posts().GetStr(`name`)
		ajax := W.Ajax{M.SX{`is_success`: true}}
		ajax.Set(`name`, name)
		ctx.AppendJson(ajax)
		return
	}
	ctx.Render(`home`, M.SX{
		`title`: `Home`,
	})
}

func Hello(ctx *W.Context) {
	ctx.Render(`hello`, M.SX{
		`title`: `Hello`,
		`name`:  ctx.ParamStr(`name`),
	})
}

func Data(ctx *W.Context) {
	params := ctx.QueryParams()
	data := M.SX{}
	params.VisitAll(func(key, value []byte) {
		data.Set(X.ToS(key), X.ToS(value))
	})
	ctx.Render(`data`, M.SX{
		`title`: `Data`,
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
	W.Routes = HANDLERS
	W.Filters = []W.Action{AuthFilter}
	runtime.GOMAXPROCS(int(L.NUM_CPU))
	// web engine
	server := W.NewEngine(DEBUG_MODE, false, PROJECT_NAME+VERSION, ROOT_DIR)
	server.StartServer(LISTEN_ADDR)
}
