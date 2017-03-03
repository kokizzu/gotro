package W

import (
	"github.com/kokizzu/gotro/M"
)

type RequestModel struct {
	Id      int64
	AppId   int64
	Posts   *Posts
	Ajax    Ajax
	DbActor int64
	Actor   int64
	Level   M.SX
	Ctx     *Context
	Ok      bool
	Action  string
}

func (rm *RequestModel) IsAjax() bool {
	return rm.Ctx.IsAjax()
}

func NewRequestModel_ById_ByDbActor_ByAjax(id, db_actor int64, ajax Ajax) *RequestModel {
	if ajax.SX == nil {
		ajax = NewAjax()
	}
	return &RequestModel{
		Id:      id,
		DbActor: db_actor,
		Ajax:    ajax,
	}
}
