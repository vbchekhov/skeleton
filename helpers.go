package skeleton

import (
	"github.com/Syfaro/telegram-bot-api"
)

// borderUse Determine the boundaries
// use methods:
// 0 - another command
// 1 - work in the bot
// 2 - work in groups
// 3 - work in channels
func BorderUse(u *tgbotapi.Update) int {

	if u.Message != nil {
		if u.Message.Chat.IsPrivate() {
			return 1
		}

		if u.Message.Chat.IsSuperGroup() || u.Message.Chat.IsGroup() {
			return 2
		}

		if u.Message.Chat.IsChannel() {
			return 3
		}
	}

	if u.EditedMessage != nil {
		if u.EditedMessage.Chat.IsPrivate() {
			return 1
		}

		if u.EditedMessage.Chat.IsSuperGroup() || u.Message.Chat.IsGroup() {
			return 2
		}

		if u.EditedMessage.Chat.IsChannel() {
			return 3
		}
	}

	if u.CallbackQuery != nil {
		if u.CallbackQuery.Message.Chat.IsPrivate() {
			return 1
		}

		if u.CallbackQuery.Message.Chat.IsSuperGroup() || u.CallbackQuery.Message.Chat.IsGroup() {
			return 2
		}

		if u.CallbackQuery.Message.Chat.IsChannel() {
			return 3
		}
	}

	if u.ChannelPost != nil {
		if u.ChannelPost.Chat.IsChannel() {
			return 3
		}
	}

	if u.EditedChannelPost != nil {
		if u.EditedChannelPost.Chat.IsChannel() {
			return 3
		}
	}

	return 0
}

// IsCallBack
func IsCallBack(u *tgbotapi.Update) bool {
	return u.CallbackQuery != nil
}

// IsMessage
func IsMessage(u *tgbotapi.Update) bool {
	return u.Message != nil
}

// IsCommand
func IsCommand(u *tgbotapi.Update) bool {
	if u.Message == nil {
		return false
	}
	return  u.Message.IsCommand()
}

// IsReplyToMessage
func IsReplyToMessage(u *tgbotapi.Update) bool {
	if u.Message == nil {
		return false
	}
	return u.Message.ReplyToMessage != nil
}

// IsEditedMessage
func IsEditedMessage(u *tgbotapi.Update) bool {
	return u.EditedMessage != nil
}

// IsChannelPost
func IsChannelPost(u *tgbotapi.Update) bool {
	return u.ChannelPost != nil
}

// IsEditedChannelPost
func IsEditedChannelPost(u *tgbotapi.Update) bool {
	return u.EditedChannelPost != nil
}

// IsInlineQuery
func IsInlineQuery(u *tgbotapi.Update) bool {
	return u.InlineQuery != nil
}

// IsChosenInlineResult
func IsChosenInlineResult(u *tgbotapi.Update) bool {
	return u.ChosenInlineResult != nil
}

// IsShippingQuery
func IsShippingQuery(u *tgbotapi.Update) bool {
	return u.ShippingQuery != nil
}

// IsPreCheckoutQuery
func IsPreCheckoutQuery(u *tgbotapi.Update) bool {
	return u.PreCheckoutQuery != nil
}

// GetCommand
func GetCommand(u *tgbotapi.Update) (command, context string, chatId int64) {

	switch true {

	// callback
	case IsCallBack(u):
		return u.CallbackQuery.Data, Callbacks, u.CallbackQuery.Message.Chat.ID

	// 	reply to message
	case IsReplyToMessage(u):
		if u.Message.Text == "" {
			return u.Message.Caption, ReplyToMessages, u.Message.Chat.ID
		}

		return u.Message.Text, ReplyToMessages, u.Message.Chat.ID

	// 	command message like a /start
	case IsCommand(u):
		if u.Message.Text == "" {
			return u.Message.Caption, Commands, u.Message.Chat.ID
		}

		return u.Message.Text, Commands, u.Message.Chat.ID

	// 	text message
	case IsMessage(u):
		if u.Message.Text == "" {
			return u.Message.Caption, Messages, u.Message.Chat.ID
		}

		return u.Message.Text, Messages, u.Message.Chat.ID

	// 	message editing
	case IsEditedMessage(u):
		if u.EditedMessage.Text == "" {
			return u.EditedMessage.Caption, EditedMessages, u.EditedMessage.Chat.ID
		}

		return u.EditedMessage.Text, EditedMessages, u.EditedMessage.Chat.ID

	// 	channel post
	case IsChannelPost(u):
		if u.ChannelPost.Text == "" {
			return u.ChannelPost.Caption, ChannelPosts, u.ChannelPost.Chat.ID
		}
		return u.ChannelPost.Text, ChannelPosts, u.ChannelPost.Chat.ID

	// 	channel post editing
	case IsEditedChannelPost(u):
		if u.EditedChannelPost.Text == "" {
			return u.EditedChannelPost.Caption, EditedChannelPosts, u.EditedChannelPost.Chat.ID
		}
		return u.EditedChannelPost.Text, EditedChannelPosts, u.EditedChannelPost.Chat.ID

	// 	inline query
	case IsInlineQuery(u):
		return u.InlineQuery.Query, InlineQuerys, int64(u.InlineQuery.From.ID)

	// 	chosen inline result
	case IsChosenInlineResult(u):
		return u.ChosenInlineResult.Query, ChosenInlineResults, int64(u.ChosenInlineResult.From.ID)

	// 	shipping query
	case IsShippingQuery(u):
		return u.ShippingQuery.InvoicePayload, ShippingQuerys, int64(u.ShippingQuery.From.ID)

	// 	checkout query
	case IsPreCheckoutQuery(u):
		return u.PreCheckoutQuery.InvoicePayload, PreCheckoutQuerys, int64(u.PreCheckoutQuery.From.ID)

	}

	return "", "", 0
}