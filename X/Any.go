package X

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/goccy/go-yaml"
	"strconv"
	"strings"
	"time"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/B"
	"github.com/kokizzu/gotro/C"
	"github.com/kokizzu/gotro/F"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
)

// Any type support package

// convert any data type to uint
//  var m interface{}
//  m = `123`
//  L.ParentDescribe(X.ToI(m)) // uint(123)
func ToU(any interface{}) uint64 {
	if any == nil {
		return 0
	}
	if val, ok := any.(uint64); ok {
		return val
	}
	switch v := any.(type) {
	case int:
		return uint64(v)
	case uint:
		return uint64(v)
	case int8:
		return uint64(v)
	case int16:
		return uint64(v)
	case int32:
		return uint64(v)
	case int64:
		return uint64(v)
	case uint8:
		return uint64(v)
	case uint16:
		return uint64(v)
	case uint32:
		return uint64(v)
	case uint64:
		return v
	case float32:
		return uint64(v)
	case float64:
		return uint64(v)
	case *int:
		if v != nil {
			return uint64(*v)
		}
	case *uint:
		if v != nil {
			return uint64(*v)
		}
	case *int8:
		if v != nil {
			return uint64(*v)
		}
	case *int16:
		if v != nil {
			return uint64(*v)
		}
	case *int32:
		if v != nil {
			return uint64(*v)
		}
	case *int64:
		if v != nil {
			return uint64(*v)
		}
	case *uint8:
		if v != nil {
			return uint64(*v)
		}
	case *uint16:
		if v != nil {
			return uint64(*v)
		}
	case *uint32:
		if v != nil {
			return uint64(*v)
		}
	case *uint64:
		if v != nil {
			return *v
		}
	case *float32:
		if v != nil {
			return uint64(*v)
		}
	case *float64:
		if v != nil {
			return uint64(*v)
		}
	case bool:
		if v {
			return 1
		}
		return 0
	case []byte:
		if val, err := strconv.ParseInt(string(v), 10, 64); err == nil {
			return uint64(val)
		}
		if val, err := strconv.ParseFloat(string(v), 64); err == nil {
			return uint64(val)
		}
		L.ParentDescribe(`Can't convert to uint64`, any)
	case string:
		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			return uint64(val)
		}
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			return uint64(val)
		}
		L.ParentDescribe(`Can't convert to uint64`, any)
	default:
		L.ParentDescribe(`Can't convert to uint64`, any)
	}
	return 0
}

// convert any data type to int8
//  var m interface{}
//  m = `123`
//  L.ParentDescribe(X.ToByte(m)) // byte(123)
func ToByte(any interface{}) byte {
	if any == nil {
		return 0
	}
	switch v := any.(type) {
	case int:
		return byte(v)
	case uint:
		return byte(v)
	case int8:
		return byte(v)
	case int16:
		return byte(v)
	case int32:
		return byte(v)
	case int64:
		return byte(v)
	case uint8:
		return v
	case uint16:
		return byte(v)
	case uint32:
		return byte(v)
	case uint64:
		return byte(v)
	case float32:
		return byte(v)
	case float64:
		return byte(v)
	case *int:
		if v != nil {
			return byte(*v)
		}
	case *uint:
		if v != nil {
			return byte(*v)
		}
	case *int8:
		if v != nil {
			return byte(*v)
		}
	case *int16:
		if v != nil {
			return byte(*v)
		}
	case *int32:
		if v != nil {
			return byte(*v)
		}
	case *int64:
		if v != nil {
			return byte(*v)
		}
	case *uint8:
		if v != nil {
			return *v
		}
	case *uint16:
		if v != nil {
			return byte(*v)
		}
	case *uint32:
		if v != nil {
			return byte(*v)
		}
	case *uint64:
		if v != nil {
			return byte(*v)
		}
	case *float32:
		if v != nil {
			return byte(*v)
		}
	case *float64:
		if v != nil {
			return byte(*v)
		}
	case bool:
		if v {
			return 1
		}
		return 0
	case []byte:
		if val, err := strconv.ParseInt(string(v), 10, 64); err == nil {
			return byte(val)
		}
		if val, err := strconv.ParseFloat(string(v), 64); err == nil {
			return byte(val)
		}
		L.ParentDescribe(`Can't convert to byte`, any)
	case string:
		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			return byte(val)
		}
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			return byte(val)
		}
		L.ParentDescribe(`Can't convert to byte`, any)
	default:
		L.ParentDescribe(`Can't convert to byte`, any)
	}
	return 0
}

