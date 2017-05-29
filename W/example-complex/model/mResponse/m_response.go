package mResponse

import (
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/W"
)

func MewAjaxResponse() W.Ajax {
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
		rm.Ajax = MewAjaxResponse()
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
