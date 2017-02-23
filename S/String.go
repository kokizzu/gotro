package S

// String support package
import (
	"encoding/json"
	"github.com/kokizzu/gotro/L"
	"strconv"
)

// simplified ternary operator (bool ? val : 0), returns second argument, if the condition (first arg) is true, returns empty string if not
//  S.If(true,`a`) // `a`
//  S.If(false,`a`) // ``
func If(b bool, yes string) string {
	if b {
		return yes
	}
	return ``
}

// ternary operator (bool ? val1 : val2), returns second argument if the condition (first arg) is true, third argument if not
//  I.IfElse(true,`a`,`b`) // `a`
//  I.IfElse(false,`a`,`b`) // `b`
func IfElse(b bool, yes, no string) string {
	if b {
		return yes
	}
	return no
}

// convert string to int64, returns 0 and silently print error if not valid
//  S.ToI(`1234`) // 1234
//  S.ToI(`1a`) // 0
func ToI(str string) int64 {
	val, err := strconv.ParseInt(str, 10, 64)
	L.IsError(err, str)
	return val
}

// convert to int64 with check
//  S.AsI(`1234`) // 1234, true
//  S.AsI(`1abc`) // 0, false
func AsI(str string) (int64, bool) {
	res, err := strconv.ParseInt(str, 10, 64)
	return res, err == nil
}

// convert string to float64, returns 0 and silently print error if not valid
//  S.ToF(`1234.5`) // 1234.5
//  S.ToF(`1a`) // 0.0
func ToF(str string) float64 {
	val, err := strconv.ParseFloat(str, 64)
	L.IsError(err, str)
	return val
}

// convert to float64 with check
//  S.AsF(`1234.5`) // 1234.5, true
//  S.AsF(`1abc`) // 0.0, false
func AsF(str string) (float64, bool) {
	res, err := strconv.ParseFloat(str, 64)
	return res, err == nil
}

// convert JSON object to map[string]interface{}, silently print and return empty map if failed
//  json_str := `{"test":123,"bla":[1,2,3,4]}`
//  map1 := S.JsonToMap(json_str)
func JsonToMap(str string) (res map[string]interface{}) {
	res = map[string]interface{}{}
	err := json.Unmarshal([]byte(str), &res)
	L.IsError(err, str)
	return
}

// convert JSON object to []interface{}, silently print and return empty slice of interface if failed
//  json_str := `[1,2,['test'],'a']`
//  arr := S.JsonToArr(json_str)
func JsonToArr(str string) (res []interface{}) {
	res = []interface{}{}
	err := json.Unmarshal([]byte(str), &res)
	L.IsError(err, str)
	return
}

// convert JSON object to map[string]interface{} with check
//  json_str := `{"test":123,"bla":[1,2,3,4]}`
//  map1, ok := S.JsonAsMap(json_str)
func JsonAsMap(str string) (res map[string]interface{}, ok bool) {
	res = map[string]interface{}{}
	err := json.Unmarshal([]byte(str), &res)
	ok = err != nil
	return
}

// convert JSON object to []interface{} with check
//  json_str := `[1,2,['test'],'a']`
//  arr, ok := S.JsonAsArr(json_str)
func JsonAsArr(str string) (res []interface{}, ok bool) {
	res = []interface{}{}
	err := json.Unmarshal([]byte(str), &res)
	ok = err != nil
	return
}
