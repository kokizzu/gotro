package M

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
)

// Map support package

// SX map with string key and any value
type SX map[string]any

// IX map with int64 key and any value
type IX map[int64]any

// SAX map with string key and array of any value
type SAX map[string][]any

// IAX map with int64 key and array of any value
type IAX map[int64][]any

// SS map with string key and string value
type SS map[string]string

// II map with int64 key and int64 value
type II map[int64]int64

// SF map with string key and float64 value
type SF map[string]float64

// SI map with string key and int64 value
type SI map[string]int64

// IS map with int64 key and string value
type IS map[int64]string

// IB map with int64 key and bool value
type IB map[int64]bool

// SB map with string key and bool value
type SB map[string]bool

// KeysConcat get concatenated integer keys
//
//	m := M.II{1: 2, 2: 567, 3:6, 5:45}
//	m.KeysConcat(`,`) // `1,2,3,5`
func (hash II) KeysConcat(with string) string {
	res := bytes.Buffer{}
	first := true
	for k := range hash {
		if first {
			first = false
		} else {
			res.WriteString(with)
		}
		res.WriteString(I.ToS(k))
	}
	return res.String()
}

// KeysConcat get concatenated string keys
//
//	m := M.SS{`tes`:`tes`,`coba`:`saja`,`lah`:`lah`}
//	m.KeysConcat(`,`) // `coba,lah,tes`
func (hash SS) KeysConcat(with string) string {
	res := bytes.Buffer{}
	first := true
	for k := range hash {
		if first {
			first = false
		} else {
			res.WriteString(with)
		}
		res.WriteString(k)
	}
	return res.String()
}

// ToScylla convert to scylla based map<text,text>
func (hash SS) ToScylla() string {
	res := bytes.Buffer{}
	res.WriteString(`{`)
	for k, v := range hash {
		res.WriteString(`'`)
		res.WriteString(S.Replace(k, `'`, `&apos;`))
		res.WriteString(`':'`)
		res.WriteString(S.Replace(v, `'`, `&apos;`))
		res.WriteString(`',`)
	}
	res.WriteString(`'':''}`)
	return res.String()

}

// ToJson convert to json string, silently print error if failed
func (hash SS) ToJson() string {
	str, err := json.Marshal(hash)
	L.IsError(err, `M.ToJson failed`, hash)
	return string(str)
}

// ToMsgp convert to msgpack string, silently print error if failed
func (hash SS) ToMsgp() []byte {
	str, err := msgpack.Marshal(hash)
	L.IsError(err, `M.ToMsgp failed`, hash)
	return str
}

// SortedKeys get sorted keys
//
//	m := M.SS{`tes`:`tes`,`coba`:`saja`,`lah`:`lah`}
//	m.SortedKeys() // []string{`coba`,`lah`,`tes`}
func (hash SS) SortedKeys() []string {
	res := make([]string, len(hash))
	idx := 0
	for k := range hash {
		res[idx] = k
		idx++
	}
	sort.Strings(res)
	return res
}

// Pretty get pretty printed values
func (hash SS) Pretty(sep string) string {
	keys := hash.SortedKeys()
	buff := bytes.Buffer{}
	for idx, key := range keys {
		buff.WriteString(key)
		buff.WriteRune(' ')
		buff.WriteString(hash[key])
		if idx > 0 {
			buff.WriteString(sep)
		}
	}
	return buff.String()
}

// PrettyFunc get pretty printed values with filter values
func (hash SS) PrettyFunc(sep string, fun func(string, string) string) string {
	keys := hash.SortedKeys()
	buff := bytes.Buffer{}
	for _, key := range keys {
		buff.WriteString(key)
		buff.WriteRune(' ')
		buff.WriteString(fun(key, hash[key]))
		if buff.Len() > 0 {
			buff.WriteString(sep)
		}
	}
	return buff.String()
}

