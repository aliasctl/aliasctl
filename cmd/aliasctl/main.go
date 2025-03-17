package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aliasctl/aliasctl/pkg/aliasctl"
)

func main() {
	am := aliasctl.NewAliasManager()

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// Version command
	if command == "version" {
		fmt.Printf("aliasctl version %s\n", aliasctl.GetVersion())
		return
	}

	// Handle shell and file setting before loading aliases
	if command == "set-shell" {
		if len(os.Args) < 3 {
			fmt.Println("Usage: aliasctl set-shell <shell-type>")
			os.Exit(1)
		}
		if err := am.SetShell(os.Args[2]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Shell set to %s\n", am.Shell)
		return
	}

	if command == "set-file" {
		if len(os.Args) < 3 {
			fmt.Println("Usage: aliasctl set-file <alias-file-path>")
			os.Exit(1)
		}
		if err := am.SetAliasFile(os.Args[2]); err != nil {
			fmt.Println("Error saving configuration:", err)
			os.Exit(1)
		}
		fmt.Printf("Alias file set to %s\n", am.AliasFile)
		return
	}

	// Handle encryption commands
	if command == "encrypt-api-keys" {
		if err := am.EncryptAPIKeys(); err != nil {
			fmt.Printf("Error encrypting API keys: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if command == "disable-encryption" {
		if err := am.DisableEncryption(); err != nil {
			fmt.Printf("Error disabling encryption: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Load existing aliases
	if err := am.LoadAliases(); err != nil {
		fmt.Println("Error loading aliases:", err)
		os.Exit(1)
	}

	// Execute the requested command
	exitCode := executeCommand(am, command, os.Args[2:])
	os.Exit(exitCode)
}

// executeCommand handles the execution of the specified command
func executeCommand(am *aliasctl.AliasManager, command string, args []string) int {
	switch command {
	case "list":
		am.ListAliases()

	case "add":
		if len(args) < 2 {
			fmt.Println("Usage: aliasctl add <name> <command>")
			return 1
		}
		am.AddAlias(args[0], strings.Join(args[1:], " "))
		if err := am.SaveAliases(); err != nil {
			fmt.Println("Error saving aliases:", err)
			return 1
		} else {
			fmt.Printf("Added alias: %s = %s\n", args[0], strings.Join(args[1:], " "))
		}

	case "remove":
		if len(args) < 1 {
			fmt.Println("Usage: aliasctl remove <name>")
			return 1
		}
		if am.RemoveAlias(args[0]) {
			if err := am.SaveAliases(); err != nil {
				fmt.Println("Error saving aliases:", err)
				return 1
			} else {
				fmt.Printf("Removed alias: %s\n", args[0])
			}
		} else {
			fmt.Printf("Alias '%s' not found\n", args[0])
			return 1
		}

	case "export":
		if len(args) < 2 {
			fmt.Println("Usage: aliasctl export <shell-type> <output-file>")
			return 1
		}
		if err := am.ExportAliases(args[0], args[1]); err != nil {
			fmt.Println("Error exporting aliases:", err)
			return 1
		} else {
			fmt.Printf("Exported aliases to %s in %s format\n", args[1], args[0])
		}

	case "configure-ollama":
		if len(args) < 2 {
			fmt.Println("Usage: aliasctl configure-ollama <endpoint> <model>")
			return 1
		}
		am.ConfigureOllama(args[0], args[1])
		fmt.Println("Ollama AI provider configured")

	case "configure-openai":
		if len(args) < 3 {
			fmt.Println("Usage: aliasctl configure-openai <endpoint> <api-key> <model>")
			return 1
		}
		am.ConfigureOpenAI(args[0], args[1], args[2])
		fmt.Println("OpenAI-compatible AI provider configured")

		// If encryption is enabled, remind the user about the key security
		if am.EncryptionUsed {
			fmt.Println("API key will be encrypted using the key stored at:", am.EncryptionKey)
			fmt.Println("WARNING: Keep this key file secure as it's needed to decrypt your API keys.")
		} else {
			fmt.Println("Warning: API key is stored in plaintext. Use 'aliasctl encrypt-api-keys' to encrypt it.")
		}

	case "configure-anthropic":
		if len(args) < 3 {
			fmt.Println("Usage: aliasctl configure-anthropic <endpoint> <api-key> <model>")
			return 1
		}
		am.ConfigureAnthropic(args[0], args[1], args[2])
		fmt.Println("Anthropic Claude AI provider configured")

		// If encryption is enabled, remind the user about the key security
		if am.EncryptionUsed {
			fmt.Println("API key will be encrypted using the key stored at:", am.EncryptionKey)
			fmt.Println("WARNING: Keep this key file secure as it's needed to decrypt your API keys.")
		} else {
			fmt.Println("Warning: API key is stored in plaintext. Use 'aliasctl encrypt-api-keys' to encrypt it.")
		}

	case "convert":
		providerName := ""
		var targetShell string
		var name string

		// Parse arguments with optional provider
		switch len(args) {
		case 2:
			name = args[0]
			targetShell = args[1]
		case 3:
			name = args[0]
			targetShell = args[1]
			providerName = args[2]
		default:
			fmt.Println("Usage: aliasctl convert <name> <target-shell> [provider]")
			return 1
		}

		if !am.AIConfigured {
			fmt.Println("AI provider not configured. Use configure-ollama, configure-openai, or configure-anthropic first.")
			return 1
		}

		converted, err := am.ConvertAlias(name, targetShell, providerName)
		if err != nil {
			fmt.Println("Error converting alias:", err)
			return 1
		} else {
			fmt.Printf("Converted alias for %s: %s\n", targetShell, converted)
		}

	case "detect-shell":
		fmt.Printf("Detected shell: %s\n", am.Shell)
		fmt.Printf("Alias file: %s\n", am.AliasFile)
		fmt.Printf("Config directory: %s\n", am.ConfigDir)

	case "import":
		if err := am.ImportAliasesFromShell(); err != nil {
			fmt.Println("Error importing aliases from shell:", err)
			return 1
		} else {
			fmt.Println("Aliases imported from shell configuration")
		}

	case "apply":
		if err := am.ApplyAliases(); err != nil {
			fmt.Println("Error applying aliases to shell:", err)
			return 1
		} else {
			fmt.Println("Aliases applied to shell configuration")
		}

	case "configure-ai":
		if len(args) < 1 {
			fmt.Println("Usage: aliasctl configure-ai <provider> [<endpoint> <model> <api-key>]")
			return 1
		}
		provider := args[0]
		switch provider {
		case "ollama":
			if len(args) < 3 {
				fmt.Println("Usage: aliasctl configure-ai ollama <endpoint> <model>")
				return 1
			}
			am.ConfigureOllama(args[1], args[2])
			fmt.Println("Ollama AI provider configured")
		case "openai":
			if len(args) < 4 {
				fmt.Println("Usage: aliasctl configure-ai openai <endpoint> <model> <api-key>")
				return 1
			}
			am.ConfigureOpenAI(args[1], args[3], args[2])
			fmt.Println("OpenAI-compatible AI provider configured")
		case "anthropic":
			if len(args) < 4 {
				fmt.Println("Usage: aliasctl configure-ai anthropic <endpoint> <model> <api-key>")
				return 1
			}
			am.ConfigureAnthropic(args[1], args[3], args[2])
			fmt.Println("Anthropic Claude AI provider configured")
		default:
			fmt.Println("Unsupported AI provider. Supported providers: ollama, openai, anthropic")
			return 1
		}

	case "list-providers":
		providers := am.GetAvailableProviders()
		if len(providers) == 0 {
			fmt.Println("No AI providers configured")
			return 1
		}
		fmt.Println("Configured AI providers:")
		for _, provider := range providers {
			fmt.Println("- " + provider)
			return 1
		}

	case "generate":
		if len(args) < 1 {
			fmt.Println("Usage: aliasctl generate <command> [provider]")
			return 1
		}

		var shellCommand, providerName string

		if len(args) == 1 {
			shellCommand = args[0]
			providerName = "" // Use default provider
		} else {
			shellCommand = args[0]
			providerName = args[1]
		}

		if !am.AIConfigured {
			fmt.Println("AI provider not configured. Use configure-ollama, configure-openai, or configure-anthropic first.")
			return 1
		}

		aliasCommand, err := am.GenerateAlias(shellCommand, providerName)
		if err != nil {
			fmt.Println("Error generating alias:", err)
			return 1
		}

		fmt.Printf("Generated alias suggestion: %s\n", aliasCommand)

		// Ask if the user wants to save this alias
		fmt.Print("Do you want to save this alias? (y/n): ")
		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
			// Parse the alias name and command
			aliasName, aliasCmd := parseAliasDefinition(aliasCommand, string(am.Shell))
			if aliasName == "" || aliasCmd == "" {
				fmt.Println("Could not parse the alias definition.")
				return 1
			}

			am.AddAlias(aliasName, aliasCmd)
			if err := am.SaveAliases(); err != nil {
				fmt.Println("Error saving alias:", err)
				return 1
			}
			fmt.Printf("Alias saved: %s = %s\n", aliasName, aliasCmd)
		} else {
			fmt.Println("Alias not saved.")
		}

	default:
		printUsage()
		return 1
	}

	return 0
}

// parseAliasDefinition attempts to extract the alias name and command from a definition
func parseAliasDefinition(definition, shellType string) (name string, command string) {
	definition = strings.TrimSpace(definition)

	switch shellType {
	case "bash", "zsh", "fish", "ksh":
		if strings.HasPrefix(definition, "alias ") {
			parts := strings.SplitN(strings.TrimPrefix(definition, "alias "), "=", 2)
			if len(parts) == 2 {
				name = strings.TrimSpace(parts[0])
				// Remove surrounding quotes
				command = strings.Trim(strings.TrimSpace(parts[1]), "'\"")
				return
			}
		}
	case "powershell", "pwsh":
		if strings.HasPrefix(definition, "Set-Alias ") {
			parts := strings.Fields(strings.TrimPrefix(definition, "Set-Alias "))
			if len(parts) >= 2 {
				name = parts[0]
				command = parts[1]
				return
			}
		} else if strings.HasPrefix(definition, "function ") {
			parts := strings.SplitN(strings.TrimPrefix(definition, "function "), " {", 2)
			if len(parts) == 2 {
				name = strings.TrimSpace(parts[0])
				command = strings.TrimSpace(strings.TrimSuffix(parts[1], "}"))
				return
			}
		}
	case "cmd":
		if strings.HasPrefix(definition, "doskey ") {
			parts := strings.SplitN(strings.TrimPrefix(definition, "doskey "), "=", 2)
			if len(parts) == 2 {
				name = strings.TrimSpace(parts[0])
				command = strings.TrimSpace(parts[1])
				return
			}
		}
	}

	// Fallback: if we couldn't parse with the specific shell format,
	// try to use a generic approach - look for the first space or equals
	if strings.Contains(definition, "=") {
		parts := strings.SplitN(definition, "=", 2)
		if len(parts) == 2 {
			name = strings.TrimSpace(parts[0])
			command = strings.TrimSpace(parts[1])
			return
		}
	}

	return "", ""
}

// printUsage prints the usage information for the aliasctl command.
func printUsage() {
	fmt.Println("AliasCtl - Cross-platform shell alias manager")
	fmt.Println("Usage: aliasctl <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  list                                   - List all aliases")
	fmt.Println("  add <name> <command>                   - Add a new alias")
	fmt.Println("  remove <name>                          - Remove an alias")
	fmt.Println("  import                                 - Import aliases from shell configuration")
	fmt.Println("  apply                                  - Apply aliases to shell configuration")
	fmt.Println("  export <shell-type> <output-file>      - Export aliases to a file for a specific shell")
	fmt.Println("  configure-ollama <endpoint> <model>    - Configure Ollama AI provider")
	fmt.Println("  configure-openai <endpoint> <api-key> <model> - Configure OpenAI-compatible AI provider")
	fmt.Println("  configure-anthropic <endpoint> <api-key> <model> - Configure Anthropic Claude AI provider")
	fmt.Println("  convert <name> <target-shell>          - Convert an alias to another shell")
	fmt.Println("  convert <name> <target-shell> [provider] - Convert an alias to another shell using the specified provider")
	fmt.Println("  detect-shell                           - Show detected shell and alias file")
	fmt.Println("  set-shell <shell-type>                 - Manually set the shell type")
	fmt.Println("  set-file <alias-file-path>             - Manually set the alias file path")
	fmt.Println("  configure-ai <provider> [<endpoint> <model> <api-key>] - Configure AI provider")
	fmt.Println("  version                                - Display version information")
	fmt.Println("  encrypt-api-keys                       - Encrypt API keys in configuration")
	fmt.Println("  disable-encryption                     - Disable API key encryption")
	fmt.Println("  list-providers                         - List all configured AI providers")
	fmt.Println("  generate <command> [provider]          - Generate alias suggestion for a command")
	fmt.Println("Supported shells: bash, zsh, fish, ksh, powershell, pwsh (PowerShell Core), cmd")
	fmt.Println("Supported AI providers: ollama, openai, anthropic")
}
