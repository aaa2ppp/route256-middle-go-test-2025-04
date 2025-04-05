package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
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
				strings.NewReader(`2
abacaa
PppP
`),
			},
			`YES
NO
`,
		},
		{
			"2",
			args{
				strings.NewReader(`2
X
kkk
`),
			},
			`YES
YES
`,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &strings.Builder{}
			run(tt.args.r, w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("run() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

var dataDir = filepath.Join("test_data", "inserting-chars")

func Test_run2(t *testing.T) {
	type args struct {
		in io.Reader
	}
	type test struct {
		name    string
		args    args
		wantOut []byte
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
			a,
		})
	}

	if len(tests) == 0 {
		t.Fatal("no any test")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &strings.Builder{}
			run(tt.args.in, w)
			if gotW := w.String(); gotW != string(tt.wantOut) {
				t.Errorf("run() = %v, want %v", gotW, string(tt.wantOut))
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