// KeysConcat get concatenated string keys
//
//	m := M.SB{`tes`:true,`coba`:true,`lah`:true}
//	m.KeysConcat(`,`) // `coba,lah,tes`
func (hash SB) KeysConcat(with string) string {
	res := bytes.Buffer{}
	first := true
	for k := range hash {
		if first {
			first = false
		} else {
			res.WriteString(with)
		}
		res.WriteString(k)
	}
	return res.String()
}

// SortedKeys get sorted keys
//
//	m := M.SS{`tes`:true,`coba`:false,`lah`:false}
//	m.SortedKeys() // []string{`coba`,`lah`,`tes`}
func (hash SB) SortedKeys() []string {
	res := make([]string, len(hash))
	idx := 0
	for k := range hash {
		res[idx] = k
		idx++
	}
	sort.Strings(res)
	return res
}

// ToJson convert to json string, silently print error if failed
func (hash SB) ToJson() string {
	str, err := json.Marshal(hash)
	L.IsError(err, `M.ToJson failed`, hash)
	return string(str)
}

// ToMsgp convert to msgpack string, silently print error if failed
func (hash SB) ToMsgp() []byte {
	str, err := msgpack.Marshal(hash)
	L.IsError(err, `M.ToMsgp failed`, hash)
	return str
}

// ToJsonPretty convert to pretty json string, silently print error if failed
func (hash SB) ToJsonPretty() string {
	str, err := json.MarshalIndent(hash, ``, `  `)
	L.IsError(err, `M.ToJsonPretty failed`, hash)
	return string(str)
}

// IntoJson convert to json string with check
func (hash SB) IntoJson() (string, bool) {
	str, err := json.Marshal(hash)
	return string(str), err == nil
}

// IntoMsgp convert to msgpack string with check
func (hash SB) IntoMsgp() ([]byte, bool) {
	str, err := msgpack.Marshal(hash)
	return str, err == nil
}

// convert to pretty json string with check
func (hash SB) IntoJsonPretty() (string, bool) {
	str, err := json.MarshalIndent(hash, ``, `  `)
	return string(str), err == nil
}

// SortedKeys get sorted keys
//
//	m := M.SX{`tes`:1,`coba`:12.4,`lah`:false}
//	m.SortedKeys() // []string{`coba`,`lah`,`tes`}
func (hash SX) SortedKeys() []string {
	res := make([]string, len(hash))
	idx := 0
	for k := range hash {
		res[idx] = k
		idx++
	}
	sort.Strings(res)
	return res
}

// KeysConcat get concatenated integer keys
//
//	m := M.IB{1: true, 2: false, 3:true, 5:false}
//	m.KeysConcat(`,`) // `1,2,3,5`
func (hash IB) KeysConcat(with string) string {
	res := bytes.Buffer{}
	first := true
	for k := range hash {
		if first {
			first = false
		} else {
			res.WriteString(with)
		}
		res.WriteString(I.ToS(k))
	}
	return res.String()
}

// FromJson convert string to map, silently print error if failed
// doesn't overwrite existing value if json doesn't have the same key
// returns true if ok
func (hash *SX) FromJson(str string) (ok bool) {
	if hash == nil {
		return false
	}
	if len(str) < 3 {
		return false
	}
	err := json.Unmarshal([]byte(str), &hash)
	return !L.IsError(err, `M.SX.FromJson failed`, str)
}

// ToJson convert to json string, silently print error if failed
func (hash SX) ToJson() string {
	str, err := json.Marshal(hash)
	L.IsError(err, `M.ToJson failed`, hash)
	return string(str)
}

// ToMsgp convert to msgpack string, silently print error if failed
func (hash SX) ToMsgp() []byte {
	str, err := msgpack.Marshal(hash)
	L.IsError(err, `M.ToMsgp failed`, hash)
	return str
}

// ToJsonPretty convert to pretty json string, silently print error if failed
func (hash SX) ToJsonPretty() string {
	str, err := json.MarshalIndent(hash, ``, `  `)
	L.IsError(err, `M.ToJsonPretty failed`, hash)
	return string(str)
}

