/*
Copyright 2019 Marius Ackerman
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Package dsp has a set of digital signal processing functions that are primarily
designed to support the discrete wavelet transform
("https://github.com/goccmack/dsp/dwt")
*/
package godsp

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"

	myioutil "github.com/goccmack/goutil/ioutil"
)

// Abs returns |x|
func Abs(x []float64) []float64 {
	x1 := make([]float64, len(x))
	for i, f := range x {
		x1[i] = math.Abs(f)
	}
	return x1
}

// AbsInt returns |x|
func AbsInt(x []int) []int {
	x1 := make([]int, len(x))
	for i, e := range x {
		if e < 0 {
			x1[i] = -e
		} else {
			x1[i] = e
		}
	}
	return x1
}

// AbsAll returns Abs(x) for every x in X
func AbsAll(X [][]float64) [][]float64 {
	x1 := make([][]float64, len(X))
	for i, x := range X {
		x1[i] = Abs(x)
	}
	return x1
}

/*
Average returns Sum(x)/len(x).
*/
func Average(x []float64) float64 {
	return Sum(x) / float64(len(x))
}

/*
DivS returns x/s where x is a vector and s a scalar.
*/
func DivS(x []float64, s float64) []float64 {
	y := make([]float64, len(x))
	for i := range x {
		y[i] = x[i] / s
	}
	return y
}

/*
DownSampleAll returns DownSample(x, len(x)/min(len(xs))) for all x in xs
*/
func DownSampleAll(xs [][]float64) [][]float64 {
	N := len(xs[0])
	for _, x := range xs {
		if len(x) < N {
			N = len(x)
		}
	}
	ys := make([][]float64, len(xs))
	for i, x := range xs {
		ys[i] = DownSample(x, len(x)/N)
	}
	return ys
}

/*
DownSample returns x downsampled by n
Function panics if len(x) is not an integer multiple of n.
*/
func DownSample(x []float64, n int) []float64 {
	if len(x)%n != 0 {
		panic(fmt.Sprintf("len(x) (%d) is not an integer multiple of n (%d)", len(x), n))
	}

	x1 := make([]float64, len(x)/n)
	for i, j := 0, 0; j < len(x1); i, j = i+n, j+1 {
		x1[j] = x[i]
	}
	return x1
}

// FindMax returns the value and index of the first element of x equal to the maximum value in x.
func FindMax(x []float64) (value float64, index int) {
	value, index = x[0], 0
	for i := 1; i < len(x)-1; i++ {
		if x[i] > value {
			value, index = x[i], i
		}
	}
	return
}

// FindMax* returns the value and index of the first element of x equal to the maximum value in x.
func FindMaxI(x []int) (value int, index int) {
	value, index = x[0], 0
	for i := 1; i < len(x)-1; i++ {
		if x[i] > value {
			value, index = x[i], i
		}
	}
	return
}

// FindMin returns the value and index of the first element of x equal to the minimum value in x.
func FindMin(x []float64) (value float64, index int) {
	value, index = x[0], 0
	for i := 1; i < len(x)-1; i++ {
		if x[i] < value {
			value, index = x[i], i
		}
	}
	return
}

/*
Float32ToFloat64 returns a copy of x with type []float64
*/
func Float32ToFloat64(x []float32) []float64 {
	y := make([]float64, len(x))
	for i, f := range x {
		y[i] = float64(f)
	}
	return y
}

func IsPowerOf2(x int) bool {
	return (x != 0) && ((x & (x - 1)) == 0)
}

/*
LoadFloats reads a text file containing one float per line.
*/
func LoadFloats(fname string) []float64 {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	rdr := bufio.NewReader(bytes.NewBuffer(data))
	x := make([]float64, 0, 1024)
	for s, err := rdr.ReadString('\n'); err == nil; s, err = rdr.ReadString('\n') {
		f, err := strconv.ParseFloat(strings.TrimSuffix(s, "\n"), 64)
		if err != nil {
			panic(err)
		}
		x = append(x, f)
	}
	return x
}

// Log2 returns the integer log base 2 of n.
// E.g.: log2(12) ~ 3.6. Log2 returns 3
func Log2(n int) int {
	return int(math.Log2(float64(n)))
}

/*
LowpassFilterAll returns LowpassFilter(x) for all x in xs.
*/
func LowpassFilterAll(xs [][]float64, alpha float64) [][]float64 {
	ys := make([][]float64, len(xs))
	for i, x := range xs {
		ys[i] = LowpassFilter(x, alpha)
	}
	return ys
}

