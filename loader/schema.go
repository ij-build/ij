package loader

import (
	"fmt"
	"os"

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
		for _, err := range result.Errors() {
			fmt.Printf("Validation error: %#v\n", err)
		}

		// TODO
		os.Exit(1)
	}

	return nil
}

func loadSchema(path string) (*gojsonschema.Schema, error) {
	json, err := readFile(path)
	if err != nil {
		return nil, err
	}

	return gojsonschema.NewSchema(gojsonschema.NewBytesLoader(json))
}