// convert any data type to int64
//  var m interface{}
//  m = `123`
//  L.ParentDescribe(X.ToI(m)) // int64(123)
func ToI(any interface{}) int64 {
	if any == nil {
		return 0
	}
	if val, ok := any.(int64); ok {
		return val
	}
	switch v := any.(type) {
	case int:
		return int64(v)
	case uint:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case *int:
		if v != nil {
			return int64(*v)
		}
	case *uint:
		if v != nil {
			return int64(*v)
		}
	case *int8:
		if v != nil {
			return int64(*v)
		}
	case *int16:
		if v != nil {
			return int64(*v)
		}
	case *int32:
		if v != nil {
			return int64(*v)
		}
	case *int64:
		if v != nil {
			return *v
		}
	case *uint8:
		if v != nil {
			return int64(*v)
		}
	case *uint16:
		if v != nil {
			return int64(*v)
		}
	case *uint32:
		if v != nil {
			return int64(*v)
		}
	case *uint64:
		if v != nil {
			return int64(*v)
		}
	case *float32:
		if v != nil {
			return int64(*v)
		}
	case *float64:
		if v != nil {
			return int64(*v)
		}
	case bool:
		if v {
			return 1
		}
		return 0
	case []byte:
		if val, err := strconv.ParseInt(string(v), 10, 64); err == nil {
			return val
		}
		if val, err := strconv.ParseFloat(string(v), 64); err == nil {
			return int64(val)
		}
		L.ParentDescribe(`Can't convert to int64`, any)
	case string:
		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			return val
		}
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			return int64(val)
		}
		L.ParentDescribe(`Can't convert to int64`, any)
	default:
		L.ParentDescribe(`Can't convert to int64`, any)
	}
	return 0
}

// Convert any data type to float64
//  var m interface{}
//  m = `123.5`
//  L.ParentDescribe(X.ToF(m)) // float64(123.5)
func ToF(any interface{}) float64 {
	if any == nil {
		return 0
	}
	if val, ok := any.(float64); ok {
		return val
	}
	switch v := any.(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case *int:
		if v != nil {
			return float64(*v)
		}
	case *uint:
		if v != nil {
			return float64(*v)
		}
	case *int8:
		if v != nil {
			return float64(*v)
		}
	case *int16:
		if v != nil {
			return float64(*v)
		}
	case *int32:
		if v != nil {
			return float64(*v)
		}
	case *int64:
		if v != nil {
			return float64(*v)
		}
	case *uint8:
		if v != nil {
			return float64(*v)
		}
	case *uint16:
		if v != nil {
			return float64(*v)
		}
	case *uint32:
		if v != nil {
			return float64(*v)
		}
	case *uint64:
		if v != nil {
			return float64(*v)
		}
	case *float32:
		if v != nil {
			return float64(*v)
		}
	case *float64:
		if v != nil {
			return float64(*v)
		}
	case bool:
		if v {
			return 1
		}
		return 0
	case []byte:
		if val, err := strconv.ParseFloat(string(v), 64); err == nil {
			return val
		}
		L.ParentDescribe(`Can't convert to float64`, any)
	case string:
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			return val
		}
		L.ParentDescribe(`Can't convert to float64`, any)
	default:
		L.ParentDescribe(`Can't convert to float64`, any)
	}
	return 0
}

