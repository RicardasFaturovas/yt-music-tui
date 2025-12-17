package config

import (
	"cmp"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gdamore/tcell/v2"
)

type Config struct {
	InvidiousUrl string `toml:"baseUrl"`
	PlaylistPath string `toml:"playlistPath"`
	Theme        Theme  `toml:"theme"`
}

type Theme struct {
	BackgroundColor    Color `toml:"primaryBackgroundColor"`
	PrimaryTextColor   Color `toml:"primaryTextColor"`
	BorderColor        Color `toml:"borderColor"`
	TitleColor         Color `toml:"titleColor"`
	SecondaryTextColor Color `toml:"secondaryTextColor"`
}

type Color tcell.Color

func (c *Color) UnmarshalText(text []byte) error {
	s := strings.TrimPrefix(string(text), "#")
	if len(s) != 6 {
		return fmt.Errorf("invalid color format: %q", text)
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	*c = Color(tcell.NewRGBColor(
		int32(b[0]), int32(b[1]), int32(b[2]),
	))
	return nil
}

func NewConfig() *Config {
	return &Config{}
}

func defaultConfig() *Config {
	return &Config{
		InvidiousUrl: "https://inv.perditum.com",
		PlaylistPath: "playlists",
		Theme: Theme{
			BackgroundColor:    Color(tcell.NewHexColor(0x2a273f)),
			PrimaryTextColor:   Color(tcell.NewHexColor(0xe0def4)),
			BorderColor:        Color(tcell.NewHexColor(0x3e8fb0)),
			TitleColor:         Color(tcell.NewHexColor(0xea9a97)),
			SecondaryTextColor: Color(tcell.NewHexColor(0x9ccfd8)),
		},
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
