package W2

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kpango/fastime"
)

type GeneratorConfig struct {
	ProjectName string // must equal go.mod header

	ModelPath     string // model directory
	WebRoutesFile string // fiber web route generated file
	CliArgsFile   string // cli args handler generated file
	ApiDocsFile   string // apidocs generated file

	ThirdParties []string
}

func GenerateFiberAndCli(c *GeneratorConfig) {
	if c.ModelPath == `` {
		c.ModelPath = `../model`
	}
	if c.WebRoutesFile == `` {
		c.WebRoutesFile = `../main_restApi_routes.GEN.go`
	}
	if c.CliArgsFile == `` {
		c.CliArgsFile = `../main_cli_args.GEN.go`
	}
	r := c.ParseRoutes(false)

	r.WriteWebRoutes(c.WebRoutesFile)

	r.WriteCliArgs(c.CliArgsFile)
}

func GenerateApiDocs(c *GeneratorConfig) {
	if c.ModelPath == `` {
		c.ModelPath = `../model`
	}
	if c.ApiDocsFile == `` {
		c.ApiDocsFile = `../svelte/src/pages/api.js`
	}
	r := c.ParseRoutes(true)

	r.WriteApiDocs(c.ApiDocsFile)

	// TODO: generate flutter client
	// TODO: generate javascript client (for Svelte)
	// TODO: generate unity client
}

//go:generate go test -run=XXX -bench=Benchmark_Generate_WebApiRoutes_CliArgs

func (c *GeneratorConfig) ParseRoutes(parseModelAndCalls bool) *RoutesArgs {
	r := &RoutesArgs{ProjectName: c.ProjectName, ThirdParties: c.ThirdParties}
	r.parseModel = parseModelAndCalls
	r.parseCalls = parseModelAndCalls

	if r.parseModel {
		r.models = &ModelArgs{}
		err := filepath.Walk(c.ModelPath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if S.EndsWith(path, `.go`) &&
					!S.EndsWith(path, `_test.go`) {
					r.models.ParseModel(path)
				}
				return nil
			})
		L.PanicIf(err, `failed parsing model`)
	}

	err := filepath.Walk(`.`,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if S.EndsWith(path, `.go`) &&
				!S.EndsWith(path, `_test.go`) {
				r.ParseDomain(path)
			}
			return nil
		})
	L.PanicIf(err, `failed parsing domain`)
	return r
}

type StructField struct {
	Name    string
	Type    string
	Tags    string
	Comment string
	IsArray bool
	IsMap   bool
}

func (f StructField) LowerName() string {
	return lowerFirstLetter(f.Name)
}

func (f StructField) ApiComment() string {
	if f.Comment == `` {
		if f.Type == `` {
			return ``
		}
		return `// ` + f.Type
	}
	return `//` + f.Type + ` | ` + f.Comment
}

type ModelArgs struct {
	// model
	lastStruct    string
	lastPackage   string
	fieldsByModel map[string][]StructField
}

type ErrPair struct {
	Code int
	Msg  string
}

func (e ErrPair) String() string {
	return fmt.Sprint(e.Code) + e.Msg
}

type ErrMap struct {
	List map[string]ErrPair
}

func (m *ErrMap) Add(pair ErrPair) *ErrMap {
	if m == nil {
		m = &ErrMap{List: map[string]ErrPair{}}
	}
	if m.List == nil {
		m.List = map[string]ErrPair{}
	}
	m.List[pair.String()] = pair
	return m
}

type CallParam struct {
	CallCount  int
	FirstParam string
}

func (p *CallParam) ToString() string {
	return p.FirstParam
}

type CallList struct {
	List map[string]*CallParam
}

func (c *CallList) SortedKeys() []string {
	res := []string{}
	for key := range c.List {
		res = append(res, key)
	}
	sort.Strings(res)
	return res
}

func NewCallList(old *CallList) *CallList {
	if old == nil {
		return &CallList{List: map[string]*CallParam{}}
	}
	return old
}

