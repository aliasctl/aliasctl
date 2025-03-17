package ai

import (
	"fmt"
)

// Manager handles interactions with AI providers
type Manager struct {
	Providers map[string]Provider
	Default   Provider
}

// NewManager creates a new AI provider manager
func NewManager() *Manager {
	return &Manager{
		Providers: make(map[string]Provider),
	}
}

// AddProvider adds an AI provider to the manager
func (m *Manager) AddProvider(name string, provider Provider) {
	m.Providers[name] = provider
	if m.Default == nil {
		m.Default = provider
	}
}

// SetDefaultProvider sets the default AI provider
func (m *Manager) SetDefaultProvider(name string) error {
	if provider, exists := m.Providers[name]; exists {
		m.Default = provider
		return nil
	}
	return fmt.Errorf("provider '%s' not configured", name)
}

// GetProvider returns the named provider or the default if no name is specified
func (m *Manager) GetProvider(name string) (Provider, error) {
	if name == "" {
		if m.Default == nil {
			return nil, fmt.Errorf("no default AI provider configured")
		}
		return m.Default, nil
	}

	if provider, exists := m.Providers[name]; exists {
		return provider, nil
	}

	return nil, fmt.Errorf("AI provider '%s' not configured. Available providers: %v",
		name, m.ListProviders())
}

// ListProviders returns a list of configured provider names
func (m *Manager) ListProviders() []string {
	providers := make([]string, 0, len(m.Providers))
	for name := range m.Providers {
		providers = append(providers, name)
	}
	return providers
}

// ConvertAlias converts an alias from one shell to another using the specified provider
func (m *Manager) ConvertAlias(alias, fromShell, toShell, providerName string) (string, error) {
	provider, err := m.GetProvider(providerName)
	if err != nil {
		return "", err
	}
	return provider.ConvertAlias(alias, fromShell, toShell)
}

// GenerateAlias generates an alias suggestion for a command using the specified provider
func (m *Manager) GenerateAlias(command, shellType, providerName string) (string, error) {
	provider, err := m.GetProvider(providerName)
	if err != nil {
		return "", err
	}
	return provider.GenerateAlias(command, shellType)
}
