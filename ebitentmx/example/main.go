package main

import (
	"os"

	"github.com/elliotmr/tmx"
	"github.com/elliotmr/tmx/ebitentmx"
	"github.com/hajimehoshi/ebiten"
	"sync"
	"image/png"
)

func noError(err error) {
	if err != nil {
		panic(err)
	}
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

	var once sync.Once

	testFunc := func (screen *ebiten.Image) error {
		tld.Draw(screen)

		once.Do(func() {
			outFile, err := os.Create("orthogonal-outside-out.png")
			noError(err)
			png.Encode(outFile, screen)
			outFile.Close()
		})

		return nil
	}

	err = ebiten.Run(
		testFunc,
		int(mapData.Width * mapData.TileWidth),
		int(mapData.Height *mapData.TileHeight),
		1.0,
		"test",
	)

}
