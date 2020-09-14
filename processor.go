package skeleton

import (
	"github.com/Syfaro/telegram-bot-api"
)

// Run application with
// current params
func (a *app) Run() {

	updates, _ := a.botAPI.GetUpdatesChan(*a.updateConfig)

	// add default func`s

	// ++ keyboard func`s
	a.HandleFunc(`prev`, prev).Border(Private).Methods(Callbacks)
	a.HandleFunc(`next`, next).Border(Private).Methods(Callbacks)
	a.HandleFunc(`list`, list).Border(Private).Methods(Callbacks)
	a.HandleFunc(`back`, back).Border(Private).Methods(Callbacks)
	a.HandleFunc(`show-(\d{0,})`, show).Border(Private).Methods(Callbacks)
	// -- keyboard func`s

	// ++ pipeline func`s
	a.HandleFunc(`abort-pipeline`, abort).Border(Private).Methods(Callbacks)
	// -- pipeline func`s

	for update := range updates {
		a.processor(&update)
	}

}

// processor()
func (a *app) processor(update *tgbotapi.Update) {

	defer func() {
		if err := recover(); err != nil {
			log.Printf("[ERR] RECOVER AFTER PANIC [%v]", err)
		}
	}()

	command, context, chatId := GetCommand(update)

	// check allow list and exist user in list
	if !a.allowList.Empty() && !a.allowList.Exist(chatId) {
		return
	}

	if rule, found := a.pipeline.get(chatId); found {

		if command == "abort-pipeline" && context == Callbacks {
			abort(&Context{
				app:    a,
				rule:   rule,
				chatId: chatId,
				BotAPI: a.botAPI,
				Update: update,
			})

			return
		}

		if ok := rule.command(&Context{
			app:    a,
			rule:   rule,
			chatId: chatId,
			BotAPI: a.botAPI,
			Update: update,
		}); ok {
			return
		}
	}

	for _, rule := range a.rules.rulesMap[context] {

		search := rule.findMatch(command)

		if len(search) == 0 {
			continue
		}
		if rule.borderUse != BorderUse(update) {
			continue
		}

		// check allow list and exist user in list
		if !rule.allowList.Empty() && !rule.allowList.Exist(chatId) {
			return
		}

		if ok := rule.command(&Context{
			app:          a,
			chatId:       chatId,
			rule:         rule,
			BotAPI:       a.botAPI,
			Update:       update,
			RegexpResult: search,
		}); !ok {
			continue
		}

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
