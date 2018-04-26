package ebitentmx

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten"
	"fmt"
)

func NewCamera() *Camera {
	return &Camera{
		updated: time.Now(),
		opts: &ebiten.DrawImageOptions{
			GeoM: ebiten.GeoM{},
			ColorM: ebiten.ColorM{},
			Filter: ebiten.FilterNearest,
		},
	}
}

type Camera struct {
	updated time.Time
	DT      time.Duration

	opts       *ebiten.DrawImageOptions
}

func (c *Camera) StartUpdate(now time.Time) {
	c.DT = now.Sub(c.updated)
}

func (c *Camera) Draw(src, dest *ebiten.Image) error {
	c.updated = c.updated.Add(c.DT)
	return dest.DrawImage(src, c.opts)
}

func (c *Camera) Pan(dir float64, rate float64) {
	var sin, cos float64
	switch dir {
	case 0:
		sin, cos = 0.0, 1.0
	case math.Pi / 2:
		sin, cos = 1.0, 0.0
	case math.Pi:
		sin, cos = 0.0, -1.0
	case 3 * math.Pi / 2:
		sin, cos = -1.0, 0.0
	default:
		sin, cos = math.Sincos(dir)
	}
	c.opts.GeoM.Translate(sin * rate * 0.001 / c.DT.Seconds(), cos * rate * 0.001 / c.DT.Seconds())
}

func (c *Camera) String() string {
	return fmt.Sprintf(
		"[x: %0.2f, y: %0.2f]",
		c.opts.GeoM.Element(0, 2),
		c.opts.GeoM.Element(1, 2),
	)
}