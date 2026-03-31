package F

import "testing"

func TestIfHelpers(t *testing.T) {
	if If(true, 1.5) != 1.5 || If(false, 1.5) != 0 {
		t.Fatalf("If mismatch")
	}
	if IfElse(true, 1.5, 2.5) != 1.5 || IfElse(false, 1.5, 2.5) != 2.5 {
		t.Fatalf("IfElse mismatch")
	}
}

func TestStringAndDateConverters(t *testing.T) {
	if got := ToS(123.456); got != `123.456` {
		t.Fatalf("ToS mismatch: %q", got)
	}
	if got := ToStr(123.456); got != `123.46` {
		t.Fatalf("ToStr mismatch: %q", got)
	}
	if got := ToIsoDateStr(0); got != `1970-01-01T00:00:00Z` {
		t.Fatalf("ToIsoDateStr mismatch: %q", got)
	}
	if got := ToDateStr(86400); got != `1970-01-02` {
		t.Fatalf("ToDateStr mismatch: %q", got)
	}
}
