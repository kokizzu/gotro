
# Example gotro/W2 Project

How to use this template?
- install Go 1.16+ and clone this repo with `--depth 1` flag
- copy this `example` directory to another folder (rename to `projectName`)
- `go mod init projectName`
- replace all word `github.com/kokizzu/gotro/W2/example` and `example` with `projectName`

How to develop?
- modify or create new `model/m*/*_tables.go`, then run `make gen-orm` (will generate ORM), you may only add column/field at the end of model.
- create a new `*_In`, `*_Out`, `*_Url`, and the business logic methods inside `domain`, then run `make gen-route` (will generate routes and API docs).
- create an integration/unit test to make sure that your code is correct
![image](https://user-images.githubusercontent.com/1061610/131266911-281090aa-f062-43eb-80cf-9ba561e019d2.png)

How to release?
- change `production/` configuration values
- setup the server, ssh to server and run `setup_server.sh`
- cd `production`, run `./sync_service.sh`
- run `./deploy_prod.sh`

## Setup

```bash
# install tools required for codegen
make setup-deps

# install reverse proxy
make setup-webserver

# install dependencies for web frontend (Svelte with Vite build system): localhost:3000
make webclient

# start dependencies (Tarantool, Clickhouse, mailhog): localhost:3301, localhost:9000, localhost:1025
make compose

# run api server (Go with Air auto-recompile): localhost:9090
make apiserver

# run reverse proxy (Caddy): localhost:80
make reverseproxy
```

## Usage

```bash
# connect to OLTP database
tarantoolctll connect 3301

# connect to OLAP database
clickhouse-client

# generate ORM (after add new table or columns on models/m*/*_tables.go)
make gen-orm

# generate route (after add new _In+_Out struct, _Url const and business logic method on domain/*.go)
make gen-route
```

## Gotchas

- Call `wc*.Set*` instead of direct assignment (`=`) before calling `wc*.DoUpdateBy*`
- Clickhouse inserts are buffered using [chTimedBuffer](//github.com/kokizzu/ch-timed-buffer), so you must wait ~1s to ensure it's flushed
- Clickhouse have eventual consistency, so you must use `FINAL` query to make sure it's committed
- You cannot change Tarantool's datatype
- You cannot change Clickhouse's ordering keys datatype
- Currently migration only allowed for adding columns/fields at the end (you cannot insert new column in the middle/begginging)
- All Tarantool's columns always set not null after migration
- Tarantool does not support client side transaction (so you must use Lua or split into SAGAs)

## TODOs

- Add SEO pre-render: [Rendora](//github.com/rendora/rendora)
- Add search-engine: [TypeSense](//typesense.org/) example
- Add persisted cache: [IceFireDB](https://github.com/gitsrc/IceFireDB) or [Aerospike](https://aerospike.com/)
- Add external storage upload example (minio? wasabi?)
- Replace LightStep with [SigNoz](https://github.com/SigNoz/signoz)
- Add more deployment script with [LXC/LXD share](https://bobcares.com/blog/how-to-setup-high-density-vps-hosting-using-lxc-linux-containers-and-lxd/) for single server multi-tenant
