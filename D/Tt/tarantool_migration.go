package Tt

import (
	"errors"
	"strings"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/X"
)

const DEBUG = false

var ErrMiscaonfiguration = errors.New(`misconfiguration`)

func Descr(args ...any) {
	if DEBUG {
		L.ParentDescribe(args...)
	}
}

type TableName string

// types
type DataType string

const (
	// https://www.tarantool.io/en/doc/latest/concepts/data_model/value_store/#data-types

	Unsigned DataType = `unsigned`
	String   DataType = `string`
	Double   DataType = `double`
	Integer  DataType = `integer`
	Boolean  DataType = `boolean`
	Array    DataType = `array`

	//Any      DataType = `any`
	//Number   DataType = `number` // use double instead
	//Decimal  DataType = `decimal` // unsupported
	//DateTime DataType = `datetime` // unsupported
	//Uuid     DataType = `string`  // `uuid` // https://github.com/tarantool/go-tarantool/issues/90
	//Scalar   DataType = `scalar`
	//Map      DataType = `map`

	//ArrayFloat    DataType = `array`
	//ArrayUnsigned DataType = `array`
	//ArrayInteger  DataType = `array`
)

var TypeToConst = map[DataType]string{
	Unsigned: `Tt.Unsigned`,
	String:   `Tt.String`,
	Integer:  `Tt.Integer`,
	Double:   `Tt.Double`,
	Boolean:  `Tt.Boolean`,
	Array:    `Tt.Array`,
}

var TypeToGoType = map[DataType]string{
	//Uuid:     `string`,
	Unsigned: `uint64`,
	String:   `string`,
	Integer:  `int64`,
	Double:   `float64`,
	Boolean:  `bool`,
	Array:    `[]any`,
}
var TypeToGoEmptyValue = map[DataType]string{
	Unsigned: `0`,
	String:   "``",
	Integer:  `0`,
	Double:   `0`,
	Boolean:  `false`,
	Array:    `[]any{}`,
}
var TypeToConvertFunc = map[DataType]string{
	Unsigned: `X.ToU`,
	String:   `X.ToS`,
	Integer:  `X.ToI`,
	Double:   `X.ToF`,
	Boolean:  `X.ToBool`,
	Array:    `X.ToArr`,
}

