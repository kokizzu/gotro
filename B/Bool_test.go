package B

import "testing"

func TestToS(t *testing.T) {
	if got := ToS(true); got != `true` {
		t.Fatalf("ToS(true) mismatch: %q", got)
	}
	if got := ToS(false); got != `false` {
		t.Fatalf("ToS(false) mismatch: %q", got)
	}
}
