package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// setShellCmd represents the set-shell command which configures the default shell type.
// Valid shell types include: bash, zsh, fish, ksh, powershell, pwsh, and cmd.
// The shell type affects how aliases are formatted and applied.
// Example usage: aliasctl set-shell zsh
var setShellCmd = &cobra.Command{
	Use:   "set-shell [shell-type]",
	Short: "Manually set the shell type",
	Long:  `Manually set the shell type to use for alias operations.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		shellType := args[0]

		if err := am.SetShell(shellType); err != nil {
			supportedShells := "bash, zsh, fish, ksh, powershell, pwsh, cmd"
			return fmt.Errorf("failed to set shell type to '%s': %w\n\nPlease use one of the supported shell types: %s", shellType, err, supportedShells)
		}

		fmt.Printf("Shell successfully set to %s\n", am.Shell)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setShellCmd)
}
