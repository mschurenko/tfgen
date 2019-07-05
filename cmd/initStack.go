package cmd

import (
	"fmt"
	"os"

	"github.com/mschurenko/tfgen/templates"
	"github.com/spf13/cobra"
)

// initStackCmd represents the init command
var initStackCmd = &cobra.Command{
	Use:   "init-stack",
	Short: "use this to initialize a new terraform directory",
	Run:   initStack,
}

func init() {
	rootCmd.AddCommand(initStackCmd)
	initStackCmd.Flags().BoolVar(&forceInitOverride, "force", false, "force overriding of existing tf files")
}

var forceInitOverride bool

func initStack(cmd *cobra.Command, args []string) {
	if err := templates.InitStack(s3Config, environments, stackRx, reqVer, forceInitOverride); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