// convert any data type to string
//  var m interface{}
//  m = `123`
//  L.ParentDescribe(X.ToS(m)) // `123`
func ToS(any interface{}) string {
	if any == nil {
		return ``
	}
	if val, ok := any.(string); ok {
		return val
	}
	if val, ok := any.([]byte); ok {
		return string(val)
	}
	switch v := any.(type) {
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatInt(int64(v), 10)
	case uint8:
		return strconv.FormatInt(int64(v), 10)
	case uint16:
		return strconv.FormatInt(int64(v), 10)
	case uint32:
		return strconv.FormatInt(int64(v), 10)
	case uint64:
		return strconv.FormatInt(int64(v), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)

	case *int:
		if v != nil {
			return strconv.FormatInt(int64(*v), 10)
		}
	case *uint:
		if v != nil {
			return strconv.FormatInt(int64(*v), 10)
		}
	case *int8:
		if v != nil {
			return strconv.FormatInt(int64(*v), 10)
		}
	case *int16:
		if v != nil {
			return strconv.FormatInt(int64(*v), 10)
		}
	case *int32:
		if v != nil {
			return strconv.FormatInt(int64(*v), 10)
		}
	case *int64:
		if v != nil {
			return strconv.FormatInt(int64(*v), 10)
		}
	case *uint8:
		if v != nil {
			return strconv.FormatInt(int64(*v), 10)
		}
	case *uint16:
		if v != nil {
			return strconv.FormatInt(int64(*v), 10)
		}
	case *uint32:
		if v != nil {
			return strconv.FormatInt(int64(*v), 10)
		}
	case *uint64:
		if v != nil {
			return strconv.FormatInt(int64(*v), 10)
		}
	case *float32:
		if v != nil {
			return strconv.FormatFloat(float64(*v), 'f', -1, 64)
		}
	case *float64:
		if v != nil {
			return strconv.FormatFloat(float64(*v), 'f', -1, 64)
		}
	case bool:
		if v {
			return `true`
		}
		return `false`
	case fmt.Stringer:
		if v == nil {
			return ``
		}
		return v.String()
	default:
		return ToJson5(v)
	}
	return ``
}

// convert any to time
func ToTime(any interface{}) time.Time {
	if any == nil {
		return time.Time{}
	}
	switch v := any.(type) {
	case time.Time:
		return v
	case *time.Time:
		if v != nil {
			return *v
		}
	case []byte: // "YYYY-MM-DD HH:MM:SS.MMMMMM"
		res, err := parseDateTime(v, time.UTC)
		L.IsError(err, `Can't convert to time.Time from []byte`, any)
		return res
	case string: // "YYYY-MM-DD HH:MM:SS.MMMMMM"
		res, err := parseDateTime([]byte(v), time.UTC)
		L.IsError(err, `Can't convert to time.Time from string`, any)
		return res
	default:
		L.CheckIf(false, `Can't convert to time.Time`, any)
	}
	return time.Time{}
}

// convert any data type to bool
//  var m interface{}
//  m = `123`
//  L.ParentDescribe(X.ToBool(m)) // bool(true)
func ToBool(any interface{}) bool {
	if any == nil {
		return false
	}
	if val, ok := any.(bool); ok {
		return val
	}
	switch v := any.(type) {
	case int:
		return v != 0
	case int8:
		return v != 0
	case int16:
		return v != 0
	case int32:
		return v != 0
	case int64:
		return v != 0
	case uint:
		return v != 0
	case uint8:
		return v != 0
	case uint16:
		return v != 0
	case uint32:
		return v != 0
	case uint64:
		return v != 0
	case float32:
		return v != 0
	case float64:
		return v != 0
	case *int:
		if v != nil {
			return *v == 0
		}
	case *uint:
		if v != nil {
			return *v == 0
		}
	case *int8:
		if v != nil {
			return *v == 0
		}
	case *int16:
		if v != nil {
			return *v == 0
		}
	case *int32:
		if v != nil {
			return *v == 0
		}
	case *int64:
		if v != nil {
			return *v == 0
		}
	case *uint8:
		if v != nil {
			return *v == 0
		}
	case *uint16:
		if v != nil {
			return *v == 0
		}
	case *uint32:
		if v != nil {
			return *v == 0
		}
	case *uint64:
		if v != nil {
			return *v == 0
		}
	case *float32:
		if v != nil {
			return *v == 0
		}
	case *float64:
		if v != nil {
			return *v == 0
		}
	case fmt.Stringer:
		val := v.String()
		val = strings.TrimSpace(strings.ToLower(val))
		return !(val == `` || val == `0` || val == `f` || val == `false` || val == `no` || val == `n`)
	case string:
		val := v
		val = strings.TrimSpace(strings.ToLower(val))
		return !(val == `` || val == `0` || val == `f` || val == `false` || val == `no` || val == `n`)
	default:
		L.ParentDescribe(`Can't convert to string`, v)
	}
	return false
}

