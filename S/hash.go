package S

import "github.com/zeebo/xxh3"

// XXH3 fastest 64-bit hash function
// https://github.com/jacobmarble/go-checksum-benchmark
// https://github.com/Cyan4973/xxHash
func XXH3(s string) uint64 {
	return xxh3.HashString(s)
}
