package domain

import (
	"time"

	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/W2/example/conf"
	"github.com/kokizzu/gotro/W2/example/model/mAuth/rqAuth"
	"github.com/kokizzu/gotro/W2/example/model/mAuth/wcAuth"
	"github.com/kokizzu/gotro/X"
	"github.com/kokizzu/id64"
	"github.com/kpango/fastime"
	"github.com/vburenin/nsync"
)

//go:generate gomodifytags -file player.go -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported --skip-unexported -w -file user.go
//go:generate replacer 'Id" form' 'Id,string" form' type user.go
//go:generate replacer 'json:"id"' 'json:id,string" form' type user.go
//go:generate replacer 'By" form' 'By,string" form' type user.go

type (
	PlayerRegister_In struct {
		RequestCommon
		UserName string `json:"userName" form:"userName" query:"userName" long:"userName" msg:"userName"`
		Email    string `json:"email" form:"email" query:"email" long:"email" msg:"email"`
		Password string `json:"password" form:"password" query:"password" long:"password" msg:"password"`
	}
	PlayerRegister_Out struct {
		ResponseCommon
		CreatedAt int64  `json:"createdAt" form:"createdAt" query:"createdAt" long:"createdAt" msg:"createdAt"`
		PlayerId  uint64 `json:"playerId,string" form:"playerId" query:"playerId" long:"playerId" msg:"playerId"`
	}
)

const PlayerRegister_Url = `/PlayerRegister`

func (d *Domain) PlayerRegister(in *PlayerRegister_In) (out PlayerRegister_Out) {

	username := S.Trim(in.UserName)
	if username == `` {
		out.SetError(400, `userName must not be empty`)
		return
	}

	if in.Email == `` {
		out.SetError(400, `email must not be empty`)
		return
	}

	user := wcAuth.NewUsersMutator(d.Taran)
	email := S.Trim(in.Email)
	user.Email = email
	if user.FindByEmail() {
		out.SetError(451, `user already exists: `+email)
		return
	}

	// create user
	user.Id = id64.UID()
	user.Email = email
	if in.Password == `` {
		in.Password = email
	}
	if !user.SetEncryptPassword(in.Password) {
		out.SetError(500, `cannot encrypt password`)
		return
	}
	if !user.DoInsert() {
		out.SetError(451, `failed to register this user: `+email)
		return
	}

	//go d.SendMail(WelcomeMail{...},user.Email)
	out.CreatedAt = user.CreatedAt
	out.PlayerId = user.Id

	return
}

type (
	PlayerLogin_In struct {
		RequestCommon
		Email    string `json:"email" form:"email" query:"email" long:"email" msg:"email"`
		Password string `json:"password" form:"password" query:"password" long:"password" msg:"password"`
	}
	PlayerLogin_Out struct {
		ResponseCommon
		WalletId string `json:"walletId,string" form:"walletId" query:"walletId" long:"walletId" msg:"walletId"` // wallet on stardust
	}
)

const PlayerLogin_Url = `/PlayerLogin`

func (d *Domain) expireSession(sessionToken string) bool {
	if sessionToken == `` {
		return false
	}
	session := wcAuth.NewSessionsMutator(d.Taran)
	session.SessionToken = sessionToken
	now := fastime.UnixNow()
	if session.FindBySessionToken() {
		if session.ExpiredAt > now {
			session.SetExpiredAt(now)
			session.DoUpdateBySessionToken()
		}
		return true
	}
	return false
}

func (d *Domain) createSession(playerId uint64, email, userAgent string) *wcAuth.SessionsMutator {
	session := wcAuth.NewSessionsMutator(d.Taran)
	session.UserId = playerId
	sess := conf.Session{
		PlayerId:  playerId,
		Email:     email,
		ExpiredAt: time.Now().AddDate(0, 0, conf.CookieDays).Unix(),
	}
	session.SessionToken = sess.Encrypt(userAgent)
	session.ExpiredAt = sess.ExpiredAt
	return session
}

func (d *Domain) PlayerLogin(in *PlayerLogin_In) (out PlayerLogin_Out) {
	player := rqAuth.NewUsers(d.Taran)
	player.Email = in.Email
	if !player.FindByEmail() {
		out.SetError(401, `wrong email or password`)
		return
	}
	if !player.CheckPassword(in.Password) {
		out.SetError(401, `wrong password`)
		return
	}

	d.expireSession(in.SessionToken)

	// create session
	session := d.createSession(player.Id, player.Email, in.UserAgent)
	if !session.DoInsert() {
		out.SetError(500, `cannot create session`)
		return
	}
	out.SessionToken = session.SessionToken
	return
}

type (
	PlayerLogout_In struct {
		RequestCommon
	}
	PlayerLogout_Out struct {
		ResponseCommon
		LoggedOut bool `json:"loggedOut" form:"loggedOut" query:"loggedOut" long:"loggedOut" msg:"loggedOut"`
	}
)

