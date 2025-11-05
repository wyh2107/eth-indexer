package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	RPCURL         string `yaml:"rpc_url"`
	ChainIdTimeout time.Duration
	DBUser         string
	DBPassword     string
	DBHost         string
	DBPort         string
	DBName         string

	StartBlock uint64
}

const (
	defaultStartBlock     = uint64(18_900_000)
	defaultChainIdTimeout = time.Second * 5
)

func LoadConfig() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		RPCURL:         os.Getenv("ETH_RPC_URL"),
		ChainIdTimeout: defaultChainIdTimeout,

		DBUser:     os.Getenv("MYSQL_USER"),
		DBPassword: os.Getenv("MYSQL_PASSWORD"),
		DBHost:     os.Getenv("MYSQL_HOST"),
		DBPort:     os.Getenv("MYSQL_PORT"),
		DBName:     os.Getenv("MYSQL_NAME"),
		StartBlock: defaultStartBlock,
	}

	// 可选: START_BLOCK
	if v := os.Getenv("START_BLOCK"); v != "" {
		if u, err := strconv.ParseUint(v, 10, 64); err == nil {
			cfg.StartBlock = u
		} else {
			return Config{}, errors.New("invalid START_BLOCK (must be uint64)")
		}
	}

	// 可选: CHAIN_ID_TIMEOUT_MS
	if v := os.Getenv("CHAIN_ID_TIMEOUT_MS"); v != "" {
		if ms, err := strconv.Atoi(v); err == nil && ms > 0 {
			cfg.ChainIdTimeout = time.Duration(ms) * time.Millisecond
		}
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil

}

func (c Config) Validate() error {
	if c.RPCURL == "" {
		return errors.New("ETH_RPC_URL not set")
	}
	// 也可在此强校验 MySQL 字段：
	// if c.DBUser == "" || c.DBHost == "" || c.DBPort == "" || c.DBName == "" {
	// 	return errors.New("mysql config incomplete")
	// }
	return nil
}
