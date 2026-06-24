package store

import (
	"testing"

	"github.com/ejoy/skm/internal/skill"
)

func TestInsertAndGetGroup(t *testing.T) {
	s := testStore(t)

	if err := s.InsertGroup("g1", "frontend", "Frontend skills"); err != nil {
		t.Fatalf("InsertGroup() error: %v", err)
	}

	g, err := s.GetGroup("frontend")
	if err != nil {
		t.Fatalf("GetGroup() error: %v", err)
	}
	if g.Name != "frontend" || g.Description != "Frontend skills" {
		t.Errorf("Group = %+v", g)
	}
}

func TestListGroups(t *testing.T) {
	s := testStore(t)
	s.InsertGroup("g1", "a-group", "")
	s.InsertGroup("g2", "b-group", "")

	groups, err := s.ListGroups()
	if err != nil {
		t.Fatalf("ListGroups() error: %v", err)
	}
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
}

func TestAddAndListGroupSkills(t *testing.T) {
	s := testStore(t)
	s.InsertGroup("g1", "my-group", "")
	s.InsertSkill(&skill.Skill{ID: "s1", Name: "skill-a", SourceType: "local", CentralPath: "/a", Enabled: true})
	s.InsertSkill(&skill.Skill{ID: "s2", Name: "skill-b", SourceType: "local", CentralPath: "/b", Enabled: true})

	s.AddSkillToGroup("g1", "s1", 0)
	s.AddSkillToGroup("g1", "s2", 1)

	skills, err := s.ListGroupSkills("g1")
	if err != nil {
		t.Fatalf("ListGroupSkills() error: %v", err)
	}
	if len(skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(skills))
	}
}

func TestDeleteGroup_CascadesGroupSkills(t *testing.T) {
	s := testStore(t)
	s.InsertGroup("g1", "cascade-test", "")
	s.InsertSkill(&skill.Skill{ID: "s1", Name: "sk", SourceType: "local", CentralPath: "/s", Enabled: true})
	s.AddSkillToGroup("g1", "s1", 0)

	if err := s.DeleteGroup("g1"); err != nil {
		t.Fatalf("DeleteGroup() error: %v", err)
	}

	skills, _ := s.ListGroupSkills("g1")
	if len(skills) != 0 {
		t.Error("group_skills should be cascaded on group delete")
	}
}

func TestRemoveSkillFromGroup(t *testing.T) {
	s := testStore(t)
	s.InsertGroup("g1", "remove-test", "")
	s.InsertSkill(&skill.Skill{ID: "s1", Name: "sk-rem", SourceType: "local", CentralPath: "/r", Enabled: true})
	s.AddSkillToGroup("g1", "s1", 0)

	if err := s.RemoveSkillFromGroup("g1", "s1"); err != nil {
		t.Fatalf("RemoveSkillFromGroup() error: %v", err)
	}

	skills, _ := s.ListGroupSkills("g1")
	if len(skills) != 0 {
		t.Errorf("expected 0 skills after remove, got %d", len(skills))
	}
}

func TestGetGroup_NotFound(t *testing.T) {
	s := testStore(t)
	_, err := s.GetGroup("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent group")
	}
}

func TestGetGroupByID(t *testing.T) {
	s := testStore(t)
	s.InsertGroup("byid-g1", "byid-group", "desc")

	g, err := s.GetGroupByID("byid-g1")
	if err != nil {
		t.Fatalf("GetGroupByID() error: %v", err)
	}
	if g.Name != "byid-group" {
		t.Errorf("Name = %q, want %q", g.Name, "byid-group")
	}
}

func TestUpdateGroup(t *testing.T) {
	s := testStore(t)
	s.InsertGroup("upd-1", "old-name", "old-desc")

	if err := s.UpdateGroup("upd-1", "new-name", "new-desc"); err != nil {
		t.Fatalf("UpdateGroup() error: %v", err)
	}

	g, err := s.GetGroup("new-name")
	if err != nil {
		t.Fatalf("GetGroup() after update error: %v", err)
	}
	if g.Description != "new-desc" {
		t.Errorf("Description = %q, want %q", g.Description, "new-desc")
	}
}
