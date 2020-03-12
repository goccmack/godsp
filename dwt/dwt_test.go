package dwt

import (
	"math"
	"testing"
)

func Test1(t *testing.T) {
	N := 1_315_840
	sections := getTransformSections(N, 4)
	sum := 0
	for _, s := range sections {
		sum += s.size
	}
	if sum != N {
		t.Errorf("Sum = %d, difference=%f", sum, math.Abs(float64(sum-N)))
	}
}
