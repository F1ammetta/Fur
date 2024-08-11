package main

import (
	"embed"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
)

//go:embed all:frontend/dist
var assets embed.FS

var path string

//go:embed build/appicon.png
var icon []byte

func main() {
	// get args
	args := os.Args

	if len(args) > 1 {
		if args[1] == "--path" {
			path = args[2]
		} else {
			path = "."
		}
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Fur",
		Width:  1024,
		Height: 768,
		Linux: &linux.Options{
			Icon:             icon,
			WebviewGpuPolicy: linux.WebviewGpuPolicyNever,
			ProgramName:      "Fur",
		},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Debug:            options.Debug{OpenInspectorOnStartup: true},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
