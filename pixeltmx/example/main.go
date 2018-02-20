package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/elliotmr/tmx"
	"github.com/elliotmr/tmx/pixeltmx"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	mapFile, err := os.Open("orthogonal-outside.tmx")
	if err != nil {
		panic(err)
	}
	mapData, err := tmx.Load(mapFile)
	if err != nil {
		panic(err)
	}
	resources, err := pixeltmx.LoadResources(mapData, "")
	if err != nil {
		panic(err)
	}

	drawer, err := pixeltmx.NewRootDrawer(resources, mapData)
	if err != nil {
		panic(err)
	}

	cfg := pixelgl.WindowConfig{
		Title:  "TMX Map Example",
		Bounds: pixel.R(0, 0, 1024, 768),
		//VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetSmooth(false)
	cameraOrigin := pixel.ZV.Add(win.Bounds().Center())
	scale := 1.0
	dragOrigin := pixel.V(0, 0)
	second := time.Tick(time.Second)
	viewMatrix := pixel.IM
	frames := 0
	for !win.Closed() {
		if win.MouseScroll().Y != 0 {
			factor := math.Pow(1.2, win.MouseScroll().Y)
			zoomDeltaStart := viewMatrix.Unproject(win.MousePosition())
			scale *= factor
			cameraOrigin = zoomDeltaStart.Add(win.Bounds().Center().Sub(win.MousePosition().Scaled(1 / scale)))
		}
		if win.JustPressed(pixelgl.MouseButton1) {
			dragOrigin = win.MousePosition().Scaled(1 / scale)
		} else if win.Pressed(pixelgl.MouseButton1) {
			newOrigin := win.MousePosition().Scaled(1 / scale)
			cameraOrigin = cameraOrigin.Sub(newOrigin.Sub(dragOrigin))
			dragOrigin = newOrigin
		}
		viewMatrix = pixel.IM.Moved(win.Bounds().Center().Sub(cameraOrigin)).Scaled(pixel.ZV, scale)
		win.Clear(colornames.Gray)
		win.SetMatrix(viewMatrix)
		drawer.Draw(win)
		win.Update()
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
