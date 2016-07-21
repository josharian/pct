package main

import (
	"container/heap"
	"math"
	"sort"

	"github.com/spaolacci/murmur3"
)

type stringCount struct {
	s string // string
	n uint64 // count
}

type stringsByCount []stringCount

func (l stringsByCount) Len() int { return len(l) }
func (l stringsByCount) Less(i, j int) bool {
	x, y := l[i], l[j]
	if x.n != y.n {
		return x.n > y.n // largest-to-smallest
	}
	return x.s < y.s
}
func (l stringsByCount) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

// topk approximately counts the top k strings in a stream.
// It uses a count min sketch implementation.
type topK struct {
	x    [][]uint64 // counts
	w, d int        // dimensions
	top  topHeap    // top items
	h    []uint64   // cached hash offsets
}

type topHeap []stringCount

func (h topHeap) Len() int            { return len(h) }
func (h topHeap) Less(i, j int) bool  { return h[i].n < h[j].n }
func (h topHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *topHeap) Push(x interface{}) { *h = append(*h, x.(stringCount)) }
func (h *topHeap) Pop() interface{} {
	x := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return x
}

// k is the number of top entries to track.
// w is the width of each array, typically small (4-8).
// d is the number of hash functions, typically medium (1024-16384).
// consult your favorite count min sketch exposition
// to learn how to interpret and pick these.
func newTopK(k, d, w int) *topK {
	t := topK{
		x:   make([][]uint64, d),
		top: make(topHeap, k),
		h:   make([]uint64, d),
		w:   w,
		d:   d,
	}
	for i := range t.x {
		t.x[i] = make([]uint64, w)
	}
	return &t
}

// Record records b. b will not be modified or retained,
// so that it is safe for use (and efficient to use) with a bufio.Scanner.
func (t *topK) Record(b []byte) {
	x, y := murmur3.Sum128(b)
	max := uint64(math.MaxUint64)
	for i := range t.h {
		t.h[i] = (x + uint64(i)*y) % uint64(t.w)
	}
	// find max
	for i, h := range t.h {
		if v := t.x[i][h]; v < max {
			max = v
		}
	}
	// allow it to increase by one
	max++
	// increment entries as long as they don't exceed max
	for i, h := range t.h {
		if v := t.x[i][h]; v < max {
			t.x[i][h] = v + 1
		}
	}
	// update top entries if needed
	if max < t.top[0].n {
		return
	}
	for i, e := range t.top {
		if string(b) == e.s && e.n > 0 {
			// already in the top, update counts
			t.top[i].n = max
			heap.Fix(&t.top, i)
			return
		}
	}
	// add to the top, ejecting the lowest element
	t.top[0] = stringCount{s: string(b), n: max}
	heap.Fix(&t.top, 0)
}

func (t *topK) Top(n int) []stringCount {
	if n != len(t.top) {
		panic("can only retrieve as many entries as you request up front")
	}
	top := make([]stringCount, len(t.top))
	copy(top, t.top)
	sort.Sort(stringsByCount(top))
	return top
}

func (t *topK) All() []stringCount {
	panic("topk cannot print all entries")
}
