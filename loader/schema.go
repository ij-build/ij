package loader

import (
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"

	"github.com/efritz/ij/schema"
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
			resultError, ok := err.(gojsonschema.ResultError)
			if !ok {
				continue
			}

			switch resultError.(type) {
			case *gojsonschema.NumberOneOfError:
				continue
			case *gojsonschema.NumberAnyOfError:
				continue
			case *gojsonschema.NumberAllOfError:
				continue
			default:
			}

			errors = append(errors, fmt.Sprintf(
				"\t%s\n\t\t%s",
				resultError.Context().String(),
				resultError.Description(),
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
