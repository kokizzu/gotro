package domain

import (
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/W2/example/conf"
	"github.com/kokizzu/gotro/W2/example/model/mAuth/rqAuth"
	"github.com/kokizzu/gotro/W2/example/model/mAuth/wcAuth"
	"github.com/kokizzu/gotro/X"
	"github.com/kokizzu/id64"
	"github.com/vburenin/nsync"
)

//go:generate gomodifytags -file user.go -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported --skip-unexported -w -file user.go
//go:generate replacer 'Id" form' 'Id,string" form' type user.go
//go:generate replacer 'json:"id"' 'json:id,string" form' type user.go
//go:generate replacer 'By" form' 'By,string" form' type user.go

type (
	UserRegister_In struct {
		RequestCommon
		UserName string `json:"userName" form:"userName" query:"userName" long:"userName" msg:"userName"`
		Email    string `json:"email" form:"email" query:"email" long:"email" msg:"email"`
		Password string `json:"password" form:"password" query:"password" long:"password" msg:"password"`
	}
	UserRegister_Out struct {
		ResponseCommon
		CreatedAt int64  `json:"createdAt" form:"createdAt" query:"createdAt" long:"createdAt" msg:"createdAt"`
		UserId    uint64 `json:"userId,string" form:"userId" query:"userId" long:"userId" msg:"userId"`
	}
)

const UserRegister_Url = `/UserRegister`

func (d *Domain) UserRegister(in *UserRegister_In) (out UserRegister_Out) {

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
	out.UserId = user.Id

	return
}

type (
	UserLogin_In struct {
		RequestCommon
		Email    string `json:"email" form:"email" query:"email" long:"email" msg:"email"`
		Password string `json:"password" form:"password" query:"password" long:"password" msg:"password"`
	}
	UserLogin_Out struct {
		ResponseCommon
	}
)

const UserLogin_Url = `/UserLogin`

func (d *Domain) UserLogin(in *UserLogin_In) (out UserLogin_Out) {
	user := rqAuth.NewUsers(d.Taran)
	user.Email = in.Email
	if !user.FindByEmail() {
		out.SetError(401, `wrong email or password`)
		return
	}
	if !user.CheckPassword(in.Password) {
		out.SetError(401, `wrong password`)
		return
	}

	d.expireSession(in.SessionToken)

	// create session
	session := d.createSession(user.Id, user.Email, in.UserAgent)
	if !session.DoInsert() {
		out.SetError(500, `cannot create session`)
		return
	}
	out.SessionToken = session.SessionToken
	return
}

type (
	UserLogout_In struct {
		RequestCommon
	}
	UserLogout_Out struct {
		ResponseCommon
		LoggedOut bool `json:"loggedOut" form:"loggedOut" query:"loggedOut" long:"loggedOut" msg:"loggedOut"`
	}
)

const UserLogout_Url = `/UserLogout`

func (d *Domain) UserLogout(in *UserLogout_In) (out UserLogout_Out) {
	loggedOut := d.expireSession(in.SessionToken)
	out.LoggedOut = loggedOut
	out.SessionToken = conf.CookieLogoutValue
	return
}

type (
	UserProfile_In struct {
		RequestCommon
	}
	UserProfile_Out struct {
		ResponseCommon
		User *rqAuth.Users `json:"user" form:"user" query:"user" long:"user" msg:"user"`
	}
)

const UserProfile_Url = `/UserProfile`

func (d *Domain) UserProfile(in *UserProfile_In) (out UserProfile_Out) {
	sess := d.mustLogin(in.SessionToken, in.UserAgent, &out.ResponseCommon)
	if sess == nil {
		return
	}

	user := rqAuth.NewUsers(d.Taran)
	user.Id = sess.UserId
	if !user.FindById() {
		out.SetError(404, `user does not exists on database: `+X.ToS(sess.UserId))
		return
	}
	user.CensorFields()
	out.User = user
	return
}

type (
	UserList_In struct {
		RequestCommon
		Limit  uint32 `json:"limit" form:"limit" query:"limit" long:"limit" msg:"limit"`
		Offset uint32 `json:"offset" form:"offset" query:"offset" long:"offset" msg:"offset"`
	}
	UserList_Out struct {
		ResponseCommon
		Limit  uint32          `json:"limit" form:"limit" query:"limit" long:"limit" msg:"limit"`
		Offset uint32          `json:"offset" form:"offset" query:"offset" long:"offset" msg:"offset"`
		Total  uint32          `json:"total" form:"total" query:"total" long:"total" msg:"total"`
		Users  []*rqAuth.Users `json:"users" form:"users" query:"users" long:"users" msg:"users"`
	}
)

const UserList_Url = `/UserList`

func (d *Domain) UserList(in *UserList_In) (out UserList_Out) {
	user := rqAuth.NewUsers(d.Taran)
	out.Users = user.FindOffsetLimit(in.Offset, in.Limit)
	out.Total = uint32(user.Total())
	out.Limit = in.Limit
	out.Offset = in.Offset
	return
}

type (
	UserForgotPassword_In struct {
		RequestCommon
		Email              string `json:"email" form:"email" query:"email" long:"email" msg:"email"`
		ChangePassCallback string `json:"changePassCallback" form:"changePassCallback" query:"changePassCallback" long:"changePassCallback" msg:"changePassCallback"`
	}
	UserForgotPassword_Out struct {
		ResponseCommon
		Ok bool `json:"ok" form:"ok" query:"ok" long:"ok" msg:"ok"`
	}
)

