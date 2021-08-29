package Tt

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
	Unsigned: `uint64`,
	Number:   `float64`,
	String:   `string`,
	Integer:  `int64`,
	Boolean:  `bool`,
}
var typeConverter = map[DataType]string{
	Unsigned: `X.ToU`,
	Number:   `X.ToF`,
	String:   `X.ToS`,
	Integer:  `X.ToI`,
	Boolean:  `X.ToBool`,
}

const connStruct = `tarantool.Adapter`
const connImport = "\n\n\tgotarantool `github.com/tarantool/go-tarantool`"
const iterEq = `gotarantool.IterEq`
const iterAll = `gotarantool.IterAll`

const warning = "// DO NOT EDIT, will be overwritten by tarantool_orm_generator.go\n\n"

func GenerateOrm(tables map[TableName]*TableProp) {
	ci := L.CallerInfo(2)
	this := L.CallerInfo()
	pkgName := S.RightOfLast(ci.PackageName, `/`)
	wcPkgName := `wc` + pkgName[1:] // write/command (mutator)
	rqPkgName := `rq` + pkgName[1:] // read/query (reader)
	mPkgName := `m` + pkgName[1:]
	L.Print(rqPkgName, wcPkgName)

	//var maxMap = map[string]string{
	//	Unsigned: `math.MaxInt32`,
	//	Number:   `math.MaxFloat32`,
	//	String:   "``",
	//}
	//var minMap = map[string]string{
	//	Unsigned: `0`,
	//	Number:   `-math.MaxFloat32`,
	//	String:   "``",
	//}

	//// do not generate when no table files changed
	//maxModTime := int64(0)
	//stat, err := os.Stat(genDir + genRqFilename)
	//if err == nil {
	//	err = filepath.Walk(genDir, func(path string, info os.FileInfo, err error) error {
	//		if strings.Contains(path, `_table_`) || strings.Contains(path, `_schema.go`) {
	//			modTime := info.ModTime().UnixNano()
	//			if maxModTime < modTime {
	//				maxModTime = modTime
	//			}
	//		}
	//		return nil
	//	})
	//	if L.IsError(err, `filepath.Walk failed: `+genDir) {
	//		return
	//	}
	//	// no table file changed
	//	if stat.ModTime().UnixNano() >= maxModTime {
	//		return
	//	}
	//}

	// generate
	rqBuf := bytes.Buffer{}
	wcBuf := bytes.Buffer{}

	RQ := func(str string) {
		_, err := rqBuf.WriteString(str)
		L.PanicIf(err, `failed rqBuf.WriteString`)
	}
	WC := func(str string) {
		_, err := wcBuf.WriteString(str)
		L.PanicIf(err, `failed wcBuf.WriteString`)
	}
	BOTH := func(str string) {
		RQ(str)
		WC(str)
	}

	//BOTH(`// generated: ` + time.Now().String() + "\n")
	RQ(`package ` + rqPkgName)
	WC(`package ` + wcPkgName)
	BOTH("\n\n")
	BOTH(warning)

	BOTH(`import (`)

	// import reader
	BOTH(qi(this.PackageName))
	RQ(qi(ci.PackageName))
	RQ(connImport)
	WC(qi(ci.PackageName + `/` + rqPkgName))

	BOTH("\n")
	BOTH(qi(`github.com/kokizzu/gotro/A`))
	BOTH(qi(`github.com/kokizzu/gotro/L`))

	RQ(qi(`github.com/kokizzu/gotro/X`))

	BOTH(`
)` + "\n\n")

	RQ(`//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file ` + rqPkgName + "__ORM.GEN.go\n")
	RQ(`//go:generate replacer 'Id" form' 'Id,string" form' type ` + rqPkgName + "__ORM.GEN.go\n")
	RQ(`//go:generate replacer 'json:"id"' 'json:"id,string"' type ` + rqPkgName + "__ORM.GEN.go\n")
	RQ(`//go:generate replacer 'By" form' 'By,string" form' type ` + rqPkgName + "__ORM.GEN.go\n")
	RQ(`// go:generate msgp -tests=false -file ` + rqPkgName + `__ORM.GEN.go -o ` + rqPkgName + `__MSG.GEN.go` + "\n\n")

	WC(`//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file ` + wcPkgName + "__ORM.GEN.go\n")
	WC(`//go:generate replacer 'Id" form' 'Id,string" form' type ` + wcPkgName + "__ORM.GEN.go\n")
	WC(`//go:generate replacer 'json:"id"' 'json:"id,string"' type ` + wcPkgName + "__ORM.GEN.go\n")
	WC(`//go:generate replacer 'By" form' 'By,string" form' type ` + wcPkgName + "__ORM.GEN.go\n")
	WC(`// go:generate msgp -tests=false -file ` + wcPkgName + `__ORM.GEN.go -o ` + wcPkgName + `__MSG.GEN.go` + "\n\n")

	// sort by table name to keep the order when regenerating structs
	tableNames := make([]string, 0, len(tables))
	for k := range tables {
		tableNames = append(tableNames, string(k))
	}
	sort.Strings(tableNames)

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

		// mutator struct
		WC(`type ` + structName + "Mutator struct {\n")
		WC(`	` + rqPkgName + `.` + structName + "\n")
		WC("	mutations []A.X\n")
		WC("}\n\n")

		// mutator struct constructor
		WC(`func New` + structName + `Mutator(adapter *` + connStruct + `) *` + structName + "Mutator {\n")
		WC(`	return &` + structName + `Mutator{` + structName + `: ` + rqPkgName + `.` + structName + "{Adapter: adapter}}\n")
		WC("}\n\n")

		// reader struct
		RQ(`type ` + structName + " struct {\n")
		RQ("	Adapter *" + connStruct + " `json:" + `"-"` + " msg:" + `"-"` + " query:" + `"-"` + " form:" + `"-"` + "`\n")
		for _, prop := range props.Fields {
			camel := S.CamelCase(prop.Name)
			RQ("	" + camel + strings.Repeat(` `, maxLen-len(camel)) + typeTranslator[prop.Type] + "\n")
		}
		RQ("}\n\n")

		// reader struct constructor
		RQ(`func New` + structName + `(adapter *` + connStruct + `) *` + structName + " {\n")
		RQ(`	return &` + structName + "{Adapter: adapter}\n")
		RQ("}\n\n")

		// table name
		receiverName := strings.ToLower(string(structName[0]))
		RQ(`func (` + receiverName + ` *` + structName + ") SpaceName() string { //nolint:dupl false positive\n")
		RQ("	return string(" + mPkgName + `.Table` + structName + ")\n")
		RQ("}\n\n")

		// sql table name
		RQ(`func (` + receiverName + ` *` + structName + ") sqlTableName() string { //nolint:dupl false positive\n")
		RQ("	return `" + `"` + tableName + `"` + "`\n")
		RQ("}\n\n")

		// have mutation
		WC(`func (` + receiverName + ` *` + structName + "Mutator) HaveMutation() bool { //nolint:dupl false positive\n")
		WC(`	return len(` + receiverName + ".mutations) > 0\n")
		WC("}\n\n")

		// unique index1
		if props.Unique1 != `` {
			uniquePropCamel := S.CamelCase(props.Unique1)
			structProp := receiverName + `.` + uniquePropCamel

			RQ(`func (` + receiverName + ` *` + structName + `) UniqueIndex` + uniquePropCamel + "() string { //nolint:dupl false positive\n")
			RQ("	return `" + props.Unique1 + "`\n")
			RQ("}\n\n")

			generateMutationByUniqueIndexumns(uniquePropCamel, structProp, receiverName, structName, RQ, WC)

			// upsert template, to be copied when need increment some field
			WC(`// func (` + receiverName + ` *` + structName + "Mutator) DoUpsert() bool { //nolint:dupl false positive\n")
			WC("//	_, err := " + receiverName + ".Adapter.Upsert(" + receiverName + ".SpaceName(), " + receiverName + ".ToArray(), A.X{\n")
			for idx, prop := range props.Fields {
				WC("//		A.X{`=`, " + X.ToS(idx) + ", " + receiverName + "." + S.CamelCase(prop.Name) + "},\n")
			}
			WC("//	})\n")
			WC("//	return !L.IsError(err, `" + structName + ".DoUpsert failed: `+" + receiverName + ".SpaceName())\n")
			WC("// }\n\n")
		}

		// unique index2
		if props.Unique2 != `` {
			uniquePropCamel := S.CamelCase(props.Unique2)
			structProp := receiverName + `.` + uniquePropCamel

			RQ(`func (` + receiverName + ` *` + structName + `) UniqueIndex` + uniquePropCamel + "() string { //nolint:dupl false positive\n")
			RQ("	return `" + props.Unique2 + "`\n")
			RQ("}\n\n")

			generateMutationByUniqueIndexumns(uniquePropCamel, structProp, receiverName, structName, RQ, WC)
		}

		// unique index3
		if props.Unique3 != `` {
			uniquePropCamel := S.CamelCase(props.Unique3)
			structProp := receiverName + `.` + uniquePropCamel

			RQ(`func (` + receiverName + ` *` + structName + `) UniqueIndex` + uniquePropCamel + "() string { //nolint:dupl false positive\n")
			RQ("	return `" + props.Unique3 + "`\n")
			RQ("}\n\n")

			generateMutationByUniqueIndexumns(uniquePropCamel, structProp, receiverName, structName, RQ, WC)
		}

		// unique indexes
		if len(props.Uniques) > 0 {
			uniquePropCamel := ``
			structProps := ``
			for _, uniq := range props.Uniques {
				uniquePropCamel += S.CamelCase(uniq)
				structProps += `, ` + receiverName + `.` + S.CamelCase(uniq)
			}
			if len(structProps) > 2 {
				structProps = structProps[2:]
			}

			RQ(`func (` + receiverName + ` *` + structName + `) UniqueIndex` + uniquePropCamel + "() string { //nolint:dupl false positive\n")
			RQ("	return `" + strings.Join(props.Uniques, `__`) + "`\n")
			RQ("}\n\n")

			generateMutationByUniqueIndexumns(uniquePropCamel, structProps, receiverName, structName, RQ, WC)
		}

		// insert, error if exists
		WC("// insert, error if exists\n")
		WC(`func (` + receiverName + ` *` + structName + "Mutator) DoInsert() bool { //nolint:dupl false positive\n")
		WC("	_, err := " + receiverName + ".Adapter.Insert(" + receiverName + ".SpaceName(), " + receiverName + ".ToArray())\n")
		WC("	return !L.IsError(err, `" + structName + ".DoInsert failed: `+" + receiverName + ".SpaceName())\n")
		WC("}\n\n")

		// replace = upsert, only error when there's unique secondary key
		WC("// replace = upsert, only error when there's unique secondary key\n")
		WC(`func (` + receiverName + ` *` + structName + "Mutator) DoReplace() bool { //nolint:dupl false positive\n")
		WC("	_, err := " + receiverName + ".Adapter.Replace(" + receiverName + ".SpaceName(), " + receiverName + ".ToArray())\n")
		WC("	return !L.IsError(err, `" + structName + ".DoReplace failed: `+" + receiverName + ".SpaceName())\n")
		WC("}\n\n")

		// sql select all fields, used when need to mutate or show every fields
		RQ(`func (` + receiverName + ` *` + structName + ") sqlSelectAllFields() string { //nolint:dupl false positive\n")
		sqlFields := ``
		for _, prop := range props.Fields {
			sqlFields += `, ` + dq(prop.Name) + "\n\t"
		}
		RQ(`	return ` + bq(sqlFields[1:]) + "\n")
		RQ("}\n\n")

		// to Update AX
		RQ(`func (` + receiverName + ` *` + structName + ") ToUpdateArray() A.X { //nolint:dupl false positive\n")
		RQ("	return A.X{\n")
		for idx, prop := range props.Fields {
			RQ("		A.X{`=`, " + X.ToS(idx) + ", " + receiverName + "." + S.CamelCase(prop.Name) + "},\n")
		}
		RQ("	}\n")
		RQ("}\n\n")

		for idx, prop := range props.Fields {
			propName := S.CamelCase(prop.Name)

			// index functions
			RQ(`func (` + receiverName + ` *` + structName + ") Idx" + propName + "() int { //nolint:dupl false positive\n")
			RQ("	return " + X.ToS(idx) + "\n")
			RQ("}\n\n")

			// column name functions
			//RQ(`func (` + receiverName + ` *` + structName + ") col" + propName + "() string { //nolint:dupl false positive\n")
			//RQ("	return `" + prop.Name + "`\n")
			//RQ("}\n\n")

			// sql column name functions
			RQ(`func (` + receiverName + ` *` + structName + ") sql" + propName + "() string { //nolint:dupl false positive\n")
			RQ("	return `" + `"` + prop.Name + `"` + "`\n")
			RQ("}\n\n")

			// mutator methods
			WC(`func (` + receiverName + ` *` + structName + "Mutator) Set" + propName + "(val " + typeTranslator[prop.Type] + ") bool { //nolint:dupl false positive\n")
			WC("	if val != " + receiverName + `.` + propName + " {\n")
			WC("		" + receiverName + ".mutations = append(" + receiverName + ".mutations, A.X{`=`, " + I.ToStr(idx) + ", val})\n")
			WC("		" + receiverName + `.` + propName + " = val\n")
			WC("		return true\n")
			WC("	}\n")
			WC("	return false\n")
			WC("}\n\n")
		}

		// to AX
		RQ(`func (` + receiverName + ` *` + structName + ") ToArray() A.X { //nolint:dupl false positive\n")
		RQ("	return A.X{\n")
		for idx, prop := range props.Fields {
			camel := S.CamelCase(prop.Name)
			RQ("		" + receiverName + "." + camel + `,` + strings.Repeat(` `, maxLen-len(camel)) + `// ` + X.ToS(idx) + "\n")
		}
		RQ("	}\n")
		RQ("}\n\n")

		// from AX
		RQ(`func (` + receiverName + ` *` + structName + `) FromArray(a A.X) *` + structName + " { //nolint:dupl false positive\n")
		for idx, prop := range props.Fields {
			RQ("	" + receiverName + "." + S.CamelCase(prop.Name) + ` = ` + typeConverter[prop.Type] + "(a[" + X.ToS(idx) + "])\n")
		}
		RQ("	return " + receiverName + "\n")
		RQ("}\n\n")

		// find many
		//RQ(`func (` + receiverName + ` *` + structName + ") FindOffsetLimit(offset, limit uint32, idx string) []*" + structName + " { //nolint:dupl false positive\n")
		//RQ("	var rows []*" + structName + "\n")
		//RQ("	res, err := " + receiverName + ".Adapter.Select(" + receiverName + ".SpaceName(), idx, offset, limit, " + iterAll + ", A.X{})\n")
		//RQ("	if L.IsError(err, `" + structName + ".FindOffsetLimit failed: `+" + receiverName + ".SpaceName()) {\n")
		//RQ("		return rows\n")
		//RQ("	}\n")
		//RQ("	for _, row := range res.Tuples() {\n")
		//RQ("		item := &" + structName + "{}\n")
		//RQ("		rows = append(rows, item.FromArray(row))\n")
		//RQ("	}\n")
		//RQ("	return rows\n")
		//RQ("}\n\n")

		// total records
		RQ(`func (` + receiverName + ` *` + structName + ") Total() int64 { //nolint:dupl false positive\n")
		RQ("	rows := " + receiverName + ".Adapter.CallBoxSpace(" + receiverName + ".SpaceName() + `:count`, A.X{})\n")
		RQ("	if len(rows) > 0 && len(rows[0]) > 0 {\n")
		RQ("		return X.ToI(rows[0][0])\n")
		RQ("	}\n")
		RQ("	return 0\n")
		RQ("}\n\n")

		//// set to min value
		//WC(`func (`+receiverName+` *` + structName + ") ResetToMax() { //nolint:dupl false positive\n")
		//for _, prop := range props.Fields {
		//	WC("	"+receiverName+"." + S.CamelCase(prop.Name) + " = " + maxMap[prop.Type] + "\n")
		//}
		//WC("}\n\n")
		//
		//// set to min value
		//WC(`func (`+receiverName+` *` + structName + ") ResetToMin() { //nolint:dupl false positive\n")
		//for _, prop := range props.Fields {
		//	WC("	"+receiverName+"." + S.CamelCase(prop.Name) + " = " + minMap[prop.Type] + "\n")
		//}
		//WC("}\n\n")
		//
		//// set if greater
		//WC(`func (`+receiverName+` *` + structName + ") SetIfLesser(l *"+structName+") { //nolint:dupl false positive\n")
		//for _, prop := range props.Fields {
		//	propName := S.CamelCase(prop.Name)
		//	WC("	if "+receiverName+"." + propName + " > l." + propName + " {\n")
		//	WC("		"+receiverName+"." + propName + " = l." + propName + "\n")
		//	WC("	}\n")
		//}
		//WC("}\n\n")
		//
		//// set if greater
		//WC(`func (`+receiverName+` *` + structName + ") SetIfGreater(l *"+structName+") { //nolint:dupl false positive\n")
		//for _, prop := range props.Fields {
		//	propName := S.CamelCase(prop.Name)
		//	WC("	if "+receiverName+"." + propName + " < l." + propName + " {\n")
		//	WC("		"+receiverName+"." + propName + " = l." + propName + "\n")
		//	WC("	}\n")
		//}
		//WC("}\n\n")
		BOTH(warning)

	}

	err := os.MkdirAll(`./`+rqPkgName, os.ModePerm)
	if L.IsError(err, `os.MkDir failed: `+rqPkgName) {
		return
	}
	rqFname := fmt.Sprintf(`./%s/%s__ORM.GEN.go`, rqPkgName, rqPkgName)
	err = os.WriteFile(rqFname, rqBuf.Bytes(), os.ModePerm)
	if L.IsError(err, `os.WriteFile failed: `+rqFname) {
		return
	}

	err = os.MkdirAll(`./`+wcPkgName, os.ModePerm)
	if L.IsError(err, `os.MkDir failed: `+wcPkgName) {
		return
	}
	wcFname := fmt.Sprintf(`./%s/%s__ORM.GEN.go`, wcPkgName, wcPkgName)
	err = os.WriteFile(wcFname, wcBuf.Bytes(), os.ModePerm)
	if L.IsError(err, `os.WriteFile failed: `+wcFname) {
		return
	}
}

