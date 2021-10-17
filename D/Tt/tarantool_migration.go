package Tt

import (
	"strings"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/X"
	"github.com/tarantool/go-tarantool"
)

const DEBUG = false

func Descr(args ...interface{}) {
	if DEBUG {
		L.ParentDescribe(args...)
	}
}

type TableName string

// types
type DataType string

const (
	Any      DataType = `any`
	Unsigned DataType = `unsigned`
	String   DataType = `string`
	Number   DataType = `number`
	Double   DataType = `double`
	Integer  DataType = `integer`
	Boolean  DataType = `boolean`
	Decimal  DataType = `decimal`
	Uuid     DataType = `string` // `uuid` // https://github.com/tarantool/go-tarantool/issues/90
	Scalar   DataType = `scalar`
	Array    DataType = `array`
	Map      DataType = `map`
)

var EmptyValue = map[DataType]string{
	Unsigned: `0`,
	String:   `''`,
	Integer:  `0`,
	Number:   `0`,
	Boolean:  `false`,
	Decimal:  `0`,
	Double:   `0`,
}

// misc
const Engine = `engine`

type EngineType string

const (
	Vinyl EngineType = `vinyl`
	Memtx EngineType = `memtx`
)
const BoxSpacePrefix = `box.space.`
const IfNotExists = `if_not_exists`

type TableProp struct {
	Fields       []Field
	Unique1      string
	Unique2      string
	Unique3      string
	Uniques      []string // multicolumn unique
	Indexes      []string
	Engine       EngineType
	HiddenFields []string
}

type Field struct { // https://godoc.org/gopkg.in/vmihailenco/msgpack.v2#pkg-examples
	Name string   `msgpack:"name"`
	Type DataType `msgpack:"type"`
}
type NullableField struct {
	Name       string   `msgpack:"name"`
	Type       DataType `msgpack:"type"`
	IsNullable bool     `msgpack:"is_nullable"`
}

type Index struct {
	Parts       []string `msgpack:"parts"`
	IfNotExists bool     `msgpack:"if_not_exists"`
	Unique      bool     `msgpack:"unique"`
}

type MSX map[string]interface{}

type Adapter struct {
	*tarantool.Connection
	Reconnect func() *tarantool.Connection
}

func (a *Adapter) UpsertTable(tableName TableName, prop *TableProp) bool {
	if DEBUG {
		L.Print(`---------------------------------------------------------`)
	}
	if prop.Engine == `` {
		prop.Engine = Vinyl
	}
	if !a.CreateSpace(string(tableName), prop.Engine) {
		return false
	}
	if !a.ReformatTable(string(tableName), prop.Fields) {
		return false // failed to create table
	}
	// create one field unique index
	a.ExecBoxSpace(string(tableName)+`:format`, A.X{})
	if prop.Unique1 != `` {
		a.ExecBoxSpace(string(tableName)+`:create_index`, A.X{
			prop.Unique1, Index{Parts: []string{prop.Unique1}, IfNotExists: true, Unique: true},
		})
	}
	if prop.Unique2 != `` {
		a.ExecBoxSpace(string(tableName)+`:create_index`, A.X{
			prop.Unique2, Index{Parts: []string{prop.Unique2}, IfNotExists: true, Unique: true},
		})
	}
	if prop.Unique3 != `` {
		a.ExecBoxSpace(string(tableName)+`:create_index`, A.X{
			prop.Unique3, Index{Parts: []string{prop.Unique3}, IfNotExists: true, Unique: true},
		})
	}
	// create multi-field unique index: [col1, col2] will named col1__col2
	if len(prop.Uniques) > 1 {
		a.ExecBoxSpace(string(tableName)+`:create_index`, A.X{
			strings.Join(prop.Uniques, `__`), Index{Parts: prop.Uniques, IfNotExists: true, Unique: true},
		})
	}
	// create other indexes
	for _, index := range prop.Indexes {
		//a.ExecBoxSpace(tableName+`.index.`+index+`:drop`, AX{index}) // TODO: remove this when index fixed
		a.ExecBoxSpace(string(tableName)+`:create_index`, A.X{
			index, Index{Parts: []string{index}, IfNotExists: true},
		})
	}
	return true
}

