package ginadapter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// TestStruct for Bind/JSON testing
type TestStruct struct {
	Name string `json:"name"`
}

func setupGinContext(method, path string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	return c, w
}

func TestGinContext_JSON(t *testing.T) {
	c, w := setupGinContext("GET", "/", nil)
	g := &GinContext{ctx: c}

	obj := TestStruct{Name: "Alice"}
	g.JSON(http.StatusOK, obj)

	require.Equal(t, http.StatusOK, w.Code)
	var res TestStruct
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	require.Equal(t, "Alice", res.Name)
}

func TestGinContext_Bind(t *testing.T) {
	payload := []byte(`{"name":"Bob"}`)
	c, _ := setupGinContext("POST", "/", payload)
	g := &GinContext{ctx: c}

	var target TestStruct
	err := g.Bind(&target)

	require.NoError(t, err)
	require.Equal(t, "Bob", target.Name)
}

func TestGinContext_ParamAndQuery(t *testing.T) {
	c, _ := setupGinContext("GET", "/user/123?role=admin", nil)
	c.Params = gin.Params{{Key: "id", Value: "123"}}
	g := &GinContext{ctx: c}

	require.Equal(t, "123", g.Param("id"))
	require.Equal(t, "admin", g.Query("role"))
}

func TestGinContext_SetHeaderAndStatus(t *testing.T) {
	c, w := setupGinContext("GET", "/", nil)
	g := &GinContext{ctx: c}

	g.SetHeader("X-Test", "ok")
	g.SetStatus(201)
	g.Body("hello")

	require.Equal(t, "ok", w.Header().Get("X-Test"))
	require.Equal(t, 201, g.Status())
	require.Contains(t, w.Body.String(), "hello")
}

func TestGinContext_MethodPathRemoteIP(t *testing.T) {
	c, _ := setupGinContext("POST", "/foo", nil)
	c.Request.RemoteAddr = "127.0.0.1:12345"
	g := &GinContext{ctx: c}

	require.Equal(t, "POST", g.Method())
	require.Equal(t, "/foo", g.Path())
	require.Equal(t, "127.0.0.1", g.RemoteIP())
}

func TestGinContext_SetGetValues(t *testing.T) {
	g := &GinContext{}

	require.Nil(t, g.Get("missing"))
	g.Set("key", 42)
	require.Equal(t, 42, g.Get("key"))
}

func TestGinContext_Header(t *testing.T) {
	c, _ := setupGinContext("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer xyz")
	g := &GinContext{ctx: c}

	require.Equal(t, "Bearer xyz", g.Header("Authorization"))
}

func TestGinContext_RequestReturnsSame(t *testing.T) {
	c, _ := setupGinContext("GET", "/", nil)
	g := &GinContext{ctx: c}
	require.Equal(t, c.Request, g.Request())
}
