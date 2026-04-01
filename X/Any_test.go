package X

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/M"
)

type testStringer struct{}

func (testStringer) String() string { return "stringer" }

func TestNumericConverters(t *testing.T) {
	i := int64(42)
	var wrapped any = `7`
	if ToU(uint64(12)) != 12 || ToU(&i) != 42 || ToU(&wrapped) != 7 || ToU(nil) != 0 {
		t.Fatalf("ToU mismatch")
	}

	b := uint8(255)
	if ToByte(`255`) != 255 || ToByte(&b) != 255 || ToByte(nil) != 0 {
		t.Fatalf("ToByte mismatch")
	}

	if ToI(`12.9`) != 12 || ToI(true) != 1 || ToI(nil) != 0 {
		t.Fatalf("ToI mismatch")
	}
	if ToF(`12.5`) != 12.5 || ToF(false) != 0 || ToF(nil) != 0 {
		t.Fatalf("ToF mismatch")
	}
}

func TestStringBoolAndArrayConverters(t *testing.T) {
	if ToS(12) != `12` || ToS(true) != `true` || ToS([]byte("abc")) != `abc` || ToS(testStringer{}) != `stringer` {
		t.Fatalf("ToS mismatch")
	}
	if ToS(nil) != `` {
		t.Fatalf("ToS(nil) should be empty")
	}

	if ToBool(`false`) || !ToBool(` yes `) || ToBool(0) || !ToBool(2) {
		t.Fatalf("ToBool mismatch")
	}

	src := []any{1, `2`}
	if got := ToArr(src); len(got) != 2 {
		t.Fatalf("ToArr mismatch: %#v", got)
	}
	if got := ToArr(`bad`); len(got) != 0 {
		t.Fatalf("ToArr invalid should return empty: %#v", got)
	}
	if got := ArrToStrArr(src); !reflect.DeepEqual(got, []string{`1`, `2`}) {
		t.Fatalf("ArrToStrArr mismatch: %#v", got)
	}
	if got := ArrToIntArr([]any{`3`, 4.8}); !reflect.DeepEqual(got, []int64{3, 4}) {
		t.Fatalf("ArrToIntArr mismatch: %#v", got)
	}
}

func TestTimeMsgpackAndMapConverters(t *testing.T) {
	expected := time.Date(2024, time.February, 29, 23, 59, 58, 123456000, time.UTC)
	if got := ToTime(`2024-02-29 23:59:58.123456`); !got.Equal(expected) {
		t.Fatalf("ToTime(string) mismatch: %v", got)
	}
	if got := ToTime([]byte(`2024-02-29`)); !got.Equal(time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("ToTime([]byte) mismatch: %v", got)
	}
	str := `2024-02-29 23:59:58`
	if got := ToTime(&str); !got.Equal(time.Date(2024, time.February, 29, 23, 59, 58, 0, time.UTC)) {
		t.Fatalf("ToTime(*string) mismatch: %v", got)
	}
	if got := ToTime(`bad`); !got.IsZero() {
		t.Fatalf("ToTime invalid should return zero time, got %v", got)
	}

	type payload struct {
		A int
	}
	raw := ToMsgp(payload{A: 7})
	var out payload
	if !FromMsgp(raw, &out) || out.A != 7 {
		t.Fatalf("FromMsgp valid payload mismatch: %#v", out)
	}
	if FromMsgp([]byte{0xff}, &out) {
		t.Fatalf("FromMsgp should fail on invalid msgpack")
	}

	if got := ToJson5(nil); got != `''` {
		t.Fatalf("ToJson5(nil) mismatch: %q", got)
	}
	if got := ToJson5(map[int64]bool{1: true}); got != `{1:true}` {
		t.Fatalf("ToJson5(map[int64]bool) mismatch: %q", got)
	}
	if got := ToJson(map[string]any{`a`: 1}); got != `{"a":1}` {
		t.Fatalf("ToJson mismatch: %q", got)
	}
	if got := ToJsonPretty(map[string]any{`a`: 1}); !strings.Contains(got, "\n") {
		t.Fatalf("ToJsonPretty should be multiline: %q", got)
	}
	if got := ToYaml(map[string]any{`a`: 1}); !strings.Contains(got, "a: 1") {
		t.Fatalf("ToYaml mismatch: %q", got)
	}

	if got := ToAX(A.X{1, 2}); len(got) != 2 {
		t.Fatalf("ToAX mismatch: %#v", got)
	}
	if got := ToAF([]any{`1.5`, 2}); !reflect.DeepEqual(got, []float64{1.5, 2}) {
		t.Fatalf("ToAF mismatch: %#v", got)
	}
	if got := ToMSX(map[string]any{`a`: 1}); got[`a`] != 1 {
		t.Fatalf("ToMSX mismatch: %#v", got)
	}
	if got := ToMSS(map[string]string{`a`: `b`}); got[`a`] != `b` {
		t.Fatalf("ToMSS mismatch: %#v", got)
	}
}

