package Tr

import (
	"github.com/francoispqt/onelog"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/X"
	"github.com/tarantool/go-tarantool"
)

type Taran struct {
	*tarantool.Connection
	Log *onelog.Logger
}

// types
const Any = `any`
const Unsigned = `unsigned`
const String = `string`
const Number = `number`
const Double = `double`
const Integer = `integer`
const Boolean = `boolean`
const Decimal = `decimal`
const Uuid = `uuid`
const Scalar = `scalar`
const Array = `array`
const Map = `map`

// calls
const createIndex = `:create_index`
const format = `:format`

// misc
const Engine = `engine`
const Vinyl = `vinyl`
const Memtx = `memtx`
const BoxSpacePrefix = `box.space.`
const IfNotExists = `if_not_exists`
const sequenceSuffix = `_seq`

type TableProp struct {
	Droppable     bool // will drop table and recreate if failed to re-format the table
	Fields        []Field
	Unique        string
	Uniques       []string // multicolumn unique
	Indexes       []string
	AutoIncrement bool   // will set first column to auto increment integer/unsigned
	SpaceName     string // table name
	conn          *Taran
}

type Field struct { // https://godoc.org/gopkg.in/vmihailenco/msgpack.v2#pkg-examples
	Name       string `msgpack:"name"`
	Type       string `msgpack:"type"`
	IsNullable bool   `msgpack:"is_nullable"`
}

type Index struct {
	Parts       []string `msgpack:"parts"`
	IfNotExists bool     `msgpack:"if_not_exists"`
	Sequence    string   `msgpack:"sequence,omitempty"`
}

func (t *TableProp) AttachConnection(connection *Taran) {
	t.conn = connection
}

func (t *TableProp) TruncateTable() bool {
	return t.conn.CallBoxSpace(t.SpaceName+`:truncate`, A.X{})
}

func (t *TableProp) DropTable() bool {
	return t.conn.CallBoxSpace(t.SpaceName+`:drop`, A.X{})
}

func (prop *TableProp) CreateTable() bool {
	tableName := prop.SpaceName
	t := prop.conn
	if t == nil {
		t.Log.Fatal(`must attach connection`)
	}
	if tableName == `` {
		t.Log.Fatal(`table name cannot be empty`)
		return false
	}
	if !t.CallTarantool(`box.schema.space.create`, A.X{
		tableName,
		map[string]interface{}{
			Engine: Memtx,
			//IfNotExists: true,
		},
	}) {
		return false
	}
	if !t.CallBoxSpace(tableName+format, A.X{
		&prop.Fields,
	}) {
		return false
	}
	// create one field unique index
	t.CallBoxSpace(tableName+format, A.X{})
	// create auto increment
	if prop.AutoIncrement {
		t.CallTarantool(`box.schema.sequence.create`, A.X{tableName + sequenceSuffix, M.SX{
			IfNotExists: true,
		}})
	}
	if prop.Unique != `` {
		idx := Index{Parts: []string{prop.Unique}, IfNotExists: true}
		if prop.AutoIncrement {
			idx.Sequence = tableName + sequenceSuffix
		}
		t.CallBoxSpace(tableName+createIndex, A.X{
			prop.Unique, idx,
		})
	}
	// create multi-field unique index
	if len(prop.Uniques) > 2 {
		t.CallBoxSpace(tableName+createIndex, A.X{
			prop.Uniques[0], Index{Parts: prop.Uniques[1:], IfNotExists: true},
		})
	}
	// create other indexes
	for _, index := range prop.Indexes {
		t.CallBoxSpace(tableName+createIndex, A.X{
			index, Index{Parts: []string{index}, IfNotExists: true},
		})
	}
	return true
}

func (t *Taran) CallTarantool(funcName string, params A.X) bool {
	L.Print(funcName)
	L.Describe(params)
	res, err := t.Call(funcName, params)
	if err != nil && (len(params) == 0 || (len(params) > 0 && err.Error() != `Space '`+X.ToS(params[0])+`' already exists (0xa)`)) {
		L.IsError(err, `error CallTarantool`)
		return false
	}
	L.Describe(res.Tuples())
	return true
}

func (t *Taran) CallBoxSpace(funcName string, params A.X) bool {
	L.Print(funcName)
	L.Describe(params)
	res, err := t.Call(BoxSpacePrefix+funcName, params)
	if L.IsError(err, `error CallBoxSpace`) {
		return false
	}
	L.Describe(res.Tuples())
	return true
}

func (t *Taran) MigrateTarantool(prop *TableProp) {
	if !prop.CreateTable() && prop.Droppable {
		// drop table and recreate if error
		prop.DropTable()
		prop.CreateTable()
	}
}
