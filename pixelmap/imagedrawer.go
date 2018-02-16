package pixelmap

import (
	"image"
	"os"

	"github.com/faiface/pixel"
	"github.com/pkg/errors"
)

type ImageDrawer struct {
	li     *LayerInfo
	sprite *pixel.Sprite
}

func NewImageDrawer(li *LayerInfo) (*ImageDrawer, error) {
	if li.layer.Image == nil {
		return nil, errors.New("image not set for image layer")
	}
	id := &ImageDrawer{
		li: li,
	}
	return id, id.Update()
}

func (id *ImageDrawer) Update() error {
	imageFile, err := os.Open(id.li.layer.Image.Source)
	if err != nil {
		return errors.Wrap(err, "unable to open image")
	}
	tilesetImg, _, err := image.Decode(imageFile)
	imageFile.Close()
	if err != nil {
		return errors.Wrap(err, "unable to decode image")
	}

	pic := pixel.PictureDataFromImage(tilesetImg)
	id.sprite = pixel.NewSprite(pic, pic.Bounds())
	return nil
}

func (id *ImageDrawer) Draw(t pixel.Target) {
	vec := id.li.TMXToPixelRect(
		id.li.offX,
		id.li.offY,
		id.sprite.Frame().W(),
		id.sprite.Frame().H(),
	).Center()
	id.sprite.DrawColorMask(t, pixel.IM.Moved(vec), id.li.color)
}
