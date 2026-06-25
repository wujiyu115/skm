package server

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) registerSettingRoutes(api fiber.Router) {
	api.Get("/settings", s.getSettings)
	api.Put("/settings", s.updateSettings)
	api.Put("/settings/:key", s.updateSingleSetting)
}

func (s *Server) getSettings(c *fiber.Ctx) error {
	// Start with config-based settings.
	result := map[string]string{
		"skills_dir": s.cfg.SkillsDir,
		"cache_dir":  s.cfg.CacheDir,
	}

	// Merge DB settings on top.
	dbSettings, err := s.store.ListSettings()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	for k, v := range dbSettings {
		result[k] = v
	}

	return c.JSON(result)
}

func (s *Server) updateSettings(c *fiber.Ctx) error {
	var body map[string]string
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid JSON body: " + err.Error(),
		})
	}

	if err := s.store.SetSettings(body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"ok": true})
}

func (s *Server) updateSingleSetting(c *fiber.Ctx) error {
	key := c.Params("key")

	var body struct {
		Value string `json:"value"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid JSON body: " + err.Error(),
		})
	}

	if err := s.store.SetSetting(key, body.Value); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"ok": true})
}
