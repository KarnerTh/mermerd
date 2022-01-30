package config

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

const fileMode = os.FileMode(0600)

type FileConfig struct {
	ConnectionStrings []string `yaml:"connectionStrings"`
}

func initialConfig() FileConfig {
	return FileConfig{
		ConnectionStrings: []string{
			"postgresql://user:password@localhost:5432/yourDb",
			"mysql://root:password@tcp(127.0.0.1:3306)/yourDb",
		},
	}
}

func LoadConfigFile() error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Could not get home directory", err.Error())
		return err
	}

	configFilePath := filepath.FromSlash(fmt.Sprintf("%s/.mermerd", dirname))
	file, err := os.OpenFile(configFilePath, os.O_RDONLY, fileMode)
	if errors.Is(err, os.ErrNotExist) {
		file, err = createInitialConfig(configFilePath)
	}

	if err != nil {
		fmt.Println("Could not open file", err.Error())
		return err
	}

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Could not read config file data", err.Error())
		return err
	}

	var config FileConfig
	err = yaml.Unmarshal(fileContent, &config)
	if err != nil {
		fmt.Println("Error reading content of config file", err.Error())
		return err
	}

	ConnectionStringSuggestions = config.ConnectionStrings
	return nil
}

func createInitialConfig(configFilePath string) (*os.File, error) {
	color.Yellow("Config file does not exist. Creating a new one in %s", configFilePath)
	file, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileMode)
	configValue, err := yaml.Marshal(initialConfig())
	if err != nil {
		fmt.Println("could not create initial config", err.Error())
		return nil, err
	}

	_, err = file.Write(configValue)
	if err != nil {
		fmt.Println("could not write to config file", err.Error())
		return nil, err
	}

	err = file.Sync()
	if err != nil {
		fmt.Println("could not save initial config", err.Error())
		return nil, err
	}

	return file, nil
}
