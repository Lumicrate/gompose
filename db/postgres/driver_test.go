package postgres

import (
	"github.com/Lumicrate/gompose/db"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"testing"
)

// Test Entity

type TestEntity struct {
	ID   string `gorm:"primaryKey"`
	Name string
}

//Test Helpers

// setupTestAdapter creates an in-memory SQLite database that mimics Postgres behavior
func setupTestAdapter(t *testing.T) *PostgresAdapter {
	dbConn, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	adapter := &PostgresAdapter{db: dbConn}
	err = adapter.Migrate([]any{&TestEntity{}})
	require.NoError(t, err)

	return adapter
}

// Tests

func TestPostgresAdapter_Init(t *testing.T) {
	adapter := New("sqlite::memory:")
	// This will fail gracefully because the DSN is not a real Postgres URL
	err := adapter.Init()
	require.Error(t, err, "Init should fail with invalid DSN")
}

func TestPostgresAdapter_Migrate(t *testing.T) {
	adapter := setupTestAdapter(t)
	err := adapter.Migrate([]any{&TestEntity{}})
	require.NoError(t, err)
}

func TestPostgresAdapter_CreateAndFindByID(t *testing.T) {
	adapter := setupTestAdapter(t)

	entity := &TestEntity{ID: "1", Name: "Alice"}
	err := adapter.Create(entity)
	require.NoError(t, err)

	var result TestEntity
	found, err := adapter.FindByID("1", &result)
	require.NoError(t, err)
	require.Equal(t, entity.ID, found.(*TestEntity).ID)
	require.Equal(t, entity.Name, found.(*TestEntity).Name)
}

func TestPostgresAdapter_Update(t *testing.T) {
	adapter := setupTestAdapter(t)

	entity := &TestEntity{ID: "1", Name: "Alice"}
	require.NoError(t, adapter.Create(entity))

	entity.Name = "Updated"
	err := adapter.Update(entity)
	require.NoError(t, err)

	var updated TestEntity
	_, err = adapter.FindByID("1", &updated)
	require.NoError(t, err)
	require.Equal(t, "Updated", updated.Name)
}

func TestPostgresAdapter_Delete(t *testing.T) {
	adapter := setupTestAdapter(t)

	entity := &TestEntity{ID: "1", Name: "ToDelete"}
	require.NoError(t, adapter.Create(entity))

	err := adapter.Delete("1", &TestEntity{})
	require.NoError(t, err)

	var result TestEntity
	_, err = adapter.FindByID("1", &result)
	require.Error(t, err, "record should be deleted")
}

func TestPostgresAdapter_FindAll(t *testing.T) {
	adapter := setupTestAdapter(t)

	require.NoError(t, adapter.Create(&TestEntity{ID: "1", Name: "Alice"}))
	require.NoError(t, adapter.Create(&TestEntity{ID: "2", Name: "Bob"}))
	require.NoError(t, adapter.Create(&TestEntity{ID: "3", Name: "Charlie"}))

	filters := map[string]any{"Name": "Alice"}
	sort := []db.Sort{{Field: "Name", Direction: "asc"}}
	pagination := db.Pagination{Limit: 2, Offset: 0}

	results, err := adapter.FindAll(TestEntity{}, filters, pagination, sort)
	require.NoError(t, err)

	slice, ok := results.([]TestEntity)
	require.True(t, ok)
	require.Len(t, slice, 1)
	require.Equal(t, "Alice", slice[0].Name)
}

func TestPostgresAdapter_FindAll_NoFilters(t *testing.T) {
	adapter := setupTestAdapter(t)

	require.NoError(t, adapter.Create(&TestEntity{ID: "1", Name: "Alice"}))
	require.NoError(t, adapter.Create(&TestEntity{ID: "2", Name: "Bob"}))

	results, err := adapter.FindAll(TestEntity{}, nil, db.Pagination{}, nil)
	require.NoError(t, err)

	slice, ok := results.([]TestEntity)
	require.True(t, ok)
	require.Len(t, slice, 2)
}

func TestPostgresAdapter_MigrateMultipleEntities(t *testing.T) {
	type AnotherEntity struct {
		Code string `gorm:"primaryKey"`
		Desc string
	}

	adapter := setupTestAdapter(t)

	err := adapter.Migrate([]any{&TestEntity{}, &AnotherEntity{}})
	require.NoError(t, err)
}
