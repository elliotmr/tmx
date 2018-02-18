package pixeltmx

import (
	"github.com/faiface/pixel"
)

type groupDrawer struct {
	info     *LayerInfo
	children []Drawer
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
	return gd, nil
}

func (gd *groupDrawer) Type() int {
	return GroupDrawer
}

func (gd *groupDrawer) Info() *LayerInfo {
	return gd.info
}

func (gd *groupDrawer) Draw(target pixel.Target) {
	for _, child := range gd.children {
		child.Draw(target)
	}
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
