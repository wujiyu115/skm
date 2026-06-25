package cli

import (
	"encoding/json"
	"fmt"

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
		newEnableCmd(),
		newDisableCmd(),
		newSyncCmd(),
		newGroupCmd(),
		newTagCmd(),
		newUpdateCmd(),
		newAgentCmd(),
		newConfigCmd(),
		newServeCmd(),
		newVersionCmd(),
		newSearchCmd(),
		newExportCmd(),
		newInfoCmd(),
	)

	return root
}

func printJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
