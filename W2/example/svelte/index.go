package svelte

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/bytebufferpool"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/W2/example/conf"
	"github.com/kokizzu/gotro/W2/example/domain"
	"github.com/kokizzu/gotro/Z"
)

var Domain *domain.Domain

var routes = map[string]*Z.TemplateChain{}

func createHandler(tc *Z.TemplateChain, path string, handler func(*fiber.Ctx, string) (bind M.SX, err error)) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		buf := bytebufferpool.Get()
		defer bytebufferpool.Put(buf)
		m, err := handler(c, path)
		tc.Render(buf, m)
		c.Set(`content-type`, `text/html; charset=utf-8`)
		_ = c.SendString(buf.String())
		return err
	}
}

func HandleSSR(app *fiber.App, viewDir string, domain *domain.Domain) {
	Domain = domain
	debug := conf.VERSION == ``
	for path, handler := range handlers {
		fileName := viewDir + `/` + path + `.html`
		tc, err := Z.ParseFile(debug, debug, fileName)
		L.IsError(err, path)
		routes[`/`+path] = tc
		handler := createHandler(tc, path, handler)
		app.Get(`/`+path, handler)
		app.Get(`/`+path+`.html`, handler)
		if S.Right(path, 5) == `index` {
			app.Get(`/`+S.LeftOfLast(path, `index`), handler)
			if S.Right(path, 6) == `/index` {
				app.Get(`/`+S.LeftOfLast(path, `/index`), handler)
			}
		}
	}
}

func index(ctx *fiber.Ctx, path string) (bind M.SX, err error) {
	// can do queries from Domain
	bind = M.SX{
		`title`: `this is generated title from server: ` + path,
		`arr`:   A.X{`1abc`, 2, 3},
		`obj`: M.SX{
			`a`: 1,
			`b`: 2.345,
			`c`: 3,
		},
	}
	return
}

func page1_subpage(c *fiber.Ctx, path string) (bind M.SX, err error) {
	bind = M.SX{
		`something`: M.SX{
			`a`: 1,
			`b`: path,
		},
		`title`: path,
	}
	return
}
func page1_subpage3_index(c *fiber.Ctx, path string) (bind M.SX, err error) {
	bind = M.SX{
		`from_server`: A.X{
			A.X{path, 2},
			A.X{`a`, 3},
		},
		`from_server2`: `whoa`,
		`title`:        path,
	}
	return
}

func page2_index(c *fiber.Ctx, path string) (bind M.SX, err error) {
	bind = M.SX{
		`title`:        `also from server`,
		`from_server2`: path,
		`from_server`:  A.X{1, 2, 3},
	}
	return
}

func page1_index(c *fiber.Ctx, path string) (bind M.SX, err error) {
	bind = M.SX{
		`title`: path,
	}
	return
}
