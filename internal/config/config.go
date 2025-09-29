package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type serverCfg struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type loggingCfg struct {
	Level string `yaml:"level"`
}

type Config struct {
	Server  serverCfg  `yaml:"server"`
	Logging loggingCfg `yaml:"logging"`
}

// Load loads configuration from configs/config.yaml if present,
// otherwise returns default configuration.
func Load() *Config {
	cfg := &Config{
		Server: serverCfg{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Logging: loggingCfg{
			Level: "info",
		},
	}

	data, err := ioutil.ReadFile("configs/config.yaml")
	if err != nil {
		// if file not found, return defaults
		return cfg
	}

	var fileCfg Config
	if err := yaml.Unmarshal(data, &fileCfg); err != nil {
		return cfg
	}

	if fileCfg.Server.Host != "" {
		cfg.Server.Host = fileCfg.Server.Host
	}
	if fileCfg.Server.Port != 0 {
		cfg.Server.Port = fileCfg.Server.Port
	}
	if fileCfg.Logging.Level != "" {
		cfg.Logging.Level = fileCfg.Logging.Level
	}

	return cfg
}

func (c *Config) LogLevel() string {
	return c.Logging.Level
}

func (c *Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
