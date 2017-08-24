package W

import (
	"bytes"
	"fmt"
	"github.com/OneOfOne/cmap"
	"github.com/buaazp/fasthttprouter"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/T"
	"github.com/kokizzu/gotro/Z"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"os"
	"time"
)

type Action func(controller *Context)

type Engine struct {
	DebugMode bool
	MultiApp  bool
	// contact if there's bug
	WebMasterAnchor string
	// view template cache
	ViewCache *cmap.CMap
	// global variables
	GlobalInt M.SI
	GlobalStr M.SS
	GlobalAny M.SX
	Router    *fasthttprouter.Router
	// project_name
	Name string
	// location of the project, will be concatenated with VIEWS_SUBDIR, PUBLIC_SUBDIR
	BaseDir   string
	PublicDir string // with slash
	// server creation time
	CreatedAt time.Time
	// assets <script and <link as string
	Assets string
}

// send debug mail
func (engine *Engine) SendDebugMail(message string) {
	mail_id := `debug`
	mailer := Mailers[mail_id]
	if mailer == nil {
		return
	}
	message = `<pre>` + message + `</pre>`
	subject := engine.Name
	if engine.DebugMode {
		subject += ` DEBUG ` + T.DateStr()
	}
	mailer.SendBCC([]string{mailer.From()}, `Internal Server Error: `+subject, message)
}

// send email and ignore error
func (engine *Engine) SendMail(mail_id string, bcc []string, subject string, message string) {
	mailer := Mailers[mail_id]
	if mailer == nil {
		L.LOG.Notice(`SendMail: Invalid mailer ID: ` + mail_id)
	}
	mailer.SendBCC(bcc, `[`+engine.Name+S.If(engine.DebugMode, `-DEBUG`)+`] `+subject, message)
}

// send email and return error
func (engine *Engine) SendMailSync(mail_id string, bcc []string, subject string, message string) string {
	mailer := Mailers[mail_id]
	if mailer == nil {
		return `SyncSendMail: Invalid mailer ID: ` + mail_id
	}
	return mailer.SendSyncBCC(bcc, `[`+engine.Name+S.If(engine.DebugMode, `-DEBUG`)+`] `+subject, message)
}

// minify assets
func (engine *Engine) MinifyAssets() {
	res := bytes.Buffer{}
	if engine.DebugMode {
		// create css-link and scripts
		for _, ext := range Assets {
			var str, path string
			switch ext[0] {
			case `js`:
				path = `lib/` + ext[1] + `.` + ext[0]
				str = `<script type='text/javascript' src='/` + path + `' ></script>`
			case `css`:
				path = `lib/` + ext[1] + `.` + ext[0]
				str = `<link type='text/css' rel='stylesheet' href='/` + path + `' />`
			case `/js`:
				path = ext[1] + `.js`
				str = `<script type='text/javascript' src='/` + path + `' ></script>`
			case `/css`:
				path = ext[1] + `.css`
				str = `<link type='text/css' rel='stylesheet' href='/` + path + `' />`
			default:
				L.LOG.Notice(`Unknown resource format %v`, ext)
			}
			str += "\n	"
			if st, err := os.Stat(engine.BaseDir + PUBLIC_SUBDIR + path); err != nil || st.IsDir() {
				L.LOG.Notice(err)
			}
			res.WriteString(str)
		}
		engine.Assets = res.String()
	} else {
		dir := engine.BaseDir + PUBLIC_SUBDIR
		m := map[string]*minify.M{}
		r := map[string]*bytes.Buffer{}
		f := func(key, mime string, mini minify.MinifierFunc) {
			m[key] = minify.New()
			m[key].AddFunc(mime, mini)
			r[key] = &bytes.Buffer{}
		}
		f(`lib.css`, `text/css`, css.Minify)
		f(`mod.css`, `text/css`, css.Minify)
		f(`lib.js`, `text/js`, js.Minify)
		f(`mod.js`, `text/js`, js.Minify)
		for _, ext := range Assets {
			var path, key, mime string
			switch ext[0] {
			case `js`, `css`:
				suffix := `.` + ext[0]
				path = `lib/` + ext[1] + suffix
				key = `lib` + suffix
				mime = `text/` + ext[0]
			case `/js`, `/css`:
				suffix := `.` + ext[0][1:]
				path = ext[1] + suffix
				key = `mod` + suffix
				mime = `text` + ext[0]
			default:
				L.LOG.Notice(`Unknown resource format ` + ext[0] + ` ` + ext[1])
				continue
			}
			dat, err := ioutil.ReadFile(dir + path)
			if L.IsError(err, `failed read asset: `+path) {
				continue
			}
			dat, err = m[key].Bytes(mime, dat)
			if L.IsError(err, `failed minify asset: `+path) {
				continue
			}
			r[key].Write(dat)
			r[key].WriteRune('\n')
		}
		// please add these files to your gitignore
		for _, fname := range []string{`lib.css`, `mod.css`, `lib.js`, `mod.js`} {
			ioutil.WriteFile(dir+fname, r[fname].Bytes(), DEFAULT_FILEDIR_PERM)
			if S.EndsWith(fname, `.js`) {
				engine.Assets += `<script type='text/javascript' src='/` + fname + `' ></script>`
			} else {
				engine.Assets += `<link type='text/css' rel='stylesheet' href='/` + fname + `' />`
			}
		}
	}
}

