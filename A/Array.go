package A

// Array support package

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
)

// X array (slice) of anything
//
//	v := A.X{}
//	v = append(v, any_value)
type X []any

// MSX array (slice) of map with string key and any value
//
//	v := A.MSX{}
//	v = append(v, map[string]{
//	  `foo`: 123,
//	  `bar`: `yay`,
//	})
type MSX []map[string]any

// ToJson convert map array of string to JSON string type
//
//	m := []any{123,`abc`}
//	L.Print(A.ToJson(m)) // [123,"abc"]
func ToJson(arr []any) string {
	str, err := json.Marshal(arr)
	L.IsError(err, `Slice.ToJson failed`, arr)
	return string(str)
}

// ToMsgp convert map array of string to MsgPack string type
//
//	m := []any{123,`abc`}
//	L.Print(string(A.ToMsgp(m))) // �{�abc
func ToMsgp(arr []any) []byte {
	str, err := msgpack.Marshal(arr)
	L.IsError(err, `Slice.ToMsgp failed`, arr)
	return str
}

// StrJoin combine strings in the array of string with the chosen string separator
//
//	m1 := []string{`satu`,`dua`}
//	A.StrJoin(m1,`-`) // satu-dua
func StrJoin(arr []string, sep string) string {
	return strings.Join(arr, sep)
}

// IntJoin combine int64s in the array of int64 with the chosen string separator
//
//	m1 := []int64{123,456}
//	A.IntJoin(m1,`|`) // 123|456
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

// UIntJoin combine uint64s in the array of int64 with the chosen string separator
//
//	m1 := []uint64{123,456}
//	A.UIntJoin(m1,`-`) // 123-456
func UIntJoin(arr []uint64, sep string) string {
	buf := bytes.Buffer{}
	len := len(arr) - 1
	for idx, v := range arr {
		buf.WriteString(I.UToS(v))
		if idx < len {
			buf.WriteString(sep)
		}
	}
	return buf.String()
}

// StrToInt Convert array of string to array of int64
//
//	func main() {
//	  m := []string{`1`,`2`}
//	  L.Print(A.StrToInt(m))//output [1 2]
//	}
//
// convert string list to integer list
func StrToInt(arr []string) []int64 {
	res := []int64{}
	for _, v := range arr {
		if v == `` {
			continue
		}
		iv, _ := strconv.ParseInt(v, 10, 64)
		res = append(res, iv)
	}
	return res
}

// StrContains Append string to array of string if not exists
func StrContains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

// IntContains Append int64 to array of string if not exists
func IntContains(arr []int64, str int64) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

// StrAppendIfNotExists Append if not exists
func StrAppendIfNotExists(arr []string, str string) []string {
	if StrContains(arr, str) {
		return arr
	}
	return append(arr, str)
}

// IntAppendIfNotExists Append if not exists
func IntAppendIfNotExists(arr []int64, str int64) []int64 {
	if IntContains(arr, str) {
		return arr
	}
	return append(arr, str)
}

// StrsAppendIfNotExists Append slices if not exists
func StrsAppendIfNotExists(arr []string, strs []string) []string {
	for _, str := range strs {
		if StrContains(arr, str) {
			continue
		}
		arr = append(arr, str)
	}
	return arr
}

// IntsAppendIfNotExists Append slices if not exists
func IntsAppendIfNotExists(arr []int64, ints []int64) []int64 {
	for _, i := range ints {
		if IntContains(arr, i) {
			continue
		}
		arr = append(arr, i)
	}
	return arr
}

// ParseEmail split, add alias, and concat emails with name
func ParseEmail(str_emails, each_name string) []string {
	temp := []string{}
	str_email := strings.Split(str_emails, `,`)
	for _, each_email := range str_email {
		each_email = strings.TrimSpace(each_email)
		if each_email == `` {
			continue
		}
		each_name = strings.ReplaceAll(each_name, `,`, `_`)
		each_name = strings.ReplaceAll(each_name, `.`, `_`)
		each_name = strings.ReplaceAll(each_name, `<`, `_`)
		each_name = strings.ReplaceAll(each_name, `>`, `_`)
		each_name = strings.ReplaceAll(each_name, `(`, `_`)
		each_name = strings.ReplaceAll(each_name, `)`, `_`)
		temp = append(temp, each_name+`<`+each_email+`>`)
	}
	return temp
}

// FloatExist check if float exists on array
func FloatExist(arr []float64, val float64) bool {
	for _, cur := range arr {
		if val == cur {
			return true
		}
	}
	return false
}
