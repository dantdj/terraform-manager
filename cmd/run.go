package cmd

import (
	"fmt"
	"os/exec"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Change this to use the actual version defined in config
		version := "1.0.9"
		app := "terraform/terraform" + version

		command := exec.Command(app, args...)
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
