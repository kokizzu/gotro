package domain

import (
	"testing"
	"time"

	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/W2"
	"github.com/kokizzu/gotro/W2/example/conf"
	"github.com/kokizzu/gotro/W2/example/model/mAuth/wcAuth"
	"github.com/kokizzu/id64"
	"github.com/kokizzu/lexid"
	"github.com/kpango/fastime"
	"github.com/stretchr/testify/assert"
)

const testDomain = `@localhost`

func TestMain(t *testing.M) {
	W2.LoadTestEnv()

	// ensure admin account exists
	d := NewDomain()
	player := wcAuth.NewUsersMutator(d.Taran)
	player.Id = 1
	if !player.FindById() {
		player.Email = conf.SuperAdmin
		player.SetEncryptPassword(player.Email)
		player.CreatedAt = fastime.UnixNow()
		if !player.DoUpdateById() {
			panic(`cannot create superadmin`)
		}
	}

	t.Run()
}

func TestDomain_PlayerLoginRegisterFlow(t *testing.T) {
	d := NewDomain()
	name := id64.ID().String()
	pass := lexid.ID()
	email := name + testDomain
	t.Run(`register should ok`, func(t *testing.T) {
		in := &PlayerRegister_In{
			Email:    email,
			Password: pass,
			UserName: name,
		}
		out := d.PlayerRegister(in)
		assert.Empty(t, out.Error)
	})

	t.Run(`reregister same email should fail`, func(t *testing.T) {
		in := &PlayerRegister_In{
			Email:    email,
			Password: pass,
			UserName: name + `1`,
		}
		out := d.PlayerRegister(in)
		assert.NotEmpty(t, out.Error)
		assert.NotEqual(t, 200, out.StatusCode)
		assert.NotEqual(t, 0, out.StatusCode)
	})

	t.Run(`reregister same userName should fail`, func(t *testing.T) {
		in := &PlayerRegister_In{
			Email:    email + `1`,
			Password: pass,
			UserName: name,
		}
		out := d.PlayerRegister(in)
		assert.NotEmpty(t, out.Error)
		assert.NotEqual(t, 200, out.StatusCode)
		assert.NotEqual(t, 0, out.StatusCode)
	})

	t.Run(`login with unregistered user should fail`, func(t *testing.T) {
		in := &PlayerLogin_In{
			Email:    name + `notExists` + testDomain,
			Password: pass,
		}
		out := d.PlayerLogin(in)
		assert.NotEmpty(t, out.Error)
	})

	t.Run(`login with wrong password should fail`, func(t *testing.T) {
		in := &PlayerLogin_In{
			Email:    email,
			Password: name,
		}
		out := d.PlayerLogin(in)
		assert.NotEmpty(t, out.Error)
	})

	sessionToken := ``
	t.Run(`login with correct password should ok`, func(t *testing.T) {
		in := &PlayerLogin_In{
			Email:    email,
			Password: pass,
		}
		out := d.PlayerLogin(in)
		assert.Empty(t, out.Error)
		assert.NotEmpty(t, out.SessionToken)
		sessionToken = out.SessionToken
	})

	t.Run(`check profile with active session should ok`, func(t *testing.T) {
		in := &PlayerProfile_In{NewRC(sessionToken)}
		out := d.PlayerProfile(in)
		assert.Empty(t, out.Error)
		assert.NotNil(t, out.Player)
		if out.Player == nil {
			t.Failed()
			return
		}
		assert.Equal(t, email, out.Player.Email)
	})

	t.Run(`change password with wrong password must fail`, func(t *testing.T) {
		in := &PlayerChangePassword_In{
			RequestCommon: NewRC(sessionToken),
			Password:      ``,
			NewPassword:   `abc`,
		}
		out := d.PlayerChangePassword(in)
		assert.NotEmpty(t, out.Error)
	})

	t.Run(`change password with correct password must ok`, func(t *testing.T) {
		in := &PlayerChangePassword_In{
			RequestCommon: NewRC(sessionToken),
			Password:      pass,
			NewPassword:   email,
		}
		out := d.PlayerChangePassword(in)
		assert.Empty(t, out.Error)
	})

	t.Run(`change userName with unique username must ok`, func(t *testing.T) {
		in := &PlayerUpdateProfile_In{
			RequestCommon: NewRC(sessionToken),
			UserName:      name + `-dummy`,
		}
		out := d.PlayerUpdateProfile(in)
		assert.Empty(t, out.Error)
	})

	t.Run(`change userName with existing username must fail`, func(t *testing.T) {
		in := &PlayerUpdateProfile_In{
			RequestCommon: NewRC(sessionToken),
			UserName:      S.LeftOf(conf.SuperAdmin, `@`),
		}
		out := d.PlayerUpdateProfile(in)
		assert.NotEmpty(t, out.Error)
	})

	t.Run(`logout should ok`, func(t *testing.T) {
		in := &PlayerLogout_In{RequestCommon{SessionToken: sessionToken}}
		out := d.PlayerLogout(in)
		assert.Equal(t, out.LoggedOut, true)
	})

	t.Run(`check profile with expired session should fail`, func(t *testing.T) {
		in := &PlayerProfile_In{NewRC(sessionToken)}
		out := d.PlayerProfile(in)
		assert.Equal(t, 403, out.StatusCode)
	})
}

