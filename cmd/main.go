package main

import (
	"github.com/doai/papsmear/internal"
	"os"
)

func main() {
	envs := os.Getenv("path")
	curDir, _ := os.Getwd()
	envs = envs + curDir + "\bin;"
	os.Setenv("path", envs)
	for {
		internal.InsertNonTrackingSlides()
		internal.SendTileService()
	}
}
