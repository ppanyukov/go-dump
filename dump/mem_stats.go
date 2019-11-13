package dump

import (
	"bytes"
	"fmt"
	"runtime"
	"text/tabwriter"
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

func newMemStatsDiff(base *MemStats, next *MemStats) *MemStatsDiff {
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

// PrintDiff creates new snapshot diff and prints it. Here to avoid pitfalls of defer etc.
func (m *MemStats) PrintDiff() {
	name := fmt.Sprintf("%s - AFTER", m.Name)
	diff := newMemStatsDiff(m, NewMemStats(name))
	diff.Print()
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
	buffer := &bytes.Buffer{}
	tw := tabwriter.NewWriter(buffer, 1, 8, 1, '\t', 0)
	_, _ = fmt.Fprintf(tw, "MEM STATS DIFF:   \t%s \t%s \t-> %s \t\n", m.Base.Name, m.Next.Name, "Delta")
	_, _ = fmt.Fprintf(tw, "    HeapAlloc  : \t%s \t%s \t-> %s \t\n", Meg(m.Base.HeapAlloc), Meg(m.Next.HeapAlloc), Meg(m.Delta.HeapAlloc))
	_, _ = fmt.Fprintf(tw, "    HeapObjects: \t%s \t%s \t-> %s \t\n", Meg(m.Base.HeapObjects), Meg(m.Next.HeapObjects), Meg(m.Delta.HeapObjects))
	tw.Flush()
	return buffer.String()
}

func (m *MemStatsDiff) Print() {
	fmt.Printf("%s\n", m)
}
