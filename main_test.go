package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"sort"
	"strings"
	"testing"
)

func Test_genFilenames(t *testing.T) {
	var tests = []struct {
		name        string
		expected    []string
		givenArgs   []string
		givenReader io.Reader
	}{
		{"should return args in channel one by one if provided",
			[]string{"foo", "bar", "buzz"},
			[]string{"foo", "bar", "buzz"},
			&bytes.Buffer{},
		},
		{"should read args from the reader if args param is not provided",
			[]string{"foo", "bar", "buzz"},
			[]string{},
			bytes.NewBufferString("foo\nbar\nbuzz"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			i := 0
			for actual := range genFilenames(tt.givenArgs, tt.givenReader) {
				expected := tt.expected[i]
				if actual != expected {
					t.Errorf("(%s): expected %s, actual %s", tt.givenArgs, expected, actual)
				}
				i++
			}
		})
	}
}

func mustReadFile(t *testing.T, name string) string {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		t.Fatalf("mustReadFile: %s", err)
	}
	return string(b)
}

func sortedLines(s string) []string {
	lines := strings.Split(s, "\n")
	sort.Strings(lines)
	return lines
}

func Test_parseSource(t *testing.T) {
	var tests = []struct {
		expected string
		given    string
	}{
		{
			"testdata/source.expected.gogr",
			"testdata/source.go",
		},
		{
			"testdata/edge_case_nil.expected.gogr",
			"testdata/edge_case_nil.go",
		},
		{
			"testdata/anonymous_functions_and_structs.expected.gogr",
			"testdata/anonymous_functions_and_structs.go",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.given, func(t *testing.T) {

			src := mustReadFile(t, tt.given)
			w := &bytes.Buffer{}
			if err := parseSource(tt.given, string(src), w); err != nil {
				t.Fatalf("parseSource: %s", err)
			}
			actualLines := sortedLines(w.String())

			expected := mustReadFile(t, tt.expected)
			expectedLines := sortedLines(expected)
			if len(actualLines) != len(expectedLines) {
				t.Errorf("different lengthes: expected %d, actual %d",
					len(expectedLines), len(actualLines))
			}
			for i, actual := range actualLines {
				if expectedLines[i] != actual {
					t.Errorf("line %d - not equal:\nexpected: %s\n  actual: %s",
						i, expectedLines[i], actual)
				}
			}
		})
	}
}
