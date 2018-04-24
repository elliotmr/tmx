package ebitentmx

import (
	"image"
	"math"
	"github.com/hajimehoshi/ebiten"
	"github.com/elliotmr/tmx"
)

func calcGeoM(tile tmx.TileInstance, rect image.Rectangle) ebiten.GeoM {
	geom := ebiten.GeoM{}
	if tile.FlippedDiagonally() {
		geom.Rotate(math.Pi)
		geom.Translate(float64(rect.Dx()), float64(rect.Dy()))
	}
	if tile.FlippedHorizontally() {
		geom.Scale(-1.0, 1.0)
		geom.Translate(float64(rect.Dx()), 0)
	}
	if tile.FlippedVertically() {
		geom.Scale(1.0, -1.0)
		geom.Translate(0, float64(rect.Dy()))
	}
	geom.Translate(
		float64(rect.Min.X),
		float64(rect.Min.Y),
	)
	return geom
}

func calcObjectGeoM(obj *tmx.Object, srcRect image.Rectangle) ebiten.GeoM {
	geom := ebiten.GeoM{}
	tile := tmx.TileInstance(*obj.GID)
	if tile.FlippedDiagonally() {
		geom.Rotate(math.Pi)
		geom.Translate(float64(srcRect.Dx()), float64(srcRect.Dy()))
	}
	if tile.FlippedHorizontally() {
		geom.Scale(-1.0, 1.0)
		geom.Translate(float64(srcRect.Dx()), 0)
	}
	if tile.FlippedVertically() {
		geom.Scale(1.0, -1.0)
		geom.Translate(0, float64(srcRect.Dy()))
	}
	scaleX := 1.0
	if obj.Width != nil {
		scaleX = *obj.Width / float64(srcRect.Dx())
	}
	scaleY := 1.0
	if obj.Height != nil {
		scaleY = *obj.Height / float64(srcRect.Dy())
	}
	geom.Scale(scaleX, scaleY)
	if obj.Rotation != nil {
		geom.Rotate(*obj.Rotation * math.Pi / 180.0)
	}
	geom.Translate(obj.X, obj.Y - float64(srcRect.Dy()) * scaleY)
	return geom
}
