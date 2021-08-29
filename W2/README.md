
# Codegen for Web Framework

This framework aims to reduce common task that happened during backend API development which are:

1. generating ORM and migrating model
2. generate for other type of input and output encoding/presentation format (json, msgpack, yaml, etc), command line, or different transport (REST, websocket, gRPC, etc)
3. generate API documentation and clients for different programming languages
4. testable business logic (by untangling `domain/` and decorator/adapter/)

This codegen currently tied to Tarantool (for OLTP use cases), Clickhouse (for OLAP use cases), and Fiber (for web routes) and only tested for generating json using REST transport and plain JS API docs. But you can always create a new codegen for other databases or other web frameworks or other API clients.

## FAQ

Why is it tied to [Tarantool](//www.tarantool.io/en/developers)?
- because it's currently the fastest SQL-OLTP database that also works as in-memory cache, so I can do [integration testing](//kokizzu.blogspot.com/2021/07/mock-vs-fake-and-classical-testing.html) without noticable overhead (can do ~200K writes per second).

Why is it tied to [Clickhouse](//clickhouse.tech)?
- because it's currently the fastest OLAP database (can do ~600K inserts per second).

Why is it tied to [Fiber](//gofiber.io/)?
- because it's currently the fastest minimalist golang framework.

Why still using `encoding/json`?
- because I couldn't find any other faster alternative that can properly parse int64/uint64 (already tried jsoniter and easyjson, both give wrong result for `{"id":"89388457092187654"}` with `json:"id,string"` tag).

## Usage

See `example/` directory for minimum framework template, or if you want to do codegen manually:

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

## Example Generated API Docs

![image](https://user-images.githubusercontent.com/1061610/131266708-44935872-e34a-4538-885a-6056946c9482.png)

