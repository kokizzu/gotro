package Sq

import (
	"github.com/kokizzu/gotro/M"
	_ "github.com/mutecomm/go-sqlcipher"
)

/*

// create new sqlite database
func NewSqConn(filename, password string) *RDBMS {
	conn := sqlx.MustConnect(`sqlite3`, `file:`+filename)
	db := &RDBMS{
		Name:    `sq::` + filename,
		Adapter: conn,
	}
	conn.MustExec(`PRAGMA key = ` + Z(password))
	conn.MustExec(`PRAGMA cipher_page_size = 4096`)
	return db
}
*/

type SqliteSession struct {
}

func (sess *SqliteSession) Del(key string) {
	// TODO: continue this
}

func (sess *SqliteSession) Expiry(key string) int64 {
	// TODO: continue this
	return 0
}

func (sess *SqliteSession) FadeStr(key, val string, sec int64) {
	// TODO: continue this
}

func (sess *SqliteSession) FadeInt(key string, val int64, sec int64) {
	// TODO: continue this
}

func (sess *SqliteSession) FadeMSX(key string, val M.SX, sec int64) {
	// TODO: continue this
}

func (sess *SqliteSession) GetStr(key string) string {
	// TODO: continue this
	return ``
}

func (sess *SqliteSession) GetInt(key string) int64 {
	// TODO: continue this
	return 0
}

func (sess *SqliteSession) GetMSX(key string) M.SX {
	// TODO: continue this
	res := M.SX{}
	return res
}

func (sess *SqliteSession) Inc(key string) int64 {
	// TODO: continue this
	return 0
}

func (sess *SqliteSession) SetStr(key, val string) {
	// TODO: continue this
}

func (sess *SqliteSession) SetInt(key string, val int64) {
	// TODO: continue this
}

func (sess *SqliteSession) SetMSX(key string, val M.SX) {
	// TODO: continue this
}
