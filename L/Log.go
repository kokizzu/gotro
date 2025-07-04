package L

// Logging support package

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/kr/pretty"
	"github.com/op/go-logging"

	"github.com/kokizzu/gotro/I"
)

var LOG *logging.Logger

var FILE_PATH string
var GO_PATH string
var GO_ROOT string

var NUM_CPU float64
var CPU_STAT2, CPU_STAT4, CPU_STAT5, LAST_STAT7 int64
var CPU_PERCENT, RAM_PERCENT float64
var LAST_CPU_CALL, LAST_RAM_CALL int64
var TIMETRACK_MIN_DURATION float64

var BgRed, BgGreen (func(format string, a ...any) string)

const WebBR = "\n<br/>"

// initialize logger
func init() {
	_, file, _, _ := runtime.Caller(0)
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

	TIMETRACK_MIN_DURATION = 100 // 100ms is slow query

	NUM_CPU = float64(runtime.NumCPU())
	PercentCPU()

	BgGreen = color.New(color.BgGreen).SprintfFunc()
	BgRed = color.New(color.BgRed).SprintfFunc()
}

// PercentCPU get CPU usage percentage
//
//	L.PercentCPU()
func PercentCPU() float64 {
	last_cpu_call := time.Now().Unix()
	if last_cpu_call <= LAST_CPU_CALL {
		return CPU_PERCENT
	}
	LAST_CPU_CALL = last_cpu_call
	fs, err := os.Open("/proc/stat") // usage
	if err != nil {
		return -1
	}
	defer fs.Close()
	scanner := bufio.NewScanner(fs)
	if !scanner.Scan() {
		return -1
	}
	cpu_stat := strings.Fields(scanner.Text())
	if len(cpu_stat) < 7 {
		return -1
	}
	l7, _ := strconv.ParseInt(cpu_stat[7], 10, 64)
	if l7 > LAST_STAT7 {
		s13, _ := strconv.ParseInt(cpu_stat[2-1], 10, 64)
		s15, _ := strconv.ParseInt(cpu_stat[4-1], 10, 64)
		s16, _ := strconv.ParseInt(cpu_stat[5-1], 10, 64)
		percent := float64(s13-CPU_STAT2+s15-CPU_STAT4) * 100 / float64(s13-CPU_STAT2+s15-CPU_STAT4+s16-CPU_STAT5) / float64(l7-LAST_STAT7)
		CPU_PERCENT = percent
		LAST_STAT7 = l7
		CPU_STAT2 = s13
		CPU_STAT4 = s15
		CPU_STAT5 = s16
	}
	return CPU_PERCENT
}

// PercentRAM get RAM usage percentage
//
//	L.PercentRAM()
func PercentRAM() float64 {
	last_ram_call := time.Now().Unix()
	if last_ram_call <= LAST_RAM_CALL {
		return RAM_PERCENT
	}
	LAST_RAM_CALL = last_ram_call
	fs, err := os.Open(`/proc/meminfo`) // usage
	if err != nil {
		return -1
	}
	defer fs.Close()
	scanner := bufio.NewScanner(fs)
	var ram [3]float64
	for x := 0; x < 3; x++ {
		if !scanner.Scan() {
			return -1
		}
		cols := strings.Fields(scanner.Text())
		ram[x], _ = strconv.ParseFloat(cols[1], 64)
	}
	RAM_PERCENT = (1 - ram[2]/ram[0]) * 100
	return RAM_PERCENT
}

// StackTrace get a stacktrace as string
//
//	L.StackTrace(0) // until current function
//	L.StackTrace(1) // until function that call this function
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
		if len(file) > len(FILE_PATH) {
			file = file[len(FILE_PATH):]
		}
		str += "\n\t" + file + `:` + I.ToStr(line) + `  ` + name
	}
	return str
}

// IsError print error
var IsError = DefaultIsError

