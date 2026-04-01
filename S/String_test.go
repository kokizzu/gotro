package S

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/vmihailenco/msgpack/v5"
	"github.com/zeebo/assert"

	"github.com/kokizzu/gotro/L"
)

// Testing positive case function EqualsIgnoreCase //
func TestPositiveEqualsIgnoreCaseFunc(t *testing.T) {
	firstInput := "Bola"
	secondInput := "bola"

	result := EqualsIgnoreCase(firstInput, secondInput)

	if !result {
		t.Errorf(`%s %s should equal`, firstInput, secondInput)
	}
}

// Testing negative case function EqualsIgnoreCase //
func TestNegativeEqualsIgnoreCaseFunc(t *testing.T) {
	firstInput := "Bola"
	secondInput := "balon"

	result := EqualsIgnoreCase(firstInput, secondInput)

	if result {
		t.Errorf(`%s %s should not equal`, firstInput, secondInput)
	}
}

func TestMidFunction(t *testing.T) {
	text := "ABCDE"
	L.Print(`1` + Mid(text, -1, 1))
	L.Print(`2` + Mid(text, 0, -1))
	L.Print(`3` + Mid(text, 3, 1))
	L.Print(`4` + Mid(text, 4, 5))
	L.Print(`5` + Mid(text, 6, 1))
	L.Print(`6` + Mid(text, 2, 3))
}

func TestMsgpToMap(t *testing.T) {
	map0 := map[string]any{"a": 1, "b": 2}
	msgpStr, _ := msgpack.Marshal(map0)
	map1 := MsgpToMap(msgpStr)
	assert.Equal(t, fmt.Sprint(map0), fmt.Sprint(map1))
}

func TestMsgpToStrStrMap(t *testing.T) {
	map0 := map[string]string{"a": "1", "b": "2"}
	msgpStr, _ := msgpack.Marshal(map0)
	map1 := MsgpToStrStrMap(msgpStr)
	assert.Equal(t, fmt.Sprint(map0), fmt.Sprint(map1))
}

func TestMsgpToArr(t *testing.T) {
	arr0 := []any{1, "test"}
	msgpStr, _ := msgpack.Marshal(arr0)
	arr1 := MsgpToArr(msgpStr)
	assert.Equal(t, fmt.Sprint(arr0), fmt.Sprint(arr1))
}

func TestMsgpToObjArr(t *testing.T) {
	arr0 := []map[string]any{{"a": 1, "b": 2}, {"c": 3, "b": 4}}
	msgpStr, _ := msgpack.Marshal(arr0)
	arr1 := MsgpToObjArr(msgpStr)
	assert.Equal(t, fmt.Sprint(arr0), fmt.Sprint(arr1))
}

func TestMsgpToStrArr(t *testing.T) {
	arr0 := []string{"a", "b"}
	msgpStr, _ := msgpack.Marshal(arr0)
	arr1 := MsgpToStrArr(msgpStr)
	assert.Equal(t, fmt.Sprint(arr0), fmt.Sprint(arr1))
}

func TestMsgpToIntArr(t *testing.T) {
	arr0 := []int{1, 2}
	msgpStr, _ := msgpack.Marshal(arr0)
	arr1 := MsgpToIntArr(msgpStr)
	assert.Equal(t, fmt.Sprint(arr0), fmt.Sprint(arr1))
}

func TestMsgpAsMap(t *testing.T) {
	map0 := map[string]any{"a": 1, "b": 2}
	msgpStr, _ := msgpack.Marshal(map0)
	map1, ok := MsgpAsMap(msgpStr)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprint(map0), fmt.Sprint(map1))
}

func TestMsgpAsArr(t *testing.T) {
	arr0 := []any{1, "test"}
	msgpStr, _ := msgpack.Marshal(arr0)
	arr1, ok := MsgpAsArr(msgpStr)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprint(arr0), fmt.Sprint(arr1))
}

