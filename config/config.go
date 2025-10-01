package config

import (
	"io"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server server
	Log    log
}

type log struct {
	Level string `yaml:"level"`
}

type server struct {
	Addr        string `yaml:"addr"`
	Port        string `yaml:"port"`
	SensorID    string `envconfig:"SENSOR_ID" yaml:"sensor_id"`
	SensorToken string `envconfig:"SENSOR_TOKEN" yaml:"sensor_token"`
}

func ReadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// Читаем переменные окружения
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}

	return cfg, err
}