// IntoJson convert to json string with check
func (hash SX) IntoJson() (string, bool) {
	str, err := json.Marshal(hash)
	return string(str), err == nil
}

// IntoMsgp convert to msgpack string with check
func (hash SX) IntoMsgp() ([]byte, bool) {
	str, err := msgpack.Marshal(hash)
	return str, err == nil
}

// IntoJsonPretty convert to pretty json string with check
func (hash SX) IntoJsonPretty() (string, bool) {
	str, err := json.MarshalIndent(hash, ``, `  `)
	return string(str), err == nil
}

// GetInt get int64 type from map
//
//	m := M.SX{`test`:234.345,`coba`:`buah`,`dia`:true,`angka`:int64(23435)}
//	m.GetInt(`test`))  // int64(234)
//	m.GetInt(`dia`))   // int64(1)
//	m.GetInt(`coba`))  // int64(0)
//	m.GetInt(`angka`)) // int64(23435)
func (json SX) GetInt(key string) int64 {
	any := json[key]
	if any == nil {
		return 0
	}
	if val, ok := any.(int64); ok {
		return val
	}
	switch v := any.(type) {
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case uint:
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
	case bool:
		if v {
			return 1
		}
		return 0
	case string:
		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			return val
		}
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			return int64(val)
		}
		L.Describe(`Property [` + key + `] is not an int64: ` + fmt.Sprintf("%T", any))
	default:
		L.Describe(`Property [` + key + `] is not an int64: ` + fmt.Sprintf("%T", any))
	}
	return 0
}

// GetUint get uint type from map
//
//	m := M.SX{`test`:234.345,`coba`:`buah`,`dia`:true,`angka`:int64(23435)}
//	m.GetInt(`test`))  // int64(234)
//	m.GetInt(`dia`))   // int64(1)
//	m.GetInt(`coba`))  // int64(0)
//	m.GetInt(`angka`)) // int64(23435)
func (json SX) GetUint(key string) uint64 {
	any := json[key]
	if any == nil {
		return 0
	}
	if val, ok := any.(uint64); ok {
		return val
	}
	switch v := any.(type) {
	case int:
		return uint64(v)
	case int8:
		return uint64(v)
	case int16:
		return uint64(v)
	case int32:
		return uint64(v)
	case int64:
		return uint64(v)
	case uint:
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
	case bool:
		if v {
			return 1
		}
		return 0
	case string:
		if val, err := strconv.ParseUint(v, 10, 64); err == nil {
			return uint64(val)
		}
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			return uint64(val)
		}
		L.Describe(`Property [` + key + `] is not an uint64: ` + fmt.Sprintf("%T", any))
	default:
		L.Describe(`Property [` + key + `] is not an uint64: ` + fmt.Sprintf("%T", any))
	}
	return 0
}

// GetFloat get float64 type from map
//
//	m := M.SX{`test`:234.345,`coba`:`buah`,`dia`:true,`angka`:23435}
//	m.GetFloat(`test`)  // float64(234.345)
//	m.GetFloat(`dia`)   // int64(1)
//	m.GetFloat(`coba`)  // 0
//	m.GetFloat(`angka`) // 0
func (json SX) GetFloat(key string) float64 {
	any := json[key]
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
	case time.Duration:
		return float64(v)
	case bool:
		if v {
			return 1
		}
		return 0
	case string:
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			return val
		}
		L.Describe(`Property [`+key+`] is not a float64: `+fmt.Sprintf("%T", any), any)
	default:
		L.Describe(`Property [`+key+`] is not a float64: `+fmt.Sprintf("%T", any), any)
	}
	return 0
}

