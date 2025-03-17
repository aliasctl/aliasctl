package ai

import (
	"encoding/json"
	"fmt"
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

	requestBody, err := json.Marshal(map[string]interface{}{
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
		return "", err
	}

	// Prepare headers
	headers := map[string]string{
		"Authorization": "Bearer " + op.APIKey,
	}

	respBody, err := MakeAPIRequest("POST", op.Endpoint+"/v1/chat/completions", headers, requestBody)
	if err != nil {
		return "", fmt.Errorf("OpenAI request failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return ExtractAliasDefinition(content), nil
				}
			}
		}
	}

	return "", fmt.Errorf("unexpected response format from OpenAI")
}

// ConvertAlias converts an alias using the OpenAI-compatible API.
func (op *OpenAIProvider) ConvertAlias(alias, fromShell, toShell string) (string, error) {
	if err := ValidateEndpoint(op.Endpoint); err != nil {
		return "", err
	}

	requestBody, err := json.Marshal(map[string]interface{}{
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
		return "", err
	}

	// Prepare headers
	headers := map[string]string{
		"Authorization": "Bearer " + op.APIKey,
	}

	respBody, err := MakeAPIRequest("POST", op.Endpoint+"/v1/chat/completions", headers, requestBody)
	if err != nil {
		return "", fmt.Errorf("OpenAI request failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return ExtractAliasDefinition(content), nil
				}
			}
		}
	}

	return "", fmt.Errorf("unexpected response format from OpenAI")
}
