package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

// exportCmd represents the export command which outputs aliases to a file in the format for a specific shell.
// This is useful for sharing aliases between different environments or systems.
// Supported shell types include: bash, zsh, fish, ksh, powershell, pwsh, and cmd.
// Example usage: aliasctl export fish ~/.config/fish/aliases.fish
var exportCmd = &cobra.Command{
	Use:   "export [shell-type] [output-file]",
	Short: "Export aliases to a file",
	Long:  `Export aliases to a file for a specific shell type.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		shellType := args[0]
		outputFile := args[1]
		absPath, _ := filepath.Abs(outputFile)

		supportedShells := []string{"bash", "zsh", "fish", "ksh", "powershell", "pwsh", "cmd"}
		validShell := false
		for _, shell := range supportedShells {
			if shellType == shell {
				validShell = true
				break
			}
		}

		if !validShell {
			return fmt.Errorf("unsupported shell type '%s'\n\nSupported shell types: bash, zsh, fish, ksh, powershell, pwsh, cmd", shellType)
		}

		if err := am.ExportAliases(shellType, outputFile); err != nil {
			return fmt.Errorf("failed to export aliases to %s: %w\n\nEnsure the directory exists and you have write permissions", absPath, err)
		}

		fmt.Printf("Successfully exported aliases to %s in %s format\n", outputFile, shellType)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
