package main

import (
	"bytes"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func Test_run(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		{
			"1",
			args{
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
			},
			`0.015
`,
		},
		// {
		// 	"2",
		// 	args{
		// 		strings.NewReader(``),
		// 	},
		// 	``,
		// },
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &strings.Builder{}
			run(tt.args.r, w)
			want := linesToFloats(lines(tt.wantW))
			got := linesToFloats(lines(w.String()))
			if !floatsIsEqual(want, got) {
				t.Errorf("run() = %v, want %v", w.String(), tt.wantW)
			}
		})
	}
}

var dataDir = filepath.Join("test_data", "three-banks")

func Test_run2(t *testing.T) {
	type args struct {
		r io.Reader
	}
	type test struct {
		name  string
		args  args
		wantW string
	}
	var tests []test

	files, err := os.ReadDir(dataDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()

		if ok, err := filepath.Match("*.a", fileName); err != nil {
			panic(err)
		} else if !ok {
			continue
		}

		testName := strings.TrimSuffix(fileName, ".a")
		q, err := os.ReadFile(filepath.Join(dataDir, testName))
		if err != nil {
			panic(err)
		}

		a, err := os.ReadFile(filepath.Join(dataDir, fileName))
		if err != nil {
			panic(err)
		}

		tests = append(tests, test{
			testName,
			args{bytes.NewBuffer(q)},
			string(a),
		})
	}

	if len(tests) == 0 {
		t.Fatal("no any test")
	}

	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		w := &strings.Builder{}
	// 		run(tt.args.in, w)
	// 		if gotW := w.String(); gotW != string(tt.wantOut) {
	// 			t.Errorf("run() = %v, want %v", gotW, string(tt.wantOut))
	// 		}
	// 	})
	// }
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &strings.Builder{}
			run(tt.args.r, w)
			// want, err := strconv.ParseFloat(strings.TrimSpace(tt.wantW), 64)
			// if err != nil {
			// 	panic(err)
			// }
			// got, err := strconv.ParseFloat(strings.TrimSpace(w.String()), 64)
			// if err != nil {
			// 	panic(err)
			// }
			// if math.Abs(want-got) > 1e-6 && math.Abs(want-got)/want > 1e-6 {
			// 	t.Errorf("run() = %v, want %v", got, want)
			// }
			want := linesToFloats(lines(tt.wantW))
			got := linesToFloats(lines(w.String()))
			if !floatsIsEqual(want, got) {
				t.Errorf("run() = %v, want %v", w.String(), tt.wantW)
			}
		})
	}
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
