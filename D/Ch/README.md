
# Clickhouse Adapter and ORM Generator

This package provides automatic migration (only when adding more columns on the last position, not for changing reordering or changing order key's data type). This package also can be used to generate ORM           

## Generated ORM example

![image](https://user-images.githubusercontent.com/1061610/131272641-2fe22b60-0f8a-47ad-b8f9-49e571946648.png)


## How to create a connection

```go
import (
	"database/sql"
	"fmt"

	_ "github.com/ClickHouse/clickhouse-go"
)

func ConnectClickhouse() *sql.DB {
	connStr := fmt.Sprintf("tcp://%s:%s?debug=true",
		CLICKHOUSE_HOST,
		CLICKHOUSE_PORT,
	)
	click, err := sql.Open(`clickhouse`, connStr)
	L.PanicIf(err, `sql.Open `+connStr)
	return click
}

// then use it like this:
ch := &Ch.Adapter{DB: conf.ConnectClickhouse(), Reconnect: conf.ConnectClickhouse}
```

## Usage

1. create a `model/` directory inside your project
2. create a `model/m[Domain]` directory, for example if the business domain is authentication, you might want to create `mAuth`
3. create a `[domain]_tables.go` something like this:

```go
package mAuth

import "github.com/kokizzu/gotro/D/Ch"

// table userlogs
const (
	TableUserLogs Ch.TableName = `userLogs`
	CreatedAt                  = `createdAt`
	RequestId                  = `requestId`
	Error                      = `error`
	ActorId                    = `actorId`
	IpAddr4                    = `ipAddr4`
	IpAddr6                    = `ipAddr6`
	UserAgent                  = `userAgent`
)

var ClickhouseTables = map[Ch.TableName]*Ch.TableProp{
	TableUserLogs: {
		Fields: []Ch.Field{
			{CreatedAt, Ch.DateTime},
			{RequestId, Ch.UInt64},
			{ActorId, Ch.UInt64},
			{Error, Ch.String},
			{IpAddr4, Ch.IPv4},
			{IpAddr6, Ch.IPv6},
			{UserAgent, Ch.String},
			// add more columns here
		},
		Orders: []string{ActorId, RequestId},
	},
}

func GenerateORM() {
	Ch.GenerateOrm(ClickhouseTables) 
}
```

4. create a `[domain]_generator_test.go` something like this:
```go
package mAuth

import (
	"testing"
)

//go:generate go test -run=XXX -bench=Benchmark_GenerateOrm

func Benchmark_GenerateOrm(b *testing.B) {
	GenerateORM()
	b.SkipNow()
}
```

5. run the test to generate new ORM, that would generate `sa[Domain]/sa[Domain]__ORM.GEN.go` file, you might want to create a helper script for that:

```bash
#!/usr/bin/env bash

cd ./model
  cat *.go | grep '//go:generate ' | cut -d ' ' -f 2- | bash -x > /tmp/1.log
  
for i in ./m*; do
  if [[ ! -d "$i" ]] ; then continue ; fi
  echo $i
  pushd .
  cd "$i"
  
  # generate ORM
  go test -bench=.
  
  for j in ./*; do 
    echo $j
    if [[ ! -d "$j" ]] ; then continue ; fi
        
    pushd .
    cd "$j" 
    echo `pwd` 
    cat *.go | grep '//go:generate ' | cut -d ' ' -f 2- | bash -x >> /tmp/1.log    
    popd 
    
  done
  
  popd
  
done
```

6. in your web server engine/domain logic (one that initializes dependencies), create methods to help initialize the buffer, something like this:

```go
type Domain struct {
	Click     *Ch.Adapter
	chBuffers map[Ch.TableName]*chBuffer.TimedBuffer
	waitGroup *sync.WaitGroup
	// add more dependency initialization here
}

func (d *Domain) InitClickhouseBuffer(preparators map[Ch.TableName]chBuffer.Preparator) {
	for tableName, preparator := range preparators {
		chb := chBuffer.NewTimedBuffer(d.Click.DB, 30000, 1*time.Second, preparator)
		chb.IgnoreInterrupt = true
		d.chBuffers[tableName] = chb
		d.waitGroup.Add(1)
	}
}

func (d *Domain) WaitInterrupt() {
	interrupt := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	//signal.Notify(interrupt, os.Interrupt, syscall.SIGKILL)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGHUP)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGQUIT)

	<-interrupt
	L.Print(`caught signal`, interrupt)
}

func (d *Domain) handleTermSignal() {
	d.WaitInterrupt()
	for tableName := range d.chBuffers {
		go func(tableName Ch.TableName) {
			chb := d.chBuffers[tableName]
			chb.Close()
			<-chb.WaitFinalFlush
			L.Print(`done waiting: ` + tableName)
			d.waitGroup.Done()
		}(tableName)
	}
	d.waitGroup.Wait()
	os.Exit(0)
}

func (d *Domain) Statistics(row AnalyticsRow) {
	tableName := row.TableName()
	res := d.chBuffers[tableName]
	if res != nil {
		res.Insert(row.SqlInsertParam())
		return
	}
	panic(`did you forgot to register InitClickhouseBuffer preparators for ` + string(tableName))
}


func NewDomain() *Domain {
	d := &Domain{
		Click: &Ch.Adapter{conf.ConnectClickhouse(), conf.ConnectClickhouse},
	}
	d.waitGroup = &sync.WaitGroup{}
	d.chBuffers = map[Ch.TableName]*chBuffer.TimedBuffer{}
	d.InitClickhouseBuffer(saAuth.Preparators)

	go d.handleTermSignal()
	// add more preparators if there's new clickhouse tables on model
	return d
}
```

7. last step is just call `Domain.Statistics` method to insert a new log, something like this:

```go

func (d *Domain) BusinessLogic1(in *BusinessLogic1_In) (out BusinessLogic1_Out) {
	
	// do something else
	
	d.Statistics(saAuth.UserLogs{
		CreatedAt: in.Now(),
		RequestId: ctx.RequestId,
		Error ctx.Error,
		ActorId: session.UserId,
		IpAddr4: ctx.RemoteAddr4,
		IpAddr6: ctx.RemoteAddr6,
		UserAgent: session.UserAgent,
	})
	
}
```

8. If you need to create an extension method for the ORM, just add a new file on `sa[Domain]/anything.go`, with a new struct method from generated ORM, something like this:
```go
package saAuth

import (
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/X"
)

type TopUser struct {
	UserId uint64
	Count int64
	Rank int64
}

func (m *UserLogs) FindTop1k(daySpan int, offset int64) (res []TopUser) {
	query := `
SELECT ` + m.sqlActorId() + `
	, COUNT(1) 
FROM ` + m.sqlTableName() + ` 
WHERE ` + m.sqlCreatedAt() + ` >= subtractDays(now(),` + I.ToStr(daySpan) + `) 
GROUP BY ` + m.sqlActorId() + ` 
ORDER BY COUNT(1) DESC
	,  MAX(` + m.sqlCreatedAt() + `)
LIMIT 1000
OFFSET ` + X.ToS(offset) + `
` // note: for string, use S.Z or S.XSS to prevent SQL injection
	rows, err := m.Adapter.Query(query)
	if L.IsError(err, `failed query: `+query) {
		return
	}
	defer rows.Close()
	rankNo := int64(1)
	for rows.Next() {
		row := TopUser{Rank: rankNo}
		err := rows.Scan(&row.UserId, &row.Count)
		if L.IsError(err, `rows.Scan`) {
			return
		}
		rankNo++
		res = append(res, row)
	}
	return
}

```

9. to initialize automatic migration, just create `model/run_migration.go`

```
func RunMigration() {
	L.Print(`run migration..`)
	ch := &Ch.Adapter{DB: ConnectClickhouse(), Reconnect: ConnectClickhouse}
	ch.MigrateTables(mAuth.ClickhouseTables)
	// add other clickhouse tables to be migrated here
}
```

then call it on `main`

```
func main() {
	model.RunMigration()
}
```