const PlayerLogout_Url = `/PlayerLogout`

func (d *Domain) PlayerLogout(in *PlayerLogout_In) (out PlayerLogout_Out) {
	loggedOut := d.expireSession(in.SessionToken)
	out.LoggedOut = loggedOut
	out.SessionToken = conf.CookieLogoutValue
	return
}

type (
	PlayerProfile_In struct {
		RequestCommon
	}
	PlayerProfile_Out struct {
		ResponseCommon
		Player *rqAuth.Users `json:"player" form:"player" query:"player" long:"player" msg:"player"`
	}
)

const PlayerProfile_Url = `/PlayerProfile`

func (d *Domain) PlayerProfile(in *PlayerProfile_In) (out PlayerProfile_Out) {
	sess := d.mustLogin(in.SessionToken, in.UserAgent, &out.ResponseCommon)
	if sess == nil {
		return
	}

	player := rqAuth.NewUsers(d.Taran)
	player.Id = sess.PlayerId
	if !player.FindById() {
		out.SetError(404, `player does not exists on database: `+X.ToS(sess.PlayerId))
		return
	}
	player.CensorFields()
	out.Player = player
	return
}

type (
	PlayerUpdateProfile_In struct {
		RequestCommon
		UserName string `json:"userName" form:"userName" query:"userName" long:"userName" msg:"userName"`
	}
	PlayerUpdateProfile_Out struct {
		ResponseCommon
		Ok bool `json:"ok" form:"ok" query:"ok" long:"ok" msg:"ok"`
	}
)

const PlayerUpdateProfile_Url = `/PlayerUpdateProfile`

func (d *Domain) PlayerUpdateProfile(in *PlayerUpdateProfile_In) (out PlayerUpdateProfile_Out) {
	sess := d.mustLogin(in.SessionToken, in.UserAgent, &out.ResponseCommon)
	if sess == nil {
		return
	}
	player := wcAuth.NewUsersMutator(d.Taran)
	player.Id = sess.PlayerId
	if !player.FindById() {
		out.SetError(400, `player not found in database. wrong env?`)
		return
	}

	// example if there's multiple unique column on table:
	//player.SetUserName(in.UserName)
	//if player.FindByUserName() && player.Id != sess.PlayerId {
	//	out.SetError(400, `userName already used by other player`)
	//	return
	//}

	if !player.DoUpdateById() {
		out.SetError(500, `failed to update profile`)
		return
	}

	out.Ok = true
	return
}

type (
	PlayerList_In struct {
		RequestCommon
		Limit  uint32 `json:"limit" form:"limit" query:"limit" long:"limit" msg:"limit"`
		Offset uint32 `json:"offset" form:"offset" query:"offset" long:"offset" msg:"offset"`
	}
	PlayerList_Out struct {
		ResponseCommon
		Limit   uint32          `json:"limit" form:"limit" query:"limit" long:"limit" msg:"limit"`
		Offset  uint32          `json:"offset" form:"offset" query:"offset" long:"offset" msg:"offset"`
		Total   uint32          `json:"total" form:"total" query:"total" long:"total" msg:"total"`
		Players []*rqAuth.Users `json:"players" form:"players" query:"players" long:"players" msg:"players"`
	}
)

const PlayerList_Url = `/PlayerList`

func (d *Domain) PlayerList(in *PlayerList_In) (out PlayerList_Out) {
	player := rqAuth.NewUsers(d.Taran)
	out.Players = player.FindOffsetLimit(in.Offset, in.Limit)
	out.Total = uint32(player.Total())
	out.Limit = in.Limit
	out.Offset = in.Offset
	for k := range out.Players {
		out.Players[k].CensorFields()
	}
	return
}

type (
	PlayerForgotPassword_In struct {
		RequestCommon
		Email              string `json:"email" form:"email" query:"email" long:"email" msg:"email"`
		ChangePassCallback string `json:"changePassCallback" form:"changePassCallback" query:"changePassCallback" long:"changePassCallback" msg:"changePassCallback"`
	}
	PlayerForgotPassword_Out struct {
		ResponseCommon
		Ok bool `json:"ok" form:"ok" query:"ok" long:"ok" msg:"ok"`
	}
)

const PlayerForgotPassword_Url = `/PlayerForgotPassword`

var forgotPasswordLock = nsync.NewNamedMutex()

