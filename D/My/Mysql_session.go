package My

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kokizzu/gotro/M"
)

/*
// create new mysql connection
func NewMyConn(user, db, ip, pass string) *RDBMS {
	opt := user + `:` + pass + `@` + ip + `/` + db
	conn := sqlx.MustConnect(`mysql`, opt)
	name := `my::` + opt
	return &RDBMS{
		Name:    name,
		Adapter: conn,
	}
}
*/

type MysqlSession struct {
}

func (sess *MysqlSession) Del(key string) {
	// TODO: continue this
}

func (sess *MysqlSession) Expiry(key string) int64 {
	// TODO: continue this
	return 0
}

func (sess *MysqlSession) FadeStr(key, val string, sec int64) {
	// TODO: continue this
}

func (sess *MysqlSession) FadeInt(key string, val int64, sec int64) {
	// TODO: continue this
}

func (sess *MysqlSession) FadeMSX(key string, val M.SX, sec int64) {
	// TODO: continue this
}

func (sess *MysqlSession) GetStr(key string) string {
	// TODO: continue this
	return ``
}

func (sess *MysqlSession) GetInt(key string) int64 {
	// TODO: continue this
	return 0
}

func (sess *MysqlSession) GetMSX(key string) M.SX {
	// TODO: continue this
	res := M.SX{}
	return res
}

func (sess *MysqlSession) Inc(key string) int64 {
	// TODO: continue this
	return 0
}

func (sess *MysqlSession) SetStr(key, val string) {
	// TODO: continue this
}

func (sess *MysqlSession) SetInt(key string, val int64) {
	// TODO: continue this
}

func (sess *MysqlSession) SetMSX(key string, val M.SX) {
	// TODO: continue this
}
