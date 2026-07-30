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

func getLifeAdvice() (string, error) {
	reqBody := GroqRequest{
		Model: "openai/gpt-oss-120b",
		Messages: []Message{
			{Role: "system", Content: "You are a life coach and wisdom expert. Provide concise, practical life advice or words of wisdom. Keep responses short and actionable. Reply in one paragraph, at least few sentences. Make replies diverse and cover different aspects of life (health, relationships, career, personal growth, mindset, habits, etc.). Avoid repeating the same advice."},
			{Role: "user", Content: "Give me 10 pieces of life advice in 10 separate paragraphs. Each paragraph should contain at least a few sentences. Avoid repeating the same advice. Don't provide titles."},
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

func getDailyHoroscope() (string, error) {
	reqBody := GroqRequest{
		Model: "llama-3.3-70b-versatile",
		Messages: []Message{
			{Role: "system", Content: "You are an astrologer. Provide daily horoscope for each zodiac sign. Keep each horoscope around 40-50 words. Be diverse and creative."},
			{Role: "user", Content: "Give me daily horoscope for all 12 zodiac signs (Aries, Taurus, Gemini, Cancer, Leo, Virgo, Libra, Scorpio, Sagittarius, Capricorn, Aquarius, Pisces) in 12 separate paragraphs. Start each paragraph with the sign name followed by colon. Each horoscope should be around 40-50 words."},
		},
		Temperature: 0.8,
		MaxTokens:   1500,
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

func getPositiveNewsPost() (string, error) {
	reqBody := GroqRequest{
		Model: "llama-3.3-70b-versatile",
		Messages: []Message{
			{Role: "system", Content: "You are a very positive and optimistic journalist who only reports uplifting, feel-good recent news. Keep each news item 40-50 words. Make it fun and funny to read. Be diverse and cover different topics."},
			{Role: "user", Content: "Give me one very recent positive news.It should be 40-50 words and fun/funny to read. Focus on recent uplifting events, breakthroughs, heartwarming stories, and funny moments. Don't provide titles."},
		},
		Temperature: 0.8,
		MaxTokens:   1000,
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
