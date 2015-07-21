// Sourced is a scripts' preprocessor for Source-based games.
// It does some solid heavy lifting for configurations in
// games like Team Fortress 2 or Counter-Strike.
package main

import (
	"flag"
	"fmt"
	"github.com/tucnak/sourced/parser"
	"io/ioutil"
)

func main() {
	inputFlag := flag.String("input", "", "input file")
	outputFlag := flag.String("output", "", "output file")
	flag.Parse()

	if *inputFlag == "" || *outputFlag == "" {
		return
	}

	input, err := ioutil.ReadFile(*inputFlag)
	if err != nil {
		fmt.Println(err)
		return
	}

	output, err := parser.Build(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(*outputFlag, output, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
