package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type SavedRoll struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Expression string `json:"expression"`
	Folder     string `json:"folder"`
	SortOrder  int    `json:"sort_order"`
}

type RollFolder struct {
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
}

type SavedRollsConfig struct {
	Folders []RollFolder `json:"folders"`
	Rolls   []SavedRoll  `json:"rolls"`
}

func savedRollsPath() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "opse", "saved_rolls.json")
}

func LoadSavedRolls() (*SavedRollsConfig, error) {
	data, err := os.ReadFile(savedRollsPath())
	if os.IsNotExist(err) {
		return &SavedRollsConfig{}, nil
	}
	if err != nil {
		return nil, err
	}
	var cfg SavedRollsConfig
	return &cfg, json.Unmarshal(data, &cfg)
}

func SaveSavedRolls(cfg *SavedRollsConfig) error {
	path := savedRollsPath()
	os.MkdirAll(filepath.Dir(path), 0755)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (c *SavedRollsConfig) Add(r SavedRoll) {
	c.Rolls = append(c.Rolls, r)
}

func (c *SavedRollsConfig) Delete(id string) {
	for i, r := range c.Rolls {
		if r.ID == id {
			c.Rolls = append(c.Rolls[:i], c.Rolls[i+1:]...)
			return
		}
	}
}

func (c *SavedRollsConfig) ByFolder() map[string][]SavedRoll {
	m := make(map[string][]SavedRoll)
	for _, r := range c.Rolls {
		m[r.Folder] = append(m[r.Folder], r)
	}
	return m
}

func (c *SavedRollsConfig) AddFolder(name string) {
	for _, f := range c.Folders {
		if f.Name == name {
			return
		}
	}
	c.Folders = append(c.Folders, RollFolder{Name: name, SortOrder: len(c.Folders)})
}

func (c *SavedRollsConfig) DeleteFolder(name string) {
	for i, f := range c.Folders {
		if f.Name == name {
			c.Folders = append(c.Folders[:i], c.Folders[i+1:]...)
			break
		}
	}
	// Remove all rolls in this folder
	kept := c.Rolls[:0]
	for _, r := range c.Rolls {
		if r.Folder != name {
			kept = append(kept, r)
		}
	}
	c.Rolls = kept
}

func (c *SavedRollsConfig) FolderRollCount(name string) int {
	n := 0
	for _, r := range c.Rolls {
		if r.Folder == name {
			n++
		}
	}
	return n
}
