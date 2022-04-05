package console

import (
	"encoding/xml"
	"strings"

	"github.com/de1phin/music-transfer/internal/interactor/interactors/console"
	"github.com/de1phin/music-transfer/internal/mux"
)

type ConsoleAdapter struct {
	console *console.ConsoleInteractor
}

func NewConsoleAdapter(console *console.ConsoleInteractor) *ConsoleAdapter {
	return &ConsoleAdapter{
		console: console,
	}
}

func (ca *ConsoleAdapter) GetMessage() mux.Message {
	msg := ca.console.GetMessage()
	msg.Content = strings.ToLower(strings.Trim(msg.Content, " \n\r\t"))
	return msg
}

func (ca *ConsoleAdapter) SendMessage(msg mux.Message) {
	content := mux.Content{}
	xml.Unmarshal([]byte(msg.Content), &content)
	text := ""
	for _, i := range content.Text {
		text += i + "\n"
	}
	for _, i := range content.Either {
		if i.Text != "" {
			text += i.Text + "\n"
		}
	}
	for _, i := range content.URL {
		text += i.Link + "\n"
	}
	for _, i := range content.Either {
		if i.Text == "" && i.URL.Link != "" {
			text += i.URL.Link + "\n"
		}
	}
	msg.Content = text
	ca.console.SendMessage(msg)
}
