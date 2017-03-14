package Sc

import (
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/T"
	"github.com/kokizzu/gotro/W"
	"github.com/kokizzu/gotro/X"
)

var OFFICE_MAIL_SUFFIX string

func InitOfficeMail(suffix string) {
	OFFICE_MAIL_SUFFIX = suffix
}

// primary table model
type Row struct {
	Row      M.SX
	Posts    *W.Posts
	Ajax     W.Ajax
	ReqModel *W.RequestModel
	Table    string
	Id       string
	Tx       *Tx
	DbActor  string
	Log      string
}

// convert Row to JSON string
func (mp *Row) ToJson() string {
	return M.ToJson(mp.Row)
}

// fetch model to be edited
func NewRow(tx *Tx, table string, rm *W.RequestModel) *Row {
	data := tx.DataJsonMap(table, rm.Id)
	return &Row{data, rm.Posts, rm.Ajax, rm, table, rm.Id, tx, rm.DbActor, ``}
}

//// fetch model to be edited from unique
//func NewRowUniq(tx *Tx, table string, unique_id string, rm *W.RequestModel) *Row {
//	// not implemented: no unique id on scylladb
//}

// insert row
func (mp *Row) InsertRow() bool {
	if mp.Ajax.HasError() {
		// ignore saving
		mp.Ajax.Info(`no record inserted..`)
		mp.Ajax.Set(`id`, ``)
		return false
	}
	label := mp.Table + `'s row ID:`
	if !mp.Tx.DoInsert(mp.DbActor, mp.Table, mp.Row) {
		mp.Ajax.Error(`Failed insert ` + label)
		L.Describe(mp.DbActor, mp.Row, mp.Id, mp.Log, mp.Table)
		return false
	}
	mp.Id = mp.Row.GetStr(`id`)
	label += mp.Id
	mp.Ajax.Info(`Created new ` + label + " with: \n" + mp.Log)
	mp.Ajax.Set(`id`, mp.Id)
	return true
}

// update row
func (mp *Row) UpdateRow() bool {
	mp.Ajax.Set(`id`, mp.Id)
	if mp.Ajax.HasError() {
		// ignore saving
		mp.Ajax.Info(`no record updated..`)
		return false
	}
	label := mp.Table + `'s row ID:` + mp.Id
	if mp.Log == `` {
		mp.Ajax.Info(`No changes detected ` + label)
		return true
	}
	if !mp.Tx.DoUpdate(mp.DbActor, mp.Table, mp.Row) {
		mp.Ajax.Error(`Failed update ` + label)
		L.Describe(mp.DbActor, mp.Row, mp.Id, mp.Log, mp.Table)
		return false
	}
	mp.Ajax.Info(`Updated ` + label + " with: \n" + mp.Log)
	return true
}

// insert or update row, if uniq ada
func (mp *Row) UpsertRow() bool {
	mp.Ajax.Set(`id`, mp.Id)
	if mp.Ajax.HasError() {
		// ignore saving
		mp.Ajax.Info(`no record upserted..`)
		return false
	}
	new_rec := mp.Id == ``
	label := mp.Table + `'s row ID:` + mp.Id
	if mp.Log == `` {
		mp.Ajax.Info(`No fields/changes detected ` + label)
		return true
	}
	if !mp.Tx.DoUpsert(mp.DbActor, mp.Table, mp.Row) {
		mp.Ajax.Error(`Failed upsert ` + label)
		L.Describe(mp.DbActor, mp.Row, mp.Id, mp.Log, mp.Table)
		return false
	}
	if new_rec {
		mp.Id = mp.Row.GetStr(`id`)
		label += mp.Id
		mp.Ajax.Info(`Created new ` + label + " with: " + S.WebBR + mp.Log)
	} else {
		mp.Ajax.Info(`Updated ` + label + " with: " + S.WebBR + mp.Log)
	}
	mp.Ajax.Set(`id`, mp.Id)
	return true
}

