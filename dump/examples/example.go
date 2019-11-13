// Copyright 2019-current Go-dump Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	"github.com/ppanyukov/go-dump/dump"
)

func Example() []int64 {
	// Write pprof heap profile at the start and end of function
	dump.WriteHeapDump(fmt.Sprintf("heap-example-before"))
	defer dump.WriteHeapDump(fmt.Sprintf("heap-example-after"))

	// Take a snapshot at the start of a function
	// Capture and print deltas at the end
	memProf := dump.NewMemProf("example")
	defer memProf.PrintDiff()

	// Similar for memStats
	memStats := dump.NewMemStats("example")
	defer memStats.PrintDiff()

	// allocate memory
	allocateMem := func () []int64 {
		return make([]int64, 100000)
	}

	var result []int64
	for i := 0; i < 1000; i++{
		result = append(result, allocateMem()...)
	}

	return result
}

func main() {
	array := Example()
	fmt.Printf("Number of int64s: %d\n", len(array))
}
