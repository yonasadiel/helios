package helios

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type A struct {
	DataX string
	DataY string
}

func (a *A) TableName() string {
	return "abc"
}

type B struct {
	DataZ string
}

func TestRegisterModel(t *testing.T) {
	App.BeforeTest()

	assert.Equal(t, 0, len(App.models), "No model registered")
	App.RegisterModel(A{})
	assert.Equal(t, 1, len(App.models), "Only one models registered")
	App.RegisterModel(B{})
	assert.Equal(t, 2, len(App.models), "Two models registered")

	App.Migrate()
	assert.True(t, DB.HasTable("abc"), "Model A should be migrated with table name abc")
	assert.True(t, DB.HasTable("bs"), "Model B should be migrated with default table name")

	App.CloseDB()
}
