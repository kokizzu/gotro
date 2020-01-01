package S

import (
	"fmt"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/B"
	"github.com/kokizzu/gotro/I"
	"runtime"
	"strconv"
)

// TODO: find out how backspace \b null \0 character processed on common SQL

// trace function, location of the caller code, replacement for ZC
func ZT(strs ...string) string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	str := A.StrJoin(strs, `|`)
	str = Replace(str, "\n", ` `)
	str = Replace(str, `\n`, ` `)
	return `-- ` + fmt.Sprintf("%s:%d %s|%s", file, line, f.Name(), str)
}

// trace function, location of 2nd level caller, parameterless, with newline
func ZT2() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	return `-- ` + fmt.Sprintf("%s:%d %s", file, line, f.Name()) + "\n"
}

// add single quote in the beginning and the end of string.
//  S.Q(`coba`) // `'coba'`
//  S.Q(`123`)  // `'123'`
func Q(str string) string {
	return `'` + str + `'`
}

// replace ` and give double quote (for table names)
//  S.ZZ(`coba"`) // `"coba&quot;"`
func ZZ(str string) string {
	str = Trim(str)
	str = Replace(str, `"`, `&quot;`)
	return `"` + str + `"`
}

// give ' to boolean value
//  S.ZB(true)  // `'true'`
//  S.ZB(false) // `'false'`
func ZB(b bool) string {
	return `'` + B.ToS(b) + `'`
}

// give ' to int64 value
//  S.ZI(23)) // '23'
//  S.ZI(03)) // '3'
func ZI(num int64) string {
	return `'` + I.ToS(num) + `'`
}

// give ' to uint value
//  S.ZI(23)) // '23'
//  S.ZI(03)) // '3'
func ZU(num uint64) string {
	return `'` + strconv.FormatUint(num, 10) + `'`
}

// double quote a json string
//  hai := `{'test':123,"bla":[1,2,3,4]}`
//  S.ZJJ(hai) // "{'test':123,\"bla\":[1,2,3,4]}"
func ZJJ(str string) string {
	str = Replace(str, `\`, `\\`)
	str = Replace(str, "\r", `\r`)
	str = Replace(str, "\n", `\n`)
	str = Replace(str, `"`, `\"`)
	return `"` + str + `"`
}

// single quote a json string
//  hai := `{'test':123,"bla":[1,2,3,4]}`
//  S.ZJ(hai) // "{'test':123,\"bla\":[1,2,3,4]}"
func ZJ(str string) string {
	str = Replace(str, `\`, `\\`)
	str = Replace(str, "\r", `\r`)
	str = Replace(str, "\n", `\n`)
	str = Replace(str, `'`, `\'`)
	return `'` + str + `'`
}

// trim, replace <, >, ', " and gives single quote
//  S.Z(`<>'"`) // `&lt;&gt;&apos;&quot;
func Z(str string) string {
	str = Trim(str)
	str = Replace(str, `<`, `&lt;`)
	str = Replace(str, `>`, `&gt;`)
	str = Replace(str, `'`, `&apos;`)
	str = Replace(str, `"`, `&quot;`)
	str = Replace(str, `\`, `\\`)
	return `'` + str + `'`
}

// replace <, >, ', " and gives single quote (without trimming)
//  S.Z(`<>'"`) // `&lt;&gt;&apos;&quot;
func ZS(str string) string {
	str = Replace(str, `<`, `&lt;`)
	str = Replace(str, `>`, `&gt;`)
	str = Replace(str, `'`, `&apos;`)
	str = Replace(str, `"`, `&quot;`)
	str = Replace(str, `\`, `\\`)
	return `'` + str + `'`
}

// replace <, >, ', ", % and gives single quote and %
//  S.ZLIKE(`coba<`))  // output '%coba&lt;%'
//  S.ZLIKE(`"coba"`)) // output '%&quot;coba&quot;%'
func ZLIKE(str string) string {
	str = Trim(str)
	str = Replace(str, `<`, `&lt;`)
	str = Replace(str, `>`, `&gt;`)
	str = Replace(str, `'`, `&apos;`)
	str = Replace(str, `"`, `&quot;`)
	str = Replace(str, `\`, `\\`)
	str = Replace(str, `%`, `\%`)
	return `'%` + str + `%'`
}

// ZLIKE but for json (not replacing double quote)
func ZJLIKE(str string) string {
	str = Trim(str)
	str = Replace(str, `<`, `&lt;`)
	str = Replace(str, `>`, `&gt;`)
	str = Replace(str, `'`, `&apos;`)
	str = Replace(str, `\`, `\\`)
	str = Replace(str, `%`, `\%`)
	return `'%` + str + `%'`
}

// replace <, >, ', ", % but without giving single quote
func XSS(str string) string {
	str = Trim(str)
	str = Replace(str, `<`, `&lt;`)
	str = Replace(str, `>`, `&gt;`)
	str = Replace(str, `'`, `&apos;`)
	str = Replace(str, `"`, `&quot;`)
	return str
}

// replace <, >, and & back, quot and apos to alternative utf8
func UZ(str string) string {
	str = Replace(str, `&apos;`, `‘`)
	str = Replace(str, `&quot;`, `ʺ`)
	str = Replace(str, `&lt;`, `<`)
	str = Replace(str, `&gt;`, `>`)
	str = Replace(str, `&amp;`, `&`)
	return str
}

// replace <, >, and & back, quot and apos to real html
func UZRAW(str string) string {
	str = Replace(str, `&apos;`, `'`)
	str = Replace(str, `&quot;`, `"`)
	str = Replace(str, `&lt;`, `<`)
	str = Replace(str, `&gt;`, `>`)
	str = Replace(str, `&amp;`, `&`)
	return str
}
