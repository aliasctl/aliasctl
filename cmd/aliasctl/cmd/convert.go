package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var providerFlag string

// convertCmd represents the convert command which transforms an alias to another shell format.
// It takes an existing alias name and a target shell type as arguments.
// The command uses AI to perform the conversion, ensuring compatibility between different shells.
// Example usage: aliasctl convert dockerup fish --provider ollama
var convertCmd = &cobra.Command{
	Use:   "convert [name] [target-shell]",
	Short: "Convert an alias to another shell",
	Long:  `Convert an alias from the current shell format to another shell format.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		targetShell := args[1]

		supportedShells := []string{"bash", "zsh", "fish", "ksh", "powershell", "pwsh", "cmd"}
		validShell := false
		for _, shell := range supportedShells {
			if targetShell == shell {
				validShell = true
				break
			}
		}

		if !validShell {
			return fmt.Errorf("unsupported target shell '%s'\n\nSupported shell types: bash, zsh, fish, ksh, powershell, pwsh, cmd", targetShell)
		}

		if !am.AIConfigured {
			return fmt.Errorf("AI provider not configured\n\nPlease first configure an AI provider using one of:\n" +
				"  aliasctl configure-ollama <endpoint> <model>\n" +
				"  aliasctl configure-openai <endpoint> <api-key> <model>\n" +
				"  aliasctl configure-anthropic <endpoint> <api-key> <model>\n\n" +
				"Example: aliasctl configure-ollama http://localhost:11434 llama2")
		}

		converted, err := am.ConvertAlias(name, targetShell, providerFlag)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return fmt.Errorf("alias '%s' not found\n\nRun 'aliasctl list' to see available aliases", name)
			}

			// Check if it's a network-related error
			if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "no such host") {
				return fmt.Errorf("failed to connect to AI provider: %w\n\nMake sure the AI service is running and accessible. If using Ollama, ensure it's started with 'ollama serve'", err)
			}

			return fmt.Errorf("failed to convert alias '%s' to %s format: %w\n\nCheck that your API key is valid and the AI service is available", name, targetShell, err)
		}

		fmt.Printf("Successfully converted alias for %s: %s\n", targetShell, converted)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Add provider flag
	convertCmd.Flags().StringVarP(&providerFlag, "provider", "p", "", "Specify AI provider for conversion")
}
