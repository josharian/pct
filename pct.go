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
	const help = `usage: ... | pct [-f] [-n] [-c]

pct calculates the distribution of lines in text.
It is similar to sort | uniq -c | sort -n -r, except
that it prints percentages as well as counts.
`

	fmt.Fprintln(os.Stderr, help)
	flag.PrintDefaults()
}

type recorder interface {
	Record([]byte)
	Top(int) []stringCount
}

type mcount map[string]uint64

func (m mcount) Record(b []byte) {
	m[string(b)]++
}

func (m mcount) Top(n int) []stringCount {
	var l []stringCount
	for k, v := range m {
		l = append(l, stringCount{n: v, s: k})
	}
	sort.Sort(stringsByCount(l))
	if n > len(l) {
		return l
	}
	return l[:n]
}

func dump(w io.Writer, tot int, r recorder) error {
	f := 100 / float64(tot)
	runtot := uint64(0)
	top := r.Top(*limit)
	for _, line := range top {
		runtot += line.n
		p := f * float64(line.n)
		var err error
		if *cum {
			_, err = fmt.Fprintf(w, "% 6.2f%% % 6.2f%%% 6d %s\n", p, f*float64(runtot), line.n, line.s)
		} else {
			_, err = fmt.Fprintf(w, "% 6.2f%%% 6d %s\n", p, line.n, line.s)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func pct(r io.Reader, w io.Writer, rec recorder) error {
	s := bufio.NewScanner(r)
	n := 0
	for s.Scan() {
		rec.Record(s.Bytes())
		n++
		if *every > 0 && n%*every == 0 {
			if err := dump(w, n, rec); err != nil {
				return err
			} else {
				fmt.Fprintln(w)
			}
		}
	}
	if err := s.Err(); err != nil {
		dump(w, n, rec)
		fmt.Fprintf(w, "Stopped at line %d: %v\n", n, err)
		return err
	}
	return dump(w, n, rec)
}

var (
	every  = flag.Int("f", 0, "print running percents every f lines, requires -n")
	limit  = flag.Int("n", 0, "only print top n lines")
	cum    = flag.Bool("c", false, "print cumulative percents as well")
	approx = flag.Bool("x", false, "use a fast approximate counter, only suitable for large input, requires -n")
)

func main() {
	flag.Usage = usage
	flag.Parse()
	if *approx && *limit == 0 {
		fmt.Fprintln(os.Stderr, "-x requires -n")
		os.Exit(2)
	}
	if *every != 0 && *limit == 0 {
		fmt.Fprintln(os.Stderr, "-f requires -n")
		os.Exit(2)
	}

	var r recorder
	if *approx {
		r = newTopK(*limit, 8, 16384)
	} else {
		r = make(mcount)
	}

	pct(os.Stdin, os.Stdout, r)
}
