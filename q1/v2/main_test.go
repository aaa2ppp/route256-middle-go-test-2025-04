package main

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"unsafe"
)

const testDataDir = "../test_data/inserting-chars"

type runTest struct {
	name    string
	in      io.Reader
	wantOut string
	debug   bool
}

func (tt runTest) do(t *testing.T, solve solveFunc) {
	t.Run(tt.name, func(t *testing.T) {
		defer func(v bool) { debugEnable = true }(debugEnable)
		debugEnable = tt.debug

		w := &strings.Builder{}
		run(tt.in, w, solve)
		if gotW := w.String(); gotW != tt.wantOut {
			t.Errorf("run() = %v, want %v", gotW, tt.wantOut)
		}
	})
}

func Test_run(t *testing.T) {
	solve := solve
	tests := []runTest{
		{
			"1",
			strings.NewReader(`2
abacaa
PppP
`),
			`YES
NO
`,
			true,
		},
		{
			"2",
			strings.NewReader(`2
X
kkk
`),
			`YES
YES
`,
			true,
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		tt.do(t, solve)
	}
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
