package tmx

import (
	"encoding/xml"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	fp, err := os.Open("resources/arena.tmx")
	require.NoError(t, err)
	m, err := Load(fp)
	defer fp.Close()
	require.NoError(t, err)
	t.Log(m)

	out, err := xml.MarshalIndent(m, "", " ")
	assert.NoError(t, err)
	t.Log(string(out))

	layer := m.Layers[1]

	iter, err := layer.Data.Iter()
	assert.NoError(t, err)

	for iter.Next() {
		fmt.Printf("%02d ", iter.Get().GID)
		if iter.GetIndex()%*layer.Width == *layer.Width-1 {
			fmt.Println("")
		}
	}
	assert.NoError(t, iter.Error())
}

func TestLoadTSX(t *testing.T) {
	fp, err := os.Open("resources/cave.tmx")
	assert.NoError(t, err)
	defer fp.Close()
	m, err := Load(fp)
	assert.NoError(t, err)
	require.Equal(t, 1, len(m.TileSets))
	assert.NotEmpty(t, m.TileSets[0].Source)         // from tmx
	assert.EqualValues(t, 1, m.TileSets[0].FirstGID) // from tmx
	assert.Equal(t, "cave", m.TileSets[0].Name)      // from tsx
}
