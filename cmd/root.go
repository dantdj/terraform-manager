package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "tfm",
		Short: "A manager for Terraform versions",
		Long:  `terraform-manager is a CLI that manages Terraform versions.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
