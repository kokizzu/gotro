package S

import (
	"fmt"
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
