package S

// String support package
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/kokizzu/rand"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/kokizzu/json5b/encoding/json5b"

	"github.com/kokizzu/gotro/C"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
)

const WebBR = "\n<br/>"

// StartsWith check whether the input string (first arg) starts with a certain character (second arg) or not.
//
//	S.StartsWith(`adakah`,`ad`) // bool(true)
//	S.StartsWith(`adakah`,`bad`) // bool(false)
func StartsWith(str, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}

// EndsWith check whether the input string (first arg) ends with a certain character (second arg) or not.
//
//	S.EndsWith(`adakah`,`ah`)) // bool(true)
//	S.EndsWith(`adakah`,`aka`)) // bool(false)
func EndsWith(str, suffix string) bool {
	return strings.HasSuffix(str, suffix)
}

// Contains check whether the input string (first arg) contains a certain sub string (second arg) or not.
//
//	S.Contains(`komputer`,`om`)) // bool(true)
//	S.Contains(`komputer`,`opu`)) // bool(false)
func Contains(str, substr string) bool {
	return strings.Contains(str, substr)
}

// Equals compare two input string (first arg) equal with another input string (second arg).
//
//	S.Equals(`komputer`,`komputer`)) // bool(true)
//	S.Equals(`komputer`,`Komputer`)) // bool(false)
func Equals(strFirst, strSecond string) bool {
	return strFirst == strSecond
}

// EqualsIgnoreCase compare two input string (first arg) equal with ignoring case another input string (second arg).
//
//	S.EqualsIgnoreCase(`komputer`,`komputer`)) // bool(true)
//	S.EqualsIgnoreCase(`komputer`,`Komputer`)) // bool(true)
func EqualsIgnoreCase(strFirst, strSecond string) bool {
	strFirst = ToLower(strFirst)
	strSecond = ToLower(strSecond)
	return strFirst == strSecond
}

// Count count how many specific character (first arg) that the string (second arg) contains
//
//	S.Count(`komputeer`,`e`))// output int(2)
func Count(str, substr string) int {
	return strings.Count(str, substr)
}

// Trim erase spaces from left and right
//
//	S.Trim(` withtrim:  `) // `withtrim:`
func Trim(str string) string {
	return strings.TrimSpace(str)
}

// TrimChars remove chars from beginning and end
//
//	S.TrimChars(`aoaaffoa`,`ao`) // `ff`
func TrimChars(str, chars string) string {
	return strings.Trim(str, chars)
}

// IndexOf get first index of
// S.IndexOf(`abcdcd`,`c) // 2, -1 if not exists
func IndexOf(str, sub string) int {
	return strings.Index(str, sub)
}

// LastIndexOf get last index of
//
//	S.LastIndexOf(`abcdcd`,`c`) // 4, -1 if not exists
func LastIndexOf(str, sub string) int {
	return strings.LastIndex(str, sub)
}

// Replace replace all substring with another substring
//
//	S.Replace(`bisa`,`is`,`us`) // `busa`
func Replace(haystack, needle, gold string) string {
	return strings.Replace(haystack, needle, gold, -1)
}

// ToLower change the characters in string to lowercase
//
//	S.ToLower(`BIsa`) // "bisa"
func ToLower(str string) string {
	return strings.ToLower(str)
}

// ToUpper change the characters in string to uppercase
// S.ToUpper(`bisa`) // "BISA"
func ToUpper(str string) string {
	return strings.ToUpper(str)
}

// CharAt get character at specific index, utf-8 safe
//
//	S.CharAt(`Halo 世界`, 5) // `世` // utf-8 example, if characters not shown, it's probably your font/editor/plugin
//	S.CharAt(`Halo`, 3) // `o`
func CharAt(str string, index int) string {
	for in, ch := range str {
		if in == index {
			return string(ch)
		}
	}
	return ``
}

// RemoveCharAt remove character at specific index, utf-8 safe
//
//	S.RemoveCharAt(`Halo 世界`, 5) // `Halo 界` --> utf-8 example, if characters not shown, it's probably your font/editor/plugin
//	S.RemoveCharAt(`Halo`, 3) // `Hal`
func RemoveCharAt(str string, index int) string {
	var chars []byte
	for in, ch := range str {
		if index != in {
			chars = append(chars, string(ch)...)
		}
	}
	return string(chars)
}

