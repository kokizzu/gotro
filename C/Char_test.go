package C

import "testing"

func TestCharClassifiers(t *testing.T) {
	if !IsDigit('9') || IsDigit('x') {
		t.Fatalf("IsDigit mismatch")
	}
	if !IsAlpha(byte('A')) || !IsAlpha('z') || IsAlpha('7') {
		t.Fatalf("IsAlpha mismatch")
	}
	if !IsIdentStart('_') || IsIdentStart('-') {
		t.Fatalf("IsIdentStart mismatch")
	}
	if !IsIdent('9') || !IsIdent('_') || IsIdent('!') {
		t.Fatalf("IsIdent mismatch")
	}
	if !IsValidFilename(' ') || !IsValidFilename('-') || IsValidFilename('/') {
		t.Fatalf("IsValidFilename mismatch")
	}
}
