package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`

	Database struct {
		Driver            string `yaml:"driver"`
		Path              string `yaml:"path"`
		MaxConnections    int    `yaml:"max_connections"`
		MaxIdleConnection int    `yaml:"max_idle_connections"`
	} `yaml:"database"`

	Auth struct {
		JWTSecret   string `yaml:"jwt_secret"`
		TokenExpiry string `yaml:"token_expiry"`
	} `yaml:"auth"`

	Storage struct {
		UploadPath  string `yaml:"upload_path"`
		MaxFileSize int    `yaml:"max_file_size"`
	} `yaml:"storage"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
		return nil, err
	}

	return &config, nil
}
