package M

import (
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v5"
)

func TestSX_FromJson(t *testing.T) {
	t.Run(`nil map`, func(t *testing.T) {
		var m SX = nil
		assert.True(t, m.FromJson(`{"a":123,"b":"abc"}`))
		if m.GetInt(`a`) != 123 || m[`b`] != `abc` {
			t.Error(`invalid value`)
		}
	})

	t.Run(`default case`, func(t *testing.T) {
		m := SX{}
		assert.True(t, m.FromJson(`{"a":123,"b":"abc"}`))
		if m.GetInt(`a`) != 123 || m[`b`] != `abc` {
			t.Error(`invalid value`)
		}
	})

	t.Run(`invalid json`, func(t *testing.T) {
		m := SX{}
		assert.False(t, m.FromJson(`{"a":123,"b":"abc"`))
		if len(m) != 2 { // goccy/go-json will still parse the valid json part
			t.Error(`invalid value`)
		}
	})

	t.Run(`overwrites`, func(t *testing.T) {
		m := SX{}
		m["a"] = 234
		assert.True(t, m.FromJson(`{"a":123,"b":"abc"}`))
		if m.GetInt(`a`) != 123 || m[`b`] != `abc` {
			t.Error(`invalid value`)
		}
	})

	t.Run(`not overwrite`, func(t *testing.T) {
		m := SX{}
		m["a"] = 234
		assert.True(t, m.FromJson(`{"b":"abc"}`))
		if m.GetInt(`a`) != 234 || m[`b`] != `abc` {
			t.Error(`invalid value`)
		}
	})

	t.Run(`empty string`, func(t *testing.T) {
		var m SX
		assert.False(t, m.FromJson(``))
		if len(m) != 0 {
			t.Error(`invalid value`)
		}
	})

	t.Run(`inside struct`, func(t *testing.T) {
		x := struct {
			Foo SX
		}{}
		assert.True(t, x.Foo.FromJson(`{"a":123,"b":"abc"}`))
		if x.Foo.GetInt(`a`) != 123 || x.Foo[`b`] != `abc` {
			t.Error(`invalid value`)
		}
	})
}

type mapStringer struct{}

func (mapStringer) String() string { return "stringer-value" }

func toSet(vals []string) map[string]bool {
	res := map[string]bool{}
	for _, v := range vals {
		res[v] = true
	}
	return res
}

func testKeysConcatSet(t *testing.T, got string, sep string, want map[string]bool) {
	t.Helper()
	if got == "" && len(want) == 0 {
		return
	}
	parts := strings.Split(got, sep)
	if !reflect.DeepEqual(toSet(parts), want) {
		t.Fatalf("keys mismatch: got=%q wantSet=%#v", got, want)
	}
}

func TestSSAndSBUtilities(t *testing.T) {
	ss := SS{"b": "2", "a": "1"}
	if !reflect.DeepEqual(ss.SortedKeys(), []string{"a", "b"}) {
		t.Fatalf("SS.SortedKeys mismatch: %#v", ss.SortedKeys())
	}
	testKeysConcatSet(t, ss.KeysConcat(","), ",", map[string]bool{"a": true, "b": true})
	if got := ss.Pretty("|"); got != "a 1b 2|" {
		t.Fatalf("SS.Pretty mismatch: %q", got)
	}
	if got := ss.PrettyFunc("|", func(k, v string) string { return "[" + v + "]" }); got != "a [1]|b [2]|" {
		t.Fatalf("SS.PrettyFunc mismatch: %q", got)
	}

	scy := SS{"a'b": "c'd"}.ToScylla()
	if !strings.Contains(scy, "'a&apos;b':'c&apos;d'") {
		t.Fatalf("SS.ToScylla mismatch: %q", scy)
	}

	jsonSS := ss.ToJson()
	ss2 := SS{}
	if err := json.Unmarshal([]byte(jsonSS), &ss2); err != nil || !reflect.DeepEqual(ss, ss2) {
		t.Fatalf("SS.ToJson mismatch: err=%v got=%#v", err, ss2)
	}
	ss3 := SS{}
	if err := msgpack.Unmarshal(ss.ToMsgp(), &ss3); err != nil || !reflect.DeepEqual(ss, ss3) {
		t.Fatalf("SS.ToMsgp mismatch: err=%v got=%#v", err, ss3)
	}

	sb := SB{"x": true, "y": false}
	if !reflect.DeepEqual(sb.SortedKeys(), []string{"x", "y"}) {
		t.Fatalf("SB.SortedKeys mismatch: %#v", sb.SortedKeys())
	}
	testKeysConcatSet(t, sb.KeysConcat("|"), "|", map[string]bool{"x": true, "y": true})

	sb2 := SB{}
	if err := json.Unmarshal([]byte(sb.ToJson()), &sb2); err != nil || !reflect.DeepEqual(sb, sb2) {
		t.Fatalf("SB.ToJson mismatch: err=%v got=%#v", err, sb2)
	}
	sb3 := SB{}
	if err := msgpack.Unmarshal(sb.ToMsgp(), &sb3); err != nil || !reflect.DeepEqual(sb, sb3) {
		t.Fatalf("SB.ToMsgp mismatch: err=%v got=%#v", err, sb3)
	}
	if pretty := sb.ToJsonPretty(); !strings.Contains(pretty, "\n") {
		t.Fatalf("SB.ToJsonPretty should be multiline: %q", pretty)
	}
	if _, ok := sb.IntoJson(); !ok {
		t.Fatalf("SB.IntoJson should succeed")
	}
	if _, ok := sb.IntoMsgp(); !ok {
		t.Fatalf("SB.IntoMsgp should succeed")
	}
	if _, ok := sb.IntoJsonPretty(); !ok {
		t.Fatalf("SB.IntoJsonPretty should succeed")
	}
}

