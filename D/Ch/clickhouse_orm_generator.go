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
//func dq(str string) string {
//	return `"` + str + `"`
//}

var typeTranslator = map[DataType]string{
	//Uuid:     `string`,
	UInt64:     `uint64`,
	Float64:    `float64`,
	String:     `string`,
	Int64:      `int64`,
	IPv4:       `net.IP`,
	IPv6:       `net.IP`,
	DateTime:   `time.Time`,
	DateTime64: `time.Time`,
	Int8:       `int8`,
	Int16:      `int16`,
	Int32:      `int32`,
	UInt8:      `uint8`,
	UInt16:     `uint16`,
	UInt32:     `uint32`,
	Float32:    `float32`,
}

//var typeConverter = map[DataType]string{
//	UInt64:  `X.ToU`,
//	Float64: `X.ToF`,
//	String:  `X.ToS`,
//	Int64:   `X.ToI`,
//	Int8:    `X.ToByte`,
//}

const connStruct = `Ch.Adapter`
const connImport = "\n\n\t_ `github.com/ClickHouse/clickhouse-go/v2`"
const connTestImport = "\n\t_ `github.com/ClickHouse/clickhouse-go/v2`"
const buffImport = "\n\tchBuffer `github.com/kokizzu/ch-timed-buffer`"

const warning = "// DO NOT EDIT, will be overwritten by github.com/kokizzu/Ch/clickhouse_orm_generator.go\n\n"

