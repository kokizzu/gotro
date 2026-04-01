package M

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/goccy/go-json"
	"github.com/hexops/autogold"
	"github.com/hexops/valast"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/B"
	"github.com/kokizzu/gotro/C"
	"github.com/kokizzu/gotro/F"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
)

func TestParseStruct(t *testing.T) {

	// to update this test:
	// cd M
	// go test -update .

	type AnonNested1 struct {
		Int int
		Str string
	}

	type Nested2 struct {
		Float float32
		Str   string
	}

	type Test1 struct {
		StrVal       string `json:"str,omitempty" ion:"foo"`
		IntValue     int    `json:"int,omitempty" ion:"bar"`
		WhateverItIs any    `json:"whatever"`
		private      float32
		Byte         byte
		I8           int8
		I16          int16
		I32          int32
		I64          int64
		U16          uint16
		U32          uint32
		U64          uint64
		F32          float32
		F64          float64
		Rune         rune
		AnonNested1
		Nested2  Nested2
		PI64     *int64
		PF64     *float64
		IntSlice []int
		SX       map[string]any
		Bool     bool
	}

	i64 := int64(567)
	f64 := float64(67.8)
	var mv Test1
	fillMv := func() {
		mv = Test1{
			StrVal:       "1",
			IntValue:     123,
			WhateverItIs: `whatever`,
			private:      234,
			Byte:         'b',
			I8:           2,
			I16:          -3,
			I32:          4,
			I64:          -5,
			U16:          6,
			U32:          7,
			U64:          8,
			F32:          -9.1,
			F64:          11.12,
			Rune:         'r',
			AnonNested1: AnonNested1{
				Int: 345,
				Str: "x",
			},
			Nested2: Nested2{
				Float: 456,
				Str:   "y",
			},
			PI64:     &i64,
			PF64:     &f64,
			IntSlice: []int{1, 2, 3},
			SX: map[string]any{
				`a`: 1,
				`b`: 2,
			},
			Bool: true,
		}
	}

	resetMv := func() {
		mv = Test1{}
	}

	t.Run(`raw`, func(t *testing.T) {
		sm := ParseStruct(&Test1{}, RawFieldName)
		want := autogold.Want(`raw`, &StructMapper{
			StructName: "M.Test1", Offset2key: map[uintptr]string{
				0:   "StrVal",
				16:  "IntValue",
				24:  "WhateverItIs",
				44:  "Byte",
				45:  "I8",
				46:  "I16",
				48:  "I32",
				56:  "I64",
				64:  "U16",
				68:  "U32",
				72:  "U64",
				80:  "F32",
				88:  "F64",
				96:  "Rune",
				104: "AnonNested1",
				128: "Nested2",
				152: "PI64",
				160: "PF64",
				168: "IntSlice",
				192: "SX",
				200: "Bool",
			},
			Key2offset: map[string]uintptr{
				"AnonNested1":  104,
				"Bool":         200,
				"Byte":         44,
				"F32":          80,
				"F64":          88,
				"I16":          46,
				"I32":          48,
				"I64":          56,
				"I8":           45,
				"IntSlice":     168,
				"IntValue":     16,
				"Nested2":      128,
				"PF64":         160,
				"PI64":         152,
				"Rune":         96,
				"SX":           192,
				"StrVal":       0,
				"U16":          64,
				"U32":          68,
				"U64":          72,
				"WhateverItIs": 24,
			},
			key2fieldName: map[string]string{
				"AnonNested1":  "AnonNested1",
				"Bool":         "Bool",
				"Byte":         "Byte",
				"F32":          "F32",
				"F64":          "F64",
				"I16":          "I16",
				"I32":          "I32",
				"I64":          "I64",
				"I8":           "I8",
				"IntSlice":     "IntSlice",
				"IntValue":     "IntValue",
				"Nested2":      "Nested2",
				"PF64":         "PF64",
				"PI64":         "PI64",
				"Rune":         "Rune",
				"SX":           "SX",
				"StrVal":       "StrVal",
				"U16":          "U16",
				"U32":          "U32",
				"U64":          "U64",
				"WhateverItIs": "WhateverItIs",
			},
		})
		want.Equal(t, sm)

		fillMv()
		m := sm.StructToMap(&mv)
		want = autogold.Want(`raw2map`, SX{
			"AnonNested1": AnonNested1{
				Int: 345,
				Str: "x",
			},
			"Bool": true,
			"Byte": 98,
			"F32":  -9.1,
			"F64":  11.12,
			"I16":  -3,
			"I32":  4,
			"I64":  -5,
			"I8":   2,
			"IntSlice": []int{
				1,
				2,
				3,
			},
			"IntValue": 123,
			"Nested2": Nested2{
				Float: 456,
				Str:   "y",
			},
			"PF64": valast.Addr(67.8).(*float64),
			"PI64": valast.Addr(int64(567)).(*int64),
			"Rune": 114,
			"SX": map[string]any{
				"a": 1,
				"b": 2,
			},
			"StrVal":       "1",
			"U16":          6,
			"U32":          7,
			"U64":          8,
			"WhateverItIs": "whatever",
		})
		want.Equal(t, m)

		resetMv()
		sm.MapToStruct(m, &mv)
		want = autogold.Want(`raw2struct`, Test1{
			StrVal: "1", IntValue: 123, WhateverItIs: "whatever",
			Byte: 98,
			I8:   2,
			I16:  -3,
			I32:  4,
			I64:  -5,
			U16:  6,
			U32:  7,
			U64:  8,
			F32:  -9.1,
			F64:  11.12,
			Rune: 114,
			AnonNested1: AnonNested1{
				Int: 345,
				Str: "x",
			},
			Nested2: Nested2{
				Float: 456,
				Str:   "y",
			},
			PI64: valast.Addr(int64(567)).(*int64),
			PF64: valast.Addr(67.8).(*float64),
			IntSlice: []int{
				1,
				2,
				3,
			},
			SX: map[string]any{
				"a": 1,
				"b": 2,
			},
			Bool: true,
		})
		want.Equal(t, mv)
	})

	t.Run(`json`, func(t *testing.T) {
		sm := ParseStruct(&Test1{}, `json`)
		want := autogold.Want(`json`, &StructMapper{
			StructName: "M.Test1", Offset2key: map[uintptr]string{
				0:  "str",
				16: "int",
				24: "whatever",
			},
			Key2offset: map[string]uintptr{
				"int":      16,
				"str":      0,
				"whatever": 24,
			},
			key2fieldName: map[string]string{
				"int":      "IntValue",
				"str":      "StrVal",
				"whatever": "WhateverItIs",
			},
		})
		want.Equal(t, sm)

		fillMv()
		m := sm.StructToMap(&mv)
		want = autogold.Want(`json2map`, SX{"": true, "int": 123, "str": "1", "whatever": "whatever"})
		want.Equal(t, m)

		resetMv()
		sm.MapToStruct(m, &mv)
		want = autogold.Want(`json2struct`, Test1{StrVal: "1", IntValue: 123, WhateverItIs: "whatever"})
		want.Equal(t, mv)
	})

	t.Run(`snake`, func(t *testing.T) {
		sm := ParseStruct(&Test1{}, SnakeFieldName)
		want := autogold.Want(`snake`, &StructMapper{
			StructName: "M.Test1", Offset2key: map[uintptr]string{
				0:   "str_val",
				16:  "int_value",
				24:  "whatever_it_is",
				44:  "byte",
				45:  "i_8",
				46:  "i_16",
				48:  "i_32",
				56:  "i_64",
				64:  "u_16",
				68:  "u_32",
				72:  "u_64",
				80:  "f_32",
				88:  "f_64",
				96:  "rune",
				104: "anon_nested_1",
				128: "nested_2",
				152: "pi_64",
				160: "pf_64",
				168: "int_slice",
				192: "sx",
				200: "bool",
			},
			Key2offset: map[string]uintptr{
				"anon_nested_1":  104,
				"bool":           200,
				"byte":           44,
				"f_32":           80,
				"f_64":           88,
				"i_16":           46,
				"i_32":           48,
				"i_64":           56,
				"i_8":            45,
				"int_slice":      168,
				"int_value":      16,
				"nested_2":       128,
				"pf_64":          160,
				"pi_64":          152,
				"rune":           96,
				"str_val":        0,
				"sx":             192,
				"u_16":           64,
				"u_32":           68,
				"u_64":           72,
				"whatever_it_is": 24,
			},
			key2fieldName: map[string]string{
				"anon_nested_1":  "AnonNested1",
				"bool":           "Bool",
				"byte":           "Byte",
				"f_32":           "F32",
				"f_64":           "F64",
				"i_16":           "I16",
				"i_32":           "I32",
				"i_64":           "I64",
				"i_8":            "I8",
				"int_slice":      "IntSlice",
				"int_value":      "IntValue",
				"nested_2":       "Nested2",
				"pf_64":          "PF64",
				"pi_64":          "PI64",
				"rune":           "Rune",
				"str_val":        "StrVal",
				"sx":             "SX",
				"u_16":           "U16",
				"u_32":           "U32",
				"u_64":           "U64",
				"whatever_it_is": "WhateverItIs",
			},
		})
		want.Equal(t, sm)

		fillMv()
		m := sm.StructToMap(&mv)
		want = autogold.Want(`snake2map`, SX{
			"anon_nested_1": AnonNested1{
				Int: 345,
				Str: "x",
			},
			"bool": true,
			"byte": 98,
			"f_32": -9.1,
			"f_64": 11.12,
			"i_16": -3,
			"i_32": 4,
			"i_64": -5,
			"i_8":  2,
			"int_slice": []int{
				1,
				2,
				3,
			},
			"int_value": 123,
			"nested_2": Nested2{
				Float: 456,
				Str:   "y",
			},
			"pf_64":   valast.Addr(67.8).(*float64),
			"pi_64":   valast.Addr(int64(567)).(*int64),
			"rune":    114,
			"str_val": "1",
			"sx": map[string]any{
				"a": 1,
				"b": 2,
			},
			"u_16":           6,
			"u_32":           7,
			"u_64":           8,
			"whatever_it_is": "whatever",
		})
		want.Equal(t, m)

		resetMv()
		sm.MapToStruct(m, &mv)
		want = autogold.Want(`snake2struct`, Test1{
			StrVal: "1", IntValue: 123, WhateverItIs: "whatever",
			Byte: 98,
			I8:   2,
			I16:  -3,
			I32:  4,
			I64:  -5,
			U16:  6,
			U32:  7,
			U64:  8,
			F32:  -9.1,
			F64:  11.12,
			Rune: 114,
			AnonNested1: AnonNested1{
				Int: 345,
				Str: "x",
			},
			Nested2: Nested2{
				Float: 456,
				Str:   "y",
			},
			PI64: valast.Addr(int64(567)).(*int64),
			PF64: valast.Addr(67.8).(*float64),
			IntSlice: []int{
				1,
				2,
				3,
			},
			SX: map[string]any{
				"a": 1,
				"b": 2,
			},
			Bool: true,
		})
		want.Equal(t, mv)
	})

	t.Run(`camel`, func(t *testing.T) {
		sm := ParseStruct(&Test1{}, CamelFieldName)
		want := autogold.Want(`camel`, &StructMapper{
			StructName: "M.Test1", Offset2key: map[uintptr]string{
				0:   "strVal",
				16:  "intValue",
				24:  "whateverItIs",
				44:  "byte",
				45:  "i8",
				46:  "i16",
				48:  "i32",
				56:  "i64",
				64:  "u16",
				68:  "u32",
				72:  "u64",
				80:  "f32",
				88:  "f64",
				96:  "rune",
				104: "anonNested1",
				128: "nested2",
				152: "pI64",
				160: "pF64",
				168: "intSlice",
				192: "sX",
				200: "bool",
			},
			Key2offset: map[string]uintptr{
				"anonNested1":  104,
				"bool":         200,
				"byte":         44,
				"f32":          80,
				"f64":          88,
				"i16":          46,
				"i32":          48,
				"i64":          56,
				"i8":           45,
				"intSlice":     168,
				"intValue":     16,
				"nested2":      128,
				"pF64":         160,
				"pI64":         152,
				"rune":         96,
				"sX":           192,
				"strVal":       0,
				"u16":          64,
				"u32":          68,
				"u64":          72,
				"whateverItIs": 24,
			},
			key2fieldName: map[string]string{
				"anonNested1":  "AnonNested1",
				"bool":         "Bool",
				"byte":         "Byte",
				"f32":          "F32",
				"f64":          "F64",
				"i16":          "I16",
				"i32":          "I32",
				"i64":          "I64",
				"i8":           "I8",
				"intSlice":     "IntSlice",
				"intValue":     "IntValue",
				"nested2":      "Nested2",
				"pF64":         "PF64",
				"pI64":         "PI64",
				"rune":         "Rune",
				"sX":           "SX",
				"strVal":       "StrVal",
				"u16":          "U16",
				"u32":          "U32",
				"u64":          "U64",
				"whateverItIs": "WhateverItIs",
			},
		})
		want.Equal(t, sm)

		fillMv()
		m := sm.StructToMap(&mv)
		want = autogold.Want(`camel2map`, SX{
			"anonNested1": AnonNested1{
				Int: 345,
				Str: "x",
			},
			"bool": true,
			"byte": 98,
			"f32":  -9.1,
			"f64":  11.12,
			"i16":  -3,
			"i32":  4,
			"i64":  -5,
			"i8":   2,
			"intSlice": []int{
				1,
				2,
				3,
			},
			"intValue": 123,
			"nested2": Nested2{
				Float: 456,
				Str:   "y",
			},
			"pF64": valast.Addr(67.8).(*float64),
			"pI64": valast.Addr(int64(567)).(*int64),
			"rune": 114,
			"sX": map[string]any{
				"a": 1,
				"b": 2,
			},
			"strVal":       "1",
			"u16":          6,
			"u32":          7,
			"u64":          8,
			"whateverItIs": "whatever",
		})
		want.Equal(t, m)

		resetMv()
		sm.MapToStruct(m, &mv)
		want = autogold.Want(`camel2struct`, Test1{
			StrVal: "1", IntValue: 123, WhateverItIs: "whatever",
			Byte: 98,
			I8:   2,
			I16:  -3,
			I32:  4,
			I64:  -5,
			U16:  6,
			U32:  7,
			U64:  8,
			F32:  -9.1,
			F64:  11.12,
			Rune: 114,
			AnonNested1: AnonNested1{
				Int: 345,
				Str: "x",
			},
			Nested2: Nested2{
				Float: 456,
				Str:   "y",
			},
			PI64: valast.Addr(int64(567)).(*int64),
			PF64: valast.Addr(67.8).(*float64),
			IntSlice: []int{
				1,
				2,
				3,
			},
			SX: map[string]any{
				"a": 1,
				"b": 2,
			},
			Bool: true,
		})
		want.Equal(t, mv)
	})

	t.Run(`ion`, func(t *testing.T) {
		sm := ParseStruct(&Test1{}, `ion`)
		want := autogold.Want(`ion`, &StructMapper{
			StructName: "M.Test1", Offset2key: map[uintptr]string{
				0:  "foo",
				16: "bar",
			},
			Key2offset: map[string]uintptr{
				"bar": 16,
				"foo": 0,
			},
			key2fieldName: map[string]string{
				"bar": "IntValue",
				"foo": "StrVal",
			},
		})
		want.Equal(t, sm)

		fillMv()
		m := sm.StructToMap(&mv)
		want = autogold.Want(`ion2map`, SX{"": true, "bar": 123, "foo": "1"})
		want.Equal(t, m)

		resetMv()
		sm.MapToStruct(m, &mv)
		want = autogold.Want(`ion2struct`, Test1{StrVal: "1", IntValue: 123})
		want.Equal(t, mv)
	})

	t.Run(`all`, func(t *testing.T) {
		sm := ParseStruct(&Test1{}, AllFieldName)
		want := autogold.Want(`all`, &StructMapper{
			StructName: "M.Test1", Offset2key: map[uintptr]string{
				0:   "StrVal",
				16:  "IntValue",
				24:  "WhateverItIs",
				44:  "Byte",
				45:  "I8",
				46:  "I16",
				48:  "I32",
				56:  "I64",
				64:  "U16",
				68:  "U32",
				72:  "U64",
				80:  "F32",
				88:  "F64",
				96:  "Rune",
				104: "AnonNested1",
				128: "Nested2",
				152: "PI64",
				160: "PF64",
				168: "IntSlice",
				192: "SX",
				200: "Bool",
			},
			Key2offset: map[string]uintptr{
				"AnonNested1":    104,
				"Bool":           200,
				"Byte":           44,
				"F32":            80,
				"F64":            88,
				"I16":            46,
				"I32":            48,
				"I64":            56,
				"I8":             45,
				"IntSlice":       168,
				"IntValue":       16,
				"Nested2":        128,
				"PF64":           160,
				"PI64":           152,
				"Rune":           96,
				"SX":             192,
				"StrVal":         0,
				"U16":            64,
				"U32":            68,
				"U64":            72,
				"WhateverItIs":   24,
				"anonNested1":    104,
				"anon_nested_1":  104,
				"bool":           200,
				"byte":           44,
				"f32":            80,
				"f64":            88,
				"f_32":           80,
				"f_64":           88,
				"i16":            46,
				"i32":            48,
				"i64":            56,
				"i8":             45,
				"i_16":           46,
				"i_32":           48,
				"i_64":           56,
				"i_8":            45,
				"intSlice":       168,
				"intValue":       16,
				"int_slice":      168,
				"int_value":      16,
				"nested2":        128,
				"nested_2":       128,
				"pF64":           160,
				"pI64":           152,
				"pf_64":          160,
				"pi_64":          152,
				"rune":           96,
				"sX":             192,
				"strVal":         0,
				"str_val":        0,
				"sx":             192,
				"u16":            64,
				"u32":            68,
				"u64":            72,
				"u_16":           64,
				"u_32":           68,
				"u_64":           72,
				"whateverItIs":   24,
				"whatever_it_is": 24,
			},
			key2fieldName: map[string]string{
				"AnonNested1":    "AnonNested1",
				"Bool":           "Bool",
				"Byte":           "Byte",
				"F32":            "F32",
				"F64":            "F64",
				"I16":            "I16",
				"I32":            "I32",
				"I64":            "I64",
				"I8":             "I8",
				"IntSlice":       "IntSlice",
				"IntValue":       "IntValue",
				"Nested2":        "Nested2",
				"PF64":           "PF64",
				"PI64":           "PI64",
				"Rune":           "Rune",
				"SX":             "SX",
				"StrVal":         "StrVal",
				"U16":            "U16",
				"U32":            "U32",
				"U64":            "U64",
				"WhateverItIs":   "WhateverItIs",
				"anonNested1":    "AnonNested1",
				"anon_nested_1":  "AnonNested1",
				"bool":           "Bool",
				"byte":           "Byte",
				"f32":            "F32",
				"f64":            "F64",
				"f_32":           "F32",
				"f_64":           "F64",
				"i16":            "I16",
				"i32":            "I32",
				"i64":            "I64",
				"i8":             "I8",
				"i_16":           "I16",
				"i_32":           "I32",
				"i_64":           "I64",
				"i_8":            "I8",
				"intSlice":       "IntSlice",
				"intValue":       "IntValue",
				"int_slice":      "IntSlice",
				"int_value":      "IntValue",
				"nested2":        "Nested2",
				"nested_2":       "Nested2",
				"pF64":           "PF64",
				"pI64":           "PI64",
				"pf_64":          "PF64",
				"pi_64":          "PI64",
				"rune":           "Rune",
				"sX":             "SX",
				"strVal":         "StrVal",
				"str_val":        "StrVal",
				"sx":             "SX",
				"u16":            "U16",
				"u32":            "U32",
				"u64":            "U64",
				"u_16":           "U16",
				"u_32":           "U32",
				"u_64":           "U64",
				"whateverItIs":   "WhateverItIs",
				"whatever_it_is": "WhateverItIs",
			},
		})
		want.Equal(t, sm)

		fillMv()
		m := sm.StructToMap(&mv)
		want = autogold.Want(`all2map`, SX{
			"AnonNested1": AnonNested1{
				Int: 345,
				Str: "x",
			},
			"Bool": true,
			"Byte": 98,
			"F32":  -9.1,
			"F64":  11.12,
			"I16":  -3,
			"I32":  4,
			"I64":  -5,
			"I8":   2,
			"IntSlice": []int{
				1,
				2,
				3,
			},
			"IntValue": 123,
			"Nested2": Nested2{
				Float: 456,
				Str:   "y",
			},
			"PF64": valast.Addr(67.8).(*float64),
			"PI64": valast.Addr(int64(567)).(*int64),
			"Rune": 114,
			"SX": map[string]any{
				"a": 1,
				"b": 2,
			},
			"StrVal":       "1",
			"U16":          6,
			"U32":          7,
			"U64":          8,
			"WhateverItIs": "whatever",
		})
		want.Equal(t, m)

		resetMv()
		sm.MapToStruct(m, &mv)
		want = autogold.Want(`all2struct`, Test1{
			StrVal: "1", IntValue: 123, WhateverItIs: "whatever",
			Byte: 98,
			I8:   2,
			I16:  -3,
			I32:  4,
			I64:  -5,
			U16:  6,
			U32:  7,
			U64:  8,
			F32:  -9.1,
			F64:  11.12,
			Rune: 114,
			AnonNested1: AnonNested1{
				Int: 345,
				Str: "x",
			},
			Nested2: Nested2{
				Float: 456,
				Str:   "y",
			},
			PI64: valast.Addr(int64(567)).(*int64),
			PF64: valast.Addr(67.8).(*float64),
			IntSlice: []int{
				1,
				2,
				3,
			},
			SX: map[string]any{
				"a": 1,
				"b": 2,
			},
			Bool: true,
		})
		want.Equal(t, mv)
	})
}

