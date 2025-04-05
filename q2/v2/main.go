package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type bank struct {
	exchangeRate []float64
}

type params struct {
	banks [3]bank
}

type results struct {
	maxDollars float64
}

const (
	// • Курс обмена рублей на доллары.
	rub2USD = iota
	// • Курс обмена рублей на евро.
	rub2EUR
	// • Курс обмена долларов на рубли.
	usd2RUB
	// • Курс обмена долларов на евро.
	usd2EUR
	// • Курс обмена евро на рубли.
	eur2RUB
	// • Курс обмена евро на доллары.
	eur2USD
)

func readParams(br *bufio.Reader) params {
	var params params
	for i := 0; i < 3; i++ {
		exchangeRate := make([]float64, 0, 6)
		for i := 0; i < 6; i++ {
			var n, m int
			if _, err := fmt.Fscanln(br, &n, &m); err != nil {
				panic(err)
			}
			exchangeRate = append(exchangeRate, float64(m)/float64(n))
		}
		params.banks[i] = bank{exchangeRate}
	}
	return params
}

func writeResults(bw *bufio.Writer, results results) {
	fmt.Fprintf(bw, "%f\n", results.maxDollars)
}

type cache struct {
	rub,
	usd,
	eur float64
}

func solve(params params) results {
	cache := cache{rub: 1, usd: 0, eur: 0}

	banks := params.banks
	a, b, c := banks[0], banks[1], banks[2]

	maxDollars := max(
		excangeChain(cache, a, b, c).usd,
		excangeChain(cache, a, c, b).usd,
		excangeChain(cache, b, a, c).usd,
		excangeChain(cache, b, c, a).usd,
		excangeChain(cache, c, a, b).usd,
		excangeChain(cache, c, b, a).usd,
	)

	return results{maxDollars}
}

func excangeChain(cache cache, banks ...bank) cache {
	for _, bank := range banks {
		cache = excange(cache, bank.exchangeRate)
	}
	return cache
}

func excange(cache cache, rate []float64) cache {
	out := cache

	out.usd = max(out.usd, cache.rub*rate[rub2USD])
	out.eur = max(out.eur, cache.rub*rate[rub2EUR])

	out.rub = max(out.rub, cache.usd*rate[usd2RUB])
	out.eur = max(out.eur, cache.usd*rate[usd2EUR])

	out.rub = max(out.rub, cache.eur*rate[eur2RUB])
	out.usd = max(out.usd, cache.eur*rate[eur2USD])

	return out
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
