package mongodb

import (
	"github.com/Lumicrate/gompose/db"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

// Mock DB Adapter

type MockDB struct {
	Entities map[string]any
}

func NewMockDB() *MockDB {
	return &MockDB{Entities: make(map[string]any)}
}

func (m *MockDB) Init() error                        { return nil }
func (m *MockDB) Migrate(entities []any) error       { return nil }
func (m *MockDB) Create(entity any) error            { return nil }
func (m *MockDB) Update(entity any) error            { return nil }
func (m *MockDB) Delete(id string, entity any) error { return nil }
func (m *MockDB) FindAll(entity any, filters map[string]any, pagination db.Pagination, sort []db.Sort) (any, error) {
	return []any{}, nil
}
func (m *MockDB) FindByID(id string, entity any) (any, error) { return entity, nil }

// Test Entities

type TestEntity struct {
	ID   string
	Name string
}

type TestIntEntity struct {
	ID   int
	Name string
}

// Tests

func TestMockDBAdapterCRUD(t *testing.T) {
	mock := NewMockDB()

	e := &TestEntity{ID: "1", Name: "Alice"}

	// Test Create
	err := mock.Create(e)
	require.NoError(t, err)

	// Test Update
	err = mock.Update(e)
	require.NoError(t, err)

	// Test Delete
	err = mock.Delete("1", e)
	require.NoError(t, err)

	// Test FindAll
	result, err := mock.FindAll(TestEntity{}, nil, db.Pagination{}, nil)
	require.NoError(t, err)
	require.IsType(t, []any{}, result)

	// Test FindByID
	entity, err := mock.FindByID("1", e)
	require.NoError(t, err)
	require.Equal(t, e, entity)
}

// MongoDB Helper Tests

func TestGetTypedId_StringID(t *testing.T) {
	typ := reflect.TypeOf(TestEntity{})
	id, err := getTypedId("abc123", typ)
	require.NoError(t, err)
	require.Equal(t, "abc123", id)
}

func TestGetTypedId_IntID(t *testing.T) {
	typ := reflect.TypeOf(TestIntEntity{})
	id, err := getTypedId("42", typ)
	require.NoError(t, err)
	require.Equal(t, 42, id)
}

func TestGetTypedId_InvalidIntID(t *testing.T) {
	typ := reflect.TypeOf(TestIntEntity{})
	_, err := getTypedId("notanint", typ)
	require.Error(t, err)
}

func TestGetEntityID_String(t *testing.T) {
	entity := &TestEntity{ID: "42"}
	id, err := getEntityID(entity)
	require.NoError(t, err)
	require.Equal(t, "42", id)
}

func TestGetEntityID_Int(t *testing.T) {
	entity := &TestIntEntity{ID: 99}
	id, err := getEntityID(entity)
	require.NoError(t, err)
	require.Equal(t, "99", id)
}

func TestGetEntityID_MissingID(t *testing.T) {
	type NoID struct {
		Name string
	}
	entity := &NoID{Name: "Bob"}
	_, err := getEntityID(entity)
	require.Error(t, err)
	require.Equal(t, "ID field not found", err.Error())
}

// Simple Filter Test

func TestMockFindAllWithFilters(t *testing.T) {
	mock := NewMockDB()
	mock.Entities["1"] = &TestEntity{ID: "1", Name: "Alice"}
	mock.Entities["2"] = &TestEntity{ID: "2", Name: "Bob"}

	filters := map[string]any{"Name": "Alice"}

	result, err := mock.FindAll(TestEntity{}, filters, db.Pagination{}, nil)
	require.NoError(t, err)
	require.IsType(t, []any{}, result)
}