// convert any data type to array of any
//  var m3 interface{}
//  m3 = []interface{}{1}   // tipe array
//  L.ParentDescribe(X.ToArr(m3)) // []interface {}{int(1),}
func ToArr(any interface{}) []interface{} {
	if any == nil {
		return []interface{}{}
	}
	val, ok := any.([]interface{})
	if L.CheckIf(!ok, `Can't convert to []interface{}`, any) {
		return []interface{}{}
	}
	return val
}

// convert array of any data type to array of string
//  var m4 []interface{}
//  m4 = []interface{}{1}     // // tipe array
//  L.ParentDescribe(X.ArrToStrArr(m4)) // []string{"1"}
func ArrToStrArr(any_arr []interface{}) []string {
	res := []string{}
	for _, val := range any_arr {
		res = append(res, ToS(val))
	}
	return res
}

// Convert array of any data type to array of int64
//  var m4 []interface{}
//  m4 = []interface{}{1}     // // tipe array
//  L.ParentDescribe(X.ArrToIntArr(m4)) // []int64{1}
func ArrToIntArr(any_arr []interface{}) []int64 {
	res := []int64{}
	for _, val := range any_arr {
		res = append(res, ToI(val))
	}
	return res
}

func json5fromMIB(orig map[int64]bool) string {
	b := bytes.Buffer{}
	b.WriteByte('{')
	first := true
	for k, v := range orig {
		if !first {
			b.WriteByte(',')
		} else {
			first = false
		}
		b.WriteString(I.ToS(k))
		b.WriteByte(':')
		b.WriteString(ToJson5(v))
	}
	b.WriteByte('}')
	return b.String()
}

func json5fromMIX(orig map[int64]interface{}) string {
	b := bytes.Buffer{}
	b.WriteByte('{')
	first := true
	for k, v := range orig {
		if !first {
			b.WriteByte(',')
		} else {
			first = false
		}
		b.WriteString(I.ToS(k))
		b.WriteByte(':')
		b.WriteString(ToJson5(v))
	}
	b.WriteByte('}')
	return b.String()
}

func json5fromMIAX(orig map[int64][]interface{}) string {
	b := bytes.Buffer{}
	b.WriteByte('{')
	first := true
	for k, v := range orig {
		if !first {
			b.WriteByte(',')
		} else {
			first = false
		}
		b.WriteString(I.ToS(k))
		b.WriteByte(':')
		b.WriteString(ToJson5(v))
	}
	b.WriteByte('}')
	return b.String()
}

func json5fromMSAX(orig map[string][]interface{}) string {
	b := bytes.Buffer{}
	b.WriteByte('{')
	first := true
	for k, v := range orig {
		if !first {
			b.WriteByte(',')
		} else {
			first = false
		}
		b.WriteString(S.ZZ(k))
		b.WriteByte(':')
		b.WriteString(ToJson5(v))
	}
	b.WriteByte('}')
	return b.String()
}

