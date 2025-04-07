package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
)

type row []byte

type params struct {
	k    int
	n, m int
	desk []row
}

type results struct {
	ok bool
}

func readParams(br *bufio.Reader) params {
	var k, n, m int
	if _, err := fmt.Fscanln(br, &k); err != nil {
		panic(err)
	}
	if _, err := fmt.Fscanln(br, &n, &m); err != nil {
		panic(err)
	}

	// NOTE: На доску добавляем барьер из 0-й серху, справа и снизу
	desk := make([]row, 0, n+2)
	desk = append(desk, make([]byte, m+1)) // барьер сверху

	for i := 0; i < n; i++ {
		s := make([]byte, m+1) // выделяем память под строку +1 байт для барьера
		// br.Read(s[:m])         // !!! не всегда читает len(s) байт
		io.ReadFull(br, s[:m])                        // читаем m байт
		if _, err := br.ReadSlice('\n'); err != nil { // вычитываем конец строки
			panic(err)
		}
		desk = append(desk, s)
	}

	desk = append(desk, make([]byte, m+1)) // барьер снизу

	return params{k, n, m, desk}
}

func writeResults(bw *bufio.Writer, results results) {
	if results.ok {
		bw.WriteString("YES\n")
	} else {
		bw.WriteString("NO\n")
	}
}

func solve(params params) results {
	k, n, m, desk := params.k, params.n, params.m, params.desk

	// NOTE(!): i индексируется начиная с 1, а j - начиная с 0

	ok := func() bool {
		ans := 0

		if m >= k {
			for i, j := 1, 0; i <= n; i++ {
				if ans |= checkLine(k, desk, i, j, 0, 1); ans == -1 {
					return false
				}
				if i <= n-k+1 {
					if ans |= checkLine(k, desk, i, j, 1, 1); ans == -1 {
						return false
					}
				}
				if i >= k {
					if ans |= checkLine(k, desk, i, j, -1, 1); ans == -1 {
						return false
					}
				}
			}
		}

		if n >= k {
			for i, j := 1, 0; j < m; j++ {
				if ans |= checkLine(k, desk, i, j, 1, 0); ans == -1 {
					return false
				}
				if j >= 1 && j <= m-k {
					if ans |= checkLine(k, desk, i, j, 1, 1); ans == -1 {
						return false
					}
					if ans |= checkLine(k, desk, i, j, -1, 1); ans == -1 {
						return false
					}
				}
			}
		}

		return ans != 0
	}()

	return results{ok}
}

// checkLine движется от (i,j) с шагом (di,dj) до барьера.
// Возвращает:
//
//	 1 если можно поставить 'X', чтобы в линию было k,
//	-1 если линия из k 'O' или 'X' уже существует,
//	 0 иначе.
func checkLine(k int, desk []row, i, j, di, dj int) (ans int) {
	olen := 0          // длина последовательности состоящей только из 'O'
	xlen := 0          // длина последовательности состоящей только из 'X' и не более одной точки
	dot := math.MinInt // дистанция до точки. если < 0, то xlen не содержит точку

	// loop:
	for {
		switch desk[i][j] {
		case 'O':
			dot = math.MaxInt
			xlen = 0
			olen++
			if olen == k {
				return -1
			}
		case '.':
			if dot > 0 {
				xlen = dot - 1
			}

			dot = 1
			olen = 0
			xlen++
			if xlen == k {
				// bingo!
				ans = 1

				// вырезаем точку
				xlen = 0
				dot = math.MinInt

				// break loop
			}
		case 'X':
			dot++
			olen = 0
			xlen++
			if xlen < k {
				break // switch
			}
			if xlen == k {
				if dot < 0 { // точка отсутствует
					return -1
				}

				// bingo!
				ans = 1

				// вырезаем точку
				xlen = dot - 1
				dot = math.MinInt

				// break loop
				break // switch
			}
			if xlen > k {
				// возможно толко, если после точки следует последовательность из k 'X'
				return -1
			}
		default: // уперлись в барьер
			return ans
		}
		i, j = i+di, j+dj
	}

	// Хвостовой цикл. Уже найден положительный ответ.
	// Осталось проверить, что дальше нет последовательностей из k 'O' или 'X'.
	// По идее, этот цикл должен быть быстрее чем предыдущий, т.к. проще.
	// Но ощутимого прироста я не почувствовал.

	i, j = i+di, j+dj
	for {
		switch desk[i][j] {
		case 'O':
			xlen = 0
			olen++
			if olen == k {
				return -1
			}
		case '.':
			olen = 0
			xlen = 0
		case 'X':
			olen = 0
			xlen++
			if xlen == k {
				return -1
			}
		default:
			return 1
		}
		i, j = i+di, j+dj
	}
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
