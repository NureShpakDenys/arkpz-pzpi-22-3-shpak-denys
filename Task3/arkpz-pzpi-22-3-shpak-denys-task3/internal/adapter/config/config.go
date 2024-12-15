package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	Http        HttpConfig `yaml:"http"`
	AuthConfig  AuthConfig `yaml:"auth"`
	DBPassword  string     `yaml:"db_password" env-required:"true"`
}

type HttpConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timout"`
}

type AuthConfig struct {
	SecretKey   string        `yaml:"secret"`
	TokenExpiry time.Duration `yaml:"token_ttl"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
