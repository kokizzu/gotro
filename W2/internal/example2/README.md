# Example2 gotro/W2 Project

How to use this template?
- install Go 1.16+ and clone this repo with `--depth 1` flag
- copy this `example2` directory to another folder (rename to `projectName`)
- `go mod init projectName`
- replace all word `github.com/kokizzu/gotro/W2/internal/example2` and `example2` with `projectName`

How to develop?
- modify or create new `model/m*/*_tables.go`, then run `make gen-orm` (will generate ORM), you may only add column/field at the end of model.
- create a new `*_In`, `*_Out`, `*_Url`, and the business logic methods inside `domain`, then run `make gen-route` (will generate routes and API docs).
- create an integration/unit test to make sure that your code is correct
![image](https://user-images.githubusercontent.com/1061610/131266911-281090aa-f062-43eb-80cf-9ba561e019d2.png)

## Project Structure

- `conf`: shared configuration
- `main`: inject all dependencies
- `svelte`: web frontend
- `deploy`: scripts for deploying to production

MVC-like structure

- `presentation` -> serialization and transport
- `domain` -> business logic, DTO
- `model` -> persistence/3rd party endpoints, DAO

## CQRS

separates read/query and write/command into different database/connection.

- `rq` = read/query (OLTP) -> to tarantool replica
- `wc` = write/command (OLTP) -> to tarantool master
- `sa` = statistics/analytics (OLAP) -> to clickhouse

## Start dev mode

```shell
# first time for development
make setup

# start docker
docker compose up # or docker-compose up

# start frontend auto build 
cd svelte
pnpm install 
bun watch 

# do migration (first time, or everytime tarantool/clickhouse docker deleted, 
# or when there's new migration)
go run main.go migrate
go run main.go import

# start golang backend server, also serving static html
air web

# manually
go run main.go web

# start reverse proxy, if you need test oauth
caddy run # foreground
caddy start # background
```

## Execute Tarantool on Docker

```shell
docker exec -it example2-tarantool1-1 tarantoolctl connect example2T:example2PT@127.0.0.1:3301
# box.execute [[ SELECT * FROM "users" LIMIT 10 ]]
```

## Execute Clickhouse on Docker

```shell
docker exec -it example2-clickhouse1-1 clickhouse-client -u example2C
# SELECT * FROM actionLogs ORDER BY createdAt DESC LIMIT 10;
```

## Generate ORM

```shell
# input: 
# - model/m*.go

./gen-orm.sh
# or
cd model
go test -bench=BenchmarkGenerateOrm
# then go generate each file

# output:
# - model/m*/rq*/*.go  # -- read/query models
# - model/m*/wc*/*.go  # -- write/command mutation models
# - model/m*/sa*/*.go  # -- statistic/analytics models
```

## Generate Views

```shell
# input: 
# - domain/*.go
# - model/m*/*/*.go
# - svelte/*.svelte

./gen-views.sh
# or
cd presentation
go test -bench=BenchmarkGenerateViews 

# output:
# - presentation/actions.GEN.go     # -- all possible commands
# - presentation/api_routes.GEN.go  # -- automatic API routes
# - presentation/web_view.GEN.go    # -- all template that can be used in web_static.go
# - presentation/cmd_run.GEN.go     # -- all CLI commands
# - svelte/jsApi.GEN.js             # -- all API client SDK 
```

## Importing data

```shell
go run main.go migrate
go run main.go import
go run main.go import_location # require google API key

```

## Note

```shell
# docker spawning failed (because test terminated improperly), run this:
alias dockill='docker kill $(docker ps -q); docker container prune -f; docker network prune -f'
```

## FAQ

- **Q**: where to put SSR?
  - **A**: `presentation/web_static.go`
- **Q**: got error `there is no space with name [tableName], table default.
  [tableName] does not exists`
  - **A**: run `go run main.go migrate` to do migration
- **Q**: got error `Command 'caddy' not found`
  - **A**: install [caddy](//caddyserver.com/docs/install)
- **Q**: got error `Command 'air' not found`
  - **A**: install [air](//github.com/cosmtrek/air) or `make setup`
- **Q**: got error `Command 'replacer' not found`
  - **A**: install [replacer](//github.com/kokizzu/replacer) or `make setup`
- **Q**: got error `Command 'gomodifytags' not found`
  - **A**: install [gomodifytags](//github.com/fatih/gomodifytags) or `make setup`
- **Q**: got error `Command 'farify' not found`
  - **A**: install [farify](//github.com/akbarfa49/farify) or `make setup`
- **Q**: got error `Command 'goimports' not found`
  - **A**: install [goimports](//cs.opensource.google/go/x/tools) or `make setup`
- **Q**: where to put secret that I don't want to commit?
  - **A**: on `.env.override` file
- **Q**: got error `.env.override` no such file or directory
  - **A**: create `.env.override` file
- **Q**: got error `failed to stat the template: index.html`
  - **A**: run `cd svelte; npm run watch` at least once
- **Q**: got error `TarantoolConf) Connect: dial tcp 127.0.0.1:3301: connect: connection refused"`
  - **A**: run `docker-compose up`
- **Q**: got error `ClickhouseConf) Connect: dial tcp 127.0.0.1:9000: connect: connection refused`
  - **A**: run `docker-compose up`
- **Q**: got error `docker.errors.DockerException: Error while fetching server API version: ('Connection aborted.', FileNotFoundError(2, 'No such file or directory'))`
  - **A**: make sure docker service is up and running
- **Q**: what's normal flow of development?
  - **A**: 
      1. create new/modify model on `model/m[schema]/[schema]_tables.go` folder, create benchmark function to generate and migrate the tables in `RunMigration` function.
      2. run `./gen-orm.sh`, create helper function on `model/w[schema]/[rq|wc|sa][schema]/[schema]_helper.go` or 3rd party wrapper in `model/x[repo]/x[provider].go`
      3. create a role in `domain/[role].go` containing all business logic for that role
      4. write test in `domain/[role]_test.go` to make sure all business requirement are met
      5. generate domain routes `cd presentation; go get -bench=BenchmarkGenerateViews`, start web service `air web`
      6. write frontend on `svelte/`, start frontend service `cd svelte; npm run watch`
      7. generate frontend helpers `cd presentation; go get -bench=BenchmarkGenerateViews`
      8. write SSR if needed on `presentation/web_static.go`
- **Q**: how to add additional tech stack?
  - **A**: put on `docker-compose.yml` and add to `domain/0_main_test.go` so it would run on integration test. Create the `conf/[provider].go` and `model/x[repo]/x[provider].go` to wrap the 3rd party connector. 
- **Q**: want to change the generated views?
  - **A**: 
      1. change the template on `presentation/1_codegen_test.go`
      2. run `cd presentation; go get -bench=BenchmarkGenerateViews`
- **Q**: want to change generated ORM or schema migration has bug?
  - **A**: create a pull request to [gotro](//github.com/kokizzu/gotro) 
- **Q**: generated html have bug?
  - **A**: create a pull request to [svelte-mpa](//github.com/kokizzu/svelte-mpa) or [svelte](//github.com/sveltejs/svelte)
- **Q**: where is the devlog?
  - **A**: on [youtube livestream](//www.youtube.com/@kokizzu/streams)
- **Q**: why secrets not encrypted?
  - **A**: it's ok for PoC phase, since it's listen to localhost
- **Q**: run test against local `docker compose` instead of `dockertest`?
  - **A**: set env or export `USE_COMPOSE=x` before running test
- **Q**: how SSR works?
  - **A**: `npm run watch` will convert `.svelte` files into `.html`, `.
    /gen-views.sh` will generate `vew_view.GEN.go` that can be called by `web_static.go` to render the `.html` files
- **Q**: how route generator works?
  - **A**: everytime you make `type XXIn`, `type XXOut`, `const XXAction`, 
    and `func (d Domain) XX(in XXIn) XXOut` inside `domain/` it will be 
    automatically 
    added 
    to `api_routes.GEN.go`
- **Q**: how orm generator works?
  - **A**: create `model/m[schema]/[schema]_tables.go` and `model/m[schema]/
    [schema]_tables_test.go` then run `./gen-orm.sh`, it would generate 
    `model/[rq|wc|sa][schema].go`, you can extend the method of generated 
    structs (`[rq|wc|sa][schema]`) on another file inside the same package.
- **Q**: how basic form and list/table works?
  - **A**: create a `[schema]Meta` on `domain/`, then just call your query 
    method based on `zCrud.Pager` (it would generate the proper SQL query), 
    then use svelte component that can render the form and table/list for you.
- **Q**: when to use each storage engine?
  - **A**: `memtx` used for kv query pattern, anything that often being read 
    and updated (eg. transactions), `vinyl` used for range queries, anything 
    that rarely being updated (eg. mutation log, history), `clickhouse` used 
    for analytics queries pattern, anything that will never being updated 
    ever (eg. action logs, events)
- **Q**: got `example2/tmpdb/*: open /*/example2/*: permission denied` on `go mod tidy`
  - **A**: run `sudo chmod a+rwx -R tmpdb` or `make modtidy`
