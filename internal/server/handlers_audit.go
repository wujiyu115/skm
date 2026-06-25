package server

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/ejoy/skm/internal/store"
)

func (s *Server) registerAuditRoutes(api fiber.Router) {
	api.Get("/audit", s.listAuditLog)
	api.Delete("/audit", s.pruneAuditLog)
}

func (s *Server) listAuditLog(c *fiber.Ctx) error {
	limit := 100
	if q := c.Query("limit"); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n > 0 {
			if n > 10000 {
				n = 10000
			}
			limit = n
		}
	}

	entries, err := s.store.ListAuditLog(limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if entries == nil {
		entries = []store.AuditEntry{}
	}
	return c.JSON(entries)
}

func (s *Server) pruneAuditLog(c *fiber.Ctx) error {
	if err := s.store.PruneAuditLog(10000); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"ok": true, "pruned": true})
}
