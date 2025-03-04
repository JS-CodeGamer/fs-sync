package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/go-units"
	"gopkg.in/yaml.v3"

	_ "embed"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`

	Database struct {
		Driver            string `yaml:"driver"`
		MaxConnections    int    `yaml:"max_connections"`
		MaxIdleConnection int    `yaml:"max_idle_connections"`
	} `yaml:"database"`

	Auth struct {
		JWTSecret          string        `yaml:"jwt_secret"`
		AccessTokenExpiry  time.Duration `yaml:"access_token_expiry"`
		RefreshTokenExpiry time.Duration `yaml:"refresh_token_expiry"`
	} `yaml:"auth"`

	Storage struct {
		BasePath         string `yaml:"base_path"`
		VersionPath      string `yaml:"versions"`
		ThumbnailPath    string `yaml:"thumbnails"`
		LogsPath         string `yaml:"logs"`
		DbPath           string `yaml:"database"`
		MaxFileSize      string `yaml:"max_file_size"`
		MaxFileSizeBytes int64  // Internal field for parsed size
	} `yaml:"storage"`

	Thumbnails struct {
		EnableService bool `yaml:"enable_service"`
		MaxWidth      int  `yaml:"max_width"`
		MaxHeight     int  `yaml:"max_height"`
	} `yaml:"thumbnails"`
}

//go:embed default_config.yaml
var defaultConfig []byte

var config Config

func InitConfig() error {
	configPath := flag.String("config", "data/config.yaml", "path to config file")
	dbPath := flag.String("db", "", "database path")
	serverHost := flag.String("host", "", "server host")
	serverPort := flag.Int("port", 0, "server port")
	flag.Parse()

	if err := ensureConfigFile(*configPath); err != nil {
		return err
	}

	err := LoadConfig(*configPath)
	if err != nil {
		return err
	}

	if *dbPath != "" {
		config.Storage.DbPath = *dbPath
	}
	if *serverHost != "" {
		config.Server.Host = *serverHost
	}
	if *serverPort != 0 {
		config.Server.Port = *serverPort
	}

	return nil
}

func LoadConfig(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
		return err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
		return err
	}

	// Parse the file size string
	size, err := units.RAMInBytes(config.Storage.MaxFileSize)
	if err != nil {
		log.Fatalf("Error parsing max file size: %v", err)
		return err
	}
	config.Storage.MaxFileSizeBytes = size

	return nil
}

func ensureConfigFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		return os.WriteFile(path, defaultConfig, 0644)
	}
	return nil
}

func GetConfig() *Config {
	return &config
}
