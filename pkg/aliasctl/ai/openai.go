package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

// OpenAIProvider implements Provider for OpenAI-compatible APIs.
type OpenAIProvider struct {
	Endpoint string // The OpenAI endpoint URL
	APIKey   string // The OpenAI API key
	Model    string // The OpenAI model name
}

// GenerateAlias generates an alias using OpenAI
func (op *OpenAIProvider) GenerateAlias(command, shellType string) (string, error) {
	if err := ValidateEndpoint(op.Endpoint); err != nil {
		return "", err
	}

	// Check API key
	if op.APIKey == "" {
		return "", fmt.Errorf("openAI API key is empty: please configure a valid API key with 'aliasctl configure-openai'")
	}

	requestBody, err := json.Marshal(map[string]any{
		"model": op.Model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": fmt.Sprintf("You are a shell alias creation expert for %s shell. Create concise, memorable aliases with proper syntax.", shellType),
			},
			{
				"role":    "user",
				"content": GenerationPrompt(command, shellType),
			},
		},
		"temperature": 0.3, // Moderate creativity
	})
	if err != nil {
		return "", fmt.Errorf("failed to create OpenAI request: %w", err)
	}

	// Prepare headers
	headers := map[string]string{
		"Authorization": "Bearer " + op.APIKey,
	}

	respBody, err := MakeAPIRequest("POST", op.Endpoint+"/v1/chat/completions", headers, requestBody)
	if err != nil {
		// Check for authentication errors
		if strings.Contains(err.Error(), "401") {
			return "", fmt.Errorf("openAI API authentication error: invalid API key. Check your API key or regenerate it in the OpenAI dashboard")
		}

		// Check for model errors
		if strings.Contains(err.Error(), "model") && strings.Contains(err.Error(), "does not exist") {
			return "", fmt.Errorf("openAI model '%s' not found: check available models in your OpenAI account", op.Model)
		}

		return "", fmt.Errorf("openAI request failed: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse OpenAI response: %w\n\nRaw response: %s", err, limitResponseText(string(respBody), 200))
	}

	// Check for error in response
	if errObj, hasErr := result["error"].(map[string]any); hasErr {
		errMsg := "unknown error"
		if msg, ok := errObj["message"].(string); ok {
			errMsg = msg
		}
		return "", fmt.Errorf("openAI API error: %s", errMsg)
	}

	if choices, ok := result["choices"].([]any); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]any); ok {
			if message, ok := choice["message"].(map[string]any); ok {
				if content, ok := message["content"].(string); ok {
					return ExtractAliasDefinition(content), nil
				}
			}
		}
	}

	return "", fmt.Errorf("unexpected response format from OpenAI: couldn't extract content from response\n\nResponse: %s", limitResponseText(string(respBody), 200))
}

// ConvertAlias converts an alias using the OpenAI-compatible API.
func (op *OpenAIProvider) ConvertAlias(alias, fromShell, toShell string) (string, error) {
	if err := ValidateEndpoint(op.Endpoint); err != nil {
		return "", err
	}

	// Check API key
	if op.APIKey == "" {
		return "", fmt.Errorf("openAI API key is empty: please configure a valid API key with 'aliasctl configure-openai'")
	}

	requestBody, err := json.Marshal(map[string]any{
		"model": op.Model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a utility that converts command line aliases between different shells. Respond only with the converted command, no explanation.",
			},
			{
				"role":    "user",
				"content": ConversionPrompt(alias, fromShell, toShell),
			},
		},
		"temperature": 0.2, // Lower temperature for more deterministic results
	})
	if err != nil {
		return "", fmt.Errorf("failed to create OpenAI request: %w", err)
	}

	// Prepare headers
	headers := map[string]string{
		"Authorization": "Bearer " + op.APIKey,
	}

	respBody, err := MakeAPIRequest("POST", op.Endpoint+"/v1/chat/completions", headers, requestBody)
	if err != nil {
		// Check for authentication errors
		if strings.Contains(err.Error(), "401") {
			return "", fmt.Errorf("openAI API authentication error: invalid API key. Check your API key or regenerate it in the OpenAI dashboard")
		}

		// Check for model errors
		if strings.Contains(err.Error(), "model") && strings.Contains(err.Error(), "does not exist") {
			return "", fmt.Errorf("openAI model '%s' not found: check available models in your OpenAI account", op.Model)
		}

		return "", fmt.Errorf("openAI request failed: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse OpenAI response: %w\n\nRaw response: %s", err, limitResponseText(string(respBody), 200))
	}

	// Check for error in response
	if errObj, hasErr := result["error"].(map[string]any); hasErr {
		errMsg := "unknown error"
		if msg, ok := errObj["message"].(string); ok {
			errMsg = msg
		}
		return "", fmt.Errorf("openAI API error: %s", errMsg)
	}

	if choices, ok := result["choices"].([]any); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]any); ok {
			if message, ok := choice["message"].(map[string]any); ok {
				if content, ok := message["content"].(string); ok {
					return ExtractAliasDefinition(content), nil
				}
			}
		}
	}

	return "", fmt.Errorf("unexpected response format from OpenAI: couldn't extract content from response\n\nResponse: %s", limitResponseText(string(respBody), 200))
}
