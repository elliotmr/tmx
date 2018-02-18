package pixeltmx

import (
	"github.com/faiface/pixel"
	"github.com/pkg/errors"
)

type Drawer interface {
	Draw(target pixel.Target)
	Update() error
}

func NewDrawer(li *LayerInfo, ts *TileSets) (Drawer, error) {
	switch li.layer.XMLName.Local {
	case "layer":
		return NewTileDrawer(li, ts)
	case "objectgroup":
		return NewObjectDrawer(li, ts)
	case "imagelayer":
		return NewImageDrawer(li)
	case "group":
		// return NewGroupDrawer(li, ts)
	}
	return nil, errors.Errorf("invalid layer type: %s", li.layer.XMLName.Local)
}
