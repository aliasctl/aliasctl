package ai

import (
	"fmt"
	"strings"
)

// Manager handles interactions with AI providers.
// It maintains a registry of available providers and handles provider selection.
type Manager struct {
	Providers map[string]Provider // Map of provider name to provider implementation
	Default   Provider            // The default provider to use when none is specified
}

// NewManager creates a new AI provider manager.
// It initializes an empty provider registry.
func NewManager() *Manager {
	return &Manager{
		Providers: make(map[string]Provider),
	}
}

// AddProvider adds an AI provider to the manager.
// If this is the first provider added, it becomes the default.
func (m *Manager) AddProvider(name string, provider Provider) {
	m.Providers[name] = provider
	if m.Default == nil {
		m.Default = provider
	}
}

// SetDefaultProvider sets the default AI provider.
// It returns an error if the specified provider is not registered.
func (m *Manager) SetDefaultProvider(name string) error {
	if provider, exists := m.Providers[name]; exists {
		m.Default = provider
		return nil
	}
	return fmt.Errorf("provider '%s' not configured. Available providers: %s", name, strings.Join(m.ListProviders(), ", "))
}

// GetProvider returns the named provider or the default if no name is specified.
// It returns an error if the named provider doesn't exist or if no default is set.
// The error message includes guidance on how to configure providers.
func (m *Manager) GetProvider(name string) (Provider, error) {
	if name == "" {
		if m.Default == nil {
			return nil, fmt.Errorf("no default AI provider configured. Use one of: 'aliasctl configure-ollama', 'aliasctl configure-openai', or 'aliasctl configure-anthropic'")
		}
		return m.Default, nil
	}

	if provider, exists := m.Providers[name]; exists {
		return provider, nil
	}

	providers := m.ListProviders()
	if len(providers) == 0 {
		return nil, fmt.Errorf("AI provider '%s' not configured and no providers are available\n\nTo configure a provider, use one of:\n  - aliasctl configure-ollama <endpoint> <model>\n  - aliasctl configure-openai <endpoint> <api-key> <model>\n  - aliasctl configure-anthropic <endpoint> <api-key> <model>", name)
	}

	return nil, fmt.Errorf("AI provider '%s' not configured\n\nAvailable providers: %s\n\nTo use a specific provider, specify it with the --provider flag", name, strings.Join(providers, ", "))
}

// ListProviders returns a list of configured provider names.
// The returned list is alphabetically sorted for consistent presentation.
func (m *Manager) ListProviders() []string {
	providers := make([]string, 0, len(m.Providers))
	for name := range m.Providers {
		providers = append(providers, name)
	}
	return providers
}

// ConvertAlias converts an alias from one shell to another using the specified provider.
// It automatically selects the default provider if none is specified.
// Returns the converted alias or an error if the conversion fails.
func (m *Manager) ConvertAlias(alias, fromShell, toShell, providerName string) (string, error) {
	provider, err := m.GetProvider(providerName)
	if err != nil {
		return "", err
	}

	result, err := provider.ConvertAlias(alias, fromShell, toShell)
	if err != nil {
		// Add more context to the error
		return "", fmt.Errorf("failed to convert alias from %s to %s: %w", fromShell, toShell, err)
	}

	return result, nil
}

// GenerateAlias generates an alias suggestion for a command using the specified provider.
// It automatically selects the default provider if none is specified.
// Returns the generated alias suggestion or an error if the generation fails.
func (m *Manager) GenerateAlias(command, shellType, providerName string) (string, error) {
	provider, err := m.GetProvider(providerName)
	if err != nil {
		return "", err
	}

	result, err := provider.GenerateAlias(command, shellType)
	if err != nil {
		// Add more context to the error
		return "", fmt.Errorf("failed to generate alias suggestion for %s shell: %w", shellType, err)
	}

	return result, nil
}
