package server

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/ejoy/skm/internal/logger"
)

func (s *Server) registerTagRoutes(api fiber.Router) {
	api.Get("/tags", s.listTags)
	api.Post("/tags/rename", s.renameTag)
	api.Delete("/tags/:tag", s.deleteTag)

	api.Get("/skills/:id/tags", s.listSkillTags)
	api.Put("/skills/:id/tags", s.setSkillTags)
}

func (s *Server) listTags(c *fiber.Ctx) error {
	tags, err := s.store.ListTags()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if tags == nil {
		tags = []string{}
	}
	return c.JSON(fiber.Map{"tags": tags})
}

func (s *Server) listSkillTags(c *fiber.Ctx) error {
	id := c.Params("id")

	// Verify the skill exists.
	if _, err := s.store.GetSkillByID(id); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	tags, err := s.store.ListSkillTags(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if tags == nil {
		tags = []string{}
	}
	return c.JSON(fiber.Map{"tags": tags})
}

func (s *Server) setSkillTags(c *fiber.Ctx) error {
	id := c.Params("id")

	// Verify the skill exists.
	if _, err := s.store.GetSkillByID(id); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	var req struct {
		Tags []string `json:"tags"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := s.store.SetSkillTags(id, req.Tags); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if err := s.store.InsertAuditLog("tags_set", id, strings.Join(req.Tags, ",")); err != nil {
		logger.Warn("audit log failed", "err", err)
	}

	return c.JSON(fiber.Map{"ok": true})
}

func (s *Server) renameTag(c *fiber.Ctx) error {
	var req struct {
		Old string `json:"old"`
		New string `json:"new"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.Old == "" || req.New == "" {
		return c.Status(400).JSON(fiber.Map{"error": "both 'old' and 'new' fields are required"})
	}

	if err := s.store.RenameTag(req.Old, req.New); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if err := s.store.InsertAuditLog("tag_rename", req.Old, fmt.Sprintf("renamed to %s", req.New)); err != nil {
		logger.Warn("audit log failed", "err", err)
	}

	return c.JSON(fiber.Map{"ok": true})
}

func (s *Server) deleteTag(c *fiber.Ctx) error {
	tag := c.Params("tag")

	if err := s.store.DeleteTag(tag); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if err := s.store.InsertAuditLog("tag_delete", tag, ""); err != nil {
		logger.Warn("audit log failed", "err", err)
	}

	return c.JSON(fiber.Map{"ok": true})
}
