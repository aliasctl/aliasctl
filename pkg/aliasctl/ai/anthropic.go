package ai

import (
	"encoding/json"
	"fmt"
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

	// Build the request payload
	requestBody, err := json.Marshal(map[string]interface{}{
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
		return "", err
	}

	// Prepare headers
	headers := map[string]string{
		"x-api-key":         ap.APIKey,
		"anthropic-version": "2023-06-01", // Use appropriate API version
	}

	respBody, err := MakeAPIRequest("POST", ap.Endpoint+"/v1/messages", headers, requestBody)
	if err != nil {
		return "", fmt.Errorf("anthropic request failed: %v", err)
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
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
		return "", fmt.Errorf("no text response found in Anthropic Claude reply")
	}

	return ExtractAliasDefinition(responseText), nil
}

// ConvertAlias converts an alias using the Anthropic Claude API.
func (ap *AnthropicProvider) ConvertAlias(alias, fromShell, toShell string) (string, error) {
	if err := ValidateEndpoint(ap.Endpoint); err != nil {
		return "", err
	}

	// Build the request payload
	requestBody, err := json.Marshal(map[string]interface{}{
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
		return "", err
	}

	// Prepare headers
	headers := map[string]string{
		"x-api-key":         ap.APIKey,
		"anthropic-version": "2023-06-01", // Use appropriate API version
	}

	respBody, err := MakeAPIRequest("POST", ap.Endpoint+"/v1/messages", headers, requestBody)
	if err != nil {
		return "", fmt.Errorf("anthropic request failed: %v", err)
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
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
		return "", fmt.Errorf("no text response found in Anthropic Claude reply")
	}

	return ExtractAliasDefinition(responseText), nil
}