func TestTimeParsingHelpers(t *testing.T) {
	if got, err := parseByteYear([]byte(`2024`)); err != nil || got != 2024 {
		t.Fatalf("parseByteYear mismatch: got=%d err=%v", got, err)
	}
	if got, err := parseByte2Digits('4', '2'); err != nil || got != 42 {
		t.Fatalf("parseByte2Digits mismatch: got=%d err=%v", got, err)
	}
	if got, err := parseByteNanoSec([]byte(`123456`)); err != nil || got != 123456000 {
		t.Fatalf("parseByteNanoSec mismatch: got=%d err=%v", got, err)
	}
	if got, err := bToi('9'); err != nil || got != 9 {
		t.Fatalf("bToi mismatch: got=%d err=%v", got, err)
	}
	if _, err := bToi('x'); err == nil {
		t.Fatalf("bToi should fail on non-digit")
	}

	if _, err := parseDateTime([]byte(`2024-02-29 23:59:58`), time.UTC); err != nil {
		t.Fatalf("parseDateTime valid value should not fail: %v", err)
	}
	if _, err := parseDateTime([]byte(`2024/02/29`), time.UTC); err == nil {
		t.Fatalf("parseDateTime should fail on bad separators")
	}
	if _, err := parseDateTime([]byte(`2024-02`), time.UTC); err == nil {
		t.Fatalf("parseDateTime should fail on invalid length")
	}

	// Keep test compilation coverage for type aliases that are commonly passed into converters.
	_ = M.SX{}
}

type ptrStringer struct{}

func (*ptrStringer) String() string { return "ptr" }

type boolStr string

func (b boolStr) String() string { return string(b) }

