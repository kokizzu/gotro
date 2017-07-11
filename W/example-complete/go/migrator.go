package main

import (
	"example-complete/sql"
	"github.com/kokizzu/gotro/D/Pg"
)

/*
 */

func main() {
	Pg.InitFunctions(sql.PG)
	sql.PG.CreateBaseTable(`users`, `users`)
	sql.PG.CreateBaseTable(`tags`, `users`)
}
