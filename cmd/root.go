package cmd

import (
	goflag "flag"
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "hc-agent",
	Short: "hc-agent is a simple healthchecks.io Agent in Go",
	Long: `
hc-agent is a simple  healthchecks.io Agent in Go.
built with love by engineering-toolkit team.
Complete documentation is available at https://github.com/engineering-toolkit/healthchecks-agent`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		// Do Stuff Here
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	goflag.CommandLine.Parse([]string{})

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		// flush before exit non-zero
		glog.Flush()
		os.Exit(-1)
	}
	// flush before exit
	glog.Flush()
}