func TestMsgpAsStrArr(t *testing.T) {
	arr0 := []string{"a", "b"}
	msgpStr, _ := msgpack.Marshal(arr0)
	arr1, ok := MsgpAsStrArr(msgpStr)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprint(arr0), fmt.Sprint(arr1))
}

func TestMsgpAsIntArr(t *testing.T) {
	arr0 := []int{1, 2}
	msgpStr, _ := msgpack.Marshal(arr0)
	arr1, ok := MsgpAsIntArr(msgpStr)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprint(arr0), fmt.Sprint(arr1))
}

func TestMsgpAsFloatArr(t *testing.T) {
	arr0 := []float64{1, 2.3}
	msgpStr, _ := msgpack.Marshal(arr0)
	arr1, ok := MsgpAsFloatArr(msgpStr)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprint(arr0), fmt.Sprint(arr1))
}

func TestDelLeft(t *testing.T) {
	assert.Equal(t, DelLeft(`test123`, 3), `t123`)
	assert.Equal(t, DelLeft(`test123`, 10), ``)
	assert.Equal(t, DelLeft(`test123`, 0), `test123`)
	assert.Equal(t, DelLeft(`test123`, -1), `test123`)
}

func TestDelRight(t *testing.T) {
	assert.Equal(t, DelRight(`test123`, 3), `test`)
	assert.Equal(t, DelRight(`test123`, 10), ``)
	assert.Equal(t, DelRight(`test123`, 0), `test123`)
	assert.Equal(t, DelRight(`test123`, -1), `test123`)
}

func TestBasicStringHelpers(t *testing.T) {
	if !StartsWith("gotro", "go") || StartsWith("gotro", "xx") {
		t.Fatalf("StartsWith mismatch")
	}
	if !EndsWith("gotro", "tro") || EndsWith("gotro", "xx") {
		t.Fatalf("EndsWith mismatch")
	}
	if !Contains("gotro", "otr") || Contains("gotro", "xyz") {
		t.Fatalf("Contains mismatch")
	}
	if !Equals("abc", "abc") || Equals("abc", "ABC") {
		t.Fatalf("Equals mismatch")
	}
	if Count("banana", "an") != 2 {
		t.Fatalf("Count mismatch")
	}
	if TrimChars("..abc..", ".") != "abc" {
		t.Fatalf("TrimChars mismatch")
	}
	if IndexOf("abcabc", "bc") != 1 || LastIndexOf("abcabc", "bc") != 4 {
		t.Fatalf("IndexOf/LastIndexOf mismatch")
	}
	if ToUpper("abC") != "ABC" || ToLower("AbC") != "abc" {
		t.Fatalf("ToUpper/ToLower mismatch")
	}
	if CharAt("Halo 世界", 5) != "世" || CharAt("abc", 9) != "" {
		t.Fatalf("CharAt mismatch")
	}
	if RemoveCharAt("Halo 世界", 5) != "Halo 界" || RemoveCharAt("abc", 9) != "abc" {
		t.Fatalf("RemoveCharAt mismatch")
	}
	if ToTitle("hello world") != "Hello World" {
		t.Fatalf("ToTitle mismatch")
	}
	if If(true, "x") != "x" || If(false, "x") != "" {
		t.Fatalf("If mismatch")
	}
	if IfElse(true, "x", "y") != "x" || IfElse(false, "x", "y") != "y" {
		t.Fatalf("IfElse mismatch")
	}
	if IfEmpty("", "y") != "y" || IfEmpty("x", "y") != "x" {
		t.Fatalf("IfEmpty mismatch")
	}
	if Coalesce("", "", "z") != "z" || Coalesce("", "") != "" {
		t.Fatalf("Coalesce mismatch")
	}
}

