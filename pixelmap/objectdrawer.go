package pixelmap

import (
	"github.com/elliotmr/tiled/tmx"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"github.com/pkg/errors"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"fmt"
	"strings"
	"strconv"
)

type ObjectDrawer struct {
	ts  *TileSets
	li  *LayerInfo
	im  []*imdraw.IMDraw
	txt []*text.Text
}

func NewObjectDrawer(li *LayerInfo, ts *TileSets) (*ObjectDrawer, error) {
	od := &ObjectDrawer{
		ts: ts,
		li: li,
		im: make([]*imdraw.IMDraw, 0),
	}
	return od, od.Update()
}

func getPosition(object *tmx.Object, li *LayerInfo) pixel.Vec {
	var v pixel.Vec
	if object.Height != nil && object.Width != nil {
		v = li.TMXToPixelRect(object.X, object.Y, *object.Width, *object.Height).Center()
	} else {
		v = li.TMXToPixelVec(object.X, object.Y)
	}
	return v
}

func getLine(points string, li *LayerInfo) ([]pixel.Vec, error) {
	ptVec := make([]pixel.Vec, 0)
	ptFields := strings.Fields(points)
	for _, ptField := range ptFields {
		pt := strings.Split(ptField, ",")
		x, err := strconv.ParseFloat(pt[0], 64)
		if err != nil {
			return nil, errors.Wrap(err, "invalid x-axis point")
		}
		y, err := strconv.ParseFloat(pt[1], 64)
		if err != nil {
			return nil, errors.Wrap(err, "invalid y-axis point")
		}
		ptVec = append(ptVec, pixel.V(x, -y))
	}
	return ptVec, nil
}

func (od *ObjectDrawer) createMatrix(object *tmx.Object) pixel.Matrix {
	v := getPosition(object, od.li)
	m := pixel.IM.Moved(v)
	if object.Rotation != nil {
		m = m.Rotated(v, *object.Rotation)
	}
	fmt.Println("Matrix: ", m)
	return m
}

func (od *ObjectDrawer) Update() error {
	// TODO: Template support
	od.im = od.im[:0]
	od.txt = od.txt[:0]
	for _, obj := range od.li.layer.Objects {
		if obj.Visible != nil && *obj.Visible == 0 {
			continue // skip invisible objects
		}
		imd := imdraw.New(nil)
		imd.Color = pixel.Alpha(0.5).Mul(pixel.ToRGBA(colornames.White))
		imd.SetMatrix(od.createMatrix(obj))
		// TODO: Something is up with rotation.
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
			l, err := getLine(obj.Polygon.Points, od.li)
			if err != nil {
				return errors.Wrap(err, "invalid polyline")
			}
			imd.Push(l...)
			imd.Polygon(0)
			od.im = append(od.im, imd)
			fmt.Println("Polygon: ", l)
		case obj.Polyline != nil:
			l, err := getLine(obj.Polyline.Points, od.li)
			if err != nil {
				return errors.Wrap(err, "invalid polyline")
			}
			imd.EndShape = imdraw.RoundEndShape
			imd.Push(l...)
			imd.Line(10)
			od.im = append(od.im, imd)
			fmt.Println("Polyline: ", l)
		case obj.Text != nil:
			// TODO: font, style handling
			imd.Push(pixel.V(0, 0))
			at := text.NewAtlas(basicfont.Face7x13, text.ASCII)
			txt := text.New(od.li.TMXToPixelRect(obj.X, obj.Y, 0, *obj.Height).Center(), at)
			fmt.Fprint(txt, obj.Text.Text)
			od.txt = append(od.txt, txt)
		default: // Box

		}
	}
	return nil
}

func (od *ObjectDrawer) Draw(t pixel.Target) {
	// TODO: Fix object/text draw order
	for _, obj := range od.im {
		obj.Draw(t)
	}
	for _, txt := range od.txt {
		txt.Draw(t, pixel.IM)
	}
}
