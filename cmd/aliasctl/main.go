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
		return
	}

	command := os.Args[1]

	// Handle shell and file setting before loading aliases
	if command == "set-shell" {
		if len(os.Args) < 3 {
			fmt.Println("Usage: aliasctl set-shell <shell-type>")
			return
		}
		if err := am.SetShell(os.Args[2]); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Shell set to %s\n", am.Shell)
		return
	}

	if command == "set-file" {
		if len(os.Args) < 3 {
			fmt.Println("Usage: aliasctl set-file <alias-file-path>")
			return
		}
		if err := am.SetAliasFile(os.Args[2]); err != nil {
			fmt.Println("Error saving configuration:", err)
			return
		}
		fmt.Printf("Alias file set to %s\n", am.AliasFile)
		return
	}

	// Load existing aliases
	if err := am.LoadAliases(); err != nil {
		fmt.Println("Error loading aliases:", err)
		return
	}

	switch command {
	case "list":
		am.ListAliases()

	case "add":
		if len(os.Args) < 4 {
			fmt.Println("Usage: aliasctl add <name> <command>")
			return
		}
		am.AddAlias(os.Args[2], strings.Join(os.Args[3:], " "))
		if err := am.SaveAliases(); err != nil {
			fmt.Println("Error saving aliases:", err)
		} else {
			fmt.Printf("Added alias: %s = %s\n", os.Args[2], strings.Join(os.Args[3:], " "))
		}

	case "remove":
		if len(os.Args) < 3 {
			fmt.Println("Usage: aliasctl remove <name>")
			return
		}
		if am.RemoveAlias(os.Args[2]) {
			if err := am.SaveAliases(); err != nil {
				fmt.Println("Error saving aliases:", err)
			} else {
				fmt.Printf("Removed alias: %s\n", os.Args[2])
			}
		} else {
			fmt.Printf("Alias '%s' not found\n", os.Args[2])
		}

	case "export":
		if len(os.Args) < 4 {
			fmt.Println("Usage: aliasctl export <shell-type> <output-file>")
			return
		}
		if err := am.ExportAliases(os.Args[2], os.Args[3]); err != nil {
			fmt.Println("Error exporting aliases:", err)
		} else {
			fmt.Printf("Exported aliases to %s in %s format\n", os.Args[3], os.Args[2])
		}

	case "configure-ollama":
		if len(os.Args) < 4 {
			fmt.Println("Usage: aliasctl configure-ollama <endpoint> <model>")
			return
		}
		am.ConfigureOllama(os.Args[2], os.Args[3])
		fmt.Println("Ollama AI provider configured")

	case "configure-openai":
		if len(os.Args) < 5 {
			fmt.Println("Usage: aliasctl configure-openai <endpoint> <api-key> <model>")
			return
		}
		am.ConfigureOpenAI(os.Args[2], os.Args[3], os.Args[4])
		fmt.Println("OpenAI-compatible AI provider configured")

	case "convert":
		if len(os.Args) < 4 {
			fmt.Println("Usage: aliasctl convert <name> <target-shell>")
			return
		}
		if !am.AIConfigured {
			fmt.Println("AI provider not configured. Use configure-ollama or configure-openai first.")
			return
		}
		converted, err := am.ConvertAlias(os.Args[2], os.Args[3])
		if err != nil {
			fmt.Println("Error converting alias:", err)
		} else {
			fmt.Printf("Converted alias for %s: %s\n", os.Args[3], converted)
		}

	case "detect-shell":
		fmt.Printf("Detected shell: %s\n", am.Shell)
		fmt.Printf("Alias file: %s\n", am.AliasFile)
		fmt.Printf("Config directory: %s\n", am.ConfigDir)

	case "import":
		if err := am.ImportAliasesFromShell(); err != nil {
			fmt.Println("Error importing aliases from shell:", err)
		} else {
			fmt.Println("Aliases imported from shell configuration")
		}

	case "apply":
		if err := am.ApplyAliases(); err != nil {
			fmt.Println("Error applying aliases to shell:", err)
		} else {
			fmt.Println("Aliases applied to shell configuration")
		}

	case "configure-ai":
		if len(os.Args) < 2 {
			fmt.Println("Usage: aliasctl configure-ai <provider> [<endpoint> <model> <api-key>]")
			return
		}
		provider := os.Args[2]
		switch provider {
		case "ollama":
			if len(os.Args) < 5 {
				fmt.Println("Usage: aliasctl configure-ai ollama <endpoint> <model>")
				return
			}
			am.ConfigureOllama(os.Args[3], os.Args[4])
			fmt.Println("Ollama AI provider configured")
		case "openai":
			if len(os.Args) < 6 {
				fmt.Println("Usage: aliasctl configure-ai openai <endpoint> <model> <api-key>")
				return
			}
			am.ConfigureOpenAI(os.Args[3], os.Args[5], os.Args[4])
			fmt.Println("OpenAI-compatible AI provider configured")
		default:
			fmt.Println("Unsupported AI provider. Supported providers: ollama, openai")
		}

	default:
		printUsage()
	}
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
	fmt.Println("  convert <name> <target-shell>          - Convert an alias to another shell")
	fmt.Println("  detect-shell                           - Show detected shell and alias file")
	fmt.Println("  set-shell <shell-type>                 - Manually set the shell type")
	fmt.Println("  set-file <alias-file-path>             - Manually set the alias file path")
	fmt.Println("  configure-ai <provider> [<endpoint> <model> <api-key>] - Configure AI provider")
	fmt.Println("Supported shells: bash, zsh, fish, ksh, powershell, pwsh (PowerShell Core), cmd")
}