var enCaser = cases.Title(language.English)

// ToTitle Change first letter for every word to uppercase
//
//	S.ToTitle(`Disa dasi`)) // output "Disa Dasi"
func ToTitle(str string) string {
	return enCaser.String(str)
}

// If simplified ternary operator (bool ? val : 0), returns second argument, if the condition (first arg) is true, returns empty string if not
//
//	S.If(true,`a`) // `a`
//	S.If(false,`a`) // ``
func If(b bool, yes string) string {
	if b {
		return yes
	}
	return ``
}

// IfElse ternary operator (bool ? val1 : val2), returns second argument if the condition (first arg) is true, third argument if not
//
//	S.IfElse(true,`a`,`b`) // `a`
//	S.IfElse(false,`a`,`b`) // `b`
func IfElse(b bool, yes, no string) string {
	if b {
		return yes
	}
	return no
}

// IfEmpty coalesce, return first non-empty string
//
//	S.IfEmpty(``,`2`) // `2`
//	S.IfEmpty(`1`,`2`) // `1`
func IfEmpty(str1, str2 string) string {
	if str1 != `` {
		return str1
	}
	return str2
}

// Coalesce coalesce, return first non-empty string
//
//	S.Coalesce(`1`,`2`) // `1`
//	S.Coalesce(``,`2`) // `2`
//	S.Coalesce(``,``,`3`) // `3`
func Coalesce(strs ...string) string {
	for _, str := range strs {
		if str != `` {
			return str
		}
	}
	return ``
}

// ToU convert string to uint64, returns 0 and silently print error if not valid
//
//	S.ToU(`1234`) // 1234
//	S.ToU(`1a`) // 0
func ToU(str string) uint64 {
	val, _ := strconv.ParseUint(str, 10, 64)
	//L.IsError(err, str)
	return val
}

// ToI convert string to int64, returns 0 and silently print error if not valid
//
//	S.ToI(`1234`) // 1234
//	S.ToI(`1a`) // 0
func ToI(str string) int64 {
	val, _ := strconv.ParseInt(str, 10, 64)
	//L.IsError(err, str)
	return val
}

// ToInt convert string to int, returns 0 and silently print error if not valid
//
//	S.ToInt(`1234`) // 1234
//	S.ToInt(`1a`) // 0
func ToInt(str string) int {
	val, _ := strconv.ParseInt(str, 10, 64)
	//L.IsError(err, str)
	return int(val)
}

// AsU convert to uint with check
//
//	S.AsU(`1234`) // 1234, true
//	S.AsU(`1abc`) // 0, false
func AsU(str string) (uint, bool) {
	res, err := strconv.ParseInt(str, 10, 64)
	return uint(res), err == nil
}

// AsI convert to int64 with check
//
//	S.AsI(`1234`) // 1234, true
//	S.AsI(`1abc`) // 0, false
func AsI(str string) (int64, bool) {
	res, err := strconv.ParseInt(str, 10, 64)
	return res, err == nil
}

// ToF convert string to float64, returns 0 and silently print error if not valid
//
//	S.ToF(`1234.5`) // 1234.5
//	S.ToF(`1a`) // 0.0
func ToF(str string) float64 {
	val, _ := strconv.ParseFloat(str, 64)
	//L.IsError(err, str)
	return val
}

// AsF convert to float64 with check
//
//	S.AsF(`1234.5`) // 1234.5, true
//	S.AsF(`1abc`) // 0.0, false
func AsF(str string) (float64, bool) {
	res, err := strconv.ParseFloat(str, 64)
	return res, err == nil
}

// JsonToMap convert JSON object to map[string]any, silently print and return empty map if failed
//
//	jsonStr := `{"test":123,"bla":[1,2,3,4]}`
//	map1 := S.JsonToMap(jsonStr)
func JsonToMap(str string) (res map[string]any) {
	res = map[string]any{}
	if len(str) == 0 {
		return
	}
	err := json.Unmarshal([]byte(str), &res)
	L.IsError(err, str)
	return
}

