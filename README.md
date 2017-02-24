# GotRo

GotRo is abbreviation of `Gotong Royong`. the meaning in `Indonesia`: "do it together", "mutual cooperation". 
This Framework is rewrite of [gokil](gitlab.com/kokizzu/gokil), that previously use [fasthttp](//github.com/julienschmidt/httprouter) but now rewritten using [fasthttprouter](github.com/buaazp/fasthttprouter)

## Design Goal
- As similar as possible to Elixir's Phoenix Framework
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
  - W - Web (the "framework")
  - X - Anything (aka `interface{}`)
  - Z - Z-Template Engine, that has syntax similar to ruby string interpolation `#{foo}` or any other that javascript friendly `{/* foo */}`, `[/* bar */]`, `/*! bar */`
- Comment and examples on each type and function, so it can be viewed using godoc, something like: `godoc github.com/kokizzu/gotro/A`

## Usage

`go get -u -v github.com/kokizzu/gotro`

## Dependencies

These dependencies automatically installed when you run `go get`

- [FastHttp](//github.com/valyala/fasthttp)
- [Logging](//github.com/op/go-logging)
- [Pretty Print Variables](//github.com/kr/pretty)
- [Terminal Color](//github.com/fatih/color)

## TODO

- Add httprouter to `W`, add logging, add panic handling (stacktrace censoring), add session loading
- List most of Phoenix features and add it to `W`
- [Review](//goo.gl/tBkfse) which databases we must support primarily for `D` (drop ones that hard to install), that can be silver bullet for extreme cases (high-write: sharding/partitioning and multi-master replication or auto-failover; full-text-search) 
- Review which databases we must support primarily for `D` (drop ones that hard to install), that can be silver bullet for extreme cases (high-write: sharding/partitioning and multi-master replication or auto-failover; full-text-search) 
  - [ArangoDB](//www.arangodb.com/)
  - [Cassandra](//cassandra.apache.org)
  - [Couchbase](//couchbase.com)
  - [CouchDB](//couchdb.apache.org/)
  - [CockroachDB](//www.cockroachlabs.com/)
  - [GunDB](//gundb.github.io)
  - [Impala](//impala.apache.org/)
  - [InfluxDB](//docs.influxdata.com/influxdb)
  - [MariaDB](//mariadb.org)
  - [OrientDB](//orientdb.com)
  - [PostgreXL](//www.postgres-xl.org/)
  - [Riak](//docs.basho.com/riak)
  - [ScyllaDB](//www.scylladb.com)
- Review which queuing service we're gonna use ([NSQ](//nsq.io), [Redis](//redis.io)), requirement: must support persistence
- Add [ExampleXxx function](//blog.golang.org/examples), getting started and more documentation 
- Create example API App
- Add graceful restart (zero downtime deployment): [grace](//github.com/facebookgo/grace) or [endless](//github.com/fvbock/endless)
- Write a book for about [Advanced Programming](//goo.gl/X4BIlM), [Database Systems](//goo.gl/uR8iVB) and [Web Programming](//goo.gl/Bl3fPE) that includes this framework