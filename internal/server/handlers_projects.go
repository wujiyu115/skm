package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/ejoy/skm/internal/agent"
	"github.com/ejoy/skm/internal/logger"
	"github.com/ejoy/skm/internal/project"
)

func (s *Server) registerProjectRoutes(api fiber.Router) {
	api.Get("/projects", s.listProjects)
	api.Post("/projects", s.createProject)
	api.Delete("/projects/:id", s.deleteProject)
	api.Get("/projects/:id/skills", s.listProjectSkills)
	api.Post("/projects/:id/skills/add", s.addProjectSkill)
	api.Put("/projects/:id/skills/toggle", s.toggleProjectSkill)
	api.Delete("/projects/:id/skills", s.deleteProjectSkill)
}

func (s *Server) listProjects(c *fiber.Ctx) error {
	projects, err := s.store.ListProjects()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(projects)
}

func (s *Server) createProject(c *fiber.Ctx) error {
	var req struct {
		Path string `json:"path"`
	}
	if err := c.BodyParser(&req); err != nil || req.Path == "" {
		return c.Status(400).JSON(fiber.Map{"error": "path required"})
	}

	if !filepath.IsAbs(req.Path) {
		return c.Status(400).JSON(fiber.Map{"error": "path must be absolute"})
	}

	info, err := os.Stat(req.Path)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("path does not exist: %s", req.Path)})
	}
	if !info.IsDir() {
		return c.Status(400).JSON(fiber.Map{"error": "path is not a directory"})
	}

	id := uuid.New().String()
	name := filepath.Base(req.Path)

	if err := s.store.InsertProject(id, name, req.Path); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if err := s.store.InsertAuditLog("project_create", name, fmt.Sprintf("path=%s", req.Path)); err != nil {
		logger.Warn("audit log failed", "err", err)
	}

	return c.Status(201).JSON(fiber.Map{"id": id, "name": name, "path": req.Path})
}

func (s *Server) deleteProject(c *fiber.Ctx) error {
	id := c.Params("id")
	proj, _ := s.store.GetProjectByID(id)

	if err := s.store.DeleteProject(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	target := id
	if proj != nil {
		target = proj.Name
	}
	if err := s.store.InsertAuditLog("project_delete", target, ""); err != nil {
		logger.Warn("audit log failed", "err", err)
	}

	return c.JSON(fiber.Map{"ok": true})
}

func (s *Server) listProjectSkills(c *fiber.Ctx) error {
	proj, err := s.store.GetProjectByID(c.Params("id"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	agents, err := s.store.ListAgents()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	skills, err := project.ScanSkills(proj.Path, agents)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(skills)
}

func (s *Server) addProjectSkill(c *fiber.Ctx) error {
	proj, err := s.store.GetProjectByID(c.Params("id"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	var req struct {
		SkillID string   `json:"skill_id"`
		Agents  []string `json:"agents"`
	}
	if err := c.BodyParser(&req); err != nil || req.SkillID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "skill_id required"})
	}

	sk, err := s.store.GetSkillByID(req.SkillID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	for _, agentName := range req.Agents {
		adapter, ok := agent.Find(agentName)
		if !ok {
			return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("unknown agent: %s", agentName)})
		}
		if err := project.AddFromCentral(sk.CentralPath, proj.Path, adapter.ProjectDir, sk.Name); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("add skill to %s: %v", agentName, err)})
		}
	}

	if err := s.store.InsertAuditLog("project_add_skill", sk.Name,
		fmt.Sprintf("project=%s agents=%s", proj.Name, strings.Join(req.Agents, ","))); err != nil {
		logger.Warn("audit log failed", "err", err)
	}

	return c.JSON(fiber.Map{"ok": true})
}

func (s *Server) toggleProjectSkill(c *fiber.Ctx) error {
	proj, err := s.store.GetProjectByID(c.Params("id"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	var req struct {
		Agent     string `json:"agent"`
		SkillName string `json:"skill_name"`
		Enabled   bool   `json:"enabled"`
	}
	if err := c.BodyParser(&req); err != nil || req.Agent == "" || req.SkillName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "agent and skill_name required"})
	}

	adapter, ok := agent.Find(req.Agent)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("unknown agent: %s", req.Agent)})
	}

	if err := project.ToggleSkill(proj.Path, adapter.ProjectDir, req.SkillName, req.Enabled); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	action := "project_enable_skill"
	if !req.Enabled {
		action = "project_disable_skill"
	}
	if err := s.store.InsertAuditLog(action, req.SkillName,
		fmt.Sprintf("project=%s agent=%s", proj.Name, req.Agent)); err != nil {
		logger.Warn("audit log failed", "err", err)
	}

	return c.JSON(fiber.Map{"ok": true})
}

func (s *Server) deleteProjectSkill(c *fiber.Ctx) error {
	proj, err := s.store.GetProjectByID(c.Params("id"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	var req struct {
		Agent     string `json:"agent"`
		SkillPath string `json:"skill_path"`
	}
	if err := c.BodyParser(&req); err != nil || req.SkillPath == "" {
		return c.Status(400).JSON(fiber.Map{"error": "skill_path required"})
	}

	// Security: prevent path traversal — skill_path must be within the project directory.
	cleanSkillPath := filepath.Clean(req.SkillPath)
	cleanProjectPath := filepath.Clean(proj.Path)
	if !strings.HasPrefix(cleanSkillPath, cleanProjectPath+string(filepath.Separator)) {
		return c.Status(403).JSON(fiber.Map{"error": "skill_path is outside project directory"})
	}

	if err := project.RemoveSkill(cleanSkillPath); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if err := s.store.InsertAuditLog("project_delete_skill", filepath.Base(cleanSkillPath),
		fmt.Sprintf("project=%s path=%s", proj.Name, cleanSkillPath)); err != nil {
		logger.Warn("audit log failed", "err", err)
	}

	return c.JSON(fiber.Map{"ok": true})
}
