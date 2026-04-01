package W2

import (
	"bytes"
	"go/ast"
	"go/parser"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/kokizzu/gotro/M"
)

func TestCodegenHelpers(t *testing.T) {
	sf := StructField{Name: "UserId", Type: "int64"}
	if sf.LowerName() != "userId" {
		t.Fatalf("LowerName mismatch: %q", sf.LowerName())
	}
	if sf.ApiComment() != "// int64" {
		t.Fatalf("ApiComment mismatch: %q", sf.ApiComment())
	}
	if (StructField{}).ApiComment() != "" {
		t.Fatalf("ApiComment for empty field should be empty")
	}
	if (StructField{Type: "int", Comment: "some comment"}).ApiComment() != "//int | some comment" {
		t.Fatalf("ApiComment with comment mismatch")
	}

	ep := ErrPair{Code: 400, Msg: `"bad"`}
	if ep.String() != `400"bad"` {
		t.Fatalf("ErrPair.String mismatch: %q", ep.String())
	}
	var em *ErrMap
	em = em.Add(ep)
	em = em.Add(ep) // duplicate key overwrite
	if len(em.List) != 1 {
		t.Fatalf("ErrMap.Add duplicate should not increase size: %#v", em.List)
	}
}

func TestCallListHelpers(t *testing.T) {
	var nilList *CallList
	if len(nilList.SortedKeys()) != 0 {
		t.Fatalf("nil CallList should have zero keys")
	}

	var cl *CallList
	cl = cl.Add("B")
	cl = cl.Add("A")
	cl.IncCallCount("A")
	cl.ZeroCallCount("B")
	keys := cl.SortedKeys()
	if !reflect.DeepEqual(keys, []string{"A", "B"}) {
		t.Fatalf("CallList keys mismatch: %#v", keys)
	}
	if cl.List["A"].CallCount != 2 || cl.List["B"].CallCount != 0 {
		t.Fatalf("CallList call counts mismatch: %#v", cl.List)
	}

	expr, err := parser.ParseExpr(`saAuth.LoginStat{}`)
	if err != nil {
		t.Fatalf("ParseExpr failed: %v", err)
	}
	cl.AddCallParam([]ast.Expr{expr})
	if cl.List["saAuth.LoginStat"] == nil || cl.List["saAuth.LoginStat"].CallCount != 1 {
		t.Fatalf("AddCallParam mismatch: %#v", cl.List)
	}

	if (&CallParam{FirstParam: "x.y"}).ToString() != "x.y" {
		t.Fatalf("CallParam.ToString mismatch")
	}
}

func TestParseStructFieldsAndSplitHelpers(t *testing.T) {
	expr, err := parser.ParseExpr(`struct{
		A string ` + "`json:\"a\"`" + `
		B *rqAuth.User
		C []int
		D []*rqAuth.User
		E []rqAuth.User
		F map[string]any
		G rqAuth.MapSlice
		H time.Time
	}`)
	if err != nil {
		t.Fatalf("ParseExpr failed: %v", err)
	}
	fields := parseStructFields(nil, expr)
	if len(fields) != 8 {
		t.Fatalf("parseStructFields count mismatch: %#v", fields)
	}
	if fields[0].Type != "string" || !strings.Contains(fields[0].Tags, `json:"a"`) {
		t.Fatalf("field A mismatch: %#v", fields[0])
	}
	if fields[1].Type != "rqAuth.User" || fields[1].IsArray || fields[1].IsMap {
		t.Fatalf("field B mismatch: %#v", fields[1])
	}
	if fields[2].Type != "int" || !fields[2].IsArray {
		t.Fatalf("field C mismatch: %#v", fields[2])
	}
	if fields[3].Type != "rqAuth.User" || !fields[3].IsArray {
		t.Fatalf("field D mismatch: %#v", fields[3])
	}
	if fields[4].Type != "rqAuth.User" || !fields[4].IsArray {
		t.Fatalf("field E mismatch: %#v", fields[4])
	}
	if !fields[5].IsMap {
		t.Fatalf("field F should be map: %#v", fields[5])
	}
	if fields[6].Type != "rqAuth.MapSlice" || !fields[6].IsMap {
		t.Fatalf("field G mismatch: %#v", fields[6])
	}
	if fields[7].Type != "time.Time" || fields[7].IsMap {
		t.Fatalf("field H mismatch: %#v", fields[7])
	}

	params, idx := unquoteAndSplit(`"/v1/:userId/:encodedId"`)
	if !reflect.DeepEqual(params, []string{"userId", "encodedId"}) {
		t.Fatalf("unquoteAndSplit params mismatch: %#v", params)
	}
	if idx["userId"] != 2 || idx["encodedId"] != 3 {
		t.Fatalf("unquoteAndSplit idx mismatch: %#v", idx)
	}

	if q("abc") != `"abc"` {
		t.Fatalf("q mismatch")
	}
	if lowerFirstLetter("") != "" || lowerFirstLetter("UserId") != "userId" {
		t.Fatalf("lowerFirstLetter mismatch")
	}
}

