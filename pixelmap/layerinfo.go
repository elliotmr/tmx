package pixelmap

import (
	"strings"
	"strconv"

	"github.com/elliotmr/tiled/tmx"
	"github.com/faiface/pixel"
	"github.com/pkg/errors"
)

type LayerInfo struct {
	mapData *tmx.Map
	layer   *tmx.Layer
	w       int
	h       int
	offX    float64
	offY    float64
	color   pixel.RGBA
}

func NewLayerInfo(mapData *tmx.Map, layerIndex int) (*LayerInfo, error) {
	if layerIndex >= len(mapData.Layers) {
		return nil, errors.Errorf("layer index out of range (%d >= %d)", layerIndex, len(mapData.Layers))
	}
	layer := mapData.Layers[layerIndex]

	li := &LayerInfo{
		mapData: mapData,
		layer:   layer,
	}

	li.w = int(mapData.Width)
	if layer.Width != nil {
		li.w = int(*layer.Width)
	}
	li.h = int(mapData.Height)
	if layer.Height != nil {
		li.h = int(*layer.Height)
	}
	li.offX = float64(mapData.TileWidth) * 0.5
	if layer.OffsetX != nil {
		li.offX  += *layer.OffsetX
	}
	li.offY = float64(mapData.TileHeight) * -0.5
	if layer.OffsetY != nil {
		li.offY -= *layer.OffsetY
	}
	var err error
	li.color, err = extractLayerColor(layer)
	return li, err
}

func extractLayerColor(layer *tmx.Layer) (pixel.RGBA, error) {
	opacity := 1.0
	if layer.Opacity != nil {
		opacity = *layer.Opacity
	}
	if layer.Visible != nil {
		if *layer.Visible == 0 {
			opacity = 0
		}
	}
	rgba := pixel.Alpha(opacity)

	if layer.Color != nil {
		s := strings.Trim(*layer.Color, "# ")
		var err error
		i := 0
		var a, r, b, g uint64
		switch len(s) {
		case 8:
			a, err = strconv.ParseUint(s[:2], 16, 8)
			if err != nil {
				return rgba, err
			}
			i = 2
			fallthrough
		case 6:
			r, err = strconv.ParseUint(s[i:i+2], 16, 8)
			if err != nil {
				return rgba, err
			}
			b, err = strconv.ParseUint(s[i+2:i+4], 16, 8)
			if err != nil {
				return rgba, err
			}
			g, err = strconv.ParseUint(s[i+4:i+6], 16, 8)
			if err != nil {
				return rgba, err
			}
		default:
			return rgba, errors.Errorf("invalid color: %s", s)
		}
		rgba.Mul(pixel.RGBA{
			float64(r) / 255.0,
			float64(g) / 255.0,
			float64(b) / 255.0,
			float64(a) / 255.0,
		})
	}
	return rgba, nil
}

func (li *LayerInfo) CellCoordinates(cell int) (float64, float64, error) {
	if cell > (li.w * li.h) {
		return 0, 0, errors.Errorf("cell out of range (%d > %d)", cell, li.w * li.h)
	}
	vx := float64(cell%li.w) * float64(li.mapData.TileWidth) + li.offX
	vy := float64(int(li.mapData.Height)-cell/li.h) * float64(li.mapData.TileHeight) + li.offY
	return vx, vy, nil
}