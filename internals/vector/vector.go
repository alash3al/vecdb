package vector

import (
	"math"
)

type Vec []float64

func (v Vec) DotProduct(b Vec) (result float64) {
	if len(v) < 1 {
		return 0
	}

	if len(v) != len(b) {
		return 0
	}

	for i := range b {
		result += v[i] * b[i]
	}

	return result
}

func (v Vec) Magnitude() float64 {
	var sum float64

	for _, n := range v {
		sum += n * n
	}

	return math.Abs(math.Sqrt(sum))
}

func (v Vec) CosineSimilarity(b Vec) float64 {
	magV := v.Magnitude()
	magB := b.Magnitude()

	if (magV == 0) || (magB == 0) {
		return 0
	}

	return v.DotProduct(b) / (magV * magB)
}
