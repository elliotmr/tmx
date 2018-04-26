package ebitentmx

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/pkg/errors"
	"github.com/elliotmr/tmx"
)

// Drawer Types
const (
	TileLayerDrawer = iota
	ObjectGroupDrawer
	ImageLayerDrawer
	GroupDrawer
)

type Drawer interface {
	Type() int
	Info() *LayerInfo
	Update() error
	Draw(image *ebiten.Image) error
	Image() *ebiten.Image
}


// NewDrawer creates a Drawer which will render the layer and recursively draw all child layers.
func NewDrawer(resources *Resources, parent Drawer, layer *tmx.Layer) (Drawer, error) {
	info, err := newLayerInfo(parent.Info(), layer)
	if err != nil {
		panic(err)
	}
	switch info.layer.XMLName.Local {
	case "layer":
		return newTileLayerDrawer(resources, info)
	case "objectgroup":
		return newObjectGroupDrawer(resources, info)
	case "imagelayer":
		return newImageLayerDriver(resources, info)
	case "group":
		return newGroupDrawer(resources, info)
	}
	return nil, errors.Errorf("invalid layer type: %s", info.layer.XMLName.Local)
}

// NewRootDrawer will create a special Drawer that will recursively draw the entire tmx map.
func NewRootDrawer(resources *Resources, mapData *tmx.Map) (Drawer, error) {
	info := &LayerInfo{
		mapData: mapData,
		layer:   nil,
		w:       int(mapData.Width),
		h:       int(mapData.Height),
		offX:    0.0,
		offY:    0.0,
		color:   ebiten.ColorM{},
	}


	gd := &groupDrawer{
		info:     info,
		children: make([]Drawer, 0),
	}
	gd.image, _ = ebiten.NewImage(
		int(mapData.Width * mapData.TileWidth),
		int(mapData.Height * mapData.TileHeight),
		ebiten.FilterNearest,
	)

	for _, l := range mapData.Layers {
		d, err := NewDrawer(resources, gd, l)
		if err != nil {
			return nil, errors.Wrap(err, "unable to create layer")
		}
		gd.children = append(gd.children, d)
	}
	return gd, nil
}
