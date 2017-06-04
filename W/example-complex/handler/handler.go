package handler

import (
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/W"
	"github.com/kokizzu/gotro/W/example-complex/model"
	"github.com/kokizzu/gotro/W/example-complex/model/mLogin"
	"github.com/kokizzu/gotro/W/example-complex/model/mResponse"
)

const NOT_REQUIRED = `not required`

func Home(ctx *W.Context) {
	rm := mResponse.PrepareVars(ctx, `Home`)
	L.Print(`masuk1`)
	if !rm.Ok {
		return
	}
	L.Print(`masuk2`)
	if rm.IsAjax() {
		// handle ajax
		//switch rm.Action {
		//case `something`: // @API
		//default: // @API-END
		ErrorHandler(rm.Ajax, rm.Action)
		//}
		ctx.AppendJson(rm.Ajax.SX)
		return
	}
	locals := W.Ajax{SX: M.SX{
		`title`: ctx.Title,
	}}
	ctx.Render(`home`, locals.SX)

}

func ErrorHandler(ajax W.Ajax, a string) {
	ajax.Error(model.ERR_402_INVALID_ACTION + a)
}

func Logout(ctx *W.Context) {
	switch NOT_REQUIRED {
	case NOT_REQUIRED: // @API
		mLogin.API_All_Logout(ctx)
	default: // @API-END
	}
}

// 2016-07-25 Prayogo
func Login(ctx *W.Context) {
	if ctx.IsAjax() {
		switch NOT_REQUIRED {
		case NOT_REQUIRED: // @API
			mLogin.API_All_Login(ctx)
		default: // @API-END
		}
	}
	ctx.Title = `Login`
	locals := M.SX{
		`title`:      ctx.Title,
		`google_url`: ``,
		`fb_url`:     ``,
		`fb_id`:      mLogin.FB_APPID,
		`csrf`:       ctx.Session.StateCSRF(),
	}
	if g_provider := mLogin.GetGPlusOAuth(ctx); g_provider != nil {
		locals[`google_url`] = g_provider.AuthCodeURL(ctx.Session.StateCSRF())
	}
	if f_provider := mLogin.GetFBOAuth(ctx); f_provider != nil {
		locals[`fb_url`] = f_provider.AuthCodeURL(ctx.Session.StateCSRF())
	}
	ctx.Render(`login/index`, locals)
}

// 2015-06-22 Prayogo
func Login_Forgot(ctx *W.Context) {
	if !ctx.IsAjax() {
		ctx.Title = `Forgot Password`
		ctx.Render(`login/forgot`, M.SX{`support_email`: model.SUPPORT_EMAIL})
		return
	}
	switch NOT_REQUIRED {
	case NOT_REQUIRED: // @API
		mLogin.API_All_LoginForgot(ctx)
	default: // @API-END
	}
}

// 2015-06-22 Prayogo
func Login_Reset(ctx *W.Context) {
	if !ctx.IsAjax() {
		ctx.Title = `Reset Password`
		ctx.Render(`login/reset`, M.SX{})
		return
	}
	switch NOT_REQUIRED {
	case NOT_REQUIRED: // @API
		mLogin.API_All_LoginReset(ctx)
	default: // @API-END
	}
}

// 2017-02-06 Prayogo
func Login_Verify(ctx *W.Context) {
	switch NOT_REQUIRED {
	case NOT_REQUIRED: // @API
		mLogin.API_All_VerifyOAuth(ctx)
	default: // @API-END
	}
}
