package Ch

import (
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/L"
)

const DEBUG = true

func Descr(args ...interface{}) {
	if DEBUG {
		L.ParentDescribe(args...)
	}
}

// https://clickhouse.tech/docs/en/engines/table-engines/

const MergeTree = `MergeTree`
const ReplacingMergeTree = `ReplacingMergeTree`
const SummingMergeTree = `SummingMergeTree`
const AggregatingMergeTree = `AggregatingMergeTree`
const CollapsingMergeTree = `CollapsingMergeTree`
const VersionedCollapsingMergeTree = `VersionedCollapsingMergeTree`
const GraphiteMergeTree = `GraphiteMergeTree`

const TinyLog = `TinyLog`
const StripeLog = `StripeLog`
const Log = `Log`

const ODBC = `ODBC`
const JDBC = `JDBC`
const MySQL = `MySQL`
const MongoDB = `MongoDB`
const HDFS = `HDFS`
const S3 = `S3`
const Kafka = `Kafka`
const EmbeddedRocksDB = `EmbeddedRocksDB`
const RabbitMQ = `RabbitMQ`
const PostgreSQL = `PostgreSQL`

const Distributed = `Distributed`
const MaterializedView = `MaterializedView`
const Dictionary = `Dictionary`
const Merge = `Merge`
const File = `File`
const Null = `Null`
const Set = `Set`
const Join = `Join`
const URL = `URL`
const View = `View`
const Memory = `Memory`
const Buffer = `Buffer`

type TableName string

/*
SELECT DISTINCT alias_to
FROM system.data_type_families
ORDER BY alias_to ASC
*/
type DataType string

const (
	DateTime    DataType = `DateTime`
	DateTime64  DataType = `DateTime64`
	Decimal     DataType = `Decimal`
	FixedString DataType = `FixedString`
	Float32     DataType = `Float32`
	Float64     DataType = `Float64`
	IPv4        DataType = `IPv4`
	IPv6        DataType = `IPv6`
	Int16       DataType = `Int16`
	Int32       DataType = `Int32`
	Int64       DataType = `Int64`
	Int8        DataType = `Int8`
	String      DataType = `String`
	UInt16      DataType = `UInt16`
	UInt32      DataType = `UInt32`
	UInt64      DataType = `UInt64`
	UInt8       DataType = `UInt8`
)

type TableProp struct {
	Fields []Field
	Engine string
	Orders []string
}

type Field struct {
	Name string
	Type DataType
}

func (a *Adapter) CreateTable(tableName TableName, props *TableProp) bool {
	query := `
CREATE TABLE IF NOT EXISTS ` + string(tableName) + ` (`
	for idx, field := range props.Fields {
		query += field.Name + ` ` + string(field.Type)
		if idx < len(props.Fields)-1 {
			query += `,`
		}
	}
	query += `
) ENGINE = ` + props.Engine + `()`
	query += `
ORDER BY (` + A.StrJoin(props.Orders, `, `) + `)`
	_, err := a.Exec(query)
	return !L.IsError(err, `Exec: `+query)
}

func (a *Adapter) UpsertTable(tableName TableName, props *TableProp) bool {
	if props.Engine == `` {
		props.Engine = ReplacingMergeTree
	}
	if !a.CreateTable(tableName, props) {
		return false
	}
	if !a.AlterMissingColumns(tableName, props) {
		return false
	}

	return true
}

func (a *Adapter) MigrateTables(tables map[TableName]*TableProp) {
	for name, props := range tables {
		Descr(name, props)
		a.UpsertTable(name, props)
	}
}

func (a *Adapter) AlterMissingColumns(tableName TableName, props *TableProp) bool {
	query := `
SELECT
    "name",
    "type",
    "position"
FROM system.columns
WHERE "table" = $1
ORDER BY "position"`
	rows, err := a.Query(query, tableName)
	if L.IsError(err, `a.Query error: `+query) {
		return false
	}
	pos := 0
	hasError := false
	for rows.Next() {
		name := ``
		typ := ``
		err := rows.Scan(&name, &typ, &pos)
		if L.IsError(err, `rows.Scan(&name, &typ, &pos) error: `+query) {
			return false
		}
		if pos-1 < len(props.Fields) {
			idx := pos - 1
			field := props.Fields[idx]
			if name != field.Name {
				query := `ALTER TABLE ` + string(tableName) + ` RENAME COLUMN ` + name + ` TO ` + field.Name
				_, err := a.Exec(query)
				hasError = hasError || L.IsError(err, query)
			}
			if typ != string(field.Type) {
				query := `ALTER TABLE ` + string(tableName) + ` MODIFY COLUMN ` + field.Name + ` ` + string(field.Type)
				_, err := a.Exec(query)
				hasError = hasError || L.IsError(err, query)

			}
		}
	}
	for ; pos < len(props.Fields); pos++ {
		field := props.Fields[pos]
		query := `ALTER TABLE ` + string(tableName) + ` ADD COLUMN ` + field.Name + ` ` + string(field.Type)
		_, err := a.Exec(query)
		hasError = hasError || L.IsError(err, query)
	}

	return !hasError
}