// GetStr get string type from map
//
//	m := M.SX{`test`:234.345,`coba`:`buah`,`angka`:int64(123)}
//	m.GetStr(`test`)  // `234.345`
//	m.GetStr(`coba`)  // `buah`
//	m.GetStr(`angka`) // `123`
func (json SX) GetStr(key string) string {
	any := json[key]
	if any == nil {
		return ``
	}
	if val, ok := any.(string); ok {
		return val
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
		return strconv.FormatInt(int64(v), 10)
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
	case bool:
		if v {
			return `true`
		}
		return `false`
	case fmt.Stringer:
		return v.String()
	case *string:
		if v == nil {
			return ``
		}
		return *v
	default:
		L.Describe(`Property [` + key + `] is not a string: ` + fmt.Sprintf("%T", any))
	}
	return ``
}

// GetBool get bool type from map (empty string, 0, `f`, `false` are false, other non empty are considered truthy value)
// m := M.SX{`test`:234.345,`coba`:`buah`,`angka`:float64(123),`salah`:123}
//
//	m.GetBool(`test`)  // bool(true)
//	m.GetBool(`coba`)  // bool(true)
//	m.GetBool(`angka`) // bool(true)
//	m.GetBool(`salah`) // bool(false)
func (json SX) GetBool(key string) bool {
	any := json[key]
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
	case fmt.Stringer:
		val := v.String()
		val = strings.TrimSpace(strings.ToLower(val))
		return !(val == `` || val == `0` || val == `f` || val == `false`)
	case string:
		val := v
		val = strings.TrimSpace(strings.ToLower(val))
		return !(val == `` || val == `0` || val == `f` || val == `false`)
	default:
		L.Describe(`Property [` + key + `] is not a bool: ` + fmt.Sprintf("%T", any))
	}
	return false
}

// GetMSB get map string bool value from map
//
//	m := M.SX{`tes`:M.SB{`1`:true,`2`:false}}
//	m.GetMSB(`tes`) // M.SB{"1":true, "2":false}
func (json SX) GetMSB(key string) SB {
	v := json[key]
	if v == nil {
		return SB{}
	}
	if val, ok := v.(map[string]bool); ok {
		return val
	} else if val, ok := v.(SB); ok {
		return val
	} else if val, ok := v.(map[string]any); ok {
		res := SB{}
		for k, vx := range val {
			if vb, ok := vx.(bool); ok {
				res[k] = vb
			}
		}
		return res
	} else {
		L.Describe(`Property [` + key + `] is not a M.SB: ` + fmt.Sprintf("%T", v))
		return SB{}
	}
}

// GetMSF get map string float64 value from map
//
//	m := M.SX{`tes`:M.SF{`satu`:32.45,`2`:12}}
//	m.GetMSF(`tes`) // M.SF{"satu":32.45, "2":12}
func (json SX) GetMSF(key string) SF {
	v := json[key]
	if v == nil {
		return SF{}
	}
	if val, ok := v.(map[string]float64); ok {
		return val
	} else if val, ok := v.(SF); ok {
		return val
	} else if val, ok := v.(map[string]any); ok {
		res := SF{}
		for k, vx := range val {
			if vf, ok := vx.(float64); ok {
				res[k] = vf
			} else if vs, ok := vx.(string); ok {
				res[k] = S.ToF(vs)
			}
		}
		return res
	} else {
		L.Describe(`Property [` + key + `] is not a M.SF: ` + fmt.Sprintf("%T", v))
		return SF{}
	}
}

// GetMSI get map string int64 value from map
//
//	m := M.SX{`tes`:M.SF{`satu`:32,`2`:12}}
//	m.GetMSI(`tes`) // M.SF{"satu":32, "2":12}
func (json SX) GetMSI(key string) SI {
	v := json[key]
	if v == nil {
		return SI{}
	}
	if val, ok := v.(map[string]int64); ok {
		return val
	} else if val, ok := v.(SI); ok {
		return val
	} else if val, ok := v.(map[string]any); ok {
		res := SI{}
		for k, vx := range val {
			if vi, ok := vx.(int64); ok {
				res[k] = vi
			} else if vs, ok := vx.(string); ok {
				res[k] = S.ToI(vs)
			}
		}
		return res
	} else {
		L.Describe(`Property [` + key + `] is not a M.SF: ` + fmt.Sprintf("%T", v))
		return SI{}
	}
}

