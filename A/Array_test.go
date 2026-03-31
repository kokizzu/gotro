package A

import (
	"reflect"
	"testing"
)

func TestToMsgp(t *testing.T) {
	m := []any{123, `abc`}
	if len(ToMsgp(m)) == 0 {
		t.Fatalf("ToMsgp should return bytes")
	}
}

func TestJoinAndContainsHelpers(t *testing.T) {
	if got := StrJoin([]string{`a`, `b`, `c`}, `-`); got != `a-b-c` {
		t.Fatalf("StrJoin mismatch: %q", got)
	}
	if got := IntJoin([]int64{1, 2, 3}, `|`); got != `1|2|3` {
		t.Fatalf("IntJoin mismatch: %q", got)
	}
	if got := UIntJoin([]uint64{4, 5}, `,`); got != `4,5` {
		t.Fatalf("UIntJoin mismatch: %q", got)
	}
	if !StrContains([]string{`foo`, `bar`}, `bar`) {
		t.Fatalf("StrContains should find existing value")
	}
	if IntContains([]int64{1, 2, 3}, 9) {
		t.Fatalf("IntContains should not find missing value")
	}
}

func TestAppendAndConvertHelpers(t *testing.T) {
	ints := StrToInt([]string{`1`, ``, `x`, `2`})
	if !reflect.DeepEqual(ints, []int64{1, 0, 2}) {
		t.Fatalf("StrToInt mismatch: %#v", ints)
	}

	gotStr := StrAppendIfNotExists([]string{`x`, `y`}, `y`)
	if !reflect.DeepEqual(gotStr, []string{`x`, `y`}) {
		t.Fatalf("StrAppendIfNotExists should keep original: %#v", gotStr)
	}
	gotStr = StrsAppendIfNotExists(gotStr, []string{`y`, `z`})
	if !reflect.DeepEqual(gotStr, []string{`x`, `y`, `z`}) {
		t.Fatalf("StrsAppendIfNotExists mismatch: %#v", gotStr)
	}

	gotInt := IntAppendIfNotExists([]int64{1}, 1)
	if !reflect.DeepEqual(gotInt, []int64{1}) {
		t.Fatalf("IntAppendIfNotExists should keep original: %#v", gotInt)
	}
	gotInt = IntsAppendIfNotExists(gotInt, []int64{1, 2, 3})
	if !reflect.DeepEqual(gotInt, []int64{1, 2, 3}) {
		t.Fatalf("IntsAppendIfNotExists mismatch: %#v", gotInt)
	}
}

func TestParseEmailAndFloatExist(t *testing.T) {
	emails := ParseEmail(` a@x.com, b@y.com , , c@z.com `, `Foo,Bar.<Q>(W)`)
	want := []string{
		`Foo_Bar__Q__W_<a@x.com>`,
		`Foo_Bar__Q__W_<b@y.com>`,
		`Foo_Bar__Q__W_<c@z.com>`,
	}
	if !reflect.DeepEqual(emails, want) {
		t.Fatalf("ParseEmail mismatch:\nwant=%#v\ngot=%#v", want, emails)
	}

	if !FloatExist([]float64{1.2, 3.4}, 3.4) {
		t.Fatalf("FloatExist should find value")
	}
	if FloatExist([]float64{1.2, 3.4}, 9.9) {
		t.Fatalf("FloatExist should not find value")
	}
}
