package domain

import (
	"testing"
	"time"

	"github.com/kokizzu/id64"
	"github.com/kokizzu/lexid"
	"github.com/kpango/fastime"
	"github.com/stretchr/testify/assert"

	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/W2/example/model/mAuth/wcAuth"
)

const testDomain = `@localhost`

func TestDomain_UserLoginRegisterFlow(t *testing.T) {
	d := NewDomain()
	name := id64.ID().String()
	pass := lexid.ID()
	email := name + testDomain
	t.Run(`register should ok`, func(t *testing.T) {
		in := &UserRegister_In{
			Email:    email,
			Password: pass,
			UserName: name,
		}
		out := d.UserRegister(in)
		assert.Empty(t, out.Error)
	})

	t.Run(`re-register same email should fail`, func(t *testing.T) {
		in := &UserRegister_In{
			Email:    email,
			Password: pass,
			UserName: name + `1`,
		}
		out := d.UserRegister(in)
		assert.NotEmpty(t, out.Error)
		assert.NotEqual(t, 200, out.StatusCode)
		assert.NotEqual(t, 0, out.StatusCode)
	})

	t.Run(`login with unregistered user should fail`, func(t *testing.T) {
		in := &UserLogin_In{
			Email:    name + `notExists` + testDomain,
			Password: pass,
		}
		out := d.UserLogin(in)
		assert.NotEmpty(t, out.Error)
	})

	t.Run(`login with wrong password should fail`, func(t *testing.T) {
		in := &UserLogin_In{
			Email:    email,
			Password: name,
		}
		out := d.UserLogin(in)
		assert.NotEmpty(t, out.Error)
	})

	sessionToken := ``
	t.Run(`login with correct password should ok`, func(t *testing.T) {
		in := &UserLogin_In{
			Email:    email,
			Password: pass,
		}
		out := d.UserLogin(in)
		assert.Empty(t, out.Error)
		assert.NotEmpty(t, out.SessionToken)
		sessionToken = out.SessionToken
	})

	t.Run(`check profile with active session should ok`, func(t *testing.T) {
		in := &UserProfile_In{NewRC(sessionToken)}
		out := d.UserProfile(in)
		assert.Empty(t, out.Error)
		assert.NotNil(t, out.User)
		if out.User == nil {
			t.Failed()
			return
		}
		assert.Equal(t, email, out.User.Email)
	})

	t.Run(`change password with wrong password must fail`, func(t *testing.T) {
		in := &UserChangePassword_In{
			RequestCommon: NewRC(sessionToken),
			Password:      ``,
			NewPassword:   `abc`,
		}
		out := d.UserChangePassword(in)
		assert.NotEmpty(t, out.Error)
	})

	t.Run(`change password with correct password must ok`, func(t *testing.T) {
		in := &UserChangePassword_In{
			RequestCommon: NewRC(sessionToken),
			Password:      pass,
			NewPassword:   email,
		}
		out := d.UserChangePassword(in)
		assert.Empty(t, out.Error)
	})

	t.Run(`logout should ok`, func(t *testing.T) {
		in := &UserLogout_In{RequestCommon{SessionToken: sessionToken}}
		out := d.UserLogout(in)
		assert.Equal(t, out.LoggedOut, true)
	})

	t.Run(`check profile with expired session should fail`, func(t *testing.T) {
		in := &UserProfile_In{NewRC(sessionToken)}
		out := d.UserProfile(in)
		assert.Equal(t, 403, out.StatusCode)
	})
}

func TestDomain_UserList(t *testing.T) {
	d := NewDomain()

	// if fail, then probably no data
	t.Run(`list user must ok`, func(t *testing.T) {
		in := &UserList_In{
			Limit:  2,
			Offset: 0,
		}
		out := d.UserList(in)
		assert.Empty(t, out.Error)
		assert.NotEmpty(t, out.Users)
		assert.Greater(t, len(out.Users), 1)
		L.Describe(out.Users)
	})
}

func dummyUser(d *Domain) (*wcAuth.UsersMutator, string) {

	p1 := wcAuth.NewUsersMutator(d.Taran)
	p1.Id = id64.UID()
	p1.Email = lexid.ID() + testDomain
	p1.CreatedAt = fastime.UnixNow()
	p1.SetEncryptPassword(p1.Email)
	if !p1.DoInsert() {
		panic(`cannot create dummy user `)
	}

	sess := d.createSession(p1.Id, p1.Email, ``)
	if !sess.DoInsert() {
		panic(`cannot create session`)
	}

	return p1, sess.SessionToken
}

func TestDomain_UserForgotReset(t *testing.T) {
	d := NewDomain()
	p1, _ := dummyUser(d)
	hash := S.EncodeCB63(int64(p1.Id), 1)
	newPass := `12345678`

	t.Run(`user forgot password must ok`, func(t *testing.T) {
		in := &UserForgotPassword_In{
			Email: p1.Email,
		}
		out := d.UserForgotPassword(in)
		time.Sleep(1 * time.Second)
		assert.Empty(t, out.Error)
	})

	t.Run(`user login with old password must ok`, func(t *testing.T) {
		loginIn := &UserLogin_In{Email: p1.Email, Password: p1.Email}
		loginOut := d.UserLogin(loginIn)
		assert.Empty(t, loginOut.Error)
	})

	t.Run(`user forgot password with invalid email must fail`, func(t *testing.T) {
		in := &UserForgotPassword_In{
			Email: id64.ID().String() + testDomain,
		}
		out := d.UserForgotPassword(in)
		assert.NotEmpty(t, out.Error)
	})

	p1.FindById() // to get SecretCode

	t.Run(`user reset password with incorrect secretCode must fail`, func(t *testing.T) {
		in := &UserResetPassword_In{
			Password:   newPass,
			SecretCode: `duar`,
			Hash:       hash,
		}
		out := d.UserResetPassword(in)
		assert.NotEmpty(t, out.Error)
	})

	t.Run(`user login with new password must fail`, func(t *testing.T) {
		loginIn := &UserLogin_In{Email: p1.Email, Password: newPass}
		loginOut := d.UserLogin(loginIn)
		assert.NotEmpty(t, loginOut.Error)
	})

	t.Run(`user reset password with incorrect hash must fail`, func(t *testing.T) {
		in := &UserResetPassword_In{
			Password:   newPass,
			SecretCode: p1.SecretCode,
			Hash:       `-`,
		}
		out := d.UserResetPassword(in)
		assert.NotEmpty(t, out.Error)
	})

	t.Run(`user login with new password must fail`, func(t *testing.T) {
		loginIn := &UserLogin_In{Email: p1.Email, Password: newPass}
		loginOut := d.UserLogin(loginIn)
		assert.NotEmpty(t, loginOut.Error)
	})

	t.Run(`user reset password must ok`, func(t *testing.T) {
		in := &UserResetPassword_In{
			Password:   newPass,
			SecretCode: p1.SecretCode,
			Hash:       hash,
		}
		out := d.UserResetPassword(in)
		assert.Empty(t, out.Error)
	})

	t.Run(`user relogin with old password must fail`, func(t *testing.T) {
		loginIn := &UserLogin_In{Email: p1.Email, Password: p1.Email}
		loginOut := d.UserLogin(loginIn)
		assert.NotEmpty(t, loginOut.Error)
	})

	t.Run(`user login with new password must ok`, func(t *testing.T) {
		loginIn := &UserLogin_In{Email: p1.Email, Password: newPass}
		loginOut := d.UserLogin(loginIn)
		assert.Empty(t, loginOut.Error)
	})
}
