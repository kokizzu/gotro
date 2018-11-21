package Pg

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
	Row       M.SX
	NonData   M.SX
	Posts     *W.Posts
	Ajax      W.Ajax
	ReqModel  *W.RequestModel
	Table     string
	Id        int64
	Tx        *Tx
	DbActor   int64
	Log       string
	UniqueId  string // set when you want to update it
	IsDeleted bool
}

// convert Row to JSON string
func (mp *Row) ToJson() string {
	return M.ToJson(mp.Row)
}

// fetch model to be edited
func NewRow(tx *Tx, table string, rm *W.RequestModel) *Row {
	id := S.ToI(rm.Id)
	data, uniq, is_deleted := tx.DataJsonMapUniqAndIsDeleted_ById(table, id)
	if rm.Ajax.SX == nil {
		rm.Ajax = W.NewAjax()
	}
	return &Row{data, M.SX{}, rm.Posts, rm.Ajax, rm, table, id, tx, S.ToI(rm.DbActor), ``, uniq, is_deleted}
}

// create empty model
func NewNonDataRow(tx *Tx, table string, rm *W.RequestModel) *Row {
	id := S.ToI(rm.Id)
	if rm.Ajax.SX == nil {
		rm.Ajax = W.NewAjax()
	}
	return &Row{M.SX{}, M.SX{}, rm.Posts, rm.Ajax, rm, table, id, tx, S.ToI(rm.DbActor), ``, ``, false}
}

// fetch model to be edited from unique
func NewRowUniq(tx *Tx, table string, unique_id string, rm *W.RequestModel) *Row {
	data, id, is_deleted := tx.DataJsonMapIdAndIsDeleted_ByUniq(table, unique_id)
	new_uniq := unique_id
	if id == 0 {
		new_uniq = ``
	}
	if rm.Ajax.SX == nil {
		rm.Ajax = W.NewAjax()
	}
	res := &Row{data, M.SX{}, rm.Posts, rm.Ajax, rm, table, id, tx, S.ToI(rm.DbActor), ``, new_uniq, is_deleted}
	if id == 0 {
		res.Set_UniqueId(unique_id)
	}
	return res
}

// insert row
func (mp *Row) InsertRow() int64 {
	if mp.Ajax.HasError() {
		// ignore saving
		mp.Ajax.Info(`no record inserted..`)
		mp.Ajax.Set(`id`, mp.Id)
		return mp.Id
	}
	params := M.SX{}
	if len(mp.Row) > 0 {
		params[`data`] = mp.ToJson()
	}
	if mp.UniqueId != `` {
		params[`unique_id`] = mp.UniqueId
	}
	for k, v := range mp.NonData {
		params[k] = v
	}
	new_id := mp.Tx.DoInsert(mp.DbActor, mp.Table, params)
	label := mp.Table + `'s row ID:` + I.ToS(new_id)
	if new_id < 1 {
		mp.Ajax.Error(`Failed saving ` + label)
		L.Describe(mp.DbActor, mp.Row, mp.Id, mp.Log, mp.Table, mp.UniqueId)
		mp.Ajax.Set(`id`, new_id)
		return new_id
	}
	mp.Id = new_id
	mp.Ajax.Info(`Created new ` + label + " with: \n" + mp.Log)
	mp.Ajax.Set(`id`, I.ToS(new_id))
	return new_id
}

// update row
func (mp *Row) UpdateRow() int64 {
	if mp.Ajax.HasError() {
		// ignore saving
		mp.Ajax.Info(`no record updated..`)
		mp.Ajax.Set(`id`, mp.Id)
		return mp.Id
	}
	label := mp.Table + `'s row ID:` + I.ToS(mp.Id)
	if mp.Log == `` {
		mp.Ajax.Info(`No changes detected ` + label)
		//L.Describe(mp)
		mp.Ajax.Set(`id`, mp.Id)
		return mp.Id
	}
	params := M.SX{}
	if len(mp.Row) > 0 {
		params[`data`] = mp.ToJson()
	}
	if mp.UniqueId != `` {
		params[`unique_id`] = mp.UniqueId
	}
	for k, v := range mp.NonData {
		params[k] = v
	}
	update_count := mp.Tx.DoUpdate(mp.DbActor, mp.Table, mp.Id, params)
	if update_count < 1 {
		mp.Ajax.Error(`Failed saving ` + label)
		L.Describe(mp.Table, mp.Id, mp.UniqueId, mp.Log, mp.Row)
		mp.Ajax.Set(`_updated`, update_count)
		return update_count
	}
	mp.Ajax.Info(`Updated ` + label + " with: \n" + mp.Log)
	mp.Ajax.Set(`_updated`, I.ToS(update_count))
	return update_count
}

