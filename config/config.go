package config

import (
	"io"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    server
	Log       log
	WebSocket websocket
	Generator generator
}

type generator struct {
	HypoxiaMode int `yaml:"hypoxia_mode" envconfig:"HYPOXIA_MODE"`
}

type log struct {
	Level string `yaml:"level"`
}

type server struct {
	Addr        string `yaml:"addr" envconfig:"SERVER_ADDR"`
	Port        string `yaml:"port" envconfig:"SERVER_PORT"`
	SensorID    string `yaml:"sensor_id" envconfig:"SENSOR_ID"`
	SensorToken string `yaml:"sensor_token" envconfig:"SENSOR_TOKEN"`
}

type websocket struct {
	Addr string `yaml:"addr" envconfig:"WEBSOCKET_ADDR"`
	Port string `yaml:"port" envconfig:"WEBSOCKET_PORT"`
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
