package ebitentmx

import (
	"testing"
	"math"
)

func BenchmarkSwap(b *testing.B) {
	var sin, cos float64
	for i := 0; i < b.N; i++ {
		sin, cos = 1.0, 0.0
		sin, cos = cos, sin
	}
}

func BenchmarkSinCosZero(b *testing.B) {
	var sin, cos float64
	for i := 0; i < b.N; i++ {
		target := math.Pi / 2
		switch target {
		case 0:
			sin, cos = 0.0, 1.0
		case math.Pi / 2:
			sin, cos = 1.0, 0.0
		case math.Pi:
			sin, cos = 0.0, -1.0
		case 3 * math.Pi / 2:
			sin, cos = -1.0, 0.0
		default:
			sin, cos = math.Sincos(target)
		}
		sin, cos = cos, sin
	}
}


func BenchmarkSinCos45(b *testing.B) {
	for i := 0; i < b.N; i++ {

		sin, cos := math.Sincos(math.Pi / 4)
		sin, cos = cos, sin
	}
}

