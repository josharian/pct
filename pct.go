package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
)

const usage = `usage: ... | pct

pct calculates the distribution of lines in text.
It is similar to sort | uniq -c | sort -n -r, except
that it prints percentages as well as counts.
`

func pct(r io.Reader, w io.Writer) error {
	s := bufio.NewScanner(r)
	m := make(map[string]int)
	n := 0
	for s.Scan() {
		m[s.Text()]++
		n++
	}
	if err := s.Err(); err != nil {
		return err
	}
	if n == 0 {
		return nil
	}

	var l lines
	for k, v := range m {
		l = append(l, line{n: v, s: k})
	}
	sort.Sort(l)

	f := float64(n)
	for _, line := range l {
		_, err := fmt.Fprintf(w, "% 6.2f%%% 6d %s\n", 100*float64(line.n)/f, line.n, line.s)
		if err != nil {
			return err
		}
	}
	return nil
}

type line struct {
	n int
	s string
}

type lines []line

func (l lines) Len() int { return len(l) }
func (l lines) Less(i, j int) bool {
	x, y := l[i], l[j]
	if x.n != y.n {
		return x.n > y.n // largest-to-smallest
	}
	return x.s < y.s
}
func (l lines) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

func main() {
	if len(os.Args) > 1 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	pct(os.Stdin, os.Stdout)
}
