# T
--
    import "github.com/kokizzu/gotro/T"


## Usage

```go
const FILE = `20060102_150405`
```

```go
const HMS = `150405`
```

```go
const HUMAN = `2-Jan-2006 15:04:05`
```

```go
const HUMAN_DATE = `2 Jan 2006`
```

```go
const ISO = `2006-01-02T15:04:05.999999`
```

```go
const YMD = `2006-01-02`
```

```go
const YMDH = `20060102.15`
```

```go
const YMDHM = `20060102.1504`
```

```go
const YMD_HM = `2006-01-02 15:04`
```

```go
const YMD_HMS = `2006-01-02 15:04:05`
```

```go
const YY = `06`
```

```go
var EMPTY = time.Time{}
```

#### func  Age

```go
func Age(birthdate time.Time) float64
```
Age returns age from current date

#### func  AgeAt

```go
func AgeAt(birthdate, point time.Time) float64
```
AgeAt returns age from within 2 date

#### func  DateHhMmStr

```go
func DateHhMmStr() string
```
DateHhMmStr current iso date and hour

    T.DateHhMmStr()// output "20160317.1059"

#### func  DateHhStr

```go
func DateHhStr() string
```
DateHhStr current iso date and hour

    T.DateHhStr()// output "20160317.10"

#### func  DateStr

```go
func DateStr() string
```
DateStr current iso date T.DateStr()) // "2016-03-17"

#### func  DateTimeStr

```go
func DateTimeStr() string
```
DateTimeStr current iso date and time

    T.ToDateTimeStr(time.Now()) // "2016-03-17 10:07:50"

#### func  DayInt

```go
func DayInt() int64
```
DayInt int64 day of current date

#### func  Epoch

```go
func Epoch() int64
```
Epoch get current unix (second) as integer

#### func  EpochAfter

```go
func EpochAfter(d time.Duration) int64
```
EpochAfter get current unix time added with a duration

#### func  EpochAfterStr

```go
func EpochAfterStr(d time.Duration) string
```
EpochAfterStr get current unix time added with a duration

#### func  EpochStr

```go
func EpochStr() string
```
EpochStr get current unix (second) as string

#### func  Filename

```go
func Filename() string
```
Filename get filename version of current date

    T.Filename()) // "20160317_102543"

#### func  HhmmssStr

```go
func HhmmssStr() string
```
HhmmssStr get filename version of current time

#### func  HourInt

```go
func HourInt() int64
```
HourInt int64 current hour

#### func  HumanStr

```go
func HumanStr() string
```
HumanStr current human date

    T.HumanStr() // "17-Mar-2016 10:06"

#### func  IsValidTimeRange

```go
func IsValidTimeRange(start, end, check time.Time) bool
```
IsValidTimeRange check if time in are in the range

    t1, _:=time.Parse(`1992-03-23`,T.DateFormat)
    t2, _:=time.Parse(`2016-03-17`,T.DateFormat)
    T.IsValidTimeRange(t1,t2,time.Now()) // bool(false)

#### func  IsoStr

```go
func IsoStr() string
```
IsoStr current iso time

    T.IsoStr() // "2016-03-17T10:07:56.418728"

#### func  LastTwoDigitYear

```go
func LastTwoDigitYear() string
```
LastTwoDigitYear return current last two digit year

#### func  MonthInt

```go
func MonthInt() int64
```
MonthInt int64 current month

#### func  RandomSleep

```go
func RandomSleep()
```
RandomSleep random 0.4-2 sec sleep

#### func  Sleep

```go
func Sleep(ns time.Duration)
```
Sleep delay for nanosecond

#### func  ToDateHourStr

```go
func ToDateHourStr(t time.Time) string
```
ToDateHourStr convert time to iso date and hour:minute

    T.ToDateHourStr(time.Now()) // "2016-03-17 10:07"

#### func  ToDateStr

```go
func ToDateStr(t time.Time) string
```
ToDateStr convert time to iso date

    T.ToDateStr(time.Now()) // output "2016-03-17"

#### func  ToDateTimeStr

```go
func ToDateTimeStr(t time.Time) string
```
ToDateTimeStr convert time to iso date and time

    T.ToDateTimeStr(time.Now()) // "2016-03-17 10:07:50"

#### func  ToEpoch

```go
func ToEpoch(date string) int64
```
ToEpoch convert string date to epoch => '2019-01-01' -->1546300800

#### func  ToHhmmssStr

```go
func ToHhmmssStr(t time.Time) string
```
ToHhmmssStr convert time to iso date and hourminutesecond

    T.ToDateHourStr(time.Now()) // "230744"

#### func  ToHumanStr

```go
func ToHumanStr(t time.Time) string
```
ToHumanStr convert time to human date

    T.ToHumanStr(time.Now()) // "17-Mar-2016 10:06"

#### func  ToIsoStr

```go
func ToIsoStr(t time.Time) string
```
ToIsoStr convert time to iso formatted time string

    T.ToIsoStr(time.Now()) // "2016-03-17T10:04:50.6489"

#### func  Track

```go
func Track(fun func()) time.Duration
```
Track measure elapsed time in nanosec

    T.Track(func(){
      x:=0
      T.Sleep(1)
    }) // "done in 1.00s"

#### func  UnixNano

```go
func UnixNano() int64
```
UnixNano get current unix nano

#### func  UnixNanoAfter

```go
func UnixNanoAfter(d time.Duration) int64
```
UnixNanoAfter get current unix nano after added with certain duration

#### func  UnixToDateStr

```go
func UnixToDateStr(epoch float64) string
```
UnixToDateStr convert from unix sconds to YYYY-MM-DD

#### func  UnixToDateTimeStr

```go
func UnixToDateTimeStr(epoch float64) string
```
UnixToDateTimeStr convert unix seconds to YYYY-MM-DD_hh:mm:ss

#### func  UnixToFile

```go
func UnixToFile(i int64) string
```
UnixToFile convert unix time to file naming

#### func  UnixToHumanDateStr

```go
func UnixToHumanDateStr(epoch float64) string
```
UnixToHumanDateStr convert from unix to human date format D MMM YYYY

#### func  UnixToHumanStr

```go
func UnixToHumanStr(epoch float64) string
```
UnixToHumanStr convert from unix to human format D-MMM-YYYY hh:mm:ss

#### func  Weekday

```go
func Weekday() int
```
Weekday get what day is it today, Sunday => 0

#### func  WeekdayStr

```go
func WeekdayStr() string
```
WeekdayStr get day's name

#### func  YearDayInt

```go
func YearDayInt() int64
```
YearDayInt int64 current day of year

#### func  YearInt

```go
func YearInt() int64
```
YearInt int64 current year
