package presentation

import (
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/Z"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"example2/conf"
	"example2/domain"
)

type WebServer struct {
	*domain.Domain
	Cfg conf.WebConf
}

var requiredHeader = M.SS{
	//domain.SomeUrl: `X-CC-Webhook-Signature`,
}

// priority:
// 1. query string
// 2. body
// 3. params
func webApiParseInput(ctx *fiber.Ctx, reqCommon *domain.RequestCommon, in any, url string) error {
	body := ctx.Body()
	reqCommon.Action = url
	reqCommon.Debug = reqCommon.Debug || conf.IsDebug()
	path := S.LeftOf(url, `?`) // without API_PREFIX
	if header, ok := requiredHeader[path]; ok {
		reqCommon.Header = ctx.Get(header)
		reqCommon.RawBody = string(body)
	}
	//L.Print(ctx.OriginalURL())
	if err := ctx.QueryParser(in); L.IsError(err, `ctx.QueryParser failed: `+url) {
		return err
	}
	if len(body) > 0 {
		retry := true
		if body[0] == '{' || ctx.Get(`content-type`) == `application/json` {
			if err := json.Unmarshal(body, in); err == nil {
				retry = false
			}
		}
		// application/x-www-form-urlencoded
		// multipart/form-data
		if retry {
			if err := ctx.BodyParser(in); L.IsError(err, `ctx.BodyParser failed: `+url) {
				_ = ctx.JSON(M.SX{
					`error`: err.Error(),
				})
				return err
			}
		}
		trimBody := S.Left(string(body), 1024)
		if reqCommon.Debug && reqCommon.RawBody == `` {
			reqCommon.RawBody = trimBody
		}
	}
	if conf.IsDebug() && reqCommon.Debug && reqCommon.RawBody != `` {
		// prevent too large dump because multipart/form-data is raw binary
		contentType := ctx.Get(`content-type`)
		if !S.StartsWith(contentType, `multipart/form-data`) {
			//log.Print(ctx.GetReqHeaders())
			log.Print(reqCommon.RawBody)
		}
	}
	reqCommon.FromFiberCtx(ctx, ctx.UserContext())
	return nil
}

func (w *WebServer) Start(log *zerolog.Logger) {
	fw := fiber.New(fiber.Config{
		ProxyHeader: `X-Real-IP`,
	})

	// check if actionLogs are there, if error, then you need to run migration: go run main.go migrate
	w.InsertActionLog(&domain.RequestCommon{
		UserAgent: "server",
		IpAddress: "127.0.0.1",
		Action:    "server/start",
	}, &domain.ResponseCommon{})

	// load svelte templates
	views = &Views{}
	views.LoadAll()

	fw.Use(limiter.New(limiter.Config{
		Max:               300,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))
	if conf.IsDebug() { // TODO: use faster logger for production
		fw.Use(logger.New())
	} else { // prevent panic on production
		copy := recover.ConfigDefault
		copy.EnableStackTrace = true
		fw.Use(recover.New(copy))
	}

	// assign static routes (GET)
	w.WebStatic(fw, w.Domain)

	// API routes (POST)
	ApiRoutes(fw, w.Domain)

	//zImport.GoogleSheetCountryDataToJson("1TmAjrclFHUwDA1487ifQjX4FzYt9y7eJ0gwyxtwZMJU")

	// finally serve static files from svelte directory if any
	fw.Static(`/`, `./svelte`, fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        false,
		CacheDuration: 5 * time.Second,
		MaxAge:        3600,
	})

	log.Err(fw.Listen(w.Cfg.ListenAddr()))
}

type Views struct {
	cache map[string]*Z.TemplateChain
}

var views *Views

func (v *Views) LoadAll() {
	if v.cache == nil {
		v.cache = map[string]*Z.TemplateChain{}
	}
	debug := conf.IsDebug()
	const svelteDir = `svelte/`
	var err error
	for svelte, html := range viewList {
		v.cache[svelte], err = Z.ParseFile(debug, debug, svelteDir+html)
		L.PanicIf(err, `failed to parse `+html)
	}
}
