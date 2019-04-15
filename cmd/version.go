package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("nightfury %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// BuildVersion return the build version
func BuildVersion() string {
	if version == "" {
		return "1.0-dev"
	}
	return version
}
