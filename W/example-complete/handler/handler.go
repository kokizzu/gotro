package handler

import (
	"example-complete/sql"
	"example-complete/sql/sLogin"
	"example-complete/sql/sResponse"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/W"
)

const NOT_REQUIRED = `not required`

func Home(ctx *W.Context) {
	rm := sResponse.Prepare(ctx, `Home`, false)
	_ = rm
	values := M.SX{
		`has_login`: ctx.Session.GetStr(`id`) != ``,
	}
	ctx.Render(`home`, values)
}

func Logout(ctx *W.Context) {
	switch NOT_REQUIRED {
	case NOT_REQUIRED: // @API
		sLogin.API_All_Logout(ctx)
	default: // @API-END
	}
}

// 2016-07-25 Prayogo
func Login(ctx *W.Context) {
	if ctx.IsAjax() {
		switch NOT_REQUIRED {
		case NOT_REQUIRED: // @API
			sLogin.API_All_Login(ctx)
		default: // @API-END
		}
		return
	}
	ctx.Title = `Login`
	locals := M.SX{
		`title`:      ctx.Title,
		`google_url`: ``,
		`fb_url`:     ``,
		`csrf`:       ctx.Session.StateCSRF(),
	}
	if g_provider := sLogin.GetGPlusOAuth(ctx); g_provider != nil {
		locals[`google_url`] = g_provider.AuthCodeURL(ctx.Session.StateCSRF())
	}
	ctx.Render(`login/index`, locals)
}

// 2017-02-06 Prayogo
func Login_Verify(ctx *W.Context) {
	switch NOT_REQUIRED {
	case NOT_REQUIRED: // @API
		sLogin.API_All_VerifyOAuth(ctx)
	default: // @API-END
	}
}

func ErrorHandler(ajax W.Ajax, a string) {
	ajax.Error(sql.ERR_402_NO_ACTION + a)
}