func TestNumericConvertersMoreBranches(t *testing.T) {
	i := int(1)
	u := uint(2)
	i8 := int8(3)
	i16 := int16(4)
	i32 := int32(5)
	i64 := int64(6)
	u8 := uint8(7)
	u16 := uint16(8)
	u32 := uint32(9)
	u64 := uint64(10)
	f32 := float32(11.5)
	f64 := float64(12.5)
	dur := time.Duration(13)
	var anyWrap any = `14`
	var nilInt *int

	if ToU(&i) != 1 || ToU(&u) != 2 || ToU(&i8) != 3 || ToU(&i16) != 4 || ToU(&i32) != 5 ||
		ToU(&i64) != 6 || ToU(&u8) != 7 || ToU(&u16) != 8 || ToU(&u32) != 9 || ToU(&u64) != 10 ||
		ToU(&f32) != 11 || ToU(&f64) != 12 || ToU(dur) != 13 || ToU(&anyWrap) != 14 ||
		ToU([]byte(`15.9`)) != 15 || ToU(false) != 0 || ToU(nilInt) != 0 || ToU([]byte(`x`)) != 0 {
		t.Fatalf("ToU branch mismatch")
	}

	if ToByte(&i) != 1 || ToByte(&u) != 2 || ToByte(&i8) != 3 || ToByte(&i16) != 4 || ToByte(&i32) != 5 ||
		ToByte(&i64) != 6 || ToByte(&u8) != 7 || ToByte(&u16) != 8 || ToByte(&u32) != 9 || ToByte(&u64) != 10 ||
		ToByte(&f32) != 11 || ToByte(&f64) != 12 || ToByte(dur) != 13 || ToByte(&anyWrap) != 14 ||
		ToByte([]byte(`15.9`)) != 15 || ToByte(false) != 0 || ToByte(nilInt) != 0 || ToByte([]byte(`x`)) != 0 {
		t.Fatalf("ToByte branch mismatch")
	}

	if ToI(&i) != 1 || ToI(&u) != 2 || ToI(&i8) != 3 || ToI(&i16) != 4 || ToI(&i32) != 5 ||
		ToI(&i64) != 6 || ToI(&u8) != 7 || ToI(&u16) != 8 || ToI(&u32) != 9 || ToI(&u64) != 10 ||
		ToI(&f32) != 11 || ToI(&f64) != 12 || ToI(dur) != 13 || ToI(&anyWrap) != 14 ||
		ToI([]byte(`15.9`)) != 15 || ToI(false) != 0 || ToI(nilInt) != 0 || ToI([]byte(`x`)) != 0 {
		t.Fatalf("ToI branch mismatch")
	}

	if ToF(&i) != 1 || ToF(&u) != 2 || ToF(&i8) != 3 || ToF(&i16) != 4 || ToF(&i32) != 5 ||
		ToF(&i64) != 6 || ToF(&u8) != 7 || ToF(&u16) != 8 || ToF(&u32) != 9 || ToF(&u64) != 10 ||
		ToF(&f32) != 11.5 || ToF(&f64) != 12.5 || ToF(dur) != 13 || ToF(&anyWrap) != 14 ||
		ToF([]byte(`15.9`)) != 15.9 || ToF(false) != 0 || ToF(nilInt) != 0 || ToF([]byte(`x`)) != 0 {
		t.Fatalf("ToF branch mismatch")
	}
}

func TestStringTimeBoolConvertersMoreBranches(t *testing.T) {
	i := int(1)
	u := uint(2)
	i8 := int8(3)
	i16 := int16(4)
	i32 := int32(5)
	i64 := int64(6)
	u8 := uint8(7)
	u16 := uint16(8)
	u32 := uint32(9)
	u64 := uint64(10)
	f32 := float32(11.5)
	f64 := float64(12.5)
	var anyWrap any = int64(13)
	var nilPS *ptrStringer

	assertToS := func(name string, got string, want string) {
		t.Helper()
		if got != want {
			t.Fatalf("ToS %s mismatch: want=%q got=%q", name, want, got)
		}
	}
	assertToS("*int", ToS(&i), `1`)
	assertToS("*uint", ToS(&u), `2`)
	assertToS("*int8", ToS(&i8), `3`)
	assertToS("*int16", ToS(&i16), `4`)
	assertToS("*int32", ToS(&i32), `5`)
	assertToS("*int64", ToS(&i64), `6`)
	assertToS("*uint8", ToS(&u8), `7`)
	assertToS("*uint16", ToS(&u16), `8`)
	assertToS("*uint32", ToS(&u32), `9`)
	assertToS("*uint64", ToS(&u64), `10`)
	assertToS("*float32", ToS(&f32), `11.5`)
	assertToS("*float64", ToS(&f64), `12.5`)
	assertToS("*any", ToS(&anyWrap), `13`)
	assertToS("bool", ToS(false), `false`)
	// nil typed pointer that implements fmt.Stringer still hits Stringer branch.
	assertToS("nil stringer pointer", ToS(nilPS), `ptr`)
	assertToS("default json", ToS(map[string]any{`x`: 1}), `{"x":1}`)

	now := time.Now().UTC().Truncate(time.Second)
	if got := ToTime(now); !got.Equal(now) {
		t.Fatalf("ToTime(time.Time) mismatch: %v", got)
	}
	if got := ToTime(&now); !got.Equal(now) {
		t.Fatalf("ToTime(*time.Time) mismatch: %v", got)
	}
	rawBytes := []byte(`2024-02-29 00:00:00`)
	if got := ToTime(&rawBytes); !got.Equal(time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("ToTime(*[]byte) mismatch: %v", got)
	}
	var wrapTime any = `2024-02-29 00:00:00`
	if got := ToTime(&wrapTime); got.IsZero() {
		t.Fatalf("ToTime(*any) should parse time")
	}
	if got := ToTime(struct{}{}); !got.IsZero() {
		t.Fatalf("ToTime(default) should return zero time, got %v", got)
	}

	// pointer branches currently use `== 0` semantics.
	pi := 0
	pi2 := 2
	if !ToBool(&pi) || ToBool(&pi2) {
		t.Fatalf("ToBool pointer-int branch mismatch")
	}
	if ToBool(boolStr(`false`)) || !ToBool(boolStr(`yes`)) {
		t.Fatalf("ToBool stringer branch mismatch")
	}
	if ToBool(`no`) || !ToBool(`ok`) {
		t.Fatalf("ToBool string branch mismatch")
	}
	if ToBool(struct{}{}) {
		t.Fatalf("ToBool default branch should be false")
	}
}

