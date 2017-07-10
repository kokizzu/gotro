package uSuperAdmin

import (
	"example-complete/handler"
	"example-complete/sql/sResponse"
	"example-complete/sql/tUsers"
	"github.com/kokizzu/gotro/D/Pg"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/W"
)

// 2017-06-04 Haries
func Users(ctx *W.Context) {
	rm := sResponse.Prepare(ctx, `Users`, false)
	if !rm.Ok {
		return
	}
	if rm.IsAjax() {
		switch rm.Action {
		case `search`: // @API
			tUsers.API_Superadmin_Search(rm)
		case `form`: // @API
			tUsers.API_Superadmin_Form(rm)
		case `save`, `delete`, `restore`: // @API
			tUsers.API_Superadmin_SaveDeleteRestore(rm)
		default: // @API-END
			handler.ErrorHandler(rm.Ajax, rm.Action)
		}
		ctx.AppendJson(rm.Ajax.SX)
		return
	}

	locals := W.Ajax{SX: M.SX{
		`title`: ctx.Title,
	}}
	qp := Pg.NewQueryParams(nil, &tUsers.TM_MASTER)
	tUsers.Search_ByQueryParams(qp)
	qp.ToMap(locals)
	ctx.Render(`superadmin/users`, locals.SX)
}
