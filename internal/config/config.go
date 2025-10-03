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

type databaseCfg struct { // <--- НОВАЯ СТРУКТУРА
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type Config struct {
	Server   serverCfg   `yaml:"server"`
	Logging  loggingCfg  `yaml:"logging"`
	Database databaseCfg `yaml:"database"` // <--- ДОБАВЛЕНО
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
		Database: databaseCfg{ // <--- ЗНАЧЕНИЯ ПО УМОЛЧАНИЮ
			Host:     "localhost",
			Port:     5432,
			User:     "user",
			Password: "password",
			DBName:   "procrastigo_db",
			SSLMode:  "disable",
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

	// ⚙️ Обновляем настройки БД, если они есть в файле
	if fileCfg.Database.Host != "" {
		cfg.Database.Host = fileCfg.Database.Host
	}
	if fileCfg.Database.Port != 0 {
		cfg.Database.Port = fileCfg.Database.Port
	}
	if fileCfg.Database.User != "" {
		cfg.Database.User = fileCfg.Database.User
	}
	if fileCfg.Database.Password != "" {
		cfg.Database.Password = fileCfg.Database.Password
	}
	if fileCfg.Database.DBName != "" {
		cfg.Database.DBName = fileCfg.Database.DBName
	}
	if fileCfg.Database.SSLMode != "" {
		cfg.Database.SSLMode = fileCfg.Database.SSLMode
	}

	return cfg
}

func (c *Config) LogLevel() string {
	return c.Logging.Level
}

func (c *Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func (c *Config) DatabaseDSN() string { // <--- НОВЫЙ МЕТОД
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.DBName, c.Database.SSLMode)
}
