package telegram

import (
	"encoding/xml"
	"log"
	"strings"

	"github.com/de1phin/music-transfer/internal/interactor/interactors/telegram"
	"github.com/de1phin/music-transfer/internal/mux"
)

type TelegramAdapter struct {
	Telegram *telegram.TelegramBot
}

func NewTelegramAdapter(tg *telegram.TelegramBot) *TelegramAdapter {
	return &TelegramAdapter{
		Telegram: tg,
	}
}

func (ta *TelegramAdapter) SendMessage(msg mux.Message) {
	content := mux.Content{}
	xml.Unmarshal([]byte(msg.Content), &content)
	for _, i := range content.Either {
		if i.Button.Text != "" {
			content.Button = append(content.Button, i.Button)
			continue
		}
		if i.URL.Link != "" {
			content.URL = append(content.URL, i.URL)
		}
	}
	ta.Telegram.SendMessage(msg.UserID, strings.Join(content.Text, "\n"), content.Button, content.URL)
}

func (ta *TelegramAdapter) GetMessage() mux.Message {
	msg := ta.Telegram.GetMessage()
	msg.Content = strings.ToLower(strings.Trim(msg.Content, " \n\r\t"))
	log.Println("Return", msg)
	return msg
}
