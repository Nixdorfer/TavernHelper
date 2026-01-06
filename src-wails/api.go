package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) GetAppInfo(token, appID string) (map[string]any, error) {
	url := fmt.Sprintf("https://aipornhub.ltd/go/api/apps/%s", appID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Origin", "https://aipornhub.ltd")
	req.Header.Set("Referer", "https://aipornhub.ltd/")
	req.Header.Set("x-language", "zh-Hans")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (a *App) OpenFileDialog(title string, filters []runtime.FileFilter) (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   title,
		Filters: filters,
	})
}

func (a *App) SaveFileDialog(title string, defaultFilename string, filters []runtime.FileFilter) (string, error) {
	return runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           title,
		DefaultFilename: defaultFilename,
		Filters:         filters,
	})
}

func (a *App) FetchWithAuth(token, url, method, body string) (string, error) {
	var reqBody io.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Origin", "https://aipornhub.ltd")
	req.Header.Set("Referer", "https://aipornhub.ltd/")
	req.Header.Set("x-language", "zh-Hans")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

func (a *App) UploadCreationImage(token string) (string, error) {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择图片",
		Filters: []runtime.FileFilter{
			{DisplayName: "图片文件", Pattern: "*.png;*.jpg;*.jpeg;*.gif;*.webp"},
		},
	})
	if err != nil {
		return "", err
	}
	if filePath == "" {
		return "", nil
	}

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	boundary := "----WebKitFormBoundary" + fmt.Sprintf("%d", time.Now().UnixNano())
	var body bytes.Buffer
	body.WriteString("--" + boundary + "\r\n")
	body.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"file\"; filename=\"%s\"\r\n", filepath.Base(filePath)))
	body.WriteString("Content-Type: image/png\r\n\r\n")
	body.Write(fileData)
	body.WriteString("\r\n--" + boundary + "--\r\n")

	req, err := http.NewRequest("POST", "https://aipornhub.ltd/go/api/app_manage/upload_cover", &body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Origin", "https://aipornhub.ltd")
	req.Header.Set("Referer", "https://aipornhub.ltd/")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	if result.Code != 100000 {
		return "", fmt.Errorf("上传失败: %s", result.Msg)
	}

	return result.Data.URL, nil
}

func (a *App) UploadGalleryImageToRemote(imageId string, token string) (string, error) {
	images, err := DBGetAllGalleryImages()
	if err != nil {
		return "", err
	}

	var targetImage *GalleryImage
	for _, img := range images {
		if img.ID == imageId {
			targetImage = &img
			break
		}
	}

	if targetImage == nil {
		return "", fmt.Errorf("图片不存在")
	}

	fileData, err := os.ReadFile(targetImage.LocalPath)
	if err != nil {
		return "", fmt.Errorf("读取图片失败: %w", err)
	}

	boundary := "----WebKitFormBoundary" + fmt.Sprintf("%d", time.Now().UnixNano())
	var body bytes.Buffer
	body.WriteString("--" + boundary + "\r\n")
	body.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"file\"; filename=\"%s\"\r\n", filepath.Base(targetImage.LocalPath)))
	body.WriteString("Content-Type: image/png\r\n\r\n")
	body.Write(fileData)
	body.WriteString("\r\n--" + boundary + "--\r\n")

	req, err := http.NewRequest("POST", "https://aipornhub.ltd/go/api/app_manage/upload_cover", &body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Origin", "https://aipornhub.ltd")
	req.Header.Set("Referer", "https://aipornhub.ltd/")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	if result.Code != 100000 {
		return "", fmt.Errorf("上传失败: %s", result.Msg)
	}

	return result.Data.URL, nil
}

func (a *App) SaveCreationAppCache(userID, appID, appName, appIcon, modelConfig string) error {
	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec(`
		INSERT INTO wt_creation (id, user_id, app_name, app_icon, model_config, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id, user_id) DO UPDATE SET
			app_name = excluded.app_name,
			app_icon = excluded.app_icon,
			model_config = excluded.model_config,
			updated_at = excluded.updated_at
	`, appID, userID, appName, appIcon, modelConfig, now)
	return err
}

func (a *App) GetCreationAppCache(userID, appID string) (string, error) {
	var modelConfig string
	err := db.QueryRow(`SELECT model_config FROM wt_creation WHERE id = ? AND user_id = ?`, appID, userID).Scan(&modelConfig)
	if err != nil {
		return "", err
	}
	return modelConfig, nil
}

