package mUsers

import (
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D/Pg"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/T"
	"github.com/kokizzu/gotro/W"
	"github.com/kokizzu/gotro/W/example-complex/model"
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
var ZT func(...string) string
var PG_W, PG_R *Pg.RDBMS

func init() {
	Z = S.Z
	ZB = S.ZB
	ZZ = S.ZZ
	ZJ = S.ZJ
	ZI = S.ZI
	ZLIKE = S.ZLIKE
	ZT = S.ZT
	PG_W = model.PG_W
	PG_R = model.PG_R
	TM_MASTER = Pg.TableModel{
		CacheName: TABLE + `_USERS_MASTER`,
		Fields: []Pg.FieldModel{
			{Key: `id`},
			{Key: `is_deleted`},
			{Key: `modified_at`},
			{Label: `E-Mail(s)`, Key: `emails`, CustomQuery: `emails_join(data)`, Type: `emails`, FormTooltip: `separate with comma`},
			{Label: `Phone`, Key: `phone`, Type: `phone`, FormHide: true},
			{Label: `Full Name`, Key: `full_name`},
		},
	}
	SELECT = TM_MASTER.Select()
}
func One_ByID(id string) M.SX {
	ram_key := ZT(id)
	query := ram_key + `
SELECT ` + SELECT + `
FROM ` + TABLE + ` x1
WHERE x1.id::TEXT = ` + Z(id)
	return PG_R.CQFirstMap(TABLE, ram_key, query)
}

func Search_ByQueryParams(qp *Pg.QueryParams) {
	qp.RamKey = ZT(qp.Term)
	if qp.Term != `` {
		qp.Where += ` AND (x1.data->>'name') LIKE ` + ZLIKE(qp.Term)
	}
	qp.From = `FROM ` + TABLE + ` x1`
	qp.OrderBy = `x1.id`
	qp.Select = SELECT
	qp.SearchQuery_ByConn(PG_W)

}

/* accessed through: {"order":["-col1","+col2"],"filter":{"is_deleted":false,"created_at":">isodate"},"limit":10,"offset":5}
this will retrieve record 6-15 order by col1 descending, col2 ascending, filtered by is_deleted=false and created_at > isodate
*/

func All_ByStartID_ByLimit_IsAsc_IsIncl(id string, limit int64, is_asc, is_incl bool) A.MSX {
	sign := S.IfElse(is_asc, `>`, `<`) + S.If(is_incl, `=`)
	ram_key := ZT(id, I.ToS(limit), sign)
	where := ``
	if id != `` {
		where = `AND x1.id ` + sign + Z(id)
	}
	query := ram_key + `
SELECT ` + SELECT + `
FROM ` + TABLE + ` x1
WHERE x1.is_deleted = false
 ` + where + `
ORDER BY x1.id ` + S.If(!is_asc, `DESC`) + `
LIMIT ` + I.ToS(limit)
	return PG_R.CQMapArray(TABLE, ram_key, query)
}

func API_Backoffice_Form(rm *W.RequestModel) {
	rm.Ajax.SX = One_ByID(rm.Id)
}

func API_Backoffice_SaveDeleteRestore(rm *W.RequestModel) {
	PG_W.DoTransaction(func(tx *Pg.Tx) string {
		dm := Pg.NewRow(tx, TABLE, rm) // NewPostlessData
		emails := rm.Posts.GetStr(`emails`)
		// rm is the requestModel, values provided by http req
		dm.Set_UserEmails(emails)
		// dm is the dataModel, row we want to update
		// we can call dm.Get* to retrieve old record values
		dm.SetStr(`full_name`)
		if !rm.Ajax.HasError() {
			dm.WipeUnwipe(rm.Action)
		}
		return rm.Ajax.LastError()
	})
}

func API_Backoffice_FormLimit(rm *W.RequestModel) {
	id := rm.Posts.GetStr(`id`)
	limit := rm.Posts.GetInt(`limit`)
	is_asc := rm.Posts.GetBool(`asc`)
	is_incl := rm.Posts.GetBool(`incl`)
	result := All_ByStartID_ByLimit_IsAsc_IsIncl(id, limit, is_asc, is_incl)
	rm.Ajax.Set(`result`, result)
}

func API_Backoffice_Search(rm *W.RequestModel) {
	qp := Pg.NewQueryParams(rm.Posts, &TM_MASTER)
	Search_ByQueryParams(qp)
	qp.ToMap(rm.Ajax)
}

// 2016-01-20 Prayogo
func FindID_ByIdent_ByPass(ident, pass string) int64 {
	ident = S.Trim(ident)
	if ident == `` || pass == `` {
		return 0
	}
	ident = Z(ident)
	hash := S.HashPassword(pass)
	ram_key := ZT(ident, hash)
	query := ram_key + `
SELECT COALESCE((
	SELECT id
	FROM accounts
	WHERE is_deleted = false
		AND (
			data->>'password' = ` + Z(hash) + `
			AND ( data->>'office_mail' = ` + ident + `
				OR data->>'gmail' = ` + ident + `
				OR data->>'yahoo' = ` + ident + `
				OR data->>'email' = ` + ident + `
				OR data->>'phone' = ` + ident + `
				OR unique_id = ` + ident + `
			)
		)
),0)`
	return PG_R.CQInt(TABLE, ram_key, query)
}

// 2016-07-26 Prayogo
func FindID_ByPhone(phone string) int64 {
	phone = Z(S.ToLower(phone))
	if phone == `''` {
		return 0
	}
	ram_key := ZT(phone)
	query := ` -- ` + ram_key + `
SELECT COALESCE((
	SELECT id
	FROM accounts
	WHERE is_deleted = false
		AND data->>'phone' = ` + phone + `
),0)`
	return PG_R.CQInt(TABLE, ram_key, query)
}

func FindID_ByEmail(email string) int64 {
	email = Z(S.ToLower(email))
	if email == `''` {
		return 0
	}
	ram_key := ZT(email)
	query := ` -- ` + ram_key + `
SELECT COALESCE((
	SELECT id
	FROM accounts
	WHERE is_deleted = false
		AND ( data->>'gmail' = ` + email + `
			OR data->>'email' = ` + email + `
			OR data->>'yahoo' = ` + email + `
			OR data->>'office_mail' = ` + email + `
		)
),0)`
	return PG_R.CQInt(TABLE, ram_key, query)
}

func UpdateLastLogin(id int64) {
	PG_W.DoTransaction(func(tx *Pg.Tx) string {
		query := ZT(I.ToS(id)) + `
UPDATE accounts SET data=JSONB_MERGE(data,'{"last_login":` + ZZ(T.EpochStr()) + `}'), updated_by=` + ZI(id) + ` WHERE id = ` + ZI(id)
		tx.DoExec(query)
		return ``
	})
}

func UpdateLastForgotPassword(id int64) {
	PG_W.DoTransaction(func(tx *Pg.Tx) string {
		query := ZT(I.ToS(id)) + `
UPDATE accounts SET data=JSONB_MERGE(data,'{"last_forgot_password":` + ZZ(T.EpochStr()) + `}'), updated_by=` + ZI(id) + ` WHERE id = ` + ZI(id)
		tx.DoExec(query)
		return ``
	})
}

func FindID_ByCompactName_ByEmail(ident, email string) int64 {
	ident = S.Trim(ident)
	email = S.Trim(email)
	if email == `` {
		return 0
	}
	ident = Z(ident)
	email = Z(S.ToLower(email))
	ram_key := ZT(ident, email)
	query := ` -- ` + ram_key + `
SELECT COALESCE((
	SELECT id
	FROM accounts
	WHERE is_deleted = false
	AND ( data->>'gmail' = ` + email + `
		OR data->>'email' = ` + email + `
		OR data->>'yahoo' = ` + email + `
		OR data->>'office_mail' = ` + email + `
	) AND alnum_str(data->>'full_name') = LOWER(` + ident + `)
),0)`
	return PG_R.CQInt(TABLE, ram_key, query)
}

func Name_Emails_ByID(id int64) (string, []string) {
	ram_key := ZT(I.ToS(id))
	query := ram_key + `
SELECT COALESCE(data->>'full_name','name_not_set')
	, emails_join(data) emails
FROM accounts
WHERE id = ` + ZI(id)
	name, mails := PG_R.QStrStr(query)
	emails := S.Split(mails, `,`)
	to := []string{}
	for _, email := range emails {
		if email == `` {
			continue
		}
		to = append(to, S.If(name != ``, ZZ(name))+` <`+email+`>`)
	}
	return name, to
}