func TestNumericAndJsonHelpers(t *testing.T) {
	if ToU("123") != 123 || ToI("123") != 123 || ToInt("123") != 123 {
		t.Fatalf("ToU/ToI/ToInt mismatch")
	}
	if v, ok := AsU("123"); !ok || v != 123 {
		t.Fatalf("AsU valid mismatch: %v %v", v, ok)
	}
	if _, ok := AsU("x"); ok {
		t.Fatalf("AsU invalid should fail")
	}
	if v, ok := AsI("123"); !ok || v != 123 {
		t.Fatalf("AsI valid mismatch: %v %v", v, ok)
	}
	if _, ok := AsI("x"); ok {
		t.Fatalf("AsI invalid should fail")
	}
	if ToF("12.5") != 12.5 {
		t.Fatalf("ToF mismatch")
	}
	if v, ok := AsF("12.5"); !ok || v != 12.5 {
		t.Fatalf("AsF valid mismatch: %v %v", v, ok)
	}
	if _, ok := AsF("x"); ok {
		t.Fatalf("AsF invalid should fail")
	}

	if got := JsonToMap(`{"a":1}`); fmt.Sprint(got) != "map[a:1]" {
		t.Fatalf("JsonToMap mismatch: %#v", got)
	}
	if got := JsonToStrStrMap(`{"a":"1"}`); got["a"] != "1" {
		t.Fatalf("JsonToStrStrMap mismatch: %#v", got)
	}
	if got := JsonToArr(`[1,"x"]`); fmt.Sprint(got) != "[1 x]" {
		t.Fatalf("JsonToArr mismatch: %#v", got)
	}
	if got := JsonToObjArr(`[{"a":1}]`); fmt.Sprint(got) != "[map[a:1]]" {
		t.Fatalf("JsonToObjArr mismatch: %#v", got)
	}
	if got := JsonToStrArr(`["a","b"]`); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Fatalf("JsonToStrArr mismatch: %#v", got)
	}
	if got := JsonToIntArr(`[1,2]`); !reflect.DeepEqual(got, []int64{1, 2}) {
		t.Fatalf("JsonToIntArr mismatch: %#v", got)
	}

	if _, ok := JsonAsMap(`{"a":1}`); !ok {
		t.Fatalf("JsonAsMap should be ok")
	}
	if _, ok := JsonAsMap(`{`); ok {
		t.Fatalf("JsonAsMap invalid should fail")
	}
	if _, ok := JsonAsArr(`[1]`); !ok {
		t.Fatalf("JsonAsArr should be ok")
	}
	if _, ok := JsonAsArr(`[`); ok {
		t.Fatalf("JsonAsArr invalid should fail")
	}
	if _, ok := JsonAsStrArr(`["a"]`); !ok {
		t.Fatalf("JsonAsStrArr should be ok")
	}
	if _, ok := JsonAsStrArr(`[1]`); ok {
		t.Fatalf("JsonAsStrArr invalid should fail")
	}
	if _, ok := JsonAsIntArr(`[1]`); !ok {
		t.Fatalf("JsonAsIntArr should be ok")
	}
	if _, ok := JsonAsIntArr(`["a"]`); ok {
		t.Fatalf("JsonAsIntArr invalid should fail")
	}
	if _, ok := JsonAsFloatArr(`[1.2]`); !ok {
		t.Fatalf("JsonAsFloatArr should be ok")
	}
	if _, ok := JsonAsFloatArr(`["a"]`); ok {
		t.Fatalf("JsonAsFloatArr invalid should fail")
	}
}