func (c *CallList) IncCallCount(funName string) {
	if c.List[funName] == nil {
		c.List[funName] = &CallParam{CallCount: 0}
	}
	c.List[funName].CallCount++
}

func (c *CallList) Add(funName string) *CallList {
	c = NewCallList(c)
	c.IncCallCount(funName)
	return c
}

func (c *CallList) ZeroCallCount(funName string) {
	c = NewCallList(c)
	if c.List[funName] == nil {
		c.List[funName] = &CallParam{CallCount: 0}
	}
	c.List[funName].CallCount = 0
}

func (c *CallList) AddCallParam(args []ast.Expr) {
	c = NewCallList(c)
	cp := CallParam{}
	if sel, ok := args[0].(*ast.CompositeLit); ok {
		if typ, ok := sel.Type.(*ast.SelectorExpr); ok {
			cp.FirstParam = typ.X.(*ast.Ident).Name + `.` + typ.Sel.Name
		} else {
			L.Describe(args[0])
			panic(`unhandled args[0] Statistics 1`)
		}
	} else {
		L.Describe(args[0])
		panic(`unhandled args[0] Statistics 2`)
	}
	if len(cp.FirstParam) > 0 {
		c.IncCallCount(cp.ToString())
	}
}

type RoutesArgs struct {
	ProjectName  string
	ThirdParties []string

	// domain
	methodsPkgMap        M.SS // method:package
	packagesRefCount     M.SI
	lastPackage          string
	lastMethod           string
	methodsUrlMap        map[string][]string       // method:urlSegments
	methodsSegmentIdxMap map[string]map[string]int // method:urlSegment:idx
	inputFieldsByMethod  map[string][]StructField
	outputFieldsByMethod map[string][]StructField

	// model
	parseModel bool
	models     *ModelArgs

	// func SetError
	parseCalls      bool
	funcSetErrorMap map[string]*ErrMap
	funcCalls       map[string]*CallList
	lastFuncDecl    string
	funcNewRqList   map[string]*CallList // read query
	funcNewWcList   map[string]*CallList // writer command
	func3rdParty    map[string]bool
	funcStatistics  map[string]*CallList
}

func (r *ModelArgs) ParseModel(path string) {
	r.lastStruct = ``
	r.lastPackage = ``
	if r.fieldsByModel == nil {
		r.fieldsByModel = map[string][]StructField{}
	}
	fs := token.NewFileSet()
	src, err := ioutil.ReadFile(path)
	L.PanicIf(err, `ioutil.ReadFile failed: `+path)
	f, err := parser.ParseFile(fs, path, src, parser.AllErrors)
	L.PanicIf(err, `parser.ParseFile failed: `+path)
	ast.Walk(r, f)
}

func (r *ModelArgs) Visit(n ast.Node) (w ast.Visitor) { // parse model
	//aa := spew.Sdump("%T %#v", n, n)
	//if S.Contains(aa, `mGame`) && !S.Contains(aa, `LevelExp`) {
	//	L.Describe(n)
	//}

	if ident, ok := n.(*ast.Ident); ok {
		val := ident.Name
		if ident.NamePos < 10 {
			if S.StartsWith(val, `rq`) {
				r.lastPackage = val
				return r
			}
			if S.StartsWith(val, `wc`) {
				r.lastPackage = val
				return r
			}
			if S.StartsWith(val, `sa`) {
				r.lastPackage = val
				return r
			}
			if S.StartsWith(val, `m`) {
				r.lastPackage = val
				return r
			}
			L.Print(val)
		}
	}

	if ts, ok := n.(*ast.TypeSpec); ok {
		r.lastStruct = ts.Name.Name
		modelName := r.lastPackage + `.` + r.lastStruct
		r.fieldsByModel[modelName] = parseStructFields(r.fieldsByModel[modelName], ts.Type)
	}

	return r
}

