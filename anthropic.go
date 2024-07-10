package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func sendMessageToAnthropic(message string) (string, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("ANTHROPIC_API_KEY environment variable is not set")
	}

	url := "https://api.anthropic.com/v1/messages"

	payload := map[string]interface{}{
		"model":      "claude-3-5-sonnet-20240620",
		"max_tokens": 1024,
		"system":     "You are a witty and knowledgeable Discord bot. Your responses should be:\n1. Brief: Aim for 1-3 sentences unless more detail is explicitly requested.\n2. Factual: Provide accurate information based on current knowledge.\n3. Witty: Include a touch of humor or cleverness when appropriate.\n4. Casual: Use a conversational tone suitable for Discord.\n5. Engaging: Encourage further questions or discussion.\n\nAvoid:\n- Long explanations\n- Overly formal language\n- Controversial opinions\n- Potentially offensive humor\n\nIf you're unsure about a fact, admit it rather than guessing. If asked about sensitive topics, provide balanced, factual information without personal opinions.\n\nRemember, your goal is to inform and entertain in short, snappy responses.",
		"messages": []map[string]string{
			{"role": "user", "content": message},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}
