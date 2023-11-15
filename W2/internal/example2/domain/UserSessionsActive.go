package domain

import (
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/S"

	"example2/model/mAuth/rqAuth"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file UserSessionsActive.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type UserSessionsActive.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type UserSessionsActive.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type UserSessionsActive.go
//go:generate farify doublequote --file UserSessionsActive.go

type (
	UserSessionsActiveIn struct {
		RequestCommon
	}
	UserSessionsActiveOut struct {
		ResponseCommon

		SessionsActive []*rqAuth.Sessions `json:"sessionsActive" form:"sessionsActive" query:"sessionsActive" long:"sessionsActive" msg:"sessionsActive"`
	}
)

const (
	UserSessionsActiveAction = `user/sessionsActive`
)

func (d *Domain) UserSessionsActive(in *UserSessionsActiveIn) (out UserSessionsActiveOut) {
	defer d.InsertActionLog(&in.RequestCommon, &out.ResponseCommon)

	sess := d.MustLogin(in.RequestCommon, &out.ResponseCommon)
	if sess == nil {
		return
	}

	logins := rqAuth.NewSessions(d.AuthOltp)
	sessionsActive := logins.AllActiveSession(sess.UserId, in.UnixNow())

	// lets hash the session token
	for i, session := range sessionsActive {
		sessionsActive[i].SessionToken = I.UToS(S.XXH3(session.SessionToken))
	}
	out.SessionsActive = sessionsActive

	return
}