// start the server
func (engine *Engine) StartServer(addressPort string) {
	engine.MinifyAssets()
	L.LOG.Notice(engine.Name + ` ` + S.IfElse(engine.DebugMode, `[DEVELOPMENT]`, `[PRODUCTION]`) + ` server with ` + I.ToStr(len(Routes)) + ` route(s) on ` + addressPort + "\n  Work Directory: " + engine.BaseDir)
	err := fasthttp.ListenAndServe(addressPort, engine.Router.Handler)
	L.IsError(err, `Failed to listen on `+addressPort)
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
	if Sessions == nil || Globals == nil {
		errors = append(errors, `W.Sessions not initialized, please call W.InitSession`)
	}
	return errors
}

func checkWebmasters(errors []string) []string {
	if Webmasters == nil {
		if Mailers == nil {
			errors = append(errors, `W.Webmasters not initialized`)
		} else {
			// set a random mailer as superadmin
			for _, m := range Mailers {
				Webmasters = M.SS{}
				Webmasters[m.Username] = m.Name
				break
			}
		}
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

func checkRequiredDir(errors []string, dir, const_name string, is_empty bool) []string {
	if is_empty {
		errors = append(errors, `W.`+const_name+` must not be empty`)
	}
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		if os.MkdirAll(dir, DEFAULT_FILEDIR_PERM) != nil {
			errors = append(errors, `W.`+const_name+` must be exists as a directory: `+dir)
		}
		if const_name == `PUBLIC_SUBDIR` {
			if os.MkdirAll(dir+`lib/`, DEFAULT_FILEDIR_PERM) != nil {
				errors = append(errors, `W.`+const_name+` must be exists as a directory: `+dir+`lib/`)
			}
		}
	}
	return errors
}

func checkRequiredDirs(errors []string, baseDir string) []string {
	errors = checkRequiredDir(errors, baseDir+PUBLIC_SUBDIR, `PUBLIC_SUBDIR`, PUBLIC_SUBDIR == ``)
	errors = checkRequiredDir(errors, baseDir+VIEWS_SUBDIR, `VIEWS_SUBDIR`, VIEWS_SUBDIR == ``)
	return errors
}

func checkRequiredFile(errors []string, path, label, default_content string) []string {
	check_err := func(err error) bool {
		if err != nil {
			errors = append(errors, label+` must be exists as a file (`+err.Error()+`): `+path)
			return true
		}
		return false
	}
	if st, err := os.Stat(path); err != nil || st.IsDir() {
		f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, DEFAULT_FILEDIR_PERM)
		if check_err(err) {
			return errors
		}
		_, err = f.WriteString(default_content)
		if check_err(err) {
			return errors

		}
		err = f.Close()
		check_err(err)
	}
	return errors
}