func (r *RoutesArgs) ParseDomain(path string) {
	r.lastMethod = ``
	r.lastPackage = ``
	r.lastFuncDecl = ``
	if r.methodsPkgMap == nil {
		r.methodsPkgMap = M.SS{}
		r.packagesRefCount = M.SI{}
		r.methodsUrlMap = map[string][]string{}
		r.methodsSegmentIdxMap = map[string]map[string]int{}
		r.inputFieldsByMethod = map[string][]StructField{}
		r.outputFieldsByMethod = map[string][]StructField{}
		if r.parseCalls {
			r.funcSetErrorMap = map[string]*ErrMap{}
			r.funcNewRqList = map[string]*CallList{}
			r.funcNewWcList = map[string]*CallList{}
			r.func3rdParty = map[string]bool{}
			r.funcCalls = map[string]*CallList{}
			r.funcStatistics = map[string]*CallList{}
		}
	}
	fs := token.NewFileSet()
	src, err := ioutil.ReadFile(path)
	L.PanicIf(err, `ioutil.ReadFile failed: `+path)
	f, err := parser.ParseFile(fs, path, src, parser.AllErrors)
	L.PanicIf(err, `parser.ParseFile failed: `+path)
	ast.Walk(r, f)

	// merge recursive function calls (eg. A call B, B call C, and C have SetError, so A should have SetError too)
	if r.parseCalls {
		changed := true
		mergeSetErrror := func(pairs *ErrMap, rqCalls, calls *CallList) (*ErrMap, *CallList) {
			if pairs == nil {
				pairs = &ErrMap{}
			}
			rqCalls = NewCallList(rqCalls)
			rqLen := len(rqCalls.List)
			for fun := range calls.List {
				// merge errors
				errList, ok := r.funcSetErrorMap[fun]
				if ok {
					origLen := len(pairs.List)
					for _, err := range errList.List {
						pairs = pairs.Add(err)
					}
					if origLen != len(pairs.List) {
						changed = true
					}
				}

				// merge all calls
				rqCalls.IncCallCount(fun)
				callList, ok := r.funcCalls[fun]
				if ok && callList != nil {
					for fun := range callList.List {
						rqCalls.IncCallCount(fun)
					}
				}
			}
			if rqLen != len(rqCalls.List) {
				changed = true
			}
			return pairs, rqCalls
		}

		for changed {
			changed = false
			for methodName, calls := range r.funcCalls {
				m := NewCallList(r.funcNewRqList[methodName])
				for funName := range calls.List {
					m.IncCallCount(funName)
				}
				r.funcNewRqList[methodName] = m
				r.funcSetErrorMap[methodName], r.funcNewRqList[methodName] =
					mergeSetErrror(r.funcSetErrorMap[methodName], r.funcNewRqList[methodName], calls)
			}
		}

		//L.Describe(r.funcNewRqList)

		// now that NewRqList contains all recursive calls
		for methodName, calls := range r.funcNewRqList {
			m := NewCallList(r.funcNewWcList[methodName])
			// copy to NewWcList but only that contains New and ends with Mutator
			// and unset NewRqList to false
			for funName, param := range calls.List {
				if S.Contains(funName, `.New`) {
					if S.EndsWith(funName, `Mutator`) {
						m.IncCallCount(funName)
						calls.ZeroCallCount(funName) // not Rq
					}
				} else {
					calls.ZeroCallCount(funName) // not Rq
					// also show 3rd party dependencies
					for _, thirdParty := range r.ThirdParties {
						if S.StartsWith(funName, thirdParty+`.`) {
							r.func3rdParty[methodName+`|`+thirdParty] = true
						}
					}
					if S.EndsWith(funName, `.Statistics`) {
						r.IncStatisticsCalls(methodName, param)
					}
				}
			}
			r.funcNewWcList[methodName] = m
		}

		// copy to r.calls, set only true to NeRqList
		r.funcCalls = r.funcNewRqList
		r.funcNewRqList = map[string]*CallList{}
		for methodName, calls := range r.funcCalls {
			m := NewCallList(r.funcNewRqList[methodName])
			for funName, param := range calls.List {
				if param.CallCount > 0 && funName != `hmac.New` {
					m.IncCallCount(funName)
				}
			}
			r.funcNewRqList[methodName] = m
		}

		//L.Describe(r.funcNewRqList)
	}
}

