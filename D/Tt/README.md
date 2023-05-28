
# Tarantool Adapter and ORM Generator

This package provides automatic schema migration (only for appending new columns, not for changing data types). This package also can be used to generate ORM.

## Dependencies

```
go install github.com/fatih/gomodifytags@latest
go install github.com/kokizzu/replacer@latest
```

## Generated ORM Example

![image](https://user-images.githubusercontent.com/1061610/131272783-0b2e0dc6-072b-47e1-854f-1f5291725255.png)
![image](https://user-images.githubusercontent.com/1061610/131272846-928a5630-4714-4092-8384-540c98772596.png)

## How to create a connection

```go
import "github.com/tarantool/go-tarantool"
import "github.com/kokizzu/gotro/L"

func ConnectTarantool() *tarantool.Connection {
	hostPort := fmt.Sprintf(`%s:%s`,
		TARANTOOL_HOST,
		TARANTOOL_PORT,
	)
	taran, err := tarantool.Connect(hostPort, tarantool.Opts{
		User: TARANTOOL_USER,
		Pass: TARANTOOL_PASS,
	})
	L.PanicIf(err, `tarantool.Connect `+hostPort)
	return taran
}

// then use it like this:
tt := &Tt.Adapter{Connection: ConnectTarantool(), Reconnect: ConnectTarantool}
```

## Usage

1. create a `model/` directory inside project
2. create a `m[Domain]` directory inside project, for example if the domain is authentication, you might want to create `mAuth`
3. create a `[domain]_tables.go` something like this:

```go
package mAuth

import "github.com/kokizzu/gotro/D/Tt"

const (
	TableUserss Tt.TableName = `users`
	Id                       = `id`
	Email                    = `email`
	Password                 = `password`
	CreatedBy                = `createdBy`
	CreatedAt                = `createdAt`
	UpdatedBy                = `updatedBy`
	UpdatedAt                = `updatedAt`
	DeletedBy                = `deletedBy`
	DeletedAt                = `deletedAt`
	IsDeleted                = `isDeleted`
	RestoredBy               = `restoredBy`
	RestoredAt               = `restoredAt`
	PasswordSetAt            = `passwordSetAt`
	SecretCode               = `secretCode`
	SecretCodeAt             = `secretCodeAt`
	VerificationSentAt       = `verificationSentAt`
	VerifiedAt               = `verifiedAt`
	LastLoginAt              = `lastLoginAt`
)

const (
	TableSessions Tt.TableName = `sessions`
	SessionToken               = `sessionToken`
	UserId                     = `userId`
	ExpiredAt                  = `expiredAt`
)

var TarantoolTables = map[Tt.TableName]*Tt.TableProp{
	// can only adding fields on back, and must IsNullable: true
	// primary key must be first field and set to Unique: fieldName
	TableUserss: {
		Fields: []Tt.Field{
			{Id, Tt.Unsigned},
			{Email, Tt.String},
			{Password, Tt.String},
			{CreatedAt, Tt.Integer},
			{CreatedBy, Tt.Unsigned},
			{UpdatedAt, Tt.Integer},
			{UpdatedBy, Tt.Unsigned},
			{DeletedAt, Tt.Integer},
			{DeletedBy, Tt.Unsigned},
			{IsDeleted, Tt.Boolean},
			{RestoredAt, Tt.Integer},
			{RestoredBy, Tt.Unsigned},
			{PasswordSetAt, Tt.Integer},
			{SecretCode, Tt.String},
			{SecretCodeAt, Tt.Integer},
			{VerificationSentAt, Tt.Integer},
			{VerifiedAt, Tt.Integer},
			{LastLoginAt, Tt.Integer},
		},
		AutoIncrementId: true,
		Unique2: Email,
		Indexes: []string{IsDeleted, SecretCode},
	},
	TableSessions: {
		Fields: []Tt.Field{
			{SessionToken, Tt.String},
			{UserId, Tt.Unsigned},
			{ExpiredAt, Tt.Integer},
		},
		Unique1: SessionToken,
	},
}

func GenerateORM() {
	Tt.GenerateOrm(TarantoolTables)
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


5. run the test to generate new ORM, that would generate `rq[Domain]/rq[Domain]__ORM.GEN.go` and `wc[Domain]/wc[Domain]__ORM.GEN.go` file, you might want to create a helper script for that:

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
	Taran     *Tt.Adapter
}

func NewDomain() *Domain {
	d := &Domain{
		Taran: &Tt.Adapter{conf.ConnectTarantool(), conf.ConnectTarantool},
	}
	return d
}
```


7. last step is just call generated method to manipulate or query, something like this:

```go

func (d *Domain) BusinessLogic1(in *BusinessLogic1_In) (out BusinessLogic1_Out) {
	
	// do something else
	
	user := wcAuth.NewUserMutator(d.Taran)
	user.Email = in.Email
	if !user.FindById() {
		user.Id = id64.UID()
		user.CreatedAt = in.UnixNow()
		if !user.DoInsert() {
			out.SetError(500, `failed to insert user record, db down?`)
			return
		}
	}
	user.SetUpdatedAt(in.UnixNow())
	// do other manipulation
	// use .Set* if you have to call DoUpdateBy*()
	if !user.DoUpdateById() {
		out.SetError(500, `failed to insert user record, db down?`)
		return		
	}
	
}

// or if you only need to read
func (d *Domain) mustLogin(token string, userAgent string, out *ResponseCommon) *conf.Session {
	sess := &conf.Session{}
	if token == `` {
		out.SetError(400, `missing session token`)
		return nil
	}
	if !sess.Decrypt(token, userAgent) {
		out.SetError(400, `invalid session token`) // if got this, possibly wrong userAgent-sessionToken pair
		return nil
	}
	if sess.ExpiredAt <= fastime.UnixNow() {
		out.SetError(400, `token expired`)
		return nil
	}

	session := rqAuth.NewSessions(d.Taran)
	session.SessionToken = token
	if !session.FindBySessionToken() {
		out.SetError(400, `session missing from database, wrong env?`)
		return nil
	}
	if session.ExpiredAt <= fastime.UnixNow() {
		out.SetError(403, `session expired or logged out`)
		return nil
	}
	return sess
}
```

8. If you need to create an extension method for the ORM, just add a new file on `rq[Domain]/anything.go`, with a new struct method from generated ORM, something like this:
```go
package rqAuth

import (
	"myProject/conf"
	
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/X"
	"golang.org/x/crypto/bcrypt"
)

func (s *Users) FindOffsetLimit(offset, limit uint32) (res []*Users) {
	query := `
SELECT ` + s.SqlSelectAllFields() + `
FROM ` + s.SqlTableName() + `
ORDER BY ` + s.SqlId() + `
LIMIT ` + X.ToS(limit) + `
OFFSET ` + X.ToS(offset) // note: for string, use S.Z or S.XSS to prevent SQL injection
	if conf.DEBUG_MODE {
		L.Print(query)
	}
	s.Adapter.QuerySql(query, func(row []any) {
		obj := &Users{}
		obj.FromArray(row)
		obj.CensorFields()
		res = append(res, obj)
	})
	return
}

func (s *Users) CheckPassword(currentPassword string) bool {
	hash := []byte(s.Password)
	pass := []byte(currentPassword)
	err := bcrypt.CompareHashAndPassword(hash, pass)

	return !L.IsError(err, `bcrypt.CompareHashAndPassword`)
}

// call before outputting to client
func (s *Users) CensorFields() {
	s.Password = ``
	s.SecretCode = ``
}
```

or in `wc[Domain]/anything.go` if you need to manipulate things

```go

func (p *UsersMutator) SetEncryptPassword(password string) bool {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	p.SetPassword(string(pass))
	return !L.IsError(err, `bcrypt.GenerateFromPassword`)
}
```


9. to initialize automatic migration, just create `model/run_migration.go`

```go
func RunMigration() {
	L.Print(`run migration..`)
	tt := &Tt.Adapter{Connection: ConnectTarantool(), Reconnect: ConnectTarantool}
	tt.MigrateTables(mAuth.ClickhouseTables)
	// add other tarantool tables to be migrated here
}
```

then call it on `main`

```go
func main() {
	model.RunMigration()
}
```
