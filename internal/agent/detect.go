package agent

import (
	"os"
	"path/filepath"
)

// Detect returns the list of built-in adapters whose detect paths
// exist under the user's home directory.
func Detect() []Adapter {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	var detected []Adapter
	for _, a := range builtins {
		if isInstalled(a, home) {
			detected = append(detected, a)
		}
	}
	return detected
}

// isInstalled checks whether any of the adapter's detect paths exist
// under the given home directory.
func isInstalled(a Adapter, home string) bool {
	for _, dp := range a.DetectPaths {
		p := filepath.Join(home, dp)
		if _, err := os.Stat(p); err == nil {
			return true
		}
	}
	return false
}
