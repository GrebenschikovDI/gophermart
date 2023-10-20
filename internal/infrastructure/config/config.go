package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	RunAddress     string
	Dsn            string
	AccrualAddress string
	Migrations     string
}

const (
	defaultRun        = "localhost:8000"
	defaultDsn        = ""
	defaultAccrual    = ""
	defaultMigrations = "migrations"
)

func LoadConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{}
	err := cfg.configureFlags()
	if err != nil {
		return nil, err
	}
	err = cfg.configureEnv()
	if err != nil {
		return nil, err
	}
	cfg.Migrations = defaultMigrations
	return cfg, nil
}

func (c *ServerConfig) configureFlags() error {
	flag.StringVar(&c.RunAddress, "a", defaultRun, "address and port for server to run")
	flag.StringVar(&c.Dsn, "d", defaultDsn, "database address")
	flag.StringVar(&c.AccrualAddress, "r", defaultAccrual, "accrual address")
	return nil
}

func (c *ServerConfig) configureEnv() error {
	if envRun := os.Getenv("RUN_ADDRESS"); envRun != "" {
		c.RunAddress = envRun
	}
	if envDsn := os.Getenv("DATABASE_URI"); envDsn != "" {
		c.Dsn = envDsn
	}
	if envAccrual := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrual != "" {
		c.AccrualAddress = envAccrual
	}
	return nil
}
