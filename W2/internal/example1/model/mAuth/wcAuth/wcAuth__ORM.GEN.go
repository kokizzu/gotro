package wcAuth

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	"example1/model/mAuth/rqAuth"

	"github.com/tarantool/go-tarantool/v2"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/X"
)

// SessionsMutator DAO writer/command struct
//
//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file wcAuth__ORM.GEN.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type wcAuth__ORM.GEN.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type wcAuth__ORM.GEN.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type wcAuth__ORM.GEN.go
type SessionsMutator struct {
	rqAuth.Sessions
	mutations *tarantool.Operations
	logs      []A.X
}

// NewSessionsMutator create new ORM writer/command object
func NewSessionsMutator(adapter *Tt.Adapter) (res *SessionsMutator) {
	res = &SessionsMutator{Sessions: rqAuth.Sessions{Adapter: adapter}}
	res.mutations = tarantool.NewOperations()
	return
}

// Logs get array of logs [field, old, new]
func (s *SessionsMutator) Logs() []A.X { //nolint:dupl false positive
	return s.logs
}

// HaveMutation check whether Set* methods ever called
func (s *SessionsMutator) HaveMutation() bool { //nolint:dupl false positive
	return len(s.logs) > 0
}

// ClearMutations clear all previously called Set* methods
func (s *SessionsMutator) ClearMutations() { //nolint:dupl false positive
	s.mutations = tarantool.NewOperations()
	s.logs = []A.X{}
}

// func (s *SessionsMutator) DoUpsert() bool { //nolint:dupl false positive
//	arr := s.ToArray()
//	_, err := s.Adapter.Upsert(s.SpaceName(), arr, A.X{
//		A.X{`=`, 0, s.SessionToken},
//		A.X{`=`, 1, s.UserId},
//		A.X{`=`, 2, s.ExpiredAt},
//	})
//	return !L.IsError(err, `Sessions.DoUpsert failed: `+s.SpaceName()+ `\n%#v`, arr)
// }

// DoOverwriteBySessionToken update all columns, error if not exists, not using mutations/Set*
func (s *SessionsMutator) DoOverwriteBySessionToken() bool { //nolint:dupl false positive
	_, err := s.Adapter.Connection.Do(tarantool.NewUpdateRequest(s.SpaceName()).
		Index(s.UniqueIndexSessionToken()).
		Key(A.X{s.SessionToken}).
		Operations(s.ToUpdateArray()),
	).Get()
	return !L.IsError(err, `Sessions.DoOverwriteBySessionToken failed: `+s.SpaceName())
}

// DoUpdateBySessionToken update only mutated fields, error if not exists, use Find* and Set* methods instead of direct assignment
func (s *SessionsMutator) DoUpdateBySessionToken() bool { //nolint:dupl false positive
	if !s.HaveMutation() {
		return true
	}
	_, err := s.Adapter.Connection.Do(
		tarantool.NewUpdateRequest(s.SpaceName()).
			Index(s.UniqueIndexSessionToken()).
			Key(A.X{s.SessionToken}).
			Operations(s.mutations),
	).Get()
	return !L.IsError(err, `Sessions.DoUpdateBySessionToken failed: `+s.SpaceName())
}

// DoDeletePermanentBySessionToken permanent delete
func (s *SessionsMutator) DoDeletePermanentBySessionToken() bool { //nolint:dupl false positive
	_, err := s.Adapter.Connection.Do(
		tarantool.NewDeleteRequest(s.SpaceName()).
			Index(s.UniqueIndexSessionToken()).
			Key(A.X{s.SessionToken}),
	).Get()
	return !L.IsError(err, `Sessions.DoDeletePermanentBySessionToken failed: `+s.SpaceName())
}

// DoInsert insert, error if already exists
func (s *SessionsMutator) DoInsert() bool { //nolint:dupl false positive
	arr := s.ToArray()
	_, err := s.Adapter.Connection.Do(
		tarantool.NewInsertRequest(s.SpaceName()).
			Tuple(arr),
	).Get()
	return !L.IsError(err, `Sessions.DoInsert failed: `+s.SpaceName()+`\n%#v`, arr)
}

