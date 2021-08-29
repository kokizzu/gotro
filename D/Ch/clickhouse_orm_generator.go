package Ch

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/X"
)

// quote import
func qi(importPath string) string {
	return `
	` + "`" + importPath + "`"
}

// backquote
func bq(importPath string) string {
	return "`" + importPath + "`"
}

// double quote
func dq(str string) string {
	return `"` + str + `"`
}

var typeTranslator = map[DataType]string{
	//Uuid:     `string`,
	UInt64:     `uint64`,
	Float64:    `float64`,
	String:     `string`,
	Int64:      `int64`,
	IPv4:       `string`,
	IPv6:       `string`,
	DateTime:   `time.Time`,
	DateTime64: `time.Time`,
	Int8:       `int8`,
}
var typeConverter = map[DataType]string{
	UInt64:  `X.ToU`,
	Float64: `X.ToF`,
	String:  `X.ToS`,
	Int64:   `X.ToI`,
	Int8:    `X.ToByte`,
}

const connStruct = `clickhouse.Adapter`
const connImport = "\n\n\t_ `github.com/ClickHouse/clickhouse-go`"
const buffImport = "\n\tchBuffer `github.com/kokizzu/ch-timed-buffer`"

const warning = "// DO NOT EDIT, will be overwritten by clickhouse_orm_generator.go\n\n"

