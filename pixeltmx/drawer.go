package pixeltmx

import (
	"github.com/elliotmr/tmx"
	"github.com/faiface/pixel"
	"github.com/pkg/errors"
)

// Drawer Types
const (
	TileLayerDrawer = iota
	ObjectGroupDrawer
	ImageLayerDrawer
	GroupDrawer
)

// Drawer is the base interface for the pixeltmx library, there are 4 concrete implementations
// for the 4 types of tmx layers (tile layer, object group, image layer, group layer). The
// underlying type can be extracted using `Type()` method. Each Layer will be updated once
// on creation and remain cached for subsequent draws. If the underlying data or resources have
// been changed, the `Update()` method must be called before the changes will be visible when
// drawing.
type Drawer interface {
	Type() int
	Info() *LayerInfo
	Update() error
	Draw(target pixel.Target)
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
		color:   pixel.Alpha(1.0),
	}
	gd := &groupDrawer{
		info:     info,
		children: make([]Drawer, 0),
	}

	for _, l := range mapData.Layers {
		d, err := NewDrawer(resources, gd, l)
		if err != nil {
			return nil, err
		}
		gd.children = append(gd.children, d)
	}
	return gd, nil
}
