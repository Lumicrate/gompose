package jwt

import (
	"testing"

	Net "net/http"

	"github.com/Lumicrate/gompose/db"
	"github.com/Lumicrate/gompose/http"
	"github.com/stretchr/testify/require"
)

// Mock DB
type MockDB struct {
	Created    []any
	FindRes    any
	FindErr    error
	MigrateErr error
}

func (m *MockDB) Init() error                        { return nil }
func (m *MockDB) Migrate(entities []any) error       { return m.MigrateErr }
func (m *MockDB) Create(entity any) error            { m.Created = append(m.Created, entity); return nil }
func (m *MockDB) Update(entity any) error            { return nil }
func (m *MockDB) Delete(id string, entity any) error { return nil }
func (m *MockDB) FindAll(entity any, filters map[string]any, pagination db.Pagination, sort []db.Sort) (any, error) {
	if m.FindErr != nil {
		return nil, m.FindErr
	}
	return m.FindRes, nil
}
func (m *MockDB) FindByID(id string, entity any) (any, error) { return entity, nil }

// Mock Context
type MockContext struct {
	Input      any
	Response   any
	StatusCode int
	Headers    map[string]string
	Aborted    bool
	Values     map[string]any
}

func (m *MockContext) JSON(code int, obj any)      { m.StatusCode = code; m.Response = obj }
func (m *MockContext) Bind(obj any) error          { return nil }
func (m *MockContext) Param(key string) string     { return "" }
func (m *MockContext) Query(key string) string     { return "" }
func (m *MockContext) SetHeader(key, value string) {}
func (m *MockContext) Method() string              { return "POST" }
func (m *MockContext) Path() string                { return "/test" }
func (m *MockContext) Status() int                 { return m.StatusCode }
func (m *MockContext) SetStatus(code int)          { m.StatusCode = code }
func (m *MockContext) RemoteIP() string            { return "127.0.0.1" }
func (m *MockContext) Abort()                      { m.Aborted = true }
func (m *MockContext) Next()                       {}
func (m *MockContext) Header(header string) string { return "Bearer token" }
func (m *MockContext) Set(key string, value any) {
	if m.Values == nil {
		m.Values = make(map[string]any)
	}
	m.Values[key] = value
}
func (m *MockContext) Get(key string) any               { return m.Values[key] }
func (m *MockContext) Body(content string)              {}
func (m *MockContext) Request() *Net.Request            { return nil }
func (m *MockContext) QueryParams() map[string][]string { return nil }
func (m *MockContext) BindJSON(obj any) error           { return nil }

// Test AuthUser Model
type TestUser struct {
	ID       string
	Email    string
	Password string
}

func (t *TestUser) GetID() string             { return t.ID }
func (t *TestUser) GetEmail() string          { return t.Email }
func (t *TestUser) GetHashedPassword() string { return t.Password }

// Tests

func TestJWTAuthProvider_Init(t *testing.T) {
	db := &MockDB{}
	provider := &JWTAuthProvider{
		SecretKey: "secret",
		DB:        db,
		UserModel: &TestUser{},
	}
	err := provider.Init()
	require.NoError(t, err)

	// SecretKey missing
	provider2 := &JWTAuthProvider{DB: db, UserModel: &TestUser{}}
	err = provider2.Init()
	require.Error(t, err)

	// UserModel nil
	provider3 := &JWTAuthProvider{DB: db, SecretKey: "secret"}
	err = provider3.Init()
	require.Error(t, err)
}

func TestJWTAuthProvider_SetUserModel(t *testing.T) {
	provider := &JWTAuthProvider{}
	provider.SetUserModel(&TestUser{})
	require.Equal(t, &TestUser{}, provider.UserModel)

	// Should panic if not AuthUser
	require.Panics(t, func() {
		provider.SetUserModel(struct{ ID string }{})
	})
}

func TestJWTAuthProvider_RegisterHandler_Success(t *testing.T) {
	db := &MockDB{}
	ctx := &MockContext{}
	provider := &JWTAuthProvider{
		SecretKey: "secret",
		DB:        db,
		UserModel: &TestUser{},
	}

	provider.registerHandler(ctx)
	require.Equal(t, 201, ctx.StatusCode)
	require.Len(t, db.Created, 1)
}

func TestJWTAuthProvider_Middleware(t *testing.T) {
	provider := &JWTAuthProvider{
		SecretKey: "secret",
	}

	called := false
	ctx := &MockContext{}
	mw := provider.Middleware()(func(c http.Context) {
		c.Set("called", true)
		called = true
	})

	// Simulate header missing
	ctx.Headers = map[string]string{}
	mw(ctx)
	require.Equal(t, 401, ctx.StatusCode)
	require.True(t, ctx.Aborted)
	require.False(t, called)
}