func GenerateOrm(tables map[TableName]*TableProp) {
	ci := L.CallerInfo(2)
	this := L.CallerInfo()
	pkgName := S.RightOfLast(ci.PackageName, `/`)
	saPkgName := `sa` + pkgName[1:] // write/command (mutator)
	L.Print(saPkgName)
	mPkgName := `m` + pkgName[1:]

	// generate
	apBuf := bytes.Buffer{}

	SA := func(str string) {
		_, err := apBuf.WriteString(str)
		L.PanicIf(err, `failed apBuf.WriteString`)
	}

	//SA(`// generated: ` + time.Now().String() + "\n")
	SA(`package ` + saPkgName)
	SA("\n\n")
	SA(warning)

	SA(`import (`)

	// import reader
	SA(qi(`database/sql`))
	SA(qi(this.PackageName)) // /3rdparty/clickhouse
	SA(qi(ci.PackageName))   // /models/m*
	SA(qi(`time`))
	SA(connImport)
	SA(buffImport)

	SA("\n")
	SA(qi(`github.com/kokizzu/gotro/A`))
	SA(qi(`github.com/kokizzu/gotro/L`))
	//SA(qi(`github.com/kokizzu/gotro/X`))

	SA(`
)` + "\n\n")

	SA(`//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file ` + saPkgName + "__ORM.GEN.go\n")
	SA(`//go:generate replacer 'Id" form' 'Id,string" form' type ` + saPkgName + "__ORM.GEN.go\n")
	SA(`//go:generate replacer 'json:"id"' 'json:"id,string"' type ` + saPkgName + "__ORM.GEN.go\n")
	SA(`//go:generate replacer 'By" form' 'By,string" form' type ` + saPkgName + "__ORM.GEN.go\n")
	SA(`// go:generate msgp -tests=false -file ` + saPkgName + `__ORM.GEN.go -o ` + saPkgName + `__MSG.GEN.go` + "\n\n")

	// sort by table name to keep the order when regenerating structs
	tableNames := make([]string, 0, len(tables))
	for k := range tables {
		tableNames = append(tableNames, string(k))
	}
	sort.Strings(tableNames)

	for _, tableName := range tableNames {
		SA(`var ` + tableName + `Dummy = ` + S.CamelCase(tableName) + "{}\n")
	}

	SA("var Preparators = map[clickhouse.TableName]chBuffer.Preparator{\n")
	for _, tableName := range tableNames {
		SA(`	` + mPkgName + `.Table` + S.CamelCase(tableName) + `: func(tx *sql.Tx) *sql.Stmt {
		query := ` + tableName + `Dummy.sqlInsert()
		stmt, err := tx.Prepare(query)
		L.IsError(err, ` + "`failed to tx.Prepare: `+query" + `)
		return stmt
	},
`)
	}
	SA("}\n")

	// for each table generate in order
	for _, tableName := range tableNames {
		props := tables[TableName(tableName)]
		structName := S.CamelCase(tableName)
		maxLen := 1
		for _, prop := range props.Fields {
			l := len(prop.Name) + 1 - strings.Count(prop.Name, `_`)
			if maxLen < l {
				maxLen = l
			}
		}

		// reader struct
		SA(`type ` + structName + " struct {\n")
		SA("	Adapter *" + connStruct + " `json:" + `"-"` + " msg:" + `"-"` + " query:" + `"-"` + " form:" + `"-"` + "`\n")
		for _, prop := range props.Fields {
			camel := S.CamelCase(prop.Name)
			SA("	" + camel + strings.Repeat(` `, maxLen-len(camel)) + typeTranslator[prop.Type] + "\n")
		}
		SA("}\n\n")

		// reader struct constructor
		SA(`func New` + structName + `(adapter *` + connStruct + `) *` + structName + " {\n")
		SA(`	return &` + structName + "{Adapter: adapter}\n")
		SA("}\n\n")

		// table name
		receiverName := strings.ToLower(string(structName[0]))
		SA(`func (` + receiverName + ` ` + structName + ") TableName() clickhouse.TableName { //nolint:dupl false positive\n")
		SA("	return " + mPkgName + `.Table` + structName + "\n")
		SA("}\n\n")

		// sql table name
		SA(`func (` + receiverName + ` *` + structName + ") sqlTableName() string { //nolint:dupl false positive\n")
		SA("	return `" + `"` + tableName + `"` + "`\n")
		SA("}\n\n")

		// insert, error if exists
		SA("// insert, error if exists\n")
		SA(`func (` + receiverName + ` *` + structName + ") sqlInsert() string { //nolint:dupl false positive\n")
		qMark := S.Repeat(`,?`, len(props.Fields))[1:]
		SA("	return `INSERT INTO ` + " + receiverName + ".sqlTableName() + `(` + " + receiverName + ".sqlAllFields() + `) VALUES (" + qMark + ")`\n")
		SA("}\n\n")

		// total records
		SA(`func (` + receiverName + ` *` + structName + ") sqlCount() string { //nolint:dupl false positive\n")
		SA("	return `SELECT COUNT(*) FROM ` + " + receiverName + ".sqlTableName()\n")
		SA("}\n\n")

		// sql select all fields, used when need to mutate or show every fields
		SA(`func (` + receiverName + ` *` + structName + ") sqlSelectAllFields() string { //nolint:dupl false positive\n")
		sqlFields := ``
		for _, prop := range props.Fields {
			sqlFields += `, ` + (prop.Name) + "\n\t"
		}
		SA(`	return ` + bq(sqlFields[1:]) + "\n")
		SA("}\n\n")

		SA(`func (` + receiverName + ` *` + structName + ") sqlAllFields() string { //nolint:dupl false positive\n")
		sqlFields = ``
		for _, prop := range props.Fields {
			sqlFields += `, ` + prop.Name + ""
		}
		SA(`	return ` + bq(sqlFields[2:]) + "\n")
		SA("}\n\n")

		// to Insert parameter
		SA(`func (` + receiverName + ` ` + structName + ") SqlInsertParam() []interface{} { //nolint:dupl false positive\n")
		SA("	return []interface{}{\n")
		for idx, prop := range props.Fields {
			SA("		" + receiverName + "." + S.CamelCase(prop.Name) + ", // " + I.ToStr(idx) + " \n")
		}
		SA("	}\n")
		SA("}\n\n")

		for idx, prop := range props.Fields {
			propName := S.CamelCase(prop.Name)

			// index functions
			SA(`func (` + receiverName + ` *` + structName + ") Idx" + propName + "() int { //nolint:dupl false positive\n")
			SA("	return " + X.ToS(idx) + "\n")
			SA("}\n\n")

			// column name functions
			//SA(`func (` + receiverName + ` *` + structName + ") col" + propName + "() string { //nolint:dupl false positive\n")
			//SA("	return `" + prop.Name + "`\n")
			//SA("}\n\n")

			// sql column name functions
			SA(`func (` + receiverName + ` *` + structName + ") sql" + propName + "() string { //nolint:dupl false positive\n")
			SA("	return `" + prop.Name + "`\n")
			SA("}\n\n")
		}

		// to AX
		SA(`func (` + receiverName + ` *` + structName + ") ToArray() A.X { //nolint:dupl false positive\n")
		SA("	return A.X{\n")
		for idx, prop := range props.Fields {
			camel := S.CamelCase(prop.Name)
			SA("		" + receiverName + "." + camel + `,` + strings.Repeat(` `, maxLen-len(camel)) + `// ` + X.ToS(idx) + "\n")
		}
		SA("	}\n")
		SA("}\n\n")

		//// from AX
		//SA(`func (` + receiverName + ` *` + structName + `) FromArray(a A.X) *` + structName + " { //nolint:dupl false positive\n")
		//for idx, prop := range props.Fields {
		//	SA("	" + receiverName + "." + S.CamelCase(prop.Name) + ` = ` + typeConverter[prop.Type] + "(a[" + X.ToS(idx) + "])\n")
		//}
		//SA("	return " + receiverName + "\n")
		//SA("}\n\n")

		SA(warning)

	}

	err := os.MkdirAll(`./`+saPkgName, os.ModePerm)
	if L.IsError(err, `os.MkDir failed: `+saPkgName) {
		return
	}
	apFname := fmt.Sprintf(`./%s/%s__ORM.GEN.go`, saPkgName, saPkgName)
	err = os.WriteFile(apFname, apBuf.Bytes(), os.ModePerm)
	if L.IsError(err, `os.WriteFile failed: `+apFname) {
		return
	}
}

