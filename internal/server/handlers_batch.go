package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/ejoy/skm/internal/agent"
	"github.com/ejoy/skm/internal/logger"
	skmsync "github.com/ejoy/skm/internal/sync"
)

func (s *Server) registerBatchRoutes(api fiber.Router) {
	api.Post("/skills/batch/delete", s.batchDelete)
	api.Post("/skills/batch/enable", s.batchEnable)
	api.Post("/skills/batch/tag", s.batchTag)
	api.Post("/skills/batch/sync", s.batchSync)
}

func (s *Server) batchDelete(c *fiber.Ctx) error {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}
	if len(req.IDs) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ids is required"})
	}

	processed := 0
	errors := []string{}

	for _, id := range req.IDs {
		sk, err := s.store.GetSkillByID(id)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", id, err))
			continue
		}

		if err := s.store.DeleteSkill(sk.ID); err != nil {
			logger.Error("batch delete skill failed", "id", sk.ID, "err", err)
			errors = append(errors, fmt.Sprintf("%s: %v", id, err))
			continue
		}

		targets, _ := s.store.ListTargets(sk.ID)
		for _, t := range targets {
			if err := skmsync.Unsync(t.TargetPath); err != nil {
				logger.Warn("batch delete unsync failed", "target", t.TargetPath, "err", err)
			}
		}
		if err := os.RemoveAll(sk.CentralPath); err != nil {
			logger.Warn("batch delete remove files failed", "path", sk.CentralPath, "err", err)
		}

		if err := s.store.InsertAuditLog("batch_delete", sk.Name, ""); err != nil {
			logger.Warn("audit log failed", "err", err)
		}
		processed++
	}

	s.writeMetadata()
	return c.JSON(fiber.Map{"ok": true, "processed": processed, "errors": errors})
}

func (s *Server) batchEnable(c *fiber.Ctx) error {
	var req struct {
		IDs     []string `json:"ids"`
		Enabled bool     `json:"enabled"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}
	if len(req.IDs) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ids is required"})
	}

	processed := 0
	errors := []string{}

	action := "batch_enable"
	if !req.Enabled {
		action = "batch_disable"
	}

	for _, id := range req.IDs {
		sk, err := s.store.GetSkillByID(id)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", id, err))
			continue
		}

		if err := s.store.SetSkillEnabled(id, req.Enabled); err != nil {
			logger.Error("batch set enabled failed", "id", id, "err", err)
			errors = append(errors, fmt.Sprintf("%s: %v", id, err))
			continue
		}

		if err := s.store.InsertAuditLog(action, sk.Name, ""); err != nil {
			logger.Warn("audit log failed", "err", err)
		}
		processed++
	}

	return c.JSON(fiber.Map{"ok": true, "processed": processed, "errors": errors})
}

func (s *Server) batchTag(c *fiber.Ctx) error {
	var req struct {
		IDs    []string `json:"ids"`
		Tags   []string `json:"tags"`
		Action string   `json:"action"` // "add" or "remove"
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}
	if len(req.IDs) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ids is required"})
	}
	if len(req.Tags) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "tags is required"})
	}
	if req.Action != "add" && req.Action != "remove" {
		return c.Status(400).JSON(fiber.Map{"error": "action must be 'add' or 'remove'"})
	}

	processed := 0
	errors := []string{}

	for _, id := range req.IDs {
		sk, err := s.store.GetSkillByID(id)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", id, err))
			continue
		}

		if req.Action == "add" {
			for _, tag := range req.Tags {
				if err := s.store.AddTag(id, tag); err != nil {
					logger.Warn("batch add tag failed", "skill", id, "tag", tag, "err", err)
				}
			}
		} else {
			existing, _ := s.store.ListSkillTags(id)
			newTags := make([]string, 0)
			for _, t := range existing {
				keep := true
				for _, rt := range req.Tags {
					if t == rt {
						keep = false
						break
					}
				}
				if keep {
					newTags = append(newTags, t)
				}
			}
			if err := s.store.SetSkillTags(id, newTags); err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", id, err))
				continue
			}
		}

		auditAction := "batch_tag_add"
		if req.Action == "remove" {
			auditAction = "batch_tag_remove"
		}
		if err := s.store.InsertAuditLog(auditAction, sk.Name, strings.Join(req.Tags, ",")); err != nil {
			logger.Warn("audit log failed", "err", err)
		}
		processed++
	}

	return c.JSON(fiber.Map{"ok": true, "processed": processed, "errors": errors})
}

func (s *Server) batchSync(c *fiber.Ctx) error {
	var req struct {
		IDs    []string `json:"ids"`
		Agents []string `json:"agents"`
		Global bool     `json:"global"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}
	if len(req.IDs) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ids is required"})
	}

	targetAgents := agent.Resolve(req.Agents)
	cwd, _ := os.Getwd()

	processed := 0
	errors := []string{}

	for _, id := range req.IDs {
		sk, err := s.store.GetSkillByID(id)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", id, err))
			continue
		}

		failed := false
		for _, ag := range targetAgents {
			targetDir, err := agent.InstallPath(ag, req.Global, cwd)
			if err != nil {
				errors = append(errors, fmt.Sprintf("%s/%s: %v", id, ag.Name, err))
				failed = true
				continue
			}
			targetPath := filepath.Join(targetDir, sk.Name)
			if err := skmsync.SyncSkill(sk.CentralPath, targetPath, "symlink"); err != nil {
				errors = append(errors, fmt.Sprintf("%s/%s: %v", id, ag.Name, err))
				failed = true
				continue
			}
			if err := s.store.UpsertTarget(sk.ID, ag.Name, targetPath, "symlink", sk.ContentHash); err != nil {
				errors = append(errors, fmt.Sprintf("%s/%s: %v", id, ag.Name, err))
				failed = true
				continue
			}
		}

		if !failed {
			processed++
		}

		agentNames := make([]string, len(targetAgents))
		for i, ag := range targetAgents {
			agentNames[i] = ag.Name
		}
		if err := s.store.InsertAuditLog("batch_sync", sk.Name, strings.Join(agentNames, ",")); err != nil {
			logger.Warn("audit log failed", "err", err)
		}
	}

	return c.JSON(fiber.Map{"ok": true, "processed": processed, "errors": errors})
}
