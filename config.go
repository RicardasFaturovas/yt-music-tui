package main

import (
	"cmp"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	InvidiousUrl string `toml:"baseUrl"`
}

var conf *Config

func loadConfig() {
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
	log.Println(c)
	conf = &c
}

func getConfig() *Config {
	if conf == nil {
		panic("getConfig called before loadConfig()")
	}
	return conf
}
