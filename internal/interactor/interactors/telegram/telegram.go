package telegram

import (
	"github.com/de1phin/music-transfer/internal/mux"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
}

func NewTelegramBot(token string) *TelegramBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	tg := TelegramBot{
		bot:     bot,
		updates: updates,
	}

	return &tg
}

func (tg *TelegramBot) GetMessage() mux.Message {
	upd := <-tg.updates
	if upd.Message == nil {
		tg.bot.Send(tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID, ""))
		return mux.Message{
			UserID:  upd.CallbackQuery.Message.Chat.ID,
			Content: upd.CallbackQuery.Data,
		}
	} else {
		return mux.Message{
			Content: upd.Message.Text,
			UserID:  upd.Message.Chat.ID,
		}
	}
}

func (tg *TelegramBot) SendMessage(chatID int64, text string, buttons []mux.Button, urls []mux.URL) {

	keyboardRows := make([]tgbotapi.InlineKeyboardButton, 0)
	for _, button := range buttons {
		keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardButtonData(button.Text, button.Metadata))
	}
	for _, url := range urls {
		keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardButtonURL(url.Text, url.Link))
	}

	msg := tgbotapi.NewMessage(chatID, text)
	if len(keyboardRows) > 0 {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboardRows)
	}
	tg.bot.Send(msg)

}