//func generateMutationByUniqueIndexumns(uniqueCamel, structProp, receiverName, structName string, AP, WC func(str string)) {
//
//	//// primary fields
//	//AP(`func (` + receiverName + ` *` + structName + ") PrimaryIndex() A.X { //nolint:dupl false positive\n")
//	//AP(`	return A.X{` + structProp + "}\n")
//	//AP("}\n\n")
//
//	// find by unique, used when need to mutate the object
//	AP(`func (` + receiverName + ` *` + structName + `) FindBy` + uniqueCamel + "() bool { //nolint:dupl false positive\n")
//	//AP("	res, err := " + receiverName + ".Adapter.Select(" + receiverName + ".SpaceName(), " + receiverName + `.UniqueIndex` + uniqueCamel + `(), 0, 1, ` + iterEq + ", A.X{" + structProp + "})\n")
//	AP("	if L.IsError(err, `" + structName + `.FindBy` + uniqueCamel + " failed: `+" + receiverName + ".SpaceName()) {\n")
//	AP("		return false\n")
//	AP("	}\n")
//	AP("	rows := res.Tuples()\n")
//	AP("	if len(rows) == 1 {\n")
//	AP("		" + receiverName + ".FromArray(rows[0])\n")
//	AP("		return true\n")
//	AP("	}\n")
//	AP("	return false\n")
//	AP("}\n\n")
//
//	// Overwrite all columns, error if not exists
//	WC("// Overwrite all columns, error if not exists\n")
//	WC(`func (` + receiverName + ` *` + structName + `Mutator) DoOverwriteBy` + uniqueCamel + "() bool { //nolint:dupl false positive\n")
//	WC("	_, err := " + receiverName + `.Adapter.Update(` + receiverName + ".SpaceName(), " + receiverName + `.UniqueIndex` + uniqueCamel + "(), A.X{" + structProp + "}, " + receiverName + ".ToUpdateArray())\n")
//	WC("	return !L.IsError(err, `" + structName + `.DoOverwriteBy` + uniqueCamel + " failed: `+" + receiverName + ".SpaceName())\n")
//	WC("}\n\n")
//
//	// Update only mutated, error if not exists
//	WC("// Update only mutated, error if not exists, use Find* and Set* methods instead of direct assignment\n")
//	WC(`func (` + receiverName + ` *` + structName + `Mutator) DoUpdateBy` + uniqueCamel + "() bool { //nolint:dupl false positive\n")
//	WC("	_, err := " + receiverName + `.Adapter.Update(` + receiverName + ".SpaceName(), " + receiverName + `.UniqueIndex` + uniqueCamel + "(), A.X{" + structProp + "}, " + receiverName + ".mutations)\n")
//	WC("	return !L.IsError(err, `" + structName + `.DoUpdateBy` + uniqueCamel + " failed: `+" + receiverName + ".SpaceName())\n")
//	WC("}\n\n")
//
//	// permanent delete
//	WC(`func (` + receiverName + ` *` + structName + `Mutator) DoDeletePermanentBy` + uniqueCamel + "() bool { //nolint:dupl false positive\n")
//	WC("	_, err := " + receiverName + ".Adapter.Delete(" + receiverName + ".SpaceName(), " + receiverName + `.UniqueIndex` + uniqueCamel + `(), A.X{` + structProp + "})\n")
//	WC("	return !L.IsError(err, `" + structName + `.DoDeletePermanentBy` + uniqueCamel + " failed: `+" + receiverName + ".SpaceName())\n")
//	WC("}\n\n")
//}
