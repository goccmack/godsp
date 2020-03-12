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

package godsp

import (
	"math"
	"sort"
)

const none = -1

type Peak struct {
	born, died, left, right int
}

func newPeak(startIdx int) *Peak {
	return &Peak{
		born:  startIdx,
		left:  startIdx,
		right: startIdx,
		died:  none,
	}
}

type Peaks struct {
	peaks []*Peak
	seq   []float64
}

func (p *Peak) getPersistence(seq []float64) float64 {
	if p.died == none {
		return math.Inf(1)
	}
	return seq[p.born] - seq[p.died]
}

func GetPeaksInt(seq []int) *Peaks {
	seq1 := ToFloat(seq)
	return GetPeaks(seq1)
}

/*
GetPeaks detects the peaks in a time series `seq` by means of persistent homology:
https://www.sthu.org/blog/13-perstopology-peakdetection/index.html.
The returned peaks are in increasing order of their indices in `seq`.
*/
func GetPeaks(seq []float64) *Peaks {
	peaks := make([]*Peak, 0, 1024)
	// Maps indices to peaks
	idxtopeak := make([]int, len(seq))
	for i := range idxtopeak {
		idxtopeak[i] = none
	}
	// Sequence indices sorted by values
	indices := Range(len(seq))
	sort.SliceStable(indices, func(i, j int) bool { return seq[indices[i]] > seq[indices[j]] })
	// Process each sample in descending order
	for _, idx := range indices {
		lftdone := (idx > 0 && idxtopeak[idx-1] != none)
		rgtdone := (idx < len(seq)-1 && idxtopeak[idx+1] != none)
		il := none
		if lftdone {
			il = idxtopeak[idx-1]
		}
		ir := none
		if rgtdone {
			ir = idxtopeak[idx+1]
		}

		// New peak born
		if !lftdone && !rgtdone {
			peaks = append(peaks, newPeak(idx))
			idxtopeak[idx] = len(peaks) - 1
		}

		// Directly merge to next peak left
		if lftdone && !rgtdone {
			peaks[il].right++
			idxtopeak[idx] = il
		}

		// Directly merge to next peak right
		if !lftdone && rgtdone {
			peaks[ir].left--
			idxtopeak[idx] = ir
		}

		// Merge left and right peaks
		if lftdone && rgtdone {
			// Left was born earlier: merge right to left
			if seq[peaks[il].born] > seq[peaks[ir].born] {
				peaks[ir].died = idx
				peaks[il].right = peaks[ir].right
				idxtopeak[peaks[il].right], idxtopeak[idx] = il, il
			} else {
				peaks[il].died = idx
				peaks[ir].left = peaks[il].left
				idxtopeak[peaks[ir].left], idxtopeak[idx] = ir, ir
			}
		}
	}

	sort.SliceStable(peaks, func(i, j int) bool { return peaks[i].born < peaks[j].born })

	return &Peaks{
		peaks: peaks,
		seq:   seq,
	}
}

/*
GetIndices returns the indices in the original time series `seq` of the peaks with
persistence/max(persitence of seq) >= `fracOfMaxPersistence`
*/
func (pks *Peaks) GetIndices(fracOfMaxPersistence float64) []int {
	indices := make([]int, 0, len(pks.peaks))
	_, maxPersistence := pks.MinMaxPersistence()
	for _, pk := range pks.peaks {
		if pk.getPersistence(pks.seq)/maxPersistence >= fracOfMaxPersistence {
			indices = append(indices, pk.born)
		}
	}
	return indices
}

/*
Max returns the index in the original time series `seq` of the peak with the
highest y-value. See GetIndices for fracOfMaxPersistence.
*/
func (pks *Peaks) Max(fracOfMaxPersistence float64) int {
	idxs := pks.GetIndices(fracOfMaxPersistence)
	biggest, val := -1, math.Inf(-1)
	for _, pk := range idxs {
		if pks.seq[pk] > val {
			biggest, val = pk, pks.seq[pk]
		}
	}
	return biggest
}

/*
MinMaxPersistence returns the minimum and maximum persistence of the peaks in `seq`.
*/
func (pks *Peaks) MinMaxPersistence() (min, max float64) {
	max, min = math.Inf(-1), math.Inf(1)
	for _, pk := range pks.peaks {
		prs := pk.getPersistence(pks.seq)
		if prs > max && prs < math.Inf(1) {
			max = prs
		}
		if prs < min {
			min = prs
		}
	}
	return
}
