package skeleton

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Run application with
// current params
func (a *app) Run() {

	updates := a.botAPI.GetUpdatesChan(*a.updateConfig)

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
			a.panic(update)
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
		a.botAPI.Send(tgbotapi.NewCallbackWithAlert(
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

// panic on the board!
func (a *app) panic(update *tgbotapi.Update) bool {

	var buf [4096]byte
	var stackFile = "./panic on the board.txt"

	n := runtime.Stack(buf[:], false)

	upd, err := json.Marshal(update)
	if err != nil {
		log.Printf("Error Marshal JSON Update: %v", err)
		return true
	}

	f, err := os.Create(stackFile)
	if err != nil {
		log.Printf("Error create stack file: %v", err)
		return true
	}

	f.Write([]byte("-------------- STACK ----------------\n"))
	f.Write(buf[:n])
	f.Write([]byte("\n\n-------------- UPDATE ----------------\n"))
	f.Write(upd)
	f.Close()

	text := fmt.Sprintf("Wow! Panic on the board!\nUpdateID: %d", update.UpdateID)

	m := tgbotapi.NewDocument(owner, tgbotapi.FileReader{
		Name: stackFile[2:],
		Reader: f,
	})
	m.Caption = text
	m.ParseMode = parseMode

	a.botAPI.Send(m)

	os.Remove(stackFile)

	return true
}
