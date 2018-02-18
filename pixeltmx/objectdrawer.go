package pixeltmx

import (
	"fmt"
	"strings"
	"strconv"
	"math"

	"github.com/elliotmr/tmx"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"github.com/pkg/errors"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type ObjectDrawer struct {
	ts              *TileSets
	li              *LayerInfo
	currentFirstGID uint64
	batches         []*pixel.Batch
}

func NewObjectDrawer(li *LayerInfo, ts *TileSets) (*ObjectDrawer, error) {
	od := &ObjectDrawer{
		ts:          ts,
		li:          li,
		currentFirstGID: math.MaxUint64,
		batches: make([]*pixel.Batch, 0),
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
		m = m.Rotated(v, *object.Rotation * -math.Pi/180.0)
	}
	fmt.Println("Matrix: ", m)
	return m
}

func (od *ObjectDrawer) createIMD(obj *tmx.Object) *imdraw.IMDraw {
	imd := imdraw.New(nil)
	imd.Color = pixel.Alpha(0.5).Mul(pixel.ToRGBA(colornames.White))
	imd.SetMatrix(od.createMatrix(obj))
	if od.currentFirstGID != 0 {
		od.batches = append(od.batches, pixel.NewBatch(&pixel.TrianglesData{}, nil))
		od.currentFirstGID = 0
	}
	return imd
}

func (od *ObjectDrawer) Update() error {
	// TODO: Template support
	od.batches = od.batches[:0] // TODO: Persist batches?
	for _, obj := range od.li.layer.Objects {
		if obj.Visible != nil && *obj.Visible == 0 {
			continue // skip invisible objects
		}
		switch {
		case obj.GID != nil:

		case obj.Ellipse != nil:
			imd := od.createIMD(obj)
			if obj.Width == nil || obj.Height == nil {
				return errors.New("ellipse without width or height set")
			}
			imd.Push(pixel.V(0, 0))
			imd.Ellipse(pixel.V(*obj.Width/2, *obj.Height/2), 0)
			imd.Draw(od.batches[len(od.batches) - 1])
		case obj.Point != nil:
			// TODO
		case obj.Polygon != nil:
			imd := od.createIMD(obj)
			l, err := getLine(obj.Polygon.Points, od.li)
			if err != nil {
				return errors.Wrap(err, "invalid polyline")
			}
			imd.Push(l...)
			imd.Polygon(0)
			if od.currentFirstGID != 0 {
				od.batches = append(od.batches, pixel.NewBatch(&pixel.TrianglesData{}, nil))
			}
			imd.Draw(od.batches[len(od.batches) - 1])
		case obj.Polyline != nil:
			imd := od.createIMD(obj)
			l, err := getLine(obj.Polyline.Points, od.li)
			if err != nil {
				return errors.Wrap(err, "invalid polyline")
			}
			imd.EndShape = imdraw.RoundEndShape
			imd.Push(l...)
			imd.Line(10)
			imd.Draw(od.batches[len(od.batches) - 1])
		case obj.Text != nil:
			// TODO: font, style handling
			at := text.NewAtlas(basicfont.Face7x13, text.ASCII)
			txt := text.New(od.li.TMXToPixelRect(obj.X, obj.Y, 0, *obj.Height).Center(), at)
			fmt.Fprint(txt, obj.Text.Text)
			if od.currentFirstGID != math.MaxUint32 + 1 {
				od.batches = append(od.batches, pixel.NewBatch(&pixel.TrianglesData{}, at.Picture()))
				od.currentFirstGID = math.MaxUint32 + 1
			}
			txt.Draw(od.batches[len(od.batches) - 1], pixel.IM)
		default: // Box
			imd := od.createIMD(obj)
			if obj.Width == nil || obj.Height == nil {
				return errors.New("ellipse without width or height set")
			}
			imd.Push(pixel.V(-(*obj.Width/2), -(*obj.Height/2)), pixel.V(*obj.Width/2, *obj.Height/2))
			imd.Rectangle(0)
			imd.Draw(od.batches[len(od.batches) - 1])
		}
	}
	return nil
}

func (od *ObjectDrawer) Draw(t pixel.Target) {
	for _, batch := range od.batches {
		batch.Draw(t)
	}
}
