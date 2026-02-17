package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".opserc"

type SessionConfig struct {
	PortraitsEnabled bool `json:"portraits_enabled"`
}

func DefaultSessionConfig() *SessionConfig {
	return &SessionConfig{
		PortraitsEnabled: true,
	}
}

// globalConfigPath returns ~/.config/opse/.opserc
func globalConfigPath() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "opse", configFileName)
}

// LoadSessionConfig tries .opserc in CWD first, then ~/.config/opse/.opserc.
// Returns defaults if neither exists.
func LoadSessionConfig() (*SessionConfig, error) {
	for _, path := range []string{configFileName, globalConfigPath()} {
		data, err := os.ReadFile(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		cfg := DefaultSessionConfig()
		return cfg, json.Unmarshal(data, cfg)
	}
	return DefaultSessionConfig(), nil
}

// SaveSessionConfig writes to .opserc in CWD.
func SaveSessionConfig(cfg *SessionConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFileName, data, 0644)
}