func TestSXSerializationAndGeneralHelpers(t *testing.T) {
	sx := SX{"b": 2, "a": "1"}
	if !reflect.DeepEqual(sx.SortedKeys(), []string{"a", "b"}) {
		t.Fatalf("SX.SortedKeys mismatch: %#v", sx.SortedKeys())
	}
	jsonSX := sx.ToJson()
	sx2 := SX{}
	if err := json.Unmarshal([]byte(jsonSX), &sx2); err != nil {
		t.Fatalf("SX.ToJson should be valid JSON: %v", err)
	}
	if len(sx.ToMsgp()) == 0 {
		t.Fatalf("SX.ToMsgp should not be empty")
	}
	if pretty := sx.ToJsonPretty(); !strings.Contains(pretty, "\n") {
		t.Fatalf("SX.ToJsonPretty should be multiline: %q", pretty)
	}
	if _, ok := sx.IntoJson(); !ok {
		t.Fatalf("SX.IntoJson should succeed")
	}
	if _, ok := sx.IntoMsgp(); !ok {
		t.Fatalf("SX.IntoMsgp should succeed")
	}
	if _, ok := sx.IntoJsonPretty(); !ok {
		t.Fatalf("SX.IntoJsonPretty should succeed")
	}

	sx.Set("c", 3)
	keys := sx.Keys()
	sort.Strings(keys)
	if !reflect.DeepEqual(keys, []string{"a", "b", "c"}) {
		t.Fatalf("SX.Keys mismatch: %#v", keys)
	}
	if got := sx.Pretty("|"); got != "a 1|b 2|c 3|" {
		t.Fatalf("SX.Pretty mismatch: %q", got)
	}

	raw := map[string]any{"x": 1}
	if got := ToJson(raw); got != `{"x":1}` {
		t.Fatalf("M.ToJson mismatch: %q", got)
	}
	if len(ToMsgp(raw)) == 0 {
		t.Fatalf("M.ToMsgp should not be empty")
	}
}

func TestSXPrimitiveGetters(t *testing.T) {
	s := "hello"
	sx := SX{
		"i":    int8(7),
		"u":    "8",
		"f":    "9.5",
		"s1":   uint16(10),
		"s2":   false,
		"s3":   &s,
		"s4":   mapStringer{},
		"b1":   " false ",
		"b2":   "yes",
		"b3":   0,
		"b4":   float64(1),
		"none": nil,
	}
	if got := sx.GetInt("i"); got != 7 {
		t.Fatalf("SX.GetInt mismatch: %d", got)
	}
	if got := sx.GetUint("u"); got != 8 {
		t.Fatalf("SX.GetUint mismatch: %d", got)
	}
	if got := sx.GetFloat("f"); got != 9.5 {
		t.Fatalf("SX.GetFloat mismatch: %v", got)
	}
	if got := sx.GetStr("s1"); got != "10" {
		t.Fatalf("SX.GetStr uint mismatch: %q", got)
	}
	if got := sx.GetStr("s2"); got != "false" {
		t.Fatalf("SX.GetStr bool mismatch: %q", got)
	}
	if got := sx.GetStr("s3"); got != "hello" {
		t.Fatalf("SX.GetStr *string mismatch: %q", got)
	}
	if got := sx.GetStr("s4"); got != "stringer-value" {
		t.Fatalf("SX.GetStr Stringer mismatch: %q", got)
	}
	if got := sx.GetInt("none"); got != 0 {
		t.Fatalf("SX.GetInt nil mismatch: %d", got)
	}
	if got := sx.GetBool("b1"); got {
		t.Fatalf("SX.GetBool false-string mismatch: %v", got)
	}
	if got := sx.GetBool("b2"); !got {
		t.Fatalf("SX.GetBool true-string mismatch: %v", got)
	}
	if got := sx.GetBool("b3"); got {
		t.Fatalf("SX.GetBool zero mismatch: %v", got)
	}
	if got := sx.GetBool("b4"); !got {
		t.Fatalf("SX.GetBool non-zero mismatch: %v", got)
	}
}

