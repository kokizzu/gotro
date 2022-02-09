package main

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/W2/example/conf"
	"github.com/kokizzu/gotro/W2/example/domain"
	"github.com/kokizzu/overseer"
	"github.com/kokizzu/overseer/fetcher"
)

var requiredHeader = M.SS{
	//domain.SomeUrl: `X-CC-Webhook-Signature`,
}

// purpose: convert HTTP request into HTTP response thru domain functions
// purpose: convert websocket message into websocket response thru domain functions

func NewWebApi() {
	overseer.Run(overseer.Config{
		Program: webApiServer(),
		Address: conf.WEBAPI_HOSTPORT,
		Fetcher: &fetcher.File{
			Path:     conf.WEBAPI_EXEPATH,
			Interval: 1 * time.Second,
		},
		//Debug: conf.DEBUG_MODE,
		// if lots of "[overseer master] proxy signal (urgent I/O condition)" shown
		// then the WEBAPI_EXEPATH is wrong
	})
}

func webApiServer() func(state overseer.State) {
	return func(state overseer.State) {
		log.Info().Str("state", state.ID).Str(`listen`, conf.WEBAPI_HOSTPORT)
		app := fiber.New()
		app.Use(logger.New())
		app.Use(recover.New())
		app.Use(cors.New()) // allow from any host
		app.Get(`/`, func(ctx *fiber.Ctx) error {
			_, _ = ctx.WriteString(`ok`)
			return nil
		})
		//seedInitialData()
		domain := webApiInitRoutes(app)
		webApiInitGraphql(app, domain)
		runCron(domain)
		L.Print(conf.AdminTestSessionToken)
		err := app.Listener(state.Listener)
		L.IsError(err, `app.Listener failed`)
	}
}

// priority:
// 1. query string
// 2. body
// 3. params
func webApiParseInput(ctx *fiber.Ctx, reqCommon *domain.RequestCommon, in interface{}, url string) error {
	body := ctx.Body()
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
				return err
			}
		}
		trimBody := S.Left(string(body), 4096)
		L.Print(trimBody)
		if reqCommon.Debug && reqCommon.RawBody == `` {
			reqCommon.RawBody = trimBody
		}
	}
	return nil
}
