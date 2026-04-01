package S

import "testing"

func TestXXH3(t *testing.T) {
	v1 := XXH3(`alpha`)
	v2 := XXH3(`alpha`)
	v3 := XXH3(`beta`)
	if v1 == 0 {
		t.Fatalf("XXH3 should not be zero for non-empty input")
	}
	if v1 != v2 {
		t.Fatalf("XXH3 should be deterministic")
	}
	if v1 == v3 {
		t.Fatalf("XXH3 should differ for different input")
	}
}
