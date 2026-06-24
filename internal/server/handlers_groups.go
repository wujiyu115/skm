package server

import (
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
	var result []groupWithCount
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
	return c.JSON(fiber.Map{"ok": true})
}

func (s *Server) deleteGroup(c *fiber.Ctx) error {
	if err := s.store.DeleteGroup(c.Params("id")); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
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
		s.store.AddSkillToGroup(c.Params("id"), sid, i)
	}
	return c.JSON(fiber.Map{"ok": true})
}

func (s *Server) removeGroupSkill(c *fiber.Ctx) error {
	s.store.RemoveSkillFromGroup(c.Params("id"), c.Params("sid"))
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

	targetAgents := resolveAgentsForAPI(req.Agents)
	cwd, _ := os.Getwd()

	for _, sk := range skills {
		for _, ag := range targetAgents {
			targetDir := agent.InstallPath(ag, req.Global, cwd)
			targetPath := filepath.Join(targetDir, sk.Name)
			skmsync.SyncSkill(sk.CentralPath, targetPath, "symlink")
			s.store.UpsertTarget(sk.ID, ag.Name, targetPath, "symlink", sk.ContentHash)
		}
	}
	return c.JSON(fiber.Map{"ok": true})
}
