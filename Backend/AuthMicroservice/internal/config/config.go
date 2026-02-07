package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string         `env:"ENV" env-default:"local"`
	Database       DatabaseConfig `env-prefix:"DB_"`
	GRPC           GRPCConfig     `env-prefix:"GRPC_"`
	AuthCode       AuthCodeConfig `env-prefix:"AUTH_CODE_"`
	Redis          RedisConfig    `env-prefix:"REDIS_"`
	MigrationsPath string         `env:"MIGRATIONS_PATH" env-default:"./migrations"`
}

type GRPCConfig struct {
	Port    int `env:"PORT" env-default:"5000"`
	Timeout int `env:"TIMEOUT" env-default:"10"`
}

type AuthCodeConfig struct {
	Length int           `env:"LENGTH" env-default:"6"`
	TTL    time.Duration `env:"TTL" env-default:"300"`
}

type DatabaseConfig struct {
	URL     string `env:"URL" env-required:"true"`
	PoolMax int    `env:"POOL_MAX" env-default:"5"`
}

type RedisConfig struct {
	Address  string `env-required:"true"    env:"ADDRESS"`
	Password string `env-required:"true"    env:"PASSWORD"`
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