// DoUpsert upsert, insert or overwrite, will error only when there's unique secondary key being violated
// replace = upsert, only error when there's unique secondary key
// previous name: DoReplace
func (s *SessionsMutator) DoUpsert() bool { //nolint:dupl false positive
	arr := s.ToArray()
	_, err := s.Adapter.Connection.Do(
		tarantool.NewReplaceRequest(s.SpaceName()).
			Tuple(arr),
	).Get()
	return !L.IsError(err, `Sessions.DoUpsert failed: `+s.SpaceName()+`\n%#v`, arr)
}

// SetSessionToken create mutations, should not duplicate
func (s *SessionsMutator) SetSessionToken(val string) bool { //nolint:dupl false positive
	if val != s.SessionToken {
		s.mutations.Assign(0, val)
		s.logs = append(s.logs, A.X{`sessionToken`, s.SessionToken, val})
		s.SessionToken = val
		return true
	}
	return false
}

// SetUserId create mutations, should not duplicate
func (s *SessionsMutator) SetUserId(val uint64) bool { //nolint:dupl false positive
	if val != s.UserId {
		s.mutations.Assign(1, val)
		s.logs = append(s.logs, A.X{`userId`, s.UserId, val})
		s.UserId = val
		return true
	}
	return false
}

// SetExpiredAt create mutations, should not duplicate
func (s *SessionsMutator) SetExpiredAt(val int64) bool { //nolint:dupl false positive
	if val != s.ExpiredAt {
		s.mutations.Assign(2, val)
		s.logs = append(s.logs, A.X{`expiredAt`, s.ExpiredAt, val})
		s.ExpiredAt = val
		return true
	}
	return false
}

