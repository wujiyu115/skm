package skill

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Skill struct {
	ID          string
	Name        string
	Description string
	SourceType  string
	SourceRef   string
	CentralPath string
	ContentHash string
	Enabled     bool
	Metadata    map[string]interface{}
}

type frontmatter struct {
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Metadata    map[string]interface{} `yaml:"metadata"`
}

func findSkillMD(dir string) (string, error) {
	for _, name := range []string{"SKILL.md", "skill.md"} {
		p := filepath.Join(dir, name)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("no SKILL.md found in %s", dir)
}

func ParseSkillMD(dir string) (*Skill, error) {
	path, err := findSkillMD(dir)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	fm, body := parseFrontmatter(string(data))

	s := &Skill{
		Name:     fm.Name,
		Metadata: fm.Metadata,
		Enabled:  true,
	}

	if s.Name == "" {
		s.Name = filepath.Base(dir)
	}

	if fm.Description != "" {
		s.Description = fm.Description
	} else {
		s.Description = extractFirstLine(body)
	}

	return s, nil
}

func parseFrontmatter(content string) (frontmatter, string) {
	var fm frontmatter

	if !strings.HasPrefix(strings.TrimSpace(content), "---") {
		return fm, content
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Scan() // skip first ---

	var yamlLines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			break
		}
		yamlLines = append(yamlLines, line)
	}

	var rest []string
	for scanner.Scan() {
		rest = append(rest, scanner.Text())
	}

	yamlStr := strings.Join(yamlLines, "\n")
	yaml.Unmarshal([]byte(yamlStr), &fm)

	return fm, strings.Join(rest, "\n")
}

func extractFirstLine(body string) string {
	scanner := bufio.NewScanner(strings.NewReader(body))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if len(line) > 100 {
			return line[:100]
		}
		return line
	}
	return ""
}

func HasSkillMD(dir string) bool {
	_, err := findSkillMD(dir)
	return err == nil
}
