package config

import (
	"github.com/BurntSushi/toml"
)

const (
	version string = "0.1"
)

type Config struct {
	Addr      string `toml:"addr"`
	Port      int    `toml:"port"`
	Database  string `toml:"database"`
	ImagePath string `toml:"image_path"`
	Version   string
}

func GetConfig(filename string) (Config, error) {
	var conf Config

	_, err := toml.DecodeFile(filename, &conf)
	if err != nil {
		return conf, err
	}
	conf.Version = version
	return conf, nil
}