func (d *Domain) PlayerForgotPassword(in *PlayerForgotPassword_In) (out PlayerForgotPassword_Out) {
	player := wcAuth.NewUsersMutator(d.Taran)
	player.Email = in.Email

	forgotPasswordLock.Lock(in.Email)
	if !player.FindByEmail() {
		out.SetError(400, `email not found`)
		return
	}
	secretCode := id64.SID()
	player.SetSecretCode(secretCode)
	player.SetSecretCodeAt(in.UnixNow())
	hash := S.EncodeCB63(int64(player.Id), 8)
	if in.ChangePassCallback == `` {
		in.ChangePassCallback = conf.WEBAPI_HOSTPORT + conf.API_PREFIX + PlayerResetPassword_Url
	}
	url := in.ChangePassCallback + `?secretCode=` + secretCode + `&hash=` + hash
	go func(email, url string) {
		defer forgotPasswordLock.Unlock(email)
		// go //go d.SendMail(ForgotPassMail{...},user.Email)
	}(in.Email, url)

	if !player.DoUpdateById() {
		out.SetError(500, `failed to update row on database`)
		return
	}
	out.Ok = true
	return
}

type (
	PlayerResetPassword_In struct {
		RequestCommon
		Password   string `json:"password" form:"password" query:"password" long:"password" msg:"password"`
		SecretCode string `query:"secretCode" json:"secretCode" form:"secretCode" long:"secretCode" msg:"secretCode"`
		Hash       string `query:"hash" json:"hash" form:"hash" long:"hash" msg:"hash"`
	}
	PlayerResetPassword_Out struct {
		ResponseCommon
		Ok bool `json:"ok" form:"ok" query:"ok" long:"ok" msg:"ok"`
	}
)

const PlayerResetPassword_Url = `/PlayerResetPassword`

func (d *Domain) PlayerResetPassword(in *PlayerResetPassword_In) (out PlayerResetPassword_Out) {
	playerId, ok := S.DecodeCB63(in.Hash)
	if !ok {
		out.SetError(400, `invalid hash`)
		return
	}
	player := wcAuth.NewUsersMutator(d.Taran)
	player.Id = uint64(playerId)
	if !player.FindById() {
		out.SetError(400, `cannot find player, wrong env?`)
		return
	}
	if player.SecretCode != in.SecretCode {
		out.SetError(400, `invalid secret code`)
		return
	}
	if in.UnixNow()-player.SecretCodeAt >= 60*45 { // 45 minutes
		out.SetError(400, `secret code expired`)
		return
	}
	if !player.SetEncryptPassword(in.Password) {
		out.SetError(500, `cannot encrypt password`)
		return
	}
	player.SetSecretCode(``)
	player.SetSecretCodeAt(0)
	player.SetPasswordSetAt(in.UnixNow())
	if !player.DoUpdateById() {
		out.SetError(500, `failed to update player password`)
		return
	}
	out.Ok = true
	return
}

type (
	PlayerChangePassword_In struct {
		RequestCommon
		Password    string `json:"password" form:"password" query:"password" long:"password" msg:"password"`
		NewPassword string `json:"newPassowrd" form:"newPassword" query:"newPassword" long:"newPassword" msg:"newPassword"`
	}
	PlayerChangePassword_Out struct {
		ResponseCommon
		UpdatedAt int64 `json:"updatedAt" form:"updatedAt" query:"updatedAt" long:"updatedAt" msg:"updatedAt"`
	}
)

const PlayerChangePassword_Url = `/PlayerChangePassword`

func (d *Domain) PlayerChangePassword(in *PlayerChangePassword_In) (out PlayerChangePassword_Out) {
	sess := d.mustLogin(in.SessionToken, in.UserAgent, &out.ResponseCommon)
	if sess == nil {
		return
	}
	player := wcAuth.NewUsersMutator(d.Taran)
	player.Id = sess.PlayerId
	if !player.FindById() {
		out.SetError(500, `player not found`)
		return
	}
	if !player.CheckPassword(in.Password) {
		out.SetError(401, `wrong password`)
		return
	}
	if !player.SetEncryptPassword(in.Password) {
		out.SetError(500, `cannot encrypt password`)
		return
	}
	player.SetUpdatedAt(in.UnixNow())
	if !player.DoUpdateById() {
		out.SetError(500, `failed to change password`)
		return
	}
	return
}

type (
	PlayerConfirmEmail_In struct {
		RequestCommon
	}
	PlayerConfirmEmail_Out struct {
		ResponseCommon
	}
)

const PlayerConfirmEmail_Url = `/PlayerConfirmEmail`

func (d *Domain) PlayerConfirmEmail(in *PlayerConfirmEmail_In) (out PlayerConfirmEmail_Out) {
	// TODO: continue this
	return
}

type (
	PlayerChangeEmail_In struct {
		RequestCommon
	}
	PlayerChangeEmail_Out struct {
		ResponseCommon
	}
)

const PlayerChangeEmail_Url = `/PlayerChangeEmail`

func (d *Domain) PlayerChangeEmail(in *PlayerChangeEmail_In) (out PlayerChangeEmail_Out) {
	// TODO: continue this
	return
}