func parseStructFields(res []StructField, specType ast.Expr) []StructField {
	//L.Describe(specType)
	if stru, ok := specType.(*ast.StructType); ok {
		fields := stru.Fields.List
		for _, field := range fields { // Doc, Tag, Comment
			if len(field.Names) > 0 {
				sf := StructField{}
				sf.Name = field.Names[0].Name
				//if sf.Name == `IsActive` {
				//	L.Describe(field)
				//}
				if field.Tag != nil {
					sf.Tags = field.Tag.Value
				}
				if typ, ok := field.Type.(*ast.Ident); ok {
					sf.Type = typ.String()
				} else if expr, ok := field.Type.(*ast.StarExpr); ok {
					if sel, ok := expr.X.(*ast.SelectorExpr); ok {
						sf.Type = sel.X.(*ast.Ident).Name + `.` + sel.Sel.Name
					} else {
						L.Describe(expr)
						// why we don't allow nested? because it would be more painful for FE guys to parse
						panic(`unhandled type, type nested or not in model 1`)
					}
				} else if arr, ok := field.Type.(*ast.ArrayType); ok {
					if expr, ok := arr.Elt.(*ast.StarExpr); ok {
						if sel, ok := expr.X.(*ast.SelectorExpr); ok {
							sf.Type = sel.X.(*ast.Ident).Name + `.` + sel.Sel.Name
						} else {
							L.Describe(arr)
							panic(`unhandled type, type nested or not in model 2`)
						}
					} else if sel, ok := arr.Elt.(*ast.SelectorExpr); ok {
						sf.Type = sel.X.(*ast.Ident).Name + `.` + sel.Sel.Name
					} else if sel, ok := arr.Elt.(*ast.Ident); ok {
						sf.Type = sel.Name
					} else {
						L.Describe(arr)
						panic(`unhandled type, type nested or not in model 3`)
					}
					sf.IsArray = true
				} else if _, ok := field.Type.(*ast.MapType); ok {
					sf.IsMap = true
				} else if sel, ok := field.Type.(*ast.SelectorExpr); ok {
					sf.Type = sel.X.(*ast.Ident).Name + `.` + sel.Sel.Name
					if sel.Sel.Name == `MapSlice` {
						sf.IsMap = true
					} else if sel.Sel.Name == `Time` {
						//L.Describe(sel)
						//panic(`not allowed to use time.Time, use int64 unix epoch instead`)
					} else {
						// other type = struct
					}
					// ignore any other type
				} else {
					L.Describe(field)
					panic(`unhandled type, type nested or not in model 4`)
				}
				res = append(res, sf)
			} // else: RequestCommon
		}
	}
	return res
}

