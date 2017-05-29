package D

import (
	"github.com/kokizzu/gotro/S"
)

type Record interface {
	GetStr(string) string
	GetFloat(string) float64
	GetInt(string) int64
	GetArr(string) []interface{}
	GetBool(string) bool
}

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
	ZJ = S.ZJJ
	ZI = S.ZI
	ZLIKE = S.ZLIKE
	ZS = S.ZS
}
