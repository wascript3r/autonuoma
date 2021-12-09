package main

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

const ConfigENV = "AUTONUOMA_CONFIG"

var (
	ErrConfigNotProvided = errors.New("config file is not provided")
)

type Config struct {
	Log struct {
		ShowTimestamp bool `json:"showTimestamp"`
	} `json:"log"`

	Database struct {
		Postgres struct {
			DSN          string   `json:"dsn"`
			QueryTimeout Duration `json:"queryTimeout"`
		} `json:"postgres"`
	} `json:"database"`

	Cipher struct {
		AES struct {
			Key string `json:"key"`
		} `json:"aes"`
	} `json:"cipher"`

	Auth struct {
		PasswordCost int `json:"passwordCost"`
		Session      struct {
			SessionLifetime Duration `json:"sessionLifetime"`
			CookieName      string   `json:"cookieName"`
			CookieLifetime  Duration `json:"cookieLifetime"`
			SecureCookie    bool     `json:"secureCookie"`
		} `json:"session"`
	} `json:"auth"`

	HTTP struct {
		Port string `json:"port"`
		CORS struct {
			Origin string `json:"origin"`
		} `json:"cors"`
		EnablePprof bool `json:"enablePprof"`
	} `json:"http"`

	WebSocket struct {
		Port         string   `json:"port"`
		ConnIdleTime Duration `json:"connIdleTime"`
	} `json:"webSocket"`
}

func getConfigPath() (string, error) {
	path := os.Getenv(ConfigENV)
	path = strings.TrimSpace(path)

	if len(path) == 0 {
		return "", ErrConfigNotProvided
	}

	return path, nil
}

func parseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = json.NewDecoder(file).Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
