package pixeltmx

import (
	"image"
	"os"

	"github.com/faiface/pixel"
	"github.com/pkg/errors"
)

type imageLayerDrawer struct {
	resources *Resources
	info      *LayerInfo
	sprite    *pixel.Sprite
}

func newImageLayerDriver(resources *Resources, info *LayerInfo) (*imageLayerDrawer, error) {
	if info.layer.Image == nil {
		return nil, errors.New("image not set for image layer")
	}
	id := &imageLayerDrawer{
		info: info,
	}
	return id, id.Update()
}

func (id *imageLayerDrawer) Type() int {
	return ImageLayerDrawer
}

func (id *imageLayerDrawer) Info() *LayerInfo {
	return id.info
}

func (id *imageLayerDrawer) Update() error {
	// TODO: Move to resources
	imageFile, err := os.Open(id.info.layer.Image.Source)
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

func (id *imageLayerDrawer) Draw(t pixel.Target) {
	vec := id.info.TMXToPixelRect(
		id.info.offX,
		id.info.offY,
		id.sprite.Frame().W(),
		id.sprite.Frame().H(),
	).Center()
	id.sprite.DrawColorMask(t, pixel.IM.Moved(vec), id.info.color)
}
