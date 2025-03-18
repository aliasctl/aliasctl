package cmd

import (
	"fmt"

	"github.com/aliasctl/aliasctl/internal/config"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize aliasctl configuration",
	Long:  `Initialize the configuration file for aliasctl with default settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := config.InitConfig()
		if err != nil {
			return err
		}
		fmt.Println("Configuration file initialized successfully")
		return nil
	},
}

// init registers the init command and its flags to the root command.
// This function is called automatically when the package is initialized
// and adds the init command to the application's command tree.
func init() {
	rootCmd.AddCommand(initCmd)
}
