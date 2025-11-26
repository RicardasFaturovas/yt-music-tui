package config

import (
	"cmp"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	InvidiousUrl string `toml:"baseUrl"`
	PlaylistPath string `toml:"playlistPath"`
}

func NewConfig() *Config {
	log.Println("Loading config")
	home, _ := os.UserHomeDir()
	xdgConfHome := cmp.Or(os.Getenv("XDG_CONFIG_HOME"), home+"/.config")

	tomlData, err := os.ReadFile(xdgConfHome + "/yt-music-tui/conf.toml")
	if err != nil {
		log.Panicln(err)
	}
	var c Config

	_, tomlErr := toml.Decode(string(tomlData), &c)

	if tomlErr != nil {
		log.Panicln("Could not decode config file")

	}

	return &c
}
