package telegram

import (
	"errors"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
}

type Message struct {
	Text   string
	ChatID int64
	Data   string
}

type Button struct {
	Text string
	Data string
}

func NewTelegramBot(config Config) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, err
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	tg := TelegramBot{
		bot:     bot,
		updates: updates,
	}

	return &tg, nil
}

func (tg *TelegramBot) GetMessage() (Message, error) {
	upd := <-tg.updates
	if upd.Message == nil {
		if upd.CallbackQuery == nil {
			return Message{}, errors.New("Empty message")
		}
		data := strings.Split(upd.CallbackQuery.Data, ":::")
		return Message{
			Text:   data[0],
			ChatID: upd.CallbackQuery.Message.Chat.ID,
			Data:   data[1],
		}, nil
	} else {
		if upd.Message.Text == "" {
			return Message{}, errors.New("Empty message")
		}
		return Message{
			Text:   upd.Message.Text,
			ChatID: upd.Message.Chat.ID,
			Data:   "",
		}, nil
	}
}

func (tg *TelegramBot) SendMessage(chatID int64, text string, buttons []Button, urls []Button) error {
	keyboardRows := make([]tgbotapi.InlineKeyboardButton, 0)
	for _, button := range buttons {
		data := fmt.Sprintf("%s:::%s", button.Text, button.Data)
		keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardButtonData(button.Text, data))
	}
	for _, url := range urls {
		keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardButtonURL(url.Text, url.Data))
	}

	msg := tgbotapi.NewMessage(chatID, text)
	if len(keyboardRows) > 0 {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboardRows)
	}
	_, err := tg.bot.Send(msg)
	return err
}
