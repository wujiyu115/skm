package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/ejoy/skm/internal/store"
)

// Server wraps a Fiber app and provides the REST API for SKM.
type Server struct {
	app   *fiber.App
	store *store.Store
	cfg   *ServerConfig
}

// ServerConfig holds the dependencies the server needs.
type ServerConfig struct {
	Store     *store.Store
	SkillsDir string
	CacheDir  string
}

// New creates a new Server with all API routes registered.
func New(cfg *ServerConfig) *Server {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(cors.New())

	s := &Server{
		app:   app,
		store: cfg.Store,
		cfg:   cfg,
	}

	api := app.Group("/api")
	s.registerSkillRoutes(api)
	s.registerGroupRoutes(api)
	s.registerAgentRoutes(api)
	s.registerSyncRoutes(api)
	s.registerSettingRoutes(api)

	// TODO: Task 13 will add go:embed static file serving here.

	return s
}

// Start begins listening on the given port.
func (s *Server) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("SKM web UI running at http://localhost:%d\n", port)
	return s.app.Listen(addr)
}

// App returns the underlying Fiber app (useful for testing).
func (s *Server) App() *fiber.App {
	return s.app
}
