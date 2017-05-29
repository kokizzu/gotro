package handler

import (
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/W"
	"github.com/kokizzu/gotro/W/example-complex/model"
	"github.com/kokizzu/gotro/W/example-complex/model/mLogin"
	"gitlab.com/kokizzu/gokil/L"
)

const NOT_REQUIRED = `not required`

func Home(ctx *W.Context) {
	rm := PrepareVars(ctx, `Home`)
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

func NewAjaxResponse() W.Ajax {
	return W.Ajax{SX: M.SX{`is_success`: true}}
}

func PrepareVars(ctx *W.Context, title string) *W.RequestModel {
	user_id := ctx.Session.GetStr(`id`)
	rm := &W.RequestModel{
		Actor:   user_id,
		DbActor: user_id,
		Level:   ctx.Session.SX,
		Ctx:     ctx,
	}
	ctx.Title = title
	is_ajax := ctx.IsAjax()
	if is_ajax {
		rm.Ajax = NewAjaxResponse()
	}
	//page := rm.Level.GetMSB(`page`)
	//first_segment := ctx.FirstPath()
	// validate if this user may access this first segment
	// check their access level, if it's not ok, set rm.Ok to false
	// then render an error, something like this:
	/*
	 if is_ajax {
	  rm.Ajax.Error(sql.ERR_403_MUST_LOGIN_HIGHER)
	  ctx.AppendJson(rm.Ajax.SX)
	  return
	 }
	 ctx.Error(403, sql.ERR_403_MUST_LOGIN_HIGHER)
	 return
	*/
	if !is_ajax {
		// render menu based on privilege
		ctx.Title = title
		menus := []string{
			// add here for more 1st segment
			`backoffice`,
		}
		email := ctx.Session.GetStr(`email`)
		logger := `Not Logged In`
		if email != `` {
			logger = email + `: `
		}
		full_name := rm.Level.GetStr(`full_name`)
		if ctx.IsWebMaster() {
			logger += ` <a class='menu' href='/system_admin/impersonate/` + rm.Actor + `'>` + full_name + ` (` + rm.Actor + `)</a>`
		} else {
			logger += full_name
		}
		values := M.SX{
			`title`:      title,
			`email`:      logger,
			`uid`:        user_id,
			`debug_mode`: S.IfElse(ctx.Engine.DebugMode, `ALPHA`, `BETA`),
		}
		empty := M.SX{}
		// TODO: check session if user may access this menu
		const user_has_access = true
		for _, menu := range menus {
			values[menu+`_menu`] = ``
			if user_has_access {
				values[menu+`_menu`] = ctx.PartialNoDebug(`menu/`+menu, empty)
			}
		}
		// L.Describe(ctx.Session.Level)
		ctx.Render(`menu/index`, values)
	} else {
		// prepare variables required for ajax response
		rm.Posts = ctx.Posts()
		rm.Action = rm.Posts.GetStr(`a`)
		id := rm.Posts.GetStr(`id`)
		rm.Id = S.If(id != `0`, id)
	}
	rm.Ok = true
	return rm
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
