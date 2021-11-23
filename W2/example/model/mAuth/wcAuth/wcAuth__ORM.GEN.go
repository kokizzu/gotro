package wcAuth

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	"github.com/kokizzu/gotro/W2/example/model/mAuth/rqAuth"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/X"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file wcAuth__ORM.GEN.go
//go:generate replacer 'Id" form' 'Id,string" form' type wcAuth__ORM.GEN.go
//go:generate replacer 'json:"id"' 'json:"id,string"' type wcAuth__ORM.GEN.go
//go:generate replacer 'By" form' 'By,string" form' type wcAuth__ORM.GEN.go
// go:generate msgp -tests=false -file wcAuth__ORM.GEN.go -o wcAuth__MSG.GEN.go

type SessionsMutator struct {
	rqAuth.Sessions
	mutations []A.X
}

func NewSessionsMutator(adapter *Tt.Adapter) *SessionsMutator {
	return &SessionsMutator{Sessions: rqAuth.Sessions{Adapter: adapter}}
}

func (s *SessionsMutator) HaveMutation() bool { //nolint:dupl false positive
	return len(s.mutations) > 0
}

// func (s *SessionsMutator) DoUpsert() bool { //nolint:dupl false positive
//	_, err := s.Adapter.Upsert(s.SpaceName(), s.ToArray(), A.X{
//		A.X{`=`, 0, s.SessionToken},
//		A.X{`=`, 1, s.UserId},
//		A.X{`=`, 2, s.ExpiredAt},
//	})
//	return !L.IsError(err, `Sessions.DoUpsert failed: `+s.SpaceName())
// }

// Overwrite all columns, error if not exists
func (s *SessionsMutator) DoOverwriteBySessionToken() bool { //nolint:dupl false positive
	_, err := s.Adapter.Update(s.SpaceName(), s.UniqueIndexSessionToken(), A.X{s.SessionToken}, s.ToUpdateArray())
	return !L.IsError(err, `Sessions.DoOverwriteBySessionToken failed: `+s.SpaceName())
}

// Update only mutated, error if not exists, use Find* and Set* methods instead of direct assignment
func (s *SessionsMutator) DoUpdateBySessionToken() bool { //nolint:dupl false positive
	if !s.HaveMutation() {
		return true
	}
	_, err := s.Adapter.Update(s.SpaceName(), s.UniqueIndexSessionToken(), A.X{s.SessionToken}, s.mutations)
	return !L.IsError(err, `Sessions.DoUpdateBySessionToken failed: `+s.SpaceName())
}

func (s *SessionsMutator) DoDeletePermanentBySessionToken() bool { //nolint:dupl false positive
	_, err := s.Adapter.Delete(s.SpaceName(), s.UniqueIndexSessionToken(), A.X{s.SessionToken})
	return !L.IsError(err, `Sessions.DoDeletePermanentBySessionToken failed: `+s.SpaceName())
}

// insert, error if exists
func (s *SessionsMutator) DoInsert() bool { //nolint:dupl false positive
	_, err := s.Adapter.Insert(s.SpaceName(), s.ToArray())
	return !L.IsError(err, `Sessions.DoInsert failed: `+s.SpaceName())
}

// replace = upsert, only error when there's unique secondary key
func (s *SessionsMutator) DoReplace() bool { //nolint:dupl false positive
	_, err := s.Adapter.Replace(s.SpaceName(), s.ToArray())
	return !L.IsError(err, `Sessions.DoReplace failed: `+s.SpaceName())
}

func (s *SessionsMutator) SetSessionToken(val string) bool { //nolint:dupl false positive
	if val != s.SessionToken {
		s.mutations = append(s.mutations, A.X{`=`, 0, val})
		s.SessionToken = val
		return true
	}
	return false
}

func (s *SessionsMutator) SetUserId(val uint64) bool { //nolint:dupl false positive
	if val != s.UserId {
		s.mutations = append(s.mutations, A.X{`=`, 1, val})
		s.UserId = val
		return true
	}
	return false
}