func TestSXNestedAndArrayGetters(t *testing.T) {
	sx := SX{
		"msb1": map[string]bool{"a": true},
		"msb2": SB{"b": false},
		"msb3": map[string]any{"x": true, "y": "no"},
		"msf":  map[string]any{"a": "1.5", "b": float64(2.3)},
		"msi":  map[string]any{"a": "2", "b": int64(3)},
		"mib":  map[int64]any{1: true, 2: "ignored"},
		"msx1": map[string]any{"k": "v"},
		"msx2": SX{"k2": 2},
		"ax":   []any{1, "2"},
		"i1":   []int64{1, 2},
		"i2":   []float64{3.9, 4.1},
		"i3":   []any{int8(5), "6", "7.8", true},
	}
	if got := sx.GetMSB("msb1"); !reflect.DeepEqual(got, SB{"a": true}) {
		t.Fatalf("SX.GetMSB map[string]bool mismatch: %#v", got)
	}
	if got := sx.GetMSB("msb2"); !reflect.DeepEqual(got, SB{"b": false}) {
		t.Fatalf("SX.GetMSB SB mismatch: %#v", got)
	}
	if got := sx.GetMSB("msb3"); !reflect.DeepEqual(got, SB{"x": true}) {
		t.Fatalf("SX.GetMSB map[string]any mismatch: %#v", got)
	}
	if got := sx.GetMSF("msf"); !reflect.DeepEqual(got, SF{"a": 1.5, "b": 2.3}) {
		t.Fatalf("SX.GetMSF mismatch: %#v", got)
	}
	if got := sx.GetMSI("msi"); !reflect.DeepEqual(got, SI{"a": 2, "b": 3}) {
		t.Fatalf("SX.GetMSI mismatch: %#v", got)
	}
	if got := sx.GetMIB("mib"); !reflect.DeepEqual(got, IB{1: true}) {
		t.Fatalf("SX.GetMIB mismatch: %#v", got)
	}
	if got := sx.GetMSX("msx1"); !reflect.DeepEqual(got, SX{"k": "v"}) {
		t.Fatalf("SX.GetMSX map mismatch: %#v", got)
	}
	if got := sx.GetMSX("msx2"); !reflect.DeepEqual(got, SX{"k2": 2}) {
		t.Fatalf("SX.GetMSX SX mismatch: %#v", got)
	}
	if got := sx.GetAX("ax"); !reflect.DeepEqual(got, []any{1, "2"}) {
		t.Fatalf("SX.GetAX mismatch: %#v", got)
	}
	if got := sx.GetIntArr("i1"); !reflect.DeepEqual(got, []int64{1, 2}) {
		t.Fatalf("SX.GetIntArr []int64 mismatch: %#v", got)
	}
	if got := sx.GetIntArr("i2"); !reflect.DeepEqual(got, []int64{3, 4}) {
		t.Fatalf("SX.GetIntArr []float64 mismatch: %#v", got)
	}
	if got := sx.GetIntArr("i3"); !reflect.DeepEqual(got, []int64{5, 6, 6, 7}) {
		t.Fatalf("SX.GetIntArr []any mismatch: %#v", got)
	}
	if got := sx.GetIntArr("missing"); len(got) != 0 {
		t.Fatalf("SX.GetIntArr missing should be empty: %#v", got)
	}
}

