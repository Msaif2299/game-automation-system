package main

import (
	"aqw-gobot/backend/bot"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	bot := bot.NewBot()

	// Create application with options
	err := wails.Run(&options.App{
		Title:         "aqw-gobot",
		Width:         600,
		Height:        300,
		DisableResize: true,
		Frameless:     true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:        bot.Startup,
		OnShutdown:       bot.Shutdown,
		Bind: []interface{}{
			bot,
		},
	})

	if err != nil {
		println("[ERROR] Error:", err.Error())
	}
}
