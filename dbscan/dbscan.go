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
Package dbscan implements the DBSCAN clustering algorithm
(https://en.wikipedia.org/wiki/DBSCAN)
*/
package dbscan

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/goccmack/goutil/ioutil"
)

const (
	noise     = -1
	undefined = 0
)

type Cluster struct {
	Min, Max int
}

/*
Histogram clusters the bins of a histogram `h`.
*/
func Histogram(h []int, eps, minPts int) []*Cluster {
	clusters := make([]int, len(h))
	C := 0 /* Cluster counter */
	for p := range h {
		if h[p] <= 0 {
			continue
		}
		if clusters[p] != undefined { /* Previously processed in inner loop */
			continue
		}
		N, S := getNeighbours(h, p, eps) /* Find neighbors */
		if len(N) < minPts {             /* Density check */
			clusters[p] = noise /* Label as noise */
			continue
		}
		C = C + 1             /* next cluster label */
		clusters[p] = C       /* Label initial point */
		for _, q := range S { /* Process every seed point */
			if clusters[q] == noise { /* Change noise to border point */
				clusters[q] = C
			}
			if clusters[q] != undefined { /* Previously processed */
				continue
			}
			clusters[q] = C                  /* Label neighbor */
			N, _ := getNeighbours(h, q, eps) /* Find neighbors */
			if len(N) >= minPts {            /* Density check */
				for _, n := range N { /* Add new neighbors to seed set */
					S = append(S, n)
				}
			}
		}
	}
	return getClusters(clusters)
}

func getClusters(cs []int) (clusters []*Cluster) {
	cmap := make(map[int]*Cluster)
	for i, c := range cs {
		if c > 0 {
			if cluster, exist := cmap[c]; exist {
				if i < cluster.Min {
					cluster.Min = i
				}
				if i > cluster.Max {
					cluster.Max = i
				}
			} else {
				cmap[c] = &Cluster{
					Min: i,
					Max: i,
				}
			}
		}
	}
	for _, c := range cmap {
		clusters = append(clusters, c)
	}
	sort.Slice(clusters,
		func(i, j int) bool { return clusters[i].Min < clusters[j].Min })
	return
}

/*
getNeighbours returns the set of neighbours of `point`, which is an index in `h`.
`neighbours` exclude `point`.
*/
func getNeighbours(h []int, point, eps int) (neighbours, nbMinPoint []int) {
	from := point - eps
	if from < 0 {
		from = 0
	}
	to := point + eps
	if to > len(h) {
		to = len(h)
	}
	for i := from; i < to; i++ {
		if h[i] > 0 {
			neighbours = append(neighbours, i)
			if i != point {
				nbMinPoint = append(nbMinPoint, i)
			}
		}
	}
	return
}

/*
WriteClusters writes the set of clusters `cs` to file `fname`.
*/
func WriteClusters(cs []*Cluster, fname string) {
	buf := new(bytes.Buffer)
	for i, c := range cs {
		fmt.Fprintf(buf, "%d, %d %d\n", i, c.Min, c.Max)
	}
	ioutil.WriteFile(fname, buf.Bytes())
}
