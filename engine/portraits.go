package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// SavedPortrait pairs a character name with its procedural portrait parameters.
type SavedPortrait struct {
	Name   string         `json:"name"`
	Params PortraitParams `json:"params"`
}

// SavedPortraitsConfig holds all saved portraits.
type SavedPortraitsConfig struct {
	Portraits []SavedPortrait `json:"portraits"`
}

func savedPortraitsPath() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "opse", "portraits.json")
}

// LoadSavedPortraits reads portraits from ~/.config/opse/portraits.json.
// Returns an empty config (not an error) if the file doesn't exist.
func LoadSavedPortraits() (*SavedPortraitsConfig, error) {
	data, err := os.ReadFile(savedPortraitsPath())
	if os.IsNotExist(err) {
		return &SavedPortraitsConfig{}, nil
	}
	if err != nil {
		return nil, err
	}
	var cfg SavedPortraitsConfig
	return &cfg, json.Unmarshal(data, &cfg)
}

// SaveSavedPortraits writes portraits to ~/.config/opse/portraits.json.
func SaveSavedPortraits(cfg *SavedPortraitsConfig) error {
	path := savedPortraitsPath()
	os.MkdirAll(filepath.Dir(path), 0755)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Add appends a portrait to the config.
func (c *SavedPortraitsConfig) Add(p SavedPortrait) {
	c.Portraits = append(c.Portraits, p)
}

// Delete removes the first portrait matching name (case-insensitive).
func (c *SavedPortraitsConfig) Delete(name string) {
	for i, p := range c.Portraits {
		if strings.EqualFold(p.Name, name) {
			c.Portraits = append(c.Portraits[:i], c.Portraits[i+1:]...)
			return
		}
	}
}

// FindByName returns a pointer to the first portrait matching name
// (case-insensitive), or nil if not found.
func (c *SavedPortraitsConfig) FindByName(name string) *SavedPortrait {
	for i, p := range c.Portraits {
		if strings.EqualFold(p.Name, name) {
			return &c.Portraits[i]
		}
	}
	return nil
}
