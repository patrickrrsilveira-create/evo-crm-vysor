package jsonloader

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func LoadJSONFromFile[T any](path string) (*T, error) {

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filepath.Join(wd, path))
	if err != nil {
		return nil, err
	}

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
