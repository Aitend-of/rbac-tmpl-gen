package main

import (
	"flag"
	"io/ioutil"
	"os"

	tmplgen "github.com/aitend-of/rbac-tmpl-gen"
	mysql "github.com/aitend-of/rbac-tmpl-gen/pkg/mysql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

func main() {
	// Define flags
	inputFilePath := flag.String("input", "./accounts.yaml", "Input YAML file path")
	outputFilePath := flag.String("output", "./rbac_template.yaml", "Output YAML file path")

	// Parse command line flags
	flag.Parse()

	// Read from input file
	inputSpec, err := ioutil.ReadFile(*inputFilePath)
	if err != nil {
		panic(err.Error())
	}

	// Parse input file
	var openAPISpec tmplgen.OpenAPISpec
	yaml.Unmarshal(inputSpec, &openAPISpec)

	// Generate output YAML
	serviceTmpl := openAPISpec.GenerateServiceTmpl()

	// Creating output file
	outputFile, err := os.Create(*outputFilePath)
	if err != nil {
		panic(err.Error())
	}

	// Encoding and writing YAML into new file
	yaml.NewEncoder(outputFile).Encode(serviceTmpl)

	// Closing new file
	err = outputFile.Close()
	if err != nil {
		panic(err.Error())
	}

	// Open DB
	db, err := mysql.Open("root:rootpassword@tcp(127.0.0.1:3306)/rbac_template")
	if err != nil {
		panic(err.Error())
	}

	// Insert output YAML to DB
	mysql.InsertServiceTmpl(db, &serviceTmpl)

	// Close DB
	err = db.Close()
	if err != nil {
		panic(err.Error())
	}
}
