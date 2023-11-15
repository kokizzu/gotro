package zCrud

import (
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D/Ch"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/S"
	"golang.org/x/exp/maps"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file pager.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type pager.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type pager.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type pager.go
//go:generate farify doublequote --file pager.go

type PagerIn struct {
	Page    int `json:"page" form:"page" query:"page" long:"page" msg:"page"`
	PerPage int `json:"perPage" form:"perPage" query:"perPage" long:"perPage" msg:"perPage"`

	// filter AND by column, if value is array then filter OR on that field
	Filters map[string][]string `json:"filters" form:"filters" query:"filters" long:"filters" msg:"filters"`

	// Order: [+col1, -col2] (+ is ascending, - is descending)
	Order []string `json:"order" form:"order" query:"order" long:"order" msg:"order"`
}

type PagerOut struct {
	Page    int `json:"page" form:"page" query:"page" long:"page" msg:"page"`
	PerPage int `json:"perPage" form:"perPage" query:"perPage" long:"perPage" msg:"perPage"`

	Pages int `json:"pages" form:"pages" query:"pages" long:"pages" msg:"pages"`
	Total int `json:"countResult" form:"countResult" query:"countResult" long:"countResult" msg:"countResult"`

	Filters map[string][]string `json:"filters" form:"filters" query:"filters" long:"filters" msg:"filters"`

	Order []string `json:"order" form:"order" query:"order" long:"order" msg:"order"`
}

func (p *PagerOut) LimitOffsetSql() string {
	offsetStr := ``
	if p.Page > 1 {
		offsetStr = fmt.Sprintf(` OFFSET %d`, (p.Page-1)*p.PerPage)
	}
	return fmt.Sprintf(`
LIMIT %d`, p.PerPage) + offsetStr
}

// OrderBySqlTt generates ORDER BY SQL for Tarantool
func (p *PagerOut) OrderBySqlTt(orders []string, fieldToType map[string]Tt.DataType) (orderBySql string) {
	var orderBy []string
	for _, dirField := range orders {
		if len(dirField) <= 2 {
			continue
		}
		dir := dirField[0]
		dirStr := ``
		if dir == '+' {
		} else if dir == '-' {
			dirStr = ` DESC`
		} else {
			continue
		}
		field := dirField[1:]
		if _, ok := fieldToType[field]; !ok {
			continue
		}
		orderBy = append(orderBy, S.QQ(field)+dirStr)
	}
	if len(orderBy) > 0 {
		orderBySql = `
ORDER BY ` + A.StrJoin(orderBy, `, `)
	}
	return
}

// OrderBySqlCh generates WHERE AND SQL for Tarantool
func (p *PagerOut) OrderBySqlCh(orders []string, fieldToType map[string]Ch.DataType) (orderBySql string) {
	var orderBy []string
	for _, dirField := range orders {
		if len(dirField) <= 2 {
			continue
		}
		dir := dirField[0]
		dirStr := ``
		if dir == '+' {
		} else if dir == '-' {
			dirStr = ` DESC`
		} else {
			continue
		}
		field := dirField[1:]
		if _, ok := fieldToType[field]; !ok {
			continue
		}
		orderBy = append(orderBy, field+dirStr)
	}
	if len(orderBy) > 0 {
		orderBySql = `
ORDER BY ` + A.StrJoin(orderBy, `, `)
	}
	return
}

// WhereAndSqlTt generates WHERE AND SQL for Tarantool
func (p *PagerOut) WhereAndSqlTt(filters map[string][]string, fieldToType map[string]Tt.DataType) (whereAndSql string) {
	var whereAnd []string
	fields := maps.Keys(filters)
	sort.Strings(fields)
	for _, field := range fields {
		typ, ok := fieldToType[field]
		if !ok {
			continue
		}
		value := filters[field]
		quotedValue, filtered := equalityQuoteValueTt(value, typ, S.QQ(field))
		if len(quotedValue) > 1 {
			whereOr := A.StrJoin(quotedValue, ` OR `)
			whereAnd = append(whereAnd, whereOr)
		} else if len(quotedValue) == 1 {
			whereAnd = append(whereAnd, quotedValue[0])
		}
		if p.Filters == nil {
			p.Filters = map[string][]string{}
		}
		p.Filters[field] = filtered
	}
	if len(whereAnd) > 0 {
		whereAndSql = `
WHERE (` + A.StrJoin(whereAnd, `)
	AND (`) + `)`
	}
	return
}

// WhereAndSqlCh generates WHERE AND SQL for ClickHouse
func (p *PagerOut) WhereAndSqlCh(filters map[string][]string, fieldToType map[string]Ch.DataType) (whereAndSql string) {
	var whereAnd []string
	fields := maps.Keys(filters)
	sort.Strings(fields)
	for _, field := range fields {
		typ, ok := fieldToType[field]
		if !ok {
			continue
		}
		value := filters[field]
		quotedValue, filtered := equalityQuoteValueCh(value, typ, field)
		if len(quotedValue) > 1 {
			whereOr := A.StrJoin(quotedValue, ` OR `)
			whereAnd = append(whereAnd, whereOr)
		} else if len(quotedValue) == 1 {
			whereAnd = append(whereAnd, quotedValue[0])
		}
		if p.Filters == nil {
			p.Filters = map[string][]string{}
		}
		p.Filters[field] = filtered
	}
	if len(whereAnd) > 0 {
		whereAndSql = `
WHERE (` + A.StrJoin(whereAnd, `)
	AND (`) + `)`
	}
	return
}

func equalityQuoteValueTt(values []string, expectTyp Tt.DataType, field string) (whereOr []string, filtered []string) {

	// TODO: if value equal then make it unique

	// allow >, <, >=, <=, <>, *LIKE* and NOT LIKE, if multiple = or <> then will use IN or NOT IN
	// currently if value equal, last write wins
	switch expectTyp {
	case Tt.Unsigned, Tt.Integer, Tt.Double:
		var equalValues, unequalValues []string
		var gte, lte, gtf, ltf string
		gtv := math.MaxFloat64
		ltv := -math.MaxFloat64
		for _, str := range values {
			operator, rhs := splitOperatorValue(str)
			v, err := strconv.ParseFloat(rhs, 64)
			if err != nil {
				continue
			}
			if operator == `=` {
				filtered = append(filtered, rhs)
				equalValues = append(equalValues, rhs)
			} else if operator == `<>` {
				filtered = append(filtered, operator+rhs)
				unequalValues = append(unequalValues, rhs)
			} else {
				if gtv >= v && operator[0] == '>' {
					gtv = v
					gte = field + operator + rhs
					gtf = operator + rhs
				}
				if ltv <= v && operator[0] == '<' {
					ltv = v
					lte = field + operator + rhs
					ltf = operator + rhs
				}
			}
		}
		if gte != `` && lte != `` {
			filtered = append(filtered, gtf, ltf)
			if gtv < ltv {
				// autodetect intersection to use AND instead of OR
				whereOr = append(whereOr, `(`+gte+` AND `+lte+`)`)
			} else {
				whereOr = append(whereOr, lte, gte)
			}
		} else if gte != `` {
			filtered = append(filtered, gtf)
			whereOr = append(whereOr, gte)
		} else if lte != `` {
			filtered = append(filtered, ltf)
			whereOr = append(whereOr, lte)
		}
		if len(equalValues) > 0 {
			whereOr = append(whereOr, field+` IN (`+A.StrJoin(equalValues, `,`)+`)`)
		}
		if len(unequalValues) > 0 {
			whereOr = append(whereOr, field+` NOT IN (`+A.StrJoin(unequalValues, `,`)+`)`)
		}
	case Tt.String:
		var equalValues, unequalValues []string
		var gte, lte, gtf, ltf string
		gtv := `~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~`
		ltv := ``
		for _, str := range values {
			operator, rhs := splitOperatorValue(str)
			hasWildcard := S.Contains(rhs, `*`)
			if operator == `=` {
				filtered = append(filtered, rhs)
				if hasWildcard {
					rhs = S.Replace(rhs, `*`, `%`)
					operator = ` LIKE `
					whereOr = append(whereOr, field+operator+S.Z(rhs))
					continue
				}
				equalValues = append(equalValues, S.Z(rhs))
			} else if operator == `<>` {
				filtered = append(filtered, operator+rhs)
				if hasWildcard {
					rhs = S.Replace(rhs, `*`, `%`)
					operator = ` NOT LIKE `
					whereOr = append(whereOr, field+operator+S.Z(rhs))
					continue
				}
				unequalValues = append(unequalValues, S.Z(rhs))
			} else {
				if gtv >= rhs && operator[0] == '>' {
					gtv = rhs
					gte = field + operator + S.Z(rhs)
					gtf = operator + rhs
				}
				if ltv <= rhs && operator[0] == '<' {
					ltv = rhs
					lte = field + operator + S.Z(rhs)
					ltf = operator + rhs
				}
			}
		}
		if gte != `` && lte != `` {
			filtered = append(filtered, gtf, ltf)
			if gtv < ltv {
				// autodetect intersection to use AND instead of OR
				whereOr = append(whereOr, `(`+gte+` AND `+lte+`)`)
			} else {
				whereOr = append(whereOr, lte, gte)
			}
		} else if gte != `` {
			filtered = append(filtered, gtf)
			whereOr = append(whereOr, gte)
		} else if lte != `` {
			filtered = append(filtered, ltf)
			whereOr = append(whereOr, lte)
		}
		if len(equalValues) > 0 {
			whereOr = append(whereOr, field+` IN (`+A.StrJoin(equalValues, `,`)+`)`)
		}
		if len(unequalValues) > 0 {
			whereOr = append(whereOr, field+` NOT IN (`+A.StrJoin(unequalValues, `,`)+`)`)
		}
		// case Tt.Array: // assume geo
		// TODO: do geoquery, but with sql: https://t.me/tarantool/15882
		// case Tt.Boolean: // ignore for now
	}
	// TODO: return debug/filtered

	return
}

func equalityQuoteValueCh(values []string, expectTyp Ch.DataType, field string) (whereOr []string, filtered []string) {
	// TODO: if value equal then make it unique

	// allow >, <, >=, <=, <>, *LIKE* and NOT LIKE, if multiple = or <> then will use IN or NOT IN
	// currently if value equal, last write wins
	switch expectTyp {
	case Ch.Int8, Ch.Int16, Ch.Int32, Ch.Int64, Ch.UInt8, Ch.UInt16, Ch.UInt32, Ch.UInt64, Ch.Float32, Ch.Float64:
		var equalValues, unequalValues []string
		var gte, lte, gtf, ltf string
		gtv := math.MaxFloat64
		ltv := -math.MaxFloat64
		for _, str := range values {
			operator, rhs := splitOperatorValue(str)
			v, err := strconv.ParseFloat(rhs, 64)
			if err != nil {
				continue
			}
			if operator == `=` {
				filtered = append(filtered, rhs)
				equalValues = append(equalValues, rhs)
			} else if operator == `<>` {
				filtered = append(filtered, operator+rhs)
				unequalValues = append(unequalValues, rhs)
			} else {
				if gtv >= v && operator[0] == '>' {
					gtv = v
					gte = field + operator + rhs
					gtf = operator + rhs
				}
				if ltv <= v && operator[0] == '<' {
					ltv = v
					lte = field + operator + rhs
					ltf = operator + rhs
				}
			}
		}
		if gte != `` && lte != `` {
			filtered = append(filtered, gtf, ltf)
			if gtv < ltv {
				// autodetect intersection to use AND instead of OR
				whereOr = append(whereOr, `(`+gte+` AND `+lte+`)`)
			} else {
				whereOr = append(whereOr, lte, gte)
			}
		} else if gte != `` {
			filtered = append(filtered, gtf)
			whereOr = append(whereOr, gte)
		} else if lte != `` {
			filtered = append(filtered, ltf)
			whereOr = append(whereOr, lte)
		}
		if len(equalValues) > 0 {
			whereOr = append(whereOr, field+` IN (`+A.StrJoin(equalValues, `,`)+`)`)
		}
		if len(unequalValues) > 0 {
			whereOr = append(whereOr, field+` NOT IN (`+A.StrJoin(unequalValues, `,`)+`)`)
		}
	case Ch.String, Ch.FixedString:
		var equalValues, unequalValues []string
		var gte, lte, gtf, ltf string
		gtv := `~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~`
		ltv := ``
		for _, str := range values {
			operator, rhs := splitOperatorValue(str)
			hasWildcard := S.Contains(rhs, `*`)
			if operator == `=` {
				filtered = append(filtered, rhs)
				if hasWildcard {
					rhs = S.Replace(rhs, `*`, `%`)
					operator = ` LIKE `
					whereOr = append(whereOr, field+operator+S.Z(rhs))
					continue
				}
				equalValues = append(equalValues, S.Z(rhs))
			} else if operator == `<>` {
				filtered = append(filtered, operator+rhs)
				if hasWildcard {
					rhs = S.Replace(rhs, `*`, `%`)
					operator = ` NOT LIKE `
					whereOr = append(whereOr, field+operator+S.Z(rhs))
					continue
				}
				unequalValues = append(unequalValues, S.Z(rhs))
			} else {
				if gtv >= rhs && operator[0] == '>' {
					gtv = rhs
					gte = field + operator + S.Z(rhs)
					gtf = operator + rhs
				}
				if ltv <= rhs && operator[0] == '<' {
					ltv = rhs
					lte = field + operator + S.Z(rhs)
					ltf = operator + rhs
				}
			}
		}
		if gte != `` && lte != `` {
			filtered = append(filtered, gtf, ltf)
			if gtv < ltv {
				// autodetect intersection to use AND instead of OR
				whereOr = append(whereOr, `(`+gte+` AND `+lte+`)`)
			} else {
				whereOr = append(whereOr, lte, gte)
			}
		} else if gte != `` {
			filtered = append(filtered, gtf)
			whereOr = append(whereOr, gte)
		} else if lte != `` {
			filtered = append(filtered, ltf)
			whereOr = append(whereOr, lte)
		}
		if len(equalValues) > 0 {
			whereOr = append(whereOr, field+` IN (`+A.StrJoin(equalValues, `,`)+`)`)
		}
		if len(unequalValues) > 0 {
			whereOr = append(whereOr, field+` NOT IN (`+A.StrJoin(unequalValues, `,`)+`)`)
		}
		// case Tt.Array: // assume geo
		// TODO: do geoquery, but with sql: https://t.me/tarantool/15882
		// case Tt.Boolean: // ignore for now
	}
	// TODO: return debug/filtered

	return
}

func splitOperatorValue(str string) (op string, rhs string) {
	l := len(str)
	if l < 1 {
		op = `=`
		return
	}
	equal := l > 1 && str[1] == '='
	startCh := 0
	if str[0] == '>' {
		startCh = 1
		if equal {
			startCh = 2
		}
		op = str[:startCh]
	} else if str[0] == '<' {
		startCh = 1
		if equal {
			startCh = 2
		} else if l > 1 && str[1] == '>' {
			startCh = 2
		}
		op = str[:startCh]
	} else {
		op = `=`
	}
	rhs = str[startCh:]
	return
}

const maxPerPage = 1000

func limit(perPage int) int {
	const defaultPerPage = 10
	if perPage <= 0 {
		return defaultPerPage
	}
	if perPage >= maxPerPage {
		return maxPerPage
	}
	return perPage
}

func offset(page, perPage, count int) (offset int) {
	if page <= 0 {
		return 0
	}
	// set to last page if overflow
	expectedOffset := (page - 1) * perPage
	maxOffset := (count / perPage) * perPage
	if expectedOffset > maxOffset {
		expectedOffset = maxOffset
	}
	return expectedOffset
}

func (p *PagerOut) CalculatePages(page, perPage, count int) {
	p.PerPage = limit(perPage)
	offset := offset(page, p.PerPage, count)
	p.Page = offset/p.PerPage + 1
	p.Total = count
	if count > 0 {
		p.Pages = (p.PerPage - 1 + count) / p.PerPage
	}
}
