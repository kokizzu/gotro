package rqAuth

import (
	"golang.org/x/crypto/bcrypt"

	"example1/conf"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/X"
)

func (s *Users) AllOffsetLimit(offset, limit uint32) (res []*Users) {
	query := `
SELECT ` + s.SqlSelectAllFields() + `
FROM ` + s.SqlTableName() + `
ORDER BY ` + s.SqlId() + `
LIMIT ` + X.ToS(limit) + `
OFFSET ` + X.ToS(offset)
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

func (s *Users) CensorFields() {
	s.Password = ``
	s.SecretCode = ``
}

func (s *Users) CheckPassword(currentPassword string) bool {
	hash := []byte(s.Password)
	pass := []byte(currentPassword)
	err := bcrypt.CompareHashAndPassword(hash, pass)

	return !L.IsError(err, `bcrypt.CompareHashAndPassword`)
}

// add more custom methods here
