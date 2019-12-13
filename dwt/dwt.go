//  Copyright 2019 Marius Ackerman
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

/*
Package DWT has functions supporting the Discrete Wavelet Transform.
*/
package dwt

import (
	"github.com/goccmack/godsp"
	"math"
)

type Transform struct {
	st      []float64
	level   int
	padding int
}

// Daubechies4 returns the DWT with Daubechies 4 coeficients to level.
func Daubechies4(s []float64, level int) *Transform {
	t := &Transform{
		st:    make([]float64, GetFrameSize(s)),
		level: level,
	}
	diff := len(t.st) - len(s)
	copy(t.st[diff:], s)
	max := len(s)
	for l := level; l > 0; l-- {
		split(t.st[:max])
		daubechies4(t.st[:max])
		max /= 2
	}

	return t
}

/*
GetFrameSize returns the size of DWT frame required for
*/
func GetFrameSize(s []float64) int {
	logLen := math.Log2(float64(len(s)))
	logLenInt := int(math.Ceil(logLen))

	return godsp.Pow2(logLenInt)
}

/*
Split s into even and odd elements,
where the even elements are in the first half
of the vector and the odd elements are in the
second half.
*/
func split(s []float64) {
	half := len(s) / 2
	odd := make([]float64, half)
	for i := 1; i < len(s); i += 2 {
		odd[i/2] = s[i]
	}
	for i := 2; i < len(s); i += 2 {
		s[i/2] = s[i]
	}
	for i, v := range odd {
		s[half+i] = v
	}
}

/*
After: Riples section 3.4
*/
func daubechies4(s []float64) {
	half := len(s) / 2

	// Update 1:
	for n := 0; n < half; n++ {
		s[n] = s[n] + math.Sqrt(3)*s[half+n]
	}

	// Predict:
	s[half] = s[half] -
		(math.Sqrt(3)/4)*s[0] -
		((math.Sqrt(3)-2)/4)*s[half-1]
	for n := 1; n < half; n++ {
		s[half+n] = s[half+n] -
			(math.Sqrt(3)/4)*s[n] -
			((math.Sqrt(3)-2)/4)*s[n-1]
	}

	// Update 2:
	for n := 0; n < half-1; n++ {
		s[n] = s[n] - s[half+n+1]
	}
	s[half-1] = s[half-1] - s[half]

	// Normalise:
	for n := 0; n < half; n++ {
		s[n] = ((math.Sqrt(3) - 1) / math.Sqrt(2)) * s[n]
		s[n+half] = ((math.Sqrt(3) + 1) / math.Sqrt(2)) * s[n+half]
	}
}

// GetCoefficients returns the coefficients of all transform levels
func (t *Transform) GetCoefficients() [][]float64 {
	cfs := make([][]float64, t.level)
	half := len(t.st) / 2
	for l := 1; l <= t.level; l++ {
		cfs[l-1] = t.st[half : 2*half]
		half /= 2
	}
	return cfs
}

// GetDownSampledCoefficients returns the coefficients of all the levels downsampled to
// the length of the deepest level of the transform.
func (t *Transform) GetDownSampledCoefficients() [][]float64 {
	cfs := t.GetCoefficients()
	minN := len(cfs[len(cfs)-1])
	for i, cf := range cfs {
		if i != len(cfs)-1 {
			ds := len(cf) / minN
			cfs[i] = godsp.DownSample(cf, ds)
		}
	}
	return cfs
}

/*
GetDecomposition returns the vector containing the DWT decomposion
*/
func (t *Transform) GetDecomposion() []float64 {
	return t.st
}
