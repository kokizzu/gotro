package W

import (
	"math/rand"
	"time"

	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/T"
	"github.com/valyala/fasthttp"
)

type SessionConnector interface {
	// delete key
	Del(key string)
	// get remaining lifespan in seconds
	Expiry(key string) int64
	// set string with remaining lifespan in seconds
	FadeStr(key, val string, ttl int64)
	// set integer with remaining lifespan in seconds
	FadeInt(key string, val int64, ttl int64)
	// set json with remaining lifespan in seconds
	FadeMSX(key string, val M.SX, ttl int64)
	// get string
	GetStr(key string) string
	// get integer
	GetInt(key string) int64
	// get string
	GetMSX(key string) M.SX
	// increment
	Inc(key string) int64
	// increment
	Dec(key string) int64
	// set string
	SetStr(key, val string)
	// set integer
	SetInt(key string, val int64)
	// set json
	SetMSX(key string, val M.SX)
	// set json
	SetMSS(key string, val M.SS)
	// get product name, eg: Pg, Sc, My, Rd
	Product() string

	Lpush(key string, val string)
	Rpush(key string, val string)
	Lrange(key string, first, last int64) []string
}

var SESS_KEY = `SK`
var EXPIRE_SEC int64
var RENEW_SEC int64

const NS2SEC = 1000 * 1000 * 1000

type Session struct {
	UserAgent string
	IpAddr    string
	Key       string
	M.SX
	Changed bool
}

func (s *Session) Logout() {
	Sessions.Del(s.Key)
	s.Key = ``
	s.SX = M.SX{}
	s.RandomKey()
}

func (s *Session) RandomKey() {
	for {
		s.Key = s.StateCSRF() + S.RandomCB63(8)
		if Sessions.GetStr(s.Key) == `` {
			Sessions.FadeMSX(s.Key, M.SX{`key_at`: T.Epoch()}, EXPIRE_SEC)
			break // no collision
		}
	}
	//L.LOG.Notice(s.Key)
	s.Changed = true
}

func (s *Session) Login(val M.SX) {
	if s.Key == `` {
		s.RandomKey()
	}
	val[`login_at`] = T.Epoch()
	s.SX = val
	s.Changed = true
	L.Print(`LOGIN`, s.Key, val)
	Sessions.FadeMSX(s.Key, val, EXPIRE_SEC)
}

// should be called after receiving request
func (s *Session) Load(ctx *Context) {
	r := ctx.RequestCtx
	s.UserAgent = string(r.Request.Header.UserAgent())
	s.IpAddr = r.RemoteAddr().String()
	cookie := string(r.Request.Header.Cookie(SESS_KEY))
	if cookie == `` {
		s.SX = M.SX{}
		s.RandomKey()
	} else if !S.StartsWith(cookie, s.StateCSRF()) {
		s.Logout() // possible incorrect cookie stealing
	} else {
		s.Key = cookie
		s.SX = Sessions.GetMSX(s.Key)
		//L.Print(`Session.Load`, s.Key,s.SX)
		if len(s.SX) == 0 {
			s.Logout() // possible using expired cookie
		}
	}
}

// should be called before writing response
func (s *Session) Save(ctx *Context) {
	//L.Print(`Session.Save?`, s.Changed)
	if s.Changed {
		rem := Sessions.Expiry(s.Key)
		expiration := time.Now().Add(time.Second * time.Duration(rem))
		//L.Print(`Session.Save`, rem, expiration, s.Key)
		cookie := &fasthttp.Cookie{}
		cookie.SetKey(SESS_KEY)
		cookie.SetValue(s.Key)
		cookie.SetPath(`/`)
		cookie.SetExpire(expiration)
		ctx.Response.Header.SetCookie(cookie)
	}
}

func (s *Session) StateCSRF() string {
	return S.HashPassword(s.UserAgent) + `|`
}

func (s *Session) Touch() {
	if s.Key == `` {
		return
	}
	if Sessions.Expiry(s.Key) < RENEW_SEC {
		s.SX[`renew_at`] = T.Epoch()
		s.Changed = true
		Sessions.FadeMSX(s.Key, s.SX, EXPIRE_SEC)
	}
}

func (s *Session) String() string {
	if len(s.SX) == 0 {
		return ``
	}
	return s.SX.Pretty(` | `)
}

func (s *Session) NewlineString() string {
	if len(s.SX) == 0 {
		return ``
	}
	return s.SX.Pretty("\n\t")
}

func (s *Session) HeaderString() string {
	return s.GetStr(`id`) + "\n\tUserAgent: " + s.UserAgent + "\n\tSessionKey: " + s.Key
}

func InitSession(sess_key string, expire_ns, renew_ns time.Duration, conn SessionConnector, glob SessionConnector) {
	rand.Seed(T.UnixNano())
	SESS_KEY = S.IfEmpty(sess_key, SESS_KEY)
	EXPIRE_SEC = int64(expire_ns / NS2SEC)
	RENEW_SEC = int64(renew_ns / NS2SEC)
	Sessions = conn
	Globals = glob
}