func TestSplitPadValidateAndBoundaries(t *testing.T) {
	if got := Split("a,b,c", ","); !reflect.DeepEqual(got, []string{"a", "b", "c"}) {
		t.Fatalf("Split mismatch: %#v", got)
	}
	if got := SplitFunc("a|b,c", func(r rune) bool { return r == '|' || r == ',' }); !reflect.DeepEqual(got, []string{"a", "b", "c"}) {
		t.Fatalf("SplitFunc mismatch: %#v", got)
	}
	if PadLeft("x", "0", 3) != "00x" || PadRight("x", "0", 3) != "x00" {
		t.Fatalf("PadLeft/PadRight mismatch")
	}

	if ValidateMailContact("Foo,Bar.<X>(Y)@") != "Foo_Bar__X__Y__" {
		t.Fatalf("ValidateMailContact mismatch")
	}
	if got := MergeMailContactEmails("Foo,Bar", "a@x.com, b@y.com , "); !reflect.DeepEqual(got, []string{"Foo_Bar<a@x.com>", "Foo_Bar<b@y.com>"}) {
		t.Fatalf("MergeMailContactEmails mismatch: %#v", got)
	}

	if ValidateIdent("ab-c_1!") != "abc_1" {
		t.Fatalf("ValidateIdent mismatch")
	}
	if ValidateEmail("a.b+1@example.com") == "" || ValidateEmail("bad@@example.com") != "" || ValidateEmail("a@exa^mple.com") != "" {
		t.Fatalf("ValidateEmail mismatch")
	}
	if ValidatePhone("+62 (812)-a3") != "+62 812-3" {
		t.Fatalf("ValidatePhone mismatch")
	}
	if ValidateFilename("ab/c:d, e.txt") != "abcd, e.txt" {
		t.Fatalf("ValidateFilename mismatch")
	}
	if len(RandomPassword(12)) != 12 {
		t.Fatalf("RandomPassword length mismatch")
	}

	if got := SplitN("abcdef", 2); !reflect.DeepEqual(got, []string{"ab", "cd", "ef"}) {
		t.Fatalf("SplitN mismatch: %#v", got)
	}
	if got := SplitN("ab", 10); !reflect.DeepEqual(got, []string{"ab"}) {
		t.Fatalf("SplitN short mismatch: %#v", got)
	}

	if LeftOf("a/b/c", "/") != "a" || LeftOf("abc", "/") != "abc" {
		t.Fatalf("LeftOf mismatch")
	}
	if RightOf("a/b/c", "/") != "b/c" || RightOf("abc", "/") != "abc" {
		t.Fatalf("RightOf mismatch")
	}
	if LeftN("abcdef", 3) != "abc..." || LeftN("ab", 3) != "ab" {
		t.Fatalf("LeftN mismatch")
	}
	if Left("abcdef", 3) != "abc" || Left("abcdef", -1) != "" {
		t.Fatalf("Left mismatch")
	}
	if Right("abcdef", 3) != "def" || Right("abcdef", -1) != "" {
		t.Fatalf("Right mismatch")
	}
	if LeftOfLast("a/b/c", "/") != "a/b" || RightOfLast("a/b/c", "/") != "c" {
		t.Fatalf("LeftOfLast/RightOfLast mismatch")
	}
	if RemoveLastN("abcdef", 2) != "abcd" || RemoveLastN("ab", 2) != "" {
		t.Fatalf("RemoveLastN mismatch")
	}
	if ConcatIfNotEmpty("a", ",") != "a," || ConcatIfNotEmpty("", ",") != "" {
		t.Fatalf("ConcatIfNotEmpty mismatch")
	}
	if LowerFirst("AbC") != "abC" || UpperFirst("abC") != "AbC" || LowerFirst("") != "" || UpperFirst("") != "" {
		t.Fatalf("LowerFirst/UpperFirst mismatch")
	}
}

func TestCaseConversions(t *testing.T) {
	if PascalCase("hello_world-test.1foo") != "HelloWorldTest1Foo" {
		t.Fatalf("PascalCase mismatch: %q", PascalCase("hello_world-test.1foo"))
	}
	if CamelCase("Hello_world-test.1foo") != "helloWorldTest1Foo" {
		t.Fatalf("CamelCase mismatch: %q", CamelCase("Hello_world-test.1foo"))
	}
	if SnakeCase("HelloWorldJSON2Data") != "hello_world_json_2_data" {
		t.Fatalf("SnakeCase mismatch: %q", SnakeCase("HelloWorldJSON2Data"))
	}
	if !FirstIsLower("abc") || FirstIsLower("Abc") || FirstIsLower("") {
		t.Fatalf("FirstIsLower mismatch")
	}
}
