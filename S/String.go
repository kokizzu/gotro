package S

// String support package
import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/kokizzu/gotro/B"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"math/rand"
)

// check whether the input string (first arg) starts with a certain character (second arg) or not.
//  S.StartsWith(`adakah`,`ad`) // bool(true)
//  S.StartsWith(`adakah`,`bad`) // bool(false)
func StartsWith(str, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}

// check whether the input string (first arg) ends with a certain character (second arg) or not.
//  S.EndsWith(`adakah`,`ah`)) // bool(true)
//  S.EndsWith(`adakah`,`aka`)) // bool(false)
func EndsWith(str, suffix string) bool {
	return strings.HasSuffix(str, suffix)
}

// check whether the input string (first arg) contains a certain sub string (second arg) or not.
//  S.Contains(`komputer`,`om`)) // bool(true)
//  S.Contains(`komputer`,`opu`)) // bool(false)
func Contains(str, substr string) bool {
	return strings.Contains(str, substr)
}

// compare two input string (first arg) equal with another input string (second arg).
//  S.Equals(`komputer`,`komputer`)) // bool(true)
//  S.Equals(`komputer`,`Komputer`)) // bool(false)
func Equals(strFirst, strSecond string) bool {
	return strFirst == strSecond
}

// count how many specific character (first arg) that the string (second arg) contains
//  S.Count(`komputeer`,`e`))// output int(2)
func Count(str, substr string) int {
	return strings.Count(str, substr)
}

// erase spaces from left and right
//  S.Trim(` withtrim:  `) // `withtrim:`
func Trim(str string) string {
	return strings.TrimSpace(str)
}

// replace all substring with another substring
//  S.Replace(`bisa`,`is`,`us`) // `busa`
func Replace(haystack, needle, gold string) string {
	return strings.Replace(haystack, needle, gold, -1)
}

// change the characters in string to lowercase
//  S.ToLower(`BIsa`) // "bisa"
func ToLower(str string) string {
	return strings.ToLower(str)
}

// change the characters in string to uppercase
// S.ToUpper(`bisa`) // "BISA"
func ToUpper(str string) string {
	return strings.ToUpper(str)
}

// get character at specific index, utf-8 safe
//  S.CharAt(`Halo 世界`, 5) // `世` // utf-8 example, if characters not shown, it's probably your font/editor/plugin
//  S.CharAt(`Halo`, 3) // `o`
func CharAt(str string, index int) string {
	for in, ch := range str {
		if in == index {
			return string(ch)
		}
	}
	return ``
}

// remove character at specific index, utf-8 safe
//  S.RemoveCharAt(`Halo 世界`, 5) // `Halo 界` --> utf-8 example, if characters not shown, it's probably your font/editor/plugin
//  S.RemoveCharAt(`Halo`, 3) // `Hal`
func RemoveCharAt(str string, index int) string {
	var chars []byte
	for in, ch := range str {
		if index != in {
			chars = append(chars, string(ch)...)
		}
	}
	return string(chars)
}

// Change first letter for every word to uppercase
//  S.ToTitle(`Disa dasi`)) // output "Disa Dasi"
func ToTitle(str string) string {
	return strings.Title(str)
}

// add single quote in the beginning and the end of string.
//  S.Q(`coba`) // `'coba'`
//  S.Q(`123`)  // `'123'`
func Q(str string) string {
	return `'` + str + `'`
}

// replace ` and give double quote (for table names)
//  S.ZZ(`coba"`) // `"coba&quot;"`
func ZZ(str string) string {
	str = Trim(str)
	str = Replace(str, `"`, `&quot;`)
	return `"` + str + `"`
}

// give ' to boolean value
//  S.ZB(true)  // `'true'`
//  S.ZB(false) // `'false'`
// boolean
func ZB(b bool) string {
	return `'` + B.ToS(b) + `'`
}

// give ' to int64 value
//  S.ZI(23)) // '23'
//  S.ZI(03)) // '3'
func ZI(num int64) string {
	return `'` + I.ToS(num) + `'`
}

// double quote a json string
//  hai := `{'test':123,"bla":[1,2,3,4]}`
//  L.Print(S.ZJ(hai))// output "{'test':123,\"bla\":[1,2,3,4]}"
func ZJ(str string) string {
	str = Replace(str, "\r", `\r`)
	str = Replace(str, "\n", `\n`)
	str = Replace(str, `"`, `\"`)
	return `"` + str + `"`
}

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

const i2c_cb63 = `-0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz`

var c2i_cb63 map[rune]int64

func init() {
	c2i_cb63 = map[rune]int64{}
	for i, ch := range i2c_cb63 {
		c2i_cb63[ch] = int64(i)
	}
}

// convert integer to custom base-63 encoding that lexicographically correct, positive integer only
//  0       -
//  1..10   0..9
//  11..36  A..Z
//  37      _
//  38..63  a..z
//  S.EncodeCB63(11) // `A`
//  S.EncodeCB63(1)) // `1`
func EncodeCB63(id int64, min_len int) string {
	if min_len < 1 {
		min_len = 1
	}
	str := make([]byte, 0, 12)
	for id > 0 {
		mod := rune(id % 64)
		str = append(str, i2c_cb63[mod])
		id /= 64
	}
	for len(str) < min_len {
		str = append(str, i2c_cb63[0])
	}
	l := len(str)
	for i, j := 0, l-1; i < l/2; i, j = i+1, j-1 {
		str[i], str[j] = str[j], str[i]
	}
	return string(str)
}

// convert custom base-63 encoding to int64
func DecodeCB63(str string) (int64, bool) {
	res := int64(0)
	for _, ch := range str {
		res *= 64
		val, ok := c2i_cb63[ch]
		if L.CheckIf(!ok, `Invalid character for CB63: `+string(ch)) {
			return -1, false
		}
		res += val
	}
	return res, true
}

// random CB63
func RandomCB63(len int64) string {
	res := ``
	for z := int64(0); z < len; z++ {
		res += EncodeCB63(rand.Int63(), 11)
	}
	return res
}
