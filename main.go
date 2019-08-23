package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/nohajc/wol-tgbot-go/bot"
)

// Device represents the computer to be woken up
type Device struct {
	Name string `yaml:"name"`
	MAC  string `yaml:"mac"`
}

// Config is the bot configuration
type Config struct {
	AllowedClients []int64  `yaml:"allowed-clients"`
	Devices        []Device `yaml:"devices"`
}

func main() {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Cannot find home directory: %v", err)
	}

	cfgPath := filepath.Join(home, ".wol-tgbot-go", "config.yaml")
	cfgFile, err := os.Open(cfgPath)
	if err != nil {
		log.Fatalf("Cannot load configuration: %v", err)
	}
	defer cfgFile.Close()

	botCfg, err := bot.ConfigFromYAML(cfgFile)
	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	tgBot, err := bot.NewBot(botCfg)
	if err != nil {
		log.Fatalf("Cannot create bot: %v", err)
	}

	tgBot.Start()
}
