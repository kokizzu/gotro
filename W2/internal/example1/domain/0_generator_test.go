package domain_test

import (
	"example1/conf"
	"github.com/kokizzu/gotro/W2"

	"testing"
)

//go:generate go test -run=XXX -bench=Benchmark_Generate_WebApiRoutes_CliArgs
//go:generate go test -run=XXX -bench=Benchmark_Generate_SvelteApiDocs

func Benchmark_Generate_WebApiRoutes_CliArgs_GraphQL(b *testing.B) {
	W2.GenerateFiberAndCli(&W2.GeneratorConfig{
		ProjectName: conf.PROJECT_NAME,
	})
	b.SkipNow()
}

func Benchmark_Generate_SvelteApiDocs(b *testing.B) {
	W2.GenerateApiDocs(&W2.GeneratorConfig{
		ProjectName: conf.PROJECT_NAME,
	})
	b.SkipNow()
}
