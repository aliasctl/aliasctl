package aliasctl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

// Completion types and templates for different shells
var (
	bashCompletionTemplate = `
# aliasctl bash completion script
_aliasctl_completions() {
	local cur prev opts
	COMPREPLY=()
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[COMP_CWORD-1]}"
	
	# List of all commands
	opts="list add remove export convert detect-shell import apply configure-ollama configure-openai configure-anthropic configure-ai list-providers generate set-shell set-file encrypt-api-keys disable-encryption version"
	
	case "${prev}" in
		add|remove|convert)
			# List aliases for these commands
			local aliases=$(aliasctl list | awk '{print $1}')
			COMPREPLY=( $(compgen -W "${aliases}" -- ${cur}) )
			return 0
			;;
		export|set-shell)
			# List shell types
			local shells="bash zsh fish ksh powershell pwsh cmd"
			COMPREPLY=( $(compgen -W "${shells}" -- ${cur}) )
			return 0
			;;
		configure-ai)
			# List provider types
			local providers="ollama openai anthropic"
			COMPREPLY=( $(compgen -W "${providers}" -- ${cur}) )
			return 0
			;;
		*)
			# Default to commands
			COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
			return 0
			;;
	esac
}

complete -F _aliasctl_completions aliasctl
`

	zshCompletionTemplate = `
# aliasctl zsh completion script
_aliasctl() {
	local -a commands
	commands=(
		'list:List all aliases'
		'add:Add a new alias'
		'remove:Remove an alias'
		'export:Export aliases to a file'
		'convert:Convert an alias to another shell'
		'detect-shell:Show detected shell and alias file'
		'import:Import aliases from shell configuration'
		'apply:Apply aliases to shell configuration'
		'configure-ollama:Configure Ollama AI provider'
		'configure-openai:Configure OpenAI-compatible AI provider'
		'configure-anthropic:Configure Anthropic Claude AI provider'
		'configure-ai:Configure AI provider'
		'version:Display version information'
		'encrypt-api-keys:Encrypt API keys in configuration'
		'disable-encryption:Disable API key encryption'
		'list-providers:List all configured AI providers'
		'generate:Generate alias suggestion for a command'
		'set-shell:Manually set the shell type'
		'set-file:Manually set the alias file path'
	)
	
	_describe -t commands 'aliasctl commands' commands
	
	case "$words[2]" in
		add|remove|convert)
			# Get list of aliases
			local -a aliases
			aliases=($(aliasctl list | awk '{print $1}'))
			_describe -t aliases 'aliases' aliases
			;;
		export|set-shell)
			# List shell types
			local -a shells
			shells=('bash' 'zsh' 'fish' 'ksh' 'powershell' 'pwsh' 'cmd')
			_describe -t shells 'shells' shells
			;;
		configure-ai)
			# List provider types
			local -a providers
			providers=('ollama' 'openai' 'anthropic')
			_describe -t providers 'providers' providers
			;;
	esac
	
	return 0
}

compdef _aliasctl aliasctl
`

	fishCompletionTemplate = `
# aliasctl fish completion script
complete -c aliasctl -f

# Command completions
complete -c aliasctl -n "__fish_use_subcommand" -a list -d "List all aliases"
complete -c aliasctl -n "__fish_use_subcommand" -a add -d "Add a new alias"
complete -c aliasctl -n "__fish_use_subcommand" -a remove -d "Remove an alias"
complete -c aliasctl -n "__fish_use_subcommand" -a export -d "Export aliases to a file"
complete -c aliasctl -n "__fish_use_subcommand" -a convert -d "Convert an alias to another shell"
complete -c aliasctl -n "__fish_use_subcommand" -a detect-shell -d "Show detected shell and alias file"
complete -c aliasctl -n "__fish_use_subcommand" -a import -d "Import aliases from shell configuration"
complete -c aliasctl -n "__fish_use_subcommand" -a apply -d "Apply aliases to shell configuration"
complete -c aliasctl -n "__fish_use_subcommand" -a configure-ollama -d "Configure Ollama AI provider"
complete -c aliasctl -n "__fish_use_subcommand" -a configure-openai -d "Configure OpenAI-compatible AI provider"
complete -c aliasctl -n "__fish_use_subcommand" -a configure-anthropic -d "Configure Anthropic Claude AI provider"
complete -c aliasctl -n "__fish_use_subcommand" -a configure-ai -d "Configure AI provider"
complete -c aliasctl -n "__fish_use_subcommand" -a version -d "Display version information"
complete -c aliasctl -n "__fish_use_subcommand" -a encrypt-api-keys -d "Encrypt API keys in configuration"
complete -c aliasctl -n "__fish_use_subcommand" -a disable-encryption -d "Disable API key encryption"
complete -c aliasctl -n "__fish_use_subcommand" -a list-providers -d "List all configured AI providers"
complete -c aliasctl -n "__fish_use_subcommand" -a generate -d "Generate alias suggestion for a command"
complete -c aliasctl -n "__fish_use_subcommand" -a set-shell -d "Manually set the shell type"
complete -c aliasctl -n "__fish_use_subcommand" -a set-file -d "Manually set the alias file path"

# Alias name completions
complete -c aliasctl -n "__fish_seen_subcommand_from remove convert" -a "(aliasctl list | string replace -r ' .*\$' '')"

# Shell type completions
complete -c aliasctl -n "__fish_seen_subcommand_from export set-shell" -a "bash zsh fish ksh powershell pwsh cmd"

# Provider completions
complete -c aliasctl -n "__fish_seen_subcommand_from configure-ai" -a "ollama openai anthropic"
`

	powershellCompletionTemplate = `
# aliasctl PowerShell completion script

function _aliasctl_completion {
    param($wordToComplete, $commandAst, $cursorPosition)
    
    # Get the current command being typed
    $command = $commandAst.ToString()
    
    # Extract the subcommand (if any)
    $subCommand = $null
    if ($command -match 'aliasctl\s+(\w+)') {
        $subCommand = $matches[1]
    }
    
    # No subcommand yet, suggest available commands
    if (-not $subCommand -or $subCommand -eq $wordToComplete) {
        @(
            "list",
            "add",
            "remove",
            "export",
            "convert",
            "detect-shell",
            "import",
            "apply",
            "configure-ollama",
            "configure-openai",
            "configure-anthropic",
            "configure-ai",
            "version",
            "encrypt-api-keys",
            "disable-encryption",
            "list-providers",
            "generate",
            "set-shell",
            "set-file",
            "completion",
            "install-completion"
        ) | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
            [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
        }
        return
    }
    
    # Provide specific completions based on the subcommand
    switch ($subCommand) {
        "remove" {
            # Get aliases from aliasctl list
            $aliases = & aliasctl list | ForEach-Object { ($_ -split '=')[0].Trim() } | Where-Object { $_ }
            $aliases | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
            }
        }
        "convert" {
            if ($command -match 'aliasctl\s+convert\s+(\S+)') {
                # If we already have an alias name, suggest shells
                $shells = @("bash", "zsh", "fish", "ksh", "powershell", "pwsh", "cmd")
                $shells | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                    [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
                }
            } else {
                # Suggest alias names
                $aliases = & aliasctl list | ForEach-Object { ($_ -split '=')[0].Trim() } | Where-Object { $_ }
                $aliases | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                    [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
                }
            }
        }
        "export" {
            # Suggest shell types
            $shells = @("bash", "zsh", "fish", "ksh", "powershell", "pwsh", "cmd")
            $shells | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
            }
        }
        "set-shell" {
            # Suggest shell types
            $shells = @("bash", "zsh", "fish", "ksh", "powershell", "pwsh", "cmd")
            $shells | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
            }
        }
        "configure-ai" {
            # Suggest provider types
            $providers = @("ollama", "openai", "anthropic")
            $providers | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
            }
        }
        "completion" {
            # Suggest shell types for completion generation
            $shells = @("bash", "zsh", "fish", "powershell", "pwsh")
            $shells | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
            }
        }
        default {
            # No specific completions for this subcommand
            @()
        }
    }
}

Register-ArgumentCompleter -Native -CommandName aliasctl -ScriptBlock $function:_aliasctl_completion
`

	pwshCompletionTemplate = `
# aliasctl PowerShell Core completion script

function _aliasctl_completion {
    param($wordToComplete, $commandAst, $cursorPosition)
    
    # Get the current command being typed
    $command = $commandAst.ToString()
    
    # Extract the subcommand (if any)
    $subCommand = $null
    if ($command -match 'aliasctl\s+(\w+)') {
        $subCommand = $matches[1]
    }
    
    # No subcommand yet, suggest available commands
    if (-not $subCommand -or $subCommand -eq $wordToComplete) {
        @(
            "list",
            "add",
            "remove",
            "export",
            "convert",
            "detect-shell",
            "import",
            "apply",
            "configure-ollama",
            "configure-openai",
            "configure-anthropic",
            "configure-ai",
            "version",
            "encrypt-api-keys",
            "disable-encryption",
            "list-providers",
            "generate",
            "set-shell",
            "set-file",
            "completion",
            "install-completion"
        ) | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
            [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
        }
        return
    }
    
    # Provide specific completions based on the subcommand
    switch ($subCommand) {
        "remove" {
            # Get aliases from aliasctl list
            $aliases = & aliasctl list | ForEach-Object { ($_ -split '=')[0].Trim() } | Where-Object { $_ }
            $aliases | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
            }
        }
        "convert" {
            if ($command -match 'aliasctl\s+convert\s+(\S+)') {
                # If we already have an alias name, suggest shells
                $shells = @("bash", "zsh", "fish", "ksh", "powershell", "pwsh", "cmd")
                $shells | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                    [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
                }
            } else {
                # Suggest alias names
                $aliases = & aliasctl list | ForEach-Object { ($_ -split '=')[0].Trim() } | Where-Object { $_ }
                $aliases | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                    [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
                }
            }
        }
        "export" {
            # Suggest shell types
            $shells = @("bash", "zsh", "fish", "ksh", "powershell", "pwsh", "cmd")
            $shells | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
            }
        }
        "set-shell" {
            # Suggest shell types
            $shells = @("bash", "zsh", "fish", "ksh", "powershell", "pwsh", "cmd")
            $shells | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
            }
        }
        "configure-ai" {
            # Suggest provider types
            $providers = @("ollama", "openai", "anthropic")
            $providers | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
            }
        }
        "completion" {
            # Suggest shell types for completion generation
            $shells = @("bash", "zsh", "fish", "powershell", "pwsh")
            $shells | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
            }
        }
        default {
            # No specific completions for this subcommand
            @()
        }
    }
}

Register-ArgumentCompleter -Native -CommandName aliasctl -ScriptBlock $function:_aliasctl_completion
`
)

