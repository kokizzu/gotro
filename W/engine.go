package W

import (
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/valyala/fasthttp"
)

type Action func(controller *Context)

type Engine struct {
	DebugMode       bool
	MultiApp        bool
	WebMasterAnchor M.SX
	GlobalInt       M.SI
	GlobalStr       M.SS
	GlobalAny       M.SX
	Router          *Router
	Name            string
	BaseDir         string
	StaticDir       string
}

func (engine *Engine) SyncSendMail(mail_id string, bcc []string, subject string, message string) string {
	// TODO: continue this
	return ``
}
func (engine *Engine) StartServer(addressPort string) {
	// TODO: continue this
	// engine.MinifyAssets()
	msg := `[DEVELOPMENT]`
	if !engine.DebugMode {
		msg = `[PRODUCTION]`
	}
	_ = msg
}

// attach a middleware on non-static files
func (engine *Engine) Use(m Action) {
	//engine.Filters = append(engine.Filters, m)
	// TODO: continue this
}

func checkMailers(errors []string) []string {
	if Mailers == nil {
		errors = append(errors, `W.Mailers not initialized`)
	}
	if Mailers[``] == nil && Mailers[`debug`] != nil {
		Mailers[``] = Mailers[`debug`]
	}
	if Mailers[`debug`] == nil && Mailers[``] != nil {
		Mailers[`debug`] = Mailers[``]
	}
	if Mailers[``] == nil || Mailers[`debug`] == nil {
		errors = append(errors, `W.Mailers[""] and W.Mailers["debug"] not initialized`)
	}
	return errors
}

func checkSessions(errors []string) []string {
	if Sessions == nil {
		errors = append(errors, `W.Sessions not initialized, please call W.InitSession`)
	}
	return errors
}

func checkWebmasters(errors []string) []string {
	if Webmasters == nil {
		errors = append(errors, `W.Webmasters not initialized`)
	}
	return errors
}

func checkRoutes(errors []string) []string {
	if Routes == nil {
		errors = append(errors, `W.Routes not initialized`)
	}
	return errors
}

func checkAssets(errors []string) []string {
	if Assets == nil {
		errors = append(errors, `W.Assets not initialized`)
	}
	return errors
}

func NewEngine(debugMode, multiApp bool, projectName, baseDir, staticDir string) *Engine {
	errors := []string{}
	errors = checkMailers(errors)
	errors = checkSessions(errors)
	errors = checkWebmasters(errors)
	errors = checkRoutes(errors)
	errors = checkAssets(errors)
	if len(errors) > 0 {
		panic(A.StrJoin(errors, "\n"))
	}
	engine := &Engine{
		Router:    New(),
		DebugMode: debugMode,
		MultiApp:  multiApp,
		Name:      projectName,
		BaseDir:   baseDir,
		StaticDir: staticDir,
	}
	//engine := &Engine{
	//	Router:          &Router{},
	//	Filters:         []Action{},
	//	pageNotFound:    Error404,
	//	ajaxNotFound:    Ajax404,
	//	DebugMode:       debugMode,
	//	BaseDir:         baseDir,
	//	PublicDir:       baseDir + `public/`,
	//	ViewDir:         baseDir + `views/`,
	//	ViewCache:       cmap.New(),
	//	Name:            projectName,
	//	WebMaster:       webMaster,
	//	WebMasterAnchor: `<a href='mailto:` + webMaster.KeysConcat(`,`) + `'>webmasters</a>`,
	//	GlobalInt:       M.SI{},
	//	GlobalStr:       M.SS{},
	//	GlobalAny:       M.SX{},
	//	BaseUrls:        baseUrls,
	//	MultiApp:        multiApp,
	//	CreatedAt:       time.Now(),
	//	Assets:          assets,
	//}

	//engine.LoadLayout()

	// initialize engine filters and static file handling
	//engine.FileServer = http.FileServer(http.Dir(engine.PublicDir))
	//return engine

	// initialize routes handlers
	for url, handler := range Routes {
		hFun := (func(url string, handler Action) fasthttp.RequestHandler {
			return func(r_ctx *fasthttp.RequestCtx) {
				n_ctx := &Context{
					RequestCtx: r_ctx,
					Engine:     engine,
					Route:      url,
					Actions:    append([]Action{}, Filters...),
				}
				n_ctx.Actions = append(n_ctx.Actions, handler)
				n_ctx.Next()(n_ctx)
			}
		})(url, handler)
		engine.Router.GET(`/`+url, hFun)
		engine.Router.POST(`/`+url, hFun)
	}
	if len(engine.GlobalInt) > 0 {
		L.Describe(engine.GlobalInt)
	}
	if len(engine.GlobalStr) > 0 {
		L.Describe(engine.GlobalStr)
	}
	if len(engine.GlobalAny) > 0 {
		L.Describe(engine.GlobalAny)
	}
	L.Describe(engine.Name+S.If(engine.DebugMode, ` DEBUG Version`), I.ToStr(len(Routes))+` route(s)`, baseDir)
	return engine
}
