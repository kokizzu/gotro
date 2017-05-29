package model

import (
	"github.com/kokizzu/gotro/D/Pg"
)

var PG_W, PG_R *Pg.RDBMS

func init() {
	PG_W = Pg.NewConn(`test1`, `test1`)
	// ^ later when scaling we replace this one
	PG_R = Pg.NewConn(`test1`, `test1`)
}
