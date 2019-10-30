package L

// Logging support package

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/kokizzu/gotro/I"
	"github.com/kr/pretty"
	"github.com/op/go-logging"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
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

var BgRed, BgGreen (func(format string, a ...interface{}) string)

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

// get CPU usage percentage
//  L.PercentCPU()
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

// get RAM usage percentage
//  L.PercentRAM()
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
		if len(file) > len(FILE_PATH) {
			file = file[len(FILE_PATH):]
		}
		str += "\n\t" + file + `:` + I.ToStr(line) + `  ` + name
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

// describe anything
func ParentDescribe(args ...interface{}) {
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
	LOG.Debug(strings.Replace(str, `%`, `%%`, -1))
}

// replacement for fmt.Println, gives line number
func Print(any ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	str := color.CyanString(file[len(FILE_PATH):] + `:` + I.ToStr(line) + `: `)
	LOG.Debug(strings.Replace(str, `%`, `%%`, -1))
	fmt.Println(any...)
}

// print but show grandparent caller function
func PrintParent(any ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	str := color.CyanString(file[len(FILE_PATH):] + `:` + I.ToStr(line) + `: `)
	LOG.Debug(strings.Replace(str, `%`, `%%`, -1))
	fmt.Println(any...)
}

// print error message and exit program
func PanicIf(err error, msg string, args ...interface{}) {
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

// return elapsed time in ms, show 1st level, returns in ms
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

// return elapsed time in ms, show 3nd level, returns in ms
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

// trace a function call
func Trace() {
	if !DEBUG {
		return
	}
	pc, file, line, _ := runtime.Caller(1)
	str := ` [TRACE] ` + file[len(FILE_PATH):] + `:` + I.ToStr(line) + `: ` + runtime.FuncForPC(pc).Name() + ` `
	fmt.Println(str)
}

// execute command and return output
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

// run cmd and pipe to stdout
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
