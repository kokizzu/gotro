package L

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFileHelpers(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, `nested`, `file.txt`)

	if FileExists(path) {
		t.Fatalf("file should not exist yet")
	}
	if !FileEmpty(path) {
		t.Fatalf("missing file should be considered empty")
	}
	if !CreateDir(filepath.Dir(path)) {
		t.Fatalf("CreateDir failed")
	}
	if !CreateFile(path, "line1\nline2\n") {
		t.Fatalf("CreateFile failed")
	}
	if !FileExists(path) {
		t.Fatalf("file should exist")
	}
	if FileEmpty(path) {
		t.Fatalf("file should not be empty")
	}
	if got := ReadFile(path); got != "line1\nline2\n" {
		t.Fatalf("ReadFile mismatch: %q", got)
	}

	var lines []string
	ok := ReadFileLines(path, func(line string) (exitEarly bool) {
		lines = append(lines, line)
		return false
	})
	if !ok || len(lines) != 2 {
		t.Fatalf("ReadFileLines mismatch: ok=%v lines=%#v", ok, lines)
	}

	calls := 0
	ok = ReadFileLines(path, func(line string) (exitEarly bool) {
		calls++
		return true
	})
	if !ok || calls != 1 {
		t.Fatalf("ReadFileLines early-exit mismatch: ok=%v calls=%d", ok, calls)
	}
}

func TestCallerAndCmdHelpers(t *testing.T) {
	ci := CallerInfo()
	if ci == nil || ci.FileName == `` || ci.Line <= 0 || ci.FuncName == `` || ci.PackageName == `` {
		t.Fatalf("CallerInfo mismatch: %#v", ci)
	}
	if !strings.Contains(ci.String(), `:`) {
		t.Fatalf("CallInfo.String should include line info")
	}

	chain := CallerChain(1, 2)
	if len(chain) == 0 {
		t.Fatalf("CallerChain should return at least one entry")
	}

	out := string(RunCmd("echo", "hello"))
	if !strings.Contains(out, "hello") {
		t.Fatalf("RunCmd output mismatch: %q", out)
	}
	if err := PipeRunCmd("echo", "world"); err != nil {
		t.Fatalf("PipeRunCmd failed: %v", err)
	}
}

func TestTrackAndGuardHelpers(t *testing.T) {
	prevMin := TIMETRACK_MIN_DURATION
	TIMETRACK_MIN_DURATION = 1_000_000_000
	defer func() {
		TIMETRACK_MIN_DURATION = prevMin
	}()

	if elapsed := TimeTrack(time.Now().Add(-5*time.Millisecond), "test"); elapsed <= 0 {
		t.Fatalf("TimeTrack should return positive elapsed, got %f", elapsed)
	}
	if elapsed := LogTrack(time.Now().Add(-5*time.Millisecond), "test"); elapsed <= 0 {
		t.Fatalf("LogTrack should return positive elapsed, got %f", elapsed)
	}

	DEBUG = false
	Trace()
	PanicIf(nil, "noop")
	PanicIf(errors.New("sql: no rows in result set"), "noop")
}
