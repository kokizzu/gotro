package mAuth

import (
	"github.com/kokizzu/gotro/D/Ch"
	"github.com/kokizzu/gotro/D/Tt"
)

// table users, sessions
const (
	TableUsers         Tt.TableName = `users`
	Id                              = `id`
	Email                           = `email`
	Password                        = `password`
	CreatedBy                       = `createdBy`
	CreatedAt                       = `createdAt`
	UpdatedBy                       = `updatedBy`
	UpdatedAt                       = `updatedAt`
	DeletedBy                       = `deletedBy`
	DeletedAt                       = `deletedAt`
	IsDeleted                       = `isDeleted`
	RestoredBy                      = `restoredBy`
	RestoredAt                      = `restoredAt`
	PasswordSetAt                   = `passwordSetAt`
	SecretCode                      = `secretCode`
	SecretCodeAt                    = `secretCodeAt`
	VerificationSentAt              = `verificationSentAt`
	VerifiedAt                      = `verifiedAt`
	LastLoginAt                     = `lastLoginAt`
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
	TableUsers: {
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
		Unique1:      Email,
		Indexes:      []string{IsDeleted, SecretCode},
		HiddenFields: []string{Password, SecretCode},
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

// table userlogs
const (
	TableUserLogs Ch.TableName = `userLogs`
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
		},
		Orders: []string{ActorId, RequestId},
	},
}

func GenerateORM() {
	Tt.GenerateOrm(TarantoolTables)
	Ch.GenerateOrm(ClickhouseTables) // find d.InitClickhouseBuffer to create chBuffer on NewDomain
}
