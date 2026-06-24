package cli

import (
	"github.com/spf13/cobra"
)

func NewRoot() *cobra.Command {
	root := &cobra.Command{
		Use:   "skm",
		Short: "AI Agent Skills Manager",
		Long:  "SKM — Install, organize, and sync skills across AI coding agents",
	}

	root.AddCommand(
		newInstallCmd(),
		newListCmd(),
		newShowCmd(),
		newRemoveCmd(),
		newSyncCmd(),
		newGroupCmd(),
	)

	return root
}
