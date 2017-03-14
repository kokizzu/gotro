package A

import (
	"bytes"
	"encoding/json"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"strings"
)

// Array support package

// array (slice) of anything
//  v := A.X{}
//  v = append(v, any_value)
type X []interface{}

// array (slice) of map with string key and any value
//  v := A.MSX{}
//  v = append(v, map[string]{
//    `foo`: 123,
//    `bar`: `yay`,
//  })
type MSX []map[string]interface{}

// convert map array of string to JSON string type
//  m:= []interface{}{123,`abc`}
//  L.Print(A.ToJson(m)) // [123,"abc"]
func ToJson(arr []interface{}) string {
	str, err := json.Marshal(arr)
	L.IsError(err, `Slice.ToJson failed`, arr)
	return string(str)
}

// combine strings in the array of string with the chosen string separator
//  m1:= []string{`satu`,`dua`}
//  A.StrJoin(m1,`-`) // satu-dua
func StrJoin(arr []string, sep string) string {
	return strings.Join(arr, sep)
}

// combine int64s in the array of int64 with the chosen string separator
//  m1:= []int64{123,456}
//  A.IntJoin(m1,`-`) // 123-456
func IntJoin(arr []int64, sep string) string {
	buf := bytes.Buffer{}
	len := len(arr) - 1
	for idx, v := range arr {
		buf.WriteString(I.ToS(v))
		if idx < len {
			buf.WriteString(sep)
		}
	}
	return buf.String()
}
