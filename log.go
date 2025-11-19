package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// GameLog represents a solo RPG game log
type GameLog struct {
	Title       string    `yaml:"title"`
	CreatedAt   time.Time `yaml:"created_at"`
	LastUpdated time.Time `yaml:"last_updated"`
	Entries     []LogEntry `yaml:"entries"`
}

// LogEntry represents a single entry in the game log
type LogEntry struct {
	Timestamp time.Time `yaml:"timestamp"`
	Content   string    `yaml:"content"`
	Type      string    `yaml:"type"` // "user", "generator", "oracle", etc.
}

// SaveLog saves a game log to a YAML file
func SaveLog(log *GameLog, filename string) error {
	log.LastUpdated = time.Now()
	
	data, err := yaml.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to marshal log: %w", err)
	}
	
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

// LoadLog loads a game log from a YAML file
func LoadLog(filename string) (*GameLog, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	var log GameLog
	err = yaml.Unmarshal(data, &log)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal log: %w", err)
	}
	
	return &log, nil
}

// ListLogs lists all log files in the current directory
func ListLogs() ([]string, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}
	
	var logs []string
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if len(name) > 5 && name[len(name)-5:] == ".yaml" {
				logs = append(logs, name)
			}
		}
	}
	
	return logs, nil
}

// AddEntry adds a new entry to the game log
func (log *GameLog) AddEntry(content string, entryType string) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Content:   content,
		Type:      entryType,
	}
	log.Entries = append(log.Entries, entry)
	log.LastUpdated = time.Now()
}

// FormatLogEntry formats a log entry for display
func FormatLogEntry(entry LogEntry) string {
	return fmt.Sprintf("[%s] %s: %s",
		entry.Timestamp.Format("2006-01-02 15:04:05"),
		entry.Type,
		entry.Content)
}

