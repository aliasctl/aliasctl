package aliasctl

import (
	"fmt"

	"github.com/aliasctl/aliasctl/pkg/aliasctl/ai"
)

// InitAIProviders initializes AI providers from configuration
func (am *AliasManager) InitAIProviders() {
	am.aiManager = ai.NewManager()
}

// ConfigureOllama sets up the Ollama AI provider.
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

// GetAvailableProviders returns a list of configured AI provider names
func (am *AliasManager) GetAvailableProviders() []string {
	if am.aiManager == nil {
		return []string{}
	}
	return am.aiManager.ListProviders()
}

// ConvertAlias converts an alias from one shell to another using the specified provider
func (am *AliasManager) ConvertAlias(name, targetShell, providerName string) (string, error) {
	if !am.AIConfigured {
		return "", fmt.Errorf("AI provider not configured")
	}

	commands, exists := am.Aliases[name]
	if !exists {
		return "", fmt.Errorf("alias '%s' not found", name)
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

// GenerateAlias generates an alias suggestion for the given command
func (am *AliasManager) GenerateAlias(command, providerName string) (string, error) {
	if !am.AIConfigured {
		return "", fmt.Errorf("AI provider not configured")
	}

	return am.aiManager.GenerateAlias(command, string(am.Shell), providerName)
}