const UserForgotPassword_Url = `/UserForgotPassword`

var forgotPasswordLock = nsync.NewNamedMutex()

func (d *Domain) UserForgotPassword(in *UserForgotPassword_In) (out UserForgotPassword_Out) {
	user := wcAuth.NewUsersMutator(d.Taran)
	user.Email = in.Email

	forgotPasswordLock.Lock(in.Email)
	if !user.FindByEmail() {
		out.SetError(400, `email not found`)
		return
	}
	secretCode := id64.SID()
	user.SetSecretCode(secretCode)
	user.SetSecretCodeAt(in.UnixNow())
	hash := S.EncodeCB63(int64(user.Id), 8)
	if in.ChangePassCallback == `` {
		in.ChangePassCallback = conf.WEBAPI_HOSTPORT + conf.API_PREFIX + UserResetPassword_Url
	}
	url := in.ChangePassCallback + `?secretCode=` + secretCode + `&hash=` + hash
	go func(email, url string) {
		defer forgotPasswordLock.Unlock(email)
		// go //go d.SendMail(ForgotPassMail{...},user.Email)
	}(in.Email, url)

	if !user.DoUpdateById() {
		out.SetError(500, `failed to update row on database`)
		return
	}
	out.Ok = true
	return
}

type (
	UserResetPassword_In struct {
		RequestCommon
		Password   string `json:"password" form:"password" query:"password" long:"password" msg:"password"`
		SecretCode string `query:"secretCode" json:"secretCode" form:"secretCode" long:"secretCode" msg:"secretCode"`
		Hash       string `query:"hash" json:"hash" form:"hash" long:"hash" msg:"hash"`
	}
	UserResetPassword_Out struct {
		ResponseCommon
		Ok bool `json:"ok" form:"ok" query:"ok" long:"ok" msg:"ok"`
	}
)

const UserResetPassword_Url = `/UserResetPassword`

func (d *Domain) UserResetPassword(in *UserResetPassword_In) (out UserResetPassword_Out) {
	userId, ok := S.DecodeCB63(in.Hash)
	if !ok {
		out.SetError(400, `invalid hash`)
		return
	}
	user := wcAuth.NewUsersMutator(d.Taran)
	user.Id = uint64(userId)
	if !user.FindById() {
		out.SetError(400, `cannot find user, wrong env?`)
		return
	}
	if user.SecretCode != in.SecretCode {
		out.SetError(400, `invalid secret code`)
		return
	}
	if in.UnixNow()-user.SecretCodeAt >= 60*45 { // 45 minutes
		out.SetError(400, `secret code expired`)
		return
	}
	if !user.SetEncryptPassword(in.Password) {
		out.SetError(500, `cannot encrypt password`)
		return
	}
	user.SetSecretCode(``)
	user.SetSecretCodeAt(0)
	user.SetPasswordSetAt(in.UnixNow())
	if !user.DoUpdateById() {
		out.SetError(500, `failed to update user password`)
		return
	}
	out.Ok = true
	return
}

type (
	UserChangePassword_In struct {
		RequestCommon
		Password    string `json:"password" form:"password" query:"password" long:"password" msg:"password"`
		NewPassword string `json:"newPassowrd" form:"newPassword" query:"newPassword" long:"newPassword" msg:"newPassword"`
	}
	UserChangePassword_Out struct {
		ResponseCommon
		UpdatedAt int64 `json:"updatedAt" form:"updatedAt" query:"updatedAt" long:"updatedAt" msg:"updatedAt"`
	}
)

const UserChangePassword_Url = `/UserChangePassword`

func (d *Domain) UserChangePassword(in *UserChangePassword_In) (out UserChangePassword_Out) {
	sess := d.mustLogin(in.SessionToken, in.UserAgent, &out.ResponseCommon)
	if sess == nil {
		return
	}
	user := wcAuth.NewUsersMutator(d.Taran)
	user.Id = sess.UserId
	if !user.FindById() {
		out.SetError(500, `user not found`)
		return
	}
	if !user.CheckPassword(in.Password) {
		out.SetError(401, `wrong password`)
		return
	}
	if !user.SetEncryptPassword(in.Password) {
		out.SetError(500, `cannot encrypt password`)
		return
	}
	user.SetUpdatedAt(in.UnixNow())
	if !user.DoUpdateById() {
		out.SetError(500, `failed to change password`)
		return
	}
	return
}

type (
	UserConfirmEmail_In struct {
		RequestCommon
	}
	UserConfirmEmail_Out struct {
		ResponseCommon
	}
)

const UserConfirmEmail_Url = `/UserConfirmEmail`

func (d *Domain) UserConfirmEmail(in *UserConfirmEmail_In) (out UserConfirmEmail_Out) {
	// TODO: continue this
	return
}

type (
	UserChangeEmail_In struct {
		RequestCommon
	}
	UserChangeEmail_Out struct {
		ResponseCommon
	}
)

const UserChangeEmail_Url = `/UserChangeEmail`

func (d *Domain) UserChangeEmail(in *UserChangeEmail_In) (out UserChangeEmail_Out) {
	// TODO: continue this
	return
}
