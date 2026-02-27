// internal/config/config_test.go

package config

import (
    "os"
    "path/filepath"
    "testing"
)

func TestLoad_Default(t *testing.T) {
    cfg, err := Load("")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if cfg.Port != 12345 {
        t.Errorf("expected default port 12345, got %d", cfg.Port)
    }
}

func TestLoad_Env(t *testing.T) {
    os.Setenv("PRISM_PORT", "54321")
    defer os.Unsetenv("PRISM_PORT")

    cfg, err := Load("")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if cfg.Port != 54321 {
        t.Errorf("expected port from env, got %d", cfg.Port)
    }
}

func TestLoad_File(t *testing.T) {
    tmp := filepath.Join(os.TempDir(), "prism_config_test.json")
    f, err := os.Create(tmp)
    if err != nil {
        t.Fatal(err)
    }
    f.WriteString("port: 9999")
    f.Close()
    defer os.Remove(tmp)

    cfg, err := Load(tmp)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if cfg.Port != 9999 {
        t.Errorf("expected port 9999, got %d", cfg.Port)
    }
}