func json5fromMIB(orig map[int64]bool) string {
	b := bytes.Buffer{}
	b.WriteByte('{')
	first := true
	for k, v := range orig {
		if !first {
			b.WriteByte(',')
		} else {
			first = false
		}
		b.WriteString(I.ToS(k))
		b.WriteByte(':')
		b.WriteString(ToJson5(v))
	}
	b.WriteByte('}')
	return b.String()
}

func json5fromMIX(orig map[int64]any) string {
	b := bytes.Buffer{}
	b.WriteByte('{')
	first := true
	for k, v := range orig {
		if !first {
			b.WriteByte(',')
		} else {
			first = false
		}
		b.WriteString(I.ToS(k))
		b.WriteByte(':')
		b.WriteString(ToJson5(v))
	}
	b.WriteByte('}')
	return b.String()
}

func json5fromMIAX(orig map[int64][]any) string {
	b := bytes.Buffer{}
	b.WriteByte('{')
	first := true
	for k, v := range orig {
		if !first {
			b.WriteByte(',')
		} else {
			first = false
		}
		b.WriteString(I.ToS(k))
		b.WriteByte(':')
		b.WriteString(ToJson5(v))
	}
	b.WriteByte('}')
	return b.String()
}

func json5fromMSAX(orig map[string][]any) string {
	b := bytes.Buffer{}
	b.WriteByte('{')
	first := true
	for k, v := range orig {
		if !first {
			b.WriteByte(',')
		} else {
			first = false
		}
		b.WriteString(S.ZZ(k))
		b.WriteByte(':')
		b.WriteString(ToJson5(v))
	}
	b.WriteByte('}')
	return b.String()
}

