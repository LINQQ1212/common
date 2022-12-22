package utils

import (
	"gopkg.in/yaml.v3"
	"os"
)

func ReadYamlConfig(f string, v interface{}) error {
	fp, err := os.Open(f)
	if err != nil {
		return err
	}
	return yaml.NewDecoder(fp).Decode(v)
}

func SaveYamlConfig(f string, v interface{}) error {
	fp, err := os.OpenFile(f, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0766)
	if err != nil {
		return err
	}
	return yaml.NewEncoder(fp).Encode(v)
}
