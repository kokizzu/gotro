package X

import (
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