func TestOtherMapHelpers(t *testing.T) {
	ss := SS{"a": "1", "b": "2.5"}
	if ss.GetInt("a") != 1 || ss.GetUint("a") != 1 || ss.GetFloat("b") != 2.5 || ss.GetStr("a") != "1" {
		t.Fatalf("SS getters mismatch")
	}
	keys := ss.Keys()
	sort.Strings(keys)
	if !reflect.DeepEqual(keys, []string{"a", "b"}) {
		t.Fatalf("SS.Keys mismatch: %#v", keys)
	}
	ss.Merge(SS{"c": "3"})
	if ss["c"] != "3" {
		t.Fatalf("SS.Merge mismatch: %#v", ss)
	}

	ix := IX{1: "a", 2: "b"}
	ixKeys := ix.Keys()
	sort.Slice(ixKeys, func(i, j int) bool { return ixKeys[i] < ixKeys[j] })
	if !reflect.DeepEqual(ixKeys, []int64{1, 2}) {
		t.Fatalf("IX.Keys mismatch: %#v", ixKeys)
	}
	if got := ix.ToSX(); !reflect.DeepEqual(got, SX{"1": "a", "2": "b"}) {
		t.Fatalf("IX.ToSX mismatch: %#v", got)
	}

	ii := II{3: 9, 1: 1}
	iiKeys := ii.Keys()
	sort.Slice(iiKeys, func(i, j int) bool { return iiKeys[i] < iiKeys[j] })
	if !reflect.DeepEqual(iiKeys, []int64{1, 3}) {
		t.Fatalf("II.Keys mismatch: %#v", iiKeys)
	}
	testKeysConcatSet(t, ii.KeysConcat(","), ",", map[string]bool{"1": true, "3": true})

	ib := IB{4: true, 2: false}
	ibKeys := ib.Keys()
	sort.Slice(ibKeys, func(i, j int) bool { return ibKeys[i] < ibKeys[j] })
	if !reflect.DeepEqual(ibKeys, []int64{2, 4}) {
		t.Fatalf("IB.Keys mismatch: %#v", ibKeys)
	}
	testKeysConcatSet(t, ib.KeysConcat(","), ",", map[string]bool{"2": true, "4": true})

	started := SSKeysStartedWith(SS{"ab": "1", "ac": "2", "zz": "3"}, "a")
	sort.Strings(started)
	if !reflect.DeepEqual(started, []string{"ab", "ac"}) {
		t.Fatalf("SSKeysStartedWith mismatch: %#v", started)
	}
}

func TestSXFromJsonNilReceiver(t *testing.T) {
	var p *SX
	if p.FromJson(`{"a":1}`) {
		t.Fatalf("nil receiver should return false")
	}
}