func json5fromMSI(orig map[string]int64) string {
	b := bytes.Buffer{}
	b.WriteByte('{')
	first := true
	for k, v := range orig {
		if !first {
			b.WriteByte(',')
		} else {
			first = false
		}
		quote := true
		if len(k) > 0 {
			ch := k[0]
			if C.IsDigit(ch) && ch != '0' {
				for _, ch := range k[1:] {
					// find non digit
					if !C.IsDigit(uint8(ch)) {
						quote = true
						break
					}
				}
			} else if C.IsIdentStart(k[0]) {
				for _, ch := range k[1:] {
					// find non identifier character
					if !C.IsIdent(uint8(ch)) {
						quote = true
						break
					}
				}
			} else {
				quote = true
			}
		}
		if quote {
			k = S.Q(k)
		}
		b.WriteString(k)
		b.WriteByte(':')
		b.WriteString(I.ToS(v))
	}
	b.WriteByte('}')
	return b.String()
}

func ToJson5(any interface{}) string {
	// bug when using map[int64]interface{}
	if any == nil {
		return `''`
	}
	switch orig := any.(type) {
	case bytes.Buffer: // return as is
		return orig.String()
	case string:
		return S.ZJJ(orig)
	case []byte:
		return S.ZJJ(string(orig))
	case int, int64, int32:
		return I.ToS(any.(int64))
	case float32, float64:
		return F.ToS(any.(float64))
	case bool:
		return B.ToS(orig)
	case M.IB:
		return json5fromMIB(orig)
	case map[int64]bool:
		return json5fromMIB(orig)
	case M.IX:
		return json5fromMIX(orig)
	case map[int64]interface{}:
		return json5fromMIX(orig)
	case M.IAX:
		return json5fromMIAX(orig)
	case map[int64][]interface{}:
		return json5fromMIAX(orig)
	case M.SAX:
		return json5fromMSAX(orig)
	case map[string][]interface{}:
		return json5fromMSAX(orig)
	case M.SX:
		return orig.ToJson()
	case map[string]interface{}:
		return M.ToJson(orig)
	//   return any.(M.SX).ToJson()
	case M.SI:
		return json5fromMSI(orig)
	case map[string]int64:
		return json5fromMSI(orig)
	case A.X:
		return A.ToJson(orig)
	case []interface{}:
		return A.ToJson(orig)
	default:
		str, err := json.Marshal(any)
		L.IsError(err, `X.ToJson5 failed`, any)
		return string(str)
	}
	// TODO: add more types (M/A) here, do not EVER TRY to use reflection in this case
}

// convert to beautiful json text
//  m:= []interface {}{true,`1`,23,`wabcd`}
//  L.Print(K.ToJsonPretty(m))
//  // [
//  //   true,
//  //   "1",
//  //   23,
//  //   "wabcd"
//  // ]
func ToJsonPretty(any interface{}) string {
	res, err := json.MarshalIndent(any, ``, `  `)
	L.IsError(err, `X.ToJsonPretty failed`, any)
	return string(res)
}

// convert to standard json text
func ToJson(any interface{}) string {
	res, err := json.Marshal(any)
	L.IsError(err, `X.ToJson failed`, any)
	return string(res)
}

func ToAX(any interface{}) A.X {
	if any == nil {
		return A.X{}
	}
	val, ok := any.([]interface{})
	if L.CheckIf(!ok, `Can't convert to A.X`, any) {
		return A.X{}
	}
	return val
}

func ToMSX(any interface{}) M.SX {
	if any == nil {
		return M.SX{}
	}
	val, ok := any.(map[string]interface{})
	if L.CheckIf(!ok, `Can't convert to M.SX`, any) {
		return M.SX{}
	}
	return val
}

func ToMSS(any interface{}) M.SS {
	if any == nil {
		return M.SS{}
	}
	val, ok := any.(map[string]string)
	if L.CheckIf(!ok, `Can't convert to M.SS`, any) {
		return M.SS{}
	}
	return val
}

func ToYaml(any interface{}) string {
	bytes, err := yaml.Marshal(any)
	L.IsError(err, `yaml.Marshal`, any)
	return string(bytes)
}
