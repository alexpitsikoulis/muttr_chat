package app

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port uint     `yaml:"port"`
	DB   DBConfig `yaml:"db"`
}

type DBConfig struct {
	Dialect  string `yaml:"dialect"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func NewConfig(filePath string) (*Config, error) {
	var cfg Config
	cfgFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(cfgFile, &cfg)
	return &cfg, nil
}
