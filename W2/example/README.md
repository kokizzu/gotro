
# Example W2 Project

## Setup

```bash
# install dependencies for web frontend
make webclient

# run api server
make apiserver

# run reverse proxy
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

