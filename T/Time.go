package T

// Time support package
import (
	"github.com/kokizzu/gotro/F"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"math/rand"
	"time"
)

const ISO = `2006-01-02T15:04:05.999999`
const YMD_HM = `2006-01-02 15:04`
const YMD_HMS = `2006-01-02 15:04:05`
const YMD = `2006-01-02`
const FILE = `20060102_150405`
const HUMAN = `2-Jan-2006 15:04:05`

var EMPTY = time.Time{}

// convert time to iso formatted time string
//  T.ToIsoStr(time.Now()) // "2016-03-17T10:04:50.6489"
func ToIsoStr(t time.Time) string {
	if t == EMPTY {
		return ``
	}
	return t.Format(ISO)
}

// current iso time
//  T.IsoStr() // "2016-03-17T10:07:56.418728"
func IsoStr() string {
	return time.Now().Format(ISO)
}

// convert time to iso date
//  T.ToDateStr(time.Now()) // output "2016-03-17"
func ToDateStr(t time.Time) string {
	if t == EMPTY {
		return ``
	}
	return t.Format(YMD)
}

// current iso date
// T.DateStr()) // "2016-03-17"
func DateStr() string {
	return time.Now().Format(YMD)
}

// convert time to human date
//  T.ToHumanStr(time.Now()) // "17-Mar-2016 10:06"
func ToHumanStr(t time.Time) string {
	if t == EMPTY {
		return ``
	}
	return t.Format(HUMAN)
}

// current human date
//  T.HumanStr() // "17-Mar-2016 10:06"
func HumanStr() string {
	return time.Now().Format(HUMAN)
}

// convert time to iso date and hour:minute
//  T.ToDateHourStr(time.Now()) // "2016-03-17 10:07"
func ToDateHourStr(t time.Time) string {
	if t == EMPTY {
		return ``
	}
	return t.Format(YMD_HM)
}

// current iso date and hour:minute
//  T.DateHourStr()// output "2016-03-17 10:07"
func DateHourStr(t time.Time) string {
	return time.Now().Format(YMD_HM)
}

// convert time to iso date and time
//  T.ToDateTimeStr(time.Now()) // "2016-03-17 10:07:50"
func ToDateTimeStr(t time.Time) string {
	if t == EMPTY {
		return ``
	}
	return t.Format(YMD_HMS)
}

// current iso date and time
//  T.ToDateTimeStr(time.Now()) // "2016-03-17 10:07:50"
func DateTimeStr(t time.Time) string {
	return time.Now().Format(YMD_HMS)
}

// int64 day of current date
func DayInt() int64 {
	return int64(time.Now().Day())
}

// int64 current hour
func HourInt() int64 {
	return int64(time.Now().Hour())
}

// int64 current month
func MonthInt() int64 {
	return int64(time.Now().Month())
}

// int64 current year
func YearInt() int64 {
	return int64(time.Now().Year())
}

// int64 current day of year
func YearDayInt() int64 {
	return int64(time.Now().YearDay())
}

// get filename version of current dat
//  T.Filename()) // "20160317_102543"
func Filename() string {
	return time.Now().Format(FILE)
}

// sleep for nanosec
func Sleep(ns time.Duration) {
	time.Sleep(ns)
}

// random 0.4-2 sec sleep
func RandomSleep() {
	dur := rand.Int63()%(1600*1000*1000) + (400 * 1000 * 1000)
	time.Sleep(time.Duration(dur))
}

// measure elapsed time in nanosec
//  T.Track(func(){
//    x:=0
//    T.Sleep(1)
//  }) // "done in 1.00s"
func Track(fun func()) time.Duration {
	start := time.Now()
	fun()
	elapsed := time.Since(start)
	L.Describe(`done in ` + F.ToStr(elapsed.Seconds()) + `s`)
	return elapsed
}

// check if time in are in the range
//  t1, _:=time.Parse(`1992-03-23`,T.DateFormat)
//  t2, _:=time.Parse(`2016-03-17`,T.DateFormat)
//  T.IsValidTimeRange(t1,t2,time.Now()) // bool(false)
func IsValidTimeRange(start, end, check time.Time) bool {
	res := check.After(start) && check.Before(end)
	return res
}

// returns age from current date
func Age(birthdate time.Time) float64 {
	return float64(time.Now().Sub(birthdate)/time.Hour) / 24 / 365.25
}

// returns age from within 2 date
func AgeAt(birthdate, point time.Time) float64 {
	return float64(point.Sub(birthdate)/time.Hour) / 24 / 365.25
}

// get current unix nano
func UnixNano() int64 {
	return time.Now().UnixNano()
}

// get current unix nano after added with certain duration
func UnixNanoAfter(d time.Duration) int64 {
	return time.Now().Add(d).UnixNano()
}

// get current unix (second) as integer
func Epoch() int64 {
	return time.Now().Unix()
}

// get current unix (second) as string
func EpochStr() string {
	return I.ToS(time.Now().Unix())
}

// get current unix time added with a duration
func EpochAfter(d time.Duration) int64 {
	return time.Now().Add(d).Unix()
}
