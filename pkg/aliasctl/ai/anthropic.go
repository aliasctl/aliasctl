package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

// AnthropicProvider implements Provider for Anthropic Claude.
type AnthropicProvider struct {
	Endpoint string // The Anthropic endpoint URL
	APIKey   string // The Anthropic API key
	Model    string // The Anthropic model name
}

// GenerateAlias generates an alias using Anthropic Claude
func (ap *AnthropicProvider) GenerateAlias(command, shellType string) (string, error) {
	if err := ValidateEndpoint(ap.Endpoint); err != nil {
		return "", err
	}

	// Check API key
	if ap.APIKey == "" {
		return "", fmt.Errorf("anthropic API key is empty: please configure a valid API key with 'aliasctl configure-anthropic'")
	}

	// Build the request payload
	requestBody, err := json.Marshal(map[string]any{
		"model": ap.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": GenerationPrompt(command, shellType),
			},
		},
		"max_tokens":  300,
		"temperature": 0.3, // Moderate creativity
	})
	if err != nil {
		return "", fmt.Errorf("failed to create Anthropic request: %w", err)
	}

	// Prepare headers
	headers := map[string]string{
		"x-api-key":         ap.APIKey,
		"anthropic-version": "2023-06-01", // Use appropriate API version
	}

	respBody, err := MakeAPIRequest("POST", ap.Endpoint+"/v1/messages", headers, requestBody)
	if err != nil {
		// Check for authentication errors
		if strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "403") {
			return "", fmt.Errorf("anthropic API authentication error: invalid API key. Check your API key or regenerate it in the Anthropic dashboard")
		}

		// Check for model errors
		if strings.Contains(err.Error(), "model") && strings.Contains(strings.ToLower(err.Error()), "not found") {
			return "", fmt.Errorf("anthropic model '%s' not found: check available models in your Anthropic account", ap.Model)
		}

		return "", fmt.Errorf("anthropic request failed: %w", err)
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Error struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse Anthropic response: %w\n\nRaw response: %s", err, limitResponseText(string(respBody), 200))
	}

	// Check if there's an error in the response
	if result.Error.Message != "" {
		return "", fmt.Errorf("anthropic API error: %s", result.Error.Message)
	}

	// Process the response content
	var responseText string
	for _, content := range result.Content {
		if content.Type == "text" {
			responseText = content.Text
			break
		}
	}

	if responseText == "" {
		return "", fmt.Errorf("no text response found in anthropic Claude reply\n\nRaw response: %s", limitResponseText(string(respBody), 200))
	}

	return ExtractAliasDefinition(responseText), nil
}

// ConvertAlias converts an alias using the Anthropic Claude API.
func (ap *AnthropicProvider) ConvertAlias(alias, fromShell, toShell string) (string, error) {
	if err := ValidateEndpoint(ap.Endpoint); err != nil {
		return "", err
	}

	// Check API key
	if ap.APIKey == "" {
		return "", fmt.Errorf("anthropic API key is empty: please configure a valid API key with 'aliasctl configure-anthropic'")
	}

	// Build the request payload
	requestBody, err := json.Marshal(map[string]any{
		"model": ap.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": ConversionPrompt(alias, fromShell, toShell),
			},
		},
		"max_tokens":  300,
		"temperature": 0.1, // Lower temperature for more deterministic results
	})
	if err != nil {
		return "", fmt.Errorf("failed to create Anthropic request: %w", err)
	}

	// Prepare headers
	headers := map[string]string{
		"x-api-key":         ap.APIKey,
		"anthropic-version": "2023-06-01", // Use appropriate API version
	}

	respBody, err := MakeAPIRequest("POST", ap.Endpoint+"/v1/messages", headers, requestBody)
	if err != nil {
		// Check for authentication errors
		if strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "403") {
			return "", fmt.Errorf("anthropic API authentication error: invalid API key. Check your API key or regenerate it in the Anthropic dashboard")
		}

		// Check for model errors
		if strings.Contains(err.Error(), "model") && strings.Contains(strings.ToLower(err.Error()), "not found") {
			return "", fmt.Errorf("anthropic model '%s' not found: check available models in your Anthropic account", ap.Model)
		}

		return "", fmt.Errorf("anthropic request failed: %w", err)
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Error struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse Anthropic response: %w\n\nRaw response: %s", err, limitResponseText(string(respBody), 200))
	}

	// Check if there's an error in the response
	if result.Error.Message != "" {
		return "", fmt.Errorf("anthropic API error: %s", result.Error.Message)
	}

	// Process the response content
	var responseText string
	for _, content := range result.Content {
		if content.Type == "text" {
			responseText = content.Text
			break
		}
	}

	if responseText == "" {
		return "", fmt.Errorf("no text response found in anthropic Claude reply\n\nRaw response: %s", limitResponseText(string(respBody), 200))
	}

	return ExtractAliasDefinition(responseText), nil
}
