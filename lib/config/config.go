package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Apns struct {
	AuthKey     string
	KeyID       string
	TeamID      string
	DeviceToken string
}
type Config struct {
	Apns Apns
}

func GetConfig() (tomlConfig Config, err error) {
	var config Config
	if _, err := toml.DecodeFile("startup.toml", &config); err != nil {
		fmt.Println(err)
		return config, err
	}
	return config, nil
}
