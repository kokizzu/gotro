package S

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
)

func TestStringAES(t *testing.T) {
	for z := 1; z <= 33; z++ {
		pass := RandomPassword(int64(z))
		for y := 1; y < 64; y++ {
			t.Run(fmt.Sprintf("keyLen=%d_textLen=%d", z, y), func(t *testing.T) {
				plain := RandomPassword(int64(y))
				crypt := EncryptAES(plain, pass)
				fmt.Println(len(crypt))
				decrypt := ``
				ok := DecryptAES(crypt, pass, &decrypt)
				assert.True(t, ok)
				assert.Equal(t, plain, decrypt)
			})
		}
	}
}

func TestStructAES(t *testing.T) {
	type Foo struct {
		A string
		B int
	}
	for z := 1; z <= 33; z++ {
		pass := RandomPassword(int64(z))
		t.Run(fmt.Sprintf("keyLen=%d", z), func(t *testing.T) {
			plain := Foo{
				A: RandomPassword(int64(z)),
				B: rand.Int(),
			}
			crypt := EncryptAES(plain, pass)
			fmt.Println(len(crypt))
			decrypt := Foo{}
			ok := DecryptAES(crypt, pass, &decrypt)
			assert.True(t, ok)
			assert.Equal(t, plain, decrypt)
		})
	}
}

func TestMapAES(t *testing.T) {
	for z := 1; z <= 33; z++ {
		pass := RandomPassword(int64(z))
		t.Run(fmt.Sprintf("keyLen=%d", z), func(t *testing.T) {
			plain := map[string]any{
				`A`: RandomPassword(int64(z)),
				`B`: rand.Uint64(), // msgpack always have length, so we cannot use int
			}
			crypt := EncryptAES(plain, pass)
			fmt.Println(len(crypt))
			decrypt := map[string]any{}
			ok := DecryptAES(crypt, pass, &decrypt)
			assert.True(t, ok)
			assert.Equal(t, plain, decrypt)
		})
	}
}
