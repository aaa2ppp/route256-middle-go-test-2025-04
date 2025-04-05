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
		in io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			"1",
			args{strings.NewReader(`7
3
ababa
ababa
ababa
3
asd
das
sda
2
abca
abc
4
aaaa
aaaa
aaaa
aaa
2
aa
aa
2
a
a
2
a
b
`)},
			`3
0
1
6
1
1
0
`,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			run(tt.args.in, out)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

var testData = filepath.Join("test_data", "even-strings")

func Test_run2(t *testing.T) {
	type args struct {
		in io.Reader
	}
	type test struct {
		name    string
		args    args
		wantOut string
	}
	var tests []test

	files, err := os.ReadDir(testData)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		if !strings.HasSuffix(fileName, ".a") {
			continue
		}

		askName := strings.TrimSuffix(fileName, ".a")
		ansName := fileName

		ask, err := os.ReadFile(filepath.Join(testData, askName))
		if err != nil {
			panic(err)
		}
		ans, err := os.ReadFile(filepath.Join(testData, ansName))
		if err != nil {
			panic(err)
		}

		tests = append(tests, test{
			askName,
			args{bytes.NewReader(ask)},
			string(ans),
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			run(tt.args.in, out)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
