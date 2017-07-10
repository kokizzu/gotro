package tUsers

import (
	"example-complete/sql"
	"github.com/kokizzu/gotro/D/Pg"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/W"
)

const TABLE = `users`

var TM_MASTER Pg.TableModel
var SELECT = ``

var Z func(string) string
var ZZ func(string) string
var ZJ func(string) string
var ZB func(bool) string
var ZI func(int64) string
var ZLIKE func(string) string
var ZT func(strs ...string) string
var PG *Pg.RDBMS

func init() {
	Z = S.Z
	ZB = S.ZB
	ZZ = S.ZZ
	ZJ = S.ZJJ
	ZI = S.ZI
	ZT = S.ZT
	ZLIKE = S.ZLIKE
	PG = sql.PG

	TM_MASTER = Pg.TableModel{
		CacheName: TABLE + `_MASTER`,
		Fields: []Pg.FieldModel{
			{Key: `id`},
			{Label: `Full Name`, Key: `full_name`},
			{Label: `Verified`, Key: `verified`, Type: `bool`},
			{Label: `E-Mail`, Key: `email`},
			{Label: `Note`, Key: `note`},
		},
	}
	SELECT = TM_MASTER.Select()
}

// 2017-05-30 Prayogo
func FindID_ByPhone(ident string) int64 {
	query := ZT(ident) + `
	SELECT COALESCE((
		SELECT id
		FROM ` + TABLE + `
		WHERE is_deleted = false
			AND phone = ` + Z(ident) + `
		LIMIT 1
	),0)`
	return PG.QInt(query)
}

// 2017-05-30 Prayogo
func FindID_ByIdent_ByPass(ident, pass string) int64 {
	pass = S.HashPassword(pass)
	query := ZT(ident, pass) + `
	SELECT COALESCE((
		SELECT id
		FROM ` + TABLE + `
		WHERE is_deleted = false
			AND email = ` + Z(ident) + `
			AND password = ` + Z(pass) + `
		LIMIT 1
	),0)`
	return PG.QInt(query)
}

// 2017-05-30 Prayogo
func FindID_ByEmail(email string) int64 {
	query := ZT(email) + `
	SELECT COALESCE((
		SELECT id
		FROM ` + TABLE + `
		WHERE is_deleted = false
			AND email = ` + Z(email) + `
		LIMIT 1
	),0)`
	return PG.QInt(query)
}

// 2017-05-30 Prayogo
func FindID_ByCompactName_ByEmail(ident, email string) int64 {
	ident = S.Trim(ident)
	email = S.Trim(email)
	if email == `` {
		return 0
	}
	ident = Z(ident)
	email = Z(S.ToLower(email))
	query := ZT(ident, email) + `
	SELECT COALESCE((
		SELECT id
		FROM ` + TABLE + `
		WHERE is_deleted = false
			AND email = ` + email + `
	),0)`
	return PG.QInt(query)
}

// 2017-05-30 Prayogo
func Name_Emails_ByID(id int64) (string, []string) {
	ram_key := ZT(I.ToS(id))
	query := ram_key + `
SELECT full_name, email
FROM ` + TABLE + `
WHERE id = ` + ZI(id)
	name, mail := PG.QStrStr(query)
	return name, []string{name + ` <` + mail + `>`}
}

// 2017-06-04 Haries
func Search_ByQueryParams(qp *Pg.QueryParams) {
	qp.RamKey = ZT(qp.Term)
	if qp.Term != `` {
		qp.Where += ` AND x1.full_name ILIKE ` + ZLIKE(qp.Term)
	}
	qp.From = `FROM ` + TABLE + ` x1`
	qp.OrderBy = `x1.id`
	qp.Select = SELECT
	qp.SearchQuery_ByConn(PG)
}

// 2017-06-04 Haries
func API_Superadmin_Search(rm *W.RequestModel) {
	qp := Pg.NewQueryParams(rm.Posts, &TM_MASTER)
	Search_ByQueryParams(qp)
	qp.ToMap(rm.Ajax)
}

// 2017-06-04 Haries
func One_ByID(id int64) M.SX {
	ram_key := ZT(I.ToS(id))
	query := ram_key + `
SELECT ` + SELECT + `
FROM ` + TABLE + ` x1
WHERE x1.id = ` + ZI(id)
	return PG.CQFirstMap(TABLE, ram_key, query)
}

// 2017-06-04 Haries
func API_Superadmin_Form(rm *W.RequestModel) {
	rm.Ajax.SX = One_ByID(S.ToI(rm.Id))
}

// 2017-06-04 Haries
func API_Superadmin_SaveDeleteRestore(rm *W.RequestModel) {
	PG.DoTransaction(func(tx *Pg.Tx) string {
		dm := Pg.NewNonDataRow(tx, TABLE, rm)
		dm.SetNonData(`full_name`)
		dm.SetNonData(`email`)
		dm.SetNonData(`group_id`)
		dm.SetNonData(`note`)
		dm.SetNonData(`phone`)
		dm.SetNonData(`verified`)
		dm.UpsertRow()
		if !rm.Ajax.HasError() {
			dm.WipeUnwipe(rm.Action)
		}
		return rm.Ajax.LastError()
	})
}