func TestWriteUrlAndGraphqlHelpers(t *testing.T) {
	buf := bytes.Buffer{}
	writeUrlSegmentToIn(&buf, []string{"userId", "encodedId", "slug"}, `ctx.Params(%s)`)
	out := buf.String()
	if !strings.Contains(out, `in.UserId = S.ToU(ctx.Params("userId"))`) {
		t.Fatalf("writeUrlSegmentToIn userId conversion missing:\n%s", out)
	}
	if !strings.Contains(out, `in.EncodedId = ctx.Params("encodedId")`) {
		t.Fatalf("writeUrlSegmentToIn encodedId conversion mismatch:\n%s", out)
	}
	if !strings.Contains(out, `in.Slug = ctx.Params("slug")`) {
		t.Fatalf("writeUrlSegmentToIn slug conversion mismatch:\n%s", out)
	}

	r := &RoutesArgs{ProjectName: "example/project", funcStatistics: map[string]*CallList{}}
	imports := M.SB{}
	types := M.SB{}

	sfIntID := StructField{Name: "CreatedBy", Type: "Int"}
	if got := r.graphqlType(&sfIntID, imports, types); got != "graphql.ID, // Int" {
		t.Fatalf("graphqlType Int/Id mismatch: %q", got)
	}

	sfCustomArr := StructField{Name: "Users", Type: "rqAuth.User", IsArray: true}
	gotCustom := r.graphqlType(&sfCustomArr, imports, types)
	if gotCustom != "graphql.NewList(rqAuth.GraphqlTypeUser), //  []rqAuth.User" {
		t.Fatalf("graphqlType custom array mismatch: %q", gotCustom)
	}
	if !imports["example/project/model/mAuth/rqAuth"] {
		t.Fatalf("graphqlType should add model import: %#v", imports)
	}

	sfUnknown := StructField{Name: "Unknown", Type: "Whatever"}
	if got := r.graphqlType(&sfUnknown, imports, types); got != "graphql.String, // Whatever" {
		t.Fatalf("graphqlType unknown mismatch: %q", got)
	}
	if r.graphqlTypeComment(&StructField{Type: "int"}) != "int" {
		t.Fatalf("graphqlTypeComment scalar mismatch")
	}
	if r.graphqlTypeComment(&StructField{Type: "int", IsArray: true}) != " []int" {
		t.Fatalf("graphqlTypeComment array mismatch")
	}
	if r.graphqlTypeComment(&StructField{Type: "int", IsMap: true}) != " map[?]int" {
		t.Fatalf("graphqlTypeComment map mismatch")
	}

	r.IncStatisticsCalls("Login", &CallParam{FirstParam: "saAuth.LoginStat", CallCount: 3})
	if r.funcStatistics["Login"] == nil || r.funcStatistics["Login"].List["saAuth.LoginStat"] == nil {
		t.Fatalf("IncStatisticsCalls mismatch: %#v", r.funcStatistics)
	}
}

