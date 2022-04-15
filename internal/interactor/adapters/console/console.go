package console

import (
	"encoding/xml"
	"strings"

	"github.com/de1phin/music-transfer/internal/interactor/interactors/console"
	"github.com/de1phin/music-transfer/internal/mux"
)

type ConsoleAdapter struct {
	console       *console.ConsoleInteractor
	defaultUserID int64
	userState     mux.UserState
}

func NewConsoleAdapter(console *console.ConsoleInteractor, defaultUserID int64) *ConsoleAdapter {
	return &ConsoleAdapter{
		console:       console,
		defaultUserID: defaultUserID,
		userState:     mux.Idle,
	}
}

func (*ConsoleAdapter) Name() string {
	return "console"
}

func (ca *ConsoleAdapter) GetMessage() (mux.Message, error) {
	text, err := ca.console.GetMessage()
	if err != nil {
		return mux.Message{}, err
	}
	msg := mux.Message{
		UserID:    ca.defaultUserID,
		UserState: ca.userState,
		Content:   strings.ToLower(strings.Trim(text, " \n\r\t")),
	}
	return msg, nil
}

func (ca *ConsoleAdapter) SendMessage(msg mux.Message) error {
	ca.userState = msg.UserState
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
	return ca.console.SendMessage(text)
}
