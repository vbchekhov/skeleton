package skeleton

import "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Context command to execution
type Context struct {
	// app - application storage
	app *app
	// chatId - current user telegram id
	chatId int64
	// rule - current execution rule
	rule *Rule
	// BotAPI - bot api (github.com/go-telegram-bot-api/telegram-bot-api )
	BotAPI *tgbotapi.BotAPI
	// Update - update from telegram bot api
	Update *tgbotapi.Update
	// RegexpResult - result exec regex
	RegexpResult []string
}

// app current application
func (c *Context) App() *app {
	return c.app
}

// ChatId current user telegram id
func (c *Context) ChatId() int64 {
	return c.chatId
}

// Keyboards current keyboard state
func (c *Context) Keyboards() *Keyboards {
	return c.app.keyboards
}

// Pipeline current pipeline storage
func (c *Context) Pipeline() *Pipeline {
	// set current rule and chat id
	c.app.pipeline.chatId = c.chatId
	c.app.pipeline.rule = c.rule

	return c.app.pipeline
}

// AllowList current app
func (c *Context) AllowList() *AllowList {
	return c.app.allowList
}
