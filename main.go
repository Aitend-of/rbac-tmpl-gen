package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"

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

	// Generate output YAML
	serviceTmpl.ServiceName = openAPISpec.Info.Title // Service name

	for pathKey, path := range openAPISpec.Paths {
		for methodKey, method := range path {
			for _, feature := range method.RBACFeature {
				var endpoints []map[string]string
				endpoint := make(map[string]string)
				endpoint[pathKey] = methodKey
				endpoints = append(endpoints, endpoint)

				serviceTmpl.Features = append(serviceTmpl.Features, Feature{
					FeatureName: strings.ReplaceAll(feature, "_", " "),
					ID:          feature,
					Description: method.Description,
					Endpoints:   endpoints,
				})
			}
		}
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