func TestJson5AndContainerConverters(t *testing.T) {
	if got := json5fromMIX(map[int64]any{1: `x`}); got != `{1:"x"}` {
		t.Fatalf("json5fromMIX mismatch: %q", got)
	}
	if got := json5fromMIAX(map[int64][]any{2: []any{1, `a`}}); got != `{2:[1,"a"]}` {
		t.Fatalf("json5fromMIAX mismatch: %q", got)
	}
	if got := json5fromMSAX(map[string][]any{`k`: []any{1}}); got != `{"k":[1]}` {
		t.Fatalf("json5fromMSAX mismatch: %q", got)
	}
	if got := json5fromMSI(map[string]int64{`abc`: 1}); got != `{'abc':1}` {
		t.Fatalf("json5fromMSI mismatch: %q", got)
	}

	var b bytes.Buffer
	b.WriteString(`raw`)
	if ToJson5(b) != `raw` {
		t.Fatalf("ToJson5(bytes.Buffer) mismatch")
	}
	if ToJson5(`x`) != `"x"` || ToJson5([]byte(`x`)) != `"x"` || ToJson5(int(1)) != `1` ||
		ToJson5(int64(2)) != `2` || ToJson5(int32(3)) != `3` || ToJson5(uint(4)) != `4` ||
		ToJson5(uint64(5)) != `5` || ToJson5(uint32(6)) != `6` || ToJson5(float32(7.5)) != `7.5` ||
		ToJson5(float64(8.5)) != `8.5` || ToJson5(true) != `true` {
		t.Fatalf("ToJson5 scalar branches mismatch")
	}
	if ToJson5(M.IB{1: true}) != `{1:true}` || ToJson5(M.IX{1: `x`}) != `{1:"x"}` ||
		ToJson5(M.IAX{1: []any{`x`}}) != `{1:["x"]}` || ToJson5(M.SAX{`k`: []any{1}}) != `{"k":[1]}` ||
		ToJson5(M.SX{`k`: 1}) != `{"k":1}` || ToJson5(map[string]any{`k`: 1}) != `{"k":1}` ||
		ToJson5(M.SI{`k`: 1}) != `{'k':1}` || ToJson5(A.X{1, `x`}) != `[1,"x"]` || ToJson5([]any{1, `x`}) != `[1,"x"]` {
		t.Fatalf("ToJson5 map/array branches mismatch")
	}
	if got := ToJson5(struct{ A int }{A: 1}); got != `{"A":1}` {
		t.Fatalf("ToJson5 default marshal mismatch: %q", got)
	}

	if got := ToAX([]any{1, 2}); !reflect.DeepEqual(got, A.X{1, 2}) {
		t.Fatalf("ToAX([]any) mismatch: %#v", got)
	}
	if got := ToAX(`bad`); len(got) != 0 {
		t.Fatalf("ToAX invalid should be empty: %#v", got)
	}
	if got := ToAF([]float64{1, 2}); !reflect.DeepEqual(got, []float64{1, 2}) {
		t.Fatalf("ToAF([]float64) mismatch: %#v", got)
	}
	if got := ToAF(`bad`); len(got) != 0 {
		t.Fatalf("ToAF invalid should be empty: %#v", got)
	}
	if got := ToMSX(M.SX{`a`: 1}); got[`a`] != 1 {
		t.Fatalf("ToMSX(M.SX) mismatch: %#v", got)
	}
	if got := ToMSX(`bad`); len(got) != 0 {
		t.Fatalf("ToMSX invalid should be empty: %#v", got)
	}
	if got := ToMSS(M.SS{`a`: `b`}); got[`a`] != `b` {
		t.Fatalf("ToMSS(M.SS) mismatch: %#v", got)
	}
	if got := ToMSS(`bad`); len(got) != 0 {
		t.Fatalf("ToMSS invalid should be empty: %#v", got)
	}
}

