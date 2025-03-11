package S

import (
	"fmt"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
		C int8
		D int16
		E int32
		F int64
		G uint
		H uint8
		I uint16
		J uint32
		K uint64
		L float32
		M float64
		N time.Time
	}
	for z := 1; z <= 33; z++ {
		pass := RandomPassword(int64(z))
		t.Run(fmt.Sprintf("keyLen=%d", z), func(t *testing.T) {
			plain := Foo{
				A: RandomPassword(int64(z)),
				B: rand.Int(),
				C: int8(rand.Int()),
				D: int16(rand.Int()),
				E: int32(rand.Int()),
				F: rand.Int64(),
				G: uint(rand.Uint32()),
				H: uint8(rand.Uint32()),
				I: uint16(rand.Uint32()),
				J: rand.Uint32(),
				K: rand.Uint64(),
				L: rand.Float32(),
				M: rand.Float64(),
				N: time.Unix(time.Now().Unix(), 0),
				// time.Now always have wall and ext that always deserialized wrongly
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
				// msgpack always have integer size, so we cannot use int or uint
				`C`: int8(rand.Int()),
				`D`: int16(rand.Int()),
				`E`: int32(rand.Int()),
				`F`: rand.Int64(),
				`H`: uint8(rand.Uint32()),
				`I`: uint16(rand.Uint32()),
				`J`: rand.Uint32(),
				`K`: rand.Uint64(),
				`L`: rand.Float32(),
				`M`: rand.Float64(),
				`N`: time.Unix(time.Now().Unix(), 0),
				// time.Now always have wall and ext that always deserialized wrongly
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