func (s *SessionsMutator) SetExpiredAt(val int64) bool { //nolint:dupl false positive
	if val != s.ExpiredAt {
		s.mutations = append(s.mutations, A.X{`=`, 2, val})
		s.ExpiredAt = val
		return true
	}
	return false
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

type UsersMutator struct {
	rqAuth.Users
	mutations []A.X
}

func NewUsersMutator(adapter *Tt.Adapter) *UsersMutator {
	return &UsersMutator{Users: rqAuth.Users{Adapter: adapter}}
}

func (u *UsersMutator) HaveMutation() bool { //nolint:dupl false positive
	return len(u.mutations) > 0
}

// Overwrite all columns, error if not exists
func (u *UsersMutator) DoOverwriteById() bool { //nolint:dupl false positive
	_, err := u.Adapter.Update(u.SpaceName(), u.UniqueIndexId(), A.X{u.Id}, u.ToUpdateArray())
	return !L.IsError(err, `Users.DoOverwriteById failed: `+u.SpaceName())
}

// Update only mutated, error if not exists, use Find* and Set* methods instead of direct assignment
func (u *UsersMutator) DoUpdateById() bool { //nolint:dupl false positive
	if !u.HaveMutation() {
		return true
	}
	_, err := u.Adapter.Update(u.SpaceName(), u.UniqueIndexId(), A.X{u.Id}, u.mutations)
	return !L.IsError(err, `Users.DoUpdateById failed: `+u.SpaceName())
}

func (u *UsersMutator) DoDeletePermanentById() bool { //nolint:dupl false positive
	_, err := u.Adapter.Delete(u.SpaceName(), u.UniqueIndexId(), A.X{u.Id})
	return !L.IsError(err, `Users.DoDeletePermanentById failed: `+u.SpaceName())
}

// func (u *UsersMutator) DoUpsert() bool { //nolint:dupl false positive
//	_, err := u.Adapter.Upsert(u.SpaceName(), u.ToArray(), A.X{
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
//	return !L.IsError(err, `Users.DoUpsert failed: `+u.SpaceName())
// }

// Overwrite all columns, error if not exists
func (u *UsersMutator) DoOverwriteByEmail() bool { //nolint:dupl false positive
	_, err := u.Adapter.Update(u.SpaceName(), u.UniqueIndexEmail(), A.X{u.Email}, u.ToUpdateArray())
	return !L.IsError(err, `Users.DoOverwriteByEmail failed: `+u.SpaceName())
}

// Update only mutated, error if not exists, use Find* and Set* methods instead of direct assignment
func (u *UsersMutator) DoUpdateByEmail() bool { //nolint:dupl false positive
	if !u.HaveMutation() {
		return true
	}
	_, err := u.Adapter.Update(u.SpaceName(), u.UniqueIndexEmail(), A.X{u.Email}, u.mutations)
	return !L.IsError(err, `Users.DoUpdateByEmail failed: `+u.SpaceName())
}

func (u *UsersMutator) DoDeletePermanentByEmail() bool { //nolint:dupl false positive
	_, err := u.Adapter.Delete(u.SpaceName(), u.UniqueIndexEmail(), A.X{u.Email})
	return !L.IsError(err, `Users.DoDeletePermanentByEmail failed: `+u.SpaceName())
}

// insert, error if exists
func (u *UsersMutator) DoInsert() bool { //nolint:dupl false positive
	row, err := u.Adapter.Insert(u.SpaceName(), u.ToArray())
	if err == nil {
		tup := row.Tuples()
		if len(tup) > 0 && len(tup[0]) > 0 && tup[0][0] != nil {
			u.Id = X.ToU(tup[0][0])
		}
	}
	return !L.IsError(err, `Users.DoInsert failed: `+u.SpaceName())
}

// replace = upsert, only error when there's unique secondary key
func (u *UsersMutator) DoReplace() bool { //nolint:dupl false positive
	_, err := u.Adapter.Replace(u.SpaceName(), u.ToArray())
	return !L.IsError(err, `Users.DoReplace failed: `+u.SpaceName())
}

func (u *UsersMutator) SetId(val uint64) bool { //nolint:dupl false positive
	if val != u.Id {
		u.mutations = append(u.mutations, A.X{`=`, 0, val})
		u.Id = val
		return true
	}
	return false
}

func (u *UsersMutator) SetEmail(val string) bool { //nolint:dupl false positive
	if val != u.Email {
		u.mutations = append(u.mutations, A.X{`=`, 1, val})
		u.Email = val
		return true
	}
	return false
}

func (u *UsersMutator) SetPassword(val string) bool { //nolint:dupl false positive
	if val != u.Password {
		u.mutations = append(u.mutations, A.X{`=`, 2, val})
		u.Password = val
		return true
	}
	return false
}

func (u *UsersMutator) SetCreatedAt(val int64) bool { //nolint:dupl false positive
	if val != u.CreatedAt {
		u.mutations = append(u.mutations, A.X{`=`, 3, val})
		u.CreatedAt = val
		return true
	}
	return false
}

func (u *UsersMutator) SetCreatedBy(val uint64) bool { //nolint:dupl false positive
	if val != u.CreatedBy {
		u.mutations = append(u.mutations, A.X{`=`, 4, val})
		u.CreatedBy = val
		return true
	}
	return false
}

func (u *UsersMutator) SetUpdatedAt(val int64) bool { //nolint:dupl false positive
	if val != u.UpdatedAt {
		u.mutations = append(u.mutations, A.X{`=`, 5, val})
		u.UpdatedAt = val
		return true
	}
	return false
}

func (u *UsersMutator) SetUpdatedBy(val uint64) bool { //nolint:dupl false positive
	if val != u.UpdatedBy {
		u.mutations = append(u.mutations, A.X{`=`, 6, val})
		u.UpdatedBy = val
		return true
	}
	return false
}

func (u *UsersMutator) SetDeletedAt(val int64) bool { //nolint:dupl false positive
	if val != u.DeletedAt {
		u.mutations = append(u.mutations, A.X{`=`, 7, val})
		u.DeletedAt = val
		return true
	}
	return false
}

func (u *UsersMutator) SetDeletedBy(val uint64) bool { //nolint:dupl false positive
	if val != u.DeletedBy {
		u.mutations = append(u.mutations, A.X{`=`, 8, val})
		u.DeletedBy = val
		return true
	}
	return false
}

func (u *UsersMutator) SetIsDeleted(val bool) bool { //nolint:dupl false positive
	if val != u.IsDeleted {
		u.mutations = append(u.mutations, A.X{`=`, 9, val})
		u.IsDeleted = val
		return true
	}
	return false
}

func (u *UsersMutator) SetRestoredAt(val int64) bool { //nolint:dupl false positive
	if val != u.RestoredAt {
		u.mutations = append(u.mutations, A.X{`=`, 10, val})
		u.RestoredAt = val
		return true
	}
	return false
}

func (u *UsersMutator) SetRestoredBy(val uint64) bool { //nolint:dupl false positive
	if val != u.RestoredBy {
		u.mutations = append(u.mutations, A.X{`=`, 11, val})
		u.RestoredBy = val
		return true
	}
	return false
}

func (u *UsersMutator) SetPasswordSetAt(val int64) bool { //nolint:dupl false positive
	if val != u.PasswordSetAt {
		u.mutations = append(u.mutations, A.X{`=`, 12, val})
		u.PasswordSetAt = val
		return true
	}
	return false
}

func (u *UsersMutator) SetSecretCode(val string) bool { //nolint:dupl false positive
	if val != u.SecretCode {
		u.mutations = append(u.mutations, A.X{`=`, 13, val})
		u.SecretCode = val
		return true
	}
	return false
}

func (u *UsersMutator) SetSecretCodeAt(val int64) bool { //nolint:dupl false positive
	if val != u.SecretCodeAt {
		u.mutations = append(u.mutations, A.X{`=`, 14, val})
		u.SecretCodeAt = val
		return true
	}
	return false
}

func (u *UsersMutator) SetVerificationSentAt(val int64) bool { //nolint:dupl false positive
	if val != u.VerificationSentAt {
		u.mutations = append(u.mutations, A.X{`=`, 15, val})
		u.VerificationSentAt = val
		return true
	}
	return false
}

func (u *UsersMutator) SetVerifiedAt(val int64) bool { //nolint:dupl false positive
	if val != u.VerifiedAt {
		u.mutations = append(u.mutations, A.X{`=`, 16, val})
		u.VerifiedAt = val
		return true
	}
	return false
}

func (u *UsersMutator) SetLastLoginAt(val int64) bool { //nolint:dupl false positive
	if val != u.LastLoginAt {
		u.mutations = append(u.mutations, A.X{`=`, 17, val})
		u.LastLoginAt = val
		return true
	}
	return false
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go
