package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

const usage = `usage: pipe text into pct

pct calculates the distribution of lines in text.
It is similar to sort | uniq -c, except that it
prints percentages as well as counts. pct's output
is unordered; use sort -n.
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

	tw := new(tabwriter.Writer)
	tw.Init(w, 0, 0, 4, ' ', 0)
	f := float64(n)
	for k, v := range m {
		fmt.Fprintf(tw, "% 6.2f%%\t%d\t%s\n", 100*float64(v)/f, v, k)
	}
	return tw.Flush()
}

func main() {
	if len(os.Args) > 1 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	if err := pct(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
