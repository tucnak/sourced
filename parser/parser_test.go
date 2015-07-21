package parser

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	testFiles, err := ioutil.ReadDir("tests")
	if err != nil {
		t.Fatalf("No tests exist!")
	}

	var tests []string
	for _, file := range testFiles {
		filename := file.Name()
		if strings.HasSuffix(filename, ".s") {
			tests = append(tests, filename[:len(filename)-2])
		}
	}

	for _, testname := range tests {
		input, err := ioutil.ReadFile("tests/" + testname + ".s")
		if err != nil {
			t.Errorf("Unable to fetch %s.s: %s", testname, err)
			continue
		}

		expected, err := ioutil.ReadFile("tests/" + testname + ".out")
		if err != nil {
			t.Errorf("Unable to fetch %s.out: %s", testname, err)
			continue
		}

		real, err := Build(input)
		if err != nil {
			t.Errorf("Failed to build %s.s: %s", testname, err)
			continue
		}

		if !bytes.Equal(real, expected) {
			t.Errorf("Cooked %s.s doesn't match expected output!", testname)
			t.Logf("Expected:\n%s", string(expected))

			if len(real) == 0 {
				real = []byte("<nothing>")
			}
			t.Logf("Recieved:\n%s", string(real))
			t.Fatal("Oops.")
		}
	}
}