func json5fromMSI(orig map[string]int64) string {
	b := bytes.Buffer{}
	b.WriteByte('{')
	first := true
	for k, v := range orig {
		if !first {
			b.WriteByte(',')
		} else {
			first = false
		}
		quote := true
		if len(k) > 0 {
			ch := k[0]
			if C.IsDigit(ch) && ch != '0' {
				for _, ch := range k[1:] {
					// find non digit
					if !C.IsDigit(uint8(ch)) {
						quote = true
						break
					}
				}
			} else if C.IsIdentStart(k[0]) {
				for _, ch := range k[1:] {
					// find non identifier character
					if !C.IsIdent(uint8(ch)) {
						quote = true
						break
					}
				}
			} else {
				quote = true
			}
		}
		if quote {
			k = S.Q(k)
		}
		b.WriteString(k)
		b.WriteByte(':')
		b.WriteString(I.ToS(v))
	}
	b.WriteByte('}')
	return b.String()
}

// ToJson5 convert to json5
func ToJson5(x any) string {
	// bug when using map[int64]any
	if x == nil {
		return `''`
	}
	switch orig := x.(type) {
	case bytes.Buffer: // return as is
		return orig.String()
	case string:
		return S.ZJJ(orig)
	case []byte:
		return S.ZJJ(string(orig))
	case int:
		return I.ToStr(orig)
	case int64:
		return I.ToS(orig)
	case int32:
		return I.ToS(int64(orig))
	case uint:
		return I.UToStr(orig)
	case uint64:
		return I.UToS(orig)
	case uint32:
		return I.UToS(uint64(orig))
	case float32:
		return F.ToS(float64(orig))
	case float64:
		return F.ToS(orig)
	case bool:
		return B.ToS(orig)
	case IB:
		return json5fromMIB(orig)
	case map[int64]bool:
		return json5fromMIB(orig)
	case IX:
		return json5fromMIX(orig)
	case map[int64]any:
		return json5fromMIX(orig)
	case IAX:
		return json5fromMIAX(orig)
	case map[int64][]any:
		return json5fromMIAX(orig)
	case SAX:
		return json5fromMSAX(orig)
	case map[string][]any:
		return json5fromMSAX(orig)
	case SX:
		return orig.ToJson()
	case map[string]any:
		return ToJson(orig)
	//   return any.(M.SX).ToJson()
	case SI:
		return json5fromMSI(orig)
	case map[string]int64:
		return json5fromMSI(orig)
	case A.X:
		return A.ToJson(orig)
	case []any:
		return A.ToJson(orig)
	default:
		str, err := json.Marshal(x)
		L.IsError(err, `X.ToJson5 failed`, x)
		return string(str)
	}
	// TODO: add more types (M/A) here, do not EVER TRY to use reflection in this case
}

