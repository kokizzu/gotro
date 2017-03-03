package My

import (
	_ "github.com/go-sql-driver/mysql"
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