func generateMutationByUniqueIndexumns(uniqueCamel, structProp, receiverName, structName string, RQ, WC func(str string)) {

	//// primary fields
	//RQ(`func (` + receiverName + ` *` + structName + ") PrimaryIndex() A.X { //nolint:dupl false positive\n")
	//RQ(`	return A.X{` + structProp + "}\n")
	//RQ("}\n\n")

	// find by unique, used when need to mutate the object
	RQ(`func (` + receiverName + ` *` + structName + `) FindBy` + uniqueCamel + "() bool { //nolint:dupl false positive\n")
	RQ("	res, err := " + receiverName + ".Adapter.Select(" + receiverName + ".SpaceName(), " + receiverName + `.UniqueIndex` + uniqueCamel + `(), 0, 1, ` + iterEq + ", A.X{" + structProp + "})\n")
	RQ("	if L.IsError(err, `" + structName + `.FindBy` + uniqueCamel + " failed: `+" + receiverName + ".SpaceName()) {\n")
	RQ("		return false\n")
	RQ("	}\n")
	RQ("	rows := res.Tuples()\n")
	RQ("	if len(rows) == 1 {\n")
	RQ("		" + receiverName + ".FromArray(rows[0])\n")
	RQ("		return true\n")
	RQ("	}\n")
	RQ("	return false\n")
	RQ("}\n\n")

	// Overwrite all columns, error if not exists
	WC("// Overwrite all columns, error if not exists\n")
	WC(`func (` + receiverName + ` *` + structName + `Mutator) DoOverwriteBy` + uniqueCamel + "() bool { //nolint:dupl false positive\n")
	WC("	_, err := " + receiverName + `.Adapter.Update(` + receiverName + ".SpaceName(), " + receiverName + `.UniqueIndex` + uniqueCamel + "(), A.X{" + structProp + "}, " + receiverName + ".ToUpdateArray())\n")
	WC("	return !L.IsError(err, `" + structName + `.DoOverwriteBy` + uniqueCamel + " failed: `+" + receiverName + ".SpaceName())\n")
	WC("}\n\n")

	// Update only mutated, error if not exists
	WC("// Update only mutated, error if not exists, use Find* and Set* methods instead of direct assignment\n")
	WC(`func (` + receiverName + ` *` + structName + `Mutator) DoUpdateBy` + uniqueCamel + "() bool { //nolint:dupl false positive\n")
	WC(`	if !` + receiverName + ".HaveMutation() {\n")
	WC("		return true\n")
	WC("	}\n")
	WC("	_, err := " + receiverName + `.Adapter.Update(` + receiverName + ".SpaceName(), " + receiverName + `.UniqueIndex` + uniqueCamel + "(), A.X{" + structProp + "}, " + receiverName + ".mutations)\n")
	WC("	return !L.IsError(err, `" + structName + `.DoUpdateBy` + uniqueCamel + " failed: `+" + receiverName + ".SpaceName())\n")
	WC("}\n\n")

	// permanent delete
	WC(`func (` + receiverName + ` *` + structName + `Mutator) DoDeletePermanentBy` + uniqueCamel + "() bool { //nolint:dupl false positive\n")
	WC("	_, err := " + receiverName + ".Adapter.Delete(" + receiverName + ".SpaceName(), " + receiverName + `.UniqueIndex` + uniqueCamel + `(), A.X{` + structProp + "})\n")
	WC("	return !L.IsError(err, `" + structName + `.DoDeletePermanentBy` + uniqueCamel + " failed: `+" + receiverName + ".SpaceName())\n")
	WC("}\n\n")
}
