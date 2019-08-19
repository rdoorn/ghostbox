package cmd

import (
	"fmt"

	"github.com/rdoorn/ghostbox/internal/ghostbox"
	"github.com/spf13/cobra"
)

func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "shows the version of Mercury",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s version: %s (build: %s)\nSha: %s\n", ghostbox.Name, ghostbox.Version, ghostbox.VersionBuild, ghostbox.VersionSha)
		},
	}
	return cmd
}