// insert or update row, if uniq ada
func (mp *Row) UpsertRow() int64 {
	if mp.Ajax.HasError() {
		// ignore saving
		mp.Ajax.Info(`no record upserted..`)
		mp.Ajax.Set(`id`, mp.Id)
		return mp.Id
	}
	new_rec := mp.Id == 0
	label := mp.Table + `'s row ID:`
	if !new_rec {
		label += I.ToS(mp.Id)
		if mp.Log == `` {
			mp.Ajax.Info(`No changes detected ` + label)
			//L.Describe(mp)
			mp.Ajax.Set(`id`, mp.Id)
			return mp.Id
		}
	}
	params := M.SX{}
	if len(mp.Row) > 0 {
		params[`data`] = mp.ToJson()
	}
	if mp.Id > 0 {
		params[`id`] = mp.Id
	}
	if mp.UniqueId != `` {
		params[`unique_id`] = mp.UniqueId
	}
	for k, v := range mp.NonData {
		params[k] = v
	}
	new_id := mp.Tx.DoUpsert(mp.DbActor, mp.Table, params)
	if new_rec {
		label += I.ToS(new_id)
	}
	if new_id < 1 {
		mp.Ajax.Error(`Failed saving ` + label)
		L.Describe(mp)
		mp.Ajax.Set(`id`, new_id)
		return new_id
	} else {
		mp.Id = new_id
	}
	if new_rec && mp.Log == `` {
		mp.Ajax.Info(`Saved new ` + label + ` with empty data`)
		mp.Ajax.Set(`id`, new_id)
		return new_id
	}
	if new_rec {
		mp.Ajax.Info(`Created new ` + label + " with: " + S.WebBR + mp.Log)
	} else {
		mp.Ajax.Info(`Updated ` + label + " with: " + S.WebBR + mp.Log)
	}
	mp.Ajax.Set(`id`, I.ToS(new_id))
	return new_id
}

// insert or update row, insert if not exists even when uinque_id exists (error)
func (mp *Row) IndateRow() int64 {
	if mp.Ajax.HasError() {
		// ignore saving
		mp.Ajax.Info(`no record upserted..`)
		return mp.Id
	}
	new_rec := mp.Id == 0
	label := mp.Table + `'s row ID:`
	if !new_rec {
		label += I.ToS(mp.Id)
		if mp.Log == `` {
			mp.Ajax.Info(`No changes detected ` + label)
			//L.Describe(mp)
			mp.Ajax.Set(`id`, mp.Id)
			return mp.Id
		}
	}
	params := M.SX{}
	if len(mp.Row) > 0 {
		params[`data`] = mp.ToJson()
	}
	if mp.Id > 0 {
		params[`id`] = mp.Id
	}
	if mp.UniqueId != `` {
		params[`unique_id`] = mp.UniqueId
	}
	for k, v := range mp.NonData {
		params[k] = v
	}
	new_id := int64(0)
	if mp.Id == 0 {
		new_id = mp.Tx.DoForcedInsert(mp.DbActor, mp.Table, params)
	} else {
		new_id = mp.Tx.DoUpsert(mp.DbActor, mp.Table, params)
	}
	if new_rec {
		label += I.ToS(new_id)
	}
	if new_id < 1 {
		mp.Ajax.Error(`Failed saving ` + label)
		L.Describe(mp)
		mp.Ajax.Set(`id`, new_id)
		return new_id
	} else {
		mp.Id = new_id
	}
	if new_rec && mp.Log == `` {
		mp.Ajax.Info(`Saved new ` + label + ` with empty data`)
		mp.Ajax.Set(`id`, new_id)
		return new_id
	}
	if new_rec {
		mp.Ajax.Info(`Created new ` + label + " with: " + S.WebBR + mp.Log)
	} else {
		mp.Ajax.Info(`Updated ` + label + " with: " + S.WebBR + mp.Log)
	}
	mp.Ajax.Set(`id`, I.ToS(new_id))
	return new_id
}

