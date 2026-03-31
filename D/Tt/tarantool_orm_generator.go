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
const connImport = "\n\n\t`github.com/tarantool/go-tarantool/v2`"
const iterEq = `tarantool.IterEq`
const iterAll = `tarantool.IterAll`
const iterNeighbor = `tarantool.IterNeighbor`

var _ = iterNeighbor

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
	rqTestBuf := bytes.Buffer{}
	wcTestBuf := bytes.Buffer{}

	RQ := func(str string) {
		_, err := rqBuf.WriteString(str)
		L.PanicIf(err, `failed rqBuf.WriteString`)
	}
	WC := func(str string) {
		_, err := wcBuf.WriteString(str)
		L.PanicIf(err, `failed wcBuf.WriteString`)
	}
	RQT := func(str string) {
		_, err := rqTestBuf.WriteString(str)
		L.PanicIf(err, `failed rqTestBuf.WriteString`)
	}
	WCT := func(str string) {
		_, err := wcTestBuf.WriteString(str)
		L.PanicIf(err, `failed wcTestBuf.WriteString`)
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
	RQT(`package ` + rqPkgName + "\n\n" + warning)
	WCT(`package ` + wcPkgName + "\n\n" + warning)

	RQT(`import (`)
	RQT(qi(`testing`))
	RQT(qi(`github.com/stretchr/testify/assert`))
	RQT(qi(this.PackageName))
	RQT(`
)` + "\n\n")

	WCT(`import (`)
	WCT(qi(`context`))
	WCT(qi(`fmt`))
	WCT(qi(`log`))
	WCT(qi(`os`))
	WCT(qi(`testing`))
	WCT(qi(`time`))
	WCT(qi(ci.PackageName))
	WCT(qi(this.PackageName))
	WCT(qi(`github.com/kokizzu/gotro/L`))
	WCT(qi(`github.com/kokizzu/gotro/S`))
	WCT(qi(`github.com/ory/dockertest/v3`))
	WCT(qi(`github.com/stretchr/testify/assert`))
	WCT(qi(`github.com/tarantool/go-tarantool/v2`))
	WCT(`
)` + "\n\n")

	WCT(`var globalPool *dockertest.Pool
var reconnect func() *tarantool.Connection
var dbConn *tarantool.Connection

func prepareDb(onReady func(db *tarantool.Connection) int) {
	const dockerRepo = "tarantool/tarantool"
	const dockerVer = "3.1"
	const ttPort = "3301/tcp"
	const dbConnStr = "127.0.0.1:%s"
	const dbUser = "guest"
	const dbPass = ""
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
	var db *tarantool.Connection
	if err := globalPool.Retry(func() error {
		var err error
		connStr := fmt.Sprintf(dbConnStr, resource.GetPort(ttPort))
		reconnect = func() *tarantool.Connection {
			db, err = tarantool.Connect(context.Background(), tarantool.NetDialer{
				Address:  connStr,
				User:     dbUser,
				Password: dbPass,
			}, tarantool.Opts{
				Timeout: 8 * time.Second,
			})
			if err != nil && !S.Contains(err.Error(), "failed to read greeting: EOF") {
				L.IsError(err, "tarantool.Connect")
			}
			return db
		}
		reconnect()
		if err != nil {
			return err
		}
		_, err = db.Do(tarantool.NewPingRequest()).Get()
		return err
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
	prepareDb(func(db *tarantool.Connection) int {
		dbConn = db
		return m.Run()
	})
}

func TestGeneratedSanity(t *testing.T) {
	if dbConn == nil {
		t.Skip("docker unavailable")
	}
	conn := dbConn
	a := &Tt.Adapter{Connection: conn, Reconnect: reconnect}
	a.MigrateTables(map[Tt.TableName]*Tt.TableProp{
		"bands": {
			Fields: []Tt.Field{
				{"id", Tt.Unsigned},
				{"band_name", Tt.String},
				{"year", Tt.Unsigned},
			},
			AutoIncrementId: true,
			Unique1:         "band_name",
			Indexes:         []string{"year"},
		},
	})
	tuples := [][]any{
		{1, "Roxette", 1986},
		{2, "Scorpions", 1965},
		{3, "Ace of Base", 1987},
		{4, "The Beatles", 1960},
	}
	for _, tuple := range tuples {
		_, err := conn.Do(tarantool.NewInsertRequest("bands").Tuple(tuple)).Get()
		assert.NoError(t, err)
	}
	_, err := conn.Do(tarantool.NewSelectRequest("bands").Limit(10).Iterator(tarantool.IterEq).Key([]any{uint(1)})).Get()
	assert.NoError(t, err)
	_, err = conn.Do(tarantool.NewSelectRequest("bands").Index("band_name").Limit(10).Iterator(tarantool.IterEq).Key([]any{"The Beatles"})).Get()
	assert.NoError(t, err)
	_, err = conn.Do(tarantool.NewUpdateRequest("bands").Key(tarantool.IntKey{2}).Operations(tarantool.NewOperations().Assign(1, "Pink Floyd"))).Get()
	assert.NoError(t, err)
	_, err = conn.Do(tarantool.NewUpsertRequest("bands").Tuple([]any{uint(5), "The Rolling Stones", 1962}).Operations(tarantool.NewOperations().Assign(1, "The Doors"))).Get()
	assert.NoError(t, err)
	_, err = conn.Do(tarantool.NewReplaceRequest("bands").Tuple([]any{1, "Queen", 1970})).Get()
	assert.NoError(t, err)
	_, err = conn.Do(tarantool.NewDeleteRequest("bands").Key([]any{uint(5)})).Get()
	assert.NoError(t, err)
}
` + "\n")

	haveString, haveAutoIncrementId := false, false
	useGraphql := len(withGraphql) > 0
	// sort by table name to keep the order when regenerating structs
	tableNames := make([]string, 0, len(tables))
	for k, v := range tables {
		tableNames = append(tableNames, string(k))
		useGraphql = useGraphql || v.GenGraphqlType // if one of them use graphql, import anyway
		for _, prop := range v.Fields {
			if prop.Type == String {
				haveString = true
			}
		}
		if v.AutoIncrementId {
			haveAutoIncrementId = true
		}
	}
	sort.Strings(tableNames)

	// import reader
	BOTH(`import (`)

	RQ(qi(ci.PackageName))
	WC(qi(ci.PackageName + `/` + rqPkgName))
	BOTH(connImport)

	BOTH(NL)
	if useGraphql {
		//RQ(qi(`github.com/   graphql-go/graphql`))
	}
	BOTH(qi(`github.com/kokizzu/gotro/A`))
	BOTH(qi(this.PackageName)) // github.com/kokizzu/gotro/D/Tt
	BOTH(qi(`github.com/kokizzu/gotro/L`))
	WC(qi(`github.com/kokizzu/gotro/M`))
	if haveString {
		WC(qi(`github.com/kokizzu/gotro/S`))
	}
	if haveAutoIncrementId {
		BOTH(qi(`github.com/kokizzu/gotro/X`))
	}

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
		propByName := map[string]Field{}
		censoredFieldsByName := map[string]bool{}
		for _, prop := range props.Fields {
			l := len(prop.Name) + 1
			if maxLen < l {
				maxLen = l
			}
			propByName[prop.Name] = prop
		}
		for _, propName := range props.AutoCensorFields {
			censoredFieldsByName[propName] = true
		}

		// mutator struct
		WC("// " + structName + "Mutator DAO writer/command struct\n")
		WC(`type ` + structName + "Mutator struct {\n")
		WC(`	` + rqPkgName + `.` + structName + NL)
		WC("	mutations *tarantool.Operations\n")
		WC("	logs	  []A.X\n")
		WC("}\n\n")

		// mutator struct constructor
		WC("// New" + structName + "Mutator create new ORM writer/command object\n")
		WC(`func New` + structName + `Mutator(adapter *` + connStruct + `) (res *` + structName + "Mutator) {\n")
		WC(`	res = &` + structName + `Mutator{` + structName + `: ` + rqPkgName + `.` + structName + "{Adapter: adapter}}\n")
		WC("	res.mutations = tarantool.NewOperations()\n")
		for _, prop := range props.Fields {
			if prop.Type == Array {
				WC(`	res.` + S.PascalCase(prop.Name) + ` = ` + TypeToGoType[prop.Type] + "{}\n")
			}
		}
		WC("	return\n")
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
		WC(`	return len(` + receiverName + ".logs) > 0\n")
		WC("}\n\n")

		// clear mutation
		WC("// ClearMutations clear all previously called Set* methods\n")
		WC(`func (` + receiverName + ` *` + structName + "Mutator) ClearMutations() { //nolint:dupl false positive\n")
		WC(`	` + receiverName + ".mutations = tarantool.NewOperations()\n")
		WC(`	` + receiverName + ".logs = []A.X{}\n")
		WC("}\n\n")

		// auto increment id
		if props.AutoIncrementId {
			uniquePropCamel := S.PascalCase(IdCol)
			structProp := receiverName + `.` + uniquePropCamel

			RQ(`func (` + receiverName + ` *` + structName + `) UniqueIndex` + uniquePropCamel + "() string { //nolint:dupl false positive\n")
			RQ("	return " + S.BT(IdCol) + NL)
			RQ("}\n\n")

			keyFunc := Field{Type: Unsigned}.KeyRenderer()
			generateMutationByUniqueIndex(uniquePropCamel, structProp, receiverName, structName, keyFunc, RQ, WC)

			if props.GenGraphqlType {
				generateGraphqlQueryField(structName, uniquePropCamel, propByName[IdCol], RQ)
			}
		}

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

			keyType := propByName[props.Unique1].KeyRenderer()
			generateMutationByUniqueIndex(uniquePropCamel, structProp, receiverName, structName, keyType, RQ, WC)

			if props.GenGraphqlType {
				generateGraphqlQueryField(structName, uniquePropCamel, propByName[props.Unique1], RQ)
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

			keyType := propByName[props.Unique2].KeyRenderer()
			generateMutationByUniqueIndex(uniquePropCamel, structProp, receiverName, structName, keyType, RQ, WC)
		}

		// unique index3
		if props.Unique3 != `` && !(props.AutoIncrementId && props.Unique3 == IdCol) {
			uniquePropCamel := S.PascalCase(props.Unique3)
			structProp := receiverName + `.` + uniquePropCamel

			RQ("// UniqueIndex" + uniquePropCamel + " return unique index name\n")
			RQ(`func (` + receiverName + ` *` + structName + `) UniqueIndex` + uniquePropCamel + "() string { //nolint:dupl false positive\n")
			RQ("	return " + S.BT(props.Unique3) + NL)
			RQ("}\n\n")

			keyType := propByName[props.Unique3].KeyRenderer()
			generateMutationByUniqueIndex(uniquePropCamel, structProp, receiverName, structName, keyType, RQ, WC)
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

			keyFunc := func(structProp string) string { return "A.X{" + structProp + "}" }
			generateMutationByUniqueIndex(uniquePropCamel, structProps, receiverName, structName, keyFunc, RQ, WC)
		}

		// insert, error if exists
		WC("// DoInsert insert, error if already exists\n")
		WC(`func (` + receiverName + ` *` + structName + "Mutator) DoInsert() bool { //nolint:dupl false positive\n")
		ret1 := S.IfElse(props.AutoIncrementId, `row`, `_`)
		WC("	arr := " + receiverName + ".ToArray()\n")
		WC("	" + ret1 + ", err := " + receiverName + ".Adapter.RetryDo(\n")
		WC("		tarantool.NewInsertRequest(" + receiverName + ".SpaceName()).\n")
		WC("		Tuple(arr),\n")
		WC("	)\n")
		if props.AutoIncrementId {
			WC("	if err == nil {\n")
			WC("		if len(row) > 0 {\n")
			WC("			if cells, ok := row[0].([]any); ok && len(cells) > 0 {\n")
			WC("				" + receiverName + ".Id = X.ToU(cells[0])\n")
			WC("			}\n")
			WC("		}\n")
			WC("	}\n")
		}
		WC("	return !L.IsError(err, `" + structName + ".DoInsert failed: `+" + receiverName + ".SpaceName() + `\\n%#v`, arr)\n")
		WC("}\n\n")

		// replace = upsert, only error when there's unique secondary key
		// https://github.com/tarantool/tarantool/issues/5732
		if props.AutoIncrementId {
			WC("// DoUpsert upsert, insert or overwrite, will error only when there's unique secondary key being violated\n")
			WC("// tarantool's replace/upsert can only match by primary key\n")
			WC("// previous name: DoReplace\n")
			WC(`func (` + receiverName + ` *` + structName + "Mutator) DoUpsertById() bool { //nolint:dupl false positive\n")
			WC("	if " + receiverName + ".Id > 0 {\n")
			WC("		return " + receiverName + ".DoUpdateById()\n")
			WC("	}\n")
			WC("	return " + receiverName + ".DoInsert()\n")
			WC("}\n\n")
		}

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
		RQ(`func (` + receiverName + ` *` + structName + ") ToUpdateArray() *tarantool.Operations { //nolint:dupl false positive\n")
		RQ("	return tarantool.NewOperations().\n")
		last := len(props.Fields) - 1
		for idx, prop := range props.Fields {
			RQ("		Assign(" + I.ToStr(idx) + ", " + receiverName + "." + S.PascalCase(prop.Name) + ")" +
				S.If(idx != last, ".") + "\n")
		}
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
				WC("		" + receiverName + ".mutations.Assign(" + I.ToStr(idx) + ", val)\n")
				if !censoredFieldsByName[prop.Name] {
					WC("		" + receiverName + ".logs = append(" + receiverName + ".logs, A.X{`" + prop.Name + "`, " + receiverName + `.` + propName + ", val})\n")
				}
				WC("		" + receiverName + `.` + propName + " = val\n")
				WC("		return true\n")
				WC("	}\n")
				WC("	return false\n")
			} else { // always overwrite for array
				WC("	" + receiverName + ".mutations.Assign(" + I.ToStr(idx) + ", val)\n")
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
		WC(`func (` + receiverName + ` *` + structName + "Mutator) SetAll(from " + rqPkgName + `.` + structName + ", excludeMap, forceMap M.SB) (changed bool) { //nolint:dupl false positive\n")
		WC("	if excludeMap == nil { // list of fields to exclude\n")
		WC("		excludeMap = M.SB{}\n")
		WC("	}\n")
		WC("	if forceMap == nil { // list of fields to force overwrite\n")
		WC("		forceMap = M.SB{}\n")
		WC("	}\n")
		for _, prop := range props.Fields {
			pascalPropName := S.PascalCase(prop.Name)

			// index functions
			WC("	if !excludeMap[`" + prop.Name + "`] && (forceMap[`" + prop.Name + "`] || from." + pascalPropName + ` != ` + TypeToGoNilValue[prop.Type] + ") {\n")
			if propByName[prop.Name].Type == String {
				WC(`		` + receiverName + `.` + pascalPropName + ` = S.Trim(from.` + pascalPropName + ")\n")
			} else {
				WC(`		` + receiverName + `.` + pascalPropName + ` = from.` + pascalPropName + "\n")
			}
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
				propType := propByName[propName].Type
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
		RQ(`func (` + receiverName + ` *` + structName + `) FromArray(ax A.X) *` + structName + " { //nolint:dupl false positive\n")
		for idx, prop := range props.Fields {
			RQ("	" + receiverName + "." + S.PascalCase(prop.Name) + ` = ` + TypeToConvertFunc[prop.Type] + "(ax[" + X.ToS(idx) + "])\n")
		}
		RQ("	return " + receiverName + NL)
		RQ("}\n\n")

		// from AX but uncensored
		RQ("// FromUncensoredArray convert slice to receiver fields\n")
		RQ(`func (` + receiverName + ` *` + structName + `) FromUncensoredArray(ax A.X) *` + structName + " { //nolint:dupl false positive\n")
		for idx, prop := range props.Fields {
			if !censoredFieldsByName[prop.Name] {
				RQ("	" + receiverName + "." + S.PascalCase(prop.Name) + ` = ` + TypeToConvertFunc[prop.Type] + "(ax[" + X.ToS(idx) + "])\n")
			}
		}
		RQ("	return " + receiverName + NL)
		RQ("}\n\n")

		// find many
		RQ("// FindOffsetLimit returns slice of struct, order by idx, eg. .UniqueIndex*()\n")
		RQ(`func (` + receiverName + ` *` + structName + ") FindOffsetLimit(offset, limit uint32, idx string) []" + structName + " { //nolint:dupl false positive\n")
		RQ("	var rows []" + structName + NL)
		RQ("	res, err := " + receiverName + ".Adapter.RetryDo(\n")
		RQ("		tarantool.NewSelectRequest(" + receiverName + ".SpaceName()).\n")
		RQ("		Index(idx).\n")
		RQ("		Offset(offset).\n")
		RQ("		Limit(limit).\n")
		RQ("		Iterator(" + iterAll + "),\n")
		RQ("	)\n")
		RQ("	if L.IsError(err, `" + structName + ".FindOffsetLimit failed: `+" + receiverName + ".SpaceName()) {\n")
		RQ("		return rows\n")
		RQ("	}\n")
		RQ("	for _, row := range res {\n")
		RQ("		item := " + structName + "{}\n")
		RQ("		row, ok := row.([]any)\n")
		RQ("		if ok {\n")
		RQ("			rows = append(rows, *item.FromArray(row))\n")
		RQ("		}\n")
		RQ("	}\n")
		RQ("	return rows\n")
		RQ("}\n\n")
		// find many

		RQ("// FindArrOffsetLimit returns as slice of slice order by idx eg. .UniqueIndex*()\n")
		RQ(`func (` + receiverName + ` *` + structName + ") FindArrOffsetLimit(offset, limit uint32, idx string) ([]A.X, Tt.QueryMeta) { //nolint:dupl false positive\n")
		RQ("	var rows []A.X" + NL)
		RQ("	resp, err := " + receiverName + ".Adapter.RetryDoResp(\n")
		RQ("		tarantool.NewSelectRequest(" + receiverName + ".SpaceName()).\n")
		RQ("		Index(idx).\n")
		RQ("		Offset(offset).\n")
		RQ("		Limit(limit).\n")
		RQ("		Iterator(" + iterAll + "),\n")
		RQ("	)\n")
		RQ("	if L.IsError(err, `" + structName + ".FindOffsetLimit failed: `+" + receiverName + ".SpaceName()) {\n")
		RQ("		return rows, Tt.QueryMetaFrom(resp, err)\n")
		RQ("	}\n")
		RQ("	res, err := resp.Decode()\n")
		RQ("	if L.IsError(err, `" + structName + ".FindOffsetLimit failed: `+" + receiverName + ".SpaceName()) {\n")
		RQ("		return rows, Tt.QueryMetaFrom(resp, err)\n")
		RQ("	}\n")
		RQ("	rows = make([]A.X, len(res))\n")
		RQ("	for _, row := range res {\n")
		RQ("		row, ok := row.([]any)\n")
		RQ("		if ok {\n")
		RQ("			rows = append(rows, row)\n")
		RQ("		}\n")
		RQ("	}\n")
		RQ("	return rows, Tt.QueryMetaFrom(resp, nil)\n")
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

		uniqueMethodSuffixes := []string{}
		if props.AutoIncrementId {
			uniqueMethodSuffixes = append(uniqueMethodSuffixes, S.PascalCase(IdCol))
		}
		if props.Unique1 != `` && !(props.AutoIncrementId && props.Unique1 == IdCol) {
			uniqueMethodSuffixes = append(uniqueMethodSuffixes, S.PascalCase(props.Unique1))
		}
		if props.Unique2 != `` && !(props.AutoIncrementId && props.Unique2 == IdCol) {
			uniqueMethodSuffixes = append(uniqueMethodSuffixes, S.PascalCase(props.Unique2))
		}
		if props.Unique3 != `` && !(props.AutoIncrementId && props.Unique3 == IdCol) {
			uniqueMethodSuffixes = append(uniqueMethodSuffixes, S.PascalCase(props.Unique3))
		}
		if len(props.Uniques) > 0 {
			uniquePropCamel := ``
			for _, uniq := range props.Uniques {
				uniquePropCamel += S.PascalCase(uniq)
			}
			uniqueMethodSuffixes = append(uniqueMethodSuffixes, uniquePropCamel)
		}

		lookupSuffix := ``
		if len(uniqueMethodSuffixes) > 0 {
			lookupSuffix = uniqueMethodSuffixes[0]
		}
		uniqueFieldMap := map[string]bool{}
		if props.Unique1 != `` {
			uniqueFieldMap[props.Unique1] = true
		}
		if props.Unique2 != `` {
			uniqueFieldMap[props.Unique2] = true
		}
		if props.Unique3 != `` {
			uniqueFieldMap[props.Unique3] = true
		}
		for _, uniq := range props.Uniques {
			uniqueFieldMap[uniq] = true
		}
		if props.AutoIncrementId {
			uniqueFieldMap[IdCol] = true
		}
		updateField := Field{}
		hasUpdateField := false
		for _, field := range props.Fields {
			if field.Name == IdCol || uniqueFieldMap[field.Name] {
				continue
			}
			updateField = field
			hasUpdateField = true
			break
		}

		// reader/unit tests (replacement of rq*_orm_test.go)
		RQT(`func TestGenerated` + structName + `OrmHelpers(t *testing.T) {` + "\n")
		RQT(`	q := New` + structName + `(nil)` + "\n")
		RQT(`	assert.NotNil(t, q)` + "\n")
		RQT(`	assert.NotEmpty(t, q.SpaceName())` + "\n")
		RQT(`	assert.NotEmpty(t, q.SqlTableName())` + "\n")
		for _, field := range props.Fields {
			camel := S.PascalCase(field.Name)
			RQT(`	q.` + camel + ` = ` + testSampleLiteral(field.Type) + "\n")
		}
		RQT(`	arr := q.ToArray()` + "\n")
		RQT(`	assert.Len(t, arr, ` + X.ToS(len(props.Fields)) + `)` + "\n")
		RQT(`	assert.NotNil(t, q.ToUpdateArray())` + "\n")
		RQT(`	decoded := (&` + structName + `{}).FromArray(arr)` + "\n")
		for _, field := range props.Fields {
			camel := S.PascalCase(field.Name)
			RQT(`	assert.Equal(t, q.` + camel + `, decoded.` + camel + `)` + "\n")
		}
		RQT(`	decoded2 := (&` + structName + `{}).FromUncensoredArray(arr)` + "\n")
		censoredMap := map[string]bool{}
		for _, f := range props.AutoCensorFields {
			censoredMap[f] = true
		}
		for _, field := range props.Fields {
			camel := S.PascalCase(field.Name)
			if censoredMap[field.Name] {
				RQT(`	assert.Equal(t, ` + TypeToGoNilValue[field.Type] + `, decoded2.` + camel + `)` + "\n")
			} else {
				RQT(`	assert.Equal(t, q.` + camel + `, decoded2.` + camel + `)` + "\n")
			}
		}
		for idx, field := range props.Fields {
			camel := S.PascalCase(field.Name)
			RQT(`	assert.Equal(t, ` + X.ToS(idx) + `, q.Idx` + camel + `())` + "\n")
			RQT(`	assert.Equal(t, ` + S.BT(dq(field.Name)) + `, q.Sql` + camel + `())` + "\n")
		}
		if len(props.Fields) > 0 {
			RQT(`	_, ok := ` + structName + `FieldTypeMap[` + S.BT(props.Fields[0].Name) + `]` + "\n")
			RQT(`	assert.True(t, ok)` + "\n")
		}
		for _, suffix := range uniqueMethodSuffixes {
			RQT(`	assert.NotEmpty(t, q.UniqueIndex` + suffix + `())` + "\n")
		}
		if props.Spatial != `` {
			spatial := S.PascalCase(props.Spatial)
			RQT(`	assert.Equal(t, ` + S.BT(props.Spatial) + `, q.SpatialIndex` + spatial + `())` + "\n")
		}
		if len(props.AutoCensorFields) > 0 {
			for _, field := range props.AutoCensorFields {
				camel := S.PascalCase(field)
				RQT(`	q.` + camel + ` = ` + testSampleLiteral(propByName[field].Type) + "\n")
			}
			RQT(`	q.CensorFields()` + "\n")
			for _, field := range props.AutoCensorFields {
				camel := S.PascalCase(field)
				RQT(`	assert.Equal(t, ` + TypeToGoEmptyValue[propByName[field].Type] + `, q.` + camel + `)` + "\n")
			}
		}
		RQT(`}` + "\n\n")

		RQT(`func TestGenerated` + structName + `DbMethodsPanic(t *testing.T) {` + "\n")
		RQT(`	q := New` + structName + `(&Tt.Adapter{})` + "\n")
		for _, field := range props.Fields {
			camel := S.PascalCase(field.Name)
			RQT(`	q.` + camel + ` = ` + testSampleLiteral(field.Type) + "\n")
		}
		for _, suffix := range uniqueMethodSuffixes {
			RQT(`	assert.Panics(t, func() { _ = q.FindBy` + suffix + `() })` + "\n")
		}
		RQT(`	assert.Panics(t, func() { _ = q.FindOffsetLimit(0, 1, "") })` + "\n")
		RQT(`	assert.Panics(t, func() { _, _ = q.FindArrOffsetLimit(0, 1, "") })` + "\n")
		RQT(`	assert.Panics(t, func() { _ = q.Total() })` + "\n")
		RQT(`}` + "\n\n")

		// mutator unit tests (replacement of wc*_unit_test.go)
		WCT(`func TestGenerated` + structName + `Unit(t *testing.T) {` + "\n")
		WCT(`	m := New` + structName + `Mutator(nil)` + "\n")
		WCT(`	assert.NotNil(t, m)` + "\n")
		WCT(`	assert.False(t, m.HaveMutation())` + "\n")
		WCT(`	assert.Empty(t, m.Logs())` + "\n")
		for _, field := range props.Fields {
			camel := S.PascalCase(field.Name)
			value := testSampleLiteral(field.Type)
			WCT(`	assert.True(t, m.Set` + camel + `(` + value + `))` + "\n")
			if field.Type != Array {
				WCT(`	assert.False(t, m.Set` + camel + `(` + value + `))` + "\n")
			}
		}
		WCT(`	assert.True(t, m.HaveMutation())` + "\n")
		WCT(`	from := m.` + structName + "\n")
		WCT(`	assert.True(t, m.SetAll(from, nil, nil))` + "\n")
		WCT(`	m2 := New` + structName + `Mutator(nil)` + "\n")
		WCT(`	fromZero := m2.` + structName + "\n")
		for _, field := range props.Fields {
			if field.Type == Array {
				WCT(`	fromZero.` + S.PascalCase(field.Name) + ` = nil` + "\n")
			}
		}
		WCT(`	assert.False(t, m.SetAll(fromZero, nil, nil))` + "\n")
		WCT(`	m.ClearMutations()` + "\n")
		WCT(`	assert.False(t, m.HaveMutation())` + "\n")
		for _, suffix := range uniqueMethodSuffixes {
			WCT(`	assert.True(t, m.DoUpdateBy` + suffix + `())` + "\n")
		}
		WCT(`}` + "\n\n")

		// mutator integration tests (replacement of wc*_test.go)
		tableConst := mPkgName + `.Table` + structName
		WCT(`func TestGenerated` + structName + `CRUD(t *testing.T) {` + "\n")
		WCT(`	if dbConn == nil {` + "\n")
		WCT(`		t.Skip("docker unavailable")` + "\n")
		WCT(`	}` + "\n")
		WCT(`	a := &Tt.Adapter{Connection: dbConn, Reconnect: reconnect}` + "\n")
		WCT(`	ok := a.UpsertTable(` + tableConst + `, ` + mPkgName + `.TarantoolTables[` + tableConst + `])` + "\n")
		WCT(`	assert.True(t, ok)` + "\n")
		WCT(`	_ = a.TruncateTable(string(` + tableConst + `))` + "\n")
		WCT(`	seed := func() *` + structName + `Mutator {` + "\n")
		WCT(`		x := New` + structName + `Mutator(a)` + "\n")
		for _, field := range props.Fields {
			if props.AutoIncrementId && field.Name == IdCol {
				continue
			}
			camel := S.PascalCase(field.Name)
			WCT(`		x.` + camel + ` = ` + testSampleLiteral(field.Type) + "\n")
		}
		WCT(`		assert.True(t, x.DoInsert())` + "\n")
		if props.AutoIncrementId {
			WCT(`		assert.Greater(t, x.Id, uint64(0))` + "\n")
		}
		WCT(`		return x` + "\n")
		WCT(`	}` + "\n")
		if lookupSuffix != `` {
			for idx, suffix := range uniqueMethodSuffixes {
				recVar := `rec` + X.ToS(idx)
				rowsVar := `rows` + X.ToS(idx)
				arrRowsVar := `arrRows` + X.ToS(idx)
				WCT(`	` + recVar + ` := seed()` + "\n")
				WCT(`	assert.True(t, ` + recVar + `.FindBy` + suffix + `())` + "\n")
				WCT(`	assert.True(t, ` + recVar + `.DoUpdateBy` + suffix + `())` + "\n")
				if hasUpdateField {
					updCamel := S.PascalCase(updateField.Name)
					WCT(`	assert.True(t, ` + recVar + `.Set` + updCamel + `(` + testSampleLiteralAlt(updateField.Type) + `))` + "\n")
					WCT(`	assert.True(t, ` + recVar + `.DoUpdateBy` + suffix + `())` + "\n")
				}
				WCT(`	_ = ` + recVar + `.DoOverwriteBy` + suffix + `()` + "\n")
				WCT(`	assert.True(t, ` + recVar + `.FindBy` + suffix + `())` + "\n")
				WCT(`	` + rowsVar + ` := ` + recVar + `.FindOffsetLimit(0, 10, ` + recVar + `.UniqueIndex` + suffix + `())` + "\n")
				WCT(`	assert.NotNil(t, ` + rowsVar + `)` + "\n")
				WCT(`	` + arrRowsVar + `, _ := ` + recVar + `.FindArrOffsetLimit(0, 10, ` + recVar + `.UniqueIndex` + suffix + `())` + "\n")
				WCT(`	assert.NotNil(t, ` + arrRowsVar + `)` + "\n")
				WCT(`	assert.GreaterOrEqual(t, ` + recVar + `.Total(), int64(0))` + "\n")
				WCT(`	assert.True(t, ` + recVar + `.DoDeletePermanentBy` + suffix + `())` + "\n")
				WCT(`	assert.False(t, ` + recVar + `.FindBy` + suffix + `())` + "\n")
			}
		}
		if props.AutoIncrementId {
			WCT(`	u := New` + structName + `Mutator(a)` + "\n")
			for _, field := range props.Fields {
				if field.Name == IdCol {
					continue
				}
				camel := S.PascalCase(field.Name)
				WCT(`	u.` + camel + ` = ` + testSampleLiteral(field.Type) + "\n")
			}
			WCT(`	assert.True(t, u.DoUpsertById())` + "\n")
			WCT(`	assert.Greater(t, u.Id, uint64(0))` + "\n")
			if hasUpdateField {
				updCamel := S.PascalCase(updateField.Name)
				WCT(`	assert.True(t, u.Set` + updCamel + `(` + testSampleLiteralAlt(updateField.Type) + `))` + "\n")
				WCT(`	assert.True(t, u.DoUpsertById())` + "\n")
			}
			WCT(`	assert.True(t, u.DoDeletePermanentById())` + "\n")
		}
		WCT(`}` + "\n\n")

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

	rqTestFname := fmt.Sprintf(`./%s/%s__ORM.GEN_test.go`, rqPkgName, rqPkgName)
	err = os.WriteFile(rqTestFname, rqTestBuf.Bytes(), os.ModePerm)
	if L.IsError(err, `os.WriteFile failed: `+rqTestFname) {
		return
	}

	wcTestFname := fmt.Sprintf(`./%s/%s__ORM.GEN_test.go`, wcPkgName, wcPkgName)
	err = os.WriteFile(wcTestFname, wcTestBuf.Bytes(), os.ModePerm)
	if L.IsError(err, `os.WriteFile failed: `+wcTestFname) {
		return
	}
}

func testSampleLiteral(typ DataType) string {
	switch typ {
	case Unsigned:
		return `uint64(1)`
	case String:
		return `"sample"`
	case Integer:
		return `int64(1)`
	case Double:
		return `1.5`
	case Boolean:
		return `true`
	case Array:
		return `[]any{1.1, 2.2}`
	default:
		return `nil`
	}
}

func testSampleLiteralAlt(typ DataType) string {
	switch typ {
	case Unsigned:
		return `uint64(2)`
	case String:
		return `"sample2"`
	case Integer:
		return `int64(2)`
	case Double:
		return `2.5`
	case Boolean:
		return `false`
	case Array:
		return `[]any{3.3, 4.4}`
	default:
		return `nil`
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

func generateMutationByUniqueIndex(uniqueCamel, structProp, receiverName, structName string, keyFunc func(string) string, RQ, WC func(str string)) {

	//// primary fields
	//RQ(`func (` + receiverName + ` *` + structName + ") PrimaryIndex() A.X { //nolint:dupl false positive\n")
	//RQ(`	return A.X{` + structProp + "}\n")
	//RQ("}\n\n")

	// find by unique, used when need to mutate the object

	RQ("// FindBy" + uniqueCamel + " Find one by " + uniqueCamel + "\n")
	RQ(`func (` + receiverName + ` *` + structName + `) FindBy` + uniqueCamel + "() bool { //nolint:dupl false positive\n")
	RQ("	res, err := " + receiverName + ".Adapter.RetryDo(\n")
	RQ("		tarantool.NewSelectRequest(" + receiverName + ".SpaceName()).\n")
	RQ("		Index(" + receiverName + ".UniqueIndex" + uniqueCamel + "()).\n")
	RQ("		Limit(1).\n")
	RQ("		Iterator(" + iterEq + ").\n")
	RQ("		Key(" + keyFunc(structProp) + "),\n")
	RQ("	)\n")
	RQ("	if L.IsError(err, `" + structName + `.FindBy` + uniqueCamel + " failed: `+" + receiverName + ".SpaceName()) {\n")
	RQ("		return false\n")
	RQ("	}\n")
	RQ("	if len(res) == 1 {\n")
	RQ("		if row, ok := res[0].([]any); ok {\n")
	RQ("			" + receiverName + ".FromArray(row)\n")
	RQ("			return true\n")
	RQ("		}\n")
	RQ("	}\n")
	RQ("	return false\n")
	RQ("}\n\n")

	// Overwrite all columns, error if not exists
	WC("// DoOverwriteBy" + uniqueCamel + " update all columns, error if not exists, not using mutations/Set*\n")
	WC(`func (` + receiverName + ` *` + structName + `Mutator) DoOverwriteBy` + uniqueCamel + "() bool { //nolint:dupl false positive\n")
	WC("	_, err := " + receiverName + ".Adapter.RetryDo(tarantool.NewUpdateRequest(" + receiverName + ".SpaceName()).\n")
	WC("		Index(" + receiverName + ".UniqueIndex" + uniqueCamel + "()).\n")
	WC("		Key(" + keyFunc(structProp) + ").\n")
	WC("		Operations(" + receiverName + ".ToUpdateArray()),\n")
	WC("	)\n")
	WC("	return !L.IsError(err, `" + structName + `.DoOverwriteBy` + uniqueCamel + " failed: `+" + receiverName + ".SpaceName())\n")
	WC("}\n\n")

	// Update only mutated, error if not exists
	WC("// DoUpdateBy" + uniqueCamel + " update only mutated fields, error if not exists, use Find* and Set* methods instead of direct assignment\n")
	WC(`func (` + receiverName + ` *` + structName + `Mutator) DoUpdateBy` + uniqueCamel + "() bool { //nolint:dupl false positive\n")
	WC(`	if !` + receiverName + ".HaveMutation() {\n")
	WC("		return true\n")
	WC("	}\n")
	WC("	_, err := " + receiverName + ".Adapter.RetryDo(\n")
	WC("		tarantool.NewUpdateRequest(" + receiverName + ".SpaceName()).\n")
	WC("		Index(" + receiverName + ".UniqueIndex" + uniqueCamel + "()).\n")
	WC("		Key(" + keyFunc(structProp) + ").\n")
	WC("		Operations(" + receiverName + ".mutations),\n")
	WC("	)\n")
	WC("	return !L.IsError(err, `" + structName + `.DoUpdateBy` + uniqueCamel + " failed: `+" + receiverName + ".SpaceName())\n")
	WC("}\n\n")

	// permanent delete
	WC("// DoDeletePermanentBy" + uniqueCamel + " permanent delete\n")
	WC(`func (` + receiverName + ` *` + structName + `Mutator) DoDeletePermanentBy` + uniqueCamel + "() bool { //nolint:dupl false positive\n")
	WC("	_, err := " + receiverName + ".Adapter.RetryDo(\n")
	WC("		tarantool.NewDeleteRequest(" + receiverName + ".SpaceName()).\n")
	WC("		Index(" + receiverName + ".UniqueIndex" + uniqueCamel + "()).\n")
	WC("		Key(" + keyFunc(structProp) + "),\n")
	WC("	)\n")
	WC("	return !L.IsError(err, `" + structName + `.DoDeletePermanentBy` + uniqueCamel + " failed: `+" + receiverName + ".SpaceName())\n")
	WC("}\n\n")

}
