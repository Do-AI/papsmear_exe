package main

import (
	"github.com/doai/papsmear/internal"
)

func main() {
	for {
		internal.InsertNonTrackingSlides()
		internal.SendTileService()
	}
}
