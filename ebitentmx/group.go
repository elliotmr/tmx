package ebitentmx

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/pkg/errors"
	"image"
)

type groupDrawer struct {
	resources *Resources
	info     *LayerInfo
	children []Drawer
	opts      *ebiten.DrawImageOptions
	image     *ebiten.Image
}

func newGroupDrawer(resources *Resources, info *LayerInfo) (*groupDrawer, error) {
	gd := &groupDrawer{
		info:     info,
		children: make([]Drawer, 0),
	}

	for _, l := range gd.info.layer.Layers {
		d, err := NewDrawer(resources, gd, l)
		if err != nil {
			return nil, err
		}
		gd.children = append(gd.children, d)
	}

	var min, max image.Point
	if len(gd.children) > 0 {
		min = gd.children[0].Image().Bounds().Min
		max = gd.children[0].Image().Bounds().Max
	}

	for _, child := range gd.children {
		bounds := child.Image().Bounds()
		if bounds.Min.X < min.X {
			min.X = bounds.Min.X
		}
		if bounds.Min.Y < min.Y {
			min.Y = bounds.Min.Y
		}
		if bounds.Max.X > max.X {
			max.X = bounds.Max.X
		}
		if bounds.Max.Y > max.Y {
			max.Y = bounds.Max.Y
		}
	}

	bounds := image.Rectangle{Min:min, Max:max}
	gd.image, _ = ebiten.NewImage(bounds.Dx(), bounds.Dy(), ebiten.FilterNearest)
	geom := ebiten.GeoM{}
	geom.Translate(info.offX, info.offY)
	gd.opts = &ebiten.DrawImageOptions{
		SourceRect: &bounds,
		GeoM:       geom,
		Filter:     ebiten.FilterNearest,
	}

	return gd, nil
}

func (gd *groupDrawer) Type() int {
	return GroupDrawer
}

func (gd *groupDrawer) Info() *LayerInfo {
	return gd.info
}

func (gd *groupDrawer) Image() *ebiten.Image {
	for _, child := range gd.children {
		err := child.Draw(gd.image)
		if err != nil {
			return nil
		}
	}
	return gd.image
}

func (gd *groupDrawer) Update() error {
	for _, child := range gd.children {
		err := child.Update()
		if err != nil {
			return err
		}
	}
	return nil
}

func (gd *groupDrawer) Draw(image *ebiten.Image) error {
	for _, child := range gd.children {
		err := child.Draw(gd.image)
		if err != nil {
			return err
		}
	}
	return errors.Wrap(image.DrawImage(gd.image, gd.opts), "unable to draw layer")
}