// log the changes
func (mp *Row) LogIt(key string, val interface{}) {
	key_label := ZZ(key)
	newv := X.ToS(val)
	new_label := ZZ(newv)
	if mp.Id == 0 {
		mp.Log += key_label + ` = ` + new_label + S.WebBR
	} else {
		oldv := X.ToS(mp.Row[key])
		if oldv != newv {
			mp.Log += key_label + ` from ` + ZZ(oldv) + ` to ` + new_label + S.WebBR
		}
	}
}

// log the changes
func (mp *Row) NotLogIt(key string, val interface{}) {
	key_label := ZZ(key)
	newv := X.ToS(val)
	new_label := ZZ(newv)
	oldv := X.ToS(mp.Row[key])
	if oldv != newv {
		mp.Log += `NOT overwriting ` + key_label + ` from ` + ZZ(oldv) + ` to ` + new_label + S.WebBR
	}
}

// log the changes
func (mp *Row) LogNonData(key string, val interface{}) {
	key_label := ZZ(key)
	newv := X.ToS(val)
	new_label := ZZ(newv)
	mp.Log += key_label + ` = ` + new_label + S.WebBR
}

// set unique id
func (mp *Row) Set_UniqueId(val string) {
	if val != `` {
		key_label := ZZ(`unique_id`)
		new_label := ZZ(val)
		if val != mp.UniqueId {
			mp.Log += key_label + ` = ` + new_label + S.WebBR
			mp.UniqueId = S.Trim(val)
		}
	}
	// TODO: unset unique id?
}

// undelete
func (mp *Row) Restore() {
	if mp.Id > 0 {
		if mp.Tx.DoRestore(mp.DbActor, mp.Table, mp.Id) {
			mp.Log += "record restored" + S.WebBR
		}
	}
}

// permanent delete (usage warning!)
func (mp *Row) PermanentErase() {
	if mp.Id > 0 {
		af, _ := mp.Tx.DoExec(`DELETE FROM ` + mp.Table + ` WHERE id = ` + I.ToS(mp.Id)).RowsAffected()
		if af > 0 {
			mp.Log += "record " + I.ToS(af) + " permanently erased" + S.WebBR
		}
	}
}

