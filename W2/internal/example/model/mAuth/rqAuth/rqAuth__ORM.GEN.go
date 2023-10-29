package rqAuth

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	"github.com/kokizzu/gotro/W2/internal/example/model/mAuth"

	"github.com/tarantool/go-tarantool"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/X"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file rqAuth__ORM.GEN.go
//go:generate replacer -afterprefix 'Id" form' 'Id,string" form' type rqAuth__ORM.GEN.go
//go:generate replacer -afterprefix 'json:"id"' 'json:"id,string"' type rqAuth__ORM.GEN.go
//go:generate replacer -afterprefix 'By" form' 'By,string" form' type rqAuth__ORM.GEN.go
// go:generate msgp -tests=false -file rqAuth__ORM.GEN.go -o rqAuth__MSG.GEN.go

// Sessions DAO reader/query struct
type Sessions struct {
	Adapter      *Tt.Adapter `json:"-" msg:"-" query:"-" form:"-"`
	SessionToken string
	UserId       uint64
	ExpiredAt    int64
}

// NewSessions create new ORM reader/query object
func NewSessions(adapter *Tt.Adapter) *Sessions {
	return &Sessions{Adapter: adapter}
}

// SpaceName returns full package and table name
func (s *Sessions) SpaceName() string { //nolint:dupl false positive
	return string(mAuth.TableSessions) // casting required to string from Tt.TableName
}

// SqlTableName returns quoted table name
func (s *Sessions) SqlTableName() string { //nolint:dupl false positive
	return `"sessions"`
}

// UniqueIndexSessionToken return unique index name
func (s *Sessions) UniqueIndexSessionToken() string { //nolint:dupl false positive
	return `sessionToken`
}

// FindBySessionToken Find one by SessionToken
func (s *Sessions) FindBySessionToken() bool { //nolint:dupl false positive
	res, err := s.Adapter.Select(s.SpaceName(), s.UniqueIndexSessionToken(), 0, 1, tarantool.IterEq, A.X{s.SessionToken})
	if L.IsError(err, `Sessions.FindBySessionToken failed: `+s.SpaceName()) {
		return false
	}
	rows := res.Tuples()
	if len(rows) == 1 {
		s.FromArray(rows[0])
		return true
	}
	return false
}

// SqlSelectAllFields generate Sql select fields
func (s *Sessions) SqlSelectAllFields() string { //nolint:dupl false positive
	return ` "sessionToken"
	, "userId"
	, "expiredAt"
	`
}

// ToUpdateArray generate slice of update command
func (s *Sessions) ToUpdateArray() A.X { //nolint:dupl false positive
	return A.X{
		A.X{`=`, 0, s.SessionToken},
		A.X{`=`, 1, s.UserId},
		A.X{`=`, 2, s.ExpiredAt},
	}
}

// IdxSessionToken return name of the index
func (s *Sessions) IdxSessionToken() int { //nolint:dupl false positive
	return 0
}

// SqlSessionToken return name of the column being indexed
func (s *Sessions) SqlSessionToken() string { //nolint:dupl false positive
	return `"sessionToken"`
}

// IdxUserId return name of the index
func (s *Sessions) IdxUserId() int { //nolint:dupl false positive
	return 1
}

// SqlUserId return name of the column being indexed
func (s *Sessions) SqlUserId() string { //nolint:dupl false positive
	return `"userId"`
}

// IdxExpiredAt return name of the index
func (s *Sessions) IdxExpiredAt() int { //nolint:dupl false positive
	return 2
}

// SqlExpiredAt return name of the column being indexed
func (s *Sessions) SqlExpiredAt() string { //nolint:dupl false positive
	return `"expiredAt"`
}

// ToArray receiver fields to slice
func (s *Sessions) ToArray() A.X { //nolint:dupl false positive
	return A.X{
		s.SessionToken, // 0
		s.UserId,       // 1
		s.ExpiredAt,    // 2
	}
}

// FromArray convert slice to receiver fields
func (s *Sessions) FromArray(a A.X) *Sessions { //nolint:dupl false positive
	s.SessionToken = X.ToS(a[0])
	s.UserId = X.ToU(a[1])
	s.ExpiredAt = X.ToI(a[2])
	return s
}

// FindOffsetLimit returns slice of struct, order by idx, eg. .UniqueIndex*()
func (s *Sessions) FindOffsetLimit(offset, limit uint32, idx string) []Sessions { //nolint:dupl false positive
	var rows []Sessions
	res, err := s.Adapter.Select(s.SpaceName(), idx, offset, limit, tarantool.IterAll, A.X{})
	if L.IsError(err, `Sessions.FindOffsetLimit failed: `+s.SpaceName()) {
		return rows
	}
	for _, row := range res.Tuples() {
		item := Sessions{}
		rows = append(rows, *item.FromArray(row))
	}
	return rows
}

