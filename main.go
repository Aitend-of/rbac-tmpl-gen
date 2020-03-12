package main

import (
	"flag"
	"fmt"
)

func main() {
	// Define flags
	inputFilePath := flag.String("input", "./accounts.yaml", "Input YAML file path")
	outputFilePath := flag.String("output", "./rbac_template.yaml", "Output YAML file path")

	//Parse command line flags
	flag.Parse()

	fmt.Printf("Input path: %s\nOutput path: %s", *inputFilePath, *outputFilePath)
}
