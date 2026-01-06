package main

import (
	"context"
	"fmt"
	"os"
	"sync"
)

type App struct {
	ctx                   context.Context
	debugMode             bool
	clipboardMonitorStop  chan struct{}
	clipboardMonitorMutex sync.Mutex
	lastClipboardContent  string
	ignoredClipboard      map[string]bool
	dirtyManager          *DirtyManager
}

func NewApp() *App {
	return &App{
		dirtyManager: NewDirtyManager(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if err := InitDB(); err != nil {
		fmt.Printf("数据库初始化失败: %v\n", err)
		panic(err)
	}
	if os.Getenv("OPEN_DEVTOOLS") == "1" {
		a.debugMode = true
	}
}

func (a *App) IsDebugMode() bool {
	return a.debugMode
}

func (a *App) LogDebug(tag string, message string) {
	fmt.Printf("[%s] %s\n", tag, message)
}