// DefaultIsError function that prints error with stacktrace
func DefaultIsError(err error, msg string, args ...any) bool {
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

// CheckIf print error
func CheckIf(isErr bool, msg string, args ...any) bool {
	if !isErr {
		return false
	}
	pc, file, line, _ := runtime.Caller(1)
	str := color.MagentaString(file[len(FILE_PATH):] + `:` + I.ToStr(line) + `: `)
	str += color.YellowString(` ` + runtime.FuncForPC(pc).Name() + `: `)
	LOG.Errorf(str+msg, args...)
	res := pretty.Formatter(isErr)
	LOG.Errorf("%# v\n", res)
	str = StackTrace(3)
	res = pretty.Formatter(isErr)
	LOG.Criticalf("%# v\n    StackTrace: %s", res, str)
	return true
}

// Describe pretty print any variable
func Describe(args ...any) {
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
	LOG.Debug(strings.ReplaceAll(str, `%`, `%%`))
}

// ParentDescribe describe anything
func ParentDescribe(args ...any) {
	pc, file, line, _ := runtime.Caller(2)
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
	LOG.Debug(strings.ReplaceAll(str, `%`, `%%`))
}

// Print replacement for fmt.Println, gives line number
func Print(any ...any) {
	_, file, line, _ := runtime.Caller(1)
	str := color.CyanString(file[len(FILE_PATH):] + `:` + I.ToStr(line) + `: `)
	LOG.Debug(strings.ReplaceAll(str, `%`, `%%`))
	fmt.Println(any...)
}

// PrintParent print but show grandparent caller function
func PrintParent(any ...any) {
	_, file, line, _ := runtime.Caller(2)
	str := color.CyanString(file[len(FILE_PATH):] + `:` + I.ToStr(line) + `: `)
	LOG.Debug(strings.ReplaceAll(str, `%`, `%%`))
	fmt.Println(any...)
}

// PanicIf print error message and exit program
func PanicIf(err error, msg string, args ...any) {
	if err == nil {
		return
	}
	if err.Error() == `sql: no rows in result set` {
		return
	}
	pc, file, line, _ := runtime.Caller(1)
	strf := file[len(FILE_PATH):] + `:` + I.ToStr(line) + `: `
	str := color.MagentaString(strf)
	strf2 := ` ` + runtime.FuncForPC(pc).Name() + `: `
	str += color.YellowString(strf2)
	LOG.Criticalf(str+msg, args...)
	stt := StackTrace(3)
	res := pretty.Formatter(err)
	LOG.Criticalf("%# v\n    StackTrace: %s", res, stt)
	panic(fmt.Errorf(err.Error()+WebBR+fmt.Sprintf("%# v"+WebBR+"    StackTrace: %s", res, stt)+WebBR+strf+strf2+WebBR+msg, args...))
}

// PanicIf print error message and exit program
func Panic(msg string, args ...any) {
	pc, file, line, _ := runtime.Caller(1)
	strf := file[len(FILE_PATH):] + `:` + I.ToStr(line) + `: `
	str := color.MagentaString(strf)
	strf2 := ` ` + runtime.FuncForPC(pc).Name() + `: `
	str += color.YellowString(strf2)
	LOG.Criticalf(str+msg, args...)
	stt := StackTrace(3)
	LOG.Criticalf("StackTrace: %s", stt)
	panic(fmt.Errorf(WebBR+fmt.Sprintf(WebBR+"    StackTrace: %s", stt)+WebBR+strf+strf2+WebBR+msg, args...))
}

// TimeTrack return elapsed time in ms, show 1st level, returns in ms
func TimeTrack(start time.Time, name string) float64 {
	_, file, line, _ := runtime.Caller(1)
	prefix := color.YellowString(file[len(FILE_PATH):] + `:` + I.ToStr(line) + `: `)
	elapsed := float64(time.Since(start).Nanoseconds()) / 1000000.0
	if elapsed < TIMETRACK_MIN_DURATION {
		return elapsed
	}
	suffix := color.GreenString(`%.2f ms`, elapsed)
	LOG.Noticef(prefix+"%s "+suffix, name)
	return elapsed
}

// LogTrack return elapsed time in ms, show 3nd level, returns in ms
func LogTrack(start time.Time, name string) float64 {
	_, file, line, _ := runtime.Caller(3)
	prefix := color.CyanString(file[len(FILE_PATH):] + `:` + I.ToStr(line) + `: `)
	elapsed := float64(time.Since(start).Nanoseconds()) / 1000000.0
	if elapsed < TIMETRACK_MIN_DURATION {
		return elapsed
	}
	suffix := color.CyanString(`%.2f ms`, elapsed)
	LOG.Noticef(prefix+"%s "+suffix, name)
	return elapsed
}

var DEBUG bool

// Trace trace a function call
func Trace() {
	if !DEBUG {
		return
	}
	pc, file, line, _ := runtime.Caller(1)
	str := ` [TRACE] ` + file[len(FILE_PATH):] + `:` + I.ToStr(line) + `: ` + runtime.FuncForPC(pc).Name() + ` `
	fmt.Println(str)
}

// RunCmd execute command and return output
func RunCmd(cmd string, args ...string) (output []byte) {
	var err error
	fullcmd := `RunCmd: ` + cmd + ` `
	if len(args) > 0 {
		fullcmd += `'` + strings.Join(args, `' '`) + `' `
	}
	Describe(fullcmd)
	output, err = exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		out_str := string(output)
		Describe(fullcmd, err)
		output = []byte(err.Error())
		fmt.Println(out_str)
	}
	return
}

