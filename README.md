# GotRo

GotRo is abbreviation of `Gotong Royong`. the meaning in `Indonesia`: "do it together", "mutual cooperation". 
This Framework is rewrite of [gokil](gitlab.com/kokizzu/gokil), that previously use [fasthttp](https://github.com/julienschmidt/httprouter) but now rewritten using [fasthttprouter](github.com/buaazp/fasthttprouter)

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

## TODO

- Add httprouter to `W`, add logging, add panic handling (stacktrace censoring), add session loading
- List most of Phoenix features and add it to `W`
- Review which databases we must support primarily for `D` (drop ones that hard to install), that can be silver bullet for extreme cases (high-write: sharding/partitioning and multi-master replication or auto-failover; full-text-search) 
  - [ArangoDB](https://www.arangodb.com/)
  - [ChronicleMap](http://chronicle.software/products/chronicle-map/)
  - [Couchbase](http://couchbase.com)
  - [CouchDB](http://couchdb.apache.org/)
  - [CockroachDB](https://www.cockroachlabs.com/)
  - [Ellasandra](https://github.com/strapdata/elassandra)
  - [GunDB](http://gundb.github.io)
  - [Impala](http://impala.apache.org/)
  - [InfluxDB](https://docs.influxdata.com/influxdb)
  - [OrientDB](http://orientdb.com)
  - [PostgreXL](http://www.postgres-xl.org/)
  - [Riak](http://docs.basho.com/riak)
  - [ScyllaDB](http://www.scylladb.com)
  - [Titan](http://titan.thinkaurelius.com/)
- Add [Example function](https://blog.golang.org/examples) and more documentation 
- Create example API App
- Write a book for about [Advanced Programming](https://goo.gl/X4BIlM), [database systems](https://goo.gl/uR8iVB) and [web programming](https://goo.gl/Bl3fPE)