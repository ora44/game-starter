package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/ktr0731/go-fuzzyfinder"
)

type Game struct {
	Name    string   `toml:"name"`
	Cmd     string   `toml:"cmd"`
	WorkDir string   `toml:"workdir"`
	Args    []string `toml:"args"`
}

type Config struct {
	Games []Game `toml:"games"`
}

func main() {
	usrconf, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	var cfg Config
	_, err = toml.DecodeFile(path.Join(usrconf, "game-starter", "config.toml"), &cfg)
	if err != nil {
		panic(err)
	}

	idx, err := fuzzyfinder.Find(cfg.Games, func(i int) string {
		return cfg.Games[i].Name
	}, fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
		if i == -1 {
			return ""
		}
		return fmt.Sprintf("will launch %s", cfg.Games[i].Name)
	}))
	if err != nil {
		panic(err)
	}
	fmt.Printf("idx: %v\n", idx)

	cmdline := strings.Split(cfg.Games[idx].Cmd, " ")
	cmd := exec.Command(cmdline[0], cmdline[1:]...)
	cmd.Run()
}