// delete
func (mp *Row) Delete() {
	if mp.Id > 0 {
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
		if val[0] != '+' || len(part) != 2 || trim != val {
			mp.Ajax.Error(`Invalid international phone format (+xx xxxxxx): ` + val)
			return X.ToS(mp.Row[key])
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
func (mp *Row) Get_Id() int64 {
	return mp.Id
}

// get unique id
func (mp *Row) Get_UniqueId() string {
	return mp.UniqueId
}

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
		ival := int64(0)
		if S.ToLower(val) == `now` {
			// special case for date
			ival = T.Epoch()
		} else {
			ival = mp.Posts.GetInt(key)
		}
		mp.LogIt(key, ival)
		mp.Row[key] = ival
	}
	return X.ToI(mp.Row[key])
}

// set float64 from Posts to Row
func (mp *Row) SetFloat(key string) float64 {
	val := mp.Posts.GetStr(key)
	if val != `` {
		fval := float64(0)
		if S.ToLower(val) == `now` {
			// special case for date
			fval = float64(T.Epoch())
		} else if S.Trim(val) == `` {
			mp.Unset(key)
			return 0
		} else {
			fval = mp.Posts.GetFloat(key)
		}
		mp.LogIt(key, fval)
		mp.Row[key] = fval
	}
	return X.ToF(mp.Row[key])
}

// set bool from Posts to Row
func (mp *Row) SetBool(key string) bool {
	val := mp.Posts.GetStr(key)
	if val != `` {
		mp.LogIt(key, val)
		mp.Row[key] = val == `true`
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

// set Row from json int[]
func (mp *Row) SetJsonIntArr(key string) []int64 {
	str := S.Trim(mp.Posts.GetStr(key))
	if str == `` || str[0] != '[' {
		return []int64{}
	}
	val := mp.Posts.GetJsonIntArr(key)
	mp.LogIt(key, str)
	mp.Row[key] = val
	return val
}

// set Row from json string[]
func (mp *Row) SetJsonStrArr(key string) []string {
	str := S.Trim(mp.Posts.GetStr(key))
	if str == `` || str[0] != '[' {
		return []string{}
	}
	val := mp.Posts.GetJsonStrArr(key)
	mp.LogIt(key, str)
	mp.Row[key] = val
	return val
}

// set Row from json object[]
func (mp *Row) SetJsonObjArr(key string) []map[string]interface{} {
	str := S.Trim(mp.Posts.GetStr(key))
	if str == `` || str[0] != '[' {
		return []map[string]interface{}{}
	}
	val := mp.Posts.GetJsonObjArr(key)
	mp.LogIt(key, str)
	mp.Row[key] = val
	return val
}

// set Row from json {object}
func (mp *Row) SetJsonMap(key string) M.SX {
	str := S.Trim(mp.Posts.GetStr(key))
	if str == `` || str[0] != '{' {
		return M.SX{}
	}
	val := mp.Posts.GetJsonMap(key)
	mp.LogIt(key, str)
	mp.Row[key] = val
	return val
}

// set Row value only if still empty, or never being saved
func (mp *Row) SetValOnce(key string, val interface{}) interface{} {
	if _, ok := mp.Row[key]; mp.Id != 0 && ok {
		mp.NotLogIt(key, val)
		return mp.Row[key]
	}
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
		emails = S.XSS(emails)
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
				ok = false
				return
			}
			if S.EndsWith(mail, OFFICE_MAIL_SUFFIX) {
				orig[`office_mail`] = mail
			} else if S.Contains(mail, `@gmail.`) {
				orig[`gmail`] = mail
			} else if S.Contains(mail, `@yahoo.`) ||
				S.Contains(mail, `@ymail.`) ||
				S.Contains(mail, `@rocketmail.`) {
				orig[`yahoo`] = mail
			} else if S.Contains(mail, `@outlook.`) ||
				S.Contains(mail, `@hotmail.`) ||
				S.Contains(mail, `@live.`) ||
				S.Contains(mail, `@windowslive.`) {
				orig[`outlook`] = mail
			} else {
				orig[`email`] = mail
			}
		}
		for _, key := range []string{`office_mail`, `gmail`, `yahoo`, `outlook`, `email`} {
			val := orig[key]
			if key != `openid_email` && key != `openid_mail` {
				mp.SetVal(key, val)
			}
			if val == `` {
				continue
			}
			lkey, rkey, id_str := Z(key), Z(val), ZI(mp.Id)
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

// set Row type with check
func (mp *Row) SetType(val string) string {
	// error check
	if mp.Id > 0 {
		old_type := mp.Row.GetStr(`type`)
		if !(old_type == `` || old_type == val) {
			return mp.Ajax.Error(`Invalid record type, trying to overwrite a ` + Z(old_type) + ` record to ` + Z(val) + `, this should not be happened, please contact Site Administrator`)
		}
	}
	mp.SetVal(`type`, val)
	return ``
}

// set raw data
func (mp *Row) SetNonDataVal(key string, val interface{}) {
	if val != nil {
		mp.NonData.Set(key, val)
		mp.LogNonData(key, val)
	}
}

func (mp *Row) SetNonData(key string) {
	nv := mp.ReqModel.Posts.GetStr(key)
	if nv != `` {
		nv = S.XSS(nv)
		mp.NonData.Set(key, nv)
		mp.LogNonData(key, nv)
	}
}
