package crud

import (
	"errors"
	"github.com/Lumicrate/gompose/db"
	"testing"

	"net/http"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock DBAdapter

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Init() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDB) Migrate(entities []any) error {
	args := m.Called(entities)
	return args.Error(0)
}

func (m *MockDB) Create(entity any) error {
	args := m.Called(entity)
	return args.Error(0)
}

func (m *MockDB) Update(entity any) error {
	args := m.Called(entity)
	return args.Error(0)
}

func (m *MockDB) Delete(id string, entity any) error {
	args := m.Called(id, entity)
	return args.Error(0)
}

func (m *MockDB) FindAll(entity any, filters map[string]any, pagination db.Pagination, sort []db.Sort) (any, error) {
	args := m.Called(entity, filters, pagination, sort)
	return args.Get(0), args.Error(1)
}

func (m *MockDB) FindByID(id string, entity any) (any, error) {
	args := m.Called(id, entity)
	return args.Get(0), args.Error(1)
}

// Mock Context

type MockContext struct {
	mock.Mock
	status int
	Resp   any
}

func (m *MockContext) JSON(code int, obj any) {
	m.status = code
	m.Resp = obj
}

func (m *MockContext) Bind(obj any) error {
	args := m.Called(obj)
	return args.Error(0)
}

func (m *MockContext) Param(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func (m *MockContext) Query(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func (m *MockContext) QueryParams() map[string][]string {
	args := m.Called()
	return args.Get(0).(map[string][]string)
}

func (m *MockContext) BindJSON(obj any) error {
	args := m.Called(obj)
	return args.Error(0)
}

func (m *MockContext) SetHeader(key, value string) {
	m.Called(key, value)
}

func (m *MockContext) Method() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockContext) Path() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockContext) SetStatus(code int) {
	m.status = code
	m.Called(code)
}

func (m *MockContext) Status() int {
	return m.status
}

func (m *MockContext) RemoteIP() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockContext) Header(header string) string {
	args := m.Called(header)
	return args.String(0)
}

func (m *MockContext) Body(key string) {
	m.Called(key)
}

func (m *MockContext) Abort() {
	m.Called()
}

func (m *MockContext) Next() {
	m.Called()
}

func (m *MockContext) Set(key string, value any) {
	m.Called(key, value)
}

func (m *MockContext) Get(key string) any {
	args := m.Called(key)
	return args.Get(0)
}

func (m *MockContext) Request() *http.Request {
	args := m.Called()
	return args.Get(0).(*http.Request)
}

type TestEntity struct {
	ID   string
	Name string
}

type HookEntity struct {
	TestEntity
	BeforeCreateErr error
	AfterCreateErr  error
	BeforeUpdateErr error
	AfterUpdateErr  error
	BeforePatchErr  error
	AfterPatchErr   error
	BeforeDeleteErr error
	AfterDeleteErr  error
}

func (h *HookEntity) BeforeCreate() error { return h.BeforeCreateErr }
func (h *HookEntity) AfterCreate() error  { return h.AfterCreateErr }
func (h *HookEntity) BeforeUpdate() error { return h.BeforeUpdateErr }
func (h *HookEntity) AfterUpdate() error  { return h.AfterUpdateErr }
func (h *HookEntity) BeforePatch() error  { return h.BeforePatchErr }
func (h *HookEntity) AfterPatch() error   { return h.AfterPatchErr }
func (h *HookEntity) BeforeDelete() error { return h.BeforeDeleteErr }
func (h *HookEntity) AfterDelete() error  { return h.AfterDeleteErr }

// Tests
// handleGetAll

func TestHandleGetAll_Success(t *testing.T) {
	mockDB := new(MockDB)
	mockCtx := new(MockContext)

	expected := []TestEntity{{ID: "1", Name: "Alice"}}
	mockDB.On("FindAll", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(expected, nil)

	mockCtx.On("QueryParams").Return(map[string][]string{})

	handleGetAll(mockCtx, mockDB, TestEntity{})

	require.Equal(t, 200, mockCtx.Status())
	require.Equal(t, expected, mockCtx.Resp)
	mockDB.AssertExpectations(t)
	mockCtx.AssertExpectations(t)
}

func TestHandleGetAll_Error(t *testing.T) {
	mockDB := new(MockDB)
	mockCtx := new(MockContext)

	mockDB.On("FindAll", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("db error"))

	mockCtx.On("QueryParams").Return(map[string][]string{})

	handleGetAll(mockCtx, mockDB, TestEntity{})

	require.Equal(t, 500, mockCtx.Status())
	require.Contains(t, mockCtx.Resp.(map[string]string)["error"], "db error")
}

// handleGetByID

func TestHandleGetByID_Success(t *testing.T) {
	mockDB := new(MockDB)
	mockCtx := new(MockContext)

	entity := &TestEntity{ID: "1", Name: "Bob"}

	mockCtx.On("Param", "id").Return("1")
	mockDB.On("FindByID", "1", mock.Anything).Return(entity, nil)

	handleGetByID(mockCtx, mockDB, TestEntity{})

	require.Equal(t, 200, mockCtx.Status())
	require.Equal(t, entity, mockCtx.Resp)
}

func TestHandleGetByID_NotFound(t *testing.T) {
	mockDB := new(MockDB)
	mockCtx := new(MockContext)

	mockCtx.On("Param", "id").Return("99")
	mockDB.On("FindByID", "99", mock.Anything).Return(nil, errors.New("not found"))

	handleGetByID(mockCtx, mockDB, TestEntity{})

	require.Equal(t, 404, mockCtx.Status())
	require.Contains(t, mockCtx.Resp.(map[string]string)["error"], "entity not found")
}

// handleCreate

func TestHandleCreate_Success(t *testing.T) {
	mockDB := new(MockDB)
	mockCtx := new(MockContext)

	newEntity := &TestEntity{Name: "Charlie"}

	mockCtx.On("Bind", mock.Anything).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*TestEntity)
		arg.Name = newEntity.Name
	}).Return(nil)

	mockDB.On("Create", mock.Anything).Return(nil)

	handleCreate(mockCtx, mockDB, TestEntity{})

	require.Equal(t, 201, mockCtx.Status())
	require.Equal(t, newEntity.Name, mockCtx.Resp.(*TestEntity).Name)
}

