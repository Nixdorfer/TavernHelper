package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type TestChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type TestChatRequest struct {
	Provider     string            `json:"provider"`
	Messages     []TestChatMessage `json:"messages"`
	SystemPrompt string            `json:"systemPrompt"`
	WorldBook    []WorldBookEntry  `json:"worldBook"`
}

type TestChatResponse struct {
	Content string `json:"content"`
	Error   string `json:"error,omitempty"`
}

func (a *App) SendTestChat(request TestChatRequest) TestChatResponse {
	systemPrompt := request.SystemPrompt
	if len(request.WorldBook) > 0 {
		matchedEntries := a.matchWorldBookEntries(request.WorldBook, request.Messages)
		if matchedEntries != "" {
			systemPrompt = systemPrompt + "\n\n" + matchedEntries
		}
	}
	var content string
	var err error
	switch request.Provider {
	case "claude":
		content, err = a.callClaudeAPI(request.Messages, systemPrompt)
	case "gemini":
		content, err = a.callGeminiAPI(request.Messages, systemPrompt)
	case "grok":
		content, err = a.callGrokAPI(request.Messages, systemPrompt)
	default:
		return TestChatResponse{Error: "未知的 API 提供商"}
	}
	if err != nil {
		return TestChatResponse{Error: err.Error()}
	}
	return TestChatResponse{Content: content}
}

func (a *App) matchWorldBookEntries(worldBook []WorldBookEntry, messages []TestChatMessage) string {
	var allText strings.Builder
	for _, msg := range messages {
		allText.WriteString(msg.Content)
		allText.WriteString(" ")
	}
	searchText := strings.ToLower(allText.String())
	var matched []string
	for _, entry := range worldBook {
		if entry.IsFolder {
			continue
		}
		keywords := entry.Keywords
		if len(keywords) == 0 && entry.Key != "" {
			keywords = strings.Split(entry.Key, ",")
		}
		for _, kw := range keywords {
			kw = strings.TrimSpace(strings.ToLower(kw))
			if kw != "" && strings.Contains(searchText, kw) {
				if entry.Value != "" {
					matched = append(matched, entry.Value)
				}
				break
			}
		}
	}
	return strings.Join(matched, "\n\n")
}

func (a *App) callClaudeAPI(messages []TestChatMessage, systemPrompt string) (string, error) {
	config, err := a.LoadConfig()
	if err != nil || config.ClaudeApiKey == "" {
		return "", fmt.Errorf("Claude API Key 未配置")
	}
	claudeMessages := make([]map[string]string, len(messages))
	for i, msg := range messages {
		claudeMessages[i] = map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		}
	}
	requestBody := map[string]any{
		"model":      "claude-sonnet-4-20250514",
		"max_tokens": 4096,
		"messages":   claudeMessages,
	}
	if systemPrompt != "" {
		requestBody["system"] = systemPrompt
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.ClaudeApiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Claude API 错误: %s", string(body))
	}
	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if result.Error != nil {
		return "", fmt.Errorf(result.Error.Message)
	}
	if len(result.Content) > 0 {
		return result.Content[0].Text, nil
	}
	return "", fmt.Errorf("Claude 返回空响应")
}

func (a *App) callGeminiAPI(messages []TestChatMessage, systemPrompt string) (string, error) {
	config, err := a.LoadConfig()
	if err != nil || config.GeminiApiKey == "" {
		return "", fmt.Errorf("Gemini API Key 未配置")
	}
	var contents []map[string]any
	for _, msg := range messages {
		role := msg.Role
		if role == "assistant" {
			role = "model"
		}
		contents = append(contents, map[string]any{
			"role": role,
			"parts": []map[string]string{
				{"text": msg.Content},
			},
		})
	}
	requestBody := map[string]any{
		"contents": contents,
	}
	if systemPrompt != "" {
		requestBody["systemInstruction"] = map[string]any{
			"parts": []map[string]string{
				{"text": systemPrompt},
			},
		}
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", config.GeminiApiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Gemini API 错误: %s", string(body))
	}
	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if result.Error != nil {
		return "", fmt.Errorf(result.Error.Message)
	}
	if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
		return result.Candidates[0].Content.Parts[0].Text, nil
	}
	return "", fmt.Errorf("Gemini 返回空响应")
}

func (a *App) callGrokAPI(messages []TestChatMessage, systemPrompt string) (string, error) {
	config, err := a.LoadConfig()
	if err != nil || config.GrokApiKey == "" {
		return "", fmt.Errorf("Grok API Key 未配置")
	}
	var grokMessages []map[string]string
	if systemPrompt != "" {
		grokMessages = append(grokMessages, map[string]string{
			"role":    "system",
			"content": systemPrompt,
		})
	}
	for _, msg := range messages {
		grokMessages = append(grokMessages, map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}
	requestBody := map[string]any{
		"model":    "grok-3-latest",
		"messages": grokMessages,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://api.x.ai/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.GrokApiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Grok API 错误: %s", string(body))
	}
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if result.Error != nil {
		return "", fmt.Errorf(result.Error.Message)
	}
	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("Grok 返回空响应")
}
