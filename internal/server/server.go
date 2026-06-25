package server

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"

	"github.com/ejoy/skm/internal/logger"
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
	MetaDir   string
	DevMode   bool
}

// New creates a new Server with all API routes registered.
func New(cfg *ServerConfig) *Server {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3721, http://localhost:5173",
	}))

	s := &Server{
		app:   app,
		store: cfg.Store,
		cfg:   cfg,
	}

	api := app.Group("/api")
	s.registerSkillRoutes(api)
	s.registerTagRoutes(api)
	s.registerGroupRoutes(api)
	s.registerAgentRoutes(api)
	s.registerSyncRoutes(api)
	s.registerSettingRoutes(api)
	s.registerAuditRoutes(api)

	// Serve embedded web UI in production mode.
	// In dev mode (SKM_DEV=1), skip this so Vite HMR works via proxy.
	if !cfg.DevMode {
		sub, err := fs.Sub(webDist, "dist")
		if err == nil {
			app.Use("/", filesystem.New(filesystem.Config{
				Root:         http.FS(sub),
				Browse:       false,
				Index:        "index.html",
				NotFoundFile: "index.html",
			}))
		}
	}

	return s
}

// Start begins listening on the given port.
func (s *Server) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	mode := "production"
	if s.cfg.DevMode {
		mode = "development (API only, use Vite for frontend)"
	}
	logger.Info("server starting", "port", port, "mode", mode)
	fmt.Printf("SKM web UI running at http://localhost:%d [%s]\n", port, mode)
	return s.app.Listen(addr)
}

// App returns the underlying Fiber app (useful for testing).
func (s *Server) App() *fiber.App {
	return s.app
}

// writeMetadata writes JSON metadata mirror files. Errors are logged
// to stderr but do not fail the request.
func (s *Server) writeMetadata() {
	if s.cfg.MetaDir == "" {
		return
	}
	if err := s.store.WriteMetadata(s.cfg.MetaDir); err != nil {
		logger.Warn("write metadata failed", "err", err)
	}
}

// IsDevMode returns true when SKM_DEV env var is set to "1" or "true".
func IsDevMode() bool {
	v := strings.ToLower(os.Getenv("SKM_DEV"))
	return v == "1" || v == "true"
}
