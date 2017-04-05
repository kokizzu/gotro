package Pg

import (
	"github.com/OneOfOne/cmap"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/C"
	"github.com/kokizzu/gotro/F"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/W"
	"github.com/kokizzu/gotro/X"
)

var SELECT_CACHE cmap.CMap
var FORM_CACHE cmap.CMap // cache form fields
var GRID_CACHE cmap.CMap // cache grid fields
var TYPE_CACHE cmap.CMap // cache key-type (M.SS)

func init() {
	SELECT_CACHE = cmap.New()
	FORM_CACHE = cmap.New()
	GRID_CACHE = cmap.New()
	TYPE_CACHE = cmap.New()
}

type FieldModel struct {
	Label       string // label in form, and grid
	HtmlLabel   string // label for both form and grid
	FormLabel   string // label in form
	GridLabel   string // label in grid
	FormTooltip string // placeholder in form
	Key         string // key in the table.data
	Type        string // for formatting sql, form, and grid: float2, integer, datetime
	SqlType     string // for overriding sql format: bigint, float
	FormType    string // for overriding form format
	GridType    string // for overriding grid format
	HtmlSubType string // for checkbox example: 'True':'False', for select: DS.DataSourceName
	Hide        bool   // for hiding in sql, form, and grid
	SqlHide     bool   // for overriding sql hide status
	HtmlHide    bool   // for hiding both form and grid
	FormHide    bool   // for overriding form hide status
	GridHide    bool   // for overriding grid hide status
	Required    bool
	CustomQuery string
	Default     string // default value
	Min         string // minimum value
	Max         string // maximum value
	SqlColPos   int    // position on SQL Select
}

func (field *FieldModel) SqlColumn() string {
	query := ``
	switch field.Key {
	case `id`, `is_deleted`, `unique_id`:
		query = `x1.` + field.Key
	case `modified_at`, `created_at`, `deleted_at`, `restored_at`, `updated_at`:
		field.CustomQuery = `EXTRACT(EPOCH FROM x1.` + field.Key + `)`
		fallthrough
	default:
		if field.CustomQuery != `` {
			query = field.CustomQuery
		} else {
			query = `x1.data->>` + Z(field.Key)
		}
		typ := S.IfEmpty(field.SqlType, field.Type)
		switch typ {
		case `float2`, `float`, `datetime`, `date`:
			query = `(` + query + `)::FLOAT`
		case `integer`, `bigint`:
			query = `(` + query + `)::BIGINT`
		case `bool`:
			query = `(` + query + ` = 'true')`
		}
	}
	return query
}

type TableModel struct {
	CacheName string
	Fields    []FieldModel
}

// generate select fields
func (tm *TableModel) Select() string {
	cache := SELECT_CACHE.Get(tm.CacheName)
	query_str, ok := cache.(string)
	if !ok {
		queries := []string{}
		pos := 1
		for idx, field := range tm.Fields {
			if field.Hide || field.SqlHide {
				continue
			}
			switch field.Key {
			case `id`:
				query := `x1.id::TEXT "id"`
				queries = append(queries, query)
				tm.Fields[idx].SqlColPos = pos
				pos += 1
				continue
			case `created_at`, `modified_at`, `updated_at`, `deleted_at`, `restored_at`:
				tm.Fields[idx].Label = S.ToTitle(S.Replace(field.Key, `_`, ` `))
				tm.Fields[idx].Type = `datetime`
			}
			query := field.SqlColumn() + ` ` + ZZ(field.Key)
			queries = append(queries, query)
			tm.Fields[idx].SqlColPos = pos
			pos += 1
		}
		query_str = A.StrJoin(queries, "\n, ")
		SELECT_CACHE.Set(tm.CacheName, query_str)
	} else if DEBUG {
		panic(`Warning: CacheName collision for ` + tm.CacheName)
	}
	return query_str
}

// generate form fields json
func (tm *TableModel) FormFields() A.MSX {
	cache := FORM_CACHE.Get(tm.CacheName)
	json_arr, ok := cache.(A.MSX)
	if !ok {
		json_arr = A.MSX{}
		for _, field := range tm.Fields {
			switch field.Key {
			case `id`, `is_deleted`, `modified_at`, `unique_id`:
				continue
			}
			if field.Hide || field.HtmlHide || field.FormHide {
				continue
			}
			json_obj := M.SX{
				`key`:     field.Key,
				`label`:   S.Coalesce(field.FormLabel, field.HtmlLabel, field.Label),
				`type`:    S.Coalesce(field.FormType, field.HtmlLabel, field.Type),
				`tooltip`: S.Coalesce(field.FormTooltip, field.FormLabel, field.HtmlLabel, field.Label),
			}
			if field.HtmlSubType != `` {
				json_obj[`sub_type`] = field.HtmlSubType
			}
			json_arr = append(json_arr, json_obj)
			if DEBUG {
				L.Print(`Creating FORM_CACHE.Select`, tm.CacheName)
			}
		}
		FORM_CACHE.Set(tm.CacheName, json_arr)
	}
	return json_arr
}

