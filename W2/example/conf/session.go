package conf

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"strings"

	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/lexid"
	"github.com/mojura/enkodo"
	"github.com/segmentio/fasthash/fnv1a"
	"github.com/zeebo/xxh3"
)

const CookieName = `auth:`
const CookieLogoutValue = `LOGOUT`
const CookieDays = 30

type Session struct {
	UserId    uint64
	ExpiredAt int64
	Email     string
}

func (u *Session) MarshalEnkodo(enc *enkodo.Encoder) (err error) {
	enc.Uint64(u.UserId)
	enc.Int64(u.ExpiredAt)
	enc.String(u.Email)
	return
}

func (u *Session) UnmarshalEnkodo(dec *enkodo.Decoder) (err error) {
	if u.UserId, err = dec.Uint64(); err != nil {
		return
	}
	if u.ExpiredAt, err = dec.Int64(); err != nil {
		return
	}
	if u.Email, err = dec.String(); err != nil {
		return
	}
	return
}

func createHash(key1, key2 string) string {
	res := xxh3.HashString128(key1 + PROJECT_NAME + key2) // PROJECT_NAME = salt, if you change this, all token will be invalidated
	const x = 256
	return hex.EncodeToString([]byte{
		byte(res.Hi >> (64 - 8) % x),
		byte(res.Hi >> (64 - 16) % x),
		byte(res.Hi >> (64 - 24) % x),
		byte(res.Hi >> (64 - 32) % x),
		byte(res.Hi >> (64 - 40) % x),
		byte(res.Hi >> (64 - 48) % x),
		byte(res.Hi >> (64 - 56) % x),
		byte(res.Hi >> (64 - 64) % x),
		byte(res.Lo >> (64 - 8) % x),
		byte(res.Lo >> (64 - 16) % x),
		byte(res.Lo >> (64 - 24) % x),
		byte(res.Lo >> (64 - 32) % x),
		byte(res.Lo >> (64 - 40) % x),
		byte(res.Lo >> (64 - 48) % x),
		byte(res.Lo >> (64 - 56) % x),
		byte(res.Lo >> (64 - 64) % x),
	})
}

const TokenSeparator = `|`

func (s *Session) Encrypt(userAgent string) string {
	key1 := lexid.NanoID()
	key2 := S.EncodeCB63(int64(fnv1a.HashString64(userAgent)), 1)
	block, err := aes.NewCipher([]byte(createHash(key1, key2)))
	L.PanicIf(err, `aes.NewCipher`)
	gcm, err := cipher.NewGCM(block)
	L.PanicIf(err, `cipher.NewGCM`)
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	L.PanicIf(err, `io.ReadFull`)
	buffer := bytes.NewBuffer(nil)
	// Create a writer
	w := enkodo.NewWriter(buffer)
	// Encode user to buffer
	err = w.Encode(s)
	L.PanicIf(err, `w.Encode`)
	ciphertext := gcm.Seal(nonce, nonce, buffer.Bytes(), nil)
	return key1 + TokenSeparator + hex.EncodeToString(ciphertext) + TokenSeparator + key2
}

func (s *Session) Decrypt(sessionToken string, userAgent string) bool {
	strs := strings.Split(sessionToken, TokenSeparator)
	if len(strs) != 3 {
		L.Print(`incorrect token segment length: ` + I.ToStr(len(strs)))
		return false
	}
	uaHash := S.EncodeCB63(int64(fnv1a.HashString64(userAgent)), 1)
	if strs[2] != uaHash {
		L.Print(`userAgent mismatch: ` + strs[2] + ` <> ` + uaHash)
		return false
	}
	data, err := hex.DecodeString(strs[1])
	if L.IsError(err, `hex.DecodeString`) {
		return false
	}
	key := []byte(createHash(strs[0], strs[2]))
	block, err := aes.NewCipher(key)
	if L.IsError(err, ` aes.NewCipher`) {
		return false
	}
	gcm, err := cipher.NewGCM(block)
	if L.IsError(err, `cipher.NewGCM`) {
		return false
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if L.IsError(err, `gcm.Open`) {
		return false
	}
	err = enkodo.Unmarshal(plaintext, s)
	return !L.IsError(err, `enkodo.Unmarshal`)
}