type fastestSrc struct {
	Name  string
	Age   int64
	Score float64
}

func TestFastestMapToStruct(t *testing.T) {
	src := map[string]any{
		"Name":  "alice",
		"Age":   int64(21),
		"Score": 98.5,
	}
	var dst fastestSrc
	FastestMapToStruct(src, &dst)
	if dst.Name != "alice" || dst.Age != 21 || dst.Score != 98.5 {
		t.Fatalf("FastestMapToStruct mismatch: %#v", dst)
	}
}

func TestFastestStructToMap(t *testing.T) {
	src := fastestSrc{
		Name:  "bob",
		Age:   33,
		Score: 77.25,
	}
	got := FastestStructToMap(src)
	want := map[string]any{
		"Name":  "bob",
		"Age":   int64(33),
		"Score": 77.25,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FastestStructToMap mismatch:\nwant=%#v\ngot=%#v", want, got)
	}
}

func TestFastestCopyStruct(t *testing.T) {
	src := fastestSrc{
		Name:  "carol",
		Age:   45,
		Score: 88.8,
	}

	var dstStruct fastestSrc
	FastestCopyStruct(src, &dstStruct)
	if !reflect.DeepEqual(src, dstStruct) {
		t.Fatalf("FastestCopyStruct to struct mismatch:\nwant=%#v\ngot=%#v", src, dstStruct)
	}

	dstMap := map[string]any{}
	FastestCopyStruct(src, &dstMap)
	wantMap := map[string]any{
		"Name":  "carol",
		"Age":   int64(45),
		"Score": 88.8,
	}
	if !reflect.DeepEqual(dstMap, wantMap) {
		t.Fatalf("FastestCopyStruct to map mismatch:\nwant=%#v\ngot=%#v", wantMap, dstMap)
	}
}

