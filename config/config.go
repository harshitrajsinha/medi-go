package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type localDB struct {
	User string `envconfig:"DB_USER"`
	Host string `envconfig:"DB_HOST"`
	Port string `envconfig:"DB_PORT"`
	Pass string `envconfig:"DB_PASS"`
	Name string `envconfig:"DB_NAME"`
}

type redisCloud struct {
	Host string `envconfig:"REDIS_HOST"`
	Port string `envconfig:"REDIS_PORT"`
	Pass string `envconfig:"REDIS_PASS"`
}

type neonDB struct {
	NeonConnStr string `envconfig:"NEON_CONNSTR"`
}

type origin struct {
	AllowedOrigin string `envconfig:"ALLOWED_ORIGIN"`
}

// helper to avoid repetition
func loadConfig[T any](cfg *T, desc string) error {
	if err := envconfig.Process("", cfg); err != nil { // load env from program's environment to declared struct
		return fmt.Errorf("failed to load config for %s: %w", desc, err)
	}
	return nil
}

func DBConfig() (*localDB, error) {
	var c localDB
	return &c, loadConfig(&c, "local database env")
}

func NeonDBConfig() (*neonDB, error) {
	var c neonDB
	return &c, loadConfig(&c, "neon db connection string")
}

func RedisConfig() (*redisCloud, error) {
	var c redisCloud
	return &c, loadConfig(&c, "redis cloud env")
}

func AllowedOrigin() (*origin, error) {
	var c origin
	return &c, loadConfig(&c, "origin website")
}
