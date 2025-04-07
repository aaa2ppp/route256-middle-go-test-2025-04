package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
)

type params struct {
	boxes     [][2]int
	paintings [][2]int
}

type results struct {
	minBoxesCount int
}

func readParams(br *bufio.Reader) params {
	boxes := readSizes(br)
	paintings := readSizes(br)
	return params{boxes, paintings}
}

func readSizes(br *bufio.Reader) [][2]int {
	var n int
	if _, err := fmt.Fscanln(br, &n); err != nil {
		panic(err)
	}

	sizes := make([][2]int, 0, n)
	for i := 0; i < n; i++ {
		var a, b int
		if _, err := fmt.Fscanln(br, &a, &b); err != nil {
			panic(err)
		}
		sizes = append(sizes, [2]int{a, b})
	}

	return sizes
}

func writeResults(bw *bufio.Writer, results results) {
	fmt.Fprintln(bw, results.minBoxesCount)
}

func solve(params params) results {
	boxes := params.boxes
	paintings := params.paintings

	normalizeSizes(boxes)
	normalizeSizes(paintings)

	count := 0
	i, j := len(boxes)-1, len(paintings)-1
	size1 := -1
	for i >= 0 && j >= 0 {
		// ищем коробку в которую влезает картина по размеру #0 такую,
		// чтобы рамер #1 был максимальный
		for i >= 0 && paintings[j][0] <= boxes[i][0] {
			if boxes[i][1] > size1 {
				size1 = boxes[i][1]
			}
			i--
		}

		// кладем картины, пока влезают по размеру #1
		var ok bool
		for j >= 0 && paintings[j][1] <= size1 {
			ok = true
			j--
		}

		if !ok {
			// не удалось положить не одной картины
			return results{-1}
		}

		// считаем использованную коробку
		count++
	}

	if j >= 0 {
		// остались неупакованные картины
	}

	return results{count}
}

// normalizeSizes разворачивает объект, так чтобы рамер #0 был меньше
// размера #1 и сортирует по рамеру #0, затем по рамеру #1
func normalizeSizes(rects [][2]int) {
	for i, rect := range rects {
		if rect[0] > rect[1] {
			rects[i] = [2]int{rect[1], rect[0]}
		}
	}
	slices.SortFunc(rects, func(a, b [2]int) int {
		if a[0] == b[0] {
			return a[1] - b[1]
		}
		return a[0] - b[0]
	})
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