func TestHandleCreate_BindError(t *testing.T) {
	mockDB := new(MockDB)
	mockCtx := new(MockContext)

	mockCtx.On("Bind", mock.Anything).Return(errors.New("bad input"))

	handleCreate(mockCtx, mockDB, TestEntity{})

	require.Equal(t, 400, mockCtx.Status())
	require.Contains(t, mockCtx.Resp.(map[string]string)["error"], "invalid input")
}

func TestHandleCreate_DBError(t *testing.T) {
	mockDB := new(MockDB)
	mockCtx := new(MockContext)

	mockCtx.On("Bind", mock.Anything).Return(nil)
	mockDB.On("Create", mock.Anything).Return(errors.New("insert failed"))

	handleCreate(mockCtx, mockDB, TestEntity{})

	require.Equal(t, 500, mockCtx.Status())
	require.Contains(t, mockCtx.Resp.(map[string]string)["error"], "insert failed")
}

// handleUpdate

func TestHandleUpdate_Success(t *testing.T) {
	mockDB := new(MockDB)
	mockCtx := new(MockContext)

	mockCtx.On("Param", "id").Return("1")
	mockCtx.On("Bind", mock.Anything).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*TestEntity)
		arg.Name = "Updated"
	}).Return(nil)

	mockDB.On("Update", mock.Anything).Return(nil)

	handleUpdate(mockCtx, mockDB, TestEntity{})

	require.Equal(t, 200, mockCtx.Status())
	require.Equal(t, "Updated", mockCtx.Resp.(*TestEntity).Name)
}

func TestHandleUpdate_BindError(t *testing.T) {
	mockDB := new(MockDB)
	mockCtx := new(MockContext)

	mockCtx.On("Param", "id").Return("1")
	mockCtx.On("Bind", mock.Anything).Return(errors.New("bad input"))

	handleUpdate(mockCtx, mockDB, TestEntity{})

	require.Equal(t, 400, mockCtx.Status())
}

func TestHandleUpdate_DBError(t *testing.T) {
	mockDB := new(MockDB)
	mockCtx := new(MockContext)

	mockCtx.On("Param", "id").Return("1")
	mockCtx.On("Bind", mock.Anything).Return(nil)
	mockDB.On("Update", mock.Anything).Return(errors.New("update failed"))

	handleUpdate(mockCtx, mockDB, TestEntity{})

	require.Equal(t, 500, mockCtx.Status())
}

// handlePatch

func TestHandlePatch_Success(t *testing.T) {
	mockDB := new(MockDB)
	mockCtx := new(MockContext)

	entity := &TestEntity{ID: "1", Name: "Old"}
	mockCtx.On("Param", "id").Return("1")
	mockDB.On("FindByID", "1", mock.Anything).Return(entity, nil)

	patchData := map[string]interface{}{"Name": "New"}
	mockCtx.On("BindJSON", mock.Anything).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*map[string]interface{})
		*arg = patchData
	}).Return(nil)

	mockDB.On("Update", mock.Anything).Return(nil)

	handlePatch(mockCtx, mockDB, TestEntity{})

	require.Equal(t, 200, mockCtx.Status())
	require.Equal(t, "New", mockCtx.Resp.(*TestEntity).Name)
}

func TestHandlePatch_NotFound(t *testing.T) {
	mockDB := new(MockDB)
	mockCtx := new(MockContext)

	mockCtx.On("Param", "id").Return("99")
	mockDB.On("FindByID", "99", mock.Anything).Return(nil, errors.New("not found"))

	handlePatch(mockCtx, mockDB, TestEntity{})

	require.Equal(t, 404, mockCtx.Status())
}

func TestHandlePatch_BindError(t *testing.T) {
	mockDB := new(MockDB)
	mockCtx := new(MockContext)

	entity := &TestEntity{ID: "1", Name: "Old"}
	mockCtx.On("Param", "id").Return("1")
	mockDB.On("FindByID", "1", mock.Anything).Return(entity, nil)

	mockCtx.On("BindJSON", mock.Anything).Return(errors.New("bad patch"))

	handlePatch(mockCtx, mockDB, TestEntity{})

	require.Equal(t, 400, mockCtx.Status())
}

// handleDelete

func TestHandleDelete_Success(t *testing.T) {
	mockDB := new(MockDB)
	mockCtx := new(MockContext)

	mockCtx.On("Param", "id").Return("1")
	mockDB.On("Delete", "1", mock.Anything).Return(nil)

	handleDelete(mockCtx, mockDB, TestEntity{})

	require.Equal(t, 204, mockCtx.Status())
	require.Nil(t, mockCtx.Resp)
}

func TestHandleDelete_DBError(t *testing.T) {
	mockDB := new(MockDB)
	mockCtx := new(MockContext)

	mockCtx.On("Param", "id").Return("1")
	mockDB.On("Delete", "1", mock.Anything).Return(errors.New("delete failed"))

	handleDelete(mockCtx, mockDB, TestEntity{})

	require.Equal(t, 500, mockCtx.Status())
	require.Contains(t, mockCtx.Resp.(map[string]string)["error"], "delete failed")
}
