package mZzz

import (
	"testing"

	"github.com/kokizzu/gotro/D/Tt"
)

func TestGenerateORM(t *testing.T) {
	Tt.GenerateOrm(TarantoolTables)
}
