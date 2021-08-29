
# Example gotro/W2 Project

How to use this template?
- install Go 1.16+ and clone this repo with `--depth 1` flag
- copy this `example` directory to another folder (rename to `projectName`)
- `go mod init projectName`
- replace all word `github.com/kokizzu/gotro/W2/example` and `example` with `projectName`

How to develop?
- modify or create new `model/m*/*_tables.go`, then run `make gen-orm`
- create a new `*_In`, `*_Out`, `*_Url`, and `*` business logic methods inside `domain`, then run `make gen-route`

## Setup

```bash
# install tools required for codegen
make setup-deps

# install reverse proxy
make setup-webserver

# install dependencies for web frontend (Svelte with Vite build system): localhost:3000
make webclient

# start dependencies (Tarantool, Clickhouse): localhost:3301, localhost:9000
make compose

# run api server (Go with Air autorecompile): localhost:9090
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

# generate route (after add new _In _Out _Url and business logic method on domain/*.go)
make gen-route
```

