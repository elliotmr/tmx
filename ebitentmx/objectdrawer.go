package ebitentmx

import (
	"fmt"

	"github.com/elliotmr/tmx"
	"github.com/hajimehoshi/ebiten"
	"github.com/pkg/errors"
	"sort"
)

type objectGroupDrawer struct {
	resources *Resources
	info      *LayerInfo
	opts      *ebiten.DrawImageOptions
	image     *ebiten.Image
}

func newObjectGroupDrawer(resources *Resources, info *LayerInfo) (*objectGroupDrawer, error) {
	img, _ := ebiten.NewImage(
		info.w*int(info.mapData.TileWidth), // TODO: this is incorrect
		info.h*int(info.mapData.TileHeight), // TODO: this is incorrect
		ebiten.FilterNearest,
	)
	bounds := img.Bounds()
	geom := ebiten.GeoM{}
	geom.Translate(info.offX, info.offY)
	opts := &ebiten.DrawImageOptions{
		SourceRect: &bounds,
		GeoM:       geom,
		Filter:     ebiten.FilterNearest,
	}
	od := &objectGroupDrawer{
		resources: resources,
		info:      info,
		opts:      opts,
		image:     img,
	}
	return od, od.Update()
}

func (ogd *objectGroupDrawer) Type() int {
	return ObjectGroupDrawer
}

func (ogd *objectGroupDrawer) Info() *LayerInfo {
	return ogd.info
}

func (ogd *objectGroupDrawer) Image() *ebiten.Image {
	return ogd.image
}

func (ogd *objectGroupDrawer) Update() error {
	// TODO: Template support
	objs := ogd.info.layer.Objects
	objIndex := make([]int, len(objs))
	for i := range objIndex {
		objIndex[i] = i
	}

	if ogd.info.layer.DrawOrder == nil || *ogd.info.layer.DrawOrder != "index" {
		sort.Slice(objIndex, func(i, j int) bool {
			return objs[objIndex[i]].Y < objs[objIndex[j]].Y
		})
	}

	for _, i := range objIndex {
		obj := objs[i]
		if obj.Visible != nil && *obj.Visible == 0 {
			continue // skip invisible objects
		}
		switch {
		case obj.GID != nil:
			tile := tmx.TileInstance(*obj.GID)
			tse := ogd.resources.entries[tile.GID()]
			entry, exists := ogd.resources.entries[tile.GID()]
			if !exists {
				fmt.Println("invalid object: ", tile.GID())
				continue
			}
			pic, exists := ogd.resources.images[entry.source]
			if !exists {
				fmt.Println("invalid object: ", tile.GID())
				continue
			}

			opts := &ebiten.DrawImageOptions{
				SourceRect: tse.rect,
				ColorM:     ogd.info.color,
				GeoM:       calcObjectGeoM(obj, *tse.rect),
				Filter:     ebiten.FilterNearest,
			}

			ogd.image.DrawImage(pic, opts)
		case obj.Ellipse != nil:
			fmt.Println("ERROR ELLIPSE OBJECT NOT IMPLEMENTED!")
			continue
		case obj.Point != nil:
			fmt.Println("ERROR POINT OBJECT NOT IMPLEMENTED!")
			continue
		case obj.Polygon != nil:
			fmt.Println("ERROR POLYGON OBJECT NOT IMPLEMENTED!")
			continue
		case obj.Polyline != nil:
			fmt.Println("ERROR POLYLINE OBJECT NOT IMPLEMENTED!")
			continue
		case obj.Text != nil:
			fmt.Println("ERROR TEXT OBJECT NOT IMPLEMENTED!")
			continue
		default: // Box
			fmt.Println("ERROR BOX OBJECT NOT IMPLEMENTED!")
			continue
		}
	}
	return nil
}

func (ogd *objectGroupDrawer) Draw(image *ebiten.Image) error {
	return errors.Wrap(image.DrawImage(ogd.image, ogd.opts), "unable to draw object layer")
}

