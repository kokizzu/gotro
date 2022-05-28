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
CheckIf print error

#### func  CreateDir

```go
func CreateDir(path string) bool
```
CreateDir create directory recursively

#### func  CreateFile

```go
func CreateFile(path string, content string) bool
```
CreateFile create file with specific content

#### func  Describe

```go
func Describe(args ...interface{})
```
Describe pretty print any variable

#### func  FileEmpty

```go
func FileEmpty(name string) bool
```
FileEmpty check file missing or has zero size

#### func  FileExists

```go
func FileExists(name string) bool
```
FileExists check file exists

#### func  IsError

```go
func IsError(err error, msg string, args ...interface{}) bool
```
IsError print error

#### func  LogTrack

```go
func LogTrack(start time.Time, name string) float64
```
LogTrack return elapsed time in ms, show 3nd level, returns in ms

#### func  PanicIf

```go
func PanicIf(err error, msg string, args ...interface{})
```
PanicIf print error message and exit program

#### func  ParentDescribe

```go
func ParentDescribe(args ...interface{})
```
ParentDescribe describe anything

#### func  PercentCPU

```go
func PercentCPU() float64
```
PercentCPU get CPU usage percentage

    L.PercentCPU()

#### func  PercentRAM

```go
func PercentRAM() float64
```
PercentRAM get RAM usage percentage

    L.PercentRAM()

#### func  PipeRunCmd

```go
func PipeRunCmd(cmd string, args ...string) error
```
PipeRunCmd run cmd and pipe to stdout

#### func  Print

```go
func Print(any ...interface{})
```
Print replacement for fmt.Println, gives line number

#### func  PrintParent

```go
func PrintParent(any ...interface{})
```
PrintParent print but show grandparent caller function

#### func  ReadFile

```go
func ReadFile(path string) string
```
ReadFile read file content as string

#### func  ReadFileLines

```go
func ReadFileLines(path string, eachLineFunc func(line string) (exitEarly bool)) (ok bool)
```
ReadFileLines read file content line by line

#### func  RunCmd

```go
func RunCmd(cmd string, args ...string) (output []byte)
```
RunCmd execute command and return output

#### func  StackTrace

```go
func StackTrace(start int) string
```
StackTrace get a stacktrace as string

    L.StackTrace(0) // until current function
    L.StackTrace(1) // until function that call this function

#### func  TimeTrack

```go
func TimeTrack(start time.Time, name string) float64
```
TimeTrack return elapsed time in ms, show 1st level, returns in ms

#### func  Trace

```go
func Trace()
```
Trace trace a function call

#### type CallInfo

```go
type CallInfo struct {
	PackageName string
	FileName    string
	FuncName    string
	Line        int
}
```


#### func  CallerChain

```go
func CallerChain(skipFrom, skipUntil int) (res []CallInfo)
```
CallerChain return caller chain until specific skipFrom 1 to 2 will return from
parent caller until grandparent

#### func  CallerInfo

```go
func CallerInfo(skip ...int) (caller *CallInfo)
```
CallerInfo return caller info default skip is 1, equal to parent caller

#### func (*CallInfo) String

```go
func (c *CallInfo) String() string
```
