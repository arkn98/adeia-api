package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config holds all of the config needed for the application.
type Config struct {
	Server ServerConfig `yaml:"server"`
}

// ServerConfig holds server-specific config.
type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// Load loads YAML from confPath into a new Config.
// confPath must be a readable file containing valid YAML.
func Load(confPath string) (*Config, error) {
	conf := &Config{}

	file, err := os.Open(confPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&conf); err != nil {
		return nil, err
	}

	return conf, nil
}