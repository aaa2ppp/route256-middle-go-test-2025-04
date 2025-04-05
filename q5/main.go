package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"slices"
)

func solve(serverThrouputs, imageWeights []int) (minTimeDifference int, imageServers []int) {
	type assignment struct {
		image  int
		server int
		time   int
	}

	// Составляем список всех возможных назначений серверов для изображений,
	// отсортированный по времени доставки

	serversCount := len(serverThrouputs)
	imagesCount := len(imageWeights)
	allAssignments := make([]assignment, 0, serversCount*imagesCount)

	for i, throuput := range serverThrouputs {
		for j, weight := range imageWeights {
			allAssignments = append(allAssignments, assignment{
				server: i,
				image:  j,
				time:   (weight + throuput - 1) / throuput},
			)
		}
	}

	slices.SortFunc(allAssignments, func(a, b assignment) int {
		return a.time - b.time
	})

	if debugEnable {
		log.Println("all possible assignments:", allAssignments)
	}

	// Ищем минимальный диапазон (по разности времени доставки), который покрывает все изображения

	minTimeDifference = math.MaxInt
	rangeStart, rangeEnd := -1, -1

	counts := make([]int, imagesCount)
	covered := 0
	l, r := 0, 0
	for l < len(allAssignments) && r < len(allAssignments) {
		for covered < imagesCount && r < len(allAssignments) {
			if debugEnable {
				log.Println(l, r, "covered:", covered)
			}
			if counts[allAssignments[r].image] == 0 {
				covered++
			}
			counts[allAssignments[r].image]++
			r++
		}

		for covered >= imagesCount && l < len(allAssignments) {
			if dif := allAssignments[r-1].time - allAssignments[l].time; dif < minTimeDifference {
				if debugEnable {
					log.Println(l, r, "covered:", covered, "bingo!", dif)
				}
				minTimeDifference = dif
				rangeStart, rangeEnd = l, r
			} else {
				if debugEnable {
					log.Println(l, r, "covered:", covered)
				}
			}
			counts[allAssignments[l].image]--
			if counts[allAssignments[l].image] == 0 {
				covered--
			}
			l++
		}
	}

	if debugEnable {
		log.Println("best assignments range:", rangeStart, rangeEnd, allAssignments[rangeStart:rangeEnd])
	}

	// Составляем список серверов для изображения (нас устроит любой подходящий сервер для изображения)

	imageServers = make([]int, imagesCount)
	for i := rangeStart; i < rangeEnd; i++ {
		imageServers[allAssignments[i].image] = allAssignments[i].server
	}

	if debugEnable {
		log.Println("result:", minTimeDifference, imageServers)
	}

	return minTimeDifference, imageServers
}

func readInts(br *bufio.Reader) []int {
	var n int
	if _, err := fmt.Fscanln(br, &n); err != nil {
		panic(err)
	}
	nums := make([]int, 0, n)
	for i := 0; i < n; i++ {
		var v int
		if _, err := fmt.Fscan(br, &v); err != nil {
			panic(err)
		}
		nums = append(nums, v)
	}
	br.ReadString('\n')
	return nums
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
		throuputs := readInts(br)
		weights := readInts(br)
		minDif, servers := solve(throuputs, weights)
		fmt.Fprintln(bw, minDif)
		for _, server := range servers {
			fmt.Fprintf(bw, "%d ", server+1) // to 1-indexing
		}
		fmt.Fprintln(bw)
	}
}

var debugEnable bool

func main() {
	_, debugEnable = os.LookupEnv("DEBUG")
	run(os.Stdin, os.Stdout)
}
