package hBackoffice

import (
	"github.com/kokizzu/gotro/D/Pg"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/W"
	"github.com/kokizzu/gotro/W/example-complex/model/mResponse"
	"github.com/kokizzu/gotro/W/example-complex/model/mUsers"
)

func Users(ctx *W.Context) {
	rm := mResponse.PrepareVars(ctx, `Users`)
	if !rm.Ok {
		return
	}
	if rm.IsAjax() {
		// handle ajax
		switch rm.Action {
		case `search`: // @API
			mUsers.API_Backoffice_Search(rm)
		case `form_limit`: // @API
			mUsers.API_Backoffice_FormLimit(rm)
		case `form`: // @API
			mUsers.API_Backoffice_Form(rm)
		case `save`, `delete`, `restore`: // @ffPI
			mUsers.API_Backoffice_SaveDeleteRestore(rm)
		default: // @API-END
			//handler.ErrorHandler(rm.Ajax, rm.Action)
		}
		ctx.AppendJson(rm.Ajax.SX)
		return
	}
	locals := W.Ajax{SX: M.SX{
		`title`: ctx.Title,
	}}
	qp := Pg.NewQueryParams(nil, &mUsers.TM_MASTER)
	mUsers.Search_ByQueryParams(qp)
	qp.ToMap(locals)
	ctx.Render(`backoffice/users`, locals.SX)
}