/*
LowpassFilter returns x filtered by alpha
*/
func LowpassFilter(x []float64, alpha float64) []float64 {
	y := make([]float64, len(x))
	y[0] = alpha * x[0]
	for i := 1; i < len(x); i++ {
		y[i] = y[i-1] + alpha*(x[i]-y[i-1])
	}
	return y
}

// Max returns the maximum value of the elements of x
func Max(x []float64) float64 {
	max := x[0]
	for _, f := range x {
		if f > max {
			max = f
		}
	}
	return max
}

// MaxInt returns the maximum value of the elements of x
func MaxInt(x []int) int {
	max := x[0]
	for _, f := range x {
		if f > max {
			max = f
		}
	}
	return max
}

/*
MovAvg returns the moving average for each x[i], given by sum(x[i-w:i+w])/(2w)
*/
func MovAvg(x []float64, w int) []float64 {
	y := make([]float64, len(x))
	for i := w; i < len(x)-w; i++ {
		y[i] = Sum(x[i-w:i+w]) / float64(2*w)
	}
	return y
}

/*
Multiplex returns on vector with the element of vs interleaved
*/
func Multiplex(channels [][]float64) []float64 {
	numChans := len(channels)
	chanLen := len(channels[0])
	buf := make([]float64, numChans*chanLen)
	for i := 0; i < chanLen; i++ {
		k := i * numChans
		for j := 0; j < numChans; j++ {
			buf[k+j] = channels[j][i]
		}
	}
	return buf
}

// Normalise returns x/max(x)
func Normalise(x []float64) []float64 {
	x1 := make([]float64, len(x))
	sum := Max(x)
	for i, f := range x {
		x1[i] = f / sum
	}
	return x1
}

// Normalise returns x/max(x) for all x in xs
func NormaliseAll(xs [][]float64) [][]float64 {
	x1 := make([][]float64, len(xs))
	for i, x := range xs {
		x1[i] = Normalise(x)
	}
	return x1
}

// Pow2 returns 2^x.
// The function panics if x < 0
func Pow2(x int) int {
	if x < 0 {
		panic(fmt.Sprintf("X = %d", x))
	}
	pw := 1
	for i := 1; i <= x; i++ {
		pw *= 2
	}
	return pw
}

// Range returns an interger range 0:1:n-1
func Range(n int) []int {
	rng := make([]int, n)
	for i := range rng {
		rng[i] = i
	}
	return rng
}

/*
RemoveAvgAllZ removes the average of all vectors x in xs. The minimum value
of any x[i] is 0.
*/
func RemoveAvgAllZ(xs [][]float64) [][]float64 {
	xs1 := make([][]float64, len(xs))
	for i, x := range xs {
		xs1[i] = RemoveAvg(x)
	}
	return xs1
}

// RemoveAvgZ returns x[i] = x[i]-sum(x)/len(x) or 0 if x[i]-sum(x)/len(x) < 0
func RemoveAvg(x []float64) []float64 {
	x1 := make([]float64, len(x))
	avg := Sum(x) / float64(len(x))
	for i, f := range x {
		x1[i] = f - avg
		if x1[i] < 0 {
			x1[i] = 0
		}
	}
	return x1
}

// Smooth smoothts x: x[i] = sum(x[i-wdw:i+wdw])/(2*wdw)
func Smooth(x []float64, wdw int) {
	for i := 0; i < wdw; i++ {
		x[i] = 0
	}
	for i := wdw; i < len(x)-wdw; i++ {
		x[i] = Sum(x[i-wdw:i+wdw]) / float64((2 * wdw))
	}
}

/*
Sub returns x - y. The function panics if len(x) != len(y).
*/
func Sub(x, y []float64) []float64 {
	if len(x) != len(y) {
		panic("len(x) != len(y)")
	}
	x1 := make([]float64, len(x))
	for i := range x {
		x1[i] = x[i] - y[i]
	}
	return x1
}

// Sum returns the sum of the elements of the vector x
func Sum(x []float64) float64 {
	sum := 0.0
	for _, f := range x {
		sum += f
	}
	return sum
}

// SumVectors returns the sum of the vectors in X.
// The function panics if all vectors don't have the same length
func SumVectors(X [][]float64) []float64 {
	N := len(X[0])
	for i, x := range X {
		if len(x) != N {
			panic(fmt.Sprintf("N=%d but len(X[%d]=%d", N, i, len(x)))
		}
	}
	sum := make([]float64, N)
	for i := 0; i < N; i++ {
		for j := range X {
			sum[i] += X[j][i]
		}
	}
	return sum
}

