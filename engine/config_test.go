package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultSessionConfig(t *testing.T) {
	cfg := DefaultSessionConfig()
	if !cfg.PortraitsEnabled {
		t.Error("expected portraits_enabled true by default")
	}
}

func TestLoadSessionConfigMissing(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	cfg, err := LoadSessionConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.PortraitsEnabled {
		t.Error("expected default portraits_enabled true")
	}
}

func TestSaveAndLoadSessionConfig(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	cfg := &SessionConfig{
		PortraitsEnabled: false,
	}
	if err := SaveSessionConfig(cfg); err != nil {
		t.Fatalf("save error: %v", err)
	}

	// Verify file is named .opserc
	if _, err := os.Stat(filepath.Join(dir, ".opserc")); err != nil {
		t.Fatalf("config file not created: %v", err)
	}

	loaded, err := LoadSessionConfig()
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded.PortraitsEnabled {
		t.Error("expected portraits_enabled false")
	}
}

func TestLoadSessionConfigCWDOverridesGlobal(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	// Write a global config
	globalDir := filepath.Join(dir, "fakehome", "opse")
	os.MkdirAll(globalDir, 0755)
	globalCfg := SessionConfig{PortraitsEnabled: true}
	data, _ := json.MarshalIndent(globalCfg, "", "  ")
	os.WriteFile(filepath.Join(globalDir, ".opserc"), data, 0644)

	// Write a CWD config with different values
	localCfg := SessionConfig{PortraitsEnabled: false}
	data, _ = json.MarshalIndent(localCfg, "", "  ")
	os.WriteFile(filepath.Join(dir, ".opserc"), data, 0644)

	// CWD should win
	loaded, err := LoadSessionConfig()
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded.PortraitsEnabled {
		t.Error("expected CWD portraits_enabled false to override global")
	}
}
