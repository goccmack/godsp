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
Package peaks finds the maxima in a vector. It works by lowering a horizontal line
across the signal, revealing peaks as it proceeds. Peaks that are closer to
each other than a minimum separation distance are merged to the left (lower index).
*/
package peaks

import (
	"math"
	"sort"

	"github.com/goccmack/godsp"
)

const (
	empty = -1
)

/*
Get returns a slice containing the indices of the peaks in x.
sep is the minimum distance between 2 peaks. Peaks closer to each other than
sep are merged to the lower index.
*/
func Get(x []float64, sep int) []int {
	pks := []int{}
	for i := range x {
		if isMax(i, i-sep, i+sep, x) {
			pks = append(pks, i)
		}
	}
	return pks
}

func getMaxIndex(x []float64) int {
	i, max := 0, math.Inf(-1)
	for j, y := range x {
		if y > max {
			i, max = j, y
		}
	}
	if max > 0 {
		return i
	}
	return -1
}

func isMax(i, min, max int, x []float64) bool {
	if min < 0 {
		min = 0
	}
	if max > len(x) {
		max = len(x)
	}
	for j := min; j < i; j++ {
		if x[j] >= x[i] {
			return false
		}
	}
	for j := i + 1; j < max; j++ {
		if x[j] > x[i] {
			return false
		}
	}
	return true
}

func getWindow(i, sep int, x []float64) (min, max int) {
	min, max = i-sep, i+sep
	if min < 0 {
		min = 0
	}
	if max > len(x) {
		max = len(x)
	}
	return
}

// func Get(x []float64, sep int) []int {
// 	si := getSortedIndices(x)
// 	pks := getEmptyPeaks(len(x))
// 	for _, xi := range si {
// 		if pks[xi] == empty {
// 			markNeighbours(xi, sep, pks)
// 		}
// 	}
// 	uniquePeaks := make([]int, 0, len(x)/(2*sep))
// 	for i, xi := range pks {
// 		if i == xi {
// 			uniquePeaks = append(uniquePeaks, xi)
// 		}
// 	}
// 	return uniquePeaks
// }

func getEmptyPeaks(n int) []int {
	epks := make([]int, n)
	for i := range epks {
		epks[i] = empty
	}
	return epks
}

func getSortedIndices(x []float64) []int {
	idx := godsp.Range(len(x))
	sort.SliceStable(idx, func(i, j int) bool { return x[i] > x[j] })
	return idx
}

func markNeighbours(xi, sep int, pks []int) {
	min := xi - sep
	if min < 0 {
		min = 0
	}
	max := xi + sep
	if max > len(pks) {
		max = len(pks)
	}
	for i := min; i < max; i++ {
		pks[i] = xi
	}
}
