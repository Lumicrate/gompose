package ginadapter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	gompose_http "github.com/Lumicrate/gompose/http"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// Test Entities
type TestEntity struct {
	ID   string
	Name string
}

// Mock Context
type MockHandler struct {
	called bool
}

func (h *MockHandler) Handler(ctx gompose_http.Context) {
	h.called = true
}

// Tests

func TestGinEngine_RegisterRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := New(8080)

	handler := &MockHandler{}

	engine.RegisterRoute("GET", "/test", handler.Handler, &TestEntity{}, false)

	// Check that route is registered
	routes := engine.Routes()
	require.Len(t, routes, 1)
	require.Equal(t, "GET", routes[0].Method)
	require.Equal(t, "/test", routes[0].Path)
	require.False(t, routes[0].Protected)

	// Simulate a request to check handler is called
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	engine.engine.ServeHTTP(w, req)

	require.True(t, handler.called)
	require.Equal(t, http.StatusOK, w.Code) // default no status code set
}

func TestGinEngine_UseMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := New(8080)

	middlewareCalled := false
	engine.Use(func(next gompose_http.HandlerFunc) gompose_http.HandlerFunc {
		return func(ctx gompose_http.Context) {
			middlewareCalled = true
			next(ctx)
		}
	})

	handlerCalled := false
	engine.RegisterRoute("GET", "/mid", func(ctx gompose_http.Context) {
		handlerCalled = true
	}, &TestEntity{}, false)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/mid", nil)
	engine.engine.ServeHTTP(w, req)

	require.True(t, middlewareCalled)
	require.True(t, handlerCalled)
}

func TestGinEngine_RoutesMethod(t *testing.T) {
	engine := New(8080)
	engine.RegisterRoute("POST", "/create", func(ctx gompose_http.Context) {}, &TestEntity{}, true)

	routes := engine.Routes()
	require.Len(t, routes, 1)
	require.Equal(t, "POST", routes[0].Method)
	require.Equal(t, "/create", routes[0].Path)
	require.True(t, routes[0].Protected)
}

func TestGinEngine_RegisterRoute_UnsupportedMethodPanics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := New(8080)

	require.Panics(t, func() {
		engine.RegisterRoute("FOO", "/foo", func(ctx gompose_http.Context) {}, &TestEntity{}, false)
	})
}
