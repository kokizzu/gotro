package Ch

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateOrmWritesExpectedFiles(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	tmpDir := t.TempDir()
	require.NoError(t, os.Chdir(tmpDir))
	defer func() {
		_ = os.Chdir(wd)
	}()

	GenerateOrm(map[TableName]*TableProp{
		`users`: {
			Fields: []Field{
				{`id`, UInt64},
				{`createdAt`, DateTime},
				{`name`, String},
				{`heightMeter`, Float64},
			},
			Engine: ReplacingMergeTree,
			Orders: []string{`id`},
		},
	})

	saFile := filepath.Join(tmpDir, `sah`, `sah__ORM.GEN.go`)
	saTestFile := filepath.Join(tmpDir, `sah`, `sah__ORM.GEN_test.go`)

	saBytes, err := os.ReadFile(saFile)
	require.NoError(t, err)
	saTestBytes, err := os.ReadFile(saTestFile)
	require.NoError(t, err)

	saText := string(saBytes)
	saTestText := string(saTestBytes)

	assert.Contains(t, saText, `package sah`)
	assert.Contains(t, saText, `type Users struct`)
	assert.Contains(t, saText, `func (u *Users) TableName() Ch.TableName`)
	assert.Contains(t, saText, `func (u *Users) SqlInsert() string`)
	assert.Contains(t, saText, `func (u *Users) ScanRowAllCols(rows *sql.Rows) (err error)`)
	assert.Contains(t, saText, `UsersFieldTypeMap`)
	assert.Contains(t, saText, `var Preparators = map[Ch.TableName]chBuffer.Preparator`)

	assert.Contains(t, saTestText, `func TestGeneratedUsersHelpers(t *testing.T)`)
	assert.Contains(t, saTestText, `func TestGeneratedUsersCRUD(t *testing.T)`)
	assert.Contains(t, saTestText, `func TestMain(m *testing.M)`)
	assert.Contains(t, saTestText, `a.UpsertTable(obj.TableName(), &Ch.TableProp{`)
}