// insert or update row, insert if not exists even when uinque_id exists (error)
func (mp *Row) IndateRow() bool {
	mp.Ajax.Set(`id`, mp.Id)
	if mp.Ajax.HasError() {
		// ignore saving
		mp.Ajax.Info(`no record indated..`)
		return false
	}
	new_rec := mp.Id == ``
	label := mp.Table + `'s row ID:` + mp.Id
	if mp.Log == `` {
		mp.Ajax.Info(`No fields/changes detected ` + label)
		return true
	}
	var ok bool
	if new_rec {
		ok = mp.Tx.DoForcedInsert(mp.DbActor, mp.Table, mp.Row)
	} else {
		ok = mp.Tx.DoUpsert(mp.DbActor, mp.Table, mp.Row)
	}
	if !ok {
		mp.Ajax.Error(`Failed indate ` + label)
		L.Describe(mp.DbActor, mp.Row, mp.Id, mp.Log, mp.Table)
		return false
	}
	if new_rec {
		mp.Id = mp.Row.GetStr(`id`)
		label += mp.Id
		mp.Ajax.Info(`Created new ` + label + " with: " + S.WebBR + mp.Log)
	} else {
		mp.Ajax.Info(`Updated ` + label + " with: " + S.WebBR + mp.Log)
	}
	mp.Ajax.Set(`id`, mp.Id)
	return true
}

// log the changes
func (mp *Row) LogIt(key string, val interface{}) {
	key_label := ZZ(key)
	newv := X.ToS(val)
	new_label := ZZ(newv)
	if mp.Id == `` {
		mp.Log += key_label + ` = ` + new_label + S.WebBR
	} else {
		oldv := X.ToS(mp.Row[key])
		if oldv != newv {
			mp.Log += key_label + `  from ` + ZZ(oldv) + ` to ` + new_label + S.WebBR
		}
	}
}

//// set unique id, equal to id
//func (mp *Row) Set_UniqueId(val string) {
//	// not implemented: since there are no unique id on scylladb
//}

// undelete
func (mp *Row) Restore() {
	if mp.Id != `` {
		if mp.Tx.DoRestore(mp.DbActor, mp.Table, mp.Id) {
			mp.Log += "record restored" + S.WebBR
		}
	}
}

// delete
func (mp *Row) Delete() {
	if mp.Id != `` {
		if mp.Tx.DoDelete(mp.DbActor, mp.Table, mp.Id) {
			mp.Log += "record deleted" + S.WebBR
		}
	}
}

// delete or restore
func (mp *Row) WipeUnwipe(a string) {
	mp.Tx.DoWipeUnwipe(a, mp.DbActor, mp.Table, mp.Id)
}

// set string
func (mp *Row) SetStr(key string) string {
	val := mp.Posts.GetStr(key)
	if val != `` {
		val = S.XSS(val)
		mp.LogIt(key, val)
		mp.Row[key] = val
	}
	return X.ToS(mp.Row[key])
}

// set string strip prefix and suffix from and letters
func (mp *Row) SetStrPhone(key string) string {
	val := mp.Posts.GetStr(key)
	val = S.ValidatePhone(val)
	if val != `` {
		mp.LogIt(key, val)
		mp.Row[key] = val
	}
	return X.ToS(mp.Row[key])
}

// set international phone, format: +xx xxxxxx
func (mp *Row) SetIntlPhone(key string) string {
	val := mp.Posts.GetStr(key)
	val = S.Trim(val)
	if val != `` {
		part := S.Split(val, ` `)
		trim := S.ValidatePhone(val)
		if val[0] != '+' || len(part) != 2 || len(trim)+1 != len(val) {
			mp.Ajax.Error(`Invalid international phone format (+xx xxxxxx): ` + val)
			return ``
		}
		mp.LogIt(key, val)
		mp.Row[key] = val
	}
	return X.ToS(mp.Row[key])
}

