package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
)

func usage() {
	const help = `usage: ... | pct [-f] [-n]

pct calculates the distribution of lines in text.
It is similar to sort | uniq -c | sort -n -r, except
that it prints percentages as well as counts.
`

	fmt.Fprintln(os.Stderr, help)
	flag.PrintDefaults()
}

func dump(w io.Writer, m map[string]int) error {
	var l lines
	var tot int
	for k, v := range m {
		l = append(l, line{n: v, s: k})
		tot += v
	}
	sort.Sort(l)

	f := 100 / float64(tot)
	lim := *limit
	for i := 0; i < len(l) && (lim <= 0 || i < lim); i++ {
		line := l[i]
		_, err := fmt.Fprintf(w, "% 6.2f%%% 6d %s\n", f*float64(line.n), line.n, line.s)
		if err != nil {
			return err
		}
	}
	return nil
}

func pct(r io.Reader, w io.Writer) error {
	s := bufio.NewScanner(r)
	m := make(map[string]int)
	n := 0
	for s.Scan() {
		m[s.Text()]++
		n++
		if *every > 0 && n%*every == 0 {
			if err := dump(w, m); err != nil {
				return err
			} else {
				fmt.Fprintln(w)
			}
		}
	}
	if err := s.Err(); err != nil {
		dump(w, m)
		fmt.Fprintf(w, "Stopped at line %d: %v\n", n, err)
		return err
	}

	return dump(w, m)
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

var (
	every = flag.Int("f", 0, "print running percents every f lines")
	limit = flag.Int("n", 0, "only print top n lines")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	pct(os.Stdin, os.Stdout)
}
