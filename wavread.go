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
	"bytes"
	"io/ioutil"

	"github.com/mjibson/go-dsp/wav"
)

/*
ReadWavFile returns the demultiplexed channels of a wav file, and the sample rate in Hz.
*/
func ReadWavFile(wavName string) (channels [][]float64, sampleRate, bitsPerSample int) {
	buf, err := ioutil.ReadFile(wavName)
	if err != nil {
		panic(err)
	}
	rdr, err := wav.New(bytes.NewBuffer(buf))
	if err != nil {
		panic(err)
	}
	numSamples, numChannels := rdr.Samples, int(rdr.NumChannels)
	sampleRate = int(rdr.SampleRate)
	bitsPerSample = int(rdr.Header.BitsPerSample)
	channels = make([][]float64, numChannels)
	chanLen := numSamples / numChannels
	for i := range channels {
		channels[i] = make([]float64, chanLen)
	}
	samples, err := rdr.ReadFloats(rdr.Samples)
	if err != nil {
		panic(err)
	}
	for i, j := 0, 0; i < len(samples); {
		for _, ch := range channels {
			ch[j] = float64(samples[i])
			i++
		}
		j++
	}
	return
}
