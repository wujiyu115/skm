package server

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ejoy/skm/internal/agent"
)

func (s *Server) registerAgentRoutes(api fiber.Router) {
	api.Get("/agents", s.listAgents)
	api.Post("/agents", s.addAgent)
	api.Delete("/agents/:name", s.removeAgent)
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

	var result []agentInfo
	for _, a := range agent.Builtin() {
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
