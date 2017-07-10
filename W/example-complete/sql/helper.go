package sql

import (
	"bytes"
	"github.com/kokizzu/gotro/S"
)

func Where_InIds(ids []string) string {
	len := len(ids) - 1
	if len < 0 {
		return ` IN ('0')` // make sure there are no zero-value id
	}
	buf := bytes.Buffer{}
	buf.WriteString(` IN (`)
	for k, v := range ids {
		buf.WriteString(S.Z(v))
		if k < len { // write except the last one
			buf.WriteRune(',')
		}
	}
	buf.WriteString(`)`)
	return buf.String()
}
