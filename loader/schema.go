package loader

import (
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"

	"github.com/efritz/pvc/schema"
)

func validateWithSchema(data []byte) error {
	schema, err := getSchema()
	if err != nil {
		return err
	}

	result, err := schema.Validate(gojsonschema.NewStringLoader(string(data)))
	if err != nil {
		return err
	}

	if !result.Valid() {
		errors := []string{}
		for _, err := range result.Errors() {
			errors = append(errors, fmt.Sprintf(
				"\t%s",
				err.(gojsonschema.ResultError),
			))
		}

		return fmt.Errorf(
			"invalid schema: \n%s",
			strings.Join(errors, "\n"),
		)
	}

	return nil
}

func getSchema() (*gojsonschema.Schema, error) {
	data, err := schema.Asset("config.yaml")
	if err != nil {
		return nil, err
	}

	json, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}

	return gojsonschema.NewSchema(gojsonschema.NewBytesLoader(json))
}
