package W

import (
	"github.com/kokizzu/gotro/M"
	"time"
)

type Session struct {
	UserId    int64
	AppId     int64
	Level     M.SX
	Email     string
	UserAgent string
	// TODO: continue this
}

func (s *Session) Logout() {
	// TODO: continue this
}

func (s *Session) Login(id int64, email string) {
	// TODO: continue this
}

func (s *Session) Get_GlobalTTL(key string) int64 {
	// TODO: continue this
	return 0
}

func (s *Session) Set_GlobalTTL(key, val string, ttl int64) {
	// TODO: continue this
}
func (s *Session) Get_Global(key string) string {
	// TODO: continue this
	return ``
}

func (s *Session) Inc_Global(key string) int64 {
	// TODO: continue this
	return 0
}
func (s *Session) Set_Global(key, val string) {
	// TODO: continue this
}

func (s *Session) Del_Global(key string) {
	// TODO: continue this
}

func (s *Session) StateCSRF() string {
	// TODO: continue this
	return ``
}

func (s *Session) Touch() {
	// TODO: update ttl
}

func InitSession(expire_ns, renew_ns time.Duration, dbNum int) {
	// TODO: continue this
}
