package cmd

import (
	"fmt"

	"github.com/dantdj/terraform-manager/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists current managed Terraform versions",
	Long:  "Lists current managed Terraform versions",
	PreRun: func(cmd *cobra.Command, args []string) {
		config.InitializeConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		for _, terraformInstance := range config.Configuration.TerraformVersions {
			versionString := "  "
			if terraformInstance.Version == config.Configuration.CurrentVersion {
				versionString = "* "
			}
			versionString += terraformInstance.Version
			fmt.Println(versionString)
		}
		return nil
	},
}
