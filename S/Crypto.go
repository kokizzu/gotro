package S

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/jxskiss/base62"

	"github.com/kokizzu/gotro/L"
)

// symmetric cryptography for string using golang

func createKey(key string) cipher.Block {
	if key == `` {
		key = `1234567890abcdef`
	}
	kLen := len(key)
	if kLen < 16 {
		for len(key) < 16 {
			key += key
		}
		key = key[:16]
	} else if kLen < 24 {
		for len(key) < 24 {
			key += key
		}
		key = key[:24]
	} else if kLen < 32 {
		for len(key) < 32 {
			key += key
		}
		key = key[:32]
	} else if kLen > 32 {
		key = key[:32]
	}
	block, err := aes.NewCipher([]byte(key))
	L.PanicIf(err, `aes.NewCipher`)
	return block
}

func EncryptAES(in any, key string) (encryptedStr string) {
	block := createKey(key)
	plaintext, err := msgpack.Marshal(in)
	L.PanicIf(err, `msgpack.Marshal`)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	L.PanicIf(err, `io.ReadFull rand.Reader`)

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base62.EncodeToString(ciphertext)
}

func DecryptAES(encryptedStr, key string, out any) bool {
	block := createKey(key)
	ciphertext, err := base62.DecodeString(encryptedStr)
	if L.IsError(err, `base64.DecodeString`) {
		return false
	}
	if len(ciphertext) < aes.BlockSize {
		return false
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	err = msgpack.Unmarshal(ciphertext, out)
	return !L.IsError(err, `msgpack.Unmarshal`)
}
