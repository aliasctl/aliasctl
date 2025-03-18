package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

// applyCmd represents the apply command which writes aliases to the shell configuration file.
// This command writes all managed aliases to the configured shell file, preserving any
// other content that might be in the file. It adds a special section marked with
// comments to identify the managed aliases section.
// Example usage: aliasctl apply
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply aliases to shell configuration",
	Long:  `Apply aliases to the current shell's configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Convert to absolute path for better error messages
		absPath, _ := filepath.Abs(am.AliasFile)

		if err := am.ApplyAliases(); err != nil {
			return fmt.Errorf("failed to apply aliases to shell configuration at %s: %w\n\nMake sure you have write permissions to this file or set a different alias file with 'aliasctl set-file'", absPath, err)
		}

		fmt.Printf("Aliases successfully applied to shell configuration at %s\n", am.AliasFile)
		fmt.Println("To use your new aliases, restart your shell or run 'source " + am.AliasFile + "'")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