func TestToBoolPointerAndNilBranches(t *testing.T) {
	i, iz := 1, 0
	u, uz := uint(1), uint(0)
	i8, i8z := int8(1), int8(0)
	i16, i16z := int16(1), int16(0)
	i32, i32z := int32(1), int32(0)
	i64, i64z := int64(1), int64(0)
	u8, u8z := uint8(1), uint8(0)
	u16, u16z := uint16(1), uint16(0)
	u32, u32z := uint32(1), uint32(0)
	u64, u64z := uint64(1), uint64(0)
	f32, f32z := float32(1), float32(0)
	f64, f64z := float64(1), float64(0)

	// pointer branches intentionally use `== 0` semantics in current implementation.
	if ToBool(&i) || !ToBool(&iz) || ToBool(&u) || !ToBool(&uz) || ToBool(&i8) || !ToBool(&i8z) ||
		ToBool(&i16) || !ToBool(&i16z) || ToBool(&i32) || !ToBool(&i32z) || ToBool(&i64) || !ToBool(&i64z) ||
		ToBool(&u8) || !ToBool(&u8z) || ToBool(&u16) || !ToBool(&u16z) || ToBool(&u32) || !ToBool(&u32z) ||
		ToBool(&u64) || !ToBool(&u64z) || ToBool(&f32) || !ToBool(&f32z) || ToBool(&f64) || !ToBool(&f64z) {
		t.Fatalf("ToBool pointer numeric branches mismatch")
	}
}

func TestToTimeAndContainerNilBranches(t *testing.T) {
	var nilBytes *[]byte
	var nilStr *string
	var nilAny *any
	if got := ToTime(nilBytes); !got.IsZero() {
		t.Fatalf("ToTime(nil *[]byte) should be zero: %v", got)
	}
	if got := ToTime(nilStr); !got.IsZero() {
		t.Fatalf("ToTime(nil *string) should be zero: %v", got)
	}
	if got := ToTime(nilAny); !got.IsZero() {
		t.Fatalf("ToTime(nil *any) should be zero: %v", got)
	}

	if got := ToAX(nil); len(got) != 0 {
		t.Fatalf("ToAX(nil) should be empty: %#v", got)
	}
	if got := ToAF(nil); len(got) != 0 {
		t.Fatalf("ToAF(nil) should be empty: %#v", got)
	}
	if got := ToMSX(nil); len(got) != 0 {
		t.Fatalf("ToMSX(nil) should be empty: %#v", got)
	}
	if got := ToMSS(nil); len(got) != 0 {
		t.Fatalf("ToMSS(nil) should be empty: %#v", got)
	}
}

