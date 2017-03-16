# GotRo

GotRo is abbreviation of `Gotong Royong`. the meaning in `Indonesia`: "do it together", "mutual cooperation". 
This Framework is rewrite of [gokil](//gitlab.com/kokizzu/gokil), that previously use [httprouter](//github.com/julienschmidt/httprouter) but now rewritten using [fasthttprouter](//github.com/buaazp/fasthttprouter)

## Design Goal
- As similar as possible to Gokil
- Opinionated (choose the best dependency), for example by default uses int64 and float64
- 1-letter supporting package so we short common function, such as: `I.ToS(1234)` to convert `int64` to `string`)
  - A - Array
  - B - Boolean
  - C - Character (or Rune)
  - D - Database
  - F - Floating Point
  - L - Logging
  - M - Map
  - I - Integer
  - S - String
  - T - Time (and Date)
  - W - Web (the "framework") **STATUS**: usable since 2017-03-08, see `W/example/` 
  - X - Anything (aka `interface{}`)
  - Z - Z-Template Engine, that has syntax similar to ruby string interpolation `#{foo}` or any other that javascript friendly `{/* foo */}`, `[/* bar */]`, `/*! bar */`
- Comment and examples on each type and function, so it can be viewed using godoc, something like: `godoc github.com/kokizzu/gotro/A`

## Benchmark

Benchmarked using [hay](//github.com/rakyll/hey) `-c 255 -n 255000 http://localhost:3001` on i7-4720HQ [gokil](//github.com/kokizzu/gotro) almost 2x faster than [gokil](//gitlab.com/kokizzu/gokil) (23k rps vs 12k rps, thanks to `fasthttp`),
this already includes session loading and template rendering (real-life use case, but with template auto-reloading which should be faster on production mode).

## Usage

`go get -u -v github.com/kokizzu/gotro`

## Dependencies

These dependencies automatically installed when you run `go get` (checked using `go list -f '{{join .Deps "\n"}}' |  xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' | uniq`)

```
github.com/OneOfOne/cmap 
github.com/aerospike/aerospike-client-go 
github.com/buaazp/fasthttprouter 
github.com/davecgh/go-spew/spew 
github.com/fatih/color 
github.com/go-sql-driver/mysql 
github.com/jmoiron/sqlx  
github.com/kr/pretty 
github.com/lib/pq 
github.com/mutecomm/go-sqlcipher 
github.com/op/go-logging 
github.com/tdewolff/minify 
github.com/valyala/fasthttp 
github.com/yosuke-furukawa/json5/encoding/json5 
gopkg.in/redis.v5 
```

## TODO

- fix mysql adapter so it becomes usable (currently copied from Postgres'), probably wait until mysql has indexable json column, or do alters like scylladb and sqlite
- possibly refactor move cachedquery, records, etc to D package since nothing different about them, wait for cassandra version
- [Review](//goo.gl/tBkfse) which databases we must support primarily for `D` (drop ones that hard to install), that can be silver bullet for extreme cases (high-write: sharding/partitioning and multi-master replication or auto-failover; full-text-search) 
  - [ArangoDB](//www.arangodb.com)
  - [Cassandra](//cassandra.apache.org) <-- high-write
  - [Couchbase](//couchbase.com)
  - [CouchDB](//couchdb.apache.org)
  - [CockroachDB](//www.cockroachlabs.com) <-- high-read
  - [GunDB](//gundb.github.io)
  - [Impala](//impala.apache.org)
  - [InfluxDB](//docs.influxdata.com/influxdb)
  - [MariaDB](//mariadb.org) <-- high-read
  - [OrientDB](//orientdb.com)
  - [PostgreSQL](//www.postgresql.org) <-- high-read
  - [PostgreXL](//www.postgres-xl.org) <-- high-write
  - [Riak](//docs.basho.com/riak)
  - [ScyllaDB](//www.scylladb.com) <-- high-write
  - [TiDB](//github.com/pingcap/tidb) <-- high-write
- Review which queuing service we're gonna use ([NATS](//nats.io)), requirement: must support persistence
- Add [ExampleXxx function](//blog.golang.org/examples), getting started and more documentation 
- Create example API App
- Add graceful restart (zero downtime deployment): [grace](//github.com/facebookgo/grace) or [endless](//github.com/fvbock/endless)
- Write a book for about [Advanced Programming](//goo.gl/X4BIlM), [Database Systems](//goo.gl/uR8iVB) and [Web Programming](//goo.gl/Bl3fPE) that includes this framework
- Add Catch NotFound (rewrite the Response.Body) if no route and static file found
- Check why the performance worse than httprouter for `siege -b`