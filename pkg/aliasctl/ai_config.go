package aliasctl

// ConfigureOllama sets up the Ollama AI provider.
func (am *AliasManager) ConfigureOllama(endpoint, model string) {
	provider := &OllamaProvider{
		Endpoint: endpoint,
		Model:    model,
	}
	am.AIProvider = provider // For backward compatibility
	am.AIProviders["ollama"] = provider
	am.AIConfigured = true
	am.SaveConfig()
}

// ConfigureOpenAI sets up the OpenAI-compatible AI provider.
func (am *AliasManager) ConfigureOpenAI(endpoint, apiKey, model string) {
	provider := &OpenAIProvider{
		Endpoint: endpoint,
		APIKey:   apiKey,
		Model:    model,
	}
	am.AIProvider = provider // For backward compatibility
	am.AIProviders["openai"] = provider
	am.AIConfigured = true
	am.SaveConfig()
}

// ConfigureAnthropic sets up the Anthropic Claude AI provider.
func (am *AliasManager) ConfigureAnthropic(endpoint, apiKey, model string) {
	provider := &AnthropicProvider{
		Endpoint: endpoint,
		APIKey:   apiKey,
		Model:    model,
	}
	am.AIProvider = provider // For backward compatibility
	am.AIProviders["anthropic"] = provider
	am.AIConfigured = true
	am.SaveConfig()
}

// GetAvailableProviders returns a list of configured AI provider names
func (am *AliasManager) GetAvailableProviders() []string {
	providers := make([]string, 0, len(am.AIProviders))
	for name := range am.AIProviders {
		providers = append(providers, name)
	}
	return providers
}