func TestConvertersAdditionalScalarAndDefaultBranches(t *testing.T) {
	var nilAny *any

	for _, tc := range []struct {
		name string
		in   any
		want uint64
	}{
		{name: "int", in: int(1), want: 1},
		{name: "uint", in: uint(2), want: 2},
		{name: "int8", in: int8(3), want: 3},
		{name: "int16", in: int16(4), want: 4},
		{name: "int32", in: int32(5), want: 5},
		{name: "int64", in: int64(6), want: 6},
		{name: "uint8", in: uint8(7), want: 7},
		{name: "uint16", in: uint16(8), want: 8},
		{name: "uint32", in: uint32(9), want: 9},
		{name: "float32", in: float32(10.9), want: 10},
		{name: "float64", in: float64(11.9), want: 11},
		{name: "bool true", in: true, want: 1},
		{name: "string int", in: `12`, want: 12},
		{name: "string float", in: `13.9`, want: 13},
		{name: "nil *any", in: nilAny, want: 0},
		{name: "default", in: struct{}{}, want: 0},
	} {
		if got := ToU(tc.in); got != tc.want {
			t.Fatalf("ToU %s mismatch: want=%d got=%d", tc.name, tc.want, got)
		}
	}

	for _, tc := range []struct {
		name string
		in   any
		want byte
	}{
		{name: "int", in: int(1), want: 1},
		{name: "uint", in: uint(2), want: 2},
		{name: "int8", in: int8(3), want: 3},
		{name: "int16", in: int16(4), want: 4},
		{name: "int32", in: int32(5), want: 5},
		{name: "int64", in: int64(6), want: 6},
		{name: "uint16", in: uint16(8), want: 8},
		{name: "uint32", in: uint32(9), want: 9},
		{name: "uint64", in: uint64(10), want: 10},
		{name: "float32", in: float32(11.9), want: 11},
		{name: "float64", in: float64(12.9), want: 12},
		{name: "bool true", in: true, want: 1},
		{name: "string int", in: `13`, want: 13},
		{name: "string float", in: `14.7`, want: 14},
		{name: "nil *any", in: nilAny, want: 0},
		{name: "default", in: struct{}{}, want: 0},
	} {
		if got := ToByte(tc.in); got != tc.want {
			t.Fatalf("ToByte %s mismatch: want=%d got=%d", tc.name, tc.want, got)
		}
	}

	for _, tc := range []struct {
		name string
		in   any
		want int64
	}{
		{name: "int64 fast path", in: int64(99), want: 99},
		{name: "int", in: int(1), want: 1},
		{name: "uint", in: uint(2), want: 2},
		{name: "int8", in: int8(3), want: 3},
		{name: "int16", in: int16(4), want: 4},
		{name: "int32", in: int32(5), want: 5},
		{name: "uint8", in: uint8(6), want: 6},
		{name: "uint16", in: uint16(7), want: 7},
		{name: "uint32", in: uint32(8), want: 8},
		{name: "uint64", in: uint64(9), want: 9},
		{name: "float32", in: float32(10.9), want: 10},
		{name: "float64", in: float64(11.9), want: 11},
		{name: "string int", in: `12`, want: 12},
		{name: "bool false", in: false, want: 0},
		{name: "nil *any", in: nilAny, want: 0},
		{name: "default", in: struct{}{}, want: 0},
	} {
		if got := ToI(tc.in); got != tc.want {
			t.Fatalf("ToI %s mismatch: want=%d got=%d", tc.name, tc.want, got)
		}
	}

	for _, tc := range []struct {
		name string
		in   any
		want float64
	}{
		{name: "float64 fast path", in: float64(99.5), want: 99.5},
		{name: "int", in: int(1), want: 1},
		{name: "uint", in: uint(2), want: 2},
		{name: "int8", in: int8(3), want: 3},
		{name: "int16", in: int16(4), want: 4},
		{name: "int32", in: int32(5), want: 5},
		{name: "int64", in: int64(6), want: 6},
		{name: "uint8", in: uint8(7), want: 7},
		{name: "uint16", in: uint16(8), want: 8},
		{name: "uint32", in: uint32(9), want: 9},
		{name: "uint64", in: uint64(10), want: 10},
		{name: "string", in: `12.5`, want: 12.5},
		{name: "bool true", in: true, want: 1},
		{name: "nil *any", in: nilAny, want: 0},
		{name: "default", in: struct{}{}, want: 0},
	} {
		if got := ToF(tc.in); got != tc.want {
			t.Fatalf("ToF %s mismatch: want=%v got=%v", tc.name, tc.want, got)
		}
	}
}

