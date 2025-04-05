package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type bank []float64

const (
	rub = iota
	usd
	eur
)

func solve(banks []bank) float64 {
	a, b, c := banks[0], banks[1], banks[2]

	cache := do3(a, b, c, []float64{1, 0, 0})
	ans := cache[usd]

	cache = do3(a, c, b, []float64{1, 0, 0})
	ans = max(ans, cache[usd])

	cache = do3(b, a, c, []float64{1, 0, 0})
	ans = max(ans, cache[usd])

	cache = do3(b, c, a, []float64{1, 0, 0})
	ans = max(ans, cache[usd])

	cache = do3(c, a, b, []float64{1, 0, 0})
	ans = max(ans, cache[usd])

	cache = do3(c, b, a, []float64{1, 0, 0})
	ans = max(ans, cache[usd])

	return ans
}

func do3(a, b, c bank, cache []float64) []float64 {
	cache = do(a, cache)
	cache = do(b, cache)
	cache = do(c, cache)
	return cache
}

func do(bank bank, cache []float64) []float64 {
	// log.Println(bank, cache)

	buf := make([]float64, 3)
	copy(buf, cache)

	// Курс обмена рублей на доллары.
	buf[usd] = max(buf[usd], bank[0]*cache[rub])

	// Курс обмена рублей на евро.
	buf[eur] = max(buf[eur], bank[1]*cache[rub])

	// Курс обмена долларов на рубли.
	buf[rub] = max(buf[rub], bank[2]*cache[usd])

	// Курс обмена долларов на евро.
	buf[eur] = max(buf[eur], bank[3]*cache[usd])

	// Курс обмена евро на рубли.
	buf[rub] = max(buf[rub], bank[4]*cache[eur])

	// Курс обмена евро на доллары.
	buf[usd] = max(buf[usd], bank[5]*cache[eur])

	// log.Println(buf)
	copy(cache, buf)
	return cache
}

func run(r io.Reader, w io.Writer) {
	br := bufio.NewReader(r)
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	var t int
	if _, err := fmt.Fscanln(br, &t); err != nil {
		panic(err)
	}

	banks := make([]bank, 3)
	for i := 0; i < 3; i++ {
		banks[i] = make([]float64, 6)
	}

	for i := 0; i < t; i++ {
		for j := 0; j < 3; j++ {
			for k := 0; k < 6; k++ {
				var n, m int
				if _, err := fmt.Fscanln(br, &n, &m); err != nil {
					panic(err)
				}
				banks[j][k] = float64(m) / float64(n)
			}
		}
		// log.Println(banks)
		ans := solve(banks)
		fmt.Fprintf(bw, "%f\n", ans)
	}
}

func main() {
	run(os.Stdin, os.Stdout)
}