// generate grid fields json
func (tm *TableModel) GridFields() A.MSX {
	cache := GRID_CACHE.Get(tm.CacheName)
	json_arr, ok := cache.(A.MSX)
	if !ok {
		json_arr = A.MSX{}
		for _, field := range tm.Fields {
			switch field.Key {
			case `id`, `is_deleted`, `modified_at`:
				continue
			}
			if field.Hide || field.HtmlHide || field.GridHide {
				continue
			}
			json_obj := M.SX{
				`key`:   field.Key,
				`label`: S.Coalesce(field.GridLabel, field.HtmlLabel, field.Label),
				`type`:  S.IfEmpty(field.GridType, field.Type),
			}
			if field.HtmlSubType != `` {
				json_obj[`sub_type`] = field.HtmlSubType
			}
			json_arr = append(json_arr, json_obj)
			if DEBUG {
				L.Print(`Creating GRID_CACHE.Select`, tm.CacheName)
			}
		}
		GRID_CACHE.Set(tm.CacheName, json_arr)
	}
	return json_arr
}

// get type of a field by key
func (tm *TableModel) FieldModel_ByKey(name string) FieldModel {
	cache := TYPE_CACHE.Get(tm.CacheName)
	kv_map, ok := cache.(map[string]FieldModel)
	if !ok {
		kv_map = map[string]FieldModel{}
		for _, field := range tm.Fields {
			kv_map[field.Key] = field
		}
		TYPE_CACHE.Set(tm.CacheName, kv_map)
		if DEBUG {
			L.Print(`Creating TYPE_CACHE.Select`, tm.CacheName)
		}
	}
	return kv_map[name]
}

const ROWS_MAX_LIMIT = 200

// 2017-01-25 Prayogo
type QueryParams struct {
	Term      string
	Offset    int64
	Limit     int64
	Count     int64
	Order     []string
	Rows      A.MSX
	Filter    M.SX
	Model     *TableModel
	IsDefault bool

	Where   string
	Select  string
	RamKey  string
	From    string
	OrderBy string
}

func NewQueryParams(posts *W.Posts, model *TableModel) *QueryParams {
	if posts == nil {
		return &QueryParams{
			Term:      ``,
			Offset:    0,
			Limit:     10,
			Model:     model,
			Order:     []string{},
			IsDefault: true,
		}
	}
	return &QueryParams{
		Term:   posts.GetStr(`term`),
		Offset: posts.GetInt(`offset`),
		Limit:  posts.GetInt(`limit`),
		Filter: posts.GetJsonMap(`filter`),
		Order:  posts.GetJsonStrArr(`order`),
		Model:  model,
	}
}

func (qp *QueryParams) ToMap(ajax W.Ajax) {
	ajax.Set(`rows`, qp.Rows)
	ajax.Set(`count`, qp.Count)
	ajax.Set(`offset`, qp.Offset)
	ajax.Set(`limit`, qp.Limit)
	if qp.IsDefault {
		// for rendering html, mostly this required
		ajax.Set(`form_fields`, qp.Model.FormFields())
		ajax.Set(`grid_fields`, qp.Model.GridFields())
	}
}

func filterCriteriaSuffix_Numeral_ByPrefix(key, typ, str string) string {
	where_or_suffix := ``
	start_parse := int64(0)
	if C.IsDigit(str[0]) {
		where_or_suffix += `=`
	} else if len(str) > 1 && (str[0] == '<' || str[0] == '>') {
		start_parse = I.IfElse(str[1] == '=', 2, 1)
		where_or_suffix += str[:start_parse]
	} else {
		L.Print(`Ignoring integer/float2: `, key, str)
		return ``
	}
	if typ == `integer` {
		return `::BIGINT` + where_or_suffix + I.ToS(S.ToI(str[start_parse:]))
	} else {
		return `::FLOAT` + where_or_suffix + F.ToS(S.ToF(str[start_parse:]))
	}
}

