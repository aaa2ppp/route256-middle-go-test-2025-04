package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func solve(s string) bool {
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

func run(r io.Reader, w io.Writer) {
	br := bufio.NewReader(r)
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	var t int
	if _, err := fmt.Fscanln(br, &t); err != nil {
		panic(err)
	}
	for i := 0; i < t; i++ {
		// TODO
		var s string
		if _, err := fmt.Fscanln(br, &s); err != nil {
			panic(err)
		}
		if solve(s) {
			fmt.Fprintln(bw, "YES")
		} else {
			fmt.Fprintln(bw, "NO")
		}
	}
}

func main() {
	run(os.Stdin, os.Stdout)
}