func (engine *Engine) LoadLayout() {
	tc, err := Z.ParseFile(engine.DebugMode, engine.DebugMode, engine.BaseDir+VIEWS_SUBDIR+`layout.html`)
	L.PanicIf(err, `Layout template not found`)
	engine.ViewCache.Set(`layout`, tc)
}

func (engine *Engine) Template(key string) *Z.TemplateChain {
	cache := engine.ViewCache.Get(key)
	var tc *Z.TemplateChain
	var ok bool
	tc, ok = cache.(*Z.TemplateChain)
	if cache == nil || tc == nil || !ok {
		var err error
		tc, err = Z.ParseFile(engine.DebugMode, engine.DebugMode, engine.BaseDir+VIEWS_SUBDIR+key+`.html`)
		L.PanicIf(err, `Layout template not found: `+key)
		engine.ViewCache.Set(key, tc)
	}
	return tc
}

func NewEngine(debugMode, multiApp bool, projectName, baseDir string) *Engine {
	errors := []string{}
	errors = checkMailers(errors)
	errors = checkSessions(errors)
	errors = checkWebmasters(errors)
	errors = checkRoutes(errors)
	errors = checkAssets(errors)
	errors = checkRequiredDirs(errors, baseDir)
	errors = checkRequiredFile(errors, baseDir+VIEWS_SUBDIR+`layout.html`, `layout template`, LAYOUT_DEFAULT_CONTENT)
	errors = checkRequiredFile(errors, baseDir+VIEWS_SUBDIR+`error.html`, `error template`, ERROR_DEFAULT_CONTENT)
	if len(errors) > 0 {
		panic(A.StrJoin(errors, "\n"))
	}
	engine := &Engine{
		Router:          fasthttprouter.New(),
		DebugMode:       debugMode,
		MultiApp:        multiApp,
		Name:            projectName,
		PublicDir:       baseDir + PUBLIC_SUBDIR,
		BaseDir:         baseDir,
		GlobalStr:       M.SS{},
		GlobalInt:       M.SI{},
		GlobalAny:       M.SX{},
		CreatedAt:       time.Now(),
		ViewCache:       cmap.New(),
		WebMasterAnchor: `<a href='mailto:` + Webmasters.KeysConcat(`,`) + `'>webmasters</a>`,
	}

	// handle static files
	fs := &fasthttp.FS{
		Root:               baseDir + PUBLIC_SUBDIR,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: false,
		AcceptByteRange:    true,
	}
	engine.Router.HandleMethodNotAllowed = false
	engine.Router.NotFound = fs.NewRequestHandler()

	engine.LoadLayout()

	// initialize routes handlers
	for url, handler := range Routes {
		hFun := (func(url string, handler Action) fasthttp.RequestHandler {
			return func(r_ctx *fasthttp.RequestCtx) {
				n_ctx := &Context{
					RequestCtx: r_ctx,
					Engine:     engine,
					Route:      url,
					Actions:    append([]Action{LogFilter, PanicFilter, SessionFilter}, Filters...),
				}
				if n_ctx.IsAjax() {
					n_ctx.ContentType = `application/json`
					n_ctx.NoLayout = true
				} else {
					n_ctx.ContentType = `text/html`
				}
				n_ctx.Actions = append(n_ctx.Actions, handler)
				n_ctx.Next()(n_ctx)
				n_ctx.Finish()
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

	// debug
	engine.Router.PanicHandler = func(ctx *fasthttp.RequestCtx, panic interface{}) {
		ref := `REF#` + S.RandomCB63(2)
		L.LOG.Criticalf(`%# v`, panic)
		L.LOG.Debug(ref)
		L.LOG.Debug(L.StackTrace(0))
		ctx.SetStatusCode(504)
		fmt.Fprint(ctx, `{"error":"Final Defense Block `+ref+`","is_success":false}`)
	}

	return engine
}
