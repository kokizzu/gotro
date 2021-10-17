package rqAuth

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	"github.com/kokizzu/gotro/W2/example/model/mAuth"

	"github.com/tarantool/go-tarantool"

	"github.com/graphql-go/graphql"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/X"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file rqAuth__ORM.GEN.go
//go:generate replacer 'Id" form' 'Id,string" form' type rqAuth__ORM.GEN.go
//go:generate replacer 'json:"id"' 'json:"id,string"' type rqAuth__ORM.GEN.go
//go:generate replacer 'By" form' 'By,string" form' type rqAuth__ORM.GEN.go
// go:generate msgp -tests=false -file rqAuth__ORM.GEN.go -o rqAuth__MSG.GEN.go

type Sessions struct {
	Adapter      *Tt.Adapter `json:"-" msg:"-" query:"-" form:"-" long:"adapter"`
	SessionToken string      `json:"sessionToken" form:"sessionToken" query:"sessionToken" long:"sessionToken" msg:"sessionToken"`
	UserId       uint64      `json:"userId,string" form:"userId" query:"userId" long:"userId" msg:"userId"`
	ExpiredAt    int64       `json:"expiredAt" form:"expiredAt" query:"expiredAt" long:"expiredAt" msg:"expiredAt"`
}

func NewSessions(adapter *Tt.Adapter) *Sessions {
	return &Sessions{Adapter: adapter}
}

func (s *Sessions) SpaceName() string { //nolint:dupl false positive
	return string(mAuth.TableSessions)
}

func (s *Sessions) sqlTableName() string { //nolint:dupl false positive
	return `"sessions"`
}

func (s *Sessions) UniqueIndexSessionToken() string { //nolint:dupl false positive
	return `sessionToken`
}

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

