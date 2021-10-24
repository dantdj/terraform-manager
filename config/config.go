package config

import (
	"encoding/json"
	"io/ioutil"
)

var Configuration Config

type Config struct {
	TerraformVersions TerraformVersionConfig `json:"terraformVersions"`
	CurrentVersion    string                 `json:"currentVersion"`
}

type TerraformVersionConfig struct {
	Version    string `json:"version"`
	PathToFile string `json:"pathToFile"`
}

func Load() {
	data, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &Configuration)
	if err != nil {
		panic(err)
	}
}

func AddVersionConfig(terraformVersion, binaryLocation string) {

}
