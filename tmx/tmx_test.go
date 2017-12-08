package tmx

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"encoding/xml"
	"fmt"
)

func TestLoad(t *testing.T) {
	fp, err := os.Open("arena.tmx")
	assert.NoError(t, err)
	m, err := Load(fp)
	assert.NoError(t, err)
	t.Log(m)

	out, err := xml.MarshalIndent(m, "", " ")
	assert.NoError(t, err)
	t.Log(string(out))

	layer := m.Layers[1]

	iter, err := layer.Data.Iter()
	assert.NoError(t, err)


	for iter.Next() {
		fmt.Printf("%02d ", iter.Get().GID)
		if iter.GetIndex() % *layer.Width == *layer.Width - 1 {
			fmt.Println("")
		}
	}
	assert.NoError(t, iter.Error())
}