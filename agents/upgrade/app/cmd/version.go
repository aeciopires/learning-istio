package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the istio-upgrade CLI version",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := fmt.Fprintf(cmd.OutOrStdout(), "istio-upgrade version %s (commit %s)\n", cliVersion, cliGitCommit)
		return err
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
