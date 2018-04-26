package nagus

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type NagusConfig struct {
	ApiKey string `json:"api_key"`
}

type ConfigReadError struct {
	Message string
}

func (e ConfigReadError) Error() string {
	return e.Message
}

// Read a config file at path
func ReadConfig(path string) (NagusConfig, error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return NagusConfig{
			ApiKey: "err",
		}, ConfigReadError{
			Message: fmt.Sprintf("File %s does not exist!", path),
		}
	}

	config := NagusConfig{
		ApiKey: "err",
	}
	err = json.Unmarshal(fileBytes, &config)
	return config, err
}
