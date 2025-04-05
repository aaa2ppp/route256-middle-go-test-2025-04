package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func solve(ss []string) int {
	var bu strings.Builder
	ms := make(map[string]int)
	ma := make(map[string]int)
	mb := make(map[string]int)

	count := 0
	for _, s := range ss {
		bu.Reset()
		bu.Grow((len(s)+1)/2 + 1)
		for i := 0; i < len(s); i += 2 {
			bu.WriteByte(s[i])
		}
		sa := bu.String()
		count += ma[sa]
		ma[sa]++

		if len(s) == 1 {
			continue
		}

		bu.Reset()
		bu.Grow((len(s)+1)/2 + 1)
		for i := 1; i < len(s); i += 2 {
			bu.WriteByte(s[i])
		}
		sb := bu.String()
		count += mb[sb]
		mb[sb]++

		count -= ms[s]
		ms[s]++
	}

	return count
}

func run(in io.Reader, out io.Writer) {
	br := bufio.NewReader(in)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	var t int
	if _, err := fmt.Fscanln(br, &t); err != nil {
		panic(err)
	}

	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscanln(br, &n); err != nil {
			panic(err)
		}
		ss := make([]string, 0, n)
		for j := 0; j < n; j++ {
			var s string
			if _, err := fmt.Fscanln(br, &s); err != nil {
				panic(err)
			}
			ss = append(ss, s)
		}
		ans := solve(ss)
		fmt.Fprintln(bw, ans)
	}
}

func main() {
	run(os.Stdin, os.Stdout)
}