type structWrapA struct {
	A int
	B string
	c bool
}

type structWrapB struct {
	A int
	B string
}

func TestStructWrappersAndInvalidPaths(t *testing.T) {
	src := &structWrapA{A: 7, B: "x", c: true}

	// StructMap cache path and wrapper FromStruct/ToStruct.
	sm1 := StructMap(src)
	sm2 := StructMap(src)
	if sm1 != sm2 {
		t.Fatalf("StructMap should return cached mapper for same type")
	}

	m := FromStruct(src)
	if m["A"] != 7 || m["B"] != "x" {
		t.Fatalf("FromStruct mismatch: %#v", m)
	}
	if _, ok := m["c"]; ok {
		t.Fatalf("FromStruct should skip unexported fields: %#v", m)
	}

	var dst structWrapA
	m.ToStruct(&dst)
	if dst.A != 7 || dst.B != "x" || dst.c {
		t.Fatalf("ToStruct mismatch: %#v", dst)
	}

	// MapToStruct should skip nil value and unknown field.
	sm := ParseStruct(&structWrapA{}, RawFieldName)
	dst2 := structWrapA{A: 9, B: "old", c: true}
	sm.MapToStruct(SX{"A": nil, "B": "new", "Unknown": 1}, &dst2)
	if dst2.A != 9 || dst2.B != "new" || !dst2.c {
		t.Fatalf("MapToStruct nil/unknown handling mismatch: %#v", dst2)
	}

	// invalid type (not struct) and different struct type paths
	mm := map[string]any{}
	sm.MapToStruct(SX{"A": 1}, &mm) // should not panic
	other := structWrapB{}
	sm.MapToStruct(SX{"A": 1, "B": "y"}, &other)
	if other.A != 0 || other.B != "" {
		t.Fatalf("MapToStruct should ignore different struct type: %#v", other)
	}

	if got := sm.StructToMap(&mm); len(got) != 0 {
		t.Fatalf("StructToMap invalid type should be empty: %#v", got)
	}
	if got := sm.StructToMap(&other); len(got) != 0 {
		t.Fatalf("StructToMap different type should be empty: %#v", got)
	}
}

func TestParseStructNilPointer(t *testing.T) {
	var nilPtr *structWrapA
	sm := ParseStruct(nilPtr, RawFieldName)
	if sm == nil {
		t.Fatalf("ParseStruct(nil pointer) should return non-nil mapper")
	}
	if sm.StructName != "" {
		t.Fatalf("ParseStruct(nil pointer) should keep empty struct name, got %q", sm.StructName)
	}
}

// TODO: test time
// TODO: setting wrong value (eg. float to int, []byte to string, etc)
