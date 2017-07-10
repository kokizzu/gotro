package tTags

import (
	"example-complete/sql"
	"github.com/kokizzu/gotro/D/Pg"
	"github.com/kokizzu/gotro/S"
)

const TABLE = `tags`

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
			{Label: `Name`, Key: `name`},
			{Label: `Type`, Key: `type`},
			{Label: `Note`, Key: `note`}, // tahun pemilihan, dapil, tag untuk labeling evidences
			{Label: `Parent`, Key: `parent_id`, Type: `bigint`},
		},
	}
	SELECT = TM_MASTER.Select()
}