// MsgpToMap convert MsgPack object to map[string]any, silently print and return empty map if failed
//
//	msgpStr := []byte(`��buah{�angka�dia`)
//	map1 := S.MsgpToMap(msgpStr)
func MsgpToMap(str []byte) (res map[string]any) {
	res = map[string]any{}
	if len(str) == 0 {
		return
	}
	err := msgpack.Unmarshal(str, &res)
	L.IsError(err, fmt.Sprint(str))
	return
}

// JsonToStrStrMap convert JSON object to map[string]string, silently print and return empty map if failed
//
//	jsonStr := `{"test":123,"bla":[1,2,3,4]}`
//	map1 := S.JsonToMap(jsonStr)
func JsonToStrStrMap(str string) (res map[string]string) {
	res = map[string]string{}
	if len(str) == 0 {
		return
	}
	err := json.Unmarshal([]byte(str), &res)
	L.IsError(err, str)
	return
}

// MsgpToStrStrMap convert Msgpack object to map[string]string, silently print and return empty map if failed
//
//	map1 := S.MsgpToStrStrMap(msgpStr)
func MsgpToStrStrMap(str []byte) (res map[string]string) {
	res = map[string]string{}
	if len(str) == 0 {
		return
	}
	err := msgpack.Unmarshal(str, &res)
	L.IsError(err, fmt.Sprint(str))
	return
}

// JsonToArr convert JSON object to []any, silently print and return empty slice of interface if failed
//
//	jsonStr := `[1,2,['test'],'a']`
//	arr := S.JsonToArr(jsonStr)
func JsonToArr(str string) (res []any) {
	res = []any{}
	if len(str) == 0 {
		return
	}
	err := json.Unmarshal([]byte(str), &res)
	L.IsError(err, str)
	return
}

// MsgpToArr convert Msgpack object to []any, silently print and return empty slice of interface if failed
//
//	arr := S.MsgpToArr(msgpStr)
func MsgpToArr(str []byte) (res []any) {
	res = []any{}
	if len(str) == 0 {
		return
	}
	err := msgpack.Unmarshal(str, &res)
	L.IsError(err, fmt.Sprint(str))
	return
}

// JsonToObjArr convert JSON object to []map[string]any, silently print and return empty slice of interface if failed
//
//	jsonStr := `[{"x":"foo"},{"y":"bar"}]`
//	arr := S.JsonToObjArr(jsonStr)
func JsonToObjArr(str string) (res []map[string]any) {
	res = []map[string]any{}
	if len(str) == 0 {
		return
	}
	err := json.Unmarshal([]byte(str), &res)
	L.IsError(err, str)
	return
}

// MsgpToObjArr convert Msgpack object to []map[string]any, silently print and return empty slice of interface if failed
//
//	arr := S.MsgpToObjArr(msgpStr)
func MsgpToObjArr(str []byte) (res []map[string]any) {
	res = []map[string]any{}
	if len(str) == 0 {
		return
	}
	err := msgpack.Unmarshal(str, &res)
	L.IsError(err, fmt.Sprint(str))
	return
}

// JsonToStrArr convert JSON object to []string, silently print and return empty slice of interface if failed
//
//	jsonStr := `["123","456",789]`
//	arr := S.JsonToStrArr(jsonStr)
func JsonToStrArr(str string) (res []string) {
	res = []string{}
	if len(str) == 0 {
		return
	}
	err := json5b.Unmarshal([]byte(str), &res)
	L.IsError(err, str)
	return
}

// MsgpToStrArr convert Msgpack object to []string, silently print and return empty slice of interface if failed
//
//	arr := S.MsgpToStrArr(msgpStr)
func MsgpToStrArr(str []byte) (res []string) {
	res = []string{}
	if len(str) == 0 {
		return
	}
	err := msgpack.Unmarshal(str, &res)
	L.IsError(err, fmt.Sprint(str))
	return
}

