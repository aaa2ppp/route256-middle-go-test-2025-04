package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strconv"
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
			args{strings.NewReader(`3
3
3 3
X..
..O
OOX
2
5 3
...
O.O
X.O
...
...
3
5 5
X.X..
.....
.OX..
..O..
...O.
`)},
			`YES
NO
NO
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

var testData = filepath.Join("../test_data", "tic-tac-toe-middle")

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

func Test_checkLine(t *testing.T) {
	type args struct {
		k      int
		desk   []string
		target byte
		i      int
		j      int
		di     int
		dj     int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"0",
			args{
				1,
				[]string{
					"X",
				},
				'X',
				0, 0,
				0, 1,
			},
			-1,
		},
		{
			"1",
			args{
				3,
				[]string{
					"XX.",
				},
				'X',
				0, 0,
				0, 1,
			},
			1,
		},
		{
			"1",
			args{
				3,
				[]string{
					"XXX.",
				},
				'X',
				0, 0,
				0, 1,
			},
			-1,
		},
		{
			"1",
			args{
				3,
				[]string{
					".XXX.",
				},
				'X',
				0, 0,
				0, 1,
			},
			-1,
		},
		{
			"1",
			args{
				3,
				[]string{
					".X.XXX.",
				},
				'X',
				0, 0,
				0, 1,
			},
			-1,
		},
		{
			"1",
			args{
				3,
				[]string{
					".X.X.",
				},
				'X',
				0, 0,
				0, 1,
			},
			1,
		},
		{
			"1",
			args{
				3,
				[]string{
					".XX",
				},
				'X',
				0, 0,
				0, 1,
			},
			1,
		},
		{
			"1",
			args{
				3,
				[]string{
					"..XX",
				},
				'X',
				0, 0,
				0, 1,
			},
			1,
		},
		{
			"1",
			args{
				3,
				[]string{
					"X.X",
				},
				'X',
				0, 0,
				0, 1,
			},
			1,
		},
		{
			"1.2",
			args{
				3,
				[]string{
					"XXX",
				},
				'X',
				0, 0,
				0, 1,
			},
			-1,
		},
		{
			"2",
			args{
				3,
				[]string{
					"X",
					".",
					"X",
				},
				'X',
				0, 0,
				1, 0,
			},
			1,
		},
		{
			"3",
			args{
				3,
				[]string{
					"X..",
					"...",
					"..X",
				},
				'X',
				0, 0,
				1, 1,
			},
			1,
		},
		{
			"4",
			args{
				3,
				[]string{
					"X..",
					"..O",
					"OOX",
				},
				'X',
				0, 0,
				1, 1,
			},
			1,
		},
		{
			"5",
			args{
				3,
				[]string{
					"X.X..",
					".....",
					".OX..",
					"..O..",
					"...O.",
				},
				'O',
				1, 0,
				1, 1,
			},
			-1,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkLine(tt.args.k, tt.args.desk, tt.args.i, tt.args.j, tt.args.di, tt.args.dj); got != tt.want {
				t.Errorf("checkLine() = %v, want %v", got, tt.want)
			}
			// if got := checkLine(tt.args.k, strDeskToByteDesk1(tt.args.desk), tt.args.i, tt.args.j, tt.args.di, tt.args.dj); got != tt.want {
			// 	t.Errorf("checkLine() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func strDeskToByteDesk1(strDesk []string) [][]byte {
	n, m := len(strDesk), len(strDesk[0])
	byteDesk := make([][]byte, n)

	for i := 0; i < n; i++ {
		byteDesk[i] = append(make([]byte, 0, m), strDesk[i]...)
	}

	return byteDesk
}

func Benchmark_run(b *testing.B) {
	const testNum = 35
	for i := 0; i < b.N; i++ {
		func() {
			testName := strconv.Itoa(testNum)
			testPath := filepath.Join(testData, testName)
			testFile, err := os.Open(testPath)
			if err != nil {
				panic(err)
			}
			defer testFile.Close()
			run(testFile, io.Discard)
		}()
	}
}