// get string from Row
func (mp *Row) GetStr(key string) string {
	return X.ToS(mp.Row[key])
}

// get boolean from Row
func (mp *Row) GetBool(key string) bool {
	return X.ToBool(mp.Row[key])
}

// get int64 from Row
func (mp *Row) GetInt(key string) int64 {
	return X.ToI(mp.Row[key])
}

// get []interface{} from Row
func (mp *Row) GetAX(key string) []interface{} {
	return X.ToArr(mp.Row[key])
}

// get float64 from Row
func (mp *Row) GetFloat(key string) float64 {
	return X.ToF(mp.Row[key])
}

// get id
func (mp *Row) Get_Id() string {
	return mp.Id
}

//// get unique id
//func (mp *Row) Get_UniqueId() string {
//	// not implemented: no unique constraint on scylladb
//}

// set time from Posts to Row
// unset when string is whitespace
func (mp *Row) SetUnsetClock(key string) string {
	val := mp.Posts.GetStr(key)
	return mp.SetUnsetValClock(key, val)
}

// set time hh:mm
func (mp *Row) SetUnsetValClock(key string, val string) string {
	//	L.Describe(key, val)
	if val != `` {
		val = S.Trim(val)
		if val == `` {
			mp.Unset(key)
			return ``
		}
		val = S.Replace(val, `.`, `:`)
		hh_mm := S.Split(val, `:`)
		if len(hh_mm) < 2 {
			mp.Ajax.Error(`invalid format for '` + key + `': ` + val + `, time format must a HH:MM`)
			return ``
		}
		// check hours
		hh := S.ToI(hh_mm[0])
		if hh < 0 || hh > 23 {
			mp.Ajax.Error(`invalid hour for '` + key + `': ` + val)
			return ``
		}
		// check minutes
		mm := S.ToI(hh_mm[1])
		if mm < 0 || mm > 59 {
			mp.Ajax.Error(`invalid minute for '` + key + `': ` + val)
			return ``
		}
		val = S.PadLeft(I.ToS(hh), `0`, 2) + `:` + S.PadLeft(I.ToS(mm), `0`, 2)
		mp.LogIt(key, val)
		mp.Row[key] = val
	}
	return X.ToS(mp.Row[key])
}

// set int64 from Posts to Row
func (mp *Row) SetInt(key string) int64 {
	val := mp.Posts.GetStr(key)
	if val != `` {
		mp.LogIt(key, val)
		val := mp.Posts.GetInt(key)
		mp.Row[key] = val
	}
	return X.ToI(mp.Row[key])
}

// set float64 from Posts to Row
func (mp *Row) SetFloat(key string) float64 {
	val := mp.Posts.GetStr(key)
	if val != `` {
		mp.LogIt(key, val)
		val := mp.Posts.GetFloat(key)
		mp.Row[key] = val
	}
	return X.ToF(mp.Row[key])
}

// set bool from Posts to Row
func (mp *Row) SetBool(key string) bool {
	val := mp.Posts.GetStr(key)
	if val != `` {
		mp.LogIt(key, val)
		mp.Row[key] = (val == `true`)
	}
	return X.ToBool(mp.Row[key])
}

// unset Row key
func (mp *Row) Unset(key string) {
	oldv, exists := mp.Row[key]
	if exists {
		mp.Log += ZZ(key) + ` removed, previously ` + ZZ(X.ToS(oldv)) + S.WebBR
		delete(mp.Row, key)
	}
}

// set unset int, returns 0 when unsetted
func (mp *Row) SetUnsetIntVal(key string, val int64) int64 {
	if val <= 0 {
		mp.Unset(key)
		return 0
	}
	mp.SetVal(key, val)
	return val
}

// set user password, skip logging
func (mp *Row) Set_UserPassword(pass string) {
	if pass != `` {
		mp.Log += ZZ(`password`) + ` changed` + S.WebBR
		mp.Row[`password`] = S.HashPassword(pass)
		mp.Row[`last_reset_password_at`] = T.Epoch()
	}
}

