package config

import (
	"encoding/json"
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

func InitializeConfig() {
	_, err := ioutil.ReadFile("config.json")
	if err != nil {
		// If we failed to open the file, create a new default one
		Configuration.TerraformVersions = []TerraformVersionConfig{}
		Configuration.CurrentVersion = ""
		data, _ := json.MarshalIndent(Configuration, "", " ")
		_ = ioutil.WriteFile("config.json", data, 0644)
	}
}

func Load() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		// If we failed to open the file, create a new default one
		Configuration.TerraformVersions = []TerraformVersionConfig{}
		Configuration.CurrentVersion = ""
		data, _ = json.MarshalIndent(Configuration, "", " ")
		_ = ioutil.WriteFile("config.json", data, 0644)
	}

	err = json.Unmarshal(data, &Configuration)
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

func UpdateCurrentVersion(terraformVersion string) {
	Configuration.CurrentVersion = terraformVersion

	data, _ := json.MarshalIndent(Configuration, "", " ")
	err := ioutil.WriteFile("config.json", data, 0644)
	if err != nil {
		panic(err)
	}
}