func (r *RoutesArgs) Visit(n ast.Node) ast.Visitor { // parse domain
	//L.Describe("%T %#v\n", n, n)
	//aa := spew.Sdump("%T %#v", n, n)
	//if S.Contains(aa, `EventUpsert`) && S.Contains(aa, `Domain`) {
	//	L.Describe(n)
	//}
	if decl, ok := n.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
		//L.Describe(decl)
		for _, spec := range decl.Specs {
			if spec, ok := spec.(*ast.TypeSpec); ok {
				val := spec.Name.Name
				l := len(val)
				if S.EndsWith(val, `_In`) {
					methodName := val[:l-3]
					r.lastMethod = methodName
					r.methodsPkgMap[methodName] = r.lastPackage
					r.packagesRefCount[r.lastPackage] += 1
					r.inputFieldsByMethod[methodName] = parseStructFields(r.inputFieldsByMethod[methodName], spec.Type)
				}
				if S.EndsWith(val, `_Out`) {
					methodName := val[:l-4]
					r.lastMethod = methodName
					r.methodsPkgMap[methodName] = r.lastPackage
					r.packagesRefCount[r.lastPackage] += 1
					r.outputFieldsByMethod[methodName] = parseStructFields(r.outputFieldsByMethod[methodName], spec.Type)
				}
			}
		}
	}
	if r.lastMethod != `` {
		if call, ok := n.(*ast.CallExpr); ok {
			if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
				methodName := r.lastMethod
				if fun.Sel.Name == `mustLogin` {
					r.inputFieldsByMethod[methodName] = append(r.inputFieldsByMethod[methodName], StructField{
						Name:    `sessionToken`,
						Type:    `string`,
						Comment: `user login token`,
					})
				}
				if fun.Sel.Name == `mustAdmin` {
					r.inputFieldsByMethod[methodName] = append(r.inputFieldsByMethod[methodName], StructField{
						Name:    `sessionToken`,
						Type:    `string`,
						Comment: `admin login token`,
					})
				}
			}
		}
		if assign, ok := n.(*ast.AssignStmt); ok {
			if call, ok := assign.Lhs[0].(*ast.SelectorExpr); ok {
				if call.Sel.Name == `SessionToken` {
					if out, ok := call.X.(*ast.Ident); ok && out.Name == `out` {
						methodName := r.lastMethod
						r.outputFieldsByMethod[methodName] = append(r.outputFieldsByMethod[methodName], StructField{
							Name:    `sessionToken`,
							Type:    `string`,
							Comment: `login token`,
						})
					}
				}
			}
		}
	}
	if r.parseCalls {
		if fun, ok := n.(*ast.FuncDecl); ok {
			r.lastFuncDecl = fun.Name.Name
			return r
		}
		if call, ok := n.(*ast.CallExpr); ok {
			if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
				funName := fun.Sel.Name
				if funName == `Statistics` {
					m := NewCallList(r.funcStatistics[r.lastMethod])
					m.AddCallParam(call.Args)
					r.funcStatistics[r.lastMethod] = m
				}
				if funName == `SetError` {
					if len(call.Args) == 2 {
						ep := ErrPair{}
						ep.Code = S.ToInt(call.Args[0].(*ast.BasicLit).Value)
						if lit, ok := call.Args[1].(*ast.BasicLit); ok {
							ep.Msg = lit.Value
						} else {
							list, _ := call.Args[1].(*ast.BinaryExpr)
							ep.Msg = list.X.(*ast.BasicLit).Value
						}
						r.funcSetErrorMap[r.lastFuncDecl] = r.funcSetErrorMap[r.lastFuncDecl].Add(ep)
					}
				} else { // calling other function
					r.funcCalls[r.lastFuncDecl] = r.funcCalls[r.lastFuncDecl].Add(fun.Sel.Name)

					if ident, ok := fun.X.(*ast.Ident); ok {
						r.funcCalls[r.lastFuncDecl] = r.funcCalls[r.lastFuncDecl].Add(ident.Name + `.` + fun.Sel.Name)
					} else {
					}
				}
			}
		}
	}

	if ident, ok := n.(*ast.Ident); ok {
		val := ident.Name
		if S.StartsWith(val, `d`) && ident.NamePos < 10 {
			r.lastPackage = val
			r.packagesRefCount[val] += 1
			return r
		}
		l := len(val)
		if S.EndsWith(val, `_Url`) {
			methodName := val[:l-4]
			r.lastMethod = methodName
			r.methodsPkgMap[methodName] = r.lastPackage
			r.packagesRefCount[r.lastPackage] += 1
			isErr := true
			if valSpec, ok := ident.Obj.Decl.(*ast.ValueSpec); ok {
				if basLit, ok := valSpec.Values[0].(*ast.BasicLit); ok && basLit.Kind == token.STRING {
					r.methodsUrlMap[methodName], r.methodsSegmentIdxMap[methodName] = unquoteAndSplit(basLit.Value)
					isErr = false
				}
			}
			if isErr {
				L.Print(`warning: constant ` + val + ` should have a string literal value`)
			}
			return r
		}
	}
	return r
}

