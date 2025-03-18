package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// configureOllamaCmd represents the configure-ollama command which sets up Ollama AI provider.
// It requires the endpoint URL and model name as arguments.
// Ollama is a local AI model server that can be used for generating and converting aliases.
// Example usage: aliasctl configure-ollama http://localhost:11434 llama2
var configureOllamaCmd = &cobra.Command{
	Use:   "configure-ollama [endpoint] [model]",
	Short: "Configure Ollama AI provider",
	Long:  `Configure the Ollama AI provider for alias generation and conversion.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		endpoint := args[0]
		model := args[1]

		am.ConfigureOllama(endpoint, model)
		fmt.Println("Ollama AI provider successfully configured")
		return nil
	},
}

// configureOpenAICmd represents the configure-openai command which sets up OpenAI-compatible API.
// It requires the endpoint URL, API key, and model name as arguments.
// This supports both OpenAI's official API and compatible third-party implementations.
// Example usage: aliasctl configure-openai https://api.openai.com YOUR_API_KEY gpt-3.5-turbo
var configureOpenAICmd = &cobra.Command{
	Use:   "configure-openai [endpoint] [api-key] [model]",
	Short: "Configure OpenAI-compatible AI provider",
	Long:  `Configure the OpenAI-compatible AI provider for alias generation and conversion.`,
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		endpoint := args[0]
		apiKey := args[1]
		model := args[2]

		am.ConfigureOpenAI(endpoint, apiKey, model)
		fmt.Println("OpenAI-compatible AI provider successfully configured")

		// If encryption is enabled, remind the user about the key security
		if am.EncryptionUsed {
			fmt.Println("API key will be encrypted using the key stored at:", am.EncryptionKey)
			fmt.Println("WARNING: Keep this key file secure as it's needed to decrypt your API keys.")
		} else {
			fmt.Println("Warning: API key is stored in plaintext. Use 'aliasctl encrypt-api-keys' to encrypt it.")
		}
		return nil
	},
}

// configureAnthropicCmd represents the configure-anthropic command which sets up Anthropic Claude API.
// It requires the endpoint URL, API key, and model name as arguments.
// Anthropic Claude is an AI service that provides high-quality language models.
// Example usage: aliasctl configure-anthropic https://api.anthropic.com YOUR_API_KEY claude-2
var configureAnthropicCmd = &cobra.Command{
	Use:   "configure-anthropic [endpoint] [api-key] [model]",
	Short: "Configure Anthropic Claude AI provider",
	Long:  `Configure the Anthropic Claude AI provider for alias generation and conversion.`,
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		endpoint := args[0]
		apiKey := args[1]
		model := args[2]

		am.ConfigureAnthropic(endpoint, apiKey, model)
		fmt.Println("Anthropic Claude AI provider successfully configured")

		// If encryption is enabled, remind the user about the key security
		if am.EncryptionUsed {
			fmt.Println("API key will be encrypted using the key stored at:", am.EncryptionKey)
			fmt.Println("WARNING: Keep this key file secure as it's needed to decrypt your API keys.")
		} else {
			fmt.Println("Warning: API key is stored in plaintext. Use 'aliasctl encrypt-api-keys' to encrypt it.")
		}
		return nil
	},
}

// configureAICmd represents the configure-ai command which is a unified interface for all AI providers.
// It takes a provider name as the first argument, followed by provider-specific arguments.
// This provides a consistent interface for all supported AI providers.
// Example usage: aliasctl configure-ai ollama http://localhost:11434 llama2
var configureAICmd = &cobra.Command{
	Use:   "configure-ai [provider] [arguments...]",
	Short: "Configure AI provider",
	Long:  `Configure an AI provider for alias generation and conversion.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := args[0]

		switch provider {
		case "ollama":
			if len(args) < 3 {
				return fmt.Errorf("insufficient arguments for ollama configuration\n\nUsage: aliasctl configure-ai ollama <endpoint> <model>\nExample: aliasctl configure-ai ollama http://localhost:11434 llama2")
			}
			am.ConfigureOllama(args[1], args[2])
			fmt.Println("Ollama AI provider successfully configured")

		case "openai":
			if len(args) < 4 {
				return fmt.Errorf("insufficient arguments for OpenAI configuration\n\nUsage: aliasctl configure-ai openai <endpoint> <model> <api-key>\nExample: aliasctl configure-ai openai https://api.openai.com gpt-3.5-turbo YOUR_API_KEY")
			}
			am.ConfigureOpenAI(args[1], args[3], args[2])
			fmt.Println("OpenAI-compatible AI provider successfully configured")

		case "anthropic":
			if len(args) < 4 {
				return fmt.Errorf("insufficient arguments for Anthropic configuration\n\nUsage: aliasctl configure-ai anthropic <endpoint> <model> <api-key>\nExample: aliasctl configure-ai anthropic https://api.anthropic.com claude-2 YOUR_API_KEY")
			}
			am.ConfigureAnthropic(args[1], args[3], args[2])
			fmt.Println("Anthropic Claude AI provider successfully configured")

		default:
			return fmt.Errorf("unsupported AI provider '%s'\n\nSupported providers: ollama, openai, anthropic\n\nExamples:\n  aliasctl configure-ai ollama http://localhost:11434 llama2\n  aliasctl configure-ai openai https://api.openai.com gpt-3.5-turbo YOUR_API_KEY\n  aliasctl configure-ai anthropic https://api.anthropic.com claude-2 YOUR_API_KEY", provider)
		}

		return nil
	},
}