func TestDomain_PlayerList(t *testing.T) {
	d := NewDomain()

	// if fail, then probably no data
	t.Run(`list player must ok`, func(t *testing.T) {
		in := &PlayerList_In{
			Limit:  2,
			Offset: 0,
		}
		out := d.PlayerList(in)
		assert.Empty(t, out.Error)
		assert.NotEmpty(t, out.Players)
		assert.Greater(t, len(out.Players), 1)
		L.Describe(out.Players)
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

func TestDomain_PlayerForgotReset(t *testing.T) {
	d := NewDomain()
	p1, _ := dummyUser(d)
	hash := S.EncodeCB63(int64(p1.Id), 1)
	newPass := `12345678`

	t.Run(`player forgot password must ok`, func(t *testing.T) {
		in := &PlayerForgotPassword_In{
			Email: p1.Email,
		}
		out := d.PlayerForgotPassword(in)
		time.Sleep(1 * time.Second)
		assert.Empty(t, out.Error)
	})

	t.Run(`player login with old password must ok`, func(t *testing.T) {
		loginIn := &PlayerLogin_In{Email: p1.Email, Password: p1.Email}
		loginOut := d.PlayerLogin(loginIn)
		assert.Empty(t, loginOut.Error)
	})

	t.Run(`player forgot password with invalid email must fail`, func(t *testing.T) {
		in := &PlayerForgotPassword_In{
			Email: id64.ID().String() + testDomain,
		}
		out := d.PlayerForgotPassword(in)
		assert.NotEmpty(t, out.Error)
	})

	p1.FindById() // to get SecretCode

	t.Run(`player reset password with incorrect secretCode must fail`, func(t *testing.T) {
		in := &PlayerResetPassword_In{
			Password:   newPass,
			SecretCode: `duar`,
			Hash:       hash,
		}
		out := d.PlayerResetPassword(in)
		assert.NotEmpty(t, out.Error)
	})

	t.Run(`player login with new password must fail`, func(t *testing.T) {
		loginIn := &PlayerLogin_In{Email: p1.Email, Password: newPass}
		loginOut := d.PlayerLogin(loginIn)
		assert.NotEmpty(t, loginOut.Error)
	})

	t.Run(`player reset password with incorrect hash must fail`, func(t *testing.T) {
		in := &PlayerResetPassword_In{
			Password:   newPass,
			SecretCode: p1.SecretCode,
			Hash:       `-`,
		}
		out := d.PlayerResetPassword(in)
		assert.NotEmpty(t, out.Error)
	})

	t.Run(`player login with new password must fail`, func(t *testing.T) {
		loginIn := &PlayerLogin_In{Email: p1.Email, Password: newPass}
		loginOut := d.PlayerLogin(loginIn)
		assert.NotEmpty(t, loginOut.Error)
	})

	t.Run(`player reset password must ok`, func(t *testing.T) {
		in := &PlayerResetPassword_In{
			Password:   newPass,
			SecretCode: p1.SecretCode,
			Hash:       hash,
		}
		out := d.PlayerResetPassword(in)
		assert.Empty(t, out.Error)
	})

	t.Run(`player relogin with old password must fail`, func(t *testing.T) {
		loginIn := &PlayerLogin_In{Email: p1.Email, Password: p1.Email}
		loginOut := d.PlayerLogin(loginIn)
		assert.NotEmpty(t, loginOut.Error)
	})

	t.Run(`player login with new password must ok`, func(t *testing.T) {
		loginIn := &PlayerLogin_In{Email: p1.Email, Password: newPass}
		loginOut := d.PlayerLogin(loginIn)
		assert.Empty(t, loginOut.Error)
	})
}
