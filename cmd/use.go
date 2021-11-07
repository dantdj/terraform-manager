package cmd

import (
	"fmt"

	"github.com/dantdj/terraform-manager/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(useCmd)
}

var useCmd = &cobra.Command{
	Use:   "use",
	Args:  cobra.MinimumNArgs(1),
	Short: "Sets the default Terraform version to use",
	Long:  "Sets the default Terraform version to use",
	RunE: func(cmd *cobra.Command, args []string) error {
		config.UpdateCurrentVersion(args[0])
		fmt.Printf("Set to use version %s\n", args[0])

		return nil
	},
}
