package M

import (
	"testing"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/L"
)

func TestFoo(t *testing.T) {
	m := []interface{}{123, `abc`}
	L.Print(A.ToMsgp(m))
}
