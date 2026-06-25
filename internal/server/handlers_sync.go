package server

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"

	"github.com/ejoy/skm/internal/agent"
	"github.com/ejoy/skm/internal/logger"
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

	logger.Info("sync triggered", "agents", req.Agents, "dryRun", req.DryRun, "global", req.Global)
	skills, _ := s.store.ListSkills()
	targetAgents := agent.Resolve(req.Agents)
	cwd, _ := os.Getwd()

	logger.Debug("sync context", "cwd", cwd, "skillCount", len(skills), "agentCount", len(targetAgents))
	if len(targetAgents) == 0 {
		logger.Warn("sync: no agents resolved", "requestedAgents", req.Agents)
	}
	for _, ag := range targetAgents {
		logger.Debug("resolved agent", "name", ag.Name, "projectDir", ag.ProjectDir, "globalDir", ag.GlobalDir)
	}

	type result struct {
		Skill  string `json:"skill"`
		Agent  string `json:"agent"`
		Status string `json:"status"`
		Detail string `json:"detail,omitempty"`
	}
	var results []result

	for _, sk := range skills {
		if !sk.Enabled {
			logger.Debug("skipping disabled skill", "skill", sk.Name)
			continue
		}
		logger.Debug("syncing skill", "skill", sk.Name, "centralPath", sk.CentralPath)
		for _, ag := range targetAgents {
			targetDir, err := agent.InstallPath(ag, req.Global, cwd)
			if err != nil {
				logger.Error("install path failed", "skill", sk.Name, "agent", ag.Name, "err", err)
				results = append(results, result{sk.Name, ag.Name, "error", err.Error()})
				continue
			}
			targetPath := filepath.Join(targetDir, sk.Name)
			logger.Debug("target path", "skill", sk.Name, "agent", ag.Name, "targetPath", targetPath)

			if skmsync.IsCurrent(sk.CentralPath, targetPath, "symlink") {
				logger.Debug("already current", "skill", sk.Name, "agent", ag.Name)
				results = append(results, result{sk.Name, ag.Name, "current", ""})
				continue
			}
			if req.DryRun {
				results = append(results, result{sk.Name, ag.Name, "would_sync", ""})
				continue
			}
			if err := skmsync.SyncSkill(sk.CentralPath, targetPath, "symlink"); err != nil {
				logger.Error("sync skill failed", "skill", sk.Name, "agent", ag.Name, "source", sk.CentralPath, "target", targetPath, "err", err)
				results = append(results, result{sk.Name, ag.Name, "error", err.Error()})
				continue
			}
			if err := s.store.UpsertTarget(sk.ID, ag.Name, targetPath, "symlink", sk.ContentHash); err != nil {
				logger.Error("upsert target failed", "skill", sk.Name, "agent", ag.Name, "err", err)
				results = append(results, result{sk.Name, ag.Name, "error", err.Error()})
				continue
			}
			logger.Info("skill synced", "skill", sk.Name, "agent", ag.Name, "target", targetPath)
			results = append(results, result{sk.Name, ag.Name, "synced", ""})
		}
	}

	return c.JSON(results)
}