// FindArrOffsetLimit returns as slice of slice order by idx eg. .UniqueIndex*()
func (s *Sessions) FindArrOffsetLimit(offset, limit uint32, idx string) ([]A.X, Tt.QueryMeta) { //nolint:dupl false positive
	var rows []A.X
	res, err := s.Adapter.Select(s.SpaceName(), idx, offset, limit, tarantool.IterAll, A.X{})
	if L.IsError(err, `Sessions.FindOffsetLimit failed: `+s.SpaceName()) {
		return rows, Tt.QueryMetaFrom(res, err)
	}
	tuples := res.Tuples()
	rows = make([]A.X, len(tuples))
	for z, row := range tuples {
		rows[z] = row
	}
	return rows, Tt.QueryMetaFrom(res, nil)
}

// Total count number of rows
func (s *Sessions) Total() int64 { //nolint:dupl false positive
	rows := s.Adapter.CallBoxSpace(s.SpaceName()+`:count`, A.X{})
	if len(rows) > 0 && len(rows[0]) > 0 {
		return X.ToI(rows[0][0])
	}
	return 0
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

// Users DAO reader/query struct
type Users struct {
	Adapter            *Tt.Adapter `json:"-" msg:"-" query:"-" form:"-"`
	Id                 uint64
	Email              string
	Password           string
	CreatedAt          int64
	CreatedBy          uint64
	UpdatedAt          int64
	UpdatedBy          uint64
	DeletedAt          int64
	DeletedBy          uint64
	IsDeleted          bool
	RestoredAt         int64
	RestoredBy         uint64
	PasswordSetAt      int64
	SecretCode         string
	SecretCodeAt       int64
	VerificationSentAt int64
	VerifiedAt         int64
	LastLoginAt        int64
}

// NewUsers create new ORM reader/query object
func NewUsers(adapter *Tt.Adapter) *Users {
	return &Users{Adapter: adapter}
}

// SpaceName returns full package and table name
func (u *Users) SpaceName() string { //nolint:dupl false positive
	return string(mAuth.TableUsers) // casting required to string from Tt.TableName
}

// SqlTableName returns quoted table name
func (u *Users) SqlTableName() string { //nolint:dupl false positive
	return `"users"`
}

func (u *Users) UniqueIndexId() string { //nolint:dupl false positive
	return `id`
}

// FindById Find one by Id
func (u *Users) FindById() bool { //nolint:dupl false positive
	res, err := u.Adapter.Select(u.SpaceName(), u.UniqueIndexId(), 0, 1, tarantool.IterEq, A.X{u.Id})
	if L.IsError(err, `Users.FindById failed: `+u.SpaceName()) {
		return false
	}
	rows := res.Tuples()
	if len(rows) == 1 {
		u.FromArray(rows[0])
		return true
	}
	return false
}

// UniqueIndexEmail return unique index name
func (u *Users) UniqueIndexEmail() string { //nolint:dupl false positive
	return `email`
}

// FindByEmail Find one by Email
func (u *Users) FindByEmail() bool { //nolint:dupl false positive
	res, err := u.Adapter.Select(u.SpaceName(), u.UniqueIndexEmail(), 0, 1, tarantool.IterEq, A.X{u.Email})
	if L.IsError(err, `Users.FindByEmail failed: `+u.SpaceName()) {
		return false
	}
	rows := res.Tuples()
	if len(rows) == 1 {
		u.FromArray(rows[0])
		return true
	}
	return false
}

// SqlSelectAllFields generate Sql select fields
func (u *Users) SqlSelectAllFields() string { //nolint:dupl false positive
	return ` "id"
	, "email"
	, "password"
	, "createdAt"
	, "createdBy"
	, "updatedAt"
	, "updatedBy"
	, "deletedAt"
	, "deletedBy"
	, "isDeleted"
	, "restoredAt"
	, "restoredBy"
	, "passwordSetAt"
	, "secretCode"
	, "secretCodeAt"
	, "verificationSentAt"
	, "verifiedAt"
	, "lastLoginAt"
	`
}

// ToUpdateArray generate slice of update command
func (u *Users) ToUpdateArray() A.X { //nolint:dupl false positive
	return A.X{
		A.X{`=`, 0, u.Id},
		A.X{`=`, 1, u.Email},
		A.X{`=`, 2, u.Password},
		A.X{`=`, 3, u.CreatedAt},
		A.X{`=`, 4, u.CreatedBy},
		A.X{`=`, 5, u.UpdatedAt},
		A.X{`=`, 6, u.UpdatedBy},
		A.X{`=`, 7, u.DeletedAt},
		A.X{`=`, 8, u.DeletedBy},
		A.X{`=`, 9, u.IsDeleted},
		A.X{`=`, 10, u.RestoredAt},
		A.X{`=`, 11, u.RestoredBy},
		A.X{`=`, 12, u.PasswordSetAt},
		A.X{`=`, 13, u.SecretCode},
		A.X{`=`, 14, u.SecretCodeAt},
		A.X{`=`, 15, u.VerificationSentAt},
		A.X{`=`, 16, u.VerifiedAt},
		A.X{`=`, 17, u.LastLoginAt},
	}
}

// IdxId return name of the index
func (u *Users) IdxId() int { //nolint:dupl false positive
	return 0
}

// SqlId return name of the column being indexed
func (u *Users) SqlId() string { //nolint:dupl false positive
	return `"id"`
}

// IdxEmail return name of the index
func (u *Users) IdxEmail() int { //nolint:dupl false positive
	return 1
}

// SqlEmail return name of the column being indexed
func (u *Users) SqlEmail() string { //nolint:dupl false positive
	return `"email"`
}

// IdxPassword return name of the index
func (u *Users) IdxPassword() int { //nolint:dupl false positive
	return 2
}

// SqlPassword return name of the column being indexed
func (u *Users) SqlPassword() string { //nolint:dupl false positive
	return `"password"`
}

// IdxCreatedAt return name of the index
func (u *Users) IdxCreatedAt() int { //nolint:dupl false positive
	return 3
}

// SqlCreatedAt return name of the column being indexed
func (u *Users) SqlCreatedAt() string { //nolint:dupl false positive
	return `"createdAt"`
}

// IdxCreatedBy return name of the index
func (u *Users) IdxCreatedBy() int { //nolint:dupl false positive
	return 4
}

// SqlCreatedBy return name of the column being indexed
func (u *Users) SqlCreatedBy() string { //nolint:dupl false positive
	return `"createdBy"`
}

// IdxUpdatedAt return name of the index
func (u *Users) IdxUpdatedAt() int { //nolint:dupl false positive
	return 5
}

// SqlUpdatedAt return name of the column being indexed
func (u *Users) SqlUpdatedAt() string { //nolint:dupl false positive
	return `"updatedAt"`
}

// IdxUpdatedBy return name of the index
func (u *Users) IdxUpdatedBy() int { //nolint:dupl false positive
	return 6
}

// SqlUpdatedBy return name of the column being indexed
func (u *Users) SqlUpdatedBy() string { //nolint:dupl false positive
	return `"updatedBy"`
}

// IdxDeletedAt return name of the index
func (u *Users) IdxDeletedAt() int { //nolint:dupl false positive
	return 7
}

// SqlDeletedAt return name of the column being indexed
func (u *Users) SqlDeletedAt() string { //nolint:dupl false positive
	return `"deletedAt"`
}

// IdxDeletedBy return name of the index
func (u *Users) IdxDeletedBy() int { //nolint:dupl false positive
	return 8
}

// SqlDeletedBy return name of the column being indexed
func (u *Users) SqlDeletedBy() string { //nolint:dupl false positive
	return `"deletedBy"`
}

// IdxIsDeleted return name of the index
func (u *Users) IdxIsDeleted() int { //nolint:dupl false positive
	return 9
}

// SqlIsDeleted return name of the column being indexed
func (u *Users) SqlIsDeleted() string { //nolint:dupl false positive
	return `"isDeleted"`
}

// IdxRestoredAt return name of the index
func (u *Users) IdxRestoredAt() int { //nolint:dupl false positive
	return 10
}

// SqlRestoredAt return name of the column being indexed
func (u *Users) SqlRestoredAt() string { //nolint:dupl false positive
	return `"restoredAt"`
}

// IdxRestoredBy return name of the index
func (u *Users) IdxRestoredBy() int { //nolint:dupl false positive
	return 11
}

// SqlRestoredBy return name of the column being indexed
func (u *Users) SqlRestoredBy() string { //nolint:dupl false positive
	return `"restoredBy"`
}

// IdxPasswordSetAt return name of the index
func (u *Users) IdxPasswordSetAt() int { //nolint:dupl false positive
	return 12
}

// SqlPasswordSetAt return name of the column being indexed
func (u *Users) SqlPasswordSetAt() string { //nolint:dupl false positive
	return `"passwordSetAt"`
}

// IdxSecretCode return name of the index
func (u *Users) IdxSecretCode() int { //nolint:dupl false positive
	return 13
}

// SqlSecretCode return name of the column being indexed
func (u *Users) SqlSecretCode() string { //nolint:dupl false positive
	return `"secretCode"`
}

// IdxSecretCodeAt return name of the index
func (u *Users) IdxSecretCodeAt() int { //nolint:dupl false positive
	return 14
}

// SqlSecretCodeAt return name of the column being indexed
func (u *Users) SqlSecretCodeAt() string { //nolint:dupl false positive
	return `"secretCodeAt"`
}

// IdxVerificationSentAt return name of the index
func (u *Users) IdxVerificationSentAt() int { //nolint:dupl false positive
	return 15
}

// SqlVerificationSentAt return name of the column being indexed
func (u *Users) SqlVerificationSentAt() string { //nolint:dupl false positive
	return `"verificationSentAt"`
}

// IdxVerifiedAt return name of the index
func (u *Users) IdxVerifiedAt() int { //nolint:dupl false positive
	return 16
}

// SqlVerifiedAt return name of the column being indexed
func (u *Users) SqlVerifiedAt() string { //nolint:dupl false positive
	return `"verifiedAt"`
}

// IdxLastLoginAt return name of the index
func (u *Users) IdxLastLoginAt() int { //nolint:dupl false positive
	return 17
}

// SqlLastLoginAt return name of the column being indexed
func (u *Users) SqlLastLoginAt() string { //nolint:dupl false positive
	return `"lastLoginAt"`
}

// ToArray receiver fields to slice
func (u *Users) ToArray() A.X { //nolint:dupl false positive
	var id any = nil
	if u.Id != 0 {
		id = u.Id
	}
	return A.X{
		id,
		u.Email,              // 1
		u.Password,           // 2
		u.CreatedAt,          // 3
		u.CreatedBy,          // 4
		u.UpdatedAt,          // 5
		u.UpdatedBy,          // 6
		u.DeletedAt,          // 7
		u.DeletedBy,          // 8
		u.IsDeleted,          // 9
		u.RestoredAt,         // 10
		u.RestoredBy,         // 11
		u.PasswordSetAt,      // 12
		u.SecretCode,         // 13
		u.SecretCodeAt,       // 14
		u.VerificationSentAt, // 15
		u.VerifiedAt,         // 16
		u.LastLoginAt,        // 17
	}
}

// FromArray convert slice to receiver fields
func (u *Users) FromArray(a A.X) *Users { //nolint:dupl false positive
	u.Id = X.ToU(a[0])
	u.Email = X.ToS(a[1])
	u.Password = X.ToS(a[2])
	u.CreatedAt = X.ToI(a[3])
	u.CreatedBy = X.ToU(a[4])
	u.UpdatedAt = X.ToI(a[5])
	u.UpdatedBy = X.ToU(a[6])
	u.DeletedAt = X.ToI(a[7])
	u.DeletedBy = X.ToU(a[8])
	u.IsDeleted = X.ToBool(a[9])
	u.RestoredAt = X.ToI(a[10])
	u.RestoredBy = X.ToU(a[11])
	u.PasswordSetAt = X.ToI(a[12])
	u.SecretCode = X.ToS(a[13])
	u.SecretCodeAt = X.ToI(a[14])
	u.VerificationSentAt = X.ToI(a[15])
	u.VerifiedAt = X.ToI(a[16])
	u.LastLoginAt = X.ToI(a[17])
	return u
}

// FindOffsetLimit returns slice of struct, order by idx, eg. .UniqueIndex*()
func (u *Users) FindOffsetLimit(offset, limit uint32, idx string) []Users { //nolint:dupl false positive
	var rows []Users
	res, err := u.Adapter.Select(u.SpaceName(), idx, offset, limit, tarantool.IterAll, A.X{})
	if L.IsError(err, `Users.FindOffsetLimit failed: `+u.SpaceName()) {
		return rows
	}
	for _, row := range res.Tuples() {
		item := Users{}
		rows = append(rows, *item.FromArray(row))
	}
	return rows
}

// FindArrOffsetLimit returns as slice of slice order by idx eg. .UniqueIndex*()
func (u *Users) FindArrOffsetLimit(offset, limit uint32, idx string) ([]A.X, Tt.QueryMeta) { //nolint:dupl false positive
	var rows []A.X
	res, err := u.Adapter.Select(u.SpaceName(), idx, offset, limit, tarantool.IterAll, A.X{})
	if L.IsError(err, `Users.FindOffsetLimit failed: `+u.SpaceName()) {
		return rows, Tt.QueryMetaFrom(res, err)
	}
	tuples := res.Tuples()
	rows = make([]A.X, len(tuples))
	for z, row := range tuples {
		rows[z] = row
	}
	return rows, Tt.QueryMetaFrom(res, nil)
}

// Total count number of rows
func (u *Users) Total() int64 { //nolint:dupl false positive
	rows := u.Adapter.CallBoxSpace(u.SpaceName()+`:count`, A.X{})
	if len(rows) > 0 && len(rows[0]) > 0 {
		return X.ToI(rows[0][0])
	}
	return 0
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go
