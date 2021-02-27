# GotRo

GotRo is abbreviation of `Gotong Royong`. the meaning in `Indonesia`: "do it together", "mutual cooperation". 
This Framework is rewrite of [gokil](//gitlab.com/kokizzu/gokil), that previously use [httprouter](//github.com/julienschmidt/httprouter) but now rewritten using [fasthttprouter](//github.com/buaazp/fasthttprouter). For tutorial, read [this blog post](//kokizzu.blogspot.com/2017/05/gotro-framework-tutorial-go-redis-and.html) (**deprecated**, do not use `W` package for now, wait for full rewrite to [fiber](//gofiber.io) or use `v1.222.1557` if you need to use the old version).

## Versioning

versioning using this format `(YEAR-2020)`.`MDD`.`HMM`,
so for example v1.213.1549 means it was released at `2021-02-13 15:49`

## Design Goal
- As similar as possible to [gokil](//gitlab.com/kokizzu/gokil) that still used by my old company
- Opinionated (choose the best dependency), for example by default uses `int64` and `float64`
- 1-letter supporting package so we only need to write a short common function, such as: `I.ToS(1234)` to convert `int64` to `string`)
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
  - W - Web (the "framework") **STATUS**: usable since 2017-03-08, see `W/example-simplified/`
  - X - Anything (aka `interface{}`)
  - Z - Z-Template Engine, that has syntax similar to ruby string interpolation `#{foo}` or any other that javascript friendly `{/* foo */}`, `[/* bar */]`, `/*! bar */`
- Comment and examples on each type and function, so it can be viewed using godoc, something like: `godoc github.com/kokizzu/gotro/A`

## Status

Usable session adapter:
  - Redis
  - ScyllaDB
  - PostgreSQL
  - AreoSpike
  
Usable database adapter:
  - PostgreSQL

Additional adapter:
  - Tarantool
  - Meilisearch
  
Other than above, you must use officially provided database adapter from respective vendors.

## Benchmark

Benchmarked using [hay](//github.com/rakyll/hey) `-c 255 -n 255000 http://localhost:3001` on i7-4720HQ [gotro](//github.com/kokizzu/gotro) almost 2x faster than [gokil](//gitlab.com/kokizzu/gokil) (23k rps vs 12k rps, thanks to `fasthttp`),
this already includes session loading and template rendering (real-life use case, but with template auto-reloading which should be faster on production mode, since unlike in development mode it doesn't stat disk at all).

## Usage

`go get -u -v github.com/kokizzu/gotro`

## Dependencies

These dependencies automatically installed when you run `go get` (checked using `go list -f '{{join .Deps "\n"}}' |  xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' | uniq`)

```
github.com/OneOfOne/cmap 
github.com/aerospike/aerospike-client-go 
github.com/buaazp/fasthttprouter 
github.com/fatih/color 
github.com/go-sql-driver/mysql 
github.com/jmoiron/sqlx  
github.com/kr/pretty 
github.com/lib/pq 
github.com/mutecomm/go-sqlcipher
github.com/meilisearch/meilisearch-go 
github.com/op/go-logging 
github.com/tdewolff/minify 
github.com/tarantool/go-tarantool
github.com/valyala/fasthttp 
github.com/yosuke-furukawa/json5/encoding/json5 
gopkg.in/redis.v5 
```

## Contributors

- Dikaimin Simon
- Dimas Yudha P
- Rizal Widyarta Gowandy
- Michael Lim
- Pham Hoang Tien
- Devin Yonas

## TODO

- fix mysql adapter so it becomes usable (currently copied from Postgres'), probably wait until mysql has indexable json column, or do alters like scylladb and sqlite
- rewrite W using [fasthttp/router](https://github.com/fasthttp/router) or [fiber](https://gofiber.io/)
- rewrite D using prepared statements, so no more `S.Z`
- use `nikoksr/notify` for notification and mail sending instead of tied to `W`
- possibly refactor move cachedquery, records, etc to D package since nothing different about them, wait for cassandra version
- [Review](//goo.gl/tBkfse) which databases we must support primarily for `D`, that can be silver bullet for extreme cases (high-write: sharding/partitioning and multi-master replication or auto-failover; full-text-search) 
  - [ActorDB](//www.actordb.com) <-- high-write
  - [Cassandra](//cassandra.apache.org) <-- high-write
  - [CockroachDB](//www.cockroachlabs.com) <-- high-read
  - [DGraph](//dgraph.io/)   
  - [CrateDB](//www.crate.io) <-- high-write
  - [GridDB](//griddb.net/en) <-- high-write
  - [GunDB](//gundb.github.io)
  - [Impala](//impala.apache.org)
  - [InfluxDB](//docs.influxdata.com/influxdb)
  - [MariaDB](//mariadb.org) <-- high-read
  - [NebulaGraph](//nebula-graph.io)
  - [OrientDB](//orientdb.com)
  - [PostgreXL](//www.postgres-xl.org) <-- high-write
  - [Riak](//docs.basho.com/riak)
  - [ScyllaDB](//www.scylladb.com) <-- high-write
  - [SingeStore](//www.singlestore.com) <-- high-write
  - [TiDB](//github.com/pingcap/tidb) <-- high-write
  - [TimeScaleDB](//www.timescale.com/) <-- high write
- Review which queuing service we're gonna use ([NATS](//nats.io)), requirement: must support persistence
- Add [ExampleXxx function](//blog.golang.org/examples), getting started and more documentation 
- Add graceful restart (zero downtime deployment): [grace](//github.com/facebookgo/grace) or [endless](//github.com/fvbock/endless) or [overseer](https://github.com/jpillora/overseer)
- Add Catch NotFound (rewrite the `Response.Body`) if no route and static file found
- rewrite router to `fiber` or `fasthttp/router` after Generics support comes up (so we can embed the database connection dependencies inside the context without casting interface)
