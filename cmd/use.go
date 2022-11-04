package cmd

import (
	"fmt"
	"os"
	"regexp"

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
	PreRun: func(cmd *cobra.Command, args []string) {
		config.InitializeConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Handle this error
		directory, _ := os.UserCacheDir()
		directory = directory + "/tfm"

		pathToBinary := fmt.Sprintf("%s/terraform/%s", directory, args[0])
		symlinkPath := directory + "/terraform/terraform"

		if _, err := os.Lstat(symlinkPath); err == nil {
			if err := os.Remove(symlinkPath); err != nil {
				return fmt.Errorf("failed to unlink: %+v", err)
			}
		}

		err := os.Symlink(pathToBinary, symlinkPath)
		if err != nil {
			return err
		}

		// TODO: Add validation that this is a version string
		if isValidVersion(args[0]) {
			config.UpdateCurrentVersion(args[0])
			fmt.Printf("Set to use version %s\n", args[0])
		} else {
			return fmt.Errorf("invalid version string provided: %s", args[0])
		}

		return nil
	},
}

func isValidVersion(version string) bool {
	match, _ := regexp.MatchString("^\\d*\\.\\d*\\.\\d*$", version)
	return match
}
