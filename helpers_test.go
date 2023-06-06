package skeleton

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestGetCommand(t *testing.T) {

	tests := []struct {
		name string
		u    *tgbotapi.Update
		want string
	}{
		{
			name: "GetCommand() #1",
			u:    &tgbotapi.Update{Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 15431451}, Text: "is text message"}},
			want: "is text message",
		},
		{
			name: "GetCommand() #2",
			u:    &tgbotapi.Update{Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 15431451}, Caption: "is caption message"}},
			want: "is caption message",
		},
		{
			name: "GetCommand() #3",
			u:    &tgbotapi.Update{Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 15431451}, Text: "is text reply to message", ReplyToMessage: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 15431451}, Text: "is text reply to message"}}},
			want: "is text reply to message",
		},
		{
			name: "GetCommand() #4",
			u:    &tgbotapi.Update{Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 15431451}, Caption: "is caption reply to message", ReplyToMessage: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 15431451}, Caption: "is caption reply to message"}}},
			want: "is caption reply to message",
		},
		{
			name: "GetCommand() #5",
			u:    &tgbotapi.Update{CallbackQuery: &tgbotapi.CallbackQuery{Data: "is callback data", Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 15431451}}}},
			want: "is callback data",
		},
		{
			name: "GetCommand() #6",
			u:    &tgbotapi.Update{EditedMessage: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 15431451}, Text: "is edit text message"}},
			want: "is edit text message",
		},
		{
			name: "GetCommand() #7",
			u:    &tgbotapi.Update{ChannelPost: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 15431451}, Text: "is text channel post"}},
			want: "is text channel post",
		},
		{
			name: "GetCommand() #8",
			u:    &tgbotapi.Update{ChannelPost: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 15431451}, Caption: "is caption channel post"}},
			want: "is caption channel post",
		},
		{
			name: "GetCommand() #9",
			u:    &tgbotapi.Update{EditedChannelPost: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 15431451}, Text: "is edit text channel post"}},
			want: "is edit text channel post",
		},
		{
			name: "GetCommand() #10",
			u:    &tgbotapi.Update{EditedChannelPost: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 15431451}, Caption: "is edit caption channel post"}},
			want: "is edit caption channel post",
		},
		{
			name: "GetCommand() #11",
			u: &tgbotapi.Update{InlineQuery: &tgbotapi.InlineQuery{
				From:  &tgbotapi.User{ID: 346245763},
				Query: "is inline query",
			}},
			want: "is inline query",
		},
		{
			name: "GetCommand() #12",
			u: &tgbotapi.Update{ChosenInlineResult: &tgbotapi.ChosenInlineResult{
				From:  &tgbotapi.User{ID: 452645364},
				Query: "is chosen inline result",
			}},
			want: "is chosen inline result",
		},
		{
			name: "GetCommand() #13",
			u: &tgbotapi.Update{ShippingQuery: &tgbotapi.ShippingQuery{
				From:           &tgbotapi.User{ID: 4315152},
				InvoicePayload: "is shipping payload",
			}},
			want: "is shipping payload",
		},
		{
			name: "GetCommand() #14",
			u: &tgbotapi.Update{PreCheckoutQuery: &tgbotapi.PreCheckoutQuery{
				From:           &tgbotapi.User{ID: 236245634},
				InvoicePayload: "is invoice payload",
			}},
			want: "is invoice payload",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _, _ := GetCommand(tt.u); got != tt.want {
				t.Errorf("GetCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
