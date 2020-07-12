package skeleton

import (
	"github.com/Syfaro/telegram-bot-api"
)

// Run application with
// current params
func (a *app) Run() {

	updates, _ := a.botAPI.GetUpdatesChan(*a.updateConfig)

	// add default func`s
	a.defaultFunc()

	for update := range updates {
		a.processor(&update)
	}

}

// defaultFunc()
func (a *app) defaultFunc()  {

	// ++ keyboard func`s
	a.HandleFunc(`prev`, prev).Border(Private).Methods(Callbacks)
	a.HandleFunc(`next`, next).Border(Private).Methods(Callbacks)
	a.HandleFunc(`list`, list).Border(Private).Methods(Callbacks)
	a.HandleFunc(`back`, back).Border(Private).Methods(Callbacks)
	a.HandleFunc(`show-(\d{0,})`, show).Border(Private).Methods(Callbacks)
	// -- keyboard func`s

}

// processor()
func (a *app) processor(update *tgbotapi.Update) {

	defer func() {
		if err := recover(); err != nil {
			log.Printf("RECOVER AFTER PANIC [%v]", err)
		}
	}()

	command, context, chatId := GetCommand(update)

	if rule, found := a.pipeline.get(chatId); found{

		if ok := rule.command(&Context{
			app:    a,
			rule:   rule,
			chatId: chatId,
			BotAPI: a.botAPI,
			Update: update,
		});

		ok { return }
	}

	for _, rule := range a.rules.rulesMap[context] {

		search := rule.findMatch(command)

		if len(search) == 0 { continue }
		if rule.borderUse != BorderUse(update) { continue }

		if ok := rule.command(&Context{
			app:a,
			chatId:chatId,
			rule:rule,
			BotAPI: a.botAPI,
			Update: update,
			RegexpResult: search,
		});

		!ok { continue }

		return
	}

	if context == Callbacks {
		a.botAPI.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(
			update.CallbackQuery.ID,
			defaultMessage))
	}

	if context == Messages ||
		context == ReplyToMessages ||
		context == Commands {

		a.botAPI.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			defaultMessage))
	}
}


