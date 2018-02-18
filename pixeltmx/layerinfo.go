package pixeltmx

import (
	"strconv"
	"strings"

	"github.com/elliotmr/tmx"
	"github.com/faiface/pixel"
	"github.com/pkg/errors"
)

// LayerInfo provides drawing information for the layer, it holds the
// recursively calculated offset, visibility, and color information, as
// well as a reference to the base map data. It prevides easy methods
// translating between tmx and pixel world coordinates.
type LayerInfo struct {
	mapData *tmx.Map
	layer   *tmx.Layer
	w       int
	h       int
	offX    float64
	offY    float64
	color   pixel.RGBA
}

func newLayerInfo(parent *LayerInfo, layer *tmx.Layer) (*LayerInfo, error) {
	li := &LayerInfo{}
	*li = *parent
	li.layer = layer

	offX, offY := extractLayerOffsets(layer)
	li.offX += offX
	li.offY += offY
	c, err := extractLayerColor(layer)
	if err != nil {
		return nil, errors.Wrap(err, "unable to extract color")
	}
	li.color = li.color.Mul(c)
	if layer.Width != nil {
		li.w = int(*layer.Width)
	}
	if layer.Height != nil {
		li.h = int(*layer.Height)
	}
	return li, nil
}

func extractLayerOffsets(layer *tmx.Layer) (float64, float64) {
	if layer == nil {
		return 0.0, 0.0
	}
	offX := 0.0
	if layer.OffsetX != nil {
		offX = *layer.OffsetX
	}
	offY := 0.0
	if layer.OffsetY != nil {
		offY = *layer.OffsetY
	}
	return offX, offY
}

func extractLayerColor(layer *tmx.Layer) (pixel.RGBA, error) {
	if layer == nil {
		return pixel.Alpha(1.0), errors.New("nil layer passed")
	}
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
			R: float64(r) / 255.0,
			G: float64(g) / 255.0,
			B: float64(b) / 255.0,
			A: float64(a) / 255.0,
		})
	}
	return rgba, nil
}

// TileRect returns the pixel.Rect of a TMX map tile in pixel world coordinates.
func (li *LayerInfo) TileRect(cell int) (pixel.Rect, error) {
	if cell > (li.w * li.h) {
		return pixel.R(0, 0, 0, 0), errors.Errorf("cell out of range (%d > %d)", cell, li.w*li.h)
	}
	tw := float64(li.mapData.TileWidth)
	th := float64(li.mapData.TileHeight)
	return li.TMXToPixelRect(
		float64(cell%li.w)*tw,
		float64(cell/li.h)*th,
		tw,
		th,
	), nil
}

// TMXToPixelVec translates TMX x and y coordinates to a pixel.Vect in pixel
// world coordinates.
func (li *LayerInfo) TMXToPixelVec(x, y float64) pixel.Vec {
	return pixel.V(x, float64(li.mapData.TileHeight*li.mapData.Height)-y)
}

// TMXToPixelRect translates a TMX four-tuple (x, y, w, h) to a pixel.Rect in
// pixel world coordinates.
func (li *LayerInfo) TMXToPixelRect(x, y, w, h float64) pixel.Rect {
	bottomLeft := li.TMXToPixelVec(x, y+h)
	topRight := li.TMXToPixelVec(x+w, y)
	return pixel.R(bottomLeft.X, bottomLeft.Y, topRight.X, topRight.Y)
}
