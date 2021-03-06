package cmd

import (
	"fmt"
	"os"

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
		cwd, _ := os.Getwd()
		pathToBinary := cwd + fmt.Sprintf("/terraform/terraform%s", args[0])
		symlinkPath := cwd + "/terraform/terraform"

		if _, err := os.Lstat(symlinkPath); err == nil {
			if err := os.Remove(symlinkPath); err != nil {
				return fmt.Errorf("failed to unlink: %+v", err)
			}
		}

		err := os.Symlink(pathToBinary, symlinkPath)
		if err != nil {
			return err
		}

		config.UpdateCurrentVersion(args[0])

		fmt.Printf("Set to use version %s\n", args[0])

		return nil
	},
}
