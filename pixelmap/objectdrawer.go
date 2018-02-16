package pixelmap

import (
	"github.com/elliotmr/tiled/tmx"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/pkg/errors"
	"golang.org/x/image/colornames"
)

type ObjectDrawer struct {
	ts *TileSets
	li *LayerInfo
	im []*imdraw.IMDraw
}

func NewObjectDrawer(li *LayerInfo, ts *TileSets) (*ObjectDrawer, error) {
	od := &ObjectDrawer{
		ts: ts,
		li: li,
		im: make([]*imdraw.IMDraw, 0),
	}
	return od, od.Update()
}

func (od *ObjectDrawer) createMatrix(object *tmx.Object) pixel.Matrix {
	var v pixel.Vec
	if object.Height != nil && object.Width != nil {
		v = od.li.TMXToPixelRect(object.X, object.Y, *object.Width, *object.Height).Center()
	} else {
		v = od.li.TMXToPixelVec(object.X, object.Y)
	}
	m := pixel.IM.Moved(v)
	if object.Rotation != nil {
		m = m.Rotated(v, *object.Rotation)
	}
	return m
}

func (od *ObjectDrawer) Update() error {
	// TODO: Template support
	od.im = od.im[:0]
	for _, obj := range od.li.layer.Objects {
		if obj.Visible != nil && *obj.Visible == 0 {
			continue // skip invisible objects
		}
		imd := imdraw.New(nil)
		imd.Color = pixel.Alpha(0.5).Mul(pixel.ToRGBA(colornames.Gray))
		imd.SetMatrix(od.createMatrix(obj))
		switch {
		case obj.GID != nil:
			// TODO
		case obj.Ellipse != nil:
			if obj.Width == nil || obj.Height == nil {
				return errors.New("ellipse without width or height set")
			}
			imd.Push(pixel.V(0, 0))
			imd.Ellipse(pixel.V(*obj.Width/2, *obj.Height/2), 0)
			od.im = append(od.im, imd)
		case obj.Point != nil:
			// TODO
		case obj.Polygon != nil:
			// TODO
		case obj.Polyline != nil:
			// TODO
		case obj.Text != nil:
			// TODO
		}
	}
	return nil
}

func (od *ObjectDrawer) Draw(t pixel.Target) {
	for _, obj := range od.im {
		obj.Draw(t)
	}
}
