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

func getTonNews() (string, error) {
	reqBody := GroqRequest{
		Model: "openai/gpt-oss-120b",
		Messages: []Message{
			{Role: "system", Content: "You are a crypto and tech news reporter. Provide concise, up-to-date news about TON (The Open Network), Telegram, Pavel Durov, Gram, and general crypto topics. Keep responses short and informative. Reply in one paragraph, at least few sentences. Make replies diverse and cover different news topics. Avoid repeating the same news."},
			{Role: "user", Content: "Give me 10 latest news items about TON, Telegram, Durov, Gram and general crypto in 10 separate paragraphs. Each paragraph should contain at least a few sentences. Avoid repeating the same news. Don't provide titles."},
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

func getCryptoNews() (string, error) {
	reqBody := GroqRequest{
		Model: "openai/gpt-oss-120b",
		Messages: []Message{
			{Role: "system", Content: "You are a crypto news reporter who specializes in the latest developments in the cryptocurrency and blockchain industry. You report on real, specific events and market movements with factual precision. Cover topics including Bitcoin, Ethereum, altcoins, DeFi protocols, NFT market trends, regulatory updates, exchange developments, institutional adoption, and major blockchain upgrades. Write in a professional news style. Each news item must be a self-contained paragraph of 2-4 sentences with specific details like names, numbers, and dates where appropriate."},
			{Role: "user", Content: "Write exactly 10 distinct crypto news items covering recent major developments in the cryptocurrency world. Separate each item by a blank line. Each item must be a substantive paragraph of 2-4 sentences with specific details. Cover diverse topics: Bitcoin, Ethereum, DeFi, NFTs, regulation, exchange news, institutional moves, Layer-2 projects, meme coins, and blockchain upgrades. Do NOT number them and do NOT include titles. Just the paragraph text separated by blank lines. You must produce all 10 items — do not stop early."},
		},
		Temperature: 0.7,
		MaxTokens:   2000,
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
