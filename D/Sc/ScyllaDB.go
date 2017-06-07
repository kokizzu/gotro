package Sc

import (
	_ "github.com/gocql/gocql"
	"github.com/kokizzu/gotro/D"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/T"
	"github.com/kokizzu/gotro/X"
	"sync/atomic"
)

var Z func(string) string
var ZZ func(string) string
var ZJ func(string) string
var ZI func(int64) string
var ZLIKE func(string) string
var ZS func(string) string

var DEBUG bool

func init() {
	Z = S.Z
	ZZ = S.ZZ
	ZJ = S.ZJJ
	ZI = S.ZI
	ZLIKE = S.ZLIKE
	ZS = S.ZS
	DEBUG = D.DEBUG
}

const NANOSEC_LEN = S.MaxStrLenCB63
const ATOMIC_LEN = 3
const SERVER_LEN = 2

var SERVER_ID string // 2 DIGIT, so maximum is 4095 writers/servers, if more needed (eg. 262143, swap the SERVER_LEN and ATOMIC_LEN)
var ATOMIC = uint32(0)

func InitServer(id uint64) {
	max := S.ModCB63[3]
	if id < 1 || id >= max {
		panic(`Sc.SERVER_ID must not be greater than ` + X.ToS(max))
	}
	SERVER_ID = S.EncodeCB63(int64(id), SERVER_LEN)
}

// return 16-byte length id (alignment issue, can be stored on 64+32 bit integer), max suggested: 22-byte (can be stored on 128 bit integer)
func NextId() string {
	if len(SERVER_ID) != 2 {
		panic(`Sc.SERVER_ID length mismatch: ` + SERVER_ID)
	}

	mod := uint32(S.ModCB63[ATOMIC_LEN])
	atom := int64(atomic.AddUint32(&ATOMIC, 1) % mod)
	return S.EncodeCB63(T.UnixNano(), S.MaxStrLenCB63) + S.EncodeCB63(atom, ATOMIC_LEN) + SERVER_ID
}
