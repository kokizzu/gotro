package S

// String support package
import (
	"encoding/json"
	"strconv"
	"strings"

	"bytes"
	"github.com/kokizzu/gotro/C"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/yosuke-furukawa/json5/encoding/json5"
	"math/rand"
	"unicode"
)

const WebBR = "\n<br/>"

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

// compare two input string (first arg) equal with ignoring case another input string (second arg).
//  S.EqualsIgnoreCase(`komputer`,`komputer`)) // bool(true)
//  S.EqualsIgnoreCase(`komputer`,`Komputer`)) // bool(true)
func EqualsIgnoreCase(strFirst, strSecond string) bool {
	strFirst = ToLower(strFirst)
	strSecond = ToLower(strSecond)
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

// remove chars from beginning and end
//  S.TrimChars(`aoaaffoa`,`ao`) // `ff`
func TrimChars(str, chars string) string {
	return strings.Trim(str, chars)
}

// get first index of
// S.IndexOf(`abcdcd`,`c) // 2, -1 if not exists
func IndexOf(str, sub string) int {
	return strings.Index(str, sub)
}

// get last index of
//  S.LastIndexOf(`abcdcd`,`c`) // 4, -1 if not exists
func LastIndexOf(str, sub string) int {
	return strings.LastIndex(str, sub)
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
//  S.IfElse(true,`a`,`b`) // `a`
//  S.IfElse(false,`a`,`b`) // `b`
func IfElse(b bool, yes, no string) string {
	if b {
		return yes
	}
	return no
}

// coalesce, return first non-empty string
//  S.IfEmpty(``,`2`) // `2`
//  S.IfEmpty(`1`,`2`) // `1`
func IfEmpty(str1, str2 string) string {
	if str1 != `` {
		return str1
	}
	return str2
}

// coalesce, return first non-empty string
//  S.Coalesce(`1`,`2`) // `1`
//  S.Coalesce(``,`2`) // `2`
//  S.Coalesce(``,``,`3`) // `3`
func Coalesce(strs ...string) string {
	for _, str := range strs {
		if str != `` {
			return str
		}
	}
	return ``
}

// convert string to uint64, returns 0 and silently print error if not valid
//  S.ToU(`1234`) // 1234
//  S.ToU(`1a`) // 0
func ToU(str string) uint64 {
	val, _ := strconv.ParseUint(str, 10, 64)
	//L.IsError(err, str)
	return val
}

// convert string to int64, returns 0 and silently print error if not valid
//  S.ToI(`1234`) // 1234
//  S.ToI(`1a`) // 0
func ToI(str string) int64 {
	val, _ := strconv.ParseInt(str, 10, 64)
	//L.IsError(err, str)
	return val
}

// convert string to int, returns 0 and silently print error if not valid
//  S.ToInt(`1234`) // 1234
//  S.ToInt(`1a`) // 0
func ToInt(str string) int {
	val, _ := strconv.ParseInt(str, 10, 64)
	//L.IsError(err, str)
	return int(val)
}

// convert to uint with check
//  S.AsU(`1234`) // 1234, true
//  S.AsU(`1abc`) // 0, false
func AsU(str string) (uint, bool) {
	res, err := strconv.ParseInt(str, 10, 64)
	return uint(res), err == nil
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
	val, _ := strconv.ParseFloat(str, 64)
	//L.IsError(err, str)
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
	if len(str) == 0 {
		return
	}
	err := json.Unmarshal([]byte(str), &res)
	L.IsError(err, str)
	return
}

// convert JSON object to map[string]string, silently print and return empty map if failed
//  json_str := `{"test":123,"bla":[1,2,3,4]}`
//  map1 := S.JsonToMap(json_str)
func JsonToStrStrMap(str string) (res map[string]string) {
	res = map[string]string{}
	if len(str) == 0 {
		return
	}
	err := json.Unmarshal([]byte(str), &res)
	L.IsError(err, str)
	return
}

// convert JSON object to []interface{}, silently print and return empty slice of interface if failed
//  json_str := `[1,2,['test'],'a']`
//  arr := S.JsonToArr(json_str)
func JsonToArr(str string) (res []interface{}) {
	res = []interface{}{}
	if len(str) == 0 {
		return
	}
	err := json.Unmarshal([]byte(str), &res)
	L.IsError(err, str)
	return
}

// convert JSON object to []map[string]interface{}, silently print and return empty slice of interface if failed
//  json_str := `[{"x":"foo"},{"y":"bar"}]`
//  arr := S.JsonToObjArr(json_str)
func JsonToObjArr(str string) (res []map[string]interface{}) {
	res = []map[string]interface{}{}
	if len(str) == 0 {
		return
	}
	err := json.Unmarshal([]byte(str), &res)
	L.IsError(err, str)
	return
}

// convert JSON object to []string, silently print and return empty slice of interface if failed
//  json_str := `["123","456",789]`
//  arr := S.JsonToStrArr(json_str)
func JsonToStrArr(str string) (res []string) {
	res = []string{}
	if len(str) == 0 {
		return
	}
	err := json5.Unmarshal([]byte(str), &res)
	L.IsError(err, str)
	return
}

// convert JSON object to []int64, silently print and return empty slice of interface if failed
//  json_str := `[1,2,['test'],'a']`
//  arr := S.JsonToArr(json_str)
func JsonToIntArr(str string) (res []int64) {
	res = []int64{}
	if len(str) == 0 {
		return
	}
	err := json5.Unmarshal([]byte(str), &res)
	L.IsError(err, str)
	return
}

// repeat string
func Repeat(str string, count int) string {
	return strings.Repeat(str, count)
}

// convert JSON object to map[string]interface{} with check
//  json_str := `{"test":123,"bla":[1,2,3,4]}`
//  map1, ok := S.JsonAsMap(json_str)
func JsonAsMap(str string) (res map[string]interface{}, ok bool) {
	res = map[string]interface{}{}
	err := json.Unmarshal([]byte(str), &res)
	ok = err == nil
	return
}

// convert JSON object to []interface{} with check
//  json_str := `[1,2,['test'],'a']`
//  arr, ok := S.JsonAsArr(json_str)
func JsonAsArr(str string) (res []interface{}, ok bool) {
	res = []interface{}{}
	err := json.Unmarshal([]byte(str), &res)
	ok = err == nil
	return
}

// convert JSON object to []string with check
//  json_str := `["a","b","c"]`
//  arr, ok := S.JsonAsStrArr(json_str)
func JsonAsStrArr(str string) (res []string, ok bool) {
	res = []string{}
	err := json.Unmarshal([]byte(str), &res)
	ok = err == nil
	return
}

// convert JSON object to []int64 with check
//  json_str := `[1,2,3]`
//  arr, ok := S.JsonAsIntArr(json_str)
func JsonAsIntArr(str string) (res []int64, ok bool) {
	res = []int64{}
	err := json.Unmarshal([]byte(str), &res)
	ok = err == nil
	return
}

// convert JSON object to []float64 with check
//  json_str := `[1,2,3]`
//  arr, ok := S.JsonAsFloatArr(json_str)
func JsonAsFloatArr(str string) (res []float64, ok bool) {
	res = []float64{}
	err := json.Unmarshal([]byte(str), &res)
	ok = err == nil
	return
}

// split a string (first arg) by characters (second arg) into array of strings (output).
//  S.Split(`biiiissssa`,func(ch rune) bool { return ch == `i` }) // output []string{"b", "", "", "", "ssssa"}
func Split(str, sep string) []string {
	return strings.Split(str, sep)
}

// split a string (first arg) based on a function
func SplitFunc(str string, fun func(rune) bool) []string {
	return strings.FieldsFunc(str, fun)
}

// append padStr to left until length is lenStr
func PadLeft(s string, padStr string, lenStr int) string {
	var padCount int
	padCount = I.MaxOf(lenStr-len(s), 0)
	return strings.Repeat(padStr, padCount) + s
}

// append padStr to right until length is lenStr
func PadRight(s string, padStr string, lenStr int) string {
	var padCount int
	padCount = I.MaxOf(lenStr-len(s), 0)
	return s + strings.Repeat(padStr, padCount)
}

// return valid version of mail contact (part before <usr@email>)
func ValidateMailContact(str string) string {
	str = Replace(str, `,`, `_`)
	str = Replace(str, `.`, `_`)
	str = Replace(str, `<`, `_`)
	str = Replace(str, `>`, `_`)
	str = Replace(str, `(`, `_`)
	str = Replace(str, `)`, `_`)
	str = Replace(str, `@`, `_`)
	return str
}

// return formatted array of mail contact <usr@email>
func MergeMailContactEmails(each_name, str_emails string) []string {
	temp := []string{}
	str_email := Split(str_emails, `,`)
	for _, each_email := range str_email {
		each_email = Trim(each_email)
		if each_email == `` {
			continue
		}
		each_name = ValidateMailContact(each_name)
		temp = append(temp, each_name+`<`+each_email+`>`)
	}
	return temp
}

// return empty string if str is not a valid email
func ValidateEmail(str string) string {
	res := strings.Split(str, `@`)
	if len(res) != 2 {
		return ``
	}
	if (strings.Trim(res[0], `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&'*+-/=?^_{|}~.`)) == `` {
		if strings.Trim(res[1], `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-.`) == `` {
			return str
		}
	}
	return ``
}

// remove invalid characters of a phone number
func ValidatePhone(str string) string {
	res := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) || r == '+' || r == ' ' || r == '-' {
			return r
		}
		return -1
	}, str)
	return res
}

