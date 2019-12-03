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
Package ioutil implements some convenience functions for working with files.
*/
package ioutil

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// FilePermission given to all created files and directories
const FilePermission = 0731

/*
MkdirAll creates all directories in `path` which don't exist.
*/
func MkdirAll(path string) error {
	if path == "" {
		return nil
	}
	return os.MkdirAll(path, FilePermission)
}

/*
WriteFile creates all non-existent directories in `path` and writes data to
the filename in `path`. The permissions are `FilePermission`.
*/
func WriteFile(path string, data []byte) error {
	dir, _ := filepath.Split(path)
	if err := MkdirAll(dir); err != nil {
		return fmt.Errorf("Error creating directory %s: %s", dir, err)
	}
	if err := ioutil.WriteFile(path, data, FilePermission); err != nil {
		return fmt.Errorf("error writing file %s: %s", path, err)
	}
	return nil
}

/*
ReadIntVector reads the elements of an integer vector from a text file, `fname`.
Each vector element is on its own line
*/
func ReadIntVector(fname string) ([]int, error) {
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	rdr := csv.NewReader(bytes.NewBuffer(buf))
	recs, err := rdr.ReadAll()
	if err != nil {
		return nil, err
	}
	x := make([]int, len(recs))
	for i, r := range recs {
		intval, err := strconv.Atoi(r[0])
		if err != nil {
			return nil, err
		}
		x[i] = intval
	}
	return x, nil
}
