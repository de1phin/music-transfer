package console

import (
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

func (ca *ConsoleAdapter) GetMessage() (msg mux.Message, err error) {
	text, err := ca.console.GetMessage()
	if err != nil {
		return msg, err
	}
	msg.UserID = ca.defaultUserID
	msg.UserState = ca.userState
	msg.Content.Text = strings.ToLower(strings.Trim(text, " \n\r\t"))
	return msg, nil
}

func (ca *ConsoleAdapter) SendMessage(msg mux.Message) error {
	ca.userState = msg.UserState
	text := msg.Content.Text + "\n"
	for _, button := range msg.Content.Buttons {
		text += button + "\n"
	}
	for _, url := range msg.Content.URLs {
		text += url.Link + "\n"
	}
	return ca.console.SendMessage(text)
}
