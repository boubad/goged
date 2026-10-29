package main

import (
	"flag"
	"fmt"
	"services.diarra.ovh/gogedinfo/infoged"
)

func main() {
	fixture := infoged.NewTestFixture()
	// Define optional flags
	input := flag.String("i", fixture.TestGedcomFilePath, "Input file (optional)")
	output := flag.String("o", fixture.TestOutputFilePath, "Output file (optional)")

	// Parse flags
	flag.Parse()
	// Print results
	fmt.Println("Input File:", *input)
	fmt.Println("Output File:", *output)
	//
	otree, err := infoged.NewTreeFromFile(*input)
	if err != nil {
		fmt.Println("Error parsing GEDCOM file:", err)
		return
	}
	if otree == nil {
		fmt.Println("Failed to parse GEDCOM file:", *input)
		return
	}
	err = otree.CleanAndExport(*output)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}