// GetMIB get map string int64 value from map
//
//	m := M.SX{`tes`:M.SB{`satu`:true,`2`:false}}
//	m.GetMSB(`tes`) // M.SB{"satu":true, "2":false}
func (json SX) GetMIB(key string) IB {
	v := json[key]
	if v == nil {
		return IB{}
	}
	if val, ok := v.(map[int64]bool); ok {
		return val
	} else if val, ok := v.(IB); ok {
		return val
	} else if val, ok := v.(map[int64]any); ok {
		res := IB{}
		for k, vx := range val {
			if vb, ok := vx.(bool); ok {
				res[k] = vb
			}
		}
		return res
	} else {
		L.Describe(`Property [` + key + `] is not a M.SB: ` + fmt.Sprintf("%T", v))
		return IB{}
	}
}

// GetMSX get map string anything value from map
//
//	m :=  M.SX{`tes`:M.SX{`satu`:234.345,`dua`:`huruf`,`tiga`:123}}
//	m.GetMSX(`tes`) // M.SX{"tiga": int(123),"satu": float64(234.345),"dua":  "huruf"}
func (json SX) GetMSX(key string) SX {
	v := json[key]
	if v == nil {
		return SX{}
	}
	if val, ok := v.(map[string]any); ok {
		return val
	} else if val, ok := v.(SX); ok {
		return val
	} else {
		L.Describe(`Property [` + key + `] is not a M.SX: ` + fmt.Sprintf("%T", v))
		return SX{}
	}
}

// GetAX get array of anything value from map
//
//	m :=  M.SX{`tes`:[]any{123,`buah`}}
//	m.GetAX(`tes`) // []interface {}{int(123),"buah"}
func (json SX) GetAX(key string) []any {
	v := json[key]
	if v == nil {
		return []any{}
	}
	if val, ok := v.([]any); ok {
		return val
	} else {
		L.Describe(`Property [` + key + `] is not a A.X: ` + fmt.Sprintf("%T", v))
		return []any{}
	}
}

// GetIntArr get array int64 value from map
//
//	m :=  M.SX{`tes`:[]int64{123,234}}
//	m.GetIntArr(`tes`)) // []int64{123, 234}
func (json SX) GetIntArr(key string) []int64 {
	v := json[key]
	if v == nil {
		return []int64{}
	}
	if val, ok := v.([]int64); ok {
		return val
	} else if val, ok := v.([]float64); ok {
		res := []int64{}
		for _, vf := range val {
			res = append(res, int64(vf))
		}
		return res
	} else if val, ok := v.([]any); ok {
		res := []int64{}
		for k, vx := range val {
			switch v := vx.(type) {
			case int:
				res = append(res, int64(v))
			case int8:
				res = append(res, int64(v))
			case int16:
				res = append(res, int64(v))
			case int32:
				res = append(res, int64(v))
			case uint:
				res = append(res, int64(v))
			case uint8:
				res = append(res, int64(v))
			case uint16:
				res = append(res, int64(v))
			case uint32:
				res = append(res, int64(v))
			case uint64:
				res = append(res, int64(v))
			case float32:
				res = append(res, int64(v))
			case float64:
				res = append(res, int64(v))
			case string:
				if val, err := strconv.ParseInt(v, 10, 64); err == nil {
					res = append(res, int64(val))
				}
				if val, err := strconv.ParseFloat(v, 64); err == nil {
					res = append(res, int64(val))
				}
				L.Describe(`Property [` + key + `][` + I.ToStr(k) + `] is not an int64: ` + fmt.Sprintf("%T", v))
			}
		}
		return res
	} else {
		L.Describe(`Property [` + key + `] is not a []int64: ` + fmt.Sprintf("%T", v))
		return []int64{}
	}
}

// GetInt get int64 from from map
func (hash SS) GetInt(key string) int64 {
	return S.ToI(hash[key])
}

