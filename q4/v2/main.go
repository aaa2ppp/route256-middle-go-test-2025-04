package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func solve(k int, desk []string) bool {
	n, m := len(desk), len(desk[0])

	total := 0
	if m >= k {
		for i := 0; i < n; i++ {
			if ans := checkLine(k, desk, i, 0, 0, 1); ans == -1 {
				return false
			} else {
				total |= ans
			}
		}
	}

	if n >= k {
		for j := 0; j < m; j++ {
			if ans := checkLine(k, desk, 0, j, 1, 0); ans == -1 {
				return false
			} else {
				total |= ans
			}
		}
	}

	for i := 0; i < n-k+1; i++ {
		if ans := checkLine(k, desk, i, 0, 1, 1); ans == -1 {
			return false
		} else {
			total |= ans
		}
	}

	for j := 1; j < m-k+1; j++ {
		if ans := checkLine(k, desk, 0, j, 1, 1); ans == -1 {
			return false
		} else {
			total |= ans
		}
	}

	for i := k - 1; i < n; i++ {
		if ans := checkLine(k, desk, i, 0, -1, 1); ans == -1 {
			return false
		} else {
			total |= ans
		}
	}

	for j := 1; j < m-k+1; j++ {
		if ans := checkLine(k, desk, n-1, j, -1, 1); ans == -1 {
			return false
		} else {
			total |= ans
		}
	}

	return total != 0
}

// checkLine движется от (i,j) с шагом (di,dj), возвращает:
//
//	 1 если можно поставить X, чтобы в линию было k,
//	-1 если линия из O или X уже существует,
//	 0 иначе.
func checkLine(k int, desk []string, i, j, di, dj int) int {
	// n, m := uint(len(desk)), uint(len(desk[0]))

	ans := 0
	lx, rx, dot := 0, 0, -1
	lo, ro := 0, 0

	// for ; uint(i) < n && uint(j) < m; i, j = i+di, j+dj {
	for ; uint(i) < uint(len(desk)) && uint(j) < uint(len(desk[i])); i, j = i+di, j+dj {
		c := desk[i][j]

		if c == 'O' {
			ro++
			lx = ro
			rx = ro
			dot = -1

			if ro-lo == k {
				return -1
			}

			continue
		}

		if c == '.' {
			if dot != -1 {
				lx = dot + 1
			}
			dot = rx
		}

		rx++
		lo = rx
		ro = rx

		if rx-lx == k {
			if dot == -1 {
				return -1
			}

			ans = 1
			lx = dot + 1
			dot = -1
			continue
		}

		if lx == dot && rx-lx-1 == k {
			return -1
		}
	}

	return ans
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
		var k int
		if _, err := fmt.Fscanln(br, &k); err != nil {
			panic(err)
		}
		var n, m int
		if _, err := fmt.Fscanln(br, &n, &m); err != nil {
			panic(err)
		}
		desk := make([]string, 0, n)
		for j := 0; j < n; j++ {
			s, err := br.ReadString('\n')
			if err != nil {
				panic(err)
			}
			s = strings.TrimRight(s, "\r\n")
			desk = append(desk, s)
		}
		if solve(k, desk) {
			fmt.Fprintln(bw, "YES")
		} else {
			fmt.Fprintln(bw, "NO")
		}
	}
}

func main() {
	run(os.Stdin, os.Stdout)
}
