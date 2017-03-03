package Sq

import (
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