func (qp *QueryParams) SearchQuery_ByConn(conn *RDBMS) {
	qp.RamKey += `:` + I.ToS(qp.Offset) + `:` + I.ToS(qp.Limit)
	if qp.Limit < 1 {
		qp.Limit = 1
	} else if qp.Limit > ROWS_MAX_LIMIT {
		qp.Limit = ROWS_MAX_LIMIT
	}
	for key, val := range qp.Filter {
		fm := qp.Model.FieldModel_ByKey(key)
		v_str := X.ToS(val)
		val_arr := S.Split(v_str, `|`)
		where_add := []string{}
		criteria := `(x1.data->>` + Z(key) + `)`
		if key == `is_deleted` {
			if v_str == `true` || v_str == `false` {
				qp.Where += ` AND x1.is_deleted = ` + v_str
			} else {
				L.Print(`Ignoring bool: `, key, v_str)
			}
			continue
		}
		if fm.Key == `` {
			L.Print(`Ignoring key: `, key)
			continue
		}
		if fm.GridType == `filled` {
			if v_str == `true` || v_str == `false` {
				criteria = ` AND COALESCE(` + criteria + `,'') ` + S.IfElse(v_str == `true`, `<>`, `=`) + ` ''`
				qp.Where += criteria
			} else {
				L.Print(`Ignoring bool: `, key, v_str)
			}
			continue
		}
		switch fm.Type {
		case `bool`:
			if len(val_arr) != 1 {
				L.Print(`Ignoring bool: `, val_arr)
				continue
			}
			if val_arr[0] == `true` {
				criteria += ` = 'true'`
			} else if val_arr[0] == `false` {
				criteria += ` <> 'true'`
			} else {
				L.Print(`Ignoring bool: `, key, val_arr[0])
				continue
			}
			where_add = append(where_add, criteria)
		case `integer`, `float2`, `date`, `datetime`:
			for _, str := range val_arr {
				str = S.Trim(str)
				if str == `` {
					continue
				}
				str2_arr := S.Split(str, ` `)
				where2_and := []string{}
				for _, str2 := range str2_arr {
					if str2 == `` {
						continue
					}
					criteria_suffix := filterCriteriaSuffix_Numeral_ByPrefix(key, fm.Type, str2)
					if criteria_suffix != `` {
						where2_and = append(where2_and, criteria+criteria_suffix)
					}
				}
				if len(where2_and) > 0 {
					where_add = append(where_add, `(`+A.StrJoin(where2_and, `) AND (`)+`)`)
				}
			}
		default:
			for _, str := range val_arr {
				str = S.Trim(str)
				if str == `` {
					continue
				}
				where_add = append(where_add, criteria+` ILIKE `+ZLIKE(str))
			}
		}
		qp.Where += ` AND ((` + A.StrJoin(where_add, `) OR (`) + `)) `
	}
	if len(qp.Order) > 0 {
		qp.OrderBy = ``
	}
	for _, order_key := range qp.Order {
		if S.Trim(order_key) == `` {
			continue
		}
		direction := order_key[0]
		order_key = order_key[1:]
		fm := qp.Model.FieldModel_ByKey(order_key)
		if fm.Key == `` {
			L.Print(`Ignoring key: ` + order_key)
			continue
		}
		if qp.OrderBy != `` {
			qp.OrderBy += `, `
		}
		if fm.SqlColPos > 0 {
			qp.OrderBy += I.ToStr(fm.SqlColPos)
		} else {
			qp.OrderBy += fm.SqlColumn()
		}
		if direction == '-' {
			qp.OrderBy += ` DESC`
		}
	}
	query_str := qp.From + `
WHERE 1=1
` + qp.Where
	query := ` -- ` + qp.RamKey + `_Count
SELECT COUNT(*)
` + query_str
	qp.Count = conn.QInt(query)
	if qp.Offset < 0 {
		qp.Offset = 0
	} else if qp.Offset >= qp.Count {
		qp.Offset = qp.Count / qp.Limit * qp.Limit
	}
	query = ` -- ` + qp.RamKey + `
SELECT ` + qp.Select + `
` + query_str + `
ORDER BY ` + S.IfEmpty(qp.OrderBy, `x1.id`) + `
LIMIT ` + I.ToS(qp.Limit) + S.If(qp.Offset > 0, ` OFFSET `+I.ToS(qp.Offset))
	if DEBUG {
		L.Print(query)
	}
	qp.Rows = conn.QMapArray(query)
}
