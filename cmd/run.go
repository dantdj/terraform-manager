package cmd

import (
	"fmt"
	"os/exec"

	"github.com/dantdj/terraform-manager/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Args:  cobra.MinimumNArgs(1),
	Short: "Runs the specified command in Terraform",
	Long:  "Runs the specified command in Terraform",
	PreRun: func(cmd *cobra.Command, args []string) {
		config.InitializeConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		versionConfig, err := config.GetCurrentVersion()
		if err != nil {
			return err
		}

		command := exec.Command(versionConfig.PathToFile, args...)
		stdout, err := command.Output()

		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		// Print the output
		fmt.Println(string(stdout))

		return nil
	},
}