// PipeRunCmd run cmd and pipe to stdout
func PipeRunCmd(cmd string, args ...string) error {
	fullcmd := `RunCmd: ` + cmd + ` `
	if len(args) > 0 {
		fullcmd += `'` + strings.Join(args, `' '`) + `' `
	}
	Describe(fullcmd)
	exe := exec.Command(cmd, args...)
	exe.Stdout = os.Stdout
	exe.Stderr = os.Stderr
	return exe.Run()
}

type CallInfo struct {
	PackageName string
	FileName    string
	FuncName    string
	Line        int
}

func (c *CallInfo) String() string {
	return fmt.Sprintf("%s:%d %s.%s", c.FileName, c.Line, c.PackageName, c.FuncName)
}

// CallerInfo return caller info
// default skip is 1, equal to parent caller
func CallerInfo(skip ...int) (caller *CallInfo) {
	caller = &CallInfo{}
	skipCount := 1
	if len(skip) > 0 {
		skipCount = skip[0]
	}

	pc, file, line, ok := runtime.Caller(skipCount)
	if !ok {
		return
	}

	caller.Line = line
	_, caller.FileName = path.Split(file)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), `.`)
	pl := len(parts)
	caller.FuncName = parts[pl-1]

	if parts[pl-2][0] == '(' {
		caller.FuncName = parts[pl-2] + `.` + caller.FuncName
		caller.PackageName = strings.Join(parts[0:pl-2], `.`)
	} else {
		caller.PackageName = strings.Join(parts[0:pl-1], `.`)
	}

	return
}

// CallerChain return caller chain until specific
// skipFrom 1 to 2 will return from parent caller until grandparent
func CallerChain(skipFrom, skipUntil int) (res []CallInfo) {
	for skipCount := skipFrom; skipCount <= skipUntil; skipCount++ {
		pc, file, line, ok := runtime.Caller(skipCount)
		if !ok {
			return
		}

		caller := CallInfo{Line: line}
		_, caller.FileName = path.Split(file)
		parts := strings.Split(runtime.FuncForPC(pc).Name(), `.`)
		pl := len(parts)
		caller.FuncName = parts[pl-1]

		if parts[pl-2][0] == '(' {
			caller.FuncName = parts[pl-2] + `.` + caller.FuncName
			caller.PackageName = strings.Join(parts[0:pl-2], `.`)
		} else {
			caller.PackageName = strings.Join(parts[0:pl-1], `.`)
		}

		res = append(res, caller)
	}
	return
}
