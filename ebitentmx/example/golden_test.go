package main

import (
	"image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTileLayerOutput(t *testing.T) {
	goldenFP, err := os.Open("orthogonal-outside.png")
	require.NoError(t, err)
	golden, err := png.Decode(goldenFP)
	require.NoError(t, err)

	mineFP, err := os.Open("orthogonal-outside-out.png")
	require.NoError(t, err)
	mine, err := png.Decode(mineFP)
	require.NoError(t, err)

	bounds := golden.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gR, gB, gG, gA := golden.At(x, y).RGBA()
			mR, mB, mG, mA := mine.At(x, y).RGBA()
			assert.Equal(t, gR, mR, "red mismatch at (%d, %d)", x, y)
			assert.Equal(t, gG, mG, "green mismatch at (%d, %d)", x, y)
			assert.Equal(t, gB, mB, "blue mismatch at (%d, %d)", x, y)
			assert.Equal(t, gA, mA, "alpha mismatch at (%d, %d)", x, y)
		}
	}

}
