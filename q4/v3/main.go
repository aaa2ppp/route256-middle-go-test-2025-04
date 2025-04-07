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

	// На доску добавляем барьер из 0-ей серху, справа и снизу
	// !!!NOTE!!!: Это меняет начало индесации по вертикали с 0 на 1!
	desk := make([]row, 0, n+2)
	desk = append(desk, make([]byte, m+1)) // барьер сверху

	for i := 0; i < n; i++ {
		s := make([]byte, m+1)                        // выделяем память под строку +1 байт для барьера слева
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

	// !!!NOTE!!!: по вертикали (i) индексируемся с 1, а по горизонтали (j) - начиная с 0
	//  (см. readParams)

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

// checkLine движется от (i,j) с шагом (di,dj) до барьера (любой символ отличный от '.', 'O' или 'X').
// Возвращает:
//
//	 1 если можно поставить 'X', чтобы в линию было k,
//	-1 если линия из k 'O' или 'X' уже существует,
//	 0 иначе.
func checkLine(k int, desk []row, i, j, di, dj int) (ans int) {
	olen := 0          // длина последовательности состоящей только из 'O'
	xlen := 0          // длина последовательности состоящей только из 'X' и не более одной точки
	dot := math.MinInt // дистанция до точки. если < 0, то xlen не содержит точку

	for {
		// TODO: Эти проверки после тестирования можно перенести внутрь кейсов.
		// Это даст -1 или -2 if'а на итерацию.
		if olen == k {
			return -1
		}
		if xlen == k && dot > 0 {
			ans = 1
			xlen = dot - 1
			dot = math.MinInt
		}
		if xlen >= k {
			return -1
		}

		switch desk[i][j] {
		case 'O':
			olen++
			xlen = 0
			dot = math.MinInt
		case '.':
			if dot > 0 {
				xlen = dot - 1
			}
			olen = 0
			xlen++
			dot = 1
		case 'X':
			olen = 0
			xlen++
			dot++
		default: // уперлись в барьер
			return ans
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
