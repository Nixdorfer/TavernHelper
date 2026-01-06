package main

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:dist
var assets embed.FS

func main() {
	openDevTools := os.Getenv("OPEN_DEVTOOLS") == "1"
	app := NewApp()

	distFS, _ := fs.Sub(assets, "dist")

	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	webviewDataPath := filepath.Join(exeDir, "webview_data")

	err := wails.Run(&options.App{
		Title:     "大纲管理器",
		Width:     1200,
		Height:    800,
		MinWidth:  1250,
		MinHeight: 800,
		AssetServer: &assetserver.Options{
			Assets: distFS,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []any{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent:  false,
			WindowIsTranslucent:   false,
			DisableWindowIcon:     false,
			WebviewUserDataPath:   webviewDataPath,
		},
		Debug: options.Debug{
			OpenInspectorOnStartup: openDevTools,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
