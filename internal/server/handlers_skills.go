package server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/ejoy/skm/internal/agent"
	"github.com/ejoy/skm/internal/fetcher"
	"github.com/ejoy/skm/internal/logger"
	"github.com/ejoy/skm/internal/skill"
	"github.com/ejoy/skm/internal/store"
	skmsync "github.com/ejoy/skm/internal/sync"
)

func (s *Server) registerSkillRoutes(api fiber.Router) {
	api.Get("/skills", s.listSkills)
	api.Post("/skills/install", s.installSkill)
	api.Get("/skills/:id", s.getSkill)
	api.Delete("/skills/:id", s.deleteSkill)
	api.Get("/skills/:id/content", s.getSkillContent)
	api.Put("/skills/:id/enable", s.setSkillEnabled)
	api.Post("/skills/:id/sync", s.syncSkill)
}

func (s *Server) listSkills(c *fiber.Ctx) error {
	skills, err := s.store.ListSkills()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	type skillWithTargets struct {
		*skill.Skill
		Targets []store.Target `json:"targets"`
	}

	result := make([]skillWithTargets, 0, len(skills))
	for _, sk := range skills {
		targets, _ := s.store.ListTargets(sk.ID)
		result = append(result, skillWithTargets{Skill: sk, Targets: targets})
	}

	return c.JSON(result)
}

func (s *Server) getSkill(c *fiber.Ctx) error {
	sk, err := s.store.GetSkillByID(c.Params("id"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	targets, _ := s.store.ListTargets(sk.ID)
	return c.JSON(fiber.Map{"skill": sk, "targets": targets})
}

func (s *Server) deleteSkill(c *fiber.Ctx) error {
	sk, err := s.store.GetSkillByID(c.Params("id"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Info("delete skill", "id", sk.ID, "name", sk.Name)
	targets, _ := s.store.ListTargets(sk.ID)
	for _, t := range targets {
		skmsync.Unsync(t.TargetPath)
	}
	os.RemoveAll(sk.CentralPath)
	if err := s.store.DeleteSkill(sk.ID); err != nil {
		logger.Error("delete skill failed", "id", sk.ID, "err", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	s.writeMetadata()
	return c.JSON(fiber.Map{"ok": true})
}

func (s *Server) installSkill(c *fiber.Ctx) error {
	var req struct {
		Source string   `json:"source"`
		Agents []string `json:"agents"`
		Global bool     `json:"global"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	logger.Info("install skill", "source", req.Source, "agents", req.Agents, "global", req.Global)

	src := fetcher.Parse(req.Source)

	var baseDir string
	if src.Type == "local" {
		absPath, err := filepath.Abs(src.Subpath)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		baseDir = absPath
	} else {
		var err error
		baseDir, err = fetcher.CloneToTemp(src)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer fetcher.CleanupTemp(baseDir)
	}

	skills, err := skill.Discover(baseDir, src.Subpath)
	if err != nil || len(skills) == 0 {
		logger.Warn("no skills found", "source", req.Source, "baseDir", baseDir, "subpath", src.Subpath, "err", err)
		return c.Status(404).JSON(fiber.Map{"error": "no skills found"})
	}

	targetAgents := agent.Resolve(req.Agents)
	cwd, _ := os.Getwd()
	var installed []string

	for _, sk := range skills {
		centralPath := filepath.Join(s.cfg.SkillsDir, sk.Name)
		os.RemoveAll(centralPath)
		fetcher.CopyLocal(sk.CentralPath, centralPath)

		hash, _ := skill.HashDir(centralPath)
		sk.ID = uuid.New().String()
		sk.CentralPath = centralPath
		sk.ContentHash = hash
		sk.SourceType = "git"
		sk.SourceRef = src.GitURL()
		if src.Type == "local" {
			sk.SourceType = "local"
			sk.SourceRef = req.Source
		}

		// DeleteSkillByName may return "not found" on first install — ignore that.
		s.store.DeleteSkillByName(sk.Name)
		if err := s.store.InsertSkill(sk); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("save skill %s: %v", sk.Name, err)})
		}

		for _, ag := range targetAgents {
			targetDir, err := agent.InstallPath(ag, req.Global, cwd)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("install path: %v", err)})
			}
			targetPath := filepath.Join(targetDir, sk.Name)
			if err := skmsync.SyncSkill(centralPath, targetPath, "symlink"); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("sync %s: %v", sk.Name, err)})
			}
			if err := s.store.UpsertTarget(sk.ID, ag.Name, targetPath, "symlink", hash); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("upsert target: %v", err)})
			}
		}
		installed = append(installed, sk.Name)
	}

	logger.Info("skills installed", "count", len(installed), "skills", installed)
	s.writeMetadata()
	return c.JSON(fiber.Map{"installed": installed})
}

func (s *Server) setSkillEnabled(c *fiber.Ctx) error {
	id := c.Params("id")

	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	sk, err := s.store.GetSkillByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	if err := s.store.SetSkillEnabled(id, req.Enabled); err != nil {
		logger.Error("set skill enabled failed", "id", id, "err", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	action := "enable"
	if !req.Enabled {
		action = "disable"
	}
	if err := s.store.InsertAuditLog(action, sk.Name, ""); err != nil {
		logger.Warn("audit log failed", "err", err)
	}

	logger.Info("skill "+action+"d", "id", id, "name", sk.Name)
	return c.JSON(fiber.Map{"ok": true})
}

func (s *Server) getSkillContent(c *fiber.Ctx) error {
	sk, err := s.store.GetSkillByID(c.Params("id"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	mdPath := filepath.Join(sk.CentralPath, "SKILL.md")
	content, err := os.ReadFile(mdPath)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "SKILL.md not found"})
	}
	return c.JSON(fiber.Map{"content": string(content)})
}

func (s *Server) syncSkill(c *fiber.Ctx) error {
	sk, err := s.store.GetSkillByID(c.Params("id"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	var req struct {
		Agents []string `json:"agents"`
		Global bool     `json:"global"`
	}
	c.BodyParser(&req)

	targetAgents := agent.Resolve(req.Agents)
	cwd, _ := os.Getwd()

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

	if err := s.store.InsertAuditLog("sync", sk.Name, ""); err != nil {
		logger.Warn("audit log failed", "err", err)
	}

	return c.JSON(fiber.Map{"ok": true})
}

