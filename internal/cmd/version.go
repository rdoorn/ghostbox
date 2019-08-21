package cmd

import (
	"fmt"

	"github.com/rdoorn/ixxi/internal/handler"
	"github.com/spf13/cobra"
)

func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "shows the version of Mercury",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s version: %s (build: %s)\nSha: %s\n", handler.Name, handler.Version, handler.VersionBuild, handler.VersionSha)
		},
	}
	return cmd
}
