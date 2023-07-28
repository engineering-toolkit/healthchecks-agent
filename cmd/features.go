package cmd

import (
	"fmt"

	"github.com/engineering-toolkit/healthchecks-agent/config"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the version of hc-agent.",
	Long:  "Shows the version of hc-agent.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hc-agent version ", config.Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
