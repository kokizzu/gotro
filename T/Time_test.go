package T

import (
	"strconv"
	"testing"
	"time"
)

func TestFixedTimeFormatters(t *testing.T) {
	fixed := time.Date(2024, time.February, 29, 23, 59, 58, 123456000, time.UTC)
	if got := ToIsoStr(fixed); got != `2024-02-29T23:59:58.123456` {
		t.Fatalf("ToIsoStr mismatch: %q", got)
	}
	if got := ToDateStr(fixed); got != `2024-02-29` {
		t.Fatalf("ToDateStr mismatch: %q", got)
	}
	if got := ToHumanStr(fixed); got != `29-Feb-2024 23:59:58` {
		t.Fatalf("ToHumanStr mismatch: %q", got)
	}
	if got := ToDateHourStr(fixed); got != `2024-02-29 23:59` {
		t.Fatalf("ToDateHourStr mismatch: %q", got)
	}
	if got := ToHhmmssStr(fixed); got != `235958` {
		t.Fatalf("ToHhmmssStr mismatch: %q", got)
	}
	if got := ToDateTimeStr(fixed); got != `2024-02-29 23:59:58` {
		t.Fatalf("ToDateTimeStr mismatch: %q", got)
	}

	if ToIsoStr(EMPTY) != `` || ToDateStr(EMPTY) != `` || ToHumanStr(EMPTY) != `` || ToDateHourStr(EMPTY) != `` || ToHhmmssStr(EMPTY) != `` || ToDateTimeStr(EMPTY) != `` {
		t.Fatalf("empty time should render empty string")
	}
}

func TestEpochAndUnixConverters(t *testing.T) {
	if got := ToEpoch(`1970-01-02`); got != 86400 {
		t.Fatalf("ToEpoch mismatch: %d", got)
	}
	if got := ToEpoch(`not-a-date`); got != 0 {
		t.Fatalf("ToEpoch invalid input should be 0, got %d", got)
	}

	if got := UnixToFile(0); got != time.Unix(0, 0).Format(FILE) {
		t.Fatalf("UnixToFile mismatch: %q", got)
	}
	if got := UnixToDateTimeStr(86400); got != time.Unix(86400, 0).Format(YMD_HMS) {
		t.Fatalf("UnixToDateTimeStr mismatch: %q", got)
	}
	if got := UnixToDateStr(86400); got != time.Unix(86400, 0).Format(YMD) {
		t.Fatalf("UnixToDateStr mismatch: %q", got)
	}
	if got := UnixToHumanDateStr(86400); got != time.Unix(86400, 0).Format(HUMAN_DATE) {
		t.Fatalf("UnixToHumanDateStr mismatch: %q", got)
	}
	if got := UnixToHumanStr(86400); got != time.Unix(86400, 0).Format(HUMAN) {
		t.Fatalf("UnixToHumanStr mismatch: %q", got)
	}
}

func TestNowBasedHelpers(t *testing.T) {
	if l := len(DateHhStr()); l != len(`20060102.15`) {
		t.Fatalf("DateHhStr length mismatch: %d", l)
	}
	if l := len(DateHhMmStr()); l != len(`20060102.1504`) {
		t.Fatalf("DateHhMmStr length mismatch: %d", l)
	}
	if l := len(Filename()); l != len(`20060102_150405`) {
		t.Fatalf("Filename length mismatch: %d", l)
	}
	if l := len(HhmmssStr()); l != len(`150405`) {
		t.Fatalf("HhmmssStr length mismatch: %d", l)
	}
	if l := len(LastTwoDigitYear()); l != 2 {
		t.Fatalf("LastTwoDigitYear length mismatch: %d", l)
	}

	if d := DayInt(); d < 1 || d > 31 {
		t.Fatalf("DayInt out of range: %d", d)
	}
	if h := HourInt(); h < 0 || h > 23 {
		t.Fatalf("HourInt out of range: %d", h)
	}
	if m := MonthInt(); m < 1 || m > 12 {
		t.Fatalf("MonthInt out of range: %d", m)
	}
	if y := YearInt(); y < 1970 {
		t.Fatalf("YearInt out of range: %d", y)
	}
	if yd := YearDayInt(); yd < 1 || yd > 366 {
		t.Fatalf("YearDayInt out of range: %d", yd)
	}
	if wd := Weekday(); wd < 0 || wd > 6 {
		t.Fatalf("Weekday out of range: %d", wd)
	}
	if WeekdayStr() == `` {
		t.Fatalf("WeekdayStr should not be empty")
	}
	if DateStr() == `` || DateTimeStr() == `` || IsoStr() == `` || HumanStr() == `` {
		t.Fatalf("date/time string helpers should not return empty")
	}

	now := Epoch()
	after := EpochAfter(2 * time.Second)
	if diff := after - now; diff < 1 || diff > 3 {
		t.Fatalf("EpochAfter diff out of range: %d", diff)
	}
	afterStr, err := strconv.ParseInt(EpochAfterStr(2*time.Second), 10, 64)
	if err != nil {
		t.Fatalf("EpochAfterStr not numeric: %v", err)
	}
	if afterStr <= now {
		t.Fatalf("EpochAfterStr should be in future: now=%d after=%d", now, afterStr)
	}
	if _, err := strconv.ParseInt(EpochStr(), 10, 64); err != nil {
		t.Fatalf("EpochStr not numeric: %v", err)
	}

	nano := UnixNano()
	nanoAfter := UnixNanoAfter(time.Millisecond)
	if nanoAfter <= nano {
		t.Fatalf("UnixNanoAfter should be greater than UnixNano")
	}
}

func TestRangeAndAgeHelpers(t *testing.T) {
	start := time.Unix(10, 0)
	end := time.Unix(20, 0)
	if IsValidTimeRange(start, end, start) {
		t.Fatalf("range should be exclusive on start")
	}
	if !IsValidTimeRange(start, end, time.Unix(15, 0)) {
		t.Fatalf("range should include middle value")
	}
	if IsValidTimeRange(start, end, end) {
		t.Fatalf("range should be exclusive on end")
	}

	now := time.Now()
	if AgeAt(now, now) != 0 {
		t.Fatalf("AgeAt same time should be 0")
	}
	if Age(now.Add(-365*24*time.Hour)) <= 0 {
		t.Fatalf("Age should be positive for past date")
	}

	Sleep(time.Nanosecond)
	if elapsed := Track(func() {}); elapsed < 0 {
		t.Fatalf("Track should be non-negative")
	}
}
