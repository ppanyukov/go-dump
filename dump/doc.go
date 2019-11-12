// package dump is a utility for printing memory allocations info using printf.
//
// Two main memory stats:
//
// - MemStats:
//   This uses `runtime.ReadMemStats`.
//
// - MemProf
//   This uses code from pprof package to show stats as would be
//   shown by pprof tool.
//
package dump


