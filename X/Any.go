package X

// Any type support package

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/goccy/go-yaml"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/B"
	"github.com/kokizzu/gotro/C"
	"github.com/kokizzu/gotro/F"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
)

// ToU convert any data type to uint
//
//	var m any
//	m = `123`
//	L.ParentDescribe(X.ToI(m)) // uint(123)
func ToU(x any) uint64 {
	if x == nil {
		return 0
	}
	if val, ok := x.(uint64); ok {
		return val
	}
	switch v := x.(type) {
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
	case time.Duration:
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
		L.ParentDescribe(`Can't convert to uint64`, x)
	case string:
		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			return uint64(val)
		}
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			return uint64(val)
		}
		L.ParentDescribe(`Can't convert to uint64`, x)
	case *any:
		if v != nil {
			return ToU(*v)
		}
	default:
		L.ParentDescribe(`Can't convert to uint64`, x)
	}
	return 0
}

// ToByte convert any data type to int8
//
//	var m any
//	m = `123`
//	L.ParentDescribe(X.ToByte(m)) // byte(123)
func ToByte(x any) byte {
	if x == nil {
		return 0
	}
	switch v := x.(type) {
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
	case time.Duration:
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
		L.ParentDescribe(`Can't convert to byte`, x)
	case string:
		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			return byte(val)
		}
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			return byte(val)
		}
		L.ParentDescribe(`Can't convert to byte`, x)
	case *any:
		if v != nil {
			return ToByte(*v)
		}
	default:
		L.ParentDescribe(`Can't convert to byte`, x)
	}
	return 0
}

// ToI convert any data type to int64
//
//	var m any
//	m = `123`
//	L.ParentDescribe(X.ToI(m)) // int64(123)
func ToI(x any) int64 {
	if x == nil {
		return 0
	}
	if val, ok := x.(int64); ok {
		return val
	}
	switch v := x.(type) {
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
	case time.Duration:
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
		L.ParentDescribe(`Can't convert to int64 from []byte`, x)
	case string:
		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			return val
		}
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			return int64(val)
		}
		L.ParentDescribe(`Can't convert to int64 from string`, x)
	case *any:
		if v != nil {
			return ToI(*v)
		}
	default:
		L.ParentDescribe(`Can't convert to int64`, x)
	}
	return 0
}

// ToF Convert any data type to float64
//
//	var m any
//	m = `123.5`
//	L.ParentDescribe(X.ToF(m)) // float64(123.5)
func ToF(x any) float64 {
	if x == nil {
		return 0
	}
	if val, ok := x.(float64); ok {
		return val
	}
	switch v := x.(type) {
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
	case time.Duration:
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
		L.ParentDescribe(`Can't convert to float64`, x)
	case string:
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			return val
		}
		L.ParentDescribe(`Can't convert to float64`, x)
	case *any:
		if v != nil {
			return ToF(*v)
		}
	default:
		L.ParentDescribe(`Can't convert to float64`, x)
	}
	return 0
}

// ToS convert any data type to string
//
//	var m any
//	m = `123`
//	L.ParentDescribe(X.ToS(m)) // `123`
func ToS(x any) string {
	if x == nil {
		return ``
	}
	if val, ok := x.(string); ok {
		return val
	}
	if val, ok := x.([]byte); ok {
		return string(val)
	}
	switch v := x.(type) {
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
	case *any:
		if v != nil {
			return ToS(*v)
		}
	default:
		return ToJson5(v)
	}
	return ``
}

// ToTime convert any to time
func ToTime(x any) time.Time {
	if x == nil {
		return time.Time{}
	}
	switch v := x.(type) {
	case time.Time:
		return v
	case *time.Time:
		if v != nil {
			return *v
		}
	case []byte: // "YYYY-MM-DD HH:MM:SS.MMMMMM"
		res, err := parseDateTime(v, time.UTC)
		L.IsError(err, `Can't convert to time.Time from []byte`, x)
		return res
	case string: // "YYYY-MM-DD HH:MM:SS.MMMMMM"
		res, err := parseDateTime([]byte(v), time.UTC)
		L.IsError(err, `Can't convert to time.Time from string`, x)
		return res
	case *[]byte: // "YYYY-MM-DD HH:MM:SS.MMMMMM"
		if v != nil {
			res, err := parseDateTime(*v, time.UTC)
			L.IsError(err, `Can't convert to time.Time from []byte`, x)
			return res
		}
	case *string: // "YYYY-MM-DD HH:MM:SS.MMMMMM"
		if v != nil {
			res, err := parseDateTime([]byte(*v), time.UTC)
			L.IsError(err, `Can't convert to time.Time from []byte`, x)
			return res
		}
	case *any:
		if v != nil {
			return ToTime(*v)
		}
	default:
		L.CheckIf(false, `Can't convert to time.Time`, x)
	}
	return time.Time{}
}

// ToBool convert any data type to bool
//
//	var m any
//	m = `123`
//	L.ParentDescribe(X.ToBool(m)) // bool(true)
func ToBool(any any) bool {
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

// ToArr convert any data type to array of any
//
//	var m3 any
//	m3 = []any{1}   // tipe array
//	L.ParentDescribe(X.ToArr(m3)) // []interface {}{int(1),}
func ToArr(x any) []any {
	if x == nil {
		return []any{}
	}
	val, ok := x.([]any)
	if L.CheckIf(!ok, `Can't convert to []any`, x) {
		return []any{}
	}
	return val
}

// ArrToStrArr convert array of any data type to array of string
//
//	var m4 []any
//	m4 = []any{1}     // // tipe array
//	L.ParentDescribe(X.ArrToStrArr(m4)) // []string{"1"}
func ArrToStrArr(any_arr []any) []string {
	res := []string{}
	for _, val := range any_arr {
		res = append(res, ToS(val))
	}
	return res
}

// ArrToIntArr Convert array of any data type to array of int64
//
//	var m4 []any
//	m4 = []any{1}     // // tipe array
//	L.ParentDescribe(X.ArrToIntArr(m4)) // []int64{1}
func ArrToIntArr(any_arr []any) []int64 {
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

func json5fromMIX(orig map[int64]any) string {
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

func json5fromMIAX(orig map[int64][]any) string {
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

func json5fromMSAX(orig map[string][]any) string {
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

// ToJson5 convert to json5
func ToJson5(x any) string {
	// bug when using map[int64]any
	if x == nil {
		return `''`
	}
	switch orig := x.(type) {
	case bytes.Buffer: // return as is
		return orig.String()
	case string:
		return S.ZJJ(orig)
	case []byte:
		return S.ZJJ(string(orig))
	case int:
		return I.ToStr(orig)
	case int64:
		return I.ToS(orig)
	case int32:
		return I.ToS(int64(orig))
	case uint:
		return I.UToStr(orig)
	case uint64:
		return I.UToS(orig)
	case uint32:
		return I.UToS(uint64(orig))
	case float32:
		return F.ToS(float64(orig))
	case float64:
		return F.ToS(orig)
	case bool:
		return B.ToS(orig)
	case M.IB:
		return json5fromMIB(orig)
	case map[int64]bool:
		return json5fromMIB(orig)
	case M.IX:
		return json5fromMIX(orig)
	case map[int64]any:
		return json5fromMIX(orig)
	case M.IAX:
		return json5fromMIAX(orig)
	case map[int64][]any:
		return json5fromMIAX(orig)
	case M.SAX:
		return json5fromMSAX(orig)
	case map[string][]any:
		return json5fromMSAX(orig)
	case M.SX:
		return orig.ToJson()
	case map[string]any:
		return M.ToJson(orig)
	//   return any.(M.SX).ToJson()
	case M.SI:
		return json5fromMSI(orig)
	case map[string]int64:
		return json5fromMSI(orig)
	case A.X:
		return A.ToJson(orig)
	case []any:
		return A.ToJson(orig)
	default:
		str, err := json.Marshal(x)
		L.IsError(err, `X.ToJson5 failed`, x)
		return string(str)
	}
	// TODO: add more types (M/A) here, do not EVER TRY to use reflection in this case
}

// ToJsonPretty convert to beautiful json text
//
//	m:= []interface {}{true,`1`,23,`wabcd`}
//	L.Print(K.ToJsonPretty(m))
//	// [
//	//   true,
//	//   "1",
//	//   23,
//	//   "wabcd"
//	// ]
func ToJsonPretty(any any) string {
	res, err := json.MarshalIndent(any, ``, `  `)
	L.IsError(err, `X.ToJsonPretty failed`, any)
	return string(res)
}

// ToJson convert to standard json text
func ToJson(any any) string {
	res, err := json.Marshal(any)
	L.IsError(err, `X.ToJson failed`, any)
	return string(res)
}

// ToAX  convert to []any
func ToAX(x any) A.X {
	if x == nil {
		return A.X{}
	}
	val, ok := x.([]any)
	if L.CheckIf(!ok, `Can't convert to A.X`, x) {
		return A.X{}
	}
	return val
}

// ToMSX convert to map[string]any
func ToMSX(x any) M.SX {
	if x == nil {
		return M.SX{}
	}
	val, ok := x.(map[string]any)
	if L.CheckIf(!ok, `Can't convert to M.SX`, x) {
		return M.SX{}
	}
	return val
}

// ToMSS convert to map[string]string
func ToMSS(any any) M.SS {
	if any == nil {
		return M.SS{}
	}
	val, ok := any.(map[string]string)
	if L.CheckIf(!ok, `Can't convert to M.SS`, any) {
		return M.SS{}
	}
	return val
}

// ToYaml convert to yaml text
func ToYaml(any any) string {
	bytes, err := yaml.Marshal(any)
	L.IsError(err, `yaml.Marshal`, any)
	return string(bytes)
}
