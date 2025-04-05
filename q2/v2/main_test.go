package main

import (
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"unsafe"
)

const testDataDir = "../test_data/three-banks"

type runTest struct {
	name    string
	in      io.Reader
	wantOut string
	debug   bool
}

func Test_run(t *testing.T) {
	solve := solve
	tests := []runTest{
		{
			"1",
			strings.NewReader(`1
100 1
100 1
1 100
3 2
1 100
2 3
100 1
100 1
1 100
3 2
1 100
2 3
100 1
100 1
1 100
3 2
1 100
2 3
`),
			`0.015
`,
			true,
		},
		// {
		// 	"2",
		// 	strings.NewReader(``),
		// 	``,
		// 	true,
		// },
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		tt.do(t, solve)
	}
}

func (tt runTest) do(t *testing.T, solve solveFunc) {
	t.Run(tt.name, func(t *testing.T) {
		defer func(v bool) { debugEnable = true }(debugEnable)
		debugEnable = tt.debug

		w := &strings.Builder{}
		run(tt.in, w, solve)
		// if gotW := w.String(); gotW != tt.wantOut {
		// 	t.Errorf("run() = %v, want %v", gotW, tt.wantOut)
		// }
		want := linesToFloats(lines(tt.wantOut))
		got := linesToFloats(lines(w.String()))
		if !floatsIsEqual(want, got) {
			t.Errorf("run() = %v, want %v", w.String(), tt.wantOut)
		}
	})
}

func lines(s string) []string {
	lines := strings.Split(strings.TrimRight(s, " \t\r\n"), "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " \t\r\n")
	}
	return lines
}

func linesToFloats(lines []string) []float64 {
	floats := make([]float64, len(lines))
	for i, s := range lines {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
		floats[i] = v
	}
	return floats
}

func floatsIsEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if math.Abs(a[i]-b[i]) > 1e-6 && math.Abs(a[i]-b[i])/a[i] > 1e-6 {
			return false
		}
	}
	return true
}

func Test_run_fullset(t *testing.T) {
	solve := solve

	files, err := os.ReadDir(testDataDir)
	if err != nil {
		panic(err)
	}

	var testNums []int
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()

		if !strings.HasSuffix(fileName, ".a") {
			continue
		}

		testName := strings.TrimSuffix(fileName, ".a")
		testNum, err := strconv.Atoi(testName)
		if err != nil {
			t.Log(err)
			continue
		}

		testNums = append(testNums, testNum)
	}

	if len(testNums) == 0 {
		t.Log("no any test")
		return
	}

	for _, testNum := range testNums {
		testName := strconv.Itoa(testNum)
		testPath := filepath.Join(testDataDir, testName)

		wantOut, err := os.ReadFile(testPath + ".a")
		if err != nil {
			t.Fatal(err)
		}

		testFile, err := os.Open(testPath)
		if err != nil {
			t.Fatal(err)
		}
		defer testFile.Close()

		runTest{
			name:    testName,
			in:      testFile,
			wantOut: unsafe.String(unsafe.SliceData(wantOut), len(wantOut)),
			debug:   false,
		}.do(t, solve)
	}
}
