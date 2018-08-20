package loader

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ghodss/yaml"
)

func readPath(path string) ([]byte, error) {
	data, err := chooseReader(path)(path)
	if err != nil {
		return nil, err
	}

	return yaml.YAMLToJSON(data)
}

func chooseReader(path string) func(string) ([]byte, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return readRemoteFile
	}

	return ioutil.ReadFile
}

func readRemoteFile(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
