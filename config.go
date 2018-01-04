package runtelchat

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

//Config represents the configuration of the chat server
type Config struct {
	Port        string `json:"port"`
	Host        string   `json:"host"`
	LogFilePath string   `json:"logFilePath"`
}

var defaultConfig = Config{
	Port:       "2300",
	Host:        "localhost",
	LogFilePath: "./messagelog.log",
}

//LoadConfig loads the config file from the current directory
//or the default config if no config file exists
func LoadConfig() (Config, error) {
	if _, err := os.Stat("./config.json"); os.IsNotExist(err) {
		return defaultConfig, nil
	}
	configFileBytes, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return Config{}, errors.Wrap(err, "Config file exists but cannot be read")
	}
	var config Config
	err = json.Unmarshal(configFileBytes, &config)
	if err != nil {
		return Config{}, errors.Wrap(err, "Config file format is invalid")
	}
	return config, nil
}
