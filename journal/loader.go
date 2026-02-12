package journal

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var timestampRe = regexp.MustCompile(`^\*(\d{2}:\d{2}) — (.+)\*$`)

func Load(filePath string) (*Journal, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	content := string(data)
	title := extractTitle(content)
	createdAt := extractCreatedAt(content)

	j := &Journal{
		Title:     title,
		CreatedAt: createdAt,
		FilePath:  filePath,
	}

	body := extractBody(content)
	j.Entries = parseEntries(body)
	return j, nil
}

func extractTitle(md string) string {
	for _, line := range strings.Split(md, "\n") {
		if after, ok := strings.CutPrefix(line, "# "); ok {
			return after
		}
	}
	return "Untitled"
}

func extractCreatedAt(md string) time.Time {
	for _, line := range strings.Split(md, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "*Started: ") && strings.HasSuffix(line, "*") {
			dateStr := strings.TrimSuffix(strings.TrimPrefix(line, "*Started: "), "*")
			if t, err := time.Parse("2006-01-02", dateStr); err == nil {
				return t
			}
		}
	}
	return time.Now()
}

func extractBody(md string) string {
	_, after, found := strings.Cut(md, "---")
	if !found {
		return md
	}
	return strings.TrimLeft(after, "\n")
}

func parseEntries(body string) []Entry {
	lines := strings.Split(body, "\n")
	var entries []Entry
	var currentTs time.Time
	var currentSource string
	var currentLines []string

	flush := func() {
		text := strings.TrimSpace(strings.Join(currentLines, "\n"))
		if text == "" {
			return
		}
		entryType := EntryNarrative
		label := ""
		if strings.HasPrefix(text, "> ") {
			entryType = classifyBlockquote(text)
		}
		if entryType == EntryNarrative {
			if currentSource != "" && currentSource != "User" {
				label = currentSource
			}
		}
		entries = append(entries, Entry{
			Timestamp: currentTs,
			Type:      entryType,
			Label:     label,
			Markdown:  text,
		})
	}

	for _, line := range lines {
		if m := timestampRe.FindStringSubmatch(strings.TrimSpace(line)); m != nil {
			flush()
			currentLines = nil
			t, _ := time.Parse("15:04", m[1])
			currentTs = t
			currentSource = m[2]
			continue
		}
		currentLines = append(currentLines, line)
	}
	flush()

	return entries
}

func classifyBlockquote(text string) EntryType {
	lower := strings.ToLower(text)
	if strings.Contains(lower, "**oracle") {
		return EntryOracle
	}
	if strings.Contains(lower, "**set the scene") {
		return EntryScene
	}
	if strings.Contains(lower, "**dice") || strings.Contains(lower, "**coin flip") ||
		strings.Contains(lower, "**card draw") || strings.Contains(lower, "**direction") ||
		strings.Contains(lower, "**weather") || strings.Contains(lower, "**color") ||
		strings.Contains(lower, "**sound") {
		return EntryTool
	}
	return EntryGenerator
}

func ListJournals(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".md" {
			files = append(files, e.Name())
		}
	}
	return files, nil
}
