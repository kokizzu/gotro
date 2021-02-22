package Tr

import (
	"bytes"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// this function generate struct of each table properties
// genUser = whoami
// genDir = `/path/to/dir/`
// genFileName = `tarantool_structs.go`
// genPkgName = `model`
func GenerateStruct(genUser, genDir, genFileName, genPkgName string, tables []TableProp) bool {
	var typeTranslator = map[string]string{
		Unsigned: `int64`,
		Number:   `float64`,
		String:   String,
	}
	var typeConverter = map[string]string{
		Unsigned: `conv.ToInt64`,
		Number:   `conv.ToFloat64`,
		String:   `conv.ToStr`,
	}
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
	if usr, err := user.Current(); err == nil {
		if usr.Username == genUser {
			// get file stats
			maxModTime := int64(0)
			stat, err := os.Stat(genDir + genFileName)
			if err == nil {
				err = filepath.Walk(genDir, func(path string, info os.FileInfo, err error) error {
					if strings.Contains(path, `tarantool_table_`) ||
						path == `tarantool_schema.go` {
						modTime := info.ModTime().UnixNano()
						if maxModTime < modTime {
							maxModTime = modTime
						}
					}
					return nil
				})
				if L.IsError(err, `walking `) {
					return false
				}
				// no table file changed
				if stat.ModTime().UnixNano() >= maxModTime {
					return false
				}
			}

			// generate
			lines := bytes.Buffer{}
			warning := `// DO NOT EDIT, will be overwritten by ` + genFileName + "\n"
			lines.WriteString(warning)
			lines.WriteString(`// generated: ` + time.Now().String() + "\npackage " + genPkgName + "\n\n")
			lines.WriteString(`import (
	"bitbucket.org/cyzainc/goc/pkg/conv"
	"bitbucket.org/cyzainc/goc/pkg/tracer"
	"github.com/tarantool/go-tarantool"
)` + "\n\n")
			// sort by table name to keep the order when regenerating structs
			var tableNames []string
			mapNameTable := map[string]TableProp{}
			for _, table := range tables {
				tableNames = append(tableNames, table.SpaceName)
				mapNameTable[table.SpaceName] = table
			}
			sort.Strings(tableNames)

			for _, tableName := range tableNames {
				props := mapNameTable[tableName]
				structName := S.CamelCase(tableName)
				maxLen := 0
				for _, prop := range props.Fields {
					if maxLen < len(prop.Name) {
						maxLen = len(prop.Name)
					}
				}
				lines.WriteString(`type ` + structName + " struct {\n")
				for _, prop := range props.Fields {
					camel := S.CamelCase(prop.Name)
					lines.WriteString("	" + camel + strings.Repeat(` `, maxLen-len(camel)) + typeTranslator[prop.Type] + "\n")
				}
				lines.WriteString("}\n\n")

				// table name
				lines.WriteString(`func (x *` + structName + ") SpaceName() string {\n")
				lines.WriteString("	return `" + tableName + "`\n")
				lines.WriteString("}\n\n")

				// primary key
				lines.WriteString(`func (x *` + structName + ") PrimaryIndex() string {\n")
				if props.Unique != `` {
					lines.WriteString("	return `" + props.Unique + "`\n")
				} else {
					lines.WriteString("	return `" + strings.Join(props.Uniques, `__`) + "`\n")
				}
				lines.WriteString("}\n\n")

				uniqueCamel := S.CamelCase(props.Unique)
				if props.Unique == `` {
					uniques := ``
					for _, uniq := range props.Uniques {
						uniques += `, x.` + S.CamelCase(uniq)
					}
					uniqueCamel = uniques[1:]
				} else {
					uniqueCamel = `x.` + uniqueCamel
				}

				// primary fields
				lines.WriteString(`func (x *` + structName + ") PrimaryFields() AX {\n")
				lines.WriteString(`	return AX{` + uniqueCamel + "}\n")
				lines.WriteString("}\n\n")

				// insert, error if exists
				lines.WriteString(`func (x *` + structName + ") Insert(a *Adapter) bool {\n")
				lines.WriteString("	_, err := a.conn.Insert(x.SpaceName(), x.ToArray())\n")
				lines.WriteString("	return !tracer.IsError(err)\n")
				lines.WriteString("}\n\n")

				// update, error if not exists
				lines.WriteString(`func (x *` + structName + ") Update(a *Adapter) bool {\n")
				lines.WriteString("	_, err := a.conn.Update(x.SpaceName(), x.PrimaryIndex(), AX{" + uniqueCamel + "}, x.ToUpdateArray())\n")
				lines.WriteString("	return !tracer.IsError(err)\n")
				lines.WriteString("}\n\n")

				// delete
				lines.WriteString(`func (x *` + structName + ") Delete(a *Adapter) bool {\n")
				lines.WriteString("	_, err := a.conn.Delete(x.SpaceName(), x.PrimaryIndex(), AX{" + uniqueCamel + "})\n")
				lines.WriteString("	return !tracer.IsError(err)\n")
				lines.WriteString("}\n\n")

				// replace = upsert, only error when there's unique secondary key
				lines.WriteString(`func (x *` + structName + ") Upsert(a *Adapter) bool {\n")
				lines.WriteString("	_, err := a.conn.Replace(x.SpaceName(), x.ToArray())\n")
				lines.WriteString("	return !tracer.IsError(err)\n")
				lines.WriteString("}\n\n")

				// upsert template, to be copied when need increment some field
				lines.WriteString(`// func (x *` + structName + ") Upsert(a *Adapter) bool {\n")
				lines.WriteString("//	_, err := a.conn.Upsert(x.SpaceName(), x.ToArray(), AX{\n")
				for idx, prop := range props.Fields {
					lines.WriteString("//		AX{`=`, " + I.ToStr(idx) + ", x." + S.CamelCase(prop.Name) + "},\n")
				}
				lines.WriteString("//	})\n")
				lines.WriteString("//	return !tracer.IsError(err)\n")
				lines.WriteString("// }\n\n")

				// to Update AX
				lines.WriteString(`func (x *` + structName + ") ToUpdateArray() AX {\n")
				lines.WriteString("	return AX{\n")
				for idx, prop := range props.Fields {
					lines.WriteString("		AX{`=`, " + I.ToStr(idx) + ", x." + S.CamelCase(prop.Name) + "},\n")
				}
				lines.WriteString("	}\n")
				lines.WriteString("}\n\n")

				// index functions
				for idx, prop := range props.Fields {
					lines.WriteString(`func (x *` + structName + ") idx" + S.CamelCase(prop.Name) + "() int {\n")
					lines.WriteString("	return " + I.ToStr(idx) + "\n")
					lines.WriteString("}\n\n")
				}

				// to AX
				lines.WriteString(`func (x *` + structName + ") ToArray() AX {\n")
				lines.WriteString("	return AX{\n")
				for idx, prop := range props.Fields {
					lines.WriteString("		x." + S.CamelCase(prop.Name) + ", // " + I.ToStr(idx) + "\n")
				}
				lines.WriteString("	}\n")
				lines.WriteString("}\n\n")

				// from AX
				lines.WriteString(`func (x *` + structName + `) FromArray(a AX) *` + structName + " {\n")
				for idx, prop := range props.Fields {
					lines.WriteString("	x." + S.CamelCase(prop.Name) + ` = ` + typeConverter[prop.Type] + "(a[" + I.ToStr(idx) + "])\n")
				}
				lines.WriteString("	return x\n")
				lines.WriteString("}\n\n")

				// find one
				lines.WriteString(`func (x *` + structName + ") FindOne(a *Adapter) bool {\n")
				lines.WriteString("	res, err := a.conn.Select(x.SpaceName(), x.PrimaryIndex(), 0, 1, tarantool.IterEq, AX{" + uniqueCamel + "})\n")
				lines.WriteString("	if tracer.IsError(err) {\n")
				lines.WriteString("		return false\n")
				lines.WriteString("	}\n")
				lines.WriteString("	rows := res.Tuples()\n")
				lines.WriteString("	if len(rows) == 1 {\n")
				lines.WriteString("		x.FromArray(rows[0])\n")
				lines.WriteString("		return true\n")
				lines.WriteString("	}\n")
				lines.WriteString("	return false\n")
				lines.WriteString("}\n\n")

				// find many
				lines.WriteString(`func (x *` + structName + ") FindOffsetLimit(a *Adapter, offset, limit uint32, idx string, val AX) []*" + structName + " {\n")
				lines.WriteString("	var rows []*" + structName + "\n")
				lines.WriteString("	res, err := a.conn.Select(x.SpaceName(), idx, offset, limit, tarantool.IterEq, val)\n")
				lines.WriteString("	if tracer.IsError(err) {\n")
				lines.WriteString("		return rows\n")
				lines.WriteString("	}\n")
				lines.WriteString("	for _, row := range res.Tuples() {\n")
				lines.WriteString("		item := &" + structName + "{}\n")
				lines.WriteString("		rows = append(rows, item.FromArray(row))\n")
				lines.WriteString("	}\n")
				lines.WriteString("	return rows\n\n")
				lines.WriteString("}\n\n")

				//// set to min value
				//lines.WriteString(`func (x *` + structName + ") ResetToMax() {\n")
				//for _, prop := range props.Fields {
				//	lines.WriteString("	x." + S.CamelCase(prop.Name) + " = " + maxMap[prop.Type] + "\n")
				//}
				//lines.WriteString("}\n\n")
				//
				//// set to min value
				//lines.WriteString(`func (x *` + structName + ") ResetToMin() {\n")
				//for _, prop := range props.Fields {
				//	lines.WriteString("	x." + S.CamelCase(prop.Name) + " = " + minMap[prop.Type] + "\n")
				//}
				//lines.WriteString("}\n\n")
				//
				//// set if greater
				//lines.WriteString(`func (x *` + structName + ") SetIfLesser(l *"+structName+") {\n")
				//for _, prop := range props.Fields {
				//	propName := S.CamelCase(prop.Name)
				//	lines.WriteString("	if x." + propName + " > l." + propName + " {\n")
				//	lines.WriteString("		x." + propName + " = l." + propName + "\n")
				//	lines.WriteString("	}\n")
				//}
				//lines.WriteString("}\n\n")
				//
				//// set if greater
				//lines.WriteString(`func (x *` + structName + ") SetIfGreater(l *"+structName+") {\n")
				//for _, prop := range props.Fields {
				//	propName := S.CamelCase(prop.Name)
				//	lines.WriteString("	if x." + propName + " < l." + propName + " {\n")
				//	lines.WriteString("		x." + propName + " = l." + propName + "\n")
				//	lines.WriteString("	}\n")
				//}
				//lines.WriteString("}\n\n")
				lines.WriteString("\n")
				lines.WriteString(warning)
			}
			err = ioutil.WriteFile(genDir+genFileName, lines.Bytes(), os.ModePerm)
			if L.IsError(err, `failed writing file: `+genDir+genFileName) {
				return false
			}
		}
	}
	return true
}