func TestSXPrimitiveGettersMoreBranches(t *testing.T) {
	duration := time.Duration(11)
	nilStr := (*string)(nil)
	sx := SX{
		"i_int":     int(1),
		"i_i16":     int16(2),
		"i_i32":     int32(3),
		"i_uint":    uint(4),
		"i_u8":      uint8(5),
		"i_u16":     uint16(6),
		"i_u32":     uint32(7),
		"i_u64":     uint64(8),
		"i_f32":     float32(9.7),
		"i_f64":     float64(10.8),
		"i_dur":     duration,
		"i_true":    true,
		"i_false":   false,
		"i_s_int":   "123",
		"i_s_float": "4.5",
		"i_s_bad":   "x",
		"i_bad":     struct{}{},

		"u_int64":   int64(12),
		"u_int":     int(13),
		"u_i8":      int8(14),
		"u_i16":     int16(15),
		"u_i32":     int32(16),
		"u_uint":    uint(17),
		"u_u8":      uint8(18),
		"u_u16":     uint16(19),
		"u_u32":     uint32(20),
		"u_u64":     uint64(21),
		"u_f32":     float32(22.2),
		"u_f64":     float64(23.3),
		"u_dur":     time.Duration(24),
		"u_true":    true,
		"u_false":   false,
		"u_s_int":   "13",
		"u_s_float": "14.6",
		"u_s_bad":   "x",
		"u_bad":     struct{}{},

		"f_int":   int(25),
		"f_i8":    int8(26),
		"f_i16":   int16(27),
		"f_i32":   int32(28),
		"f_int64": int64(15),
		"f_uint":  uint(16),
		"f_u8":    uint8(17),
		"f_u16":   uint16(18),
		"f_u32":   uint32(19),
		"f_u64":   uint64(20),
		"f_f32":   float32(21.5),
		"f_dur":   time.Duration(22),
		"f_true":  true,
		"f_false": false,
		"f_s":     "23.25",
		"f_s_bad": "xx",
		"f_bad":   struct{}{},

		"s_i8":       int8(24),
		"s_i16":      int16(25),
		"s_i32":      int32(26),
		"s_i64":      int64(27),
		"s_u":        uint(28),
		"s_u8":       uint8(29),
		"s_u16":      uint16(30),
		"s_u32":      uint32(31),
		"s_u64":      uint64(32),
		"s_f32":      float32(33.5),
		"s_f64":      float64(34.5),
		"s_true":     true,
		"s_false":    false,
		"s_stringer": mapStringer{},
		"s_nilptr":   nilStr,
		"s_bad":      struct{}{},

		"b_int":      int(1),
		"b_int_zero": int(0),
		"b_i8":       int8(1),
		"b_i16":      int16(1),
		"b_i32":      int32(1),
		"b_i64":      int64(1),
		"b_u":        uint(1),
		"b_u8":       uint8(1),
		"b_u16":      uint16(1),
		"b_u32":      uint32(1),
		"b_u64":      uint64(1),
		"b_f32":      float32(1),
		"b_f64":      float64(1),
		"b_s_false":  " false ",
		"b_s_f":      "f",
		"b_s_zero":   "0",
		"b_s_empty":  "",
		"b_s_true":   "yes",
		"b_str_obj":  codegenStringer{},
		"b_bad":      struct{}{},
	}

	if sx.GetInt("i_int") != 1 || sx.GetInt("i_i16") != 2 || sx.GetInt("i_i32") != 3 ||
		sx.GetInt("i_uint") != 4 || sx.GetInt("i_u8") != 5 || sx.GetInt("i_u16") != 6 ||
		sx.GetInt("i_u32") != 7 || sx.GetInt("i_u64") != 8 || sx.GetInt("i_f32") != 9 ||
		sx.GetInt("i_f64") != 10 || sx.GetInt("i_dur") != 11 || sx.GetInt("i_true") != 1 ||
		sx.GetInt("i_false") != 0 || sx.GetInt("i_s_int") != 123 || sx.GetInt("i_s_float") != 4 ||
		sx.GetInt("i_s_bad") != 0 || sx.GetInt("i_bad") != 0 {
		t.Fatalf("GetInt branch mismatch")
	}

	if sx.GetUint("u_int64") != 12 || sx.GetUint("u_int") != 13 || sx.GetUint("u_i8") != 14 ||
		sx.GetUint("u_i16") != 15 || sx.GetUint("u_i32") != 16 || sx.GetUint("u_uint") != 17 ||
		sx.GetUint("u_u8") != 18 || sx.GetUint("u_u16") != 19 || sx.GetUint("u_u32") != 20 ||
		sx.GetUint("u_u64") != 21 || sx.GetUint("u_f32") != 22 || sx.GetUint("u_f64") != 23 ||
		sx.GetUint("u_dur") != 24 || sx.GetUint("u_true") != 1 || sx.GetUint("u_false") != 0 ||
		sx.GetUint("u_s_int") != 13 || sx.GetUint("u_s_float") != 14 || sx.GetUint("u_s_bad") != 0 ||
		sx.GetUint("u_bad") != 0 {
		t.Fatalf("GetUint branch mismatch")
	}

	if sx.GetFloat("f_int") != 25 || sx.GetFloat("f_i8") != 26 || sx.GetFloat("f_i16") != 27 ||
		sx.GetFloat("f_i32") != 28 || sx.GetFloat("f_int64") != 15 || sx.GetFloat("f_uint") != 16 ||
		sx.GetFloat("f_u8") != 17 || sx.GetFloat("f_u16") != 18 || sx.GetFloat("f_u32") != 19 ||
		sx.GetFloat("f_u64") != 20 || sx.GetFloat("f_f32") != 21.5 || sx.GetFloat("f_dur") != 22 ||
		sx.GetFloat("f_true") != 1 || sx.GetFloat("f_false") != 0 || sx.GetFloat("f_s") != 23.25 ||
		sx.GetFloat("f_s_bad") != 0 || sx.GetFloat("f_bad") != 0 {
		t.Fatalf("GetFloat branch mismatch")
	}

	if sx.GetStr("s_i8") != "24" || sx.GetStr("s_i16") != "25" || sx.GetStr("s_i32") != "26" ||
		sx.GetStr("s_i64") != "27" || sx.GetStr("s_u") != "28" || sx.GetStr("s_u8") != "29" ||
		sx.GetStr("s_u16") != "30" || sx.GetStr("s_u32") != "31" || sx.GetStr("s_u64") != "32" ||
		sx.GetStr("s_f32") != "33.5" || sx.GetStr("s_f64") != "34.5" || sx.GetStr("s_true") != "true" ||
		sx.GetStr("s_false") != "false" || sx.GetStr("s_stringer") != "stringer-value" ||
		sx.GetStr("s_nilptr") != "" || sx.GetStr("s_bad") != "" {
		t.Fatalf("GetStr branch mismatch")
	}

	if !sx.GetBool("b_int") || sx.GetBool("b_int_zero") || !sx.GetBool("b_i8") || !sx.GetBool("b_i16") ||
		!sx.GetBool("b_i32") || !sx.GetBool("b_i64") || !sx.GetBool("b_u") || !sx.GetBool("b_u8") ||
		!sx.GetBool("b_u16") || !sx.GetBool("b_u32") || !sx.GetBool("b_u64") || !sx.GetBool("b_f32") ||
		!sx.GetBool("b_f64") || sx.GetBool("b_s_false") || sx.GetBool("b_s_f") || sx.GetBool("b_s_zero") ||
		sx.GetBool("b_s_empty") || !sx.GetBool("b_s_true") || !sx.GetBool("b_str_obj") || sx.GetBool("b_bad") {
		t.Fatalf("GetBool branch mismatch")
	}
}

