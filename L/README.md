# L
--
    import "github.com/kokizzu/gotro/L"


## Usage

```go
const WebBR = "\n<br/>"
```

```go
var BgRed, BgGreen (func(format string, a ...interface{}) string)
```

```go
var CPU_PERCENT, RAM_PERCENT float64
```

```go
var CPU_STAT2, CPU_STAT4, CPU_STAT5, LAST_STAT7 int64
```

```go
var DEBUG bool
```

```go
var FILE_PATH string
```

```go
var GO_PATH string
```

```go
var GO_ROOT string
```

```go
var LAST_CPU_CALL, LAST_RAM_CALL int64
```

```go
var LOG *logging.Logger
```

```go
var NUM_CPU float64
```

```go
var TIMETRACK_MIN_DURATION float64
```

#### func  CheckIf

```go
func CheckIf(is_err bool, msg string, args ...interface{}) bool
```
print error

#### func  CreateDir

```go
func CreateDir(path string) bool
```

#### func  CreateFile

```go
func CreateFile(path string, content string) bool
```

#### func  Describe

```go
func Describe(args ...interface{})
```
describe anything

#### func  FileEmpty

```go
func FileEmpty(name string) bool
```

#### func  FileExists

```go
func FileExists(name string) bool
```

#### func  IsError

```go
func IsError(err error, msg string, args ...interface{}) bool
```
print error

#### func  LogTrack

```go
func LogTrack(start time.Time, name string) float64
```
return elapsed time in ms, show 3nd level, returns in ms

#### func  PanicIf

```go
func PanicIf(err error, msg string, args ...interface{})
```
print error message and exit program

#### func  ParentDescribe

```go
func ParentDescribe(args ...interface{})
```
describe anything

#### func  PercentCPU

```go
func PercentCPU() float64
```
get CPU usage percentage

    L.PercentCPU()

#### func  PercentRAM

```go
func PercentRAM() float64
```
get RAM usage percentage

    L.PercentRAM()

#### func  PipeRunCmd

```go
func PipeRunCmd(cmd string, args ...string) error
```
run cmd and pipe to stdout

#### func  Print

```go
func Print(any ...interface{})
```
replacement for fmt.Println, gives line number

#### func  PrintParent

```go
func PrintParent(any ...interface{})
```
print but show grandparent caller function

#### func  ReadFile

```go
func ReadFile(path string) string
```

#### func  RunCmd

```go
func RunCmd(cmd string, args ...string) (output []byte)
```
execute command and return output

#### func  StackTrace

```go
func StackTrace(start int) string
```
get a stacktrace as string

    L.StackTrace(0) // until current function
    L.StackTrace(1) // until function that call this function

#### func  TimeTrack

```go
func TimeTrack(start time.Time, name string) float64
```
return elapsed time in ms, show 1st level, returns in ms

#### func  Trace

```go
func Trace()
```
trace a function call

#### type CallInfo

```go
type CallInfo struct {
	PackageName string
	FileName    string
	FuncName    string
	Line        int
}
```


#### func  CallerInfo

```go
func CallerInfo(skip ...int) *CallInfo
```
