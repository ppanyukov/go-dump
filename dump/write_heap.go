package dump

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"runtime/pprof"
)

// HeapDumpDir is the directory where the heap dump will be written. By default empty (current working dir).
// Set to whatever is required.
var HeapDumpDir = func() string {
	// ignore error
	dir, _ := os.Getwd()
	return dir
}()

// WriteHeapDump is a convenience method to write a pprof heap dump which can be
// examined with pprof tool. The file `{name}.pb.gz` is written to `HeapDumpDir`.
//
// Call at the start and end of the function to see how much and where things
// were allocated within that function.
func WriteHeapDump(name string) {
	var fName = fmt.Sprintf("%s.pb.gz", name)
	fName = path.Join(HeapDumpDir, fName)

	f, _ := os.Create(fName)
	defer f.Close()
	runtime.GC()
	err := pprof.WriteHeapProfile(f)
	if err != nil {
		fmt.Printf("ERROR WRITING HEAP DUMP TO %s: %v\n", fName, err)
		return
	}

	fmt.Printf("WRITTEN HEAP DUMP TO %s\n", fName)
}
