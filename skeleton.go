package skeleton

import (
	"github.com/Syfaro/telegram-bot-api"
)

type app struct {
	botAPI       *tgbotapi.BotAPI
	rules        *rules
	pipeline     *Pipeline
	keyboards    *Keyboards
	updateConfig *tgbotapi.UpdateConfig
}

// NewBot
func NewBot(token string) *app {

	// set logger
	tgbotapi.SetLogger(log)

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Printf("Error create Telegram Bot API: %s", err)
		return nil
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	var updateConfig = tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	return &app{
		botAPI:       bot,
		rules:        newRules(),
		pipeline:     newPipeline(),
		keyboards:    newKeyboards(),
		updateConfig: &updateConfig,
	}
}

// Debug ON debug bot
func (a *app) Debug() {
	if a.botAPI == nil {
		return
	}
	a.botAPI.Debug = true
}

// parse message mode default
var parseMode = tgbotapi.ModeHTML

// SetParseMode
func SetParseMode(mode string) {
	parseMode = mode
}

// default message if not found rule
var defaultMessage = "I dont understand you 😔"

// SetDefaultMessage
func SetDefaultMessage(text string) {
	defaultMessage = text
}
