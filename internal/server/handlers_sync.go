package server

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"

	"github.com/ejoy/skm/internal/agent"
	skmsync "github.com/ejoy/skm/internal/sync"
)

func (s *Server) registerSyncRoutes(api fiber.Router) {
	api.Get("/sync/status", s.syncStatus)
	api.Post("/sync", s.triggerSync)
}

func (s *Server) syncStatus(c *fiber.Ctx) error {
	skills, _ := s.store.ListSkills()
	total, synced, stale := 0, 0, 0

	for _, sk := range skills {
		targets, _ := s.store.ListTargets(sk.ID)
		for _, t := range targets {
			total++
			if skmsync.IsCurrent(sk.CentralPath, t.TargetPath, t.Mode) {
				synced++
			} else {
				stale++
			}
		}
	}

	return c.JSON(fiber.Map{"total": total, "synced": synced, "stale": stale})
}

func (s *Server) triggerSync(c *fiber.Ctx) error {
	var req struct {
		Agents []string `json:"agents"`
		DryRun bool     `json:"dry_run"`
		Global bool     `json:"global"`
	}
	c.BodyParser(&req)

	skills, _ := s.store.ListSkills()
	targetAgents := resolveAgentsForAPI(req.Agents)
	cwd, _ := os.Getwd()

	type result struct {
		Skill  string `json:"skill"`
		Agent  string `json:"agent"`
		Status string `json:"status"`
	}
	var results []result

	for _, sk := range skills {
		if !sk.Enabled {
			continue
		}
		for _, ag := range targetAgents {
			targetDir := agent.InstallPath(ag, req.Global, cwd)
			targetPath := filepath.Join(targetDir, sk.Name)

			if skmsync.IsCurrent(sk.CentralPath, targetPath, "symlink") {
				results = append(results, result{sk.Name, ag.Name, "current"})
				continue
			}
			if req.DryRun {
				results = append(results, result{sk.Name, ag.Name, "would_sync"})
				continue
			}
			if err := skmsync.SyncSkill(sk.CentralPath, targetPath, "symlink"); err != nil {
				results = append(results, result{sk.Name, ag.Name, "error"})
				continue
			}
			s.store.UpsertTarget(sk.ID, ag.Name, targetPath, "symlink", sk.ContentHash)
			results = append(results, result{sk.Name, ag.Name, "synced"})
		}
	}

	return c.JSON(results)
}
