package bot

import (
	"io"
	"time"

	"github.com/linde12/gowol"
	tb "gopkg.in/tucnak/telebot.v2"
	"gopkg.in/yaml.v2"
)

// Device represents the computer to be woken up
type Device struct {
	Name string `yaml:"name"`
	MAC  string `yaml:"mac"`
	IP   string `yaml:"ip"`
}

// Config is the bot configuration
type Config struct {
	BotToken       string   `yaml:"bot-token"`
	AllowedClients []int64  `yaml:"allowed-clients"`
	Devices        []Device `yaml:"devices"`
}

// ConfigFromYAML loads bot configuration from a YAML file
func ConfigFromYAML(input io.Reader) (*Config, error) {
	yamlDec := yaml.NewDecoder(input)

	var botCfg = &Config{}
	err := yamlDec.Decode(botCfg)
	return botCfg, err
}

// Bot is the Telegram bot
type Bot struct {
	cfg *Config
	bot *tb.Bot
}

// NewBot returns a new instance of Bot
func NewBot(cfg *Config) (*Bot, error) {
	bot, err := tb.NewBot(tb.Settings{
		Token:  cfg.BotToken,
		Poller: &tb.LongPoller{Timeout: 30 * time.Second},
	})

	if err != nil {
		return nil, err
	}

	return &Bot{
		cfg: cfg,
		bot: bot,
	}, nil
}

func (b *Bot) clientIsAllowed(chatID int64) bool {
	for _, id := range b.cfg.AllowedClients {
		if id == chatID {
			return true
		}
	}
	return false
}

func (b *Bot) getDeviceByName(name string) *Device {
	for _, d := range b.cfg.Devices {
		if d.Name == name {
			return &d
		}
	}
	return nil
}

// Start the bot
func (b *Bot) Start() {
	b.bot.Handle("/wol", func(m *tb.Message) {
		if !b.clientIsAllowed(m.Chat.ID) {
			b.bot.Send(m.Sender, "request not allowed")
			return
		}

		dev := b.getDeviceByName(m.Payload)
		if dev == nil {
			b.bot.Send(m.Sender, "unknown device '"+m.Payload+"'")
			return
		}

		packet, err := gowol.NewMagicPacket(dev.MAC)
		if err != nil {
			b.bot.Send(m.Sender, "error creating magic packet: "+err.Error())
			return
		}

		err = packet.Send(dev.IP)
		if err != nil {
			b.bot.Send(m.Sender, "error sending magic packet: "+err.Error())
			return
		}
		b.bot.Send(m.Sender, "magic packet sent!")
	})

	b.bot.Start()
}
