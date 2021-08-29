
# Codegen for Web Framework

## Usage

1. create a test file `0_generator_test.go` inside your `domain/` project folder

```go
package domain_test

import (
	"testing"
	`github.com/kokizzu/gotro/WW`
)

//go:generate go test -run=XXX -bench=Benchmark_Generate_WebApiRoutes_CliArgs
//go:generate go test -run=XXX -bench=Benchmark_Generate_SvelteApiDocs

func Benchmark_Generate_WebApiRoutes_CliArgs(b *testing.B) {
	WW.GenerateFiberAndCli(&WW.GeneratorConfig{
		ProjectName: PROJECT_NAME,
	})
	b.SkipNow()
}

func Benchmark_Generate_SvelteApiDocs(b *testing.B) {
	WW.GenerateApiDocs(&WW.GeneratorConfig{
		ProjectName: PROJECT_NAME,
	})
	b.SkipNow()
}

```

2. create a makefile to do the codegen
```Makefile

gen-route:
	cd domain ; rm -f *MSG.GEN.go 
	cd domain ; go test -bench=Benchmark_Generate_WebApiRoutes_CliArgs 0_generator_test.go
	cd domain ; cat *.go | grep '//go:generate ' | cut -d ' ' -f 2- | sh -x
	cd domain ; go test -bench=Benchmark_Generate_SvelteApiDocs 0_generator_test.go

```

3. run `make gen-route`

this would create few generated file:

```
main_cli_args.GEN.go --> cli arguments handler
main_restApi_routes.GEN.go --> used to generating fiber route handlers
svelte/src/pages/api.js --> used for generating API client
```