func ToFloat(x []int) []float64 {
	y := make([]float64, len(x))
	for i, e := range x {
		y[i] = float64(e) / float64(math.MaxInt64)
	}
	return y
}

/*
ToInt returns y * math.MaxInt64.
The range of x is [-1.0,1.0].
The function panics if bitsPerSample is not one of 8,16,32.
*/
func ToInt(x []float64, bitsPerSample int) []int {
	y := make([]int, len(x))
	if bitsPerSample != 8 && bitsPerSample != 16 && bitsPerSample != 32 {
		panic(fmt.Sprintf("Invalid bitsPerSample %d", bitsPerSample))
	}
	max := float64(int(1)<<bitsPerSample - 1)
	for i, f := range x {
		y[i] = int(f * max)
	}
	return y
}

func ToIntS(x float64, bitsPerSample int) int {
	max := float64(int(1)<<bitsPerSample - 1)
	return int(x * max)
}

func findLocalMax(x []float64, from, wdw, step int) (maxI, slopeEnd int) {
	i, slp := from+wdw, 0
	for slp >= 0 && i < len(x)-wdw {
		slp = slope(x[i : i+wdw])
		i += step
	}
	_, maxI = FindMax(x[from:i])
	maxI += from
	slopeEnd = i
	return
}

func findLocalMin(x []float64, from, wdw, step int) (minI, slopeEnd int) {
	i, slp := from+wdw, 0
	for slp <= 0 && i < len(x)-wdw {
		slp = slope(x[i : i+wdw])
		i += step
	}
	_, minI = FindMin(x[from:i])
	minI += from
	slopeEnd = i
	return
}

func findNon0Slope(x []float64, from, wdw int) (slp, end int) {
	for i := from; i < len(x)-wdw; i++ {
		slp := slope(x[i : i+wdw])
		if slp != 0 {
			return slp, i
		}
	}
	return 0, len(x)
}

// slope returns +1, 0, -1
func slope(x []float64) int {
	end := len(x) - 1
	if x[0] < x[end] {
		return -1
	}
	if x[0] == x[end] {
		return 0
	}
	return 1
}

func ivecContain(x []int, v int) bool {
	for _, v1 := range x {
		if v1 == v {
			return true
		}
	}
	return false
}

// WriteAllDataFile writes each xs[i] in xs to a test file `fname_i.txt`
func WriteAllDataFile(xs [][]float64, fname string) {
	for i, xs := range xs {
		WriteDataFile(xs, fmt.Sprintf("%s_%d", fname, i))
	}
}

// WriteDataFile writes x to a text file `fname.txt`
func WriteDataFile(x []float64, fname string) {
	buf := new(bytes.Buffer)
	for _, f := range x {
		fmt.Fprintf(buf, "%f\n", f)
	}
	if err := myioutil.WriteFile(fname+".txt", buf.Bytes()); err != nil {
		panic(err)
	}
}

// WriteIntDataFile writes x to a text file `fname.txt`
func WriteIntDataFile(x []int, fname string) {
	buf := new(bytes.Buffer)
	for _, f := range x {
		fmt.Fprintf(buf, "%d\n", f)
	}
	if err := myioutil.WriteFile(fname+".txt", buf.Bytes()); err != nil {
		panic(err)
	}
}

/*
WriteIntMatrixDataFile writes an integer matrix to a text file `fname.csv`
*/
func WriteIntMatrixDataFile(x [][]int, fname string) {
	buf := new(bytes.Buffer)
	for _, row := range x {
		for i, col := range row {
			if i > 0 {
				fmt.Fprint(buf, ",")
			}
			fmt.Fprintf(buf, "%d", col)
		}
		fmt.Fprintln(buf)
	}
	if err := myioutil.WriteFile(fname+".csv", buf.Bytes()); err != nil {
		panic(err)
	}
}

/*
Xcorr returns the cross correlation of x with y for maxDelay.
*/
func Xcorr(x, y []float64, maxDelay int) (corr []float64) {
	N := len(x)
	corr = make([]float64, maxDelay)
	for k := 0; k < maxDelay; k++ {
		for n := 0; n < N-k; n++ {
			corr[k] += x[n] * y[n+k]
		}
		corr[k] /= float64(N)
	}
	return
}
