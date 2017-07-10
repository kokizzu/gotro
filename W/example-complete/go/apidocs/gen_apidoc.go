package main

import (
	"bufio"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

type Action struct {
	name     string
	function string
	comment  string
}

type Record struct {
	posts     M.SB
	gets      M.SB
	responses M.SB
	actions   []Action
	functions M.SB
	models    M.SB
	errors    M.SB
	files     M.SB
}

func NewRecord() *Record {
	return &Record{
		posts:     M.SB{},
		gets:      M.SB{},
		responses: M.SB{},
		functions: M.SB{},
		models:    M.SB{},
		errors:    M.SB{},
		files:     M.SB{},
	}
}

var DIRS = [...]string{`../../handler`, `../../sql`, `../../../github.com/kokizzu/gotro`}
var ROUTERS = map[string]string{
	`MAIN`: `../../router.go`,
}
var FUNCTIONS M.SB
var OUTPUT = `../../public/apidocs.html`
var DOCS []string

func init() {
	FUNCTIONS = M.SB{}
	DOCS = []string{}
}

var RE = M.SX{
	`post`:     regexp.MustCompile(`(osts(?P<post>\.Get.+\)))|(dm\.(?P<post>Set.+\)))|(tx\.(?P<post>UploadedFile.+\)))|(osts\(\)(?P<post>\.Get.+\)))`),
	`get`:      regexp.MustCompile(`aram\.(?P<get>Get.+\))`),
	`response`: regexp.MustCompile(`jax\.(?P<response>Set.+\))`),
	`action`:   regexp.MustCompile(`case(?P<action>.+)// @API.*`),
	`model`:    regexp.MustCompile(`(?P<model>m[a-zA-Z]+\.TM_[^\)\.\,]+)`),
	`error`:    regexp.MustCompile(`sql\.(?P<error>ERR[^\)]+)`),
}

var params_re = regexp.MustCompile(`:[a-zA-Z_]*`)

var records = map[string]Record{}

func record_all_functions(path string, fi os.FileInfo, err error) error {
	if filepath.Ext(path) == `.go` {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(f)
		var package_name = ``
		for scanner.Scan() {
			line := S.Trim(scanner.Text())
			if S.StartsWith(scanner.Text(), `package`) {
				package_name = S.Split(line, ` `)[1]
				switch package_name {
				case `X`, `S`, `I`, `F`, `M`, `A`: // ignore these packages
					return nil
				}
			} else if S.StartsWith(line, `func`) {
				s := S.SplitFunc(line[len(`func `):], func(c rune) bool {
					return c == '(' || c == ')'
				})
				func_name := s[0]
				FUNCTIONS[package_name+`.`+func_name] = true
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}
	return nil
}

func mergeMapStr(dst M.SB, str string) {
	str = S.Trim(str)
	if str != `` {
		dst[str] = true
	}
}

func visit_all_functions(path string, fi os.FileInfo, err error) error {
	if filepath.Ext(path) == `.go` {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(f)
		var package_name = ""
		var table = ""
		for scanner.Scan() {
			line := S.Trim(scanner.Text())
			if S.StartsWith(line, `package`) {
				package_name = S.Split(line, ` `)[1]
			} else if S.StartsWith(line, `const TABLE = `) {
				table = package_name
			} else if S.StartsWith(line, `func`) {
				s := S.SplitFunc(line[len(`func `):], func(c rune) bool {
					return c == '(' || c == ')'
				})
				func_name := s[0]
				record := NewRecord()
				has_table := false
				for scanner.Scan() && !S.StartsWith(scanner.Text(), `}`) {
					line = S.Trim(scanner.Text())
					if S.Contains(line, `TABLE`) {
						has_table = true
					}
					if S.Contains(line, `(`) {
						s2 := S.SplitFunc(line, func(c rune) bool {
							return c == '(' || c == '=' || c == ' ' || c == '\t'
						})
						for _, fname := range s2 {
							if func_name == fname {
								continue
							}
							if FUNCTIONS[package_name+`.`+fname] {
								mergeMapStr(record.functions, package_name+`.`+fname)
							} else if S.Contains(fname, `.`) && FUNCTIONS[fname] {
								mergeMapStr(record.functions, fname)
							}
						}
					}
					for k, v := range RE {
						re := v.(*regexp.Regexp)
						if match := re.FindString(line); match != `` {
							split := S.Split(line, `//`)
							comment := ``
							if len(split) > 1 {
								comment = spanCls(COMMENT, ` // `+split[1])
							}
							if k == `post` || k == `response` {
								match = A.StrJoin(S.Split(match, ".")[1:], `.`)
							}
							if k == `post` {
								mergeMapStr(record.posts, match+comment)
							} else if k == `get` {
								mergeMapStr(record.gets, match+comment)
							} else if k == `response` {
								mergeMapStr(record.responses, match+comment)
							} else if k == `model` {
								mergeMapStr(record.models, match+comment)
							} else if k == `error` {
								mergeMapStr(record.errors, S.SplitFunc(match, func(c rune) bool {
									return c == ' ' || c == ')'
								})[0]+comment)
							} else if k == `action` {
								//L.Print(`action`,scanner.Text())
								//fmt.Println(`next line`,scanner.Text())
								comment = ``
								ss := S.Split(match, "`") // case `foo`, `bar`, `baz`: // @API hel
								name := ``
								if len(ss) > 2 {
									name = ss[1] // foo
									ss = ss[2:]  // : // [bar, baz, : // @API hel]
									for len(ss) > 2 {
										name += `, ` + ss[1]
										ss = ss[2:]
									}
								} else {
									ss = S.Split(match, ` `) // case NOT_REQUIRED: // @API
									if len(ss) > 1 {
										name = ss[1] // NOT_REQUIRED:
									}
								}
								ss = S.Split(match, `@API`)
								if len(ss) > 1 {
									comment = ss[1]
								}
								// read next line, that should be a function call mSomething.FuncName()
								scanner.Scan()
								line := S.Trim(scanner.Text())
								text := S.Split(S.Trim(line), "(")[0]
								ss = S.Split(text, ` `)
								record.actions = append(record.actions, Action{
									name,
									ss[len(ss)-1],
									comment,
								})
								//fmt.Println(`name`,name,ss[len(ss) - 1])
							}
							mergeMapStr(record.files, path)
						}
					}
				}
				if has_table && table != `` {
					mergeMapStr(record.models, table)
				}
				records[package_name+`.`+func_name] = *record
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}
	return nil
}

const ACTION = `action`
const COMMENT = `comment`

func spanCls(cls, str string) string {
	if str == `` {
		return ``
	}
	return ` <span class="` + cls + `">` + str + `</span>`
}

func aHref(href, str string) string {
	return `<a href="#` + href + `">` + str + `</a>`
}

type HandleParams struct {
	Handler string
	Params  []string
}

func explore_router() {
	fo, err := os.Create(OUTPUT)
	if err != nil {
		panic(err)
	}
	defer fo.Close()
	w := bufio.NewWriter(fo)
	w.WriteString(`<!DOCTYPE html><html><head><title>example-complete-example-cron's API Docs</title><style>
body {font-family: sans-serif}
b {color: blue; padding-left: 4px}
p {margin: 0}
ol {counter-reset: item; display: table}
ol > li {list-style-type: none; display: table-row}
ol > li:before {content: counters(item, ".") " "; counter-increment: item; color:red; font-weight: bold; display: table-cell}
span.action, a {padding-left: 3px; padding-right: 3px; color: green}
span.comment {color: brown}
</style></head>`)

	w.WriteString(`<body>`)
	w.WriteString(`<h1>example-complete-example-cron API Documentation</h1>`)
	w.WriteString(`<h2 id="toc">Table of Contents</h2>`)

	domain_route_handler_params := map[string]map[string]HandleParams{}

	w.WriteString(`<ol>`)
	sorted_domains := []string{}
	for domain := range ROUTERS {
		sorted_domains = append(sorted_domains, domain)
	}
	sort.Strings(sorted_domains)
	for _, domain := range sorted_domains {
		router := ROUTERS[domain]
		domain_route_handler_params[domain] = map[string]HandleParams{}
		fi, err := os.Open(router)
		if err != nil {
			panic(err)
		}
		defer fi.Close()
		scanner := bufio.NewScanner(fi)
		for scanner.Scan() && !S.Contains(scanner.Text(), `var HANDLERS = map[string]W.Action{`) {
			if S.Contains(scanner.Text(), `// @DOC`) {
				// @DOC may only on router
				DOCS = append(DOCS, S.Split(scanner.Text(), `@DOC`)[1])
			}
		}
		for scanner.Scan() && !S.StartsWith(scanner.Text(), `}`) {
			if ss := S.Split(scanner.Text(), "`"); len(ss) > 1 {
				route := ss[1]
				handler := S.Trim(ss[2][1 : len(ss[2])-1])
				params := []string{}
				for _, p := range params_re.FindAllString(route, -1) {
					params = append(params, p[1:])
				}
				domain_route_handler_params[domain][route] = HandleParams{handler, params}
				record := records[handler]
				full_route := domain + `/` + route
				w.WriteString(`<li id="toc|` + full_route + `">`)
				w.WriteString(full_route)
				w.WriteString(aHref(full_route, `&#x25BC;`))
				for _, action := range record.actions {
					w.WriteString(aHref(full_route+`@`+action.name, action.name))
				}
				w.WriteString(`</li>`)
			}
		}
	}
	w.WriteString(`</ol>`)
	w.WriteString(`<h3>Notes</h3>`)
	writeArr(w, ``, DOCS)
	w.WriteString(`<h2>Details</h2>`)
	w.WriteString(`<ol>`)
	for _, domain := range sorted_domains {
		routes := domain_route_handler_params[domain]
		sorted_routes := []string{}
		for route := range routes {
			sorted_routes = append(sorted_routes, route)
		}
		sort.Strings(sorted_routes)
		for _, route := range sorted_routes {
			hp := routes[route]
			record := records[hp.Handler]
			w.WriteString(`<li>`)
			full_route := domain + `/` + route
			writeNoP(w, `Route`, `<u id="`+full_route+`">`+full_route+`</u>`+aHref(`toc|`+full_route, `&#x25B2;`))
			writeWithP(w, `Handler`, hp.Handler)
			writeArr(w, `Params`, hp.Params)
			writeMap(w, `Model`, record.models)
			writeMap(w, `Possible Error`, record.errors)

			if len(record.actions) > 0 {
				w.WriteString(`<p><b>Actions</b>: `)
				w.WriteString(`<ol>`)
				for _, action := range record.actions {
					name, function, comment := action.name, action.function, action.comment
					w.WriteString(`<li id="` + full_route + `@` + name + `">`)
					writeNoP(w, `Action`, spanCls(ACTION, name)+aHref(`toc|`+full_route, `&#x25B2;`))
					writeWithP(w, `Description`, spanCls(COMMENT, comment))
					tranverse_func(function)
					if record, ok := records[function]; ok {
						writeMap(w, `Post`, record.posts)
						writeMap(w, `Get`, record.gets)
						writeMap(w, `Response`, record.responses)
						writeMap(w, `Model`, record.models)
						writeMap(w, `Possible Error`, record.errors)
						//writeMapCensor(w, `Files`, record.files)
						writeMap(w, `Functions`, record.functions)
					}
					w.WriteString(`</li>`)
					w.WriteByte('\n')
				}
				w.WriteString(`</ol>`)
				w.WriteString(`</p>`)
			}

			w.WriteString(`</li>`)
			w.WriteByte('\n')
		}
	}

	w.WriteString(`</ol></body></html>`)
	w.Flush()
}

func writeWithP(w *bufio.Writer, name, value string) {
	if value != `` {
		w.WriteString(`<p><b>`)
		w.WriteString(name)
		w.WriteString(`</b>: `)
		w.WriteString(value)
		w.WriteString(`</p>`)
		w.WriteByte('\n')
	}
}
func writeNoP(w *bufio.Writer, name, value string) {
	if value != `` {
		w.WriteString(`<b>`)
		w.WriteString(name)
		w.WriteString(`</b>: `)
		w.WriteString(value)
		w.WriteByte('\n')
	}
}

func writeMap(w *bufio.Writer, name string, values M.SB) {
	if len(values) > 0 {
		w.WriteString(`<p><b>`)
		w.WriteString(name)
		w.WriteString(`</b>: <ul><li>`)
		w.WriteString(A.StrJoin(values.SortedKeys(), `</li><li>`))
		w.WriteString(`</li></ul></p>`)
		w.WriteByte('\n')
	}
}

func writeMapCensor(w *bufio.Writer, name string, values M.SB) {
	if len(values) > 0 {
		w.WriteString(`<p><b>`)
		w.WriteString(name)
		w.WriteString(`</b>: <ul><li>`)
		keys := values.SortedKeys()
		for k, v := range keys {
			keys[k] = v[:len(v)-3]
		}
		w.WriteString(A.StrJoin(keys, `</li><li>`))
		w.WriteString(`</li></ul></p>`)
		w.WriteByte('\n')
	}
}

func writeArr(w *bufio.Writer, name string, values []string) {
	if len(values) > 0 {
		w.WriteString(`<p>`)
		if name != `` {
			w.WriteString(`<b>`)
			w.WriteString(name)
			w.WriteString(`</b>:`)
		}
		w.WriteString(`<ul><li>`)
		w.WriteString(A.StrJoin(values, `</li><li>`))
		w.WriteString(`</li></ul></p>`)
		w.WriteByte('\n')
	}
}

func mergeMapMap(dst, src M.SB) {
	for k := range src {
		if k != `` {
			dst[k] = true
		}
	}
}

func tranverse_func(function_name string) Record {
	if record, ok := records[function_name]; ok {
		for f := range record.functions {
			ret := tranverse_func(f)
			mergeMapMap(record.posts, ret.posts)
			mergeMapMap(record.gets, ret.gets)
			mergeMapMap(record.responses, ret.responses)
			record.actions = append(record.actions, ret.actions...)
			mergeMapMap(record.models, ret.models)
			mergeMapMap(record.errors, ret.errors)
			mergeMapMap(record.files, ret.files)
		}
		record.functions = M.SB{}
		records[function_name] = record
		return record
	}
	return Record{}
}

func main() {
	for _, dir := range DIRS {
		err := filepath.Walk(dir, record_all_functions)
		L.PanicIf(err, `record_all_functions`)
	}
	for _, dir := range DIRS {
		err := filepath.Walk(dir, visit_all_functions)
		L.PanicIf(err, `visit`)
	}
	explore_router()
}
