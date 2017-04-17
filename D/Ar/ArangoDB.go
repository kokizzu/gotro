package Ar

import (
	"github.com/kokizzu/gotro/D"
	"github.com/kokizzu/gotro/S"
)

var Z func(string) string
var ZZ func(string) string
var ZJ func(string) string
var ZI func(int64) string
var ZLIKE func(string) string
var ZS func(string) string

var DEBUG bool

func init() {
	Z = S.Z
	ZZ = S.ZZ
	ZJ = S.ZJ
	ZI = S.ZI
	ZLIKE = S.ZLIKE
	ZS = S.ZS
	DEBUG = D.DEBUG
}

const SQL_FUNCTIONS = ``

/*
func InitFunctions(conn *RDBMS) {
	conn.InitTrigger()
	conn.DoTransaction(func(tx *Tx) string {
		tx.DoExec(SQL_FUNCTIONS)
		return ``
	})
}
*/