func TestGenerateFiberCliApiDocsAndGraphql(t *testing.T) {
	tmp := t.TempDir()
	domainDir := filepath.Join(tmp, "domain")
	modelRqDir := filepath.Join(tmp, "model", "rqAuth")
	modelSaDir := filepath.Join(tmp, "model", "saAuth")
	svelteDir := filepath.Join(tmp, "svelte", "src", "pages")
	if err := os.MkdirAll(domainDir, 0755); err != nil {
		t.Fatalf("mkdir domain failed: %v", err)
	}
	if err := os.MkdirAll(modelRqDir, 0755); err != nil {
		t.Fatalf("mkdir model rq failed: %v", err)
	}
	if err := os.MkdirAll(modelSaDir, 0755); err != nil {
		t.Fatalf("mkdir model sa failed: %v", err)
	}
	if err := os.MkdirAll(svelteDir, 0755); err != nil {
		t.Fatalf("mkdir svelte failed: %v", err)
	}

	domainSrc := `package dAuth

const (
	Login_Url = "/login/:userId/:encodedId"
	ErrConst = "invalid"
)

type Login_In struct {
	UserId uint64
	Name string ` + "`json:\",string\"`" + `
	Profile rqAuth.User
}

type Login_Out struct {
	SessionToken string
	Ok bool
}

type Domain struct {}

func (d *Domain) mustLogin() {}
func (d *Domain) helper() {}

func (d *Domain) Login(in *Login_In) (out Login_Out) {
	d.mustLogin()
	d.helper()
	out.SetError(400, "bad")
	out.SetError(401, ErrConst)
	out.SetError(402, "bad" + ErrConst)
	rqAuth.NewFoo()
	wcAuth.NewBarMutator()
	metrics.Statistics(saAuth.LoginStat{})
	external.Call()
	out.SessionToken = "tok"
	return
}
`
	if err := os.WriteFile(filepath.Join(domainDir, "login.go"), []byte(domainSrc), 0644); err != nil {
		t.Fatalf("write domain file failed: %v", err)
	}

	modelRqSrc := `package rqAuth
type User struct {
	Name string
	Age int64
}
`
	if err := os.WriteFile(filepath.Join(modelRqDir, "user.go"), []byte(modelRqSrc), 0644); err != nil {
		t.Fatalf("write model rq file failed: %v", err)
	}
	modelSaSrc := `package saAuth
type LoginStat struct {
	Count int64
}
`
	if err := os.WriteFile(filepath.Join(modelSaDir, "stat.go"), []byte(modelSaSrc), 0644); err != nil {
		t.Fatalf("write model sa file failed: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}
	defer func() {
		_ = os.Chdir(wd)
	}()
	if err := os.Chdir(domainDir); err != nil {
		t.Fatalf("Chdir failed: %v", err)
	}

	cfg := &GeneratorConfig{
		ProjectName:  "example/project",
		GenGraphQl:   true,
		ThirdParties: []string{"external"},
	}
	GenerateFiberAndCli(cfg)
	GenerateApiDocs(cfg)

	webRoutesPath := filepath.Join(tmp, "main_restApi_routes.GEN.go")
	webGraphqlPath := filepath.Join(tmp, "main_restApi_graphql.GEN.go")
	cliArgsPath := filepath.Join(tmp, "main_cli_args.GEN.go")
	apiDocsPath := filepath.Join(tmp, "svelte", "src", "pages", "api.js")

	webRoutes, err := os.ReadFile(webRoutesPath)
	if err != nil {
		t.Fatalf("read web routes failed: %v", err)
	}
	webRoutesStr := string(webRoutes)
	if !strings.Contains(webRoutesStr, `in.UserId = S.ToU(ctx.Params("userId"))`) ||
		!strings.Contains(webRoutesStr, `in.EncodedId = ctx.Params("encodedId")`) {
		t.Fatalf("web routes content mismatch:\n%s", webRoutesStr)
	}

	cliArgs, err := os.ReadFile(cliArgsPath)
	if err != nil {
		t.Fatalf("read cli args failed: %v", err)
	}
	cliArgsStr := string(cliArgs)
	if !strings.Contains(cliArgsStr, `cliSegmentFromIdx(args[0],patterns[pattern]["userId"])`) {
		t.Fatalf("cli args content mismatch:\n%s", cliArgsStr)
	}

	apiDocs, err := os.ReadFile(apiDocsPath)
	if err != nil {
		t.Fatalf("read api docs failed: %v", err)
	}
	apiDocsStr := string(apiDocs)
	if !strings.Contains(apiDocsStr, `Login:`) ||
		!strings.Contains(apiDocsStr, `"Auth.Foo"`) ||
		!strings.Contains(apiDocsStr, `"Auth.Bar"`) ||
		!strings.Contains(apiDocsStr, `"Auth.LoginStat"`) ||
		!strings.Contains(apiDocsStr, `"external"`) ||
		!strings.Contains(apiDocsStr, `[400, "bad"]`) {
		t.Fatalf("api docs content mismatch:\n%s", apiDocsStr)
	}

	webGraphql, err := os.ReadFile(webGraphqlPath)
	if err != nil {
		t.Fatalf("read graphql failed: %v", err)
	}
	webGraphqlStr := string(webGraphql)
	if !strings.Contains(webGraphqlStr, `graphqlTypeLoginOut`) ||
		!strings.Contains(webGraphqlStr, "`Login`: &graphql.Field") ||
		!strings.Contains(webGraphqlStr, "`example/project/model/mAuth/rqAuth`") {
		t.Fatalf("graphql content mismatch:\n%s", webGraphqlStr)
	}
}
