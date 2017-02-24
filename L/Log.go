package L

// Logging support package

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/kokizzu/gotro/I"
	"github.com/kr/pretty"
	"github.com/op/go-logging"
	"os"
	"runtime"
	"strings"
)

var LOG *logging.Logger

var FILE_PATH string
var GO_PATH string
var GO_ROOT string

// initialize logger
func init() {
	_, file, _, _ := runtime.Caller(1)
	FILE_PATH = file[:4+strings.Index(file, `/src/`)]
	GO_PATH = os.Getenv(`GOPATH`)
	GO_ROOT = os.Getenv(`GOROOT`)
	LOG = logging.MustGetLogger(`[GotRo]`)
	backend := logging.NewLogBackend(os.Stderr, ``, 0)
	format := logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05.000} %{shortfunc} ▶%{color:reset} %{message}`,
	)
	formatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(formatter)
}

// get a stacktrace as string
//  L.StackTrace(0) // until current function
//  L.StackTrace(1) // until function that call this function
func StackTrace(start int) string {
	str := ``
	for {
		pc, file, line, ok := runtime.Caller(start)
		name := runtime.FuncForPC(pc).Name()
		if !ok || strings.HasPrefix(name, `main.`) {
			break
		}
		start += 1
		if strings.HasPrefix(name, `runtime.`) || strings.HasPrefix(name, `github.com/kokizzu/gotro/L.`) {
			continue
		}
		str += "\n\t" + file[len(FILE_PATH):] + `:` + I.ToStr(line) + `  ` + name
	}
	return str
}

// print error
func IsError(err error, msg string, args ...interface{}) bool {
	if err == nil {
		return false
	}
	pc, file, line, _ := runtime.Caller(1)
	str := color.MagentaString(file[len(FILE_PATH):] + `:` + I.ToStr(line) + `: `)
	str += color.YellowString(` ` + runtime.FuncForPC(pc).Name() + `: `)
	LOG.Errorf(str+msg, args...)
	res := pretty.Formatter(err)
	LOG.Errorf("%# v\n", res)
	str = StackTrace(3)
	res = pretty.Formatter(err)
	LOG.Criticalf("%# v\n    StackTrace: %s", res, str)
	return true
}

// print error
func CheckIf(is_err bool, msg string, args ...interface{}) bool {
	if !is_err {
		return false
	}
	pc, file, line, _ := runtime.Caller(1)
	str := color.MagentaString(file[len(FILE_PATH):] + `:` + I.ToStr(line) + `: `)
	str += color.YellowString(` ` + runtime.FuncForPC(pc).Name() + `: `)
	LOG.Errorf(str+msg, args...)
	res := pretty.Formatter(is_err)
	LOG.Errorf("%# v\n", res)
	str = StackTrace(3)
	res = pretty.Formatter(is_err)
	LOG.Criticalf("%# v\n    StackTrace: %s", res, str)
	return true
}

// describe anything
func Describe(args ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	prefix := ``
	if len(file) >= len(FILE_PATH) {
		prefix = file[len(FILE_PATH):]
	}
	str := color.CyanString(prefix + `:` + I.ToStr(line) + `: `)
	str += color.YellowString(` ` + runtime.FuncForPC(pc).Name() + "\n")
	for _, arg := range args {
		//res, _ := json.MarshalIndent(variable, `   `, `  `)
		res := pretty.Formatter(arg)
		str += fmt.Sprintf("\t%# v\n", res)
	}
	LOG.Debug(strings.Replace(str, `%`, `%%`, -1))
}

// replacement for fmt.Println, gives line number
func Print(any ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	str := color.CyanString(file[len(FILE_PATH):] + `:` + I.ToStr(line) + `: `)
	LOG.Info(strings.Replace(str, `%`, `%%`, -1))
	fmt.Println(any...)
}