// validate file name
func ValidateFilename(str string) string {
	res := strings.Map(func(r rune) rune {
		if C.IsValidFilename(byte(r)) {
			return r
		}
		return -1
	}, str)
	return res
}

// create a random password
func RandomPassword(strlen int64) string {
	const chars = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ123456789" // l and I removed
	result := make([]byte, strlen)
	for i := int64(0); i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// split to substrings with maximum n characters
func SplitN(str string, n int) []string {
	if len(str) < n {
		return []string{str}
	}
	sub := ``
	subs := []string{}
	runes := bytes.Runes([]byte(str))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ``
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}
	return subs
}

// substring before first `substr`
func LeftOf(str, substr string) string {
	len := strings.Index(str, substr)
	if len < 0 {
		return str
	}
	return str[:len]
}

// substring after first `substr`
func RightOf(str, substr string) string {
	pos := strings.Index(str, substr)
	if pos < 0 {
		return str
	}
	return str[pos+len(substr):]
}

// substring at most n characters
func LeftN(str string, n int) string {
	if len(str) < n {
		return str
	}
	return str[:n] + `...`
}

// substring at most n characters
func Left(str string, n int) string {
	if len(str) < n {
		return str
	}
	if n < 0 {
		n = 0
	}
	return str[:n]
}

// substring at right most n characters
func Right(str string, n int) string {
	if len(str) < n {
		return str
	}
	if n < 0 {
		n = 0
	}
	return str[(len(str) - n):]
}

// substring at set left right n characters
func Mid(str string, left int, length int) string {
	if len(str) < left {
		return str
	}
	if left < 0 {
		left = 0
	}
	if length < 0 {
		return ``
	}
	if len(str) < (left + length) {
		return str[left:]
	}
	return str[left : left+length]
}

// substring before last `substr`
func LeftOfLast(str, substr string) string {
	len := strings.LastIndex(str, substr)
	if len < 0 {
		return str
	}
	return str[:len]
}

// substring after last `substr`
func RightOfLast(str, substr string) string {
	pos := strings.LastIndex(str, substr)
	if pos < 0 {
		return str
	}
	return str[pos+len(substr):]
}

// remove last n character, not UTF-8 friendly
func RemoveLastN(str string, n int) string {
	m := len(str)
	if n >= m {
		return ``
	}
	return str[0 : m-n]
}

// concat if not empty with additional separator
func ConcatIfNotEmpty(str, sep string) string {
	if str == `` {
		return ``
	}
	return str + sep
}