func TestToSBoolToArrAndJson5AdditionalBranches(t *testing.T) {
	var nilAny *any

	if ToS(`plain`) != `plain` || ToS(uint(2)) != `2` || ToS(uint64(3)) != `3` || ToS(float64(4.5)) != `4.5` {
		t.Fatalf("ToS scalar branches mismatch")
	}
	if ToS(nilAny) != `` {
		t.Fatalf("ToS(nil *any) should be empty")
	}

	if !ToBool(true) || ToBool(false) || !ToBool(int8(1)) || ToBool(int16(0)) || !ToBool(uint8(1)) ||
		ToBool(uint16(0)) || !ToBool(float32(1)) || ToBool(float64(0)) {
		t.Fatalf("ToBool scalar branches mismatch")
	}
	for _, v := range []string{``, `0`, `f`, `false`, `no`, `n`, ` N `} {
		if ToBool(v) {
			t.Fatalf("ToBool false-string mismatch for %q", v)
		}
	}
	for _, v := range []boolStr{``, `0`, `f`, `false`, `no`, `n`, ` yes `} {
		got := ToBool(v)
		want := strings.TrimSpace(strings.ToLower(v.String()))
		expect := !(want == `` || want == `0` || want == `f` || want == `false` || want == `no` || want == `n`)
		if got != expect {
			t.Fatalf("ToBool stringer mismatch for %q: want=%v got=%v", v, expect, got)
		}
	}

	if got := ToArr(nil); len(got) != 0 {
		t.Fatalf("ToArr(nil) should be empty: %#v", got)
	}

	assertMapLike := func(name, got string, contains ...string) {
		t.Helper()
		if !strings.HasPrefix(got, "{") || !strings.HasSuffix(got, "}") {
			t.Fatalf("%s malformed map json5: %q", name, got)
		}
		if len(contains) > 1 && !strings.Contains(got, ",") {
			t.Fatalf("%s should contain comma for multi-entry map: %q", name, got)
		}
		for _, must := range contains {
			if !strings.Contains(got, must) {
				t.Fatalf("%s missing %q in %q", name, must, got)
			}
		}
	}

	assertMapLike("json5fromMIB", json5fromMIB(map[int64]bool{1: true, 2: false}), `1:true`, `2:false`)
	assertMapLike("json5fromMIX", json5fromMIX(map[int64]any{1: `x`, 2: 3}), `1:"x"`, `2:3`)
	assertMapLike("json5fromMIAX", json5fromMIAX(map[int64][]any{1: {`a`}, 2: {1}}), `1:["a"]`, `2:[1]`)
	assertMapLike("json5fromMSAX", json5fromMSAX(map[string][]any{`a`: {1}, `b`: {`x`}}), `"a":[1]`, `"b":["x"]`)
	assertMapLike("json5fromMSI", json5fromMSI(map[string]int64{`1`: 1, `abc`: 2, `a-b`: 3, ``: 4}), `'1':1`, `'abc':2`, `'a-b':3`, `'':4`)

	assertMapLike("ToJson5 map[int64]any", ToJson5(map[int64]any{1: `x`, 2: true}), `1:"x"`, `2:true`)
	assertMapLike("ToJson5 map[int64][]any", ToJson5(map[int64][]any{1: {`x`}, 2: {1}}), `1:["x"]`, `2:[1]`)
	assertMapLike("ToJson5 map[string][]any", ToJson5(map[string][]any{`a`: {1}, `b`: {`x`}}), `"a":[1]`, `"b":["x"]`)
	assertMapLike("ToJson5 map[string]int64", ToJson5(map[string]int64{`1`: 1, `abc`: 2}), `'1':1`, `'abc':2`)
}