func unquoteAndSplit(url string) ([]string, map[string]int) {
	url = strings.TrimFunc(url, func(r rune) bool {
		return r == '"' || r == '`'
	})
	segments := strings.Split(url, `/`)
	params := []string{}
	segmentIdxMap := map[string]int{}
	for idx, param := range segments {
		if len(param) > 1 && param[0] == ':' {
			segment := param[1:]
			params = append(params, segment)
			segmentIdxMap[segment] = idx
		}
	}
	return params, segmentIdxMap
}

func q(str string) string {
	return `"` + str + `"`
}

func (r *RoutesArgs) WriteWebRoutes(path string) {
	buf := bytes.Buffer{}
	buf.WriteString(`package main

import (
	` + q(r.ProjectName+`/conf`) + `
	` + q(r.ProjectName+`/domain`) + `

	` + q(`go.opentelemetry.io/otel/trace`) + `

	` + q(`github.com/gofiber/fiber/v2`) + `
	//` + q(`github.com/kokizzu/gotro/S`))
	r.writeImports(&buf)

	buf.WriteString(`
)
`)

	buf.WriteString(`
func webApiInitRoutes(app *fiber.App) *domain.Domain {
	var (`)
	r.writeDomainInitialization(&buf)

	buf.WriteString(`
	)
`)

	// fiber
	sortedMethods := r.methodsPkgMap.SortedKeys()
	for _, method := range sortedMethods {
		pkg := r.methodsPkgMap[method]

		pm := pkg + `.` + method
		pmUrl := pm + `_Url`
		buf.WriteString(`
	app.All(conf.API_PREFIX+` + pmUrl + `, func(ctx *fiber.Ctx) error {
		url := ` + pmUrl + `
		tracerCtx, span := conf.T.Start(ctx.Context(), url, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		in := ` + pm + `_In{}
		if err := webApiParseInput(ctx, &in.RequestCommon, &in, url); err != nil {
			return err
		}`)
		writeUrlSegmentToIn(&buf, r.methodsUrlMap[method], `ctx.Params(%s)`)
		buf.WriteString(`
		in.FromFiberCtx(ctx, tracerCtx)
		out := v` + pm + `(&in)
		out.ToFiberCtx(ctx, &in.RequestCommon, &in)
		return in.ToFiberCtx(ctx, out)
	})
`)
	}
	buf.WriteString(`
	return vdomain
}
`)

	err := ioutil.WriteFile(path, buf.Bytes(), 0644)
	L.PanicIf(err, `ioutil.WriteFile failed: `+path)
}

func writeUrlSegmentToIn(buf *bytes.Buffer, segments []string, parseMethod string) {
	for _, param := range segments {
		inputMethod := fmt.Sprintf(parseMethod, q(param))
		if S.EndsWith(param, `Id`) && !S.StartsWith(param, `encoded`) {
			buf.WriteString(`
		in.` + S.CamelCase(param) + ` = S.ToU(` + inputMethod + `)`)
		} else {
			buf.WriteString(`
		in.` + S.CamelCase(param) + ` = ` + inputMethod)
		}
	}
}

