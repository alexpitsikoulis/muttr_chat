package app

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port   uint         `yaml:"port"`
	DB     DBConfig     `yaml:"db"`
	Server ServerConfig `yaml:"server"`
}

type DBConfig struct {
	Dialect  string `yaml:"dialect"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port uint   `yaml:"port"`
}

func NewConfig(filePath string) (*Config, error) {
	var cfg Config
	cfgFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(cfgFile, &cfg)
	return &cfg, err
}
