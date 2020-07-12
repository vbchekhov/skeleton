/*
Skeleton ☠️


Skeleton is a telegram bot framework based by github.com/Syfaro/telegram-bot-api

You can use all methods Syfaro library.


Examples


Simple app:
	func main() {
		// create app
		app := skeleton.NewBot("TELERAM-TOKEN-BOT", )

		// create new handler
		// seach commad /start
		// work only private
		app.HandleFunc("^/start$", func (c *skeleton.Context) bool {

		c.BotAPI.Send(tgbotapi.NewMessage(c.ChatId(),
		"Hello, " + c.Update.Message.Chat.FirstName))

		c.Pipeline().Next()

		return true

		}).Border(sk.Private).Methods(sk.Commands)

		// start app
		app.Run()
	}

Pipeline app:
	func main() {

		// create app
		app := skeleton.NewBot("TELERAM-TOKEN-BOT")

		// create pipeline
		// enter point command start
		// work only bot
		//
		// Append() add next
		pl := app.HandleFunc("/start", func(c *skeleton.Context) bool {
			c.BotAPI.Send(tgbotapi.NewMessage(c.ChatId(),
				"Hello, "+c.Update.Message.Chat.FirstName+" how old are you?"))

			// set next command
			c.Pipeline().Next()

			return true
		}).Border(skeleton.Private).Methods(skeleton.Commands).Append()

		// add next command
		// set timeout 20 second
		pl = pl.Func(func(c *skeleton.Context) bool {
			c.BotAPI.Send(tgbotapi.NewMessage(c.ChatId(),
				"Ok! You "+c.Update.Message.Text+" old! Where you were born?"))

			// set next command
			c.Pipeline().Next()

			return true
		}).Timeout(time.Second * 20).Append()

		// add next command
		// set timeout 20 second
		pl = pl.Func(func(c *skeleton.Context) bool {
			c.BotAPI.Send(tgbotapi.NewMessage(c.ChatId(),
				"Ok! You born "+c.Update.Message.Text+"! Nice to meet you"))

			// stop pipeline
			c.Pipeline().Stop()

			return true
		}).Timeout(time.Second * 20)

		app.Run()
	}

Keyboard app:
	func main() {

		// create app
		app := skeleton.NewBot(
			"TELERAM-TOKEN-BOT"")

		// create new handler
		// seach commad /start
		// work only private
		app.HandleFunc("^/start$", func(c *skeleton.Context) bool {

			// create new text message
			msg := tgbotapi.NewMessage(c.ChatId(),
				"Hello, "+c.Update.Message.Chat.FirstName+", tap to button")

			// create new inline keyboard generator
			k := skeleton.NewInlineKeyboard(1, 3)
			// add buttons
			k.Buttons.Add("1", "is-1")
			k.Buttons.Add("2", "is-2")
			k.Buttons.Add("3", "is-3")
			k.Buttons.Add("4", "is-4")
			k.Buttons.Add("5", "is-5")
			k.Buttons.Add("6", "is-6")
			k.Buttons.Add("7", "is-7")
			// add control buttons
			k.ControlButtons.AddPagesControl()

			// generate keyboard and add in message markup
			msg.ReplyMarkup = k.Generate().InlineKeyboardMarkup()

			// if message send dont have errors - save keyboard state
			if m, err := c.BotAPI.Send(msg); err == nil {
				c.Keyboards().SaveState(&m, k)
			}

			return true
		}).Border(skeleton.Private).Methods(skeleton.Commands)

		// add new handler callback from inline keyboard
		app.HandleFunc(`is-(\d{1,})`, func(c *skeleton.Context) bool {

			// RegexpResult using FindStringSubmatch()
			// more in processor.go func (a *app) processor(update *tgbotapi.Update) {]
			c.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(
				c.Update.CallbackQuery.ID,
				"Is "+c.RegexpResult[1]+" number!"))

			return true
		}).Border(skeleton.Private).Methods(skeleton.Callbacks)

		// start app
		app.Run()
	}


*/
package skeleton
