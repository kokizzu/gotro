package sResponse

import (
	"example-complete/sql"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/W"
)

func NewAjax() W.Ajax {
	return W.Ajax{
		SX: M.SX{
			`is_success`: true,
		},
	}
}
func Prepare(ctx *W.Context, title string, must_login bool) (rm *W.RequestModel) {
	user_id := ctx.Session.GetInt(`user_id`)
	rm = &W.RequestModel{
		Actor:   I.ToS(user_id),
		DbActor: I.ToS(user_id),
		Level:   ctx.Session.GetMSX(`level`),
		Ctx:     ctx,
	}
	is_ajax := ctx.IsAjax()
	if is_ajax {
		rm.Ajax = NewAjaxResponse()
	}
	page := rm.Level.GetMSB(`page`)
	first_segment := ctx.FirstPath()
	is_wm := ctx.IsWebMaster()
	rm.Ok = is_wm || page[first_segment] || first_segment == `guest`
	if rm.Ok {
		// TODO: check access level, or use AuthFilter
	}
	if !rm.Ok {
		// 403
		if is_ajax {
			rm.Ajax.Error(sql.ERR_001_MUST_LOGIN)
			ctx.AppendJson(rm.Ajax.SX)
			return
		}
		ctx.Title = title
		if must_login {
			Render403(ctx)
		}
	}
	if !is_ajax {
		// GET
		ctx.Title = title
		menus := []string{
			`guest`,
			`superadmin`,
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
		for _, menu := range menus {
			values[menu+`_menu`] = ``
			if menu == `guest` || page[menu] || is_wm {
				values[menu+`_menu`] = ctx.PartialNoDebug(`menu/`+menu, empty)
			}
		}
		// L.Describe(ctx.Session.Level)
		ctx.Render(`menu/index`, values)
	} else {
		// POST
		rm.Posts = ctx.Posts()
		rm.Action = rm.Posts.GetStr(`a`)
		rm.Id = rm.Posts.GetStr(`id`)
	}
	return
}

func Render403(ctx *W.Context) {
	values := M.SX{
		`requested_path`: ctx.Path,
		`webmaster`:      sql.SUPPORT_EMAIL,
	}
	ctx.Render(`403`, values)
}

func NewAjaxResponse() W.Ajax {
	return W.Ajax{SX: M.SX{`is_success`: true}}
}
