package tmplgen

import "strings" // Is it ok to use standart packages here?

// Model of input YAML

type Info struct {
	Title       string `yaml:"title"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

type Server struct {
	URL string `yaml:"url"`
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

// Model of output YAML

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

// Iterate input YAML struct
func (inputSpec *OpenAPISpec) IterateInputSpec(action func(pathKey string, path map[string]Method, methodKey string, method Method, feature string)) {
	for pathKey, path := range inputSpec.Paths {
		for methodKey, method := range path {
			for _, feature := range method.RBACFeature {
				action(pathKey, path, methodKey, method, feature)
			}
		}
	}
}

// Generate output YAML struct
func (inputSpec *OpenAPISpec) GenerateServiceTmpl() ServiceTmpl {
	var serviceTmpl ServiceTmpl

	serviceTmpl.ServiceName = inputSpec.Info.Title

	GenerateServiceFeatures := func(pathKey string, path map[string]Method, methodKey string, method Method, feature string) {
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
	}

	inputSpec.IterateInputSpec(GenerateServiceFeatures)

	return serviceTmpl
}
