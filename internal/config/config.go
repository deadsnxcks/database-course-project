package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env 	string 		`yaml:"env" env-default:"local"`
	GRPC 	GRPCConfig	`yaml:"grpc"`
}

type GRPCConfig struct {
	Port 	int				`yaml:"port"`
	Timeout time.Duration	`yaml:"timeout"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("config path is empty")
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg) 
	if err != nil {
		panic("failed read config: " + err.Error())
	}

	return &cfg
}