package D

import (
	"bytes"
	"github.com/kokizzu/gotro/S"
)

type Record interface {
	GetStr(string) string
	GetFloat(string) float64
	GetInt(string) int64
	GetArr(string) []interface{}
	GetBool(string) bool
}

const REDIS = `Rd`
const SCYLLA = `Sc`
const POSTGRE = `Pg`
const DUMMY = `Du`
const ARANGO = `Ar`
const AEROSP = `As`

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

func WhereIn(vals []string) string {
	len := len(vals) - 1
	if len < 0 {
		return ` IN ('')`
	}
	buf := bytes.Buffer{}
	buf.WriteString(` IN (`)
	for k, v := range vals {
		buf.WriteString(S.Z(v))
		if k < len { // write except the last one
			buf.WriteRune(',')
		}
	}
	buf.WriteString(`)`)
	return buf.String()
}

func WhereInStrIds(ids []string) string {
	len := len(ids) - 1
	if len < 0 {
		return ` IN ('0')` // make sure there are no zero-value id
	}
	buf := bytes.Buffer{}
	buf.WriteString(` IN (`)
	for k, v := range ids {
		buf.WriteString(S.Z(v))
		if k < len { // write except the last one
			buf.WriteRune(',')
		}
	}
	buf.WriteString(`)`)
	return buf.String()
}

func WhereInIds(ids []int64) string {
	len := len(ids) - 1
	if len < 0 {
		return ` IN ('0')` // make sure there are no zero-value id
	}
	buf := bytes.Buffer{}
	buf.WriteString(` IN (`)
	for k, v := range ids {
		buf.WriteString(S.ZI(v))
		if k < len { // write except the last one
			buf.WriteRune(',')
		}
	}
	buf.WriteString(`)`)
	return buf.String()
}
