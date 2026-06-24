package server

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) registerSettingRoutes(api fiber.Router) {
	api.Get("/settings", s.getSettings)
	api.Put("/settings", s.updateSettings)
}

func (s *Server) getSettings(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"skills_dir": s.cfg.SkillsDir,
		"cache_dir":  s.cfg.CacheDir,
	})
}

func (s *Server) updateSettings(c *fiber.Ctx) error {
	// TODO: Implement settings persistence when needed.
	return c.JSON(fiber.Map{"ok": true})
}