// JsonToIntArr convert JSON object to []int64, silently print and return empty slice of interface if failed
//
//	jsonStr := `[1,2,['test'],'a']`
//	arr := S.JsonToArr(jsonStr)
func JsonToIntArr(str string) (res []int64) {
	res = []int64{}
	if len(str) == 0 {
		return
	}
	err := json5b.Unmarshal([]byte(str), &res)
	L.IsError(err, str)
	return
}

// MsgpToIntArr convert Msgpack object to []int64, silently print and return empty slice of interface if failed
//
//	arr := S.MsgpToIntArr(msgpStr)
func MsgpToIntArr(str []byte) (res []int64) {
	res = []int64{}
	if len(str) == 0 {
		return
	}
	err := msgpack.Unmarshal(str, &res)
	L.IsError(err, fmt.Sprint(str))
	return
}

// repeat string
func Repeat(str string, count int) string {
	return strings.Repeat(str, count)
}

// JsonAsMap convert JSON object to map[string]any with check
//
//	jsonStr := `{"test":123,"bla":[1,2,3,4]}`
//	map1, ok := S.JsonAsMap(jsonStr)
func JsonAsMap(str string) (res map[string]any, ok bool) {
	res = map[string]any{}
	err := json.Unmarshal([]byte(str), &res)
	ok = err == nil
	return
}

// MsgpAsMap convert Msgpack object to map[string]any with check
//
//	map1, ok := S.MsgpAsMap(msgpStr)
func MsgpAsMap(str []byte) (res map[string]any, ok bool) {
	res = map[string]any{}
	err := msgpack.Unmarshal(str, &res)
	ok = err == nil
	return
}

// JsonAsArr convert JSON object to []any with check
//
//	jsonStr := `[1,2,['test'],'a']`
//	arr, ok := S.JsonAsArr(jsonStr)
func JsonAsArr(str string) (res []any, ok bool) {
	res = []any{}
	err := json.Unmarshal([]byte(str), &res)
	ok = err == nil
	return
}

// MsgpAsArr convert Msgpack object to []any with check
//
//	arr, ok := S.MsgpAsArr(msgpStr)
func MsgpAsArr(str []byte) (res []any, ok bool) {
	res = []any{}
	err := msgpack.Unmarshal(str, &res)
	ok = err == nil
	return
}

// JsonAsStrArr convert JSON object to []string with check
//
//	jsonStr := `["a","b","c"]`
//	arr, ok := S.JsonAsStrArr(jsonStr)
func JsonAsStrArr(str string) (res []string, ok bool) {
	res = []string{}
	err := json.Unmarshal([]byte(str), &res)
	ok = err == nil
	return
}

// MsgpAsStrArr convert Msgpack object to []string with check
//
//	arr, ok := S.MsgpAsStrArr(msgpStr)
func MsgpAsStrArr(str []byte) (res []string, ok bool) {
	res = []string{}
	err := msgpack.Unmarshal(str, &res)
	ok = err == nil
	return
}

// JsonAsIntArr convert JSON object to []int64 with check
//
//	jsonStr := `[1,2,3]`
//	arr, ok := S.JsonAsIntArr(jsonStr)
func JsonAsIntArr(str string) (res []int64, ok bool) {
	res = []int64{}
	err := json.Unmarshal([]byte(str), &res)
	ok = err == nil
	return
}

// MsgpAsIntArr convert Msgpack object to []int64 with check
//
//	arr, ok := S.MsgpAsIntArr(msgpStr)
func MsgpAsIntArr(str []byte) (res []int64, ok bool) {
	res = []int64{}
	err := msgpack.Unmarshal(str, &res)
	ok = err == nil
	return
}

// JsonAsFloatArr convert JSON object to []float64 with check
//
//	jsonStr := `[1,2,3]`
//	arr, ok := S.JsonAsFloatArr(jsonStr)
func JsonAsFloatArr(str string) (res []float64, ok bool) {
	res = []float64{}
	err := json.Unmarshal([]byte(str), &res)
	ok = err == nil
	return
}