func (r *RoutesArgs) WriteCliArgs(path string) {
	buf := bytes.Buffer{}
	buf.WriteString(`package main

import (
	` + q(`context`) + `
	` + q(r.ProjectName+`/conf`) + `
	` + q(r.ProjectName+`/domain`) + `

	` + q(`os`))
	r.writeImports(&buf)

	buf.WriteString(`
)
`)

	buf.WriteString(`
func cliArgsRunner(args []string) {
	tracerCtx, span := conf.T.Start(context.Background(), args[0])
	defer span.End()

	var (`)
	r.writeDomainInitialization(&buf)

	buf.WriteString(`
	)

	patterns := map[string]map[string]int{`)

	sortedMethods := r.methodsPkgMap.SortedKeys()
	maxLen := 0
	for _, method := range sortedMethods {
		if maxLen < len(method) {
			maxLen = len(method)
		}
	}
	for _, method := range sortedMethods {
		pkg := r.methodsPkgMap[method]

		pm := pkg + `.` + method
		space := S.Repeat(` `, maxLen-len(method))
		buf.WriteString(`
		` + pm + `_Url: ` + space + `{`)
		segments := r.methodsUrlMap[method]
		segmentIdxMap := r.methodsSegmentIdxMap[method]
		for z, param := range segments {
			buf.WriteString(q(param))
			buf.WriteRune(':')
			buf.WriteString(I.ToStr(segmentIdxMap[param]))
			if z != len(segments)-1 {
				buf.WriteRune(',')
			}
		}
		buf.WriteString(`},`)
	}

	buf.WriteString(`
	}
	switch pattern := cliUrlPattern(args[0], patterns); pattern {
`)

	for _, method := range sortedMethods {
		pkg := r.methodsPkgMap[method]

		pm := pkg + `.` + method
		pmUrl := pm + `_Url`
		buf.WriteString(`
	case ` + pmUrl + `:
		in := ` + pm + `_In{}
		in.FromCli(os.Stdin, tracerCtx)`)
		writeUrlSegmentToIn(&buf, r.methodsUrlMap[method], `cliSegmentFromIdx(args[0],patterns[pattern][%s])`)
		buf.WriteString(`
		out := v` + pm + `(&in)
		out.ToCli(os.Stdout)
		in.ToCli(os.Stdout, &out)
`)
	}
	buf.WriteString(`
	}
}
`)

	err := ioutil.WriteFile(path, buf.Bytes(), 0644)
	L.PanicIf(err, `ioutil.WriteFile failed: `+path)
}

func (r *RoutesArgs) writeDomainInitialization(buf *bytes.Buffer) {
	//for pkg, used := range r.packagesRefCount {
	//	if used < 4 { // if not have _In, _Out, _Url
	//		continue
	//	}
	//	buf.WriteString(`
	//	v` + pkg + ` = ` + pkg + `.New(deps)`)
	//}
	buf.WriteString(`
		vdomain = domain.NewDomain()`)

}

func (r *RoutesArgs) writeImports(buf *bytes.Buffer) {
	//imports := []string{``}
	//for pkg, used := range r.packagesRefCount {
	//	if used < 4 { // if not have _In, _Out, _Url
	//		continue
	//	}
	//	imports = append(imports, `
	//`+q(`domain`+pkg))
	//}
	//sort.Strings(imports)
	//strImports := strings.Join(imports, ``)
	//buf.WriteString(strImports)
}