// GetUint get uint from map
func (hash SS) GetUint(key string) uint64 {
	return S.ToU(hash[key])
}

// GetFloat get float64 type from map
func (hash SS) GetFloat(key string) float64 {
	return S.ToF(hash[key])
}

// GetStr get string type from map
func (hash SS) GetStr(key string) string {
	return hash[key]
}

// Keys get array of string keys
//
//	m :=  M.SS{`satu`:`1`,`dua`:`2`}
//	m.Keys() // []string{"satu", "dua"}
func (hash SS) Keys() []string {
	res := []string{}
	for k := range hash {
		res = append(res, k)
	}
	return res
}

// Merge merge from another map
func (hash SS) Merge(src SS) {
	for k, v := range src {
		hash[k] = v
	}
}

// Keys get array of string keys
//
//	m :=  M.SS{`satu`:`1`,`dua`:`2`}
//	m.Keys() // []string{"satu", "dua"}
func (hash SX) Keys() []string {
	res := []string{}
	for k := range hash {
		res = append(res, k)
	}
	return res
}

// Keys get array of int64 keys
//
//	m :=  M.IX{1:1,2:`DIA`}
//	m.Keys()) // []int64{1, 2}
func (hash IX) Keys() []int64 {
	res := []int64{}
	for k := range hash {
		res = append(res, k)
	}
	return res
}

// Keys get array of int64 keys
//
//	m :=  M.II{1:1,2:3}
//	m.Keys() // []int64{1, 2}
func (hash II) Keys() []int64 {
	res := []int64{}
	for k := range hash {
		res = append(res, k)
	}
	return res
}

// Keys get array of int64 keys
//
//	m :=  M.IB{1:true,2:false}
//	m.Keys() // []int64{1, 2}
//
// get all integer keys
func (hash IB) Keys() []int64 {
	res := []int64{}
	for k := range hash {
		res = append(res, k)
	}
	return res
}

// ToSX convert keys to string
//
//	m :=  M.IX{1:1,2:`DUA`}
//	m.ToSX() // M.SX{"1": int(1),"2": "DUA"}
//
// convert integer keys to string keys
func (hash IX) ToSX() SX {
	res := SX{}
	for k, v := range hash {
		res[I.ToS(k)] = v
	}
	return res
}

// ToJson convert map[string]any to json
//
//	m :=  map[string]any{`buah`:123,`angka`:`dia`}
//	M.ToJson(m) // {"angka":"dia","buah":123}
func ToJson(hash map[string]any) string {
	str, err := json.Marshal(hash)
	L.IsError(err, `M.ToJson failed`, hash)
	return string(str)
}

// ToMsgp convert map[string]any to json
//
//	m :=  map[string]any{`buah`:123,`angka`:`dia`}
//	M.ToMsgp(m) // []byte{0x82, 0xa5, 0x61, 0x6e, 0x67, 0x6b, 0x61, 0xa3, 0x64, 0x69, 0x61, 0xa5, 0x62, 0x75, 0x61, 0x68, 0xcd, 0x7b}
func ToMsgp(hash map[string]any) []byte {
	str, err := msgpack.Marshal(hash)
	L.IsError(err, `M.ToMsgp failed`, hash)
	return str
}

// Set set key with any value
func (hash SX) Set(key string, val any) {
	hash[key] = val
}

// Pretty get pretty printed values
func (hash SX) Pretty(sep string) string {
	keys := hash.SortedKeys()
	buff := bytes.Buffer{}
	for _, key := range keys {
		buff.WriteString(key)
		buff.WriteRune(' ')
		buff.WriteString(fmt.Sprintf("%+v", hash[key]))
		buff.WriteString(sep)
	}
	return buff.String()
}

// SSKeysStartedWith retrieve all keys started with
func SSKeysStartedWith(m SS, prefix string) []string {
	res := []string{}
	for k := range m {
		if S.StartsWith(k, prefix) {
			res = append(res, k)
		}
	}
	return res
}
