package M

import (
	"testing"

	"github.com/hexops/autogold"
	"github.com/hexops/valast"
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
		StrVal       string      `json:"str,omitempty" ion:"foo"`
		IntValue     int         `json:"int,omitempty" ion:"bar"`
		WhateverItIs interface{} `json:"whatever"`
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
		SX       map[string]interface{}
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
			SX: map[string]interface{}{
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
			"SX": map[string]interface{}{
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
			SX: map[string]interface{}{
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
			"sx": map[string]interface{}{
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
			SX: map[string]interface{}{
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
			"sX": map[string]interface{}{
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
			SX: map[string]interface{}{
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
			"SX": map[string]interface{}{
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
			SX: map[string]interface{}{
				"a": 1,
				"b": 2,
			},
			Bool: true,
		})
		want.Equal(t, mv)
	})
}

// TODO: test time
// TODO: setting wrong value (eg. float to int, []byte to string, etc)