// GenerateCompletionScript generates a shell completion script for the given shell
func (am *AliasManager) GenerateCompletionScript(shellType string) (string, error) {
	var tmplContent string

	switch shellType {
	case "bash":
		tmplContent = bashCompletionTemplate
	case "zsh":
		tmplContent = zshCompletionTemplate
	case "fish":
		tmplContent = fishCompletionTemplate
	case "powershell":
		tmplContent = powershellCompletionTemplate
	case "pwsh":
		tmplContent = pwshCompletionTemplate // Uses same script as PowerShell
	default:
		return "", fmt.Errorf("completion script not available for shell: %s", shellType)
	}

	tmpl, err := template.New("completion").Parse(tmplContent)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, nil); err != nil {
		return "", err
	}

	return result.String(), nil
}

// InstallCompletionScript installs the completion script for the current shell
func (am *AliasManager) InstallCompletionScript() error {
	shellType := string(am.Shell)

	script, err := am.GenerateCompletionScript(shellType)
	if err != nil {
		return err
	}

	// Determine where to install the completion script
	var completionPath string

	switch shellType {
	case "bash":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		completionPath = filepath.Join(homeDir, ".bash_completion.d", "aliasctl.bash")
	case "zsh":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		completionPath = filepath.Join(homeDir, ".zsh", "completion", "_aliasctl")
	case "fish":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		completionPath = filepath.Join(homeDir, ".config", "fish", "completions", "aliasctl.fish")
	case "powershell", "pwsh":
		// PowerShell uses a profile directory
		profileDir, err := getPowerShellProfileDir(shellType == "pwsh")
		if err != nil {
			return err
		}
		completionPath = filepath.Join(profileDir, "aliases", "aliasctl.ps1")
	default:
		return fmt.Errorf("completion installation not supported for shell: %s", shellType)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(completionPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write completion script
	if err := os.WriteFile(completionPath, []byte(script), 0644); err != nil {
		return err
	}

	// Provide instructions for sourcing the completion script
	fmt.Printf("Completion script installed to %s\n", completionPath)

	switch shellType {
	case "bash":
		fmt.Printf("Add the following line to your ~/.bashrc file:\n")
		fmt.Printf("  source %s\n", completionPath)
	case "zsh":
		fmt.Printf("Add the following line to your ~/.zshrc file:\n")
		fmt.Printf("  fpath=(%s $fpath)\n", filepath.Dir(completionPath))
		fmt.Printf("  autoload -U compinit && compinit\n")
	case "fish":
		fmt.Printf("Fish will automatically load completions from %s\n", completionPath)
	case "powershell", "pwsh":
		profileFile := "$PROFILE"
		if shellType == "pwsh" {
			profileFile = "$PROFILE" // Same variable works for both
		}
		fmt.Printf("Add the following line to your PowerShell profile (%s):\n", profileFile)
		fmt.Printf("  . '%s'\n", completionPath)
	}

	return nil
}

// getPowerShellProfileDir returns the directory for the PowerShell profile
func getPowerShellProfileDir(isCore bool) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var profileDir string
	if runtime.GOOS == "windows" {
		// Windows uses different paths for PowerShell and PowerShell Core
		if isCore {
			// PowerShell Core on Windows
			profileDir = filepath.Join(homeDir, "Documents", "PowerShell")
		} else {
			// Windows PowerShell
			profileDir = filepath.Join(homeDir, "Documents", "WindowsPowerShell")
		}
	} else {
		// On Unix systems, both use ~/.config/powershell
		profileDir = filepath.Join(homeDir, ".config", "powershell")
	}

	return profileDir, nil
}
