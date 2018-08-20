package loader

import (
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

const SchemaPath = "schemas/config.yaml"

func validateWithSchema(data []byte) error {
	// TODO - need bindata or something
	schema, err := loadSchema(SchemaPath)
	if err != nil {
		return err
	}

	loader := gojsonschema.NewStringLoader(string(data))

	result, err := schema.Validate(loader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		errors := []string{}
		for _, err := range result.Errors() {
			errors = append(errors, fmt.Sprintf("\t%s", err.(gojsonschema.ResultError)))
		}

		return fmt.Errorf("invalid schema: \n%s", strings.Join(errors, "\n"))
	}

	return nil
}

func loadSchema(path string) (*gojsonschema.Schema, error) {
	json, err := readPath(path)
	if err != nil {
		return nil, err
	}

	return gojsonschema.NewSchema(gojsonschema.NewBytesLoader(json))
}
