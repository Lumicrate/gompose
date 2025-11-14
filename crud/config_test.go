package crud

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if len(cfg.ProtectedMethods) != 0 {
		t.Errorf("expected empty ProtectedMethods, got %v", cfg.ProtectedMethods)
	}
}

func TestProtect(t *testing.T) {
	cfg := DefaultConfig()
	Protect("GET", "POST")(cfg)

	if !cfg.ProtectedMethods["GET"] || !cfg.ProtectedMethods["POST"] {
		t.Errorf("methods not protected: %v", cfg.ProtectedMethods)
	}
}

func TestProtectAll(t *testing.T) {
	cfg := DefaultConfig()
	ProtectAll()(cfg)

	expected := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	for _, m := range expected {
		if !cfg.ProtectedMethods[m] {
			t.Errorf("method %s not protected", m)
		}
	}
}
