package store

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type skillMeta struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SourceType  string `json:"source_type"`
	SourceRef   string `json:"source_ref"`
	ContentHash string `json:"content_hash"`
}

type groupMeta struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Skills      []string `json:"skills"`
}

// WriteMetadata writes JSON metadata files for all skills and groups into metaDir.
// This provides a human-readable, VCS-friendly mirror of the database state.
func (s *Store) WriteMetadata(metaDir string) error {
	skillsDir := filepath.Join(metaDir, "skills")
	groupsDir := filepath.Join(metaDir, "groups")
	os.MkdirAll(skillsDir, 0755)
	os.MkdirAll(groupsDir, 0755)

	skills, err := s.ListSkills()
	if err != nil {
		return err
	}
	for _, sk := range skills {
		m := skillMeta{
			ID:          sk.ID,
			Name:        sk.Name,
			Description: sk.Description,
			SourceType:  sk.SourceType,
			SourceRef:   sk.SourceRef,
			ContentHash: sk.ContentHash,
		}
		if err := writeJSON(filepath.Join(skillsDir, sk.Name+".json"), m); err != nil {
			return err
		}
	}

	groups, err := s.ListGroups()
	if err != nil {
		return err
	}
	for _, g := range groups {
		gSkills, _ := s.ListGroupSkills(g.ID)
		var skillNames []string
		for _, sk := range gSkills {
			skillNames = append(skillNames, sk.Name)
		}
		m := groupMeta{
			ID:          g.ID,
			Name:        g.Name,
			Description: g.Description,
			Skills:      skillNames,
		}
		if err := writeJSON(filepath.Join(groupsDir, g.Name+".json"), m); err != nil {
			return err
		}
	}

	return nil
}

func writeJSON(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
