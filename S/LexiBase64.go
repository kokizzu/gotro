package S

import (
	"github.com/kokizzu/gotro/L"
	"math/rand"
)

const i2c_cb63 = `-0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz`

var c2i_cb63 map[rune]int64

func init() {
	c2i_cb63 = map[rune]int64{}
	for i, ch := range i2c_cb63 {
		c2i_cb63[ch] = int64(i)
	}
}

// convert integer to custom base-63 encoding that lexicographically correct, positive integer only
//  0       -
//  1..10   0..9
//  11..36  A..Z
//  37      _
//  38..63  a..z
//  S.EncodeCB63(11) // `A`
//  S.EncodeCB63(1)) // `1`
func EncodeCB63(id int64, min_len int) string {
	if min_len < 1 {
		min_len = 1
	}
	str := make([]byte, 0, 12)
	for id > 0 {
		mod := rune(id % 64)
		str = append(str, i2c_cb63[mod])
		id /= 64
	}
	for len(str) < min_len {
		str = append(str, i2c_cb63[0])
	}
	l := len(str)
	for i, j := 0, l-1; i < l/2; i, j = i+1, j-1 {
		str[i], str[j] = str[j], str[i]
	}
	return string(str)
}

// convert custom base-63 encoding to int64
func DecodeCB63(str string) (int64, bool) {
	res := int64(0)
	for _, ch := range str {
		res *= 64
		val, ok := c2i_cb63[ch]
		if L.CheckIf(!ok, `Invalid character for CB63: `+string(ch)) {
			return -1, false
		}
		res += val
	}
	return res, true
}

// random CB63 n-times, the result is n*11 bytes
func RandomCB63(len int64) string {
	res := ``
	for z := int64(0); z < len; z++ {
		res += EncodeCB63(rand.Int63(), 11)
	}
	return res
}
