package A

import (
	"testing"

	"github.com/kokizzu/gotro/L"
)

func TestToMsgp(t *testing.T) {
	m := []any{123, `abc`}
	L.Print(ToMsgp(m))
}
