package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Database struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	URL  string `yaml:"url"`
}

type Server struct {
	Swagger   bool `yaml:"swagger"`
	Websocket bool `yaml:"websocket"`
	Port      int  `yaml:"port"`
}

type Device struct {
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
	Token    string `yaml:"token"`
}

type Config struct {
	Database Database `yaml:"database"`
	Server   Server   `yaml:"server"`
	Devices  []Device `yaml:"devices"`
}

func Load(path string) (*Config, error) {
	bytes, err := os.ReadFile(fmt.Sprintf("%s/config.yml", path))
	if err != nil {
		return nil, fmt.Errorf("error reading config file by path %s: %w", path, err)
	}

	config := &Config{}
	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file by path %s: %w", path, err)
	}

	return config, nil
}
