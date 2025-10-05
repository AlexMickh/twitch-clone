package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env    string       `yaml:"env" env-default:"prod"`
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
	Redis  RedisConfig  `yaml:"redis"`
	Mail   MailConfig   `yaml:"mail"`
}

type ServerConfig struct {
	Addr        string        `yaml:"addr" env-default:"0.0.0.0:50070"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	Session     SessionConfig `yaml:"session"`
}

type DBConfig struct {
	Host        string            `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port        int               `yaml:"port" env:"DB_PORT" env-default:"27017"`
	User        string            `yaml:"user" env:"DB_USER" env-default:"mongo"`
	Password    string            `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	Database    string            `yaml:"database" env:"DB_NAME" env-default:"users"`
	Collections map[string]string `yaml:"collections" env-default:"users"`
}

type RedisConfig struct {
	Host       string        `env:"REDIS_HOST" yaml:"host" env-default:"localhost"`
	Port       int           `env:"REDIS_PORT" yaml:"port" env-default:"6379"`
	User       string        `env:"REDIS_USER" yaml:"user" env-default:"root"`
	Password   string        `env:"REDIS_PASSWORD" yaml:"password" env-default:"root"`
	DB         int           `env:"REDIS_DB" yaml:"db" env-default:"0"`
	Expiration time.Duration `env:"REDIS_EXPIRATION" yaml:"expire_time" env-default:"24h"`
}

type MailConfig struct {
	Host     string `env:"MAIL_HOST" yaml:"host" env-required:"true"`
	Port     int    `env:"MAIL_PORT" yaml:"port" env-required:"true"`
	FromAddr string `env:"MAIL_FROM_ADDR" yaml:"from_addr" env-required:"true"`
	Password string `env:"MAIL_PASSWORD" yaml:"password" env-required:"true"`
}

type SessionConfig struct {
	Name     string `yaml:"name" env-default:"session_id"`
	HttpOnly bool   `yaml:"http_only" env-default:"true"`
	Secure   bool   `yaml:"secure" env-default:"true"`
	MaxAge   int    `yaml:"max_age" env-default:"432000"`
}

func MustLoad() *Config {
	path := fetchPath()
	cfg, err := Load(path)
	if err != nil {
		panic(err)
	}
	return cfg
}

func Load(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", path)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

func fetchPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
