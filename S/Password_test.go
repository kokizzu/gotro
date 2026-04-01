package S

import "testing"

func TestHashPassword(t *testing.T) {
	if got := HashPassword(`abc`); got != `ungWv48Bz+pBQUDeXa4iI7ADYaOWF3qctBD/YfIAFa0=` {
		t.Fatalf("HashPassword mismatch: %q", got)
	}
	if got := HashPassword(``); got != `47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=` {
		t.Fatalf("HashPassword(empty) mismatch: %q", got)
	}
}

func TestEncryptAndCheckPassword(t *testing.T) {
	raw := `myS3cret!`
	hash := EncryptPassword(raw)
	if hash == `` || hash == raw {
		t.Fatalf("EncryptPassword should return bcrypt hash")
	}
	if err := CheckPassword(hash, raw); err != nil {
		t.Fatalf("CheckPassword valid input should pass: %v", err)
	}
	if err := CheckPassword(hash, raw+`x`); err == nil {
		t.Fatalf("CheckPassword invalid input should fail")
	}
}
