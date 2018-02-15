package pixelmap

import "github.com/faiface/pixel"

type Drawer interface {
	Draw(target pixel.Target) error
}
