package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var Configuration Config

type Config struct {
	TerraformVersions []TerraformVersionConfig `json:"terraformVersions"`
	CurrentVersion    string                   `json:"currentVersion"`
}

type TerraformVersionConfig struct {
	Version    string `json:"version"`
	PathToFile string `json:"pathToFile"`
}

// Reads the config and returns the data wihin.
func InitializeConfig() []byte {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		// If we failed to open the file, create a new default one
		Configuration.TerraformVersions = []TerraformVersionConfig{}
		Configuration.CurrentVersion = ""
		data, _ = json.MarshalIndent(Configuration, "", " ")
		_ = ioutil.WriteFile("config.json", data, 0644)
	}

	return data
}

func Load() {
	data := InitializeConfig()

	err := json.Unmarshal(data, &Configuration)
	if err != nil {
		panic(err)
	}
}

func AddVersionConfig(terraformVersion, binaryLocation string) {
	Configuration.TerraformVersions = append(Configuration.TerraformVersions, TerraformVersionConfig{
		Version:    terraformVersion,
		PathToFile: binaryLocation,
	})

	data, _ := json.MarshalIndent(Configuration, "", " ")
	err := ioutil.WriteFile("config.json", data, 0644)
	if err != nil {
		panic(err)
	}
}

func GetCurrentVersion() (TerraformVersionConfig, error) {
	if Configuration.CurrentVersion == "" {
		panic(fmt.Errorf("current version is not set"))
	}

	for _, versionConfig := range Configuration.TerraformVersions {
		if versionConfig.Version == Configuration.CurrentVersion {
			return versionConfig, nil
		}
	}

	return TerraformVersionConfig{}, fmt.Errorf("the currentVersion field does not match any known version in config")
}

func UpdateCurrentVersion(terraformVersion string) {
	Configuration.CurrentVersion = terraformVersion

	data, _ := json.MarshalIndent(Configuration, "", " ")
	err := ioutil.WriteFile("config.json", data, 0644)
	if err != nil {
		panic(err)
	}
}
