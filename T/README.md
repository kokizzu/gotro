# T
--
    import "github.com/kokizzu/gotro/T"


## Usage

```go
const FILE = `20060102_150405`
```go

```go
const HMS = `150405`
```go

```go
const HUMAN = `2-Jan-2006 15:04:05`
```go

```go
const HUMAN_DATE = `2 Jan 2006`
```go

```go
const ISO = `2006-01-02T15:04:05.999999`
```go

```go
const YMD = `2006-01-02`
```go

```go
const YMDH = `20060102.15`
```go

```go
const YMDHM = `20060102.1504`
```go

```go
const YMD_HM = `2006-01-02 15:04`
```go

```go
const YMD_HMS = `2006-01-02 15:04:05`
```go

```go
const YY = `06`
```go

```go
var EMPTY = time.Time{}
```go

#### func  Age

```go
func Age(birthdate time.Time) float64
```go
returns age from current date

#### func  AgeAt

```go
func AgeAt(birthdate, point time.Time) float64
```go
returns age from within 2 date

#### func  DateHhMmStr

```go
func DateHhMmStr() string
```go
current iso date and hour

    T.DateHhMmStr()// output "20160317.1059"

#### func  DateHhStr

```go
func DateHhStr() string
```go
current iso date and hour

    T.DateHhStr()// output "20160317.10"

#### func  DateStr

```go
func DateStr() string
```go
current iso date T.DateStr()) // "2016-03-17"

#### func  DateTimeStr

```go
func DateTimeStr() string
```go
current iso date and time

    T.ToDateTimeStr(time.Now()) // "2016-03-17 10:07:50"

#### func  DayInt

```go
func DayInt() int64
```go
int64 day of current date

#### func  Epoch

```go
func Epoch() int64
```go
get current unix (second) as integer

#### func  EpochAfter

```go
func EpochAfter(d time.Duration) int64
```go
get current unix time added with a duration

#### func  EpochAfterStr

```go
func EpochAfterStr(d time.Duration) string
```go
get current unix time added with a duration

#### func  EpochStr

```go
func EpochStr() string
```go
get current unix (second) as string

#### func  Filename

```go
func Filename() string
```go
get filename version of current date

    T.Filename()) // "20160317_102543"

#### func  HhmmssStr

```go
func HhmmssStr() string
```go
get filename version of current time

#### func  HourInt

```go
func HourInt() int64
```go
int64 current hour

#### func  HumanStr

```go
func HumanStr() string
```go
current human date

    T.HumanStr() // "17-Mar-2016 10:06"

#### func  IsValidTimeRange

```go
func IsValidTimeRange(start, end, check time.Time) bool
```go
check if time in are in the range

    t1, _:=time.Parse(`1992-03-23`,T.DateFormat)
    t2, _:=time.Parse(`2016-03-17`,T.DateFormat)
    T.IsValidTimeRange(t1,t2,time.Now()) // bool(false)

#### func  IsoStr

```go
func IsoStr() string
```go
current iso time

    T.IsoStr() // "2016-03-17T10:07:56.418728"

#### func  LastTwoDigitYear

```go
func LastTwoDigitYear() string
```go
return current last two digit year

#### func  MonthInt

```go
func MonthInt() int64
```go
int64 current month

#### func  RandomSleep

```go
func RandomSleep()
```go
random 0.4-2 sec sleep

#### func  Sleep

```go
func Sleep(ns time.Duration)
```go
sleep for nanosec

#### func  ToDateHourStr

```go
func ToDateHourStr(t time.Time) string
```go
convert time to iso date and hour:minute

    T.ToDateHourStr(time.Now()) // "2016-03-17 10:07"

#### func  ToDateStr

```go
func ToDateStr(t time.Time) string
```go
convert time to iso date

    T.ToDateStr(time.Now()) // output "2016-03-17"

#### func  ToDateTimeStr

```go
func ToDateTimeStr(t time.Time) string
```go
convert time to iso date and time

    T.ToDateTimeStr(time.Now()) // "2016-03-17 10:07:50"

#### func  ToEpoch

```go
func ToEpoch(date string) int64
```go
2019-07-16 Yonas convert string date to epoch => '2019-01-01' -->1546300800

#### func  ToHhmmssStr

```go
func ToHhmmssStr(t time.Time) string
```go
convert time to iso date and hourminutesecond

    T.ToDateHourStr(time.Now()) // "230744"

#### func  ToHumanStr

```go
func ToHumanStr(t time.Time) string
```go
convert time to human date

    T.ToHumanStr(time.Now()) // "17-Mar-2016 10:06"

#### func  ToIsoStr

```go
func ToIsoStr(t time.Time) string
```go
convert time to iso formatted time string

    T.ToIsoStr(time.Now()) // "2016-03-17T10:04:50.6489"

#### func  Track

```go
func Track(fun func()) time.Duration
```go
measure elapsed time in nanosec

    T.Track(func(){
      x:=0
      T.Sleep(1)
    }) // "done in 1.00s"

#### func  UnixNano

```go
func UnixNano() int64
```go
get current unix nano

#### func  UnixNanoAfter

```go
func UnixNanoAfter(d time.Duration) int64
```go
get current unix nano after added with certain duration

#### func  UnixToDateStr

```go
func UnixToDateStr(epoch float64) string
```go
convert from unix to date format

#### func  UnixToDateTimeStr

```go
func UnixToDateTimeStr(epoch float64) string
```go
convert

#### func  UnixToFile

```go
func UnixToFile(i int64) string
```go
convert unix time to file naming

#### func  UnixToHumanDateStr

```go
func UnixToHumanDateStr(epoch float64) string
```go
convert from unix to human date format

#### func  UnixToHumanStr

```go
func UnixToHumanStr(epoch float64) string
```go
convert from unix to human format

#### func  Weekday

```go
func Weekday() int
```go
get what day is it today, Sunday => 0

#### func  WeekdayStr

```go
func WeekdayStr() string
```go
get day's name

#### func  YearDayInt

```go
func YearDayInt() int64
```go
int64 current day of year

#### func  YearInt

```go
func YearInt() int64
```go
int64 current year