// MsgpAsFloatArr convert Msgpack object to []float64 with check
//
//	arr, ok := S.MsgpAsFloatArr(msgpStr)
func MsgpAsFloatArr(str []byte) (res []float64, ok bool) {
	res = []float64{}
	err := msgpack.Unmarshal(str, &res)
	ok = err == nil
	return
}

// Split split a string (first arg) by characters (second arg) into array of strings (output).
//
//	S.Split(`biiiissssa`,func(ch rune) bool { return ch == `i` }) // output []string{"b", "", "", "", "ssssa"}
func Split(str, sep string) []string {
	return strings.Split(str, sep)
}

// SplitFunc split a string (first arg) based on a function
func SplitFunc(str string, fun func(rune) bool) []string {
	return strings.FieldsFunc(str, fun)
}

// PadLeft append padStr to left until length is lenStr
func PadLeft(s string, padStr string, lenStr int) string {
	padCount := I.MaxOf(lenStr-len(s), 0)
	return strings.Repeat(padStr, padCount) + s
}

// PadRight append padStr to right until length is lenStr
func PadRight(s string, padStr string, lenStr int) string {
	padCount := I.MaxOf(lenStr-len(s), 0)
	return s + strings.Repeat(padStr, padCount)
}

// ValidateMailContact return valid version of mail contact (part before <usr@email>)
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

// MergeMailContactEmails return formatted array of mail contact <usr@email>
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

// ValidateIdent return valid identifier (a-zA-Z0-9_)
func ValidateIdent(str string) string {
	res := strings.Map(func(r rune) rune {
		if C.IsIdent(r) {
			return r
		}
		return -1
	}, str)
	return res
}

// ValidateEmail return empty string if str is not a valid email
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

// ValidatePhone remove invalid characters of a phone number
func ValidatePhone(str string) string {
	res := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) || r == '+' || r == ' ' || r == '-' {
			return r
		}
		return -1
	}, str)
	return res
}

// ValidateFilename validate file name
func ValidateFilename(str string) string {
	res := strings.Map(func(r rune) rune {
		if C.IsValidFilename(byte(r)) {
			return r
		}
		return -1
	}, str)
	return res
}

