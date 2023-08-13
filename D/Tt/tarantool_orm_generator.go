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
	` + S.BT(importPath)
}

// double quote
func dq(str string) string {
	return `"` + str + `"`
}

var typeGraphql = map[DataType]string{
	Unsigned: `Int`,
	//Number:   `Float`,
	String:  `String`,
	Double:  `Float`,
	Integer: `Int`,
	Boolean: `Boolean`,
}

func TypeGraphql(field Field) string {
	typ := typeGraphql[field.Type]
	if typ == `Int` {
		if S.EndsWith(field.Name, `By`) || S.EndsWith(field.Name, `Id`) {
			return `ID`
		}
	}
	return typ
}

const connStruct = `Tt.Adapter`
const connImport = "\n\n\t`github.com/tarantool/go-tarantool`"
const iterEq = `tarantool.IterEq`
const iterAll = `tarantool.IterAll`
const iterNeighbor = `tarantool.IterNeighbor`

var _ = iterNeighbor

//const iterAll = `tarantool.IterAll`

const NL = "\n"

const warning = "// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go\n\n"

func GenerateOrm(tables map[TableName]*TableProp, withGraphql ...bool) {
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

	//BOTH(`// generated: ` + time.Now().String() + NL)
	RQ(`package ` + rqPkgName)
	WC(`package ` + wcPkgName)
	BOTH("\n\n")
	BOTH(warning)

	useGraphql := len(withGraphql) > 0
	// sort by table name to keep the order when regenerating structs
	tableNames := make([]string, 0, len(tables))
	for k, v := range tables {
		tableNames = append(tableNames, string(k))
		useGraphql = useGraphql || v.GenGraphqlType // if one of them use graphql, import anyway
	}
	sort.Strings(tableNames)

	// import reader
	BOTH(`import (`)

	RQ(qi(ci.PackageName))
	RQ(connImport)
	WC(qi(ci.PackageName + `/` + rqPkgName))

	BOTH(NL)
	if useGraphql {
		//RQ(qi(`github.com/   graphql-go/graphql`))
	}
	BOTH(qi(`github.com/kokizzu/gotro/A`))
	BOTH(qi(this.PackageName)) // github.com/kokizzu/gotro/D/Tt
	BOTH(qi(`github.com/kokizzu/gotro/L`))
	WC(qi(`github.com/kokizzu/gotro/M`))
	BOTH(qi(`github.com/kokizzu/gotro/X`))

	BOTH(`
)` + "\n\n")

	RQ(`//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file ` + rqPkgName + "__ORM.GEN.go\n")
	RQ(`//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type ` + rqPkgName + "__ORM.GEN.go\n")
	RQ(`//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type ` + rqPkgName + "__ORM.GEN.go\n")
	RQ(`//go:generate replacer -afterprefix "By\" form" "By,string\" form" type ` + rqPkgName + "__ORM.GEN.go\n")
	//RQ(`//go:generate msgp -tests=false -file ` + rqPkgName + `__ORM.GEN.go -o ` + rqPkgName + `__MSG.GEN.go` + "\n\n")

	WC(`//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file ` + wcPkgName + "__ORM.GEN.go\n")
	WC(`//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type ` + wcPkgName + "__ORM.GEN.go\n")
	WC(`//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type ` + wcPkgName + "__ORM.GEN.go\n")
	WC(`//go:generate replacer -afterprefix "By\" form" "By,string\" form" type ` + wcPkgName + "__ORM.GEN.go\n")
	//WC(` //go:generate msgp -tests=false -file ` + wcPkgName + `__ORM.GEN.go -o ` + wcPkgName + `__MSG.GEN.go` + "\n\n")

	// for each table generate in order
	for _, tableName := range tableNames {
		props := tables[TableName(tableName)]
		structName := S.PascalCase(tableName)
		maxLen := 1
		propTypeByName := map[string]Field{}
		censoredFieldsByName := map[string]bool{}
		for _, prop := range props.Fields {
			l := len(prop.Name) + 1 - strings.Count(prop.Name, `_`)
			if maxLen < l {
				maxLen = l
			}
			propTypeByName[prop.Name] = prop
		}
		for _, propName := range props.AutoCensorFields {
			censoredFieldsByName[propName] = true
		}

		// mutator struct
		WC("// " + structName + "Mutator DAO writer/command struct\n")
		WC(`type ` + structName + "Mutator struct {\n")
		WC(`	` + rqPkgName + `.` + structName + NL)
		WC("	mutations []A.X\n")
		WC("	logs      []A.X\n")
		WC("}\n\n")

		// mutator struct constructor
		WC("// New" + structName + "Mutator create new ORM writer/command object\n")
		WC(`func New` + structName + `Mutator(adapter *` + connStruct + `) *` + structName + "Mutator {\n")
		WC(`	return &` + structName + `Mutator{` + structName + `: ` + rqPkgName + `.` + structName + "{Adapter: adapter}}\n")
		WC("}\n\n")

		// reader struct
		RQ("// " + structName + " DAO reader/query struct\n")
		RQ(`type ` + structName + " struct {\n")
		const none = `"-"`
		RQ("	Adapter *" + connStruct + " " + S.BT("json:"+none+" msg:"+none+" query:"+none+" form:"+none) + NL)
		for _, prop := range props.Fields {
			camel := S.PascalCase(prop.Name)
			RQ("	" + camel + strings.Repeat(` `, maxLen-len(camel)) + TypeToGoType[prop.Type] + NL)
		}
		RQ("}\n\n")

		// reader struct constructor
		RQ("// New" + structName + " create new ORM reader/query object\n")
		RQ(`func New` + structName + `(adapter *` + connStruct + `) *` + structName + " {\n")
		RQ(`	return &` + structName + "{Adapter: adapter}\n")
		RQ("}\n\n")

		// table name
		receiverName := strings.ToLower(string(structName[0]))
		RQ("// SpaceName returns full package and table name\n")
		RQ(`func (` + receiverName + ` *` + structName + ") SpaceName() string { //nolint:dupl false positive\n")
		RQ("	return string(" + mPkgName + `.Table` + structName + ") // casting required to string from Tt.TableName\n")
		RQ("}\n\n")

		// Sql table name
		RQ("// SqlTableName returns quoted table name\n")
		RQ(`func (` + receiverName + ` *` + structName + ") SqlTableName() string { //nolint:dupl false positive\n")
		RQ("	return " + S.BT(S.QQ(tableName)) + NL)
		RQ("}\n\n")

		// have mutation
		WC("// Logs get array of logs [field, old, new]\n")
		WC(`func (` + receiverName + ` *` + structName + "Mutator) Logs() []A.X { //nolint:dupl false positive\n")
		WC(`	return ` + receiverName + ".logs\n")
		WC("}\n\n")

		// have mutation
		WC("// HaveMutation check whether Set* methods ever called\n")
		WC(`func (` + receiverName + ` *` + structName + "Mutator) HaveMutation() bool { //nolint:dupl false positive\n")
		WC(`	return len(` + receiverName + ".mutations) > 0\n")
		WC("}\n\n")

		// clear mutation
		WC("// ClearMutations clear all previously called Set* methods\n")
		WC(`func (` + receiverName + ` *` + structName + "Mutator) ClearMutations() { //nolint:dupl false positive\n")
		WC(`	` + receiverName + ".mutations = []A.X{}\n")
		WC(`	` + receiverName + ".logs = []A.X{}\n")
		WC("}\n\n")

		// auto increment id
		if props.AutoIncrementId {
			uniquePropCamel := S.PascalCase(IdCol)
			structProp := receiverName + `.` + uniquePropCamel

			RQ(`func (` + receiverName + ` *` + structName + `) UniqueIndex` + uniquePropCamel + "() string { //nolint:dupl false positive\n")
			RQ("	return " + S.BT(IdCol) + NL)
			RQ("}\n\n")

			generateMutationByUniqueIndexumns(uniquePropCamel, structProp, receiverName, structName, RQ, WC)

			if props.GenGraphqlType {
				generateGraphqlQueryField(structName, uniquePropCamel, propTypeByName[IdCol], RQ)
			}
		}

		// upsert template, to be copied when need increment some field
		WC(`// func (` + receiverName + ` *` + structName + "Mutator) DoUpsert() bool { //nolint:dupl false positive\n")
		WC("//	_, err := " + receiverName + ".Adapter.Upsert(" + receiverName + ".SpaceName(), " + receiverName + ".ToArray(), A.X{\n")
		for idx, prop := range props.Fields {
			WC("//		A.X{`=`, " + X.ToS(idx) + ", " + receiverName + "." + S.PascalCase(prop.Name) + "},\n")
		}
		WC("//	})\n")
		WC("//	return !L.IsError(err, `" + structName + ".DoUpsert failed: `+" + receiverName + ".SpaceName())\n")
		WC("// }\n\n")

		// spatial index
		if props.Spatial != `` {
			uniquePropCamel := S.PascalCase(props.Spatial)

			RQ("// SpatialIndex" + uniquePropCamel + " return spatial index name\n")
			RQ(`func (` + receiverName + ` *` + structName + `) SpatialIndex` + uniquePropCamel + "() string { //nolint:dupl false positive\n")
			RQ("	return " + S.BT(props.Spatial) + NL)
			RQ("}\n\n")
		}

		// unique index1
		if props.Unique1 != `` && !(props.AutoIncrementId && props.Unique1 == IdCol) {
			uniquePropCamel := S.PascalCase(props.Unique1)
			structProp := receiverName + `.` + uniquePropCamel

			RQ("// UniqueIndex" + uniquePropCamel + " return unique index name\n")
			RQ(`func (` + receiverName + ` *` + structName + `) UniqueIndex` + uniquePropCamel + "() string { //nolint:dupl false positive\n")
			RQ("	return " + S.BT(props.Unique1) + NL)
			RQ("}\n\n")

			generateMutationByUniqueIndexumns(uniquePropCamel, structProp, receiverName, structName, RQ, WC)

			if props.GenGraphqlType {
				generateGraphqlQueryField(structName, uniquePropCamel, propTypeByName[props.Unique1], RQ)
			}
		}

		// unique index2
		if props.Unique2 != `` && !(props.AutoIncrementId && props.Unique2 == IdCol) {
			uniquePropCamel := S.PascalCase(props.Unique2)
			structProp := receiverName + `.` + uniquePropCamel

			RQ("// UniqueIndex" + uniquePropCamel + " return unique index name\n")
			RQ(`func (` + receiverName + ` *` + structName + `) UniqueIndex` + uniquePropCamel + "() string { //nolint:dupl false positive\n")
			RQ("	return " + S.BT(props.Unique2) + NL)
			RQ("}\n\n")

			generateMutationByUniqueIndexumns(uniquePropCamel, structProp, receiverName, structName, RQ, WC)
		}

		// unique index3
		if props.Unique3 != `` && !(props.AutoIncrementId && props.Unique3 == IdCol) {
			uniquePropCamel := S.PascalCase(props.Unique3)
			structProp := receiverName + `.` + uniquePropCamel

			RQ("// UniqueIndex" + uniquePropCamel + " return unique index name\n")
			RQ(`func (` + receiverName + ` *` + structName + `) UniqueIndex` + uniquePropCamel + "() string { //nolint:dupl false positive\n")
			RQ("	return " + S.BT(props.Unique3) + NL)
			RQ("}\n\n")

			generateMutationByUniqueIndexumns(uniquePropCamel, structProp, receiverName, structName, RQ, WC)
		}

		// unique indexes
		if len(props.Uniques) > 0 {
			uniquePropCamel := ``
			structProps := ``
			for _, uniq := range props.Uniques {
				uniquePropCamel += S.PascalCase(uniq)
				structProps += `, ` + receiverName + `.` + S.PascalCase(uniq)
			}
			if len(structProps) > 2 {
				structProps = structProps[2:]
			}

			RQ("// UniqueIndex" + uniquePropCamel + " return unique index name\n")
			RQ(`func (` + receiverName + ` *` + structName + `) UniqueIndex` + uniquePropCamel + "() string { //nolint:dupl false positive\n")
			RQ("	return " + S.BT(strings.Join(props.Uniques, `__`)) + NL)
			RQ("}\n\n")

			generateMutationByUniqueIndexumns(uniquePropCamel, structProps, receiverName, structName, RQ, WC)
		}

		// insert, error if exists
		WC("// DoInsert insert, error if already exists\n")
		WC(`func (` + receiverName + ` *` + structName + "Mutator) DoInsert() bool { //nolint:dupl false positive\n")
		ret1 := S.IfElse(props.AutoIncrementId, `row`, `_`)
		WC("	" + ret1 + ", err := " + receiverName + ".Adapter.Insert(" + receiverName + ".SpaceName(), " + receiverName + ".ToArray())\n")
		if props.AutoIncrementId {
			WC("	if err == nil {\n")
			WC("		tup := row.Tuples()\n")
			WC("		if len(tup) > 0 && len(tup[0]) > 0 && tup[0][0] != nil {\n")
			WC("			" + receiverName + ".Id = X.ToU(tup[0][0])\n")
			WC("		}\n")
			WC("	}\n")
		}
		WC("	return !L.IsError(err, `" + structName + ".DoInsert failed: `+" + receiverName + ".SpaceName())\n")
		WC("}\n\n")

		// replace = upsert, only error when there's unique secondary key
		WC("// DoUpsert upsert, insert or overwrite, will error only when there's unique secondary key being violated\n")
		WC("// replace = upsert, only error when there's unique secondary key\n")
		WC("// previous name: DoReplace\n")
		WC(`func (` + receiverName + ` *` + structName + "Mutator) DoUpsert() bool { //nolint:dupl false positive\n")
		WC("	_, err := " + receiverName + ".Adapter.Replace(" + receiverName + ".SpaceName(), " + receiverName + ".ToArray())\n")
		WC("	return !L.IsError(err, `" + structName + ".DoUpsert failed: `+" + receiverName + ".SpaceName())\n")
		WC("}\n\n")

		// Sql select all fields, used when need to mutate or show every fields
		RQ("// SqlSelectAllFields generate Sql select fields\n")
		RQ(`func (` + receiverName + ` *` + structName + ") SqlSelectAllFields() string { //nolint:dupl false positive\n")
		sqlFields := ``
		for _, prop := range props.Fields {
			sqlFields += `, ` + dq(prop.Name) + "\n\t"
		}
		RQ(`	return ` + S.BT(sqlFields[1:]) + NL)
		RQ("}\n\n")

		// Sql select all fields, used when need to mutate or show only uncensored fields
		RQ("// SqlSelectAllUncensoredFields generate Sql select fields\n")
		RQ(`func (` + receiverName + ` *` + structName + ") SqlSelectAllUncensoredFields() string { //nolint:dupl false positive\n")
		sqlUncenFields := ``
		for _, prop := range props.Fields {
			if !censoredFieldsByName[prop.Name] {
				sqlUncenFields += `, ` + dq(prop.Name) + "\n\t"
			}
		}
		RQ(`	return ` + S.BT(sqlFields[1:]) + NL)
		RQ("}\n\n")

		// to Update AX
		RQ("// ToUpdateArray generate slice of update command\n")
		RQ(`func (` + receiverName + ` *` + structName + ") ToUpdateArray() A.X { //nolint:dupl false positive\n")
		RQ("	return A.X{\n")
		for idx, prop := range props.Fields {
			RQ("		A.X{`=`, " + X.ToS(idx) + ", " + receiverName + "." + S.PascalCase(prop.Name) + "},\n")
		}
		RQ("	}\n")
		RQ("}\n\n")

		for idx, prop := range props.Fields {
			propName := S.PascalCase(prop.Name)

			// index functions
			RQ("// Idx" + propName + " return name of the index\n")
			RQ(`func (` + receiverName + ` *` + structName + ") Idx" + propName + "() int { //nolint:dupl false positive\n")
			RQ("	return " + X.ToS(idx) + NL)
			RQ("}\n\n")

			// column name functions
			//RQ(`func (` + receiverName + ` *` + structName + ") col" + propName + "() string { //nolint:dupl false positive\n")
			//RQ("	return " + S.BT(prop.Name) + NL)
			//RQ("}\n\n")

			// Sql column name functions
			RQ("// Sql" + propName + " return name of the column being indexed\n")
			RQ(`func (` + receiverName + ` *` + structName + ") Sql" + propName + "() string { //nolint:dupl false positive\n")
			RQ("	return " + S.BT(S.QQ(prop.Name)) + NL)
			RQ("}\n\n")

			// mutator methods
			WC("// Set" + propName + " create mutations, should not duplicate\n")
			propType := TypeToGoType[prop.Type]
			WC(`func (` + receiverName + ` *` + structName + "Mutator) Set" + propName + "(val " + propType + ") bool { //nolint:dupl false positive\n")
			if prop.Type != Array {
				WC("	if val != " + receiverName + `.` + propName + " {\n")
				WC("		" + receiverName + ".mutations = append(" + receiverName + ".mutations, A.X{`=`, " + I.ToStr(idx) + ", val})\n")
				if !censoredFieldsByName[prop.Name] {
					WC("		" + receiverName + ".logs = append(" + receiverName + ".logs, A.X{`" + prop.Name + "`, " + receiverName + `.` + propName + ", val})\n")
				}
				WC("		" + receiverName + `.` + propName + " = val\n")
				WC("		return true\n")
				WC("	}\n")
				WC("	return false\n")
			} else { // always overwrite for array
				WC("	" + receiverName + ".mutations = append(" + receiverName + ".mutations, A.X{`=`, " + I.ToStr(idx) + ", val})\n")
				WC("	" + receiverName + ".logs = append(" + receiverName + ".logs, A.X{`" + prop.Name + "`, " + receiverName + `.` + propName + ", val})\n")
				if !censoredFieldsByName[prop.Name] {
					WC("	" + receiverName + `.` + propName + " = val\n")
				}
				WC("	return true\n")
			}
			WC("}\n\n")
		}

		// SetAll
		WC("// SetAll set all from another source, only if another property is not empty/nil/zero or in forceMap\n")
		WC(`func (` + receiverName + ` *` + structName + "Mutator) SetAll(in " + rqPkgName + `.` + structName + ", excludeMap, forceMap M.SB) (changed bool) { //nolint:dupl false positive\n")
		WC("	if excludeMap == nil { // list of fields to exclude\n")
		WC("		excludeMap = M.SB{}\n")
		WC("	}\n")
		WC("	if forceMap == nil { // list of fields to force overwrite\n")
		WC("		forceMap = M.SB{}\n")
		WC("	}\n")
		for _, prop := range props.Fields {
			propName := S.PascalCase(prop.Name)

			// index functions
			WC("	if !excludeMap[`" + prop.Name + "`] && (forceMap[`" + prop.Name + "`] || " + receiverName + `.` + propName + ` != ` + TypeToGoNilValue[prop.Type] + ") {\n")
			WC(`		` + receiverName + `.` + propName + ` = in.` + propName + "\n")
			WC("		changed = true\n")
			WC("	}\n")
		}
		WC("	return\n")
		WC("}\n\n")

		// CensorFields
		if len(props.AutoCensorFields) > 0 {
			RQ("// CensorFields remove sensitive fields for output\n")
			RQ(`func (` + receiverName + ` *` + structName + ") CensorFields() { //nolint:dupl false positive\n")
			for _, propName := range props.AutoCensorFields {
				propType := propTypeByName[propName].Type
				RQ("	" + receiverName + "." + S.PascalCase(propName) + " = " + TypeToGoEmptyValue[propType] + "\n")
			}
			RQ("	}\n")
		}

		// to AX
		RQ("// ToArray receiver fields to slice\n")
		RQ(`func (` + receiverName + ` *` + structName + ") ToArray() A.X { //nolint:dupl false positive\n")
		if props.AutoIncrementId {
			RQ("	var " + IdCol + " any = nil\n")
			idProp := receiverName + "." + S.PascalCase(IdCol)
			RQ("	if " + idProp + " != 0 {\n")
			RQ("		" + IdCol + " = " + idProp + "\n")
			RQ("	}\n")
		}
		RQ("	return A.X{\n")
		for idx, prop := range props.Fields {
			camel := S.PascalCase(prop.Name)
			if props.AutoIncrementId && IdCol == prop.Name {
				RQ("		" + IdCol + ",\n")
			} else {
				RQ("		" + receiverName + "." + camel + `,` + strings.Repeat(` `, maxLen-len(camel)) + `// ` + X.ToS(idx) + NL)
			}
		}
		RQ("	}\n")
		RQ("}\n\n")

		// from AX
		RQ("// FromArray convert slice to receiver fields\n")
		RQ(`func (` + receiverName + ` *` + structName + `) FromArray(a A.X) *` + structName + " { //nolint:dupl false positive\n")
		for idx, prop := range props.Fields {
			RQ("	" + receiverName + "." + S.PascalCase(prop.Name) + ` = ` + TypeToConvertFunc[prop.Type] + "(a[" + X.ToS(idx) + "])\n")
		}
		RQ("	return " + receiverName + NL)
		RQ("}\n\n")

		// from AX but uncensored
		RQ("// FromUncensoredArray convert slice to receiver fields\n")
		RQ(`func (` + receiverName + ` *` + structName + `) FromUncensoredArray(a A.X) *` + structName + " { //nolint:dupl false positive\n")
		for idx, prop := range props.Fields {
			if !censoredFieldsByName[prop.Name] {
				RQ("	" + receiverName + "." + S.PascalCase(prop.Name) + ` = ` + TypeToConvertFunc[prop.Type] + "(a[" + X.ToS(idx) + "])\n")
			}
		}
		RQ("	return " + receiverName + NL)
		RQ("}\n\n")

		// find many
		RQ("// FindOffsetLimit returns slice of struct, order by idx, eg. .UniqueIndex*()\n")
		RQ(`func (` + receiverName + ` *` + structName + ") FindOffsetLimit(offset, limit uint32, idx string) []" + structName + " { //nolint:dupl false positive\n")
		RQ("	var rows []" + structName + NL)
		RQ("	res, err := " + receiverName + ".Adapter.Select(" + receiverName + ".SpaceName(), idx, offset, limit, " + iterAll + ", A.X{})\n")
		RQ("	if L.IsError(err, `" + structName + ".FindOffsetLimit failed: `+" + receiverName + ".SpaceName()) {\n")
		RQ("		return rows\n")
		RQ("	}\n")
		RQ("	for _, row := range res.Tuples() {\n")
		RQ("		item := " + structName + "{}\n")
		RQ("		rows = append(rows, *item.FromArray(row))\n")
		RQ("	}\n")
		RQ("	return rows\n")
		RQ("}\n\n")
		// find many

		RQ("// FindArrOffsetLimit returns as slice of slice order by idx eg. .UniqueIndex*()\n")
		RQ(`func (` + receiverName + ` *` + structName + ") FindArrOffsetLimit(offset, limit uint32, idx string) ([]A.X, Tt.QueryMeta) { //nolint:dupl false positive\n")
		RQ("	var rows []A.X" + NL)
		RQ("	res, err := " + receiverName + ".Adapter.Select(" + receiverName + ".SpaceName(), idx, offset, limit, " + iterAll + ", A.X{})\n")
		RQ("	if L.IsError(err, `" + structName + ".FindOffsetLimit failed: `+" + receiverName + ".SpaceName()) {\n")
		RQ("		return rows, Tt.QueryMetaFrom(res, err)\n")
		RQ("	}\n")
		RQ("	tuples := res.Tuples()\n")
		RQ("	rows = make([]A.X, len(tuples))\n")
		RQ("	for z, row := range tuples {\n")
		RQ("		rows[z] = row\n")
		RQ("	}\n")
		RQ("	return rows, Tt.QueryMetaFrom(res, nil)\n")
		RQ("}\n\n")

		// total records
		RQ("// Total count number of rows\n")
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
		//	WC("	"+receiverName+"." + S.PascalCase(prop.Name) + " = " + maxMap[prop.Type] + NL)
		//}
		//WC("}\n\n")
		//
		//// set to min value
		//WC(`func (`+receiverName+` *` + structName + ") ResetToMin() { //nolint:dupl false positive\n")
		//for _, prop := range props.Fields {
		//	WC("	"+receiverName+"." + S.PascalCase(prop.Name) + " = " + minMap[prop.Type] + NL)
		//}
		//WC("}\n\n")
		//
		//// set if greater
		//WC(`func (`+receiverName+` *` + structName + ") SetIfLesser(l *"+structName+") { //nolint:dupl false positive\n")
		//for _, prop := range props.Fields {
		//	propName := S.PascalCase(prop.Name)
		//	WC("	if "+receiverName+"." + propName + " > l." + propName + " {\n")
		//	WC("		"+receiverName+"." + propName + " = l." + propName + NL)
		//	WC("	}\n")
		//}
		//WC("}\n\n")
		//
		//// set if greater
		//WC(`func (`+receiverName+` *` + structName + ") SetIfGreater(l *"+structName+") { //nolint:dupl false positive\n")
		//for _, prop := range props.Fields {
		//	propName := S.PascalCase(prop.Name)
		//	WC("	if "+receiverName+"." + propName + " < l." + propName + " {\n")
		//	WC("		"+receiverName+"." + propName + " = l." + propName + NL)
		//	WC("	}\n")
		//}
		//WC("}\n\n")

		// field type map
		RQ("// " + structName + "FieldTypeMap returns key value of field name and key\n")
		RQ("var " + structName + "FieldTypeMap = map[string]Tt.DataType { //nolint:dupl false positive\n")
		for _, field := range props.Fields {
			RQ("	" + S.BT(field.Name) + `:` + strings.Repeat(` `, maxLen-len(field.Name)) + TypeToConst[field.Type] + ",\n")
		}
		RQ("}\n\n")

		// graphql type
		if props.GenGraphqlType {
			RQ(`var GraphqlType` + structName + " = graphql.NewObject(\n")
			RQ("	graphql.ObjectConfig{\n")
			RQ("		Name: " + S.BT(tableName) + ",\n")
			RQ("		Fields: graphql.Fields{\n")

			hiddenFields := map[string]bool{}
			for _, fieldName := range props.HiddenFields {
				hiddenFields[fieldName] = true
			}
			for _, field := range props.Fields {
				if hiddenFields[field.Name] {
					continue
				}
				RQ("			" + S.BT(field.Name) + ": &graphql.Field{\n")
				RQ("				Type: graphql." + TypeGraphql(field) + ",\n")
				RQ("			},\n")
			}

			RQ("		},\n")
			RQ("	},\n")
			RQ(")\n\n")

			//// graphql field list
			//RQ(`var GraphqlField` + structName + "List = &graphql.Field{\n")
			//RQ("	Type: GraphqlType" + structName + ",\n")
			//RQ("	Description: " + S.BT(`list of `+structName) + ",\n")
			//
			//RQ("}\n\n")
		}

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

func generateGraphqlQueryField(structName string, uniqueFieldName string, field Field, RQ func(str string)) {

	// graphql field
	RQ(`var GraphqlField` + structName + "By" + uniqueFieldName + " = &graphql.Field{\n")
	RQ("	Type: GraphqlType" + structName + ",\n")
	RQ("	Description: " + S.BT(`list of `+structName) + ",\n")
	RQ("	Args: graphql.FieldConfigArgument{\n")
	RQ("		" + S.BT(uniqueFieldName) + ": &graphql.ArgumentConfig{\n")
	RQ("			Type: graphql." + TypeGraphql(field) + ",\n")
	RQ("		},\n")
	RQ("	},\n")
	RQ("}\n\n")

	// graphql field resolver
	RQ(`func (g *` + structName + `) GraphqlField` + structName + "By" + uniqueFieldName + "WithResolver() *graphql.Field {\n")
	RQ("	field := *GraphqlField" + structName + "By" + uniqueFieldName + "\n")
	RQ("	field.Resolve = func(p graphql.ResolveParams) (any, error) {\n")
	RQ("		q := g\n")
	RQ("		v, ok := p.Args[" + S.BT(S.LowerFirst(uniqueFieldName)) + "]\n")
	RQ("		if !ok {\n")
	RQ("			v, _ = p.Args[" + S.BT(uniqueFieldName) + "]\n")
	RQ("		}\n")
	RQ("		q." + uniqueFieldName + " = " + TypeToConvertFunc[field.Type] + "(v)\n")
	RQ("		if q.FindBy" + uniqueFieldName + "() {\n")
	RQ("			return q, nil\n")
	RQ("		}\n")
	RQ("		return nil, nil\n")
	RQ("	}\n")
	RQ("	return &field\n")
	RQ("}\n\n")

}

func generateMutationByUniqueIndexumns(uniqueCamel, structProp, receiverName, structName string, RQ, WC func(str string)) {

	//// primary fields
	//RQ(`func (` + receiverName + ` *` + structName + ") PrimaryIndex() A.X { //nolint:dupl false positive\n")
	//RQ(`	return A.X{` + structProp + "}\n")
	//RQ("}\n\n")

	// find by unique, used when need to mutate the object

	RQ("// FindBy" + uniqueCamel + " Find one by " + uniqueCamel + "\n")
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
	WC("// DoOverwriteBy" + uniqueCamel + " update all columns, error if not exists, not using mutations/Set*\n")
	WC(`func (` + receiverName + ` *` + structName + `Mutator) DoOverwriteBy` + uniqueCamel + "() bool { //nolint:dupl false positive\n")
	WC("	_, err := " + receiverName + `.Adapter.Update(` + receiverName + ".SpaceName(), " + receiverName + `.UniqueIndex` + uniqueCamel + "(), A.X{" + structProp + "}, " + receiverName + ".ToUpdateArray())\n")
	WC("	return !L.IsError(err, `" + structName + `.DoOverwriteBy` + uniqueCamel + " failed: `+" + receiverName + ".SpaceName())\n")
	WC("}\n\n")

	// Update only mutated, error if not exists
	WC("// DoUpdateBy" + uniqueCamel + " update only mutated fields, error if not exists, use Find* and Set* methods instead of direct assignment\n")
	WC(`func (` + receiverName + ` *` + structName + `Mutator) DoUpdateBy` + uniqueCamel + "() bool { //nolint:dupl false positive\n")
	WC(`	if !` + receiverName + ".HaveMutation() {\n")
	WC("		return true\n")
	WC("	}\n")
	WC("	_, err := " + receiverName + `.Adapter.Update(` + receiverName + ".SpaceName(), " + receiverName + `.UniqueIndex` + uniqueCamel + "(), A.X{" + structProp + "}, " + receiverName + ".mutations)\n")
	WC("	return !L.IsError(err, `" + structName + `.DoUpdateBy` + uniqueCamel + " failed: `+" + receiverName + ".SpaceName())\n")
	WC("}\n\n")

	// permanent delete
	WC("// DoDeletePermanentBy" + uniqueCamel + " permanent delete\n")
	WC(`func (` + receiverName + ` *` + structName + `Mutator) DoDeletePermanentBy` + uniqueCamel + "() bool { //nolint:dupl false positive\n")
	WC("	_, err := " + receiverName + ".Adapter.Delete(" + receiverName + ".SpaceName(), " + receiverName + `.UniqueIndex` + uniqueCamel + `(), A.X{` + structProp + "})\n")
	WC("	return !L.IsError(err, `" + structName + `.DoDeletePermanentBy` + uniqueCamel + " failed: `+" + receiverName + ".SpaceName())\n")
	WC("}\n\n")

}