var GraphqlFieldSessionsBySessionToken = &graphql.Field{
	Type:        GraphqlTypeSessions,
	Description: `list of Sessions`,
	Args: graphql.FieldConfigArgument{
		`SessionToken`: &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
}

func (g *Sessions) GraphqlFieldSessionsBySessionTokenWithResolver() *graphql.Field {
	field := *GraphqlFieldSessionsBySessionToken
	field.Resolve = func(p graphql.ResolveParams) (interface{}, error) {
		q := g
		v, ok := p.Args[`sessionToken`]
		if !ok {
			v, _ = p.Args[`SessionToken`]
		}
		q.SessionToken = X.ToS(v)
		if q.FindBySessionToken() {
			return q, nil
		}
		return nil, nil
	}
	return &field
}

func (s *Sessions) sqlSelectAllFields() string { //nolint:dupl false positive
	return ` "sessionToken"
	, "userId"
	, "expiredAt"
	`
}

func (s *Sessions) ToUpdateArray() A.X { //nolint:dupl false positive
	return A.X{
		A.X{`=`, 0, s.SessionToken},
		A.X{`=`, 1, s.UserId},
		A.X{`=`, 2, s.ExpiredAt},
	}
}

func (s *Sessions) IdxSessionToken() int { //nolint:dupl false positive
	return 0
}

func (s *Sessions) sqlSessionToken() string { //nolint:dupl false positive
	return `"sessionToken"`
}

func (s *Sessions) IdxUserId() int { //nolint:dupl false positive
	return 1
}

func (s *Sessions) sqlUserId() string { //nolint:dupl false positive
	return `"userId"`
}

func (s *Sessions) IdxExpiredAt() int { //nolint:dupl false positive
	return 2
}

func (s *Sessions) sqlExpiredAt() string { //nolint:dupl false positive
	return `"expiredAt"`
}

func (s *Sessions) ToArray() A.X { //nolint:dupl false positive
	return A.X{
		s.SessionToken, // 0
		s.UserId,       // 1
		s.ExpiredAt,    // 2
	}
}

func (s *Sessions) FromArray(a A.X) *Sessions { //nolint:dupl false positive
	s.SessionToken = X.ToS(a[0])
	s.UserId = X.ToU(a[1])
	s.ExpiredAt = X.ToI(a[2])
	return s
}

func (s *Sessions) Total() int64 { //nolint:dupl false positive
	rows := s.Adapter.CallBoxSpace(s.SpaceName()+`:count`, A.X{})
	if len(rows) > 0 && len(rows[0]) > 0 {
		return X.ToI(rows[0][0])
	}
	return 0
}

var GraphqlTypeSessions = graphql.NewObject(
	graphql.ObjectConfig{
		Name: `sessions`,
		Fields: graphql.Fields{
			`sessionToken`: &graphql.Field{
				Type: graphql.String,
			},
			`userId`: &graphql.Field{
				Type: graphql.ID,
			},
			`expiredAt`: &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

type Users struct {
	Adapter            *Tt.Adapter `json:"-" msg:"-" query:"-" form:"-" long:"adapter"`
	Id                 uint64      `json:"id,string" form:"id" query:"id" long:"id" msg:"id"`
	Email              string      `json:"email" form:"email" query:"email" long:"email" msg:"email"`
	Password           string      `json:"password" form:"password" query:"password" long:"password" msg:"password"`
	CreatedAt          int64       `json:"createdAt" form:"createdAt" query:"createdAt" long:"createdAt" msg:"createdAt"`
	CreatedBy          uint64      `json:"createdBy,string" form:"createdBy" query:"createdBy" long:"createdBy" msg:"createdBy"`
	UpdatedAt          int64       `json:"updatedAt" form:"updatedAt" query:"updatedAt" long:"updatedAt" msg:"updatedAt"`
	UpdatedBy          uint64      `json:"updatedBy,string" form:"updatedBy" query:"updatedBy" long:"updatedBy" msg:"updatedBy"`
	DeletedAt          int64       `json:"deletedAt" form:"deletedAt" query:"deletedAt" long:"deletedAt" msg:"deletedAt"`
	DeletedBy          uint64      `json:"deletedBy,string" form:"deletedBy" query:"deletedBy" long:"deletedBy" msg:"deletedBy"`
	IsDeleted          bool        `json:"isDeleted" form:"isDeleted" query:"isDeleted" long:"isDeleted" msg:"isDeleted"`
	RestoredAt         int64       `json:"restoredAt" form:"restoredAt" query:"restoredAt" long:"restoredAt" msg:"restoredAt"`
	RestoredBy         uint64      `json:"restoredBy,string" form:"restoredBy" query:"restoredBy" long:"restoredBy" msg:"restoredBy"`
	PasswordSetAt      int64       `json:"passwordSetAt" form:"passwordSetAt" query:"passwordSetAt" long:"passwordSetAt" msg:"passwordSetAt"`
	SecretCode         string      `json:"secretCode" form:"secretCode" query:"secretCode" long:"secretCode" msg:"secretCode"`
	SecretCodeAt       int64       `json:"secretCodeAt" form:"secretCodeAt" query:"secretCodeAt" long:"secretCodeAt" msg:"secretCodeAt"`
	VerificationSentAt int64       `json:"verificationSentAt" form:"verificationSentAt" query:"verificationSentAt" long:"verificationSentAt" msg:"verificationSentAt"`
	VerifiedAt         int64       `json:"verifiedAt" form:"verifiedAt" query:"verifiedAt" long:"verifiedAt" msg:"verifiedAt"`
	LastLoginAt        int64       `json:"lastLoginAt" form:"lastLoginAt" query:"lastLoginAt" long:"lastLoginAt" msg:"lastLoginAt"`
}

func NewUsers(adapter *Tt.Adapter) *Users {
	return &Users{Adapter: adapter}
}

func (u *Users) SpaceName() string { //nolint:dupl false positive
	return string(mAuth.TableUsers)
}

func (u *Users) sqlTableName() string { //nolint:dupl false positive
	return `"users"`
}

func (u *Users) UniqueIndexId() string { //nolint:dupl false positive
	return `id`
}

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

var GraphqlFieldUsersById = &graphql.Field{
	Type:        GraphqlTypeUsers,
	Description: `list of Users`,
	Args: graphql.FieldConfigArgument{
		`Id`: &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
	},
}

func (g *Users) GraphqlFieldUsersByIdWithResolver() *graphql.Field {
	field := *GraphqlFieldUsersById
	field.Resolve = func(p graphql.ResolveParams) (interface{}, error) {
		q := g
		v, ok := p.Args[`id`]
		if !ok {
			v, _ = p.Args[`Id`]
		}
		q.Id = X.ToU(v)
		if q.FindById() {
			return q, nil
		}
		return nil, nil
	}
	return &field
}

func (u *Users) UniqueIndexEmail() string { //nolint:dupl false positive
	return `email`
}

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

func (u *Users) sqlSelectAllFields() string { //nolint:dupl false positive
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

func (u *Users) IdxId() int { //nolint:dupl false positive
	return 0
}

func (u *Users) sqlId() string { //nolint:dupl false positive
	return `"id"`
}

func (u *Users) IdxEmail() int { //nolint:dupl false positive
	return 1
}

func (u *Users) sqlEmail() string { //nolint:dupl false positive
	return `"email"`
}

func (u *Users) IdxPassword() int { //nolint:dupl false positive
	return 2
}

func (u *Users) sqlPassword() string { //nolint:dupl false positive
	return `"password"`
}

func (u *Users) IdxCreatedAt() int { //nolint:dupl false positive
	return 3
}

func (u *Users) sqlCreatedAt() string { //nolint:dupl false positive
	return `"createdAt"`
}

func (u *Users) IdxCreatedBy() int { //nolint:dupl false positive
	return 4
}

func (u *Users) sqlCreatedBy() string { //nolint:dupl false positive
	return `"createdBy"`
}

func (u *Users) IdxUpdatedAt() int { //nolint:dupl false positive
	return 5
}

func (u *Users) sqlUpdatedAt() string { //nolint:dupl false positive
	return `"updatedAt"`
}

func (u *Users) IdxUpdatedBy() int { //nolint:dupl false positive
	return 6
}

func (u *Users) sqlUpdatedBy() string { //nolint:dupl false positive
	return `"updatedBy"`
}

func (u *Users) IdxDeletedAt() int { //nolint:dupl false positive
	return 7
}

func (u *Users) sqlDeletedAt() string { //nolint:dupl false positive
	return `"deletedAt"`
}

func (u *Users) IdxDeletedBy() int { //nolint:dupl false positive
	return 8
}

func (u *Users) sqlDeletedBy() string { //nolint:dupl false positive
	return `"deletedBy"`
}

func (u *Users) IdxIsDeleted() int { //nolint:dupl false positive
	return 9
}

func (u *Users) sqlIsDeleted() string { //nolint:dupl false positive
	return `"isDeleted"`
}

func (u *Users) IdxRestoredAt() int { //nolint:dupl false positive
	return 10
}

func (u *Users) sqlRestoredAt() string { //nolint:dupl false positive
	return `"restoredAt"`
}

func (u *Users) IdxRestoredBy() int { //nolint:dupl false positive
	return 11
}

func (u *Users) sqlRestoredBy() string { //nolint:dupl false positive
	return `"restoredBy"`
}

func (u *Users) IdxPasswordSetAt() int { //nolint:dupl false positive
	return 12
}

func (u *Users) sqlPasswordSetAt() string { //nolint:dupl false positive
	return `"passwordSetAt"`
}

func (u *Users) IdxSecretCode() int { //nolint:dupl false positive
	return 13
}

func (u *Users) sqlSecretCode() string { //nolint:dupl false positive
	return `"secretCode"`
}

func (u *Users) IdxSecretCodeAt() int { //nolint:dupl false positive
	return 14
}

func (u *Users) sqlSecretCodeAt() string { //nolint:dupl false positive
	return `"secretCodeAt"`
}

func (u *Users) IdxVerificationSentAt() int { //nolint:dupl false positive
	return 15
}

func (u *Users) sqlVerificationSentAt() string { //nolint:dupl false positive
	return `"verificationSentAt"`
}

func (u *Users) IdxVerifiedAt() int { //nolint:dupl false positive
	return 16
}

func (u *Users) sqlVerifiedAt() string { //nolint:dupl false positive
	return `"verifiedAt"`
}

func (u *Users) IdxLastLoginAt() int { //nolint:dupl false positive
	return 17
}

func (u *Users) sqlLastLoginAt() string { //nolint:dupl false positive
	return `"lastLoginAt"`
}

func (u *Users) ToArray() A.X { //nolint:dupl false positive
	return A.X{
		u.Id,                 // 0
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

func (u *Users) Total() int64 { //nolint:dupl false positive
	rows := u.Adapter.CallBoxSpace(u.SpaceName()+`:count`, A.X{})
	if len(rows) > 0 && len(rows[0]) > 0 {
		return X.ToI(rows[0][0])
	}
	return 0
}

var GraphqlTypeUsers = graphql.NewObject(
	graphql.ObjectConfig{
		Name: `users`,
		Fields: graphql.Fields{
			`id`: &graphql.Field{
				Type: graphql.Int,
			},
			`email`: &graphql.Field{
				Type: graphql.String,
			},
			`createdAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`createdBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`updatedAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`updatedBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`deletedAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`deletedBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`isDeleted`: &graphql.Field{
				Type: graphql.Boolean,
			},
			`restoredAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`restoredBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`passwordSetAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`secretCodeAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`verificationSentAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`verifiedAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`lastLoginAt`: &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go
