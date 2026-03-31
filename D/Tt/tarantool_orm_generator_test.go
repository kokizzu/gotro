package Tt

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
				{`id`, Unsigned},
				{`name`, String},
				{`createdAt`, Integer},
				{`score`, Double},
				{`isActive`, Boolean},
				{`coords`, Array},
			},
			Unique1:         `id`,
			Unique2:         `name`,
			Unique3:         `createdAt`,
			Uniques:         []string{`name`, `createdAt`},
			Indexes:         []string{`score`},
			Spatial:         `coords`,
			Engine:          Memtx,
			AutoIncrementId: true,
			AutoCensorFields: []string{
				`name`,
			},
			GenGraphqlType: true,
		},
	})

	rqFile := filepath.Join(tmpDir, `rqt`, `rqt__ORM.GEN.go`)
	wcFile := filepath.Join(tmpDir, `wct`, `wct__ORM.GEN.go`)
	rqTestFile := filepath.Join(tmpDir, `rqt`, `rqt__ORM.GEN_test.go`)
	wcTestFile := filepath.Join(tmpDir, `wct`, `wct__ORM.GEN_test.go`)

	rqBytes, err := os.ReadFile(rqFile)
	require.NoError(t, err)
	wcBytes, err := os.ReadFile(wcFile)
	require.NoError(t, err)
	rqTestBytes, err := os.ReadFile(rqTestFile)
	require.NoError(t, err)
	wcTestBytes, err := os.ReadFile(wcTestFile)
	require.NoError(t, err)

	rqText := string(rqBytes)
	wcText := string(wcBytes)
	rqTestText := string(rqTestBytes)
	wcTestText := string(wcTestBytes)
	assert.Contains(t, rqText, `package rqt`)
	assert.Contains(t, rqText, `type Users struct`)
	assert.Contains(t, rqText, `func (u *Users) UniqueIndexName() string`)
	assert.Contains(t, rqText, `func (u *Users) SpatialIndexCoords() string`)
	assert.Contains(t, rqText, `func (u *Users) FindById() bool`)
	assert.Contains(t, rqText, `SqlSelectAllUncensoredFields`)
	assert.Contains(t, rqText, `UsersFieldTypeMap`)

	assert.Contains(t, wcText, `package wct`)
	assert.Contains(t, wcText, `type UsersMutator struct`)
	assert.Contains(t, wcText, `func (u *UsersMutator) DoInsert() bool`)
	assert.Contains(t, wcText, `func (u *UsersMutator) SetAll(`)
	assert.Contains(t, wcText, `func (u *UsersMutator) SetCoords(`)
	assert.Contains(t, rqTestText, `func TestGeneratedUsersOrmHelpers(t *testing.T)`)
	assert.Contains(t, rqTestText, `func TestGeneratedUsersDbMethodsPanic(t *testing.T)`)
	assert.Contains(t, wcTestText, `func TestGeneratedUsersUnit(t *testing.T)`)
	assert.Contains(t, wcTestText, `func TestGeneratedUsersCRUD(t *testing.T)`)
}
