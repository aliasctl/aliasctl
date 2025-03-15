package aliasctl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// ConvertAlias converts an alias from one shell to another using the configured AI provider.
func (am *AliasManager) ConvertAlias(name, targetShell string) (string, error) {
	if !am.AIConfigured {
		return "", fmt.Errorf("AI provider not configured")
	}

	command, exists := am.Aliases[name]
	if !exists {
		return "", fmt.Errorf("alias '%s' not found", name)
	}

	return am.AIProvider.ConvertAlias(command, string(am.Shell), targetShell)
}

// ConvertAlias for OllamaProvider converts an alias using the Ollama AI service.
func (op *OllamaProvider) ConvertAlias(alias, fromShell, toShell string) (string, error) {
	prompt := fmt.Sprintf("Convert the following command from %s shell to %s shell: %s", fromShell, toShell, alias)

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":  op.Model,
		"prompt": prompt,
	})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(op.Endpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if response, ok := result["response"].(string); ok {
		return response, nil
	}

	return "", fmt.Errorf("unexpected response format from Ollama")
}

// ConvertAlias for OpenAIProvider converts an alias using the OpenAI-compatible API.
func (op *OpenAIProvider) ConvertAlias(alias, fromShell, toShell string) (string, error) {
	prompt := fmt.Sprintf("Convert the following command from %s shell to %s shell: %s", fromShell, toShell, alias)

	requestBody, err := json.Marshal(map[string]interface{}{
		"model": op.Model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are an assistant that converts command line aliases between different shells.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", op.Endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+op.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return content, nil
				}
			}
		}
	}

	return "", fmt.Errorf("unexpected response format from OpenAI")
}
