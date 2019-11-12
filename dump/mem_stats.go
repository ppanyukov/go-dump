package dump

import (
	"bytes"
	"fmt"
	"runtime"
)

func NewMemStats(name string) *MemStats {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	return &MemStats{
		Name:        name,
		HeapAlloc:   int64(stats.HeapAlloc),
		HeapObjects: int64(stats.HeapObjects),
	}
}

func NewMemStatsDiff(base *MemStats, next *MemStats) *MemStatsDiff {
	return &MemStatsDiff{
		Base: base,
		Next: next,
		Delta: &MemStats{
			Name:        fmt.Sprintf("Delta: %s -> %s", base.Name, next.Name),
			HeapAlloc:   next.HeapAlloc - base.HeapAlloc,
			HeapObjects: next.HeapObjects - base.HeapObjects,
		},
	}
}

type MemStats struct {
	Name        string
	HeapAlloc   int64
	HeapObjects int64
}

// NewDiff creates new MemProfDiff with this MemProf as a base.
func (m *MemStats) NewDiff(name string) *MemStatsDiff {
	return NewMemStatsDiff(m, NewMemStats(name))
}

func (m *MemStats) String() string {
	buf := &bytes.Buffer{}
	_ = fmt.Sprint(buf, "MEM STATS: %s\n", m.Name)
	_ = fmt.Sprint(buf, "  HeapAlloc  : %s\n", Meg(m.HeapAlloc))
	_ = fmt.Sprint(buf, "  HeapObjects: %s", Meg(m.HeapObjects))
	return buf.String()
}

func (m *MemStats) Print() {
	fmt.Printf("%s\n", m)
}

type MemStatsDiff struct {
	Base  *MemStats
	Next  *MemStats
	Delta *MemStats
}

func (m *MemStatsDiff) String() string {
	buf := &bytes.Buffer{}
	_ = fmt.Sprint(buf, "MEM STATS DIFF: %s -> %s\n", m.Base.Name, m.Next.Name)
	_ = fmt.Sprint(buf, "  HeapAlloc  : %s -> %s -> %s\n", Meg(m.Base.HeapAlloc), Meg(m.Next.HeapAlloc), Meg(m.Delta.HeapAlloc))
	_ = fmt.Sprint(buf, "  HeapObjects: %s -> %s - %s", Meg(m.Base.HeapObjects), Meg(m.Next.HeapObjects), Meg(m.Delta.HeapObjects))
	return buf.String()
}

func (m *MemStatsDiff) Print() {
	fmt.Printf("%s\n", m)
}