// check password
func (mp *Row) Check_UserPassword(pass string) bool {
	return mp.Row.GetStr(`password`) == S.HashPassword(pass)
}

// set Row value
func (mp *Row) SetVal(key string, val interface{}) interface{} {
	switch v := val.(type) {
	case string:
		val = S.XSS(v)
	}
	mp.LogIt(key, val)
	mp.Row[key] = val
	return val
}

func (mp *Row) IsChanged() bool {
	return mp.Log != ``
}

// set val without XSS filtering
func (mp *Row) SetValNoXSS(key string, val interface{}) interface{} {
	mp.LogIt(key, val)
	mp.Row[key] = val
	return val
}

// set Row value if ok
func (mp *Row) SetValIf(ok bool, key string, val interface{}) {
	if ok {
		mp.SetVal(key, val)
	}
}

// set Row value from string
func (mp *Row) SetValStr(key, val string) {
	if val != `` {
		val = S.XSS(val)
		mp.LogIt(key, val)
		mp.Row[key] = val
	}
}

// set Row by current date epoch as float
func (mp *Row) SetValEpoch(key string) float64 {
	val := T.Epoch()
	mp.LogIt(key, val)
	mp.Row[key] = val
	return float64(val)
}

// set Row by current date epoch as float
func (mp *Row) SetValEpochOnce(key string) float64 {
	old_val := X.ToF(mp.Row[key])
	if old_val > 0 {
		return old_val
	}
	val := T.Epoch()
	mp.LogIt(key, val)
	mp.Row[key] = val
	return float64(val)
}

// set Row office_mail, gmail, yahoo and email
func (mp *Row) Set_UserEmails(emails string) (ok bool) {
	return mp.Set_UserEmails_ByTable(emails, `users`)
}

func (mp *Row) Set_UserEmails_ByTable(emails string, table string) (ok bool) {
	ok = true
	if emails != `` {
		emails = S.Z(emails)
		emails = S.ToLower(emails)
		orig := M.SS{}
		orig[`office_mail`] = ``
		orig[`gmail`] = ``
		orig[`yahoo`] = ``
		orig[`email`] = ``
		mails := S.Split(emails, `,`)
		for _, mail := range mails {
			orig_mail := S.Trim(mail)
			mail = S.ValidateEmail(orig_mail)
			if mail == `` {
				L.Describe(`invalid e-mail`, orig_mail)
				continue
			}
			if S.EndsWith(mail, OFFICE_MAIL_SUFFIX) {
				orig[`office_mail`] = mail
			} else if S.Contains(mail, `@gmail.`) {
				orig[`gmail`] = mail
			} else if S.Contains(mail, `@yahoo.`) ||
				S.Contains(mail, `@ymail.`) ||
				S.Contains(mail, `@rocketmail.`) {
				orig[`yahoo`] = mail
			} else {
				orig[`email`] = mail
			}
		}
		for _, key := range []string{`office_mail`, `gmail`, `yahoo`, `email`} {
			val := orig[key]
			mp.SetVal(key, val)
			if val == `` {
				continue
			}
			lkey, rkey, id_str := Z(key), Z(val), Z(mp.Id)
			query := `SELECT COALESCE((SELECT id FROM ` + ZZ(table) + ` WHERE data->>` + lkey + ` = ` + rkey + ` AND id <> ` + id_str + ` AND is_deleted = false LIMIT 1),0)`
			dup_id := mp.Tx.QInt(query)
			if dup_id == 0 {
				continue
			}
			msg := `The ` + lkey + ` is being used by another user account: ` + rkey + `, if you think this should not be happened, please send a bug report to the WebMaster.` // ` = ` + ZI(dup_id) + ` <> ` + id_str
			L.Describe(msg)
			mp.Ajax.Error(msg)
			ok = false
		}
	}
	return
}
