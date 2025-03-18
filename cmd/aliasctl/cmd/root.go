package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aliasctl/aliasctl/pkg/aliasctl"
	"github.com/spf13/cobra"
)

var am *aliasctl.AliasManager
var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "aliasctl",
	Short:   "Cross-platform shell alias manager",
	Long:    `AliasCtl is a powerful tool that helps you manage shell aliases across different operating systems and shell environments.`,
	Version: aliasctl.GetVersion(),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip loading for certain setup commands
		cmdName := cmd.Name()
		if cmdName == "set-shell" || cmdName == "set-file" || cmdName == "version" {
			return nil
		}

		// Load aliases for all other commands
		if err := am.LoadAliases(); err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "Error loading aliases: %v\n", err)
				// Check if the error is due to a missing file
				if os.IsNotExist(err) {
					dir := filepath.Dir(am.AliasStore)
					return fmt.Errorf("alias storage file not found at %s\n\nSuggestions:\n- Make sure the directory %s exists\n- Check file permissions\n- Use 'aliasctl set-file' to specify a different location", am.AliasStore, dir)
				}

				// Check if it might be a permissions issue
				if os.IsPermission(err) {
					return fmt.Errorf("permission denied when accessing %s\n\nSuggestions:\n- Check file permissions\n- Run with elevated privileges\n- Use 'aliasctl set-file' to specify a different location", am.AliasStore)
				}

				return fmt.Errorf("error loading aliases: %w\n\nUse --verbose flag for more details", err)
			}
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initAliasManager)

	// Add global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}

// initAliasManager initializes the alias manager
func initAliasManager() {
	am = aliasctl.NewAliasManager()

	// Check if we can access the config directory
	if _, err := os.Stat(am.ConfigDir); os.IsNotExist(err) {
		// Try to create it
		if err := os.MkdirAll(am.ConfigDir, 0755); err != nil && verbose {
			fmt.Fprintf(os.Stderr, "Warning: Couldn't create config directory %s: %v\n", am.ConfigDir, err)
		}
	}
}
