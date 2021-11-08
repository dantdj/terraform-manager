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

func Load() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
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
	_ = ioutil.WriteFile("config.json", data, 0644)
}

func UpdateCurrentVersion(terraformVersion string) {
	Configuration.CurrentVersion = terraformVersion

	data, _ := json.MarshalIndent(Configuration, "", " ")
	_ = ioutil.WriteFile("config.json", data, 0644)
}
