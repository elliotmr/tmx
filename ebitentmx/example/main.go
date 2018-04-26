package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/elliotmr/tmx"
	"github.com/elliotmr/tmx/ebitentmx"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/pkg/errors"
)

func noError(err error) {
	if err != nil {
		panic(err)
	}
}

type Example struct {
	mapDrawer ebitentmx.Drawer
	cam       *ebitentmx.Camera
}

func (e *Example) Update(screen *ebiten.Image) error {
	if ebiten.IsRunningSlowly() {
		return nil
	}
	start := time.Now()
	e.cam.StartUpdate(start)
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		e.cam.Pan(3*math.Pi/2, 30.0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		e.cam.Pan(math.Pi/2, 30.0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		e.cam.Pan(0, 30.0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		e.cam.Pan(math.Pi, 30)
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("quitting")
	}
	e.cam.Draw(e.mapDrawer.Image(), screen)
	end := time.Now()
	msg := fmt.Sprintf("FPS: %0.2f, Pos: %v, DT: %0.3f ms, Draw: %0.3f ms", ebiten.CurrentFPS(), e.cam, e.cam.DT.Seconds() * 1000, end.Sub(start).Seconds() * 1000)
	ebitenutil.DebugPrint(screen, msg)
	return nil
}

func main() {
	mapFile, err := os.Open("orthogonal-outside.tmx")
	noError(err)
	mapData, err := tmx.Load(mapFile)
	noError(err)
	resources, err := ebitentmx.LoadResources(mapData, "")
	noError(err)
	tld, err := ebitentmx.NewRootDrawer(resources, mapData)
	noError(err)
	c := ebitentmx.NewCamera()
	e := &Example{cam: c, mapDrawer: tld}
	err = ebiten.Run(
		e.Update,
		int(mapData.Width*mapData.TileWidth),
		int(mapData.Height*mapData.TileHeight),
		1.0,
		"test",
	)

}
