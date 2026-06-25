package server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/ejoy/skm/internal/agent"
	"github.com/ejoy/skm/internal/store"
	skmsync "github.com/ejoy/skm/internal/sync"
)

func (s *Server) registerGroupRoutes(api fiber.Router) {
	api.Get("/groups", s.listGroups)
	api.Post("/groups", s.createGroup)
	api.Get("/groups/:id", s.getGroup)
	api.Put("/groups/:id", s.updateGroup)
	api.Delete("/groups/:id", s.deleteGroup)
	api.Post("/groups/:id/skills", s.addGroupSkills)
	api.Delete("/groups/:id/skills/:sid", s.removeGroupSkill)
	api.Post("/groups/:id/install", s.installGroup)
}

func (s *Server) listGroups(c *fiber.Ctx) error {
	groups, err := s.store.ListGroups()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	type groupWithCount struct {
		store.Group
		SkillCount int `json:"skill_count"`
	}
	result := make([]groupWithCount, 0, len(groups))
	for _, g := range groups {
		skills, _ := s.store.ListGroupSkills(g.ID)
		result = append(result, groupWithCount{Group: g, SkillCount: len(skills)})
	}
	return c.JSON(result)
}

func (s *Server) createGroup(c *fiber.Ctx) error {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&req); err != nil || req.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name required"})
	}
	id := uuid.New().String()
	if err := s.store.InsertGroup(id, req.Name, req.Description); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	s.writeMetadata()
	return c.Status(201).JSON(fiber.Map{"id": id, "name": req.Name})
}

func (s *Server) getGroup(c *fiber.Ctx) error {
	g, err := s.store.GetGroupByID(c.Params("id"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	skills, _ := s.store.ListGroupSkills(g.ID)
	return c.JSON(fiber.Map{"group": g, "skills": skills})
}

func (s *Server) updateGroup(c *fiber.Ctx) error {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	if err := s.store.UpdateGroup(c.Params("id"), req.Name, req.Description); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	s.writeMetadata()
	return c.JSON(fiber.Map{"ok": true})
}

func (s *Server) deleteGroup(c *fiber.Ctx) error {
	if err := s.store.DeleteGroup(c.Params("id")); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	s.writeMetadata()
	return c.JSON(fiber.Map{"ok": true})
}

func (s *Server) addGroupSkills(c *fiber.Ctx) error {
	var req struct {
		SkillIDs []string `json:"skill_ids"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	for i, sid := range req.SkillIDs {
		if err := s.store.AddSkillToGroup(c.Params("id"), sid, i); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("add skill to group: %v", err)})
		}
	}
	s.writeMetadata()
	return c.JSON(fiber.Map{"ok": true})
}

func (s *Server) removeGroupSkill(c *fiber.Ctx) error {
	if err := s.store.RemoveSkillFromGroup(c.Params("id"), c.Params("sid")); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("remove skill from group: %v", err)})
	}
	s.writeMetadata()
	return c.JSON(fiber.Map{"ok": true})
}

func (s *Server) installGroup(c *fiber.Ctx) error {
	var req struct {
		Agents []string `json:"agents"`
		Global bool     `json:"global"`
	}
	c.BodyParser(&req)

	skills, err := s.store.ListGroupSkills(c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	targetAgents := agent.Resolve(req.Agents)
	cwd, _ := os.Getwd()

	for _, sk := range skills {
		for _, ag := range targetAgents {
			targetDir, err := agent.InstallPath(ag, req.Global, cwd)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("install path: %v", err)})
			}
			targetPath := filepath.Join(targetDir, sk.Name)
			if err := skmsync.SyncSkill(sk.CentralPath, targetPath, "symlink"); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("sync %s: %v", sk.Name, err)})
			}
			if err := s.store.UpsertTarget(sk.ID, ag.Name, targetPath, "symlink", sk.ContentHash); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("upsert target: %v", err)})
			}
		}
	}
	return c.JSON(fiber.Map{"ok": true})
}
