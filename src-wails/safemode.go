package main

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) GetSafeModeTxtFiles() ([]string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	exeDir := filepath.Dir(exePath)
	textDir := filepath.Join(exeDir, "text")

	if _, err := os.Stat(textDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	files, err := os.ReadDir(textDir)
	if err != nil {
		return nil, err
	}

	var txtFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".txt") {
			txtFiles = append(txtFiles, file.Name())
		}
	}
	return txtFiles, nil
}

func (a *App) ReadSafeModeTemplate(filename string) (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exePath)
	filePath := filepath.Join(exeDir, "text", filename)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

var (
	user32DLL   = syscall.NewLazyDLL("user32.dll")
	getKeyState = user32DLL.NewProc("GetKeyState")
)

const VK_CAPITAL = 0x14

func (a *App) IsCapsLockOn() bool {
	ret, _, _ := getKeyState.Call(uintptr(VK_CAPITAL))
	return ret&1 != 0
}

func (a *App) HideWindow() {
	wailsRuntime.WindowHide(a.ctx)
}

func (a *App) ShowWindow() {
	wailsRuntime.WindowShow(a.ctx)
}

func (a *App) ReadPromptFile(filename string) (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exePath)
	filePath := filepath.Join(exeDir, "..", "..", "prompts", filename)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (a *App) WritePromptFile(filename string, content string) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)
	promptsDir := filepath.Join(exeDir, "..", "..", "prompts")

	if err := os.MkdirAll(promptsDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(promptsDir, filename)
	return os.WriteFile(filePath, []byte(content), 0644)
}
