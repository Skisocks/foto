package config

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/skisocks/foto/pkg/helpers/files"
	"os"
	"path"
	"path/filepath"
)

const (
	configName = "config.json"
)

var questions = []*survey.Question{
	{
		Name:     "RootPhotoDir",
		Prompt:   &survey.Input{Message: "Where is your photos root directory?"},
		Validate: survey.Required,
	},
}

// Config stores all the information that the CLI needs to run
type Config struct {
	// Directories stores the locations of relevant folders
	Directories struct {
		RootPhotoDir string `json:"rootPhotoDirectory"`
	} `json:"directories"`
}

func (c *Config) init() error {
	// Init Dirs
	err := survey.Ask(questions, &c.Directories)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Write(configPath string) error {
	data, err := files.MarshalJSON(c)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0700)
}

func LoadConfigOrCreateIt() (*Config, error) {
	cfg := new(Config)
	configDir, err := GetConfigDirectoryOrCreateIt()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(configDir, configName)

	var data []byte
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		fmt.Println("No config was found. Initialising config.")
		err := cfg.init()
		if err != nil {
			return nil, err
		}
		err = cfg.Write(configPath)
		if err != nil {
			return nil, err
		}
	} else {
		data, err = os.ReadFile(configPath)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(data, cfg)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func GetConfigDirectoryOrCreateIt() (string, error) {
	userHomeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	configDir := path.Join(userHomeDir, ".foto")
	_, err = os.Stat(configDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(configDir, os.FileMode(0700))
		if err != nil {
			return "", err
		}
	}
	return configDir, nil
}
