// package dump is a utility for printing memory allocation info using printf.
//
// Memory stats:
// 	dump.Mem("some_func - BEFORE")
//  some_func()
//  dump.Mem("some_func - AFTER")
//
// Output is like this:
// 	2019/10/10 14:37:48 some_func - BEFORE
// 	2019/10/10 14:37:48   HeapAlloc  : 1.99G, delta: 408.46M
// 	2019/10/10 14:37:48   HeapObjects: 18.36M, delta: 4.58M
// 	2019/10/10 14:37:51 some_func - AFTER
// 	2019/10/10 14:37:51   HeapAlloc  : 3.15G, delta: 1.15G
// 	2019/10/10 14:37:51   HeapObjects: 22.98M, delta: 4.62M
package dump

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"sync"
)

const (
	kilo float64 = 1024.0
	mega float64 = 1024.0 * kilo
	giga float64 = 1024.0 * mega
)

// Printf calls standard log.Printf, here to avoid importing "log".
func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// prevM keeps previous memory status to calculate deltas.
var prevM = &runtime.MemStats{}

// lock for thread safe Mem calls.
var lock = sync.Mutex{}

// Mem prints memory status to stderr with given message.
// Thread safe.
func Mem(message string) {
	lock.Lock()
	defer lock.Unlock()

	currentM := &runtime.MemStats{}
	runtime.ReadMemStats(currentM)
	printMem(message, prevM, currentM)
	prevM = currentM
}

func printMem(message string, prevM, currentM *runtime.MemStats) {
	heapDelta := delta64(currentM.HeapAlloc, prevM.HeapAlloc)
	objDelta := delta64(currentM.HeapObjects, prevM.HeapObjects)

	log.Print(message)
	log.Printf("  HeapAlloc  : %s, delta: %s", Meg(currentM.HeapAlloc), Meg(heapDelta))
	log.Printf("  HeapObjects: %s, delta: %s", Meg(currentM.HeapObjects), Meg(objDelta))
}

// Meg prints out a number in human readable form, e.g. 20, 20K, 20M, 20G.
// Returns "NaN" if the input is not a number (int, float).
func Meg(n interface{}) string {
	// in the order of likelihood
	switch t := n.(type) {
	case uint64:
		return megInt64(float64(t))
	case int64:
		return megInt64(float64(t))

	case int:
		return megInt64(float64(t))

	case float64:
		return megInt64(float64(t))
	case float32:
		return megInt64(float64(t))

	case int32:
		return megInt64(float64(t))
	case uint32:
		return megInt64(float64(t))

	case int16:
		return megInt64(float64(t))

	case uint16:
		return megInt64(float64(t))

	case int8:
		return megInt64(float64(t))
	case uint8: // also byte
		return megInt64(float64(t))

	default:
		return "NaN"
	}
}

func megInt64(x float64) string {
	xAbs := math.Abs(x)
	switch {
	case xAbs > giga:
		return fmt.Sprintf("%.2fG", x/giga)
	case xAbs > mega:
		return fmt.Sprintf("%.2fM", x/mega)
	case xAbs > kilo:
		return fmt.Sprintf("%.2fK", x/kilo)
	default:
		return fmt.Sprintf("%d", int(x))
	}
}

// delta64 calculates delta from two uint64.
func delta64(current, prev uint64) int64 {
	if current > prev {
		return int64(current - prev)
	} else {
		return -int64(prev - current)
	}
}
