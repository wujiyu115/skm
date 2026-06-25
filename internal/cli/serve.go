package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/ejoy/skm/internal/logger"
	"github.com/ejoy/skm/internal/server"
)

func newServeCmd() *cobra.Command {
	var (
		port        int
		openBrowser bool
		debug       bool
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the web UI",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			devMode := debug || server.IsDevMode()

			srv := server.New(&server.ServerConfig{
				Store:     cfg.Store,
				SkillsDir: cfg.SkillsDir,
				CacheDir:  cfg.CacheDir,
				MetaDir:   cfg.MetaDir,
				DevMode:   devMode,
			})

			if debug {
				go startViteDev()
			}

			browseURL := fmt.Sprintf("http://localhost:%d", port)
			if debug {
				browseURL = "http://localhost:5173"
			}

			if openBrowser {
				go openURL(browseURL)
			}

			return srv.Start(port)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 3721, "Port to listen on")
	cmd.Flags().BoolVar(&openBrowser, "open", false, "Open browser automatically")
	cmd.Flags().BoolVar(&debug, "debug", false, "Debug mode: start Vite dev server for hot reload")
	return cmd
}

func startViteDev() {
	exe, _ := os.Executable()
	webDir := filepath.Join(filepath.Dir(exe), "web")
	if _, err := os.Stat(filepath.Join(webDir, "package.json")); err != nil {
		cwd, _ := os.Getwd()
		webDir = filepath.Join(cwd, "web")
	}
	if _, err := os.Stat(filepath.Join(webDir, "package.json")); err != nil {
		logger.Warn("web/ directory not found, skipping Vite dev server")
		fmt.Fprintln(os.Stderr, "warning: web/ directory not found, Vite dev server not started")
		return
	}

	logger.Info("starting Vite dev server", "dir", webDir)
	fmt.Println("Starting Vite dev server at http://localhost:5173")

	cmd := exec.Command("npx", "vite", "--host")
	cmd.Dir = webDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Error("vite dev server exited", "err", err)
	}
}

func openURL(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	}
	if cmd != nil {
		cmd.Start()
	}
}
