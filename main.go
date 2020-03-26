package main

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

// Models of input YAML

type Info struct {
	Title       string `yaml:"title"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

type Server struct {
	URL string "json:`url`"
}

type Method struct {
	Description string                 `yaml:"description"`
	OperationID string                 `yaml:"operationId"`
	Responses   map[string]interface{} `yaml:"responses"`
	RBACFeature []string               `yaml:"x-rbac-feature"`
}

type OpenAPISpec struct {
	Version string                       `yaml:"openapi"`
	Info    Info                         `yaml:"info"`
	Servers []Server                     `yaml:"servers"`
	Paths   map[string]map[string]Method `yaml:"paths"`
}

// Models of output YAML

type Feature struct {
	FeatureName string              `yaml:"feature"`
	ID          string              `yaml:"id"`
	Description string              `yaml:"description"`
	Endpoints   []map[string]string `yaml:"endpoints"`
}

type ServiceTmpl struct {
	ServiceName string `yaml:"serviceName"`
	Features    []Feature
}

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
	var openAPISpec OpenAPISpec
	var serviceTmpl ServiceTmpl

	yaml.Unmarshal(inputSpec, &openAPISpec)

	// Open DB and Prepare request
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/rbac_template")
	if err != nil {
		panic(err.Error())
	}

	stmt, err := db.Prepare("INSERT INTO features(feature_id, feature_name, feature_descr, service_name, endpoint_path, endpoint_method) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}

	// Generate output YAML
	serviceTmpl.ServiceName = openAPISpec.Info.Title // Service name

	for pathKey, path := range openAPISpec.Paths {
		for methodKey, method := range path {
			for _, feature := range method.RBACFeature {
				var endpoints []map[string]string
				endpoint := make(map[string]string)
				endpoint[pathKey] = methodKey
				endpoints = append(endpoints, endpoint)

				currFeature := Feature{
					FeatureName: strings.ReplaceAll(feature, "_", " "),
					ID:          feature,
					Description: method.Description,
					Endpoints:   endpoints,
				}

				serviceTmpl.Features = append(serviceTmpl.Features, currFeature)

				// Write output YAML to DB features table
				_, err := stmt.Exec(currFeature.ID, currFeature.FeatureName, currFeature.Description, serviceTmpl.ServiceName, pathKey, methodKey)
				if err != nil {
					panic(err.Error())
				}

			}
		}
	}

	// Close DB
	err = db.Close()
	if err != nil {
		panic(err.Error())
	}

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
}
