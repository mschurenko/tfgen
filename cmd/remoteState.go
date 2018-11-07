package cmd

import (
  "fmt"
  "github.com/mschurenko/tfgen/templates"
  "github.com/spf13/cobra"
  "os"
)

var noVerifyKey bool

// remoteStateCmd represents the remoteState command
var remoteStateCmd = &cobra.Command{
  Use:   "remote-state <stack path>",
  Short: "<stack path> is everything after the 'stacks/' prefix ",
  Args: func(cmd *cobra.Command, args []string) error {
    if len(args) != 1 {
      return fmt.Errorf("missing <existing stack> argument")
    }
    return nil
  },
  Run: remoteState,
}

func init() {
  rootCmd.AddCommand(remoteStateCmd)
  remoteStateCmd.Flags().BoolVar(&noVerifyKey, "no-verify-key", false, "don't check that key exits in remote backend")
}

func remoteState(cmd *cobra.Command, args []string) {
  if err := templates.RemoteState(s3Config, args[0], noVerifyKey); err != nil {
    fmt.Println("Error:", err)
    os.Exit(1)
  }
}
