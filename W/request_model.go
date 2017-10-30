package W

import (
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
)

type RequestModel struct {
	Id      string
	Uniq    string
	AppId   string
	Posts   *Posts
	Ajax    Ajax
	DbActor string
	Actor   string
	Level   M.SX
	Ctx     *Context
	Ok      bool
	Action  string
}

func (rm *RequestModel) IsAjax() bool {
	return rm.Ctx.IsAjax()
}

func NewRequestModel_ById_ByDbActor_ByAjax(id, db_actor string, ajax Ajax) *RequestModel {
	if ajax.SX == nil {
		ajax = NewAjax()
	}
	return &RequestModel{
		Id:      id,
		DbActor: db_actor,
		Ajax:    ajax,
	}
}

func NewRequestModel_ByUniq_ByDbActor_ByAjax(uniq_id string, db_actor string, ajax Ajax) *RequestModel {
	if ajax.SX == nil {
		ajax = NewAjax()
	}
	return &RequestModel{
		Uniq:    uniq_id,
		DbActor: db_actor,
		Ajax:    ajax,
	}
}

func (rm *RequestModel) IdInt() int64 {
	return S.ToI(rm.Id)
}

func (rm *RequestModel) HasAjaxError(err error, errmsg string) bool {
	if L.IsError(err, errmsg) {
		rm.Ajax.Error(errmsg)
	}
	return rm.Ajax.HasError()
}
