package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	// Define flags
	inputFilePath := flag.String("input", "./accounts.yaml", "Input YAML file path")
	outputFilePath := flag.String("output", "./rbac_template.yaml", "Output YAML file path")

	// Parse command line flags
	flag.Parse()

	fmt.Printf("Input path: %s\nOutput path: %s\n", *inputFilePath, *outputFilePath)

	// Reading from file
	inputSpec, err := ioutil.ReadFile(*inputFilePath)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Input Open API spec:\n%s", string(inputSpec))

	// Creating new file
	outputFile, err := os.Create(*outputFilePath)
	if err != nil {
		panic(err.Error())
	}

	// Writing into new file
	_, err = outputFile.WriteString(string(inputSpec))
	if err != nil {
		panic(err.Error())
	}

	// Closing new file
	err = outputFile.Close()
	if err != nil {
		panic(err.Error())
	}
}
