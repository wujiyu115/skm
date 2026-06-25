package server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"

	"github.com/ejoy/skm/internal/agent"
	"github.com/ejoy/skm/internal/logger"
	"github.com/ejoy/skm/internal/project"
)

func (s *Server) registerAgentRoutes(api fiber.Router) {
	api.Get("/agents", s.listAgents)
	api.Post("/agents", s.addAgent)
	api.Delete("/agents/:name", s.removeAgent)
	api.Get("/agents/:name/skills", s.listAgentGlobalSkills)
	api.Post("/agents/:name/skills/add", s.addAgentGlobalSkill)
	api.Put("/agents/:name/skills/toggle", s.toggleAgentGlobalSkill)
	api.Delete("/agents/:name/skills", s.deleteAgentGlobalSkill)
}

func (s *Server) listAgents(c *fiber.Ctx) error {
	detected := agent.Detect()
	detectedMap := map[string]bool{}
	for _, a := range detected {
		detectedMap[a.Name] = true
	}

	type agentInfo struct {
		agent.Adapter
		Detected bool `json:"detected"`
	}

	builtins := agent.Builtin()
	result := make([]agentInfo, 0, len(builtins))
	for _, a := range builtins {
		result = append(result, agentInfo{Adapter: a, Detected: detectedMap[a.Name]})
	}
	return c.JSON(result)
}

func (s *Server) addAgent(c *fiber.Ctx) error {
	var a agent.Adapter
	if err := c.BodyParser(&a); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	if err := s.store.InsertAgent(a); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"ok": true})
}

func (s *Server) removeAgent(c *fiber.Ctx) error {
	s.store.DeleteAgent(c.Params("name"))
	return c.JSON(fiber.Map{"ok": true})
}

func (s *Server) agentGlobalDir(name string) (agent.Adapter, string, error) {
	a, ok := agent.Find(name)
	if !ok {
		return agent.Adapter{}, "", fmt.Errorf("unknown agent: %s", name)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return agent.Adapter{}, "", err
	}
	return a, filepath.Join(home, a.GlobalDir), nil
}

func (s *Server) listAgentGlobalSkills(c *fiber.Ctx) error {
	a, globalDir, err := s.agentGlobalDir(c.Params("name"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	home, _ := os.UserHomeDir()
	skills, err := project.ScanSkills(home, []agent.Adapter{a})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	_ = globalDir
	return c.JSON(skills)
}

func (s *Server) addAgentGlobalSkill(c *fiber.Ctx) error {
	a, _, err := s.agentGlobalDir(c.Params("name"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	var req struct {
		SkillID string `json:"skill_id"`
	}
	if err := c.BodyParser(&req); err != nil || req.SkillID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "skill_id required"})
	}

	sk, err := s.store.GetSkillByID(req.SkillID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	home, _ := os.UserHomeDir()
	if err := project.AddFromCentral(sk.CentralPath, home, a.GlobalDir, sk.Name); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if err := s.store.InsertAuditLog("global_add_skill", sk.Name, fmt.Sprintf("agent=%s", a.Name)); err != nil {
		logger.Warn("audit log failed", "err", err)
	}

	return c.JSON(fiber.Map{"ok": true})
}

func (s *Server) toggleAgentGlobalSkill(c *fiber.Ctx) error {
	a, _, err := s.agentGlobalDir(c.Params("name"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	var req struct {
		SkillName string `json:"skill_name"`
		Enabled   bool   `json:"enabled"`
	}
	if err := c.BodyParser(&req); err != nil || req.SkillName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "skill_name required"})
	}

	home, _ := os.UserHomeDir()
	if err := project.ToggleSkill(home, a.GlobalDir, req.SkillName, req.Enabled); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	action := "global_enable_skill"
	if !req.Enabled {
		action = "global_disable_skill"
	}
	if err := s.store.InsertAuditLog(action, req.SkillName, fmt.Sprintf("agent=%s", a.Name)); err != nil {
		logger.Warn("audit log failed", "err", err)
	}

	return c.JSON(fiber.Map{"ok": true})
}

func (s *Server) deleteAgentGlobalSkill(c *fiber.Ctx) error {
	a, globalDir, err := s.agentGlobalDir(c.Params("name"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	var req struct {
		SkillPath string `json:"skill_path"`
	}
	if err := c.BodyParser(&req); err != nil || req.SkillPath == "" {
		return c.Status(400).JSON(fiber.Map{"error": "skill_path required"})
	}

	cleanPath := filepath.Clean(req.SkillPath)
	cleanGlobal := filepath.Clean(globalDir)
	disabledGlobal := filepath.Clean(project.DisabledDirFor(a.GlobalDir))
	home, _ := os.UserHomeDir()
	cleanDisabled := filepath.Clean(filepath.Join(home, disabledGlobal))

	if !hasPrefix(cleanPath, cleanGlobal) && !hasPrefix(cleanPath, cleanDisabled) {
		return c.Status(403).JSON(fiber.Map{"error": "skill_path is outside agent directory"})
	}

	if err := project.RemoveSkill(cleanPath); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if err := s.store.InsertAuditLog("global_delete_skill", filepath.Base(cleanPath), fmt.Sprintf("agent=%s", a.Name)); err != nil {
		logger.Warn("audit log failed", "err", err)
	}

	return c.JSON(fiber.Map{"ok": true})
}

func hasPrefix(path, prefix string) bool {
	return path == prefix || len(path) > len(prefix) && path[:len(prefix)+1] == prefix+string(filepath.Separator)
}
