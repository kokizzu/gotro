
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
- setup the server, ssh to the production/staging server and run `setup_server.sh`
- cd `production`, run `./sync_service.sh` 
- run `./deploy_prod.sh`

How to do multi server?

- replace [id64](//github.com/kokizzu/lexid) with string, for example [lexid](//github.com/kokizzu/lexid) or standard `uuid`
- add an environment variable for SERVER_ID, and init it as `lexid.ServerId`
- before running deployment script, make sure to append environment variable SERVER_ID that are must unique per server

How to do multi database?

- see [this](//kokizzu.blogspot.com/2021/05/easy-tarantool-clickhouse-replication-setup.html) blog post

## Directory Structure

- `3rdparty` - all third party wrapper should be here as a subfolder
- `conf` - all configuration constants
- `domain` - contains your business logic, these are the one that should be integration/unit tested
- `model` - contains your domains' data store
  - `m[Domain]` - contains data store that should be grouped inside that domain
    - `rq[Domain]` - read query (R from CQRS), you can add a new file here to extend the default ORM
    - `sa[Domain]` - statistics analytics (event source), you can add a new file here to extend the default ORM
    - `wc[Domain]` - write command (C from CQRS), you can add a new file here to extend the default ORM
    - `*_table.go` - the schema file for that domain, to generate the ORM and as an input for migration
- `production` - scripts and env for deploying to production
- `svelte` - frontend (can be replaced with any framework) 

outer files:
- `main_*.GEN.go` - will be generated per transport/presentation/adapter (eg. gRPC, REST, WebSocket, CLI, etc)

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
tarantoolctl connect 3301

# connect to OLAP database
clickhouse-client

# generate ORM (after add new table or columns on models/m*/*_tables.go)
make gen-orm

# generate route (after add new _In+_Out struct, _Url const and business logic method on domain/*.go)
make gen-route
```

## Gotchas

- Calling direct assignment (`=`) instead of `wc*.Set*()` before calling `wc*.DoUpdateBy*()` will do nothing, as direct assignment does not append mutation property
```
# proper way to update
x := mAuth.NewUsersMutator(s.Taran)
x.Id = ...
if !x.FindById() {
   // not found
   return // or x.DoInsert() then continue with update
}
x.SetBla(..) 
x.SetFoo(..)
x.SetBar(..)
x.SetBaz(..)
x.SetBaz(..) // calling twice on the same column will cause DoUpdate to fail
if !x.DoUpdateById() {
   // failed to update
}

# but if you need only insert or replace, you can use = directly
x := mAuth.NewUsersMutator(s.Taran)
x.Bla = ..
x.Foo = ..
x.Bar = ..
x.Baz = ..
x.DoInsert() or x.DoReplace() // calling DoUpdateBy*() will do nothing, since mutation property only set when calling .Set*() method
```
- Clickhouse inserts are buffered using [chTimedBuffer](//github.com/kokizzu/ch-timed-buffer), so you must wait ~1s to ensure it's flushed
- Clickhouse have eventual consistency, so you must use `FINAL` query to make sure it's force-committed
- You cannot change Tarantool's datatype
- You cannot change Clickhouse's ordering keys datatype
- Currently migration only allowed for adding columns/fields at the end (you cannot insert new column in the middle/begginging)
- All Tarantool's columns always set not null after migration (I hate null values XD)
- Tarantool does not support client side transaction (so you must use [Lua](//www.tarantool.io/en/doc/latest/book/box/atomic/) or 2PC or split into SAGAs)
- Current parser/codegen does not allow calling SetError with more than 1 concatenation or complex expression or non constant left-hand-side, eg. `d.SetError(500, "error on" + Bla(bar) + Yay(baz))`, you must repharase the error detail into something like this: `d.SetError(500, "error on " + msg)`

## TODOs

- Add SEO pre-render: [Rendora](//github.com/rendora/rendora)
- Add search-engine: [TypeSense](//typesense.org/) example
- Add persisted cache: [IceFireDB](//github.com/gitsrc/IceFireDB) or [Aerospike](//aerospike.com/) 
- Add external storage upload example (minio? wasabi?)
- Replace LightStep with [SigNoz](//github.com/SigNoz/signoz) and/or [datav](//github.com/datav-io/datav) 
- Add more deployment script with [LXC/LXD share](//bobcares.com/blog/how-to-setup-high-density-vps-hosting-using-lxc-linux-containers-and-lxd/) for single server multi-tenant
- Add backup scripts for Tarantool and Clickhouse

## File Upload Example

the schema (`/model/m[Something]/[something]_tables.go`), after creating this, run `make gen-orm`:

```go
const (
	TableMediaUploads Tt.TableName = `mediaUploads`
	Id            = `id`
	CreatedBy     = `createdBy`
	CreatedAt     = `createdAt`
	UpdatedBy     = `updatedBy`
	UpdatedAt     = `updatedAt`
	DeletedBy     = `deletedBy`
	DeletedAt     = `deletedAt`
	IsDeleted     = `isDeleted`
	RestoredBy    = `restoredBy`
	RestoredAt    = `restoredAt`
	SizeByte      = `sizeByte`
	FilePath      = `filePath`
	ContentType   = `contentType`
	OrigName      = `origName`
)

var TarantoolTables = map[Tt.TableName]*Tt.TableProp{
	TableMediaUploads: {
		Fields: []Tt.Field{
			{Id, Tt.Unsigned},
			{CreatedAt, Tt.Integer},
			{CreatedBy, Tt.Unsigned},
			{UpdatedAt, Tt.Integer},
			{UpdatedBy, Tt.Unsigned},
			{DeletedAt, Tt.Integer},
			{DeletedBy, Tt.Unsigned},
			{IsDeleted, Tt.Boolean},
			{RestoredAt, Tt.Integer},
			{RestoredBy, Tt.Unsigned},
			{SizeByte, Tt.Unsigned},
			{FilePath, Tt.String},
			{ContentType, Tt.String},
			{OrigName, Tt.String},
		},
		Unique1: Id,
		Unique2: FilePath,
	},
}

func GenerateORM() {
	Tt.GenerateOrm(TarantoolTables)
}

// don't forget to add migration on model.go:
//	m.Taran.MigrateTables(mSomething.TarantoolTables)
```

the code for domain/business logic `/domain/media.go`, after creating this, run `make gen-route`: 

```go
type (
	MediaUpload_In struct {
		RequestCommon
		UploadId   uint64 `json:"uploadId,string" form:"uploadId" query:"uploadId" long:"uploadId" msg:"uploadId"`
		FileBinary string `json:"fileBinary" form:"fileBinary" query:"fileBinary" long:"fileBinary" msg:"fileBinary"`
	}
	MediaUpload_Out struct {
		ResponseCommon
		MediaUpload *rqSomething.MediaUploads `json:"mediaUpload" form:"mediaUpload" query:"mediaUpload" long:"mediaUpload" msg:"mediaUpload"`
	}
)

const MediaUpload_Url = `/MediaUpload`

func (d *Domain) MediaUpload(in *MediaUpload_In) (out MediaUpload_Out) {
	sess := d.mustAdmin(in.SessionToken, in.UserAgent, &out.ResponseCommon)
	if sess == nil {
		return
	}

	if len(in.Uploads) == 0 {
		out.SetError(400, `missing fileBinary, enctype not multipart/form-data?`)
		return
	}

	up := wcSomething.NewMediaUploadsMutator(d.Taran)
	up.Id = in.UploadId
	if in.UploadId > 0 {
		if !up.FindById() {
			out.SetError(404, `upload not found, wrong env?`)
			return
		}
	}
	if up.CreatedAt == 0 {
		up.Id = id64.UID()
		up.CreatedAt = in.UnixNow()
		up.CreatedBy = sess.PlayerId
	}
	up.UpdatedAt = in.UnixNow()
	up.UpdatedBy = sess.PlayerId
	for fileName, tmpFile := range in.Uploads {
		up.OrigName = fileName
		oldPath := up.FilePath
		uriPath := conf.UPLOAD_URI + conf.MEDIA_SUBDIR
		dir := conf.UPLOAD_DIR + conf.MEDIA_SUBDIR
		idStr := I.UToS(up.Id)
		if S.StartsWith(oldPath, uriPath) {
			oldPath = dir + S.RightOf(oldPath, uriPath)
		} else if oldPath != `` {
			L.Print(`ERROR weird name format to be replaced for mediaUpload.id: ` + idStr + `:` + oldPath)
		}

		mtype, err := mimetype.DetectFile(tmpFile)
		if L.IsError(err, `cannot detect file type: `+tmpFile) {
			out.SetError(500, `cannot detect file type: `+up.OrigName)
			return
		}

		err = os.MkdirAll(dir, 0755)
		if L.IsError(err, `failed to create upload directory: `+dir) {
			out.SetError(500, `cannot create upload directory`)
			return
		}
		ext := S.ToLower(filepath.Ext(fileName))
		newName := idStr + ext
		newPath := dir + newName
		err = os.Rename(tmpFile, newPath)
		if L.IsError(err, `failed to rename `+tmpFile+` to `+newPath) {
			out.SetError(500, `failed moving uploaded file`)
			return
		}
		in.Uploads[fileName] = oldPath // delete old file later
		up.FilePath = uriPath + newName
		stat, err := os.Stat(newPath)
		if L.IsError(err, `failed to stat moved file: `+newPath) {
			out.SetError(500, `failed to stat moved file`)
			return
		}
		up.SizeByte = uint64(stat.Size())
		up.ContentType = mtype.String()

		// ignore if upload more than one
		break
	}
	//if in.DoDelete {
	//	up.IsDeleted = true
	//	up.DeletedAt = in.UnixNow()
	//	up.DeletedBy = sess.PlayerId
	//}
	//if in.DoRestore {
	//	up.IsDeleted = false
	//	up.RestoredAt = in.UnixNow()
	//	up.RestoredBy = sess.PlayerId
	//}

	if !up.DoReplace() {
		out.SetError(500, `cannot upsert media`)
		return
	}
	out.MediaUpload = &up.MediaUploads
	return
}
```
