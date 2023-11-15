package domain

import (
	"fmt"

	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"

	"example2/conf"
	"example2/model/mAuth/rqAuth"
	"example2/model/mAuth/wcAuth"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file GuestAutoLogin.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type GuestAutoLogin.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type GuestAutoLogin.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type GuestAutoLogin.go
//go:generate farify doublequote --file GuestAutoLogin.go

type (
	GuestAutoLoginIn struct {
		RequestCommon

		Uid   string `json:"uid" form:"uid" query:"uid" long:"uid" msg:"uid"`
		Token string `json:"token" form:"token" query:"token" long:"token" msg:"token"`
		Path  string `json:"path" form:"path" query:"path" long:"path" msg:"path"`
	}
	GuestAutoLoginOut struct {
		ResponseCommon
		User *rqAuth.Users `json:"user" form:"user" query:"user" long:"user" msg:"user"`

		Segments M.SB `json:"segments" form:"segments" query:"segments" long:"segments" msg:"segments"`
	}
)

const (
	GuestAutoLoginAction = `guest/autoLogin`

	ErrGuestAutoLoginInvalidUid           = `autologin invalid uid`
	ErrGuestAutoLoginUserNotFound         = `autologin user not found`
	ErrGuestAutoLoginInvalidToken         = `autologin invalid token`
	ErrGuestAutoLoginFailedStoringSession = `failed storing session for autologin`
)

func (d *Domain) GuestAutoLogin(in *GuestAutoLoginIn) (out GuestAutoLoginOut) {
	defer d.InsertActionLog(&in.RequestCommon, &out.ResponseCommon)

	userId, ok := S.DecodeCB63[uint64](in.Uid)
	if !ok {
		out.SetError(400, ErrGuestAutoLoginInvalidUid)
		return
	}

	user := wcAuth.NewUsersMutator(d.AuthOltp)
	user.Id = userId
	if !user.FindById() {
		out.SetError(400, ErrGuestAutoLoginUserNotFound)
		return
	}

	sess := &Session{
		UserId:       userId,
		ExpiredAt:    0,
		Email:        ``,
		IsSuperAdmin: false,
		Segments:     nil,
	}

	if !sess.Decrypt(in.Token, conf.AutoLoginUA+fmt.Sprint(user.UpdatedAt)+in.Path) {
		out.SetError(400, ErrGuestAutoLoginInvalidToken)
		return
	}

	user.SetLastLoginAt(in.UnixNow())
	user.SetUpdatedAt(in.UnixNow()) // to expire the autologin
	if !user.DoUpsert() {
		out.AddTrace(WarnFailedSetLastLoginAt)
		return
	}
	user.CensorFields()
	out.User = &user.Users
	session, sess := d.CreateSession(user.Id, user.Email, in.UserAgent, in.IpAddress)

	// TODO: set list of roles in the session
	if !session.DoInsert() {
		out.SetError(500, ErrGuestAutoLoginFailedStoringSession)
		return
	}
	out.SessionToken = session.SessionToken
	out.Segments = d.segmentsFromSession(sess)

	out.SetRedirect(in.Path)
	return
}
