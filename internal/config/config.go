package config

import (
	"cmp"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	InvidiousUrl string `toml:"baseUrl"`
	PlaylistPath string `toml:"playlistPath"`
}

func NewConfig() *Config {
	return &Config{}
}

func defaultConfig() *Config {
	return &Config{
		InvidiousUrl: "https://inv.perditum.com",
		PlaylistPath: "playlists",
	}
}

func (c *Config) Load() error {
	log.Println("Loading config")
	*c = *defaultConfig()

	home, _ := os.UserHomeDir()
	xdgConfHome := cmp.Or(os.Getenv("XDG_CONFIG_HOME"), home+"/.config")

	tomlData, err := os.ReadFile(xdgConfHome + "/yt-music-tui/conf.toml")
	if err != nil {
		return fmt.Errorf("Read file failed: %w", err)
	}
	_, tomlErr := toml.Decode(string(tomlData), &c)

	if tomlErr != nil {
		return fmt.Errorf("Decode Toml failed: %w", err)
	}

	return nil
}
