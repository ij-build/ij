package loader

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

func readFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return readData(data)
}

func readData(data []byte) ([]byte, error) {
	return yaml.YAMLToJSON(data)
}
