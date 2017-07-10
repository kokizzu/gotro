package main

import (
	"example-complete/sql"
	"example-complete/sql/nAccounts"
	"example-complete/sql/nCompanies"
	"example-complete/sql/nSensors"
	"github.com/kokizzu/gotro/D/Pg"
)

/*
 */

func main() {
	Pg.InitFunctions(sql.PG)
	sql.PG.CreateBaseTable(`users`, `users`)
	sql.PG.CreateBaseTable(`tags`, `users`)
}