func (r *RoutesArgs) WriteApiDocs(path string) {
	buf := bytes.Buffer{}
	buf.WriteString(`// can be hit using with /api/[ApiName]
export const LastUpdatedAt = ` + I.ToS(fastime.UnixNow()) + `
export const APIs = {`)

	writeField := func(field StructField) {
		if field.IsArray {
			buf.WriteString(` [], ` + field.ApiComment())
			return
		}
		if field.IsMap {
			buf.WriteString(` {key:'value'}, ` + field.ApiComment())
			return
		}
		if field.Type == `string` || S.Contains(field.Tags, `,string`) {
			buf.WriteString(` '', ` + field.ApiComment())
			return
		}
		if field.Type == `bool` {
			buf.WriteString(` false, ` + field.ApiComment())
			return
		}
		if field.Type == `float64` || field.Type == `float32` {
			buf.WriteString(` 0.0, ` + field.ApiComment())
			return
		}
		buf.WriteString(` 0, ` + field.ApiComment())
	}

	writeInnerField := func(f2 StructField) {
		buf.WriteString(`
				` + f2.LowerName() + `: `)
		if f2.Type == `string` || S.Contains(f2.Tags, `,string`) {
			buf.WriteString(` '', ` + f2.ApiComment())
		} else if f2.Type == `bool` {
			buf.WriteString(` false, ` + f2.ApiComment())
		} else if f2.Type == `float64` || f2.Type == `float32` {
			buf.WriteString(` 0.0, ` + f2.ApiComment())
		} else {
			buf.WriteString(` 0, ` + f2.ApiComment())
		}
	}

	writeInOutStruct := func(field StructField) {
		buf.WriteString(`
			` + field.LowerName() + `:`)
		nested := r.models.fieldsByModel[field.Type]
		if nested == nil {
			writeField(field)
			return
		}
		if field.IsArray {
			buf.WriteString(` [{`)
		} else {
			buf.WriteString(` {`)
		}
		//L.Describe(nested)
		for _, f2 := range nested {
			if f2.Name == `Adapter` {
				continue
			}
			writeInnerField(f2)
		}
		if field.IsArray {
			buf.WriteString(`
			}],`)
		} else {
			buf.WriteString(`
			},`)
		}
	}

	sortedMethods := r.methodsPkgMap.SortedKeys()
	for _, method := range sortedMethods {
		if method == `XXX` {
			continue
		}

		buf.WriteString(`
	` + method + `: {
		in: {`)
		inFields := r.inputFieldsByMethod[method]
		for _, field := range inFields {
			writeInOutStruct(field)
		} // in}

		buf.WriteString(`
		}, out: {`)
		outFields := r.outputFieldsByMethod[method]
		for _, field := range outFields {
			writeInOutStruct(field)
		} // out}

		buf.WriteString(`
		}, read: [`)
		rqDeps := r.funcNewRqList[method]
		if rqDeps != nil {
			dbDeps := rqDeps.SortedKeys()
			for _, name := range dbDeps {
				name = S.Replace(name, `.New`, `.`)
				name = S.Replace(name, `rq`, ``)
				buf.WriteString(`
			"` + name + `",`)
			} // read}
		}

		buf.WriteString(`
		], write: [`)
		wcDeps := r.funcNewWcList[method]
		if wcDeps != nil {
			dbDeps := wcDeps.SortedKeys()
			for _, name := range dbDeps {
				name = S.Replace(name, `.New`, `.`)
				name = S.Replace(name, `Mutator`, ``)
				name = S.Replace(name, `wc`, ``)
				buf.WriteString(`
			"` + name + `",`)
			} // write}
		}

		buf.WriteString(`
		], stat: [`)
		saDeps := r.funcStatistics[method]
		if saDeps != nil {
			chDeps := saDeps.SortedKeys()
			for _, name := range chDeps {
				if len(name) <= 2 {
					continue
				}
				name = name[2:]
				buf.WriteString(`
			"` + name + `",`)
			} // stat}
		}

		buf.WriteString(`
		], deps: [`)
		for _, thirdP := range r.ThirdParties {
			if r.func3rdParty[method+`|`+thirdP] {
				buf.WriteString(`
			"` + thirdP + `",`)
			}

		}
		buf.WriteString(`
		], err: [`)
		errList, ok := r.funcSetErrorMap[method]
		if ok && len(errList.List) > 0 {
			errPairs := []ErrPair{}
			for _, pair := range errList.List {
				errPairs = append(errPairs, pair)
			}
			sort.Slice(errPairs, func(a, b int) bool {
				A := errPairs[a]
				B := errPairs[b]
				if A.Code == B.Code {
					return A.Msg < B.Msg
				}
				return A.Code < B.Code
			})
			for _, err := range errPairs {
				buf.WriteString(`
			[` + I.ToStr(err.Code) + `, ` + err.Msg + `],`)
			}
			buf.WriteString(`
		]`) // err]
		} else {
			buf.WriteString(`]`) // err]
		}
		buf.WriteString(`
	},`)
	}

	buf.WriteString(`
}`)

	err := ioutil.WriteFile(path, buf.Bytes(), 0644)
	L.PanicIf(err, `ioutil.WriteFile failed: `+path)
}

func (r *RoutesArgs) IncStatisticsCalls(methodName string, param *CallParam) {
	m := NewCallList(r.funcStatistics[methodName])
	m.List[param.ToString()] = param
	r.funcStatistics[methodName] = m
}

func lowerFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return S.ToLower(string(s[0])) + s[1:]
}
