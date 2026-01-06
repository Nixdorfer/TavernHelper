package main

import (
	"encoding/json"
	"strconv"
)

func (a *App) LoadConfig() (*Config, error) {
	config := &Config{
		Theme:          "default",
		Language:       "zh-CN",
		ColorScheme:    1,
		ColorMode:      "dark",
		SafeModeAction: "randomChars",
	}
	if theme, _ := DBGetConfig("theme"); theme != "" {
		config.Theme = theme
	}
	if lang, _ := DBGetConfig("language"); lang != "" {
		config.Language = lang
	}
	if lastProject, _ := DBGetConfig("lastOpenedProject"); lastProject != "" {
		config.LastOpenedProject = lastProject
	}
	if colorScheme, _ := DBGetConfig("colorScheme"); colorScheme != "" {
		if val, err := strconv.Atoi(colorScheme); err == nil {
			config.ColorScheme = val
		}
	}
	if colorMode, _ := DBGetConfig("colorMode"); colorMode != "" {
		config.ColorMode = colorMode
	}
	if systemPrompt, _ := DBGetConfig("systemPrompt"); systemPrompt != "" {
		config.SystemPrompt = systemPrompt
	}
	if systemPromptType, _ := DBGetConfig("systemPromptType"); systemPromptType != "" {
		config.SystemPromptType = systemPromptType
	} else {
		config.SystemPromptType = "system"
	}
	if debugMode, _ := DBGetConfig("debugMode"); debugMode == "true" {
		config.DebugMode = true
	}
	if safeMode, _ := DBGetConfig("safeMode"); safeMode == "true" {
		config.SafeMode = true
	}
	if safeModeAction, _ := DBGetConfig("safeModeAction"); safeModeAction != "" {
		config.SafeModeAction = safeModeAction
	}
	if safeModeTemplate, _ := DBGetConfig("safeModeTemplate"); safeModeTemplate != "" {
		config.SafeModeTemplate = safeModeTemplate
	}
	if debugTestReply, _ := DBGetConfig("debugTestReply"); debugTestReply != "" {
		config.DebugTestReply = debugTestReply
	}
	if bytedanceApiKey, _ := DBGetConfig("bytedanceApiKey"); bytedanceApiKey != "" {
		config.BytedanceApiKey = bytedanceApiKey
	}
	if autoGenerateImage, _ := DBGetConfig("autoGenerateImage"); autoGenerateImage == "true" {
		config.AutoGenerateImage = true
	}
	if noImageMode, _ := DBGetConfig("noImageMode"); noImageMode == "true" {
		config.NoImageMode = true
	}
	if strictMode, _ := DBGetConfig("strictMode"); strictMode == "true" {
		config.StrictMode = true
	}
	if claudeApiKey, _ := DBGetConfig("claudeApiKey"); claudeApiKey != "" {
		config.ClaudeApiKey = claudeApiKey
	}
	if geminiApiKey, _ := DBGetConfig("geminiApiKey"); geminiApiKey != "" {
		config.GeminiApiKey = geminiApiKey
	}
	if grokApiKey, _ := DBGetConfig("grokApiKey"); grokApiKey != "" {
		config.GrokApiKey = grokApiKey
	}
	return config, nil
}

func (a *App) SaveConfig(config Config) error {
	DBSetConfig("theme", config.Theme)
	DBSetConfig("language", config.Language)
	DBSetConfig("lastOpenedProject", config.LastOpenedProject)
	DBSetConfig("colorScheme", strconv.Itoa(config.ColorScheme))
	DBSetConfig("colorMode", config.ColorMode)
	DBSetConfig("systemPrompt", config.SystemPrompt)
	DBSetConfig("systemPromptType", config.SystemPromptType)
	if config.DebugMode {
		DBSetConfig("debugMode", "true")
		a.debugMode = true
	} else {
		DBSetConfig("debugMode", "false")
		a.debugMode = false
	}
	if config.SafeMode {
		DBSetConfig("safeMode", "true")
	} else {
		DBSetConfig("safeMode", "false")
	}
	DBSetConfig("safeModeAction", config.SafeModeAction)
	DBSetConfig("safeModeTemplate", config.SafeModeTemplate)
	DBSetConfig("debugTestReply", config.DebugTestReply)
	DBSetConfig("bytedanceApiKey", config.BytedanceApiKey)
	if config.AutoGenerateImage {
		DBSetConfig("autoGenerateImage", "true")
	} else {
		DBSetConfig("autoGenerateImage", "false")
	}
	if config.NoImageMode {
		DBSetConfig("noImageMode", "true")
	} else {
		DBSetConfig("noImageMode", "false")
	}
	if config.StrictMode {
		DBSetConfig("strictMode", "true")
	} else {
		DBSetConfig("strictMode", "false")
	}
	DBSetConfig("claudeApiKey", config.ClaudeApiKey)
	DBSetConfig("geminiApiKey", config.GeminiApiKey)
	DBSetConfig("grokApiKey", config.GrokApiKey)
	return nil
}

func (a *App) UpdateConfig(key string, value string) error {
	return DBSetConfig(key, value)
}

func (a *App) SaveSessionState(state SessionState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return DBSetConfig("sessionState", string(data))
}

func (a *App) LoadSessionState() (*SessionState, error) {
	stateStr, err := DBGetConfig("sessionState")
	if err != nil || stateStr == "" {
		return nil, nil
	}

	var state SessionState
	if err := json.Unmarshal([]byte(stateStr), &state); err != nil {
		return nil, err
	}
	return &state, nil
}
