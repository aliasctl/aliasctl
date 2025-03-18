package aliasctl

import (
	"fmt"

	"github.com/aliasctl/aliasctl/pkg/aliasctl/ai"
)

// InitAIProviders initializes AI providers from configuration.
// This sets up the AI manager and prepares it to handle provider requests.
// It should be called during AliasManager initialization.
func (am *AliasManager) InitAIProviders() {
	am.aiManager = ai.NewManager()
}

// ConfigureOllama sets up the Ollama AI provider.
// It creates and configures an Ollama provider with the specified endpoint and model,
// then adds it to the AI manager and sets it as the default provider.
// The configuration is saved after setup.
func (am *AliasManager) ConfigureOllama(endpoint, model string) {
	provider := &ai.OllamaProvider{
		Endpoint: endpoint,
		Model:    model,
	}

	if am.aiManager == nil {
		am.aiManager = ai.NewManager()
	}

	am.aiManager.AddProvider("ollama", provider)
	am.aiManager.SetDefaultProvider("ollama")
	am.AIConfigured = true
	am.SaveConfig()
}

// ConfigureOpenAI sets up the OpenAI-compatible AI provider.
// It creates and configures an OpenAI provider with the specified endpoint, API key, and model,
// then adds it to the AI manager and sets it as the default provider.
// The configuration is saved after setup.
func (am *AliasManager) ConfigureOpenAI(endpoint, apiKey, model string) {
	provider := &ai.OpenAIProvider{
		Endpoint: endpoint,
		APIKey:   apiKey,
		Model:    model,
	}

	if am.aiManager == nil {
		am.aiManager = ai.NewManager()
	}

	am.aiManager.AddProvider("openai", provider)
	am.aiManager.SetDefaultProvider("openai")
	am.AIConfigured = true
	am.SaveConfig()
}

// ConfigureAnthropic sets up the Anthropic Claude AI provider.
// It creates and configures an Anthropic provider with the specified endpoint, API key, and model,
// then adds it to the AI manager and sets it as the default provider.
// The configuration is saved after setup.
func (am *AliasManager) ConfigureAnthropic(endpoint, apiKey, model string) {
	provider := &ai.AnthropicProvider{
		Endpoint: endpoint,
		APIKey:   apiKey,
		Model:    model,
	}

	if am.aiManager == nil {
		am.aiManager = ai.NewManager()
	}

	am.aiManager.AddProvider("anthropic", provider)
	am.aiManager.SetDefaultProvider("anthropic")
	am.AIConfigured = true
	am.SaveConfig()
}

// GetAvailableProviders returns a list of configured AI provider names.
// It queries the AI manager for all registered providers.
// Returns an empty slice if no providers are configured.
func (am *AliasManager) GetAvailableProviders() []string {
	if am.aiManager == nil {
		return []string{}
	}
	return am.aiManager.ListProviders()
}

// ConvertAlias converts an alias from one shell to another using the specified provider.
// It retrieves the alias definition for the current shell and asks the AI to convert it
// to the target shell format.
// Returns an error if the alias doesn't exist, no AI provider is configured, or the conversion fails.
func (am *AliasManager) ConvertAlias(name, targetShell, providerName string) (string, error) {
	if !am.AIConfigured {
		return "", fmt.Errorf("AI provider not configured. Use 'aliasctl configure-ollama', 'aliasctl configure-openai', or 'aliasctl configure-anthropic' to set up an AI provider")
	}

	commands, exists := am.Aliases[name]
	if !exists {
		return "", fmt.Errorf("alias '%s' not found. Run 'aliasctl list' to see available aliases", name)
	}

	var command string
	switch am.Shell {
	case ShellBash:
		command = commands.Bash
	case ShellZsh:
		command = commands.Zsh
	case ShellFish:
		command = commands.Fish
	case ShellKsh:
		command = commands.Ksh
	case ShellPowerShell:
		command = commands.PowerShell
	case ShellPowerShellCore:
		command = commands.PowerShellCore
	case ShellCmd:
		command = commands.Cmd
	}

	if command == "" {
		return "", fmt.Errorf("command for shell '%s' not found", am.Shell)
	}

	return am.aiManager.ConvertAlias(command, string(am.Shell), targetShell, providerName)
}

// GenerateAlias generates an alias suggestion for the given command.
// It uses the configured AI provider to suggest a shell-appropriate alias name and format
// for the provided command.
// Returns an error if no AI provider is configured or the generation fails.
func (am *AliasManager) GenerateAlias(command, providerName string) (string, error) {
	if !am.AIConfigured {
		return "", fmt.Errorf("AI provider not configured. Use 'aliasctl configure-ollama', 'aliasctl configure-openai', or 'aliasctl configure-anthropic' to set up an AI provider")
	}

	return am.aiManager.GenerateAlias(command, string(am.Shell), providerName)
}