func (a *App) GetAllCreationAppCache(userID string) ([]map[string]any, error) {
	rows, err := db.Query(`SELECT id, app_name, app_icon, model_config, updated_at FROM wt_creation WHERE user_id = ? ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []map[string]any
	for rows.Next() {
		var id, appName, appIcon, modelConfig, updatedAt string
		if err := rows.Scan(&id, &appName, &appIcon, &modelConfig, &updatedAt); err != nil {
			continue
		}
		results = append(results, map[string]any{
			"id":           id,
			"app_name":     appName,
			"app_icon":     appIcon,
			"model_config": modelConfig,
			"updated_at":   updatedAt,
		})
	}
	return results, nil
}

func (a *App) DeleteCreationAppCache(userID, appID string) error {
	_, err := db.Exec(`DELETE FROM wt_creation WHERE id = ? AND user_id = ?`, appID, userID)
	return err
}

func (a *App) GetConversations(token, appID string, page, limit int) (map[string]any, error) {
	url := fmt.Sprintf("https://aipornhub.ltd/console/api/installed-apps/%s/conversations?page=%d&limit=%d", appID, page, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Origin", "https://aipornhub.ltd")
	req.Header.Set("Referer", "https://aipornhub.ltd/")
	req.Header.Set("x-language", "zh-Hans")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (a *App) DeleteConversation(token, appID, conversationID string) error {
	url := fmt.Sprintf("https://aipornhub.ltd/console/api/installed-apps/%s/conversations/%s", appID, conversationID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Origin", "https://aipornhub.ltd")
	req.Header.Set("Referer", "https://aipornhub.ltd/")
	req.Header.Set("x-language", "zh-Hans")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("删除对话失败: %s", string(respBody))
	}
	return nil
}

func (a *App) RenameConversation(token, appID, conversationID, newName string) error {
	url := fmt.Sprintf("https://aipornhub.ltd/console/api/installed-apps/%s/conversations/%s/name", appID, conversationID)

	bodyData := map[string]string{"name": newName}
	bodyJSON, err := json.Marshal(bodyData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Origin", "https://aipornhub.ltd")
	req.Header.Set("Referer", "https://aipornhub.ltd/")
	req.Header.Set("x-language", "zh-Hans")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("重命名对话失败: %s", string(respBody))
	}

	return nil
}

func (a *App) CreateNewConversation(token, appID, query, conversationName string) (map[string]any, error) {
	url := fmt.Sprintf("https://aipornhub.ltd/console/api/installed-apps/%s/chat-messages", appID)

	body := map[string]any{
		"response_mode":     "blocking",
		"conversation_id":   "",
		"query":             query,
		"inputs":            map[string]any{},
		"conversation_name": conversationName,
		"history_start_at":  nil,
	}

	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Origin", "https://aipornhub.ltd")
	req.Header.Set("Referer", "https://aipornhub.ltd/")
	req.Header.Set("x-language", "zh-Hans")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("创建对话失败: %s", string(respBody))
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err == nil {
		if convId, ok := result["conversation_id"].(string); ok && convId != "" {
			return map[string]any{
				"conversation_id": convId,
			}, nil
		}
	}

	return nil, fmt.Errorf("未能从响应中获取 conversation_id")
}

func (a *App) GetConversationDetail(token, appID, conversationID string) (map[string]any, error) {
	url := fmt.Sprintf("https://aipornhub.ltd/console/api/installed-apps/%s/conversations/%s", appID, conversationID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Origin", "https://aipornhub.ltd")
	req.Header.Set("Referer", "https://aipornhub.ltd/")
	req.Header.Set("x-language", "zh-Hans")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("获取会话详情失败: %s", string(respBody))
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (a *App) GetUserAppModelConfig(token, appID string) (map[string]any, error) {
	url := fmt.Sprintf("https://aipornhub.ltd/console/api/apps/%s/user_app_model_config", appID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-language", "zh-Hans")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("获取应用配置失败: %s", string(respBody))
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (a *App) ExportUserAppModelConfig(token, appID string) (map[string]any, error) {
	url := fmt.Sprintf("https://aipornhub.ltd/console/api/apps/%s/user_app_model_config/export", appID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-language", "zh-Hans")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("导出应用配置失败: %s", string(respBody))
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (a *App) UpdateUserAppModelConfig(token, appID string, payload map[string]any) (map[string]any, error) {
	url := fmt.Sprintf("https://aipornhub.ltd/console/api/apps/%s/user_app_model_config", appID)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://aipornhub.ltd")
	req.Header.Set("Referer", "https://aipornhub.ltd/")
	req.Header.Set("x-language", "zh-Hans")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("更新应用配置失败: %s", string(respBody))
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (a *App) GetCreationsDir() (string, error) {
	execPath, _ := os.Executable()
	creationsDir := filepath.Join(filepath.Dir(execPath), "creations")
	if err := os.MkdirAll(creationsDir, 0755); err != nil {
		return "", err
	}
	return creationsDir, nil
}

func (a *App) GetLocalCreations() ([]map[string]any, error) {
	creationsDir, err := a.GetCreationsDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(creationsDir)
	if err != nil {
		return nil, err
	}
	var results []map[string]any
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		projectDir := filepath.Join(creationsDir, entry.Name())
		configPath := filepath.Join(projectDir, "config.json")
		configData, err := os.ReadFile(configPath)
		if err != nil {
			continue
		}
		var config map[string]any
		if err := json.Unmarshal(configData, &config); err != nil {
			continue
		}
		config["folder_name"] = entry.Name()
		config["folder_path"] = projectDir
		results = append(results, config)
	}
	return results, nil
}

func (a *App) GetLocalCreationByRemoteId(remoteId string) (map[string]any, error) {
	creationsDir, err := a.GetCreationsDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(creationsDir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		projectDir := filepath.Join(creationsDir, entry.Name())
		configPath := filepath.Join(projectDir, "config.json")
		configData, err := os.ReadFile(configPath)
		if err != nil {
			continue
		}
		var config map[string]any
		if err := json.Unmarshal(configData, &config); err != nil {
			continue
		}
		if rid, ok := config["remote_id"].(string); ok && rid == remoteId {
			config["folder_name"] = entry.Name()
			config["folder_path"] = projectDir
			return config, nil
		}
	}
	return nil, nil
}

func (a *App) SaveLocalCreation(folderName string, config map[string]any) error {
	creationsDir, err := a.GetCreationsDir()
	if err != nil {
		return err
	}
	projectDir := filepath.Join(creationsDir, folderName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return err
	}
	pagesDir := filepath.Join(projectDir, "pages")
	if err := os.MkdirAll(pagesDir, 0755); err != nil {
		return err
	}
	configPath := filepath.Join(projectDir, "config.json")
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, configData, 0644)
}

func (a *App) RenameLocalCreation(oldName, newName string) error {
	creationsDir, err := a.GetCreationsDir()
	if err != nil {
		return err
	}
	oldPath := filepath.Join(creationsDir, oldName)
	newPath := filepath.Join(creationsDir, newName)
	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("目标文件夹已存在")
	}
	return os.Rename(oldPath, newPath)
}

func (a *App) DeleteLocalCreation(folderName string) error {
	creationsDir, err := a.GetCreationsDir()
	if err != nil {
		return err
	}
	projectDir := filepath.Join(creationsDir, folderName)
	return os.RemoveAll(projectDir)
}

func (a *App) GetLocalCreationConfig(folderName string) (map[string]any, error) {
	creationsDir, err := a.GetCreationsDir()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(creationsDir, folderName, "config.json")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config map[string]any
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, err
	}
	return config, nil
}

func (a *App) SaveLocalCreationPage(folderName, pageName, content string) error {
	creationsDir, err := a.GetCreationsDir()
	if err != nil {
		return err
	}
	pagesDir := filepath.Join(creationsDir, folderName, "pages")
	if err := os.MkdirAll(pagesDir, 0755); err != nil {
		return err
	}
	pagePath := filepath.Join(pagesDir, pageName+".html")
	return os.WriteFile(pagePath, []byte(content), 0644)
}

func (a *App) GetLocalCreationPage(folderName, pageName string) (string, error) {
	creationsDir, err := a.GetCreationsDir()
	if err != nil {
		return "", err
	}
	pagePath := filepath.Join(creationsDir, folderName, "pages", pageName+".html")
	content, err := os.ReadFile(pagePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (a *App) CreateNewLocalCreation(name string) (string, error) {
	creationsDir, err := a.GetCreationsDir()
	if err != nil {
		return "", err
	}
	folderName := fmt.Sprintf("%s_%d", name, time.Now().UnixNano())
	projectDir := filepath.Join(creationsDir, folderName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return "", err
	}
	pagesDir := filepath.Join(projectDir, "pages")
	if err := os.MkdirAll(pagesDir, 0755); err != nil {
		return "", err
	}
	config := map[string]any{
		"name":                      name,
		"description":               "",
		"language":                  "",
		"summary":                   "",
		"world_book":                []any{},
		"cg_book":                   []any{},
		"pre_text":                  "",
		"post_text":                 "",
		"is_anonymous":              false,
		"banned_words":              []any{},
		"pre_prompt":                nil,
		"cover":                     nil,
		"cover_tiny":                nil,
		"opening_statement":         nil,
		"suggested_questions":       []any{},
		"suggested_questions_after": map[string]any{"enabled": false},
		"bg_image":                  nil,
		"bg_mobile":                 nil,
		"category":                  0,
		"builtInCss":                nil,
		"mod_permission":            4,
		"shortcut_commands":         []any{},
		"is_available_not_public":   true,
		"preset_type":               1,
		"preset_chats":              []any{},
		"tags":                      []any{},
		"created_at":                time.Now().Format(time.RFC3339),
		"updated_at":                time.Now().Format(time.RFC3339),
	}
	configPath := filepath.Join(projectDir, "config.json")
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		return "", err
	}
	return folderName, nil
}

func (a *App) GenerateImageByteDance(apiKey, prompt string) (*ByteDanceImageResponse, error) {
	url := "https://ark.ap-southeast.bytepluses.com/api/v3/images/generations"

	body := map[string]any{
		"model":                       "seedream-4-0-250828",
		"prompt":                      prompt,
		"sequential_image_generation": "disabled",
		"response_format":             "url",
		"size":                        "2K",
		"stream":                      false,
		"watermark":                   true,
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result ByteDanceImageResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %s", string(respBody))
	}

	if result.Error != nil {
		return nil, fmt.Errorf("API错误: %s - %s", result.Error.Code, result.Error.Message)
	}

	return &result, nil
}
