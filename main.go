package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/nohajc/wol-tgbot-go/bot"
)

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