// SetAll set all from another source, only if another property is not empty/nil/zero or in forceMap
func (s *SessionsMutator) SetAll(from rqAuth.Sessions, excludeMap, forceMap M.SB) (changed bool) { //nolint:dupl false positive
	if excludeMap == nil { // list of fields to exclude
		excludeMap = M.SB{}
	}
	if forceMap == nil { // list of fields to force overwrite
		forceMap = M.SB{}
	}
	if !excludeMap[`sessionToken`] && (forceMap[`sessionToken`] || from.SessionToken != ``) {
		s.SessionToken = S.Trim(from.SessionToken)
		changed = true
	}
	if !excludeMap[`userId`] && (forceMap[`userId`] || from.UserId != 0) {
		s.UserId = from.UserId
		changed = true
	}
	if !excludeMap[`expiredAt`] && (forceMap[`expiredAt`] || from.ExpiredAt != 0) {
		s.ExpiredAt = from.ExpiredAt
		changed = true
	}
	return
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

// UsersMutator DAO writer/command struct
type UsersMutator struct {
	rqAuth.Users
	mutations *tarantool.Operations
	logs      []A.X
}

// NewUsersMutator create new ORM writer/command object
func NewUsersMutator(adapter *Tt.Adapter) (res *UsersMutator) {
	res = &UsersMutator{Users: rqAuth.Users{Adapter: adapter}}
	res.mutations = tarantool.NewOperations()
	return
}

// Logs get array of logs [field, old, new]
func (u *UsersMutator) Logs() []A.X { //nolint:dupl false positive
	return u.logs
}

// HaveMutation check whether Set* methods ever called
func (u *UsersMutator) HaveMutation() bool { //nolint:dupl false positive
	return len(u.logs) > 0
}

// ClearMutations clear all previously called Set* methods
func (u *UsersMutator) ClearMutations() { //nolint:dupl false positive
	u.mutations = tarantool.NewOperations()
	u.logs = []A.X{}
}

// DoOverwriteById update all columns, error if not exists, not using mutations/Set*
func (u *UsersMutator) DoOverwriteById() bool { //nolint:dupl false positive
	_, err := u.Adapter.Connection.Do(tarantool.NewUpdateRequest(u.SpaceName()).
		Index(u.UniqueIndexId()).
		Key(A.X{u.Id}).
		Operations(u.ToUpdateArray()),
	).Get()
	return !L.IsError(err, `Users.DoOverwriteById failed: `+u.SpaceName())
}

// DoUpdateById update only mutated fields, error if not exists, use Find* and Set* methods instead of direct assignment
func (u *UsersMutator) DoUpdateById() bool { //nolint:dupl false positive
	if !u.HaveMutation() {
		return true
	}
	_, err := u.Adapter.Connection.Do(
		tarantool.NewUpdateRequest(u.SpaceName()).
			Index(u.UniqueIndexId()).
			Key(A.X{u.Id}).
			Operations(u.mutations),
	).Get()
	return !L.IsError(err, `Users.DoUpdateById failed: `+u.SpaceName())
}

// DoDeletePermanentById permanent delete
func (u *UsersMutator) DoDeletePermanentById() bool { //nolint:dupl false positive
	_, err := u.Adapter.Connection.Do(
		tarantool.NewDeleteRequest(u.SpaceName()).
			Index(u.UniqueIndexId()).
			Key(A.X{u.Id}),
	).Get()
	return !L.IsError(err, `Users.DoDeletePermanentById failed: `+u.SpaceName())
}

// func (u *UsersMutator) DoUpsert() bool { //nolint:dupl false positive
//	arr := u.ToArray()
//	_, err := u.Adapter.Upsert(u.SpaceName(), arr, A.X{
//		A.X{`=`, 0, u.Id},
//		A.X{`=`, 1, u.Email},
//		A.X{`=`, 2, u.Password},
//		A.X{`=`, 3, u.CreatedAt},
//		A.X{`=`, 4, u.CreatedBy},
//		A.X{`=`, 5, u.UpdatedAt},
//		A.X{`=`, 6, u.UpdatedBy},
//		A.X{`=`, 7, u.DeletedAt},
//		A.X{`=`, 8, u.DeletedBy},
//		A.X{`=`, 9, u.IsDeleted},
//		A.X{`=`, 10, u.RestoredAt},
//		A.X{`=`, 11, u.RestoredBy},
//		A.X{`=`, 12, u.PasswordSetAt},
//		A.X{`=`, 13, u.SecretCode},
//		A.X{`=`, 14, u.SecretCodeAt},
//		A.X{`=`, 15, u.VerificationSentAt},
//		A.X{`=`, 16, u.VerifiedAt},
//		A.X{`=`, 17, u.LastLoginAt},
//	})
//	return !L.IsError(err, `Users.DoUpsert failed: `+u.SpaceName()+ `\n%#v`, arr)
// }

// DoOverwriteByEmail update all columns, error if not exists, not using mutations/Set*
func (u *UsersMutator) DoOverwriteByEmail() bool { //nolint:dupl false positive
	_, err := u.Adapter.Connection.Do(tarantool.NewUpdateRequest(u.SpaceName()).
		Index(u.UniqueIndexEmail()).
		Key(A.X{u.Email}).
		Operations(u.ToUpdateArray()),
	).Get()
	return !L.IsError(err, `Users.DoOverwriteByEmail failed: `+u.SpaceName())
}

// DoUpdateByEmail update only mutated fields, error if not exists, use Find* and Set* methods instead of direct assignment
func (u *UsersMutator) DoUpdateByEmail() bool { //nolint:dupl false positive
	if !u.HaveMutation() {
		return true
	}
	_, err := u.Adapter.Connection.Do(
		tarantool.NewUpdateRequest(u.SpaceName()).
			Index(u.UniqueIndexEmail()).
			Key(A.X{u.Email}).
			Operations(u.mutations),
	).Get()
	return !L.IsError(err, `Users.DoUpdateByEmail failed: `+u.SpaceName())
}

// DoDeletePermanentByEmail permanent delete
func (u *UsersMutator) DoDeletePermanentByEmail() bool { //nolint:dupl false positive
	_, err := u.Adapter.Connection.Do(
		tarantool.NewDeleteRequest(u.SpaceName()).
			Index(u.UniqueIndexEmail()).
			Key(A.X{u.Email}),
	).Get()
	return !L.IsError(err, `Users.DoDeletePermanentByEmail failed: `+u.SpaceName())
}

// DoInsert insert, error if already exists
func (u *UsersMutator) DoInsert() bool { //nolint:dupl false positive
	arr := u.ToArray()
	row, err := u.Adapter.Connection.Do(
		tarantool.NewInsertRequest(u.SpaceName()).
			Tuple(arr),
	).Get()
	if err == nil {
		if len(row) > 0 {
			if cells, ok := row[0].([]any); ok && len(cells) > 0 {
				u.Id = X.ToU(cells[0])
			}
		}
	}
	return !L.IsError(err, `Users.DoInsert failed: `+u.SpaceName()+`\n%#v`, arr)
}

// DoUpsert upsert, insert or overwrite, will error only when there's unique secondary key being violated
// replace = upsert, only error when there's unique secondary key
// previous name: DoReplace
func (u *UsersMutator) DoUpsert() bool { //nolint:dupl false positive
	arr := u.ToArray()
	row, err := u.Adapter.Connection.Do(
		tarantool.NewReplaceRequest(u.SpaceName()).
			Tuple(arr),
	).Get()
	if err == nil {
		if len(row) > 0 {
			if cells, ok := row[0].([]any); ok && len(cells) > 0 {
				u.Id = X.ToU(cells[0])
			}
		}
	}
	return !L.IsError(err, `Users.DoUpsert failed: `+u.SpaceName()+`\n%#v`, arr)
}

// SetId create mutations, should not duplicate
func (u *UsersMutator) SetId(val uint64) bool { //nolint:dupl false positive
	if val != u.Id {
		u.mutations.Assign(0, val)
		u.logs = append(u.logs, A.X{`id`, u.Id, val})
		u.Id = val
		return true
	}
	return false
}

// SetEmail create mutations, should not duplicate
func (u *UsersMutator) SetEmail(val string) bool { //nolint:dupl false positive
	if val != u.Email {
		u.mutations.Assign(1, val)
		u.logs = append(u.logs, A.X{`email`, u.Email, val})
		u.Email = val
		return true
	}
	return false
}

// SetPassword create mutations, should not duplicate
func (u *UsersMutator) SetPassword(val string) bool { //nolint:dupl false positive
	if val != u.Password {
		u.mutations.Assign(2, val)
		u.logs = append(u.logs, A.X{`password`, u.Password, val})
		u.Password = val
		return true
	}
	return false
}

// SetCreatedAt create mutations, should not duplicate
func (u *UsersMutator) SetCreatedAt(val int64) bool { //nolint:dupl false positive
	if val != u.CreatedAt {
		u.mutations.Assign(3, val)
		u.logs = append(u.logs, A.X{`createdAt`, u.CreatedAt, val})
		u.CreatedAt = val
		return true
	}
	return false
}

// SetCreatedBy create mutations, should not duplicate
func (u *UsersMutator) SetCreatedBy(val uint64) bool { //nolint:dupl false positive
	if val != u.CreatedBy {
		u.mutations.Assign(4, val)
		u.logs = append(u.logs, A.X{`createdBy`, u.CreatedBy, val})
		u.CreatedBy = val
		return true
	}
	return false
}

// SetUpdatedAt create mutations, should not duplicate
func (u *UsersMutator) SetUpdatedAt(val int64) bool { //nolint:dupl false positive
	if val != u.UpdatedAt {
		u.mutations.Assign(5, val)
		u.logs = append(u.logs, A.X{`updatedAt`, u.UpdatedAt, val})
		u.UpdatedAt = val
		return true
	}
	return false
}

// SetUpdatedBy create mutations, should not duplicate
func (u *UsersMutator) SetUpdatedBy(val uint64) bool { //nolint:dupl false positive
	if val != u.UpdatedBy {
		u.mutations.Assign(6, val)
		u.logs = append(u.logs, A.X{`updatedBy`, u.UpdatedBy, val})
		u.UpdatedBy = val
		return true
	}
	return false
}

// SetDeletedAt create mutations, should not duplicate
func (u *UsersMutator) SetDeletedAt(val int64) bool { //nolint:dupl false positive
	if val != u.DeletedAt {
		u.mutations.Assign(7, val)
		u.logs = append(u.logs, A.X{`deletedAt`, u.DeletedAt, val})
		u.DeletedAt = val
		return true
	}
	return false
}

// SetDeletedBy create mutations, should not duplicate
func (u *UsersMutator) SetDeletedBy(val uint64) bool { //nolint:dupl false positive
	if val != u.DeletedBy {
		u.mutations.Assign(8, val)
		u.logs = append(u.logs, A.X{`deletedBy`, u.DeletedBy, val})
		u.DeletedBy = val
		return true
	}
	return false
}

// SetIsDeleted create mutations, should not duplicate
func (u *UsersMutator) SetIsDeleted(val bool) bool { //nolint:dupl false positive
	if val != u.IsDeleted {
		u.mutations.Assign(9, val)
		u.logs = append(u.logs, A.X{`isDeleted`, u.IsDeleted, val})
		u.IsDeleted = val
		return true
	}
	return false
}

// SetRestoredAt create mutations, should not duplicate
func (u *UsersMutator) SetRestoredAt(val int64) bool { //nolint:dupl false positive
	if val != u.RestoredAt {
		u.mutations.Assign(10, val)
		u.logs = append(u.logs, A.X{`restoredAt`, u.RestoredAt, val})
		u.RestoredAt = val
		return true
	}
	return false
}

// SetRestoredBy create mutations, should not duplicate
func (u *UsersMutator) SetRestoredBy(val uint64) bool { //nolint:dupl false positive
	if val != u.RestoredBy {
		u.mutations.Assign(11, val)
		u.logs = append(u.logs, A.X{`restoredBy`, u.RestoredBy, val})
		u.RestoredBy = val
		return true
	}
	return false
}

// SetPasswordSetAt create mutations, should not duplicate
func (u *UsersMutator) SetPasswordSetAt(val int64) bool { //nolint:dupl false positive
	if val != u.PasswordSetAt {
		u.mutations.Assign(12, val)
		u.logs = append(u.logs, A.X{`passwordSetAt`, u.PasswordSetAt, val})
		u.PasswordSetAt = val
		return true
	}
	return false
}

// SetSecretCode create mutations, should not duplicate
func (u *UsersMutator) SetSecretCode(val string) bool { //nolint:dupl false positive
	if val != u.SecretCode {
		u.mutations.Assign(13, val)
		u.logs = append(u.logs, A.X{`secretCode`, u.SecretCode, val})
		u.SecretCode = val
		return true
	}
	return false
}

// SetSecretCodeAt create mutations, should not duplicate
func (u *UsersMutator) SetSecretCodeAt(val int64) bool { //nolint:dupl false positive
	if val != u.SecretCodeAt {
		u.mutations.Assign(14, val)
		u.logs = append(u.logs, A.X{`secretCodeAt`, u.SecretCodeAt, val})
		u.SecretCodeAt = val
		return true
	}
	return false
}

// SetVerificationSentAt create mutations, should not duplicate
func (u *UsersMutator) SetVerificationSentAt(val int64) bool { //nolint:dupl false positive
	if val != u.VerificationSentAt {
		u.mutations.Assign(15, val)
		u.logs = append(u.logs, A.X{`verificationSentAt`, u.VerificationSentAt, val})
		u.VerificationSentAt = val
		return true
	}
	return false
}

// SetVerifiedAt create mutations, should not duplicate
func (u *UsersMutator) SetVerifiedAt(val int64) bool { //nolint:dupl false positive
	if val != u.VerifiedAt {
		u.mutations.Assign(16, val)
		u.logs = append(u.logs, A.X{`verifiedAt`, u.VerifiedAt, val})
		u.VerifiedAt = val
		return true
	}
	return false
}

// SetLastLoginAt create mutations, should not duplicate
func (u *UsersMutator) SetLastLoginAt(val int64) bool { //nolint:dupl false positive
	if val != u.LastLoginAt {
		u.mutations.Assign(17, val)
		u.logs = append(u.logs, A.X{`lastLoginAt`, u.LastLoginAt, val})
		u.LastLoginAt = val
		return true
	}
	return false
}

// SetAll set all from another source, only if another property is not empty/nil/zero or in forceMap
func (u *UsersMutator) SetAll(from rqAuth.Users, excludeMap, forceMap M.SB) (changed bool) { //nolint:dupl false positive
	if excludeMap == nil { // list of fields to exclude
		excludeMap = M.SB{}
	}
	if forceMap == nil { // list of fields to force overwrite
		forceMap = M.SB{}
	}
	if !excludeMap[`id`] && (forceMap[`id`] || from.Id != 0) {
		u.Id = from.Id
		changed = true
	}
	if !excludeMap[`email`] && (forceMap[`email`] || from.Email != ``) {
		u.Email = S.Trim(from.Email)
		changed = true
	}
	if !excludeMap[`password`] && (forceMap[`password`] || from.Password != ``) {
		u.Password = S.Trim(from.Password)
		changed = true
	}
	if !excludeMap[`createdAt`] && (forceMap[`createdAt`] || from.CreatedAt != 0) {
		u.CreatedAt = from.CreatedAt
		changed = true
	}
	if !excludeMap[`createdBy`] && (forceMap[`createdBy`] || from.CreatedBy != 0) {
		u.CreatedBy = from.CreatedBy
		changed = true
	}
	if !excludeMap[`updatedAt`] && (forceMap[`updatedAt`] || from.UpdatedAt != 0) {
		u.UpdatedAt = from.UpdatedAt
		changed = true
	}
	if !excludeMap[`updatedBy`] && (forceMap[`updatedBy`] || from.UpdatedBy != 0) {
		u.UpdatedBy = from.UpdatedBy
		changed = true
	}
	if !excludeMap[`deletedAt`] && (forceMap[`deletedAt`] || from.DeletedAt != 0) {
		u.DeletedAt = from.DeletedAt
		changed = true
	}
	if !excludeMap[`deletedBy`] && (forceMap[`deletedBy`] || from.DeletedBy != 0) {
		u.DeletedBy = from.DeletedBy
		changed = true
	}
	if !excludeMap[`isDeleted`] && (forceMap[`isDeleted`] || from.IsDeleted != false) {
		u.IsDeleted = from.IsDeleted
		changed = true
	}
	if !excludeMap[`restoredAt`] && (forceMap[`restoredAt`] || from.RestoredAt != 0) {
		u.RestoredAt = from.RestoredAt
		changed = true
	}
	if !excludeMap[`restoredBy`] && (forceMap[`restoredBy`] || from.RestoredBy != 0) {
		u.RestoredBy = from.RestoredBy
		changed = true
	}
	if !excludeMap[`passwordSetAt`] && (forceMap[`passwordSetAt`] || from.PasswordSetAt != 0) {
		u.PasswordSetAt = from.PasswordSetAt
		changed = true
	}
	if !excludeMap[`secretCode`] && (forceMap[`secretCode`] || from.SecretCode != ``) {
		u.SecretCode = S.Trim(from.SecretCode)
		changed = true
	}
	if !excludeMap[`secretCodeAt`] && (forceMap[`secretCodeAt`] || from.SecretCodeAt != 0) {
		u.SecretCodeAt = from.SecretCodeAt
		changed = true
	}
	if !excludeMap[`verificationSentAt`] && (forceMap[`verificationSentAt`] || from.VerificationSentAt != 0) {
		u.VerificationSentAt = from.VerificationSentAt
		changed = true
	}
	if !excludeMap[`verifiedAt`] && (forceMap[`verifiedAt`] || from.VerifiedAt != 0) {
		u.VerifiedAt = from.VerifiedAt
		changed = true
	}
	if !excludeMap[`lastLoginAt`] && (forceMap[`lastLoginAt`] || from.LastLoginAt != 0) {
		u.LastLoginAt = from.LastLoginAt
		changed = true
	}
	return
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go