// ignore return value
func (a *Adapter) ExecTarantool(funcName string, params A.X) bool {
	return a.ExecTarantoolVerbose(funcName, params) == ``
}

func (a *Adapter) ExecTarantoolVerbose(funcName string, params A.X) string {
	Descr(funcName)
	Descr(params)
	res, err := a.Call(funcName, params)
	if err != nil && (len(params) == 0 || (len(params) > 0 && err.Error() != `Space '`+X.ToS(params[0])+`' already exists (0xa)`)) {
		L.IsError(err, `ExecTarantool failed: `+funcName)
		return err.Error()
	}
	Descr(res.Tuples())
	return ``
}

// ignore return value
func (a *Adapter) ExecBoxSpace(funcName string, params A.X) bool {
	return a.ExecBoxSpaceVerbose(funcName, params) == ``
}

func (a *Adapter) ExecBoxSpaceVerbose(funcName string, params A.X) string {
	Descr(funcName)
	Descr(params)
	res, err := a.Call(BoxSpacePrefix+funcName, params)
	if L.IsError(err, `ExecBoxSpace failed: `+funcName) {
		L.Describe(params)
		return err.Error()
	}
	Descr(res.Tuples())
	return ``
}

// ignore return value
func (a *Adapter) CallBoxSpace(funcName string, params A.X) (rows [][]interface{}) {
	Descr(funcName)
	Descr(params)
	res, err := a.Call(BoxSpacePrefix+funcName, params)
	Descr(res)
	if L.IsError(err, `ExecBoxSpace failed: `+funcName) {
		L.Describe(params)
		return
	}
	rows = res.Tuples()
	return
}

func (a *Adapter) DropTable(tableName string) bool {
	return a.ExecBoxSpace(tableName+`:drop`, A.X{})
}

func (a *Adapter) TruncateTable(tableName string) bool {
	return a.ExecBoxSpace(tableName+`:truncate`, A.X{})
}

func (a *Adapter) ReformatTable(tableName string, fields []Field) bool {
	// check old schema
	a.Connection = a.Reconnect() // need reconnect after creating space or a.Schema.Spaces will be empty
	schema := a.Schema
	table := schema.Spaces[tableName]
	if len(table.Fields) == 0 { // new table, create anyway
		if !a.ExecBoxSpace(tableName+`:format`, A.X{
			&fields,
		}) {
			return false // failed to create table
		}
		return true
	}
	// table already exists
	newFields := []NullableField{}
	nullFields := map[string]DataType{}
	for idx, newField := range fields { // diff and create nullable field
		origField := table.FieldsById[uint32(idx)]
		if origField != nil && origField.Type == string(newField.Type) {
			newFields = append(newFields, NullableField{Name: newField.Name, Type: newField.Type})
		} else {
			newFields = append(newFields, NullableField{Name: newField.Name, Type: newField.Type, IsNullable: true})
			nullFields[newField.Name] = newField.Type
		}
	}
	res := a.ExecBoxSpaceVerbose(tableName+`:format`, A.X{
		&newFields,
	})
	if res != `` {
		L.Describe(res)
	}
	a.Connection = a.Reconnect()
	// update all column
	if len(nullFields) > 0 {
		updateCols := []string{}
		for col, typ := range nullFields {
			defaultvalue := EmptyValue[typ]
			if defaultvalue == `` {
				panic(`please set EmptyValue for: ` + typ)
			}
			updateCols = append(updateCols, dq(col)+` = `+defaultvalue)
		}
		a.ExecSql(`UPDATE ` + dq(tableName) + ` SET ` + A.StrJoin(updateCols, `, `))
	}

	return a.ExecBoxSpace(tableName+`:format`, A.X{
		&fields,
	})
}

func (a *Adapter) CreateSpace(tableName string, engine EngineType) bool {
	err := a.ExecTarantoolVerbose(`box.schema.space.create`, A.X{
		tableName,
		map[string]interface{}{
			Engine: string(engine),
			//IfNotExists: true,
		},
	})
	if err == `` || S.StartsWith(err, `unsupported Lua type 'function'`) {
		return true // ignore
	}
	return true
}

func (a *Adapter) MigrateTables(tables map[TableName]*TableProp) {
	for name, props := range tables {
		Descr(name, props)
		a.UpsertTable(name, props)
	}
}
