package pixeltmx

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/elliotmr/tmx"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"github.com/pkg/errors"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type objectGroupDrawer struct {
	resources       *Resources
	info            *LayerInfo
	currentFirstGID uint64
	batches         []*pixel.Batch
}

func newObjectGroupDrawer(resources *Resources, info *LayerInfo) (*objectGroupDrawer, error) {
	od := &objectGroupDrawer{
		resources:       resources,
		info:            info,
		currentFirstGID: math.MaxUint64,
		batches:         make([]*pixel.Batch, 0),
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

func (ogd *objectGroupDrawer) createMatrix(object *tmx.Object) pixel.Matrix {
	v := getPosition(object, ogd.info)
	m := pixel.IM.Moved(v)
	topLeft := v
	if object.Width != nil && object.Height != nil {
		topLeft = v.Sub(pixel.V(*object.Width/2, -*object.Height/2))
	}
	if object.Rotation != nil {
		m = m.Rotated(topLeft, *object.Rotation*-math.Pi/180.0)
	}
	return m
}

func flipDiagonal(around pixel.Vec, m pixel.Matrix) pixel.Matrix {
	return pixel.Matrix{m[1], -m[0], m[3], -m[2], m[5] - around.Y + around.X, -m[4] + around.X + around.Y}
}

func (ogd *objectGroupDrawer) createMatrixTile(tile tmx.TileInstance, frame pixel.Rect, object *tmx.Object) pixel.Matrix {
	// Get Initial Position
	v := getPosition(object, ogd.info).Add(pixel.V(0.0, *object.Height))
	m := pixel.IM.Moved(v)

	// Rotate 90 deg around center for diagonal flip
	if tile.FlippedDiagonally() {
		m = flipDiagonal(v, m)
	}

	xScale := 1.0
	if object.Width != nil {
		xScale = *object.Width / frame.W()
	}
	if tile.FlippedHorizontally() {
		xScale = -xScale
	}
	yScale := 1.0
	if object.Height != nil {
		yScale = *object.Height / frame.H()
	}
	if tile.FlippedVertically() {
		yScale = -yScale
	}

	m = m.ScaledXY(v, pixel.V(xScale, yScale))

	bottomLeft := v.Sub(pixel.V(*object.Width/2, *object.Height/2))
	if object.Rotation != nil {
		m = m.Rotated(bottomLeft, *object.Rotation*-math.Pi/180.0)
	}
	return m
}

func (ogd *objectGroupDrawer) createIMD(obj *tmx.Object) *imdraw.IMDraw {
	imd := imdraw.New(nil)
	imd.Color = pixel.Alpha(0.5).Mul(pixel.ToRGBA(colornames.White))
	imd.SetMatrix(ogd.createMatrix(obj))
	if ogd.currentFirstGID != 0 {
		ogd.batches = append(ogd.batches, pixel.NewBatch(&pixel.TrianglesData{}, nil))
		ogd.currentFirstGID = 0
	}
	return imd
}

func (ogd *objectGroupDrawer) Type() int {
	return ObjectGroupDrawer
}

func (ogd *objectGroupDrawer) Info() *LayerInfo {
	return ogd.info
}

func (ogd *objectGroupDrawer) Update() error {
	// TODO: Template support
	ogd.batches = ogd.batches[:0] // TODO: Persist batches?
	for _, obj := range ogd.info.layer.Objects {
		if obj.Visible != nil && *obj.Visible == 0 {
			continue // skip invisible objects
		}
		switch {
		case obj.GID != nil:
			tile := tmx.TileInstance(*obj.GID)
			entry, exists := ogd.resources.entries[tile.GID()]
			if !exists {
				fmt.Println("invalid object: ", tile.GID())
				continue
			}
			pic := ogd.resources.images[entry.source]
			if ogd.currentFirstGID != uint64(entry.firstGID) {
				ogd.batches = append(ogd.batches, pixel.NewBatch(&pixel.TrianglesData{}, pic))
				ogd.currentFirstGID = uint64(entry.firstGID)
			}
			sprite := pixel.NewSprite(pic, entry.frame)
			if obj.Width == nil || obj.Height == nil {
				return errors.New("tile object without width or height set")
			}
			m := ogd.createMatrixTile(tile, entry.frame, obj)
			sprite.Draw(ogd.batches[len(ogd.batches)-1], m)
		case obj.Ellipse != nil:
			imd := ogd.createIMD(obj)
			if obj.Width == nil || obj.Height == nil {
				return errors.New("ellipse without width or height set")
			}
			imd.Push(pixel.V(0, 0))
			imd.Ellipse(pixel.V(*obj.Width/2, *obj.Height/2), 0)
			imd.Draw(ogd.batches[len(ogd.batches)-1])
		case obj.Point != nil:
			// TODO
		case obj.Polygon != nil:
			imd := ogd.createIMD(obj)
			l, err := getLine(obj.Polygon.Points, ogd.info)
			if err != nil {
				return errors.Wrap(err, "invalid polyline")
			}
			imd.Push(l...)
			imd.Polygon(0)
			// BUG(elliotmr): something strange is happening with polygon rendering.
			if ogd.currentFirstGID != 0 {
				ogd.batches = append(ogd.batches, pixel.NewBatch(&pixel.TrianglesData{}, nil))
			}
			imd.Draw(ogd.batches[len(ogd.batches)-1])
		case obj.Polyline != nil:
			imd := ogd.createIMD(obj)
			l, err := getLine(obj.Polyline.Points, ogd.info)
			if err != nil {
				return errors.Wrap(err, "invalid polyline")
			}
			imd.EndShape = imdraw.RoundEndShape
			imd.Push(l...)
			imd.Line(10)
			imd.Draw(ogd.batches[len(ogd.batches)-1])
		case obj.Text != nil:
			// TODO: font, style handling
			at := text.NewAtlas(basicfont.Face7x13, text.ASCII)
			txt := text.New(ogd.info.TMXToPixelRect(obj.X, obj.Y, 0, *obj.Height).Center(), at)
			fmt.Fprint(txt, obj.Text.Text)
			if ogd.currentFirstGID != math.MaxUint32+1 {
				ogd.batches = append(ogd.batches, pixel.NewBatch(&pixel.TrianglesData{}, at.Picture()))
				ogd.currentFirstGID = math.MaxUint32 + 1
			}
			txt.Draw(ogd.batches[len(ogd.batches)-1], pixel.IM)
		default: // Box
			imd := ogd.createIMD(obj)
			if obj.Width == nil || obj.Height == nil {
				return errors.New("ellipse without width or height set")
			}
			imd.Push(pixel.V(-(*obj.Width/2), -(*obj.Height/2)), pixel.V(*obj.Width/2, *obj.Height/2))
			imd.Rectangle(0)
			imd.Draw(ogd.batches[len(ogd.batches)-1])
		}
	}
	return nil
}

func (ogd *objectGroupDrawer) Draw(t pixel.Target) {
	for _, batch := range ogd.batches {
		batch.Draw(t)
	}
}
