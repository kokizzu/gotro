package Pg

import (
	"database/sql"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	// TODO: replace with faster one `github.com/jackc/pgx/stdlib`, see https://github.com/jackc/pgx/issues/81
	// https://jmoiron.github.io/sqlx/
	// https://github.com/jmoiron/sqlx
	// https://sourcegraph.com/github.com/jmoiron/sqlx
)

func rowsAffected(rs sql.Result) int64 {
	ra, err := rs.RowsAffected()
	if L.IsError(err, `failed to get rows affected`, ra) {
		return 0
	}
	return ra
}

var Z func(string) string
var ZZ func(string) string
var ZJ func(string) string
var ZI func(int64) string
var ZLIKE func(string) string
var ZS func(string) string

func init() {
	Z = S.Z
	ZZ = S.ZZ
	ZJ = S.ZJ
	ZI = S.ZI
	ZLIKE = S.ZLIKE
	ZS = S.ZS
}

var DEBUG bool

type Record interface {
	Get_Str(string) string
	Get_Float(string) float64
	Get_Int(string) int64
	Get_Arr(string) []interface{}
	Get_Bool(string) bool
	Get_Id() int64
	Get_UniqueId() string
}
