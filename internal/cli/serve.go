package cli

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/ejoy/skm/internal/server"
)

func newServeCmd() *cobra.Command {
	var (
		port        int
		openBrowser bool
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

			srv := server.New(&server.ServerConfig{
				Store:     cfg.Store,
				SkillsDir: cfg.SkillsDir,
				CacheDir:  cfg.CacheDir,
				MetaDir:   cfg.MetaDir,
				DevMode:   server.IsDevMode(),
			})

			if openBrowser {
				url := fmt.Sprintf("http://localhost:%d", port)
				go openURL(url)
			}

			return srv.Start(port)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 3721, "Port to listen on")
	cmd.Flags().BoolVar(&openBrowser, "open", false, "Open browser automatically")
	return cmd
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
