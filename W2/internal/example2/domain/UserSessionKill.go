package domain

import (
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/S"

	"example2/model/mAuth/rqAuth"
	"example2/model/mAuth/wcAuth"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file UserSessionKill.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type UserSessionKill.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type UserSessionKill.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type UserSessionKill.go
//go:generate farify doublequote --file UserSessionKill.go

type (
	UserSessionKillIn struct {
		RequestCommon

		SessionTokenHash string `json:"sessionTokenHash" form:"sessionTokenHash" query:"sessionTokenHash" long:"sessionTokenHash" msg:"sessionTokenHash"`
	}
	UserSessionKillOut struct {
		ResponseCommon
		LogoutAt          int64 `json:"loggedOut" form:"loggedOut" query:"loggedOut" long:"loggedOut" msg:"loggedOut"`
		SessionTerminated int64 `json:"sessionTerminated" form:"sessionTerminated" query:"sessionTerminated" long:"sessionTerminated" msg:"sessionTerminated"`
	}
)

const (
	UserSessionKillAction = `user/sessionKill`

	ErrUserSessionTerminationFailed = `user session termination failed`
)

func (d *Domain) UserSessionKill(in *UserSessionKillIn) (out UserSessionKillOut) {
	defer d.InsertActionLog(&in.RequestCommon, &out.ResponseCommon)

	sess := d.MustLogin(in.RequestCommon, &out.ResponseCommon)
	if sess == nil {
		return
	}

	logins := rqAuth.NewSessions(d.AuthOltp)
	sessionList := logins.AllActiveSession(sess.UserId, in.UnixNow())

	now := in.UnixNow()
	for _, session := range sessionList {
		if I.UToS(S.XXH3(session.SessionToken)) == in.SessionTokenHash {

			if session.ExpiredAt > now {

				// create mutator
				session.Adapter = d.AuthOltp
				toUpdate := wcAuth.NewSessionsMutator(d.AuthOltp)
				toUpdate.Sessions = *session
				toUpdate.SetExpiredAt(now)

				// make it expired
				if toUpdate.DoUpdateBySessionToken() {
					out.LogoutAt = toUpdate.ExpiredAt
					out.SessionTerminated++
				} else {
					out.SetError(500, ErrUserSessionTerminationFailed)
					return
				}
			}
		}
	}

	return
}
