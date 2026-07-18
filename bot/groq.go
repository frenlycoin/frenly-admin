package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type GroqRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func getProgrammingAdvice() (string, error) {
	reqBody := GroqRequest{
		Model: "openai/gpt-oss-120b",
		Messages: []Message{
			{Role: "system", Content: "You are a programming expert. Provide concise, practical programming advice or best practices. Keep responses short and actionable. Reply in one paragraph, at least few sentences. Do not include any code examples. Make replies diverse and cover different programming topics. Avoid repeating the same advice."},
			{Role: "user", Content: "Give me 10 pieces of programming advice in 10 separate paragraphs. Each paragraph should contain at least a few sentences. Avoid repeating the same advice. Don't provide titles."},
		},
		Temperature: 0.7,
		MaxTokens:   700,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+conf.GroqAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var groqResp GroqResponse
	if err := json.NewDecoder(resp.Body).Decode(&groqResp); err != nil {
		return "", err
	}

	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf("empty response from API")
	}

	return groqResp.Choices[0].Message.Content, nil
}