// RandomPassword create a random password
func RandomPassword(strlen int64) string {
	const chars = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ123456789" // l and I removed
	result := make([]byte, strlen)
	for i := int64(0); i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// SplitN split to substrings with maximum n characters
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

// LeftOf substring before first `substr`
func LeftOf(str, substr string) string {
	len := strings.Index(str, substr)
	if len < 0 {
		return str
	}
	return str[:len]
}

// RightOf substring after first `substr`
func RightOf(str, substr string) string {
	pos := strings.Index(str, substr)
	if pos < 0 {
		return str
	}
	return str[pos+len(substr):]
}

// LeftN substring at most n characters
func LeftN(str string, n int) string {
	if len(str) < n {
		return str
	}
	return str[:n] + `...`
}

// Left substring at most n characters
func Left(str string, n int) string {
	if len(str) < n {
		return str
	}
	if n < 0 {
		n = 0
	}
	return str[:n]
}

// Right substring at right most n characters
func Right(str string, n int) string {
	if len(str) < n {
		return str
	}
	if n < 0 {
		n = 0
	}
	return str[(len(str) - n):]
}

// DelLeft return remainder after N first characters removed, return empty if N more than len, not unicode-safe
func DelLeft(str string, n int) string {
	if len(str) < n {
		return ``
	}
	if n <= 0 {
		return str
	}
	return str[n:]
}

// DelRight return remainder after N last characters deleted, return empty if N more than len, not unicode-safe
func DelRight(str string, n int) string {
	if len(str) < n {
		return ``
	}
	if n <= 0 {
		return str
	}
	return str[:len(str)-n]
}

// Mid substring at set left right n characters
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

// LeftOfLast substring before last `substr`
func LeftOfLast(str, substr string) string {
	len := strings.LastIndex(str, substr)
	if len < 0 {
		return str
	}
	return str[:len]
}

// RightOfLast substring after last `substr`
func RightOfLast(str, substr string) string {
	pos := strings.LastIndex(str, substr)
	if pos < 0 {
		return str
	}
	return str[pos+len(substr):]
}

// RemoveLastN remove last n character, not UTF-8 friendly
func RemoveLastN(str string, n int) string {
	m := len(str)
	if n >= m {
		return ``
	}
	return str[0 : m-n]
}

// ConcatIfNotEmpty concat if not empty with additional separator
func ConcatIfNotEmpty(str, sep string) string {
	if str == `` {
		return ``
	}
	return str + sep
}

// LowerFirst convert to lower only first char
func LowerFirst(s string) string {
	for i, v := range s {
		return string(unicode.ToLower(v)) + s[i+1:]
	}
	return ``
}

// UpperFirst convert to lower only first char
func UpperFirst(s string) string {
	for i, v := range s {
		return string(unicode.ToUpper(v)) + s[i+1:]
	}
	return ``
}

// PascalCase convert to PascalCase
// source: https://github.com/iancoleman/strcase
func PascalCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	n := strings.Builder{}
	n.Grow(len(s))
	capNext := true
	for i, v := range []byte(s) {
		isCap := v >= 'A' && v <= 'Z'
		isLow := v >= 'a' && v <= 'z'
		if capNext {
			if isLow {
				v += 'A'
				v -= 'a'
			}
		} else if i == 0 {
			if isCap {
				v += 'a'
				v -= 'A'
			}
		}
		if isCap || isLow {
			n.WriteByte(v)
			capNext = false
		} else if vIsNum := v >= '0' && v <= '9'; vIsNum {
			n.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return n.String()
}

// CamelCase convert to camelCase
func CamelCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	n := strings.Builder{}
	n.Grow(len(s))
	capNext := false
	for i, v := range []byte(s) {
		isCap := v >= 'A' && v <= 'Z'
		isLow := v >= 'a' && v <= 'z'
		if capNext {
			if isLow {
				v += 'A'
				v -= 'a'
			}
		} else if i == 0 {
			if isCap {
				v += 'a'
				v -= 'A'
			}
		}
		if isCap || isLow {
			n.WriteByte(v)
			capNext = false
		} else if vIsNum := v >= '0' && v <= '9'; vIsNum {
			n.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return n.String()
}

// SnakeCase convert to snake case
// source: https://github.com/iancoleman/strcase
func SnakeCase(s string) string {
	s = strings.TrimSpace(s)
	n := strings.Builder{}
	const delimiter = '_'
	ignore := byte(0)
	n.Grow(len(s) + 2) // nominal 2 bytes of extra space for inserted delimiters
	for i, v := range []byte(s) {
		isCap := v >= 'A' && v <= 'Z'
		isLow := v >= 'a' && v <= 'z'
		if isCap {
			v += 'a'
			v -= 'A'
		}

		// treat acronyms as words, eg for JSONData -> JSON is a whole word
		if i+1 < len(s) {
			next := s[i+1]
			vIsNum := v >= '0' && v <= '9'
			nextIsCap := next >= 'A' && next <= 'Z'
			nextIsLow := next >= 'a' && next <= 'z'
			nextIsNum := next >= '0' && next <= '9'
			// add underscore if next letter case type is changed
			if (isCap && (nextIsLow || nextIsNum)) || (isLow && (nextIsCap || nextIsNum)) || (vIsNum && (nextIsCap || nextIsLow)) {
				if prevIgnore := ignore > 0 && i > 0 && s[i-1] == ignore; !prevIgnore {
					if isCap && nextIsLow {
						if prevIsCap := i > 0 && s[i-1] >= 'A' && s[i-1] <= 'Z'; prevIsCap {
							n.WriteByte(delimiter)
						}
					}
					n.WriteByte(v)
					if isLow || vIsNum || nextIsNum {
						n.WriteByte(delimiter)
					}
					continue
				}
			}
		}

		if (v == ' ' || v == '_' || v == '-') && v != ignore {
			// replace space/underscore/hyphen with delimiter
			n.WriteByte(delimiter)
		} else {
			n.WriteByte(v)
		}
	}

	return n.String()
}

// FirstIsLower check first character is lowercase
func FirstIsLower(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[0] >= 'a' && s[0] <= 'z'
}
