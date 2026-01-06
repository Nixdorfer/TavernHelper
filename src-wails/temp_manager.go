package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type TempFileInfo struct {
	Filename     string `json:"filename"`
	Type         string `json:"type"`
	ID           string `json:"id"`
	Size         int64  `json:"size"`
	CreatedAt    string `json:"createdAt"`
	LastModified string `json:"lastModified"`
}

type SessionRecoveryInfo struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	LastModified string `json:"lastModified"`
	Preview      string `json:"preview"`
	FilePath     string `json:"filePath"`
}

func (a *App) GetTempDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	tempDir := filepath.Join(filepath.Dir(execPath), "temp")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", err
	}
	return tempDir, nil
}

func (a *App) SaveTempFile(filename string, content string) error {
	tempDir, err := a.GetTempDir()
	if err != nil {
		return err
	}
	filePath := filepath.Join(tempDir, filename)
	return os.WriteFile(filePath, []byte(content), 0644)
}

func (a *App) LoadTempFile(filename string) (string, error) {
	tempDir, err := a.GetTempDir()
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(tempDir, filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (a *App) DeleteTempFile(filename string) error {
	tempDir, err := a.GetTempDir()
	if err != nil {
		return err
	}
	filePath := filepath.Join(tempDir, filename)
	return os.Remove(filePath)
}

func (a *App) ListTempFiles() ([]TempFileInfo, error) {
	tempDir, err := a.GetTempDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return nil, err
	}
	var files []TempFileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		filename := entry.Name()
		fileType := ""
		fileID := ""
		if strings.HasPrefix(filename, "creation_") {
			fileType = "creation"
			parts := strings.Split(strings.TrimSuffix(filename, ".json"), "_")
			if len(parts) >= 2 {
				fileID = parts[1]
			}
		} else if strings.HasPrefix(filename, "worldtree_") {
			fileType = "worldtree"
			parts := strings.Split(strings.TrimSuffix(filename, ".json"), "_")
			if len(parts) >= 2 {
				fileID = parts[1]
			}
		} else if strings.HasPrefix(filename, "session_") {
			fileType = "session"
		}
		files = append(files, TempFileInfo{
			Filename:     filename,
			Type:         fileType,
			ID:           fileID,
			Size:         info.Size(),
			CreatedAt:    info.ModTime().Format(time.RFC3339),
			LastModified: info.ModTime().Format(time.RFC3339),
		})
	}
	return files, nil
}

func (a *App) CleanupTempFiles(olderThanMinutes int) (int, error) {
	tempDir, err := a.GetTempDir()
	if err != nil {
		return 0, err
	}
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return 0, err
	}
	cutoff := time.Now().Add(-time.Duration(olderThanMinutes) * time.Minute)
	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			filePath := filepath.Join(tempDir, entry.Name())
			if err := os.Remove(filePath); err == nil {
				count++
			}
		}
	}
	return count, nil
}

func (a *App) GetUnsavedSessions() ([]SessionRecoveryInfo, error) {
	tempDir, err := a.GetTempDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return nil, err
	}
	var sessions []SessionRecoveryInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filename := entry.Name()
		if !strings.HasSuffix(filename, ".json") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		filePath := filepath.Join(tempDir, filename)
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}
		var tempData map[string]any
		if err := json.Unmarshal(data, &tempData); err != nil {
			continue
		}
		session := SessionRecoveryInfo{
			LastModified: info.ModTime().Format(time.RFC3339),
			FilePath:     filePath,
		}
		if strings.HasPrefix(filename, "creation_") {
			session.Type = "creation"
			parts := strings.Split(strings.TrimSuffix(filename, ".json"), "_")
			if len(parts) >= 2 {
				session.ID = parts[1]
			}
			if name, ok := tempData["name"].(string); ok {
				session.Name = name
			} else {
				session.Name = "未命名创作"
			}
			if desc, ok := tempData["description"].(string); ok && len(desc) > 100 {
				session.Preview = desc[:100] + "..."
			} else if desc, ok := tempData["description"].(string); ok {
				session.Preview = desc
			}
		} else if strings.HasPrefix(filename, "worldtree_") {
			session.Type = "worldtree"
			parts := strings.Split(strings.TrimSuffix(filename, ".json"), "_")
			if len(parts) >= 2 {
				session.ID = parts[1]
			}
			if name, ok := tempData["name"].(string); ok {
				session.Name = name
			} else {
				session.Name = "未命名世界树"
			}
		} else {
			continue
		}
		sessions = append(sessions, session)
	}
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].LastModified > sessions[j].LastModified
	})
	return sessions, nil
}

func (a *App) RecoverTempSession(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取临时文件失败: %w", err)
	}
	return string(data), nil
}

func (a *App) DiscardTempSession(filePath string) error {
	return os.Remove(filePath)
}

func (a *App) DiscardAllTempSessions() error {
	tempDir, err := a.GetTempDir()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filename := entry.Name()
		if strings.HasPrefix(filename, "creation_") || strings.HasPrefix(filename, "worldtree_") {
			filePath := filepath.Join(tempDir, filename)
			os.Remove(filePath)
		}
	}
	return nil
}

func (a *App) SaveCreationTemp(creationId string, content string) error {
	filename := fmt.Sprintf("creation_%s_%d.json", creationId, time.Now().UnixNano())
	tempDir, err := a.GetTempDir()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(tempDir)
	if err == nil {
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), fmt.Sprintf("creation_%s_", creationId)) {
				os.Remove(filepath.Join(tempDir, entry.Name()))
			}
		}
	}
	return a.SaveTempFile(filename, content)
}

func (a *App) SaveWorldTreeTemp(projectId string, content string) error {
	filename := fmt.Sprintf("worldtree_%s_%d.json", projectId, time.Now().UnixNano())
	tempDir, err := a.GetTempDir()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(tempDir)
	if err == nil {
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), fmt.Sprintf("worldtree_%s_", projectId)) {
				os.Remove(filepath.Join(tempDir, entry.Name()))
			}
		}
	}
	return a.SaveTempFile(filename, content)
}

func (a *App) ClearCreationTemp(creationId string) error {
	tempDir, err := a.GetTempDir()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), fmt.Sprintf("creation_%s_", creationId)) {
			os.Remove(filepath.Join(tempDir, entry.Name()))
		}
	}
	return nil
}

func (a *App) ClearWorldTreeTemp(projectId string) error {
	tempDir, err := a.GetTempDir()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), fmt.Sprintf("worldtree_%s_", projectId)) {
			os.Remove(filepath.Join(tempDir, entry.Name()))
		}
	}
	return nil
}