type codegenStringer struct{}

func (codegenStringer) String() string { return "non-empty" }

func TestSXNestedAndArrayGettersMoreBranches(t *testing.T) {
	sx := SX{
		"msb_invalid": 123,
		"msf_alias":   SF{"a": 1.2},
		"msf_invalid": 1,
		"msi_alias":   SI{"a": 2},
		"msi_invalid": 1,
		"mib_alias":   IB{3: true},
		"mib_direct":  map[int64]bool{4: false},
		"mib_invalid": 1,
		"msx_invalid": 1,
		"ax_invalid":  1,
		"int_invalid": "nope",
		"int_many": []any{
			int(1), int8(2), int16(3), int32(4),
			uint(5), uint8(6), uint16(7), uint32(8), uint64(9),
			float32(10.2), float64(11.2),
			"12", "13.9", "xx",
		},
	}
	if got := sx.GetMSB("msb_invalid"); len(got) != 0 {
		t.Fatalf("GetMSB invalid should be empty: %#v", got)
	}
	if got := sx.GetMSF("msf_alias"); !reflect.DeepEqual(got, SF{"a": 1.2}) {
		t.Fatalf("GetMSF alias mismatch: %#v", got)
	}
	if got := sx.GetMSF("msf_invalid"); len(got) != 0 {
		t.Fatalf("GetMSF invalid should be empty: %#v", got)
	}
	if got := sx.GetMSI("msi_alias"); !reflect.DeepEqual(got, SI{"a": 2}) {
		t.Fatalf("GetMSI alias mismatch: %#v", got)
	}
	if got := sx.GetMSI("msi_invalid"); len(got) != 0 {
		t.Fatalf("GetMSI invalid should be empty: %#v", got)
	}
	if got := sx.GetMIB("mib_alias"); !reflect.DeepEqual(got, IB{3: true}) {
		t.Fatalf("GetMIB alias mismatch: %#v", got)
	}
	if got := sx.GetMIB("mib_direct"); !reflect.DeepEqual(got, IB{4: false}) {
		t.Fatalf("GetMIB direct mismatch: %#v", got)
	}
	if got := sx.GetMIB("mib_invalid"); len(got) != 0 {
		t.Fatalf("GetMIB invalid should be empty: %#v", got)
	}
	if got := sx.GetMSX("msx_invalid"); len(got) != 0 {
		t.Fatalf("GetMSX invalid should be empty: %#v", got)
	}
	if got := sx.GetAX("ax_invalid"); len(got) != 0 {
		t.Fatalf("GetAX invalid should be empty: %#v", got)
	}
	if got := sx.GetIntArr("int_invalid"); len(got) != 0 {
		t.Fatalf("GetIntArr invalid should be empty: %#v", got)
	}
	if got := sx.GetIntArr("int_many"); !reflect.DeepEqual(got, []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 12, 13}) {
		t.Fatalf("GetIntArr many types mismatch: %#v", got)
	}
}