// listProvidersCmd represents the list-providers command which shows all configured AI providers.
// It lists the names of all AI providers that have been set up and are available for use.
// The command will return an error if no providers are configured.
// Example usage: aliasctl list-providers
var listProvidersCmd = &cobra.Command{
	Use:   "list-providers",
	Short: "List all configured AI providers",
	Long:  `List all AI providers that have been configured.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		providers := am.GetAvailableProviders()
		if len(providers) == 0 {
			return fmt.Errorf("no AI providers are configured\n\nTo configure a provider, use one of:\n" +
				"  aliasctl configure-ollama <endpoint> <model>\n" +
				"  aliasctl configure-openai <endpoint> <api-key> <model>\n" +
				"  aliasctl configure-anthropic <endpoint> <api-key> <model>\n\n" +
				"Example for Ollama: aliasctl configure-ollama http://localhost:11434 llama2")
		}

		fmt.Println("Configured AI providers:")
		for _, provider := range providers {
			fmt.Println("- " + provider)
		}
		return nil
	},
}

var generateProvider string

// generateCmd represents the generate command which uses AI to suggest an alias for a shell command.
// It takes a shell command as an argument and uses the configured AI provider to generate a suitable alias.
// The user can specify a particular AI provider using the --provider flag.
// Example usage: aliasctl generate "docker-compose up -d" --provider ollama
var generateCmd = &cobra.Command{
	Use:   "generate [command]",
	Short: "Generate alias suggestion for a command",
	Long:  `Use AI to generate an alias suggestion for a shell command.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		shellCommand := args[0]

		if !am.AIConfigured {
			return fmt.Errorf("AI provider not configured\n\nPlease first configure an AI provider using one of:\n" +
				"  aliasctl configure-ollama <endpoint> <model>\n" +
				"  aliasctl configure-openai <endpoint> <api-key> <model>\n" +
				"  aliasctl configure-anthropic <endpoint> <api-key> <model>\n\n" +
				"Example: aliasctl configure-ollama http://localhost:11434 llama2")
		}

		aliasCommand, err := am.GenerateAlias(shellCommand, generateProvider)
		if err != nil {
			// Check if it's a network-related error
			if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "no such host") {
				return fmt.Errorf("failed to connect to AI provider: %w\n\nMake sure the AI service is running and accessible. If using Ollama, ensure it's started with 'ollama serve'", err)
			}
			return fmt.Errorf("failed to generate alias: %w\n\nCheck that your API key is valid and the AI service is available", err)
		}

		fmt.Printf("Generated alias suggestion: %s\n", aliasCommand)

		// Parse the alias name and command
		aliasName, aliasCmd := parseAliasDefinition(aliasCommand, string(am.Shell))
		if aliasName == "" || aliasCmd == "" {
			return fmt.Errorf("failed to parse the generated alias definition: %s", aliasCommand)
		}

		// Ask if user wants to use suggested name or provide a different one
		fmt.Printf("Use suggested alias name '%s'? [Y/n/custom name]: ", aliasName)
		var nameResponse string
		fmt.Scanln(&nameResponse)

		nameResponse = strings.TrimSpace(nameResponse)
		if nameResponse != "" && strings.ToLower(nameResponse) != "y" && strings.ToLower(nameResponse) != "yes" {
			// If response isn't yes/y and isn't empty, use the response as the custom name
			if strings.ToLower(nameResponse) != "n" && strings.ToLower(nameResponse) != "no" {
				aliasName = nameResponse
			} else {
				// User entered n/no, so prompt for the name explicitly
				fmt.Print("Enter custom alias name: ")
				fmt.Scanln(&aliasName)
				aliasName = strings.TrimSpace(aliasName)

				if aliasName == "" {
					return fmt.Errorf("alias name cannot be empty")
				}
			}
		}

		// Ask if the user wants to save this alias
		fmt.Print("Save this alias? [Y/n]: ")
		var saveResponse string
		fmt.Scanln(&saveResponse)

		if saveResponse == "" || strings.ToLower(saveResponse) == "y" || strings.ToLower(saveResponse) == "yes" {
			am.AddAlias(aliasName, aliasCmd)
			if err := am.SaveAliases(); err != nil {
				return fmt.Errorf("failed to save the new alias: %w", err)
			}
			fmt.Printf("Alias successfully saved: %s = %s\n", aliasName, aliasCmd)
		} else {
			fmt.Println("Alias not saved")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configureOllamaCmd)
	rootCmd.AddCommand(configureOpenAICmd)
	rootCmd.AddCommand(configureAnthropicCmd)
	rootCmd.AddCommand(configureAICmd)
	rootCmd.AddCommand(listProvidersCmd)
	rootCmd.AddCommand(generateCmd)

	// Add provider flag to generate command
	generateCmd.Flags().StringVarP(&generateProvider, "provider", "p", "", "Specify AI provider for generation")
}
