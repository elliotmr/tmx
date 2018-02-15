package pixelmap

import (
	"github.com/faiface/pixel"
	"image"
	"os"
	"github.com/pkg/errors"
)

type ImageDrawer struct {
	li      *LayerInfo
	sprite  *pixel.Sprite
	drawers map[uint32]*pixel.Drawer
}

func NewImageDrawer(li *LayerInfo) (*ImageDrawer, error) {
	if li.layer.Image == nil {
		return nil, errors.New("image not set for image layer")
	}
	id := &ImageDrawer{
		li:      li,
		drawers: make(map[uint32]*pixel.Drawer),
	}

	imageFile, err := os.Open(li.layer.Image.Source)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open image")
	}
	tilesetImg, _, err := image.Decode(imageFile)
	imageFile.Close()
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode image")
	}

	pic := pixel.PictureDataFromImage(tilesetImg)
	id.sprite = pixel.NewSprite(pic, pic.Bounds())
	return id, nil
}

func (id *ImageDrawer) Draw(t pixel.Target) error {
	vec := id.li.TMXToPixelRect(
		id.li.offX,
		id.li.offY,
		id.sprite.Frame().W(),
		id.sprite.Frame().H(),
	).Center()
	id.sprite.DrawColorMask(t, pixel.IM.Moved(vec), id.li.color)
	return nil
}
