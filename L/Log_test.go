package L

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func helperForCallerInfo() *CallInfo {
	return CallerInfo()
}

func TestResourcePercentAndStackHelpers(t *testing.T) {
	prevCPUPercent := CPU_PERCENT
	prevLastCPUCall := LAST_CPU_CALL
	CPU_PERCENT = 12.34
	LAST_CPU_CALL = time.Now().Unix() + 10
	if got := PercentCPU(); got != 12.34 {
		t.Fatalf("PercentCPU cached mismatch: %v", got)
	}
	CPU_PERCENT = prevCPUPercent
	LAST_CPU_CALL = prevLastCPUCall
	if got := PercentCPU(); got < -1 || got > 1000 {
		t.Fatalf("PercentCPU range mismatch: %v", got)
	}

	prevRAMPercent := RAM_PERCENT
	prevLastRAMCall := LAST_RAM_CALL
	RAM_PERCENT = 56.78
	LAST_RAM_CALL = time.Now().Unix() + 10
	if got := PercentRAM(); got != 56.78 {
		t.Fatalf("PercentRAM cached mismatch: %v", got)
	}
	RAM_PERCENT = prevRAMPercent
	LAST_RAM_CALL = prevLastRAMCall
	if got := PercentRAM(); got < -1 || got > 1000 {
		t.Fatalf("PercentRAM range mismatch: %v", got)
	}

	trace := StackTrace(0)
	if trace == `` || !strings.Contains(trace, `testing.tRunner`) {
		t.Fatalf("StackTrace unexpected value: %q", trace)
	}
}

func TestErrorAndPanicBranches(t *testing.T) {
	if DefaultIsError(nil, "x") {
		t.Fatalf("DefaultIsError(nil) should be false")
	}
	if !DefaultIsError(errors.New("boom"), "x") {
		t.Fatalf("DefaultIsError(err) should be true")
	}
	if CheckIf(false, "x") {
		t.Fatalf("CheckIf(false) should be false")
	}
	if !CheckIf(true, "x") {
		t.Fatalf("CheckIf(true) should be true")
	}

	func() {
		defer func() {
			if recover() == nil {
				t.Fatalf("Panic should panic")
			}
		}()
		Panic("panic message")
	}()

	func() {
		defer func() {
			if recover() == nil {
				t.Fatalf("PanicIf should panic on non-nil error")
			}
		}()
		PanicIf(errors.New("boom"), "panicif message")
	}()
}

func TestCallerInfoSkipAndTrace(t *testing.T) {
	ci := helperForCallerInfo()
	if ci == nil || ci.FuncName == `` || ci.PackageName == `` || ci.FileName == `` || ci.Line <= 0 {
		t.Fatalf("CallerInfo mismatch: %#v", ci)
	}
	if skipSelf := CallerInfo(0); skipSelf == nil || skipSelf.FuncName == `` {
		t.Fatalf("CallerInfo(0) mismatch: %#v", skipSelf)
	}
	if chain := CallerChain(1, 1); len(chain) != 1 {
		t.Fatalf("CallerChain(1,1) mismatch: %#v", chain)
	}

	DEBUG = true
	Trace()
	DEBUG = false
}
