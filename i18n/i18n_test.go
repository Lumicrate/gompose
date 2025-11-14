package i18n_test

import (
	"github.com/Lumicrate/gompose/i18n"
	Net "net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// Mock Context Implementation

type mockContext struct {
	req    *Net.Request
	params map[string]string
	store  map[string]any
	status int
	header map[string]string
}

func newMockContext(req *Net.Request) *mockContext {
	return &mockContext{
		req:    req,
		params: map[string]string{},
		store:  map[string]any{},
		header: map[string]string{},
	}
}

func (m *mockContext) JSON(code int, obj any)           {}
func (m *mockContext) Bind(obj any) error               { return nil }
func (m *mockContext) BindJSON(obj any) error           { return nil }
func (m *mockContext) Param(key string) string          { return m.params[key] }
func (m *mockContext) Query(key string) string          { return "" }
func (m *mockContext) QueryParams() map[string][]string { return nil }
func (m *mockContext) SetHeader(k, v string)            { m.header[k] = v }
func (m *mockContext) Method() string                   { return m.req.Method }
func (m *mockContext) Path() string                     { return m.req.URL.Path }
func (m *mockContext) SetStatus(code int)               { m.status = code }
func (m *mockContext) Status() int                      { return m.status }
func (m *mockContext) RemoteIP() string                 { return "" }
func (m *mockContext) Header(h string) string           { return m.req.Header.Get(h) }
func (m *mockContext) Body(string)                      {}
func (m *mockContext) Abort()                           {}
func (m *mockContext) Next()                            {}
func (m *mockContext) Set(k string, v any)              { m.store[k] = v }
func (m *mockContext) Get(k string) any                 { return m.store[k] }
func (m *mockContext) Request() *Net.Request            { return m.req }

// Tests

func TestCookieLanguageExtractor(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&Net.Cookie{Name: "lang", Value: "fr"})

	ctx := newMockContext(req)

	langs := i18n.CookieLanguageExtractor(ctx, i18n.LanguageExtractorOptions{
		"CookieName": "lang",
	})

	if len(langs) != 1 || langs[0] != "fr" {
		t.Fatalf("expected [fr], got %+v", langs)
	}
}

func TestHeaderLanguageExtractor(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Language", "es-MX,es;q=0.9")

	ctx := newMockContext(req)

	langs := i18n.HeaderLanguageExtractor(ctx, nil)

	if langs[0] != "es-MX" || langs[1] != "es" {
		t.Fatalf("unexpected result: %+v", langs)
	}
}

func TestURLPrefixLanguageExtractor(t *testing.T) {
	req := httptest.NewRequest("GET", "/de/products", nil)
	ctx := newMockContext(req)
	ctx.params["lang"] = "de"

	langs := i18n.URLPrefixLanguageExtractor(ctx, i18n.LanguageExtractorOptions{
		"URLPrefixName": "lang",
	})

	if len(langs) != 1 || langs[0] != "de" {
		t.Fatalf("expected [de], got %+v", langs)
	}
}

func TestNewI18n(t *testing.T) {
	tmp := t.TempDir()

	yaml := []byte(`
hello:
  other: "Hello Test"
`)

	err := os.WriteFile(filepath.Join(tmp, "en.yaml"), yaml, 0644)
	if err != nil {
		t.Fatalf("write temp yaml error: %v", err)
	}

	tr, err := i18n.NewI18n(tmp, "en")
	if err != nil {
		t.Fatalf("NewI18n error: %v", err)
	}

	msg := tr.T("hello")
	if msg != "Hello Test" {
		t.Fatalf("expected 'Hello Test', got '%s'", msg)
	}
}
