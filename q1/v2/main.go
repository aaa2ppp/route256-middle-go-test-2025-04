package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type params struct {
	s string
}

type results struct {
	ok bool
}

func readParams(br *bufio.Reader) params {
	s, err := br.ReadString('\n')
	if err != nil {
		panic(err)
	}
	s = strings.TrimSpace(s)
	return params{s}
}

func writeResults(bw *bufio.Writer, results results) {
	if results.ok {
		bw.WriteString("YES\n")
	} else {
		bw.WriteString("NO\n")
	}
}

func solve(params params) results {
	return results{_solve(params.s)}
}

func _solve(s string) bool {
	if s[0] != s[len(s)-1] {
		return false
	}
	for i := 2; i < len(s); i++ {
		if s[i] != s[0] {
			continue
		}
		if s[i] != s[i-1] && s[i] != s[i-2] {
			return false
		}
	}
	return true
}

// ----------------------------------------------------------------------------

type solveFunc func(params) results

func runTask(br *bufio.Reader, bw *bufio.Writer, solve solveFunc) {
	writeResults(bw, solve(readParams(br)))
}

func run(r io.Reader, w io.Writer, solve solveFunc) {
	br := bufio.NewReader(r)
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	var t int
	if _, err := fmt.Fscanln(br, &t); err != nil {
		panic(err)
	}
	for i := 0; i < t; i++ {
		runTask(br, bw, solve)
	}
}

var debugEnable bool

func main() {
	_, debugEnable = os.LookupEnv("DEBUG")
	run(os.Stdin, os.Stdout, solve)
}
