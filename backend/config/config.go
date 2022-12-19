package config

import (
	"github.com/pelletier/go-toml"
	"io"
)

type DBConfig struct {
	Address  string
	UserName string
	Password string
	DBName   string
}

type MailConfig struct {
	ServerHost string
	ServerPort int
	UserName   string
	Password   string
}

type Config struct {
	DB   map[string]*DBConfig
	Mail map[string]*MailConfig
}

func Read(r io.Reader) *Config {
	tree, err := toml.LoadReader(r)
	if err != nil {
		panic(err)
	}
	var cfg = new(Config)
	err = tree.Unmarshal(cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
