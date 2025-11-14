package crud

import (
	"testing"

	gomposehttp "github.com/Lumicrate/gompose/http"
	"github.com/stretchr/testify/require"
)

//  Mock Engine

type MockEngine struct {
	RoutesRegistered []gomposehttp.Route
}

func (m *MockEngine) Init(port int) error {
	return nil
}

func (m *MockEngine) RegisterRoute(method string, path string, handler gomposehttp.HandlerFunc, entity any, isProtected bool) {
	m.RoutesRegistered = append(m.RoutesRegistered, gomposehttp.Route{
		Method:    method,
		Path:      path,
		Entity:    entity,
		Protected: isProtected,
	})
}

func (m *MockEngine) Use(middleware gomposehttp.MiddlewareFunc) {}
func (m *MockEngine) Start() error                              { return nil }
func (m *MockEngine) Routes() []gomposehttp.Route               { return m.RoutesRegistered }

//  Mock Auth

type MockAuth struct{}

func (m *MockAuth) Init() error                                  { return nil }
func (m *MockAuth) RegisterRoutes(engine gomposehttp.HTTPEngine) {}
func (m *MockAuth) Middleware() gomposehttp.MiddlewareFunc {
	return func(next gomposehttp.HandlerFunc) gomposehttp.HandlerFunc {
		return func(ctx gomposehttp.Context) {
			// For testing, just call the next handler
			next(ctx)
		}
	}
}

//  Test

func TestRegisterCRUDRoutes(t *testing.T) {
	engine := &MockEngine{}
	dbAdapter := &MockDB{} // You can reuse your previous MockDB
	authProvider := &MockAuth{}

	config := DefaultConfig()
	config.ProtectedMethods["POST"] = true // protect POST only for testing

	RegisterCRUDRoutes(engine, dbAdapter, TestEntity{}, config, authProvider)

	routes := engine.Routes()
	require.Len(t, routes, 6, "Expected 6 CRUD routes registered")

	expected := map[string]string{
		"GET":    "/testentities",
		"GETID":  "/testentities/:id",
		"POST":   "/testentities",
		"PUT":    "/testentities/:id",
		"PATCH":  "/testentities/:id",
		"DELETE": "/testentities/:id",
	}

	for _, r := range routes {
		switch r.Method {
		case "GET":
			if r.Path == expected["GET"] {
				require.False(t, r.Protected)
			} else if r.Path == expected["GETID"] {
				require.False(t, r.Protected)
			}
		case "POST":
			require.Equal(t, expected["POST"], r.Path)
			require.True(t, r.Protected)
		case "PUT":
			require.Equal(t, expected["PUT"], r.Path)
			require.False(t, r.Protected)
		case "PATCH":
			require.Equal(t, expected["PATCH"], r.Path)
			require.False(t, r.Protected)
		case "DELETE":
			require.Equal(t, expected["DELETE"], r.Path)
			require.False(t, r.Protected)
		}
	}
}
