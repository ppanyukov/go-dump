package dump

import (
	"bytes"
	"fmt"
	"math"
	"runtime"
	"sort"
	"text/tabwriter"
)

// NewMemProf creates memory heap profile as would shown by pprof tool.
// Code here is base largely on pprof package with lots of copy/paste.
func NewMemProf(name string) *MemProf {
	// NOTE: this does seem to be required for profiles to work
	runtime.GC()
	runtime.GC()

	memProf := getHeapInternal()
	allocObjects, allocBytes := scaleHeapSample(memProf.AllocObjects, memProf.AllocBytes, int64(runtime.MemProfileRate))
	inUseObjects, inUseBytes := scaleHeapSample(memProf.InUseObjects(), memProf.InUseBytes(), int64(runtime.MemProfileRate))
	return &MemProf{
		Name:         name,
		AllocObjects: allocObjects,
		AllocBytes:   allocBytes,
		InUseObjects: inUseObjects,
		InUseBytes:   inUseBytes,
	}
}

func newMemProfDiff(base *MemProf, next *MemProf) *MemProfDiff {
	return &MemProfDiff{
		Base: base,
		Next: next,
		Delta: &MemProf{
			Name:         fmt.Sprintf("Delta: %s -> %s", base.Name, next.Name),
			AllocObjects: next.AllocObjects - base.AllocObjects,
			AllocBytes:   next.AllocBytes - base.AllocBytes,
			InUseObjects: next.InUseObjects - base.InUseObjects,
			InUseBytes:   next.InUseBytes - base.InUseBytes,
		},
	}
}

// MemProf is memory heap profile as would be shown by pprof tool.
type MemProf struct {
	Name         string
	AllocObjects int64
	AllocBytes   int64
	InUseObjects int64
	InUseBytes   int64
}

// PrintDiff creates new snapshot diff and prints it. Here to avoid pitfalls of defer etc.
func (m *MemProf) PrintDiff() {
	name := fmt.Sprintf("%s - AFTER", m.Name)
	diff := newMemProfDiff(m, NewMemProf(name))
	diff.Print()
}

func (m *MemProf) String() string {
	buffer := &bytes.Buffer{}
	_, _ = fmt.Fprintf(buffer, "MEM PROF DIFF: %s\n", m.Name)
	_, _ = fmt.Fprintf(buffer, "    AllocObjects: %s\n", Meg(m.AllocObjects))
	_, _ = fmt.Fprintf(buffer, "    AllocBytes  : %s\n", Meg(m.AllocBytes))
	_, _ = fmt.Fprintf(buffer, "    InUseObjects: %s\n", Meg(m.InUseObjects))
	_, _ = fmt.Fprintf(buffer, "    InUseBytes  : %s", Meg(m.InUseBytes))
	return buffer.String()
}

func (m *MemProf) Print() {
	fmt.Printf("%s\n", m)
}

type MemProfDiff struct {
	Base  *MemProf
	Next  *MemProf
	Delta *MemProf
}

func (m *MemProfDiff) String() string {
	buffer := &bytes.Buffer{}
	tw := tabwriter.NewWriter(buffer, 1, 8, 1, '\t', 0)
	_, _ = fmt.Fprintf(tw, "MEM PROF DIFF:    \t%s \t%s \t-> %s \t\n", m.Base.Name, m.Next.Name, "Delta")
	_, _ = fmt.Fprintf(tw, "    AllocObjects: \t%s \t%s \t-> %s \t\n", Meg(m.Base.AllocObjects), Meg(m.Next.AllocObjects), Meg(m.Delta.AllocObjects))
	_, _ = fmt.Fprintf(tw, "    AllocBytes  : \t%s \t%s \t-> %s \t\n", Meg(m.Base.AllocBytes), Meg(m.Next.AllocBytes), Meg(m.Delta.AllocBytes))
	_, _ = fmt.Fprintf(tw, "    InUseObjects: \t%s \t%s \t-> %s \t\n", Meg(m.Base.InUseObjects), Meg(m.Next.InUseObjects), Meg(m.Delta.InUseObjects))
	_, _ = fmt.Fprintf(tw, "    InUseBytes  : \t%s \t%s \t-> %s \t\n", Meg(m.Base.InUseBytes), Meg(m.Next.InUseBytes), Meg(m.Delta.InUseBytes))
	tw.Flush()
	return buffer.String()
}

func (m *MemProfDiff) Print() {
	fmt.Printf("%s\n", m)
}

// scaleHeapSample adjusts the data from a heap Sample to
// account for its probability of appearing in the collected
// data. heap profiles are a sampling of the memory allocations
// requests in a program. We estimate the unsampled value by dividing
// each collected sample by its probability of appearing in the
// profile. heap profiles rely on a poisson process to determine
// which samples to collect, based on the desired average collection
// rate R. The probability of a sample of size S to appear in that
// profile is 1-exp(-S/R).
//
// copypasta from pprof
func scaleHeapSample(count, size, rate int64) (int64, int64) {
	if count == 0 || size == 0 {
		return 0, 0
	}

	if rate <= 1 {
		// if rate==1 all samples were collected so no adjustment is needed.
		// if rate<1 treat as unknown and skip scaling.
		return count, size
	}

	avgSize := float64(size) / float64(count)
	scale := 1 / (1 - math.Exp(-avgSize/float64(rate)))

	return int64(float64(count) * scale), int64(float64(size) * scale)
}

// copypasta from pprof.WriteHeapProfile.
// we only care for the total
func getHeapInternal() runtime.MemProfileRecord {
	// Find out how many records there are (MemProfile(nil, true)),
	// allocate that many records, and get the data.
	// There's a race—more records might be added between
	// the two calls—so allocate a few extra records for safety
	// and also try again if we're very unlucky.
	// The loop should only execute one iteration in the common case.
	var p []runtime.MemProfileRecord
	n, ok := runtime.MemProfile(nil, true)
	for {
		// Allocate room for a slightly bigger profile,
		// in case a few more entries have been added
		// since the call to MemProfile.
		p = make([]runtime.MemProfileRecord, n+50)
		n, ok = runtime.MemProfile(p, true)
		if ok {
			p = p[0:n]
			break
		}
		// Profile grew; try again.
	}

	sort.Slice(p, func(i, j int) bool { return p[i].InUseBytes() > p[j].InUseBytes() })

	var total runtime.MemProfileRecord
	for i := range p {
		r := &p[i]
		total.AllocBytes += r.AllocBytes
		total.AllocObjects += r.AllocObjects
		total.FreeBytes += r.FreeBytes
		total.FreeObjects += r.FreeObjects
	}

	return total
}
