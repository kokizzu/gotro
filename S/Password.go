package S

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/kokizzu/gotro/L"
	"golang.org/x/crypto/bcrypt"
)

// hash password with sha256 (without salt)
func HashPassword(pass string) string {
	res1 := []byte(pass)
	res2 := sha256.Sum256(res1)
	res3 := res2[:]
	return base64.StdEncoding.EncodeToString(res3)
}

// hash password (with salt)
func EncryptPassword(s string) string {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	L.PanicIf(err, `failed to encrypt password`)

	return string(hashedBytes[:])
}

// check encrypted password
//
func CheckPassword(hash string, rawPassword string) error {
	incoming := []byte(rawPassword)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}
