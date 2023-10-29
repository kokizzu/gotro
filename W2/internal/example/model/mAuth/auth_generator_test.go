package mAuth

import "testing"

//go:generate go test -run=XXX -bench=Benchmark_GenerateOrm

func Benchmark_GenerateOrm(b *testing.B) {
	GenerateORM()
	b.SkipNow()
}