var Type2TarantoolDefault = map[DataType]string{
	Unsigned: `0`,
	String:   `''`,
	Integer:  `0`,
	Boolean:  `false`,
	Double:   `0`,
	Array:    `nil`,
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

const IdCol = `id`

type TableProp struct {
	Fields []Field

	// indexes
	Unique1 string
	Unique2 string
	Unique3 string
	Uniques []string // multicolumn unique
	Indexes []string
	Spatial string

	Engine          EngineType
	HiddenFields    []string
	AutoIncrementId bool // "id" column will be used to generate sequence, can only be created at beginning
	GenGraphqlType  bool

	// hook
	PreReformatMigrationHook func(*Adapter)
	PreUnique1MigrationHook  func(*Adapter)
	PreUnique2MigrationHook  func(*Adapter)
	PreUnique3MigrationHook  func(*Adapter)
	PreUniquesMigrationHook  func(*Adapter)
	PreSpatialMigrationHook  func(*Adapter)

	AutoCensorFields []string // fields to automatically censor
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
	Sequence    string   `msgpack:"sequence,omitempty"`
	Type        string   `msgpack:"type,omitempty"`
}

type MSX map[string]any

func (a *Adapter) UpsertTable(tableName TableName, prop *TableProp) bool {
	if DEBUG {
		L.Print(`---UpsertTable-Start-` + tableName + `----------------------------------------------------`)
		defer L.Print(`---UpsertTable-End-` + tableName + `-------------------------------------------------------`)
	}
	if prop.Engine == `` {
		prop.Engine = Vinyl
	}
	if !a.CreateSpace(string(tableName), prop.Engine) {
		return false
	}
	if prop.PreReformatMigrationHook != nil {
		prop.PreReformatMigrationHook(a)
	}
	if !a.ReformatTable(string(tableName), prop.Fields) {
		return false // failed to create table
	}
	// create one field unique index
	a.ExecBoxSpace(string(tableName)+`:format`, A.X{})
	if prop.AutoIncrementId {
		if len(prop.Fields) < 1 || prop.Fields[0].Name != IdCol || prop.Fields[0].Type != Unsigned {
			panic(`must create Unsigned id field on first field to use AutoIncrementId`)
		}

		seqName := string(tableName) + `_` + IdCol
		a.ExecTarantoolVerbose(`box.schema.sequence.create`, A.X{
			seqName,
		})
		a.ExecBoxSpace(string(tableName)+`:create_index`, A.X{
			IdCol, Index{
				Sequence:    seqName,
				Parts:       []string{IdCol},
				IfNotExists: true,
				Unique:      true,
			},
		})
	}
	if prop.PreUnique3MigrationHook != nil {
		prop.PreUnique3MigrationHook(a)
	}
	// only create unique if not "id"
	if prop.Unique1 != `` && !(prop.AutoIncrementId && prop.Unique1 == IdCol) {
		a.ExecBoxSpace(string(tableName)+`:create_index`, A.X{
			prop.Unique1, Index{Parts: []string{prop.Unique1}, IfNotExists: true, Unique: true},
		})
		if prop.Unique2 != `` && prop.Unique1 == prop.Unique2 {
			panic(`Unique1 and Unique2 must be unique`)
		}
		if prop.Unique3 != `` && prop.Unique1 == prop.Unique3 {
			panic(`Unique1 and Unique3 must be unique`)
		}
	}
	if prop.PreUnique2MigrationHook != nil {
		prop.PreUnique2MigrationHook(a)
	}
	if prop.Unique2 != `` && !(prop.AutoIncrementId && prop.Unique2 == IdCol) {
		a.ExecBoxSpace(string(tableName)+`:create_index`, A.X{
			prop.Unique2, Index{Parts: []string{prop.Unique2}, IfNotExists: true, Unique: true},
		})
		if prop.Unique3 != `` && prop.Unique2 == prop.Unique3 {
			panic(`Unique2 and Unique3 must be unique`)
		}
	}
	if prop.PreUnique3MigrationHook != nil {
		prop.PreUnique3MigrationHook(a)
	}
	if prop.Unique3 != `` && !(prop.AutoIncrementId && prop.Unique3 == IdCol) {
		a.ExecBoxSpace(string(tableName)+`:create_index`, A.X{
			prop.Unique3, Index{Parts: []string{prop.Unique3}, IfNotExists: true, Unique: true},
		})
	}
	if prop.PreUniquesMigrationHook != nil {
		prop.PreUniquesMigrationHook(a)
	}
	// create multi-field unique index: [col1, col2] will named col1__col2
	if len(prop.Uniques) > 1 {
		a.ExecBoxSpace(string(tableName)+`:create_index`, A.X{
			strings.Join(prop.Uniques, `__`), Index{Parts: prop.Uniques, IfNotExists: true, Unique: true},
		})
	}
	if prop.PreSpatialMigrationHook != nil {
		prop.PreSpatialMigrationHook(a)
	}
	// create spatial index (only works for memtx)
	if prop.Spatial != `` {
		a.ExecBoxSpace(string(tableName)+`:create_index`, A.X{
			prop.Spatial, Index{Parts: []string{prop.Spatial}, IfNotExists: true, Type: `RTREE`},
		})
	}
	// create other indexes
	for _, index := range prop.Indexes {
		//a.ExecBoxSpace(tableName+`.index.`+index+`:drop`, AX{index}) // TODO: remove this when index fixed
		a.ExecBoxSpace(string(tableName)+`:create_index`, A.X{
			index, Index{Parts: []string{index}, IfNotExists: true},
		})
	}
	// need refresh index after migrate
	// https://github.com/tarantool/go-tarantool/pull/259#pullrequestreview-1242058107
	a.Connection = a.Reconnect()
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
	if err != nil {
		if len(params) > 0 {
			errStr := err.Error()
			if errStr == `Space '`+X.ToS(params[0])+`' already exists` ||
				errStr == `Sequence '`+X.ToS(params[0])+`' already exists` {
				L.IsError(err, `ExecTarantool failed: `+funcName)
				return err.Error()
			}
		}
		if len(params) == 0 {
			L.IsError(err, `ExecTarantool failed: `+funcName)
			return err.Error()
		}
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
func (a *Adapter) CallBoxSpace(funcName string, params A.X) (rows [][]any) {
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
		return a.ExecBoxSpace(tableName+`:format`, A.X{
			&fields,
		}) // failed to create table
	}
	// table already exists
	var newFields []NullableField
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
			defaultvalue := Type2TarantoolDefault[typ]
			if defaultvalue == `` {
				panic(`please set Type2TarantoolDefault for: ` + typ)
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
		map[string]any{
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
		// validity check
		if props.Engine == Vinyl && props.Spatial != `` {
			L.PanicIf(ErrMiscaonfiguration, `spatial index is not supported in vinyl engine`)
		}
		a.UpsertTable(name, props)
	}
}