func GenerateOrm(tables map[TableName]*TableProp) {
	ci := L.CallerInfo(2)
	this := L.CallerInfo()
	pkgName := S.RightOfLast(ci.PackageName, `/`)
	saPkgName := `sa` + pkgName[1:] // write/command (mutator)
	L.Print(saPkgName)
	mPkgName := `m` + pkgName[1:]

	// generate
	apBuf := bytes.Buffer{}
	apTestBuf := bytes.Buffer{}

	SA := func(str string) {
		_, err := apBuf.WriteString(str)
		L.PanicIf(err, `failed apBuf.WriteString`)
	}
	SAT := func(str string) {
		_, err := apTestBuf.WriteString(str)
		L.PanicIf(err, `failed apTestBuf.WriteString`)
	}

	// sort by table name to keep the order when regenerating structs
	tableNames := make([]string, 0, len(tables))
	for k := range tables {
		tableNames = append(tableNames, string(k))
	}
	sort.Strings(tableNames)

	// check have IP
	haveIp := false
	for _, tableName := range tableNames {
		props := tables[TableName(tableName)]
		for _, prop := range props.Fields {
			if prop.Type == IPv4 || prop.Type == IPv6 {
				haveIp = true
				break
			}
		}
	}

	//SA(`// generated: ` + time.Now().String() + "\n")
	SA(`package ` + saPkgName)
	SA("\n\n")
	SA(warning)
	SAT(`package ` + saPkgName)
	SAT("\n\n")
	SAT(warning)

	SA(`import (`)

	// import reader
	SA(qi(`database/sql`))
	if haveIp {
		SA(qi(`net`))
	}
	SA(qi(`time`))
	SA("\n")
	SA(qi(ci.PackageName)) // /models/m*
	SA(connImport)
	SA(buffImport)

	SA("\n")
	SA(qi(`github.com/kokizzu/gotro/A`))
	SA(qi(this.PackageName)) // github.com/kokizzu/gotro/D/Ch
	SA(qi(`github.com/kokizzu/gotro/L`))
	//SA(qi(`github.com/kokizzu/gotro/X`))

	SA(`
)` + "\n\n")

	SAT(`import (`)
	SAT(qi(`database/sql`))
	SAT(qi(`fmt`))
	SAT(qi(`log`))
	SAT(qi(`os`))
	SAT(qi(`testing`))
	SAT(qi(`time`))
	SAT("\n")
	SAT(connTestImport)
	SAT(qi(`github.com/ory/dockertest/v3`))
	SAT(qi(`github.com/stretchr/testify/assert`))
	SAT(qi(this.PackageName)) // github.com/kokizzu/gotro/D/Ch
	SAT(qi(`github.com/kokizzu/gotro/L`))
	if haveIp {
		SAT(qi(`net`))
	}
	SAT(`
)` + "\n\n")

	SAT(`var globalPool *dockertest.Pool
var dbConn *sql.DB
var reconnect func() *sql.DB

func prepareDb(onReady func(db *sql.DB) int) {
	const dockerRepo = "yandex/clickhouse-server"
	const dockerVer = "latest"
	const chPort = "9000/tcp"
	const dbDriver = "clickhouse"
	const dbConnStr = "tcp://127.0.0.1:%s?debug=true"
	var err error
	if globalPool == nil {
		globalPool, err = dockertest.NewPool("")
		if err != nil {
			log.Printf("Could not connect to docker: %s\n", err)
			os.Exit(onReady(nil))
		}
	}
	resource, err := globalPool.Run(dockerRepo, dockerVer, []string{})
	if err != nil {
		log.Printf("Could not start resource: %s\n", err)
		os.Exit(onReady(nil))
	}
	var db *sql.DB
	if err := globalPool.Retry(func() error {
		var err error
		connStr := fmt.Sprintf(dbConnStr, resource.GetPort(chPort))
		reconnect = func() *sql.DB {
			db, err = sql.Open(dbDriver, connStr)
			L.IsError(err, "sql.Open: "+connStr)
			return db
		}
		reconnect()
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Printf("Could not connect to docker: %s\n", err)
		os.Exit(onReady(nil))
	}
	code := onReady(db)
	if err := globalPool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(code)
}

func TestMain(m *testing.M) {
	prepareDb(func(db *sql.DB) int {
		dbConn = db
		return m.Run()
	})
}

` + "\n")

	SA(`//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file ` + saPkgName + "__ORM.GEN.go\n")
	SA(`//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type ` + saPkgName + "__ORM.GEN.go\n")
	SA(`//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type ` + saPkgName + "__ORM.GEN.go\n")
	SA(`//go:generate replacer -afterprefix "By\" form" "By,string\" form" type ` + saPkgName + "__ORM.GEN.go\n")
	SA(`// go:generate msgp -tests=false -file ` + saPkgName + `__ORM.GEN.go -o ` + saPkgName + `__MSG.GEN.go` + "\n\n")

	for _, tableName := range tableNames {
		SA(`var ` + tableName + `Dummy = ` + S.PascalCase(tableName) + "{}\n")
	}

	SA("var Preparators = map[Ch.TableName]chBuffer.Preparator{\n")
	for _, tableName := range tableNames {
		SA(`	` + mPkgName + `.Table` + S.PascalCase(tableName) + `: func(tx *sql.Tx) *sql.Stmt {
		query := ` + tableName + `Dummy.SqlInsert()
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
		structName := S.PascalCase(tableName)
		maxLen := 1
		for _, prop := range props.Fields {
			l := len(prop.Name) + 1
			if maxLen < l {
				maxLen = l
			}
		}

		// reader struct
		SA(`type ` + structName + " struct {\n")
		SA("	Adapter *" + connStruct + " `json:" + `"-"` + " msg:" + `"-"` + " query:" + `"-"` + " form:" + `"-"` + "`\n")
		for _, prop := range props.Fields {
			camel := S.PascalCase(prop.Name)
			SA("	" + camel + strings.Repeat(` `, maxLen-len(camel)) + typeTranslator[prop.Type] + "\n")
		}
		SA("}\n\n")

		// reader struct constructor
		SA(`func New` + structName + `(adapter *` + connStruct + `) *` + structName + " {\n")
		SA(`	return &` + structName + "{Adapter: adapter}\n")
		SA("}\n\n")

		// field type map
		SA("// " + structName + "FieldTypeMap returns key value of field name and key\n")
		SA("var " + structName + "FieldTypeMap = map[string]Ch.DataType { //nolint:dupl false positive\n")
		for _, field := range props.Fields {
			SA("	" + S.BT(field.Name) + `:` + strings.Repeat(` `, maxLen-len(field.Name)) + TypeToConst[field.Type] + ",\n")
		}
		SA("}\n\n")

		// table name
		receiverName := strings.ToLower(string(structName[0]))
		SA(`func (` + receiverName + ` *` + structName + ") TableName() Ch.TableName { //nolint:dupl false positive\n")
		SA("	return " + mPkgName + `.Table` + structName + "\n")
		SA("}\n\n")

		// Sql table name
		SA(`func (` + receiverName + ` *` + structName + ") SqlTableName() string { //nolint:dupl false positive\n")
		SA("	return `" + `"` + tableName + `"` + "`\n")
		SA("}\n\n")

		// ScanRowAllColumns
		SA(`func (` + receiverName + ` *` + structName + ") ScanRowAllCols(rows *sql.Rows) (err error) { //nolint:dupl false positive\n")
		SA("	return rows.Scan(\n")
		for _, prop := range props.Fields {
			SA("		&" + receiverName + "." + S.PascalCase(prop.Name) + ",\n")
		}
		SA("	)\n")
		SA("}\n\n")

		// ScanRowsAllColumns
		SA(`func (` + receiverName + ` *` + structName + ") ScanRowsAllCols(rows *sql.Rows, estimateRows int) (res []" + structName + ", err error) { //nolint:dupl false positive\n")
		SA("	res = make([]" + structName + ", 0, estimateRows)\n")
		SA("	defer rows.Close()\n")
		SA("	for rows.Next() {\n")
		SA("		var row " + structName + "\n")
		SA("		err = row.ScanRowAllCols(rows)\n")
		SA("		if err != nil {\n")
		SA("			return\n")
		SA("		}\n")
		SA("		res = append(res, row)\n")
		SA("	}\n")
		SA("	return\n")
		SA("}\n\n")

		// insert, error if exists
		SA("// insert, error if exists\n")
		SA(`func (` + receiverName + ` *` + structName + ") SqlInsert() string { //nolint:dupl false positive\n")
		qMark := S.Repeat(`,?`, len(props.Fields))[1:]
		SA("	return `INSERT INTO ` + " + receiverName + ".SqlTableName() + ` (` + " + receiverName + ".SqlAllFields() + `) VALUES (" + qMark + ")`\n")
		SA("}\n\n")

		// total records
		SA(`func (` + receiverName + ` *` + structName + ") SqlCount() string { //nolint:dupl false positive\n")
		SA("	return `SELECT COUNT(*) FROM ` + " + receiverName + ".SqlTableName()\n")
		SA("}\n\n")

		// Sql select all fields, used when need to mutate or show every fields
		SA(`func (` + receiverName + ` *` + structName + ") SqlSelectAllFields() string { //nolint:dupl false positive\n")
		sqlFields := ``
		for _, prop := range props.Fields {
			sqlFields += `, ` + (prop.Name) + "\n\t"
		}
		SA(`	return ` + bq(sqlFields[1:]) + "\n")
		SA("}\n\n")

		SA(`func (` + receiverName + ` *` + structName + ") SqlAllFields() string { //nolint:dupl false positive\n")
		sqlFields = ``
		for _, prop := range props.Fields {
			sqlFields += `, ` + prop.Name + ""
		}
		SA(`	return ` + bq(sqlFields[2:]) + "\n")
		SA("}\n\n")

		// to Insert parameter
		SA(`func (` + receiverName + ` ` + structName + ") SqlInsertParam() []any { //nolint:dupl false positive\n")
		SA("	return []any{\n")
		for idx, prop := range props.Fields {
			SA("		" + receiverName + "." + S.PascalCase(prop.Name) + ", // " + I.ToStr(idx) + " \n")
		}
		SA("	}\n")
		SA("}\n\n")

		for idx, prop := range props.Fields {
			propName := S.PascalCase(prop.Name)

			// index functions
			SA(`func (` + receiverName + ` *` + structName + ") Idx" + propName + "() int { //nolint:dupl false positive\n")
			SA("	return " + X.ToS(idx) + "\n")
			SA("}\n\n")

			// column name functions
			//SA(`func (` + receiverName + ` *` + structName + ") col" + propName + "() string { //nolint:dupl false positive\n")
			//SA("	return `" + prop.Name + "`\n")
			//SA("}\n\n")

			// Sql column name functions
			SA(`func (` + receiverName + ` *` + structName + ") Sql" + propName + "() string { //nolint:dupl false positive\n")
			SA("	return `" + prop.Name + "`\n")
			SA("}\n\n")
		}

		// to AX
		SA(`func (` + receiverName + ` *` + structName + ") ToArray() A.X { //nolint:dupl false positive\n")
		SA("	return A.X{\n")
		for idx, prop := range props.Fields {
			camel := S.PascalCase(prop.Name)
			SA("		" + receiverName + "." + camel + `,` + strings.Repeat(` `, maxLen-len(camel)) + `// ` + X.ToS(idx) + "\n")
		}
		SA("	}\n")
		SA("}\n\n")

		//// from AX
		//SA(`func (` + receiverName + ` *` + structName + `) FromArray(a A.X) *` + structName + " { //nolint:dupl false positive\n")
		//for idx, prop := range props.Fields {
		//	SA("	" + receiverName + "." + S.PascalCase(prop.Name) + ` = ` + typeConverter[prop.Type] + "(a[" + X.ToS(idx) + "])\n")
		//}
		//SA("	return " + receiverName + "\n")
		//SA("}\n\n")

		// generated tests
		SAT(`func TestGenerated` + structName + `Helpers(t *testing.T) {` + "\n")
		SAT(`	obj := New` + structName + `(nil)` + "\n")
		SAT(`	assert.NotNil(t, obj)` + "\n")
		SAT(`	assert.NotEmpty(t, obj.TableName())` + "\n")
		SAT(`	assert.NotEmpty(t, obj.SqlTableName())` + "\n")
		for _, field := range props.Fields {
			camel := S.PascalCase(field.Name)
			SAT(`	obj.` + camel + ` = ` + chTestSampleLiteral(field.Type, false) + "\n")
		}
		SAT(`	assert.NotEmpty(t, obj.SqlInsert())` + "\n")
		SAT(`	assert.NotEmpty(t, obj.SqlCount())` + "\n")
		SAT(`	assert.NotEmpty(t, obj.SqlSelectAllFields())` + "\n")
		SAT(`	assert.NotEmpty(t, obj.SqlAllFields())` + "\n")
		SAT(`	arr := obj.ToArray()` + "\n")
		SAT(`	assert.Len(t, arr, ` + X.ToS(len(props.Fields)) + `)` + "\n")
		SAT(`	params := obj.SqlInsertParam()` + "\n")
		SAT(`	assert.Len(t, params, ` + X.ToS(len(props.Fields)) + `)` + "\n")
		if len(props.Fields) > 0 {
			SAT(`	_, ok := ` + structName + `FieldTypeMap[` + S.BT(props.Fields[0].Name) + `]` + "\n")
			SAT(`	assert.True(t, ok)` + "\n")
		}
		for idx, field := range props.Fields {
			camel := S.PascalCase(field.Name)
			SAT(`	assert.Equal(t, ` + X.ToS(idx) + `, obj.Idx` + camel + `())` + "\n")
			SAT(`	assert.Equal(t, ` + S.BT(field.Name) + `, obj.Sql` + camel + `())` + "\n")
		}
		SAT(`	prep, ok := Preparators[obj.TableName()]` + "\n")
		SAT(`	assert.True(t, ok)` + "\n")
		SAT(`	assert.NotNil(t, prep)` + "\n")
		SAT(`}` + "\n\n")

		SAT(`func TestGenerated` + structName + `CRUD(t *testing.T) {` + "\n")
		SAT(`	if dbConn == nil {` + "\n")
		SAT(`		t.Skip("docker unavailable")` + "\n")
		SAT(`	}` + "\n")
		SAT(`	a := &Ch.Adapter{DB: dbConn, Reconnect: reconnect}` + "\n")
		SAT(`	obj := New` + structName + `(a)` + "\n")
		SAT(`	ok := a.UpsertTable(obj.TableName(), ` + renderChTablePropLiteral(props) + `)` + "\n")
		SAT(`	assert.True(t, ok)` + "\n")
		SAT(`	_, _ = a.Exec("TRUNCATE TABLE " + string(obj.TableName()))` + "\n")
		SAT(`	row1 := New` + structName + `(a)` + "\n")
		for _, field := range props.Fields {
			camel := S.PascalCase(field.Name)
			SAT(`	row1.` + camel + ` = ` + chTestSampleLiteral(field.Type, false) + "\n")
		}
		SAT(`	_, err := a.Exec(row1.SqlInsert(), row1.SqlInsertParam()...)` + "\n")
		SAT(`	assert.NoError(t, err)` + "\n")
		SAT(`	rows, err := a.Query("SELECT " + row1.SqlSelectAllFields() + " FROM " + row1.SqlTableName() + " LIMIT 1")` + "\n")
		SAT(`	assert.NoError(t, err)` + "\n")
		SAT(`	assert.True(t, rows.Next())` + "\n")
		SAT(`	got := New` + structName + `(a)` + "\n")
		SAT(`	assert.NoError(t, got.ScanRowAllCols(rows))` + "\n")
		SAT(`	assert.NoError(t, rows.Close())` + "\n")
		SAT(`	row2 := New` + structName + `(a)` + "\n")
		for _, field := range props.Fields {
			camel := S.PascalCase(field.Name)
			SAT(`	row2.` + camel + ` = ` + chTestSampleLiteral(field.Type, true) + "\n")
		}
		SAT(`	_, err = a.Exec(row2.SqlInsert(), row2.SqlInsertParam()...)` + "\n")
		SAT(`	assert.NoError(t, err)` + "\n")
		SAT(`	rows, err = a.Query("SELECT " + row1.SqlSelectAllFields() + " FROM " + row1.SqlTableName() + " LIMIT 10")` + "\n")
		SAT(`	assert.NoError(t, err)` + "\n")
		SAT(`	parsed, err := row1.ScanRowsAllCols(rows, 10)` + "\n")
		SAT(`	assert.NoError(t, err)` + "\n")
		SAT(`	assert.GreaterOrEqual(t, len(parsed), 2)` + "\n")
		SAT(`	var cnt int64` + "\n")
		SAT(`	assert.NoError(t, a.QueryRow(row1.SqlCount()).Scan(&cnt))` + "\n")
		SAT(`	assert.GreaterOrEqual(t, cnt, int64(1))` + "\n")
		SAT(`}` + "\n\n")

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
	apTestFname := fmt.Sprintf(`./%s/%s__ORM.GEN_test.go`, saPkgName, saPkgName)
	err = os.WriteFile(apTestFname, apTestBuf.Bytes(), os.ModePerm)
	if L.IsError(err, `os.WriteFile failed: `+apTestFname) {
		return
	}
}

func chTypeConst(typ DataType) string {
	if c, ok := TypeToConst[typ]; ok && c != `` {
		return c
	}
	return `Ch.` + string(typ)
}

func chTestSampleLiteral(typ DataType, alt bool) string {
	switch typ {
	case UInt64:
		if alt {
			return `uint64(2)`
		}
		return `uint64(1)`
	case UInt32:
		if alt {
			return `uint32(2)`
		}
		return `uint32(1)`
	case UInt16:
		if alt {
			return `uint16(2)`
		}
		return `uint16(1)`
	case UInt8:
		if alt {
			return `uint8(2)`
		}
		return `uint8(1)`
	case Int64:
		if alt {
			return `int64(2)`
		}
		return `int64(1)`
	case Int32:
		if alt {
			return `int32(2)`
		}
		return `int32(1)`
	case Int16:
		if alt {
			return `int16(2)`
		}
		return `int16(1)`
	case Int8:
		if alt {
			return `int8(2)`
		}
		return `int8(1)`
	case Float64:
		if alt {
			return `2.5`
		}
		return `1.5`
	case Float32:
		if alt {
			return `float32(2.5)`
		}
		return `float32(1.5)`
	case DateTime, DateTime64:
		if alt {
			return `time.Now().UTC().Add(time.Second).Truncate(time.Second)`
		}
		return `time.Now().UTC().Truncate(time.Second)`
	case IPv4:
		if alt {
			return `net.ParseIP("10.0.0.2")`
		}
		return `net.ParseIP("10.0.0.1")`
	case IPv6:
		if alt {
			return `net.ParseIP("2001:db8::2")`
		}
		return `net.ParseIP("2001:db8::1")`
	case Decimal:
		if alt {
			return `"2.5"`
		}
		return `"1.5"`
	case FixedString:
		if alt {
			return `"sample2"`
		}
		return `"sample"`
	case String:
		if alt {
			return `"sample2"`
		}
		return `"sample"`
	default:
		if alt {
			return `"sample2"`
		}
		return `"sample"`
	}
}

func renderChTablePropLiteral(props *TableProp) string {
	buf := bytes.Buffer{}
	buf.WriteString("&Ch.TableProp{\n")
	buf.WriteString("Fields: []Ch.Field{\n")
	for _, field := range props.Fields {
		buf.WriteString("{`" + field.Name + "`, " + chTypeConst(field.Type) + "},\n")
	}
	buf.WriteString("},\n")
	if props.Engine != `` {
		buf.WriteString("Engine: `" + props.Engine + "`,\n")
	}
	if len(props.Partitions) > 0 {
		buf.WriteString("Partitions: []string{")
		for idx, partition := range props.Partitions {
			if idx > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString("`" + partition + "`")
		}
		buf.WriteString("},\n")
	}
	if len(props.Orders) > 0 {
		buf.WriteString("Orders: []string{")
		for idx, order := range props.Orders {
			if idx > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString("`" + order + "`")
		}
		buf.WriteString("},\n")
	}
	if props.DefaultCodec != `` {
		buf.WriteString("DefaultCodec: `" + props.DefaultCodec + "`,\n")
	}
	buf.WriteString("}")
	return buf.String()
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
