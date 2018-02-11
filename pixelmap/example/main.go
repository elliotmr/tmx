package main

import (
	"os"

	"github.com/elliotmr/tiled/pixelmap"
	"github.com/elliotmr/tiled/tmx"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/pkg/profile"
	"golang.org/x/image/colornames"
	"time"
	"fmt"
)

func run() {
	mapFile, err := os.Open("arena.tmx")
	if err != nil {
		panic(err)
	}
	mapData, err := tmx.Load(mapFile)
	if err != nil {
		panic(err)
	}
	ts, err := pixelmap.LoadTileSets(mapData)
	if err != nil {
		panic(err)
	}

	ld, err := pixelmap.NewLayerDrawer(mapData, 1, ts)
	if err != nil {
		panic(err)
	}

	cfg := pixelgl.WindowConfig{
		Title:  "Tiled Map Example",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetSmooth(false)
	cameraOrigin := pixel.ZV
	dragOrigin := pixel.V(0, 0)
	second := time.Tick(time.Second)
	frames := 0
	for !win.Closed() {
		if win.JustPressed(pixelgl.MouseButton1) {
			dragOrigin = win.MousePosition()
		} else if win.Pressed(pixelgl.MouseButton1) {
			cameraOrigin = cameraOrigin.Sub(win.MousePosition().Sub(dragOrigin))
			dragOrigin = win.MousePosition()
		}
		cam := pixel.IM.Moved(win.Bounds().Center().Sub(cameraOrigin))
		win.SetMatrix(cam)

		win.Clear(colornames.Whitesmoke)
		//iter, err := baseLayer.Data.Iter()
		//if err != nil {
		//	panic(err)
		//}
		//for iter.Next() {
		//	cellSprite := pixel.NewSprite(tsr.GetCellFromGID(iter.Get().GID))
		//	vx := float64(iter.GetIndex()%mapData.Width) * float64(mapData.TileWidth)
			// we have to flip the row index because pixel and tmx have opposite y dim definitions
		//	vy := float64(mapData.Height-iter.GetIndex()/mapData.Width) * float64(mapData.TileHeight)
		//	cellSprite.Draw(win, pixel.IM.Moved(pixel.V(vx, vy)).Scaled(pixel.ZV, 2))
		//}
		//if iter.Error() != nil {
		//	panic(iter.Error())
		//}
		ld.Draw(win)
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
	defer profile.Start().Stop()
	pixelgl.Run(run)
}
