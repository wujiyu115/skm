package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "dev"

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print skm version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("skm %s\n", Version)
		},
	}
}
