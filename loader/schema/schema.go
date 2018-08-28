package schema

import (
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"

	"github.com/efritz/ij/asset"
)

func Validate(name string, data []byte) error {
	schema, err := getSchema(name)
	if err != nil {
		return err
	}

	result, err := schema.Validate(gojsonschema.NewStringLoader(string(data)))
	if err != nil {
		return err
	}

	if !result.Valid() {
		for _, err := range result.Errors() {
			resultError, ok := err.(gojsonschema.ResultError)
			if !ok {
				continue
			}

			return fmt.Errorf(resultError.Description())
		}

		return fmt.Errorf("invalid data")
	}

	return nil
}

func getSchema(name string) (*gojsonschema.Schema, error) {
	data, err := asset.Asset(name)
	if err != nil {
		return nil, err
	}

	json, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}

	return gojsonschema.NewSchema(gojsonschema.NewBytesLoader(json))
}
