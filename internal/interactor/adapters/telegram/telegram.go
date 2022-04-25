package telegram

import (
	"strconv"
	"strings"

	"github.com/de1phin/music-transfer/internal/interactor/interactors/telegram"
	"github.com/de1phin/music-transfer/internal/mux"
	"github.com/de1phin/music-transfer/internal/storage"
)

type TelegramAdapter struct {
	Telegram         *telegram.TelegramBot
	userStateStorage storage.Storage[int64, mux.UserState]
}

func NewTelegramAdapter(tg *telegram.TelegramBot, userStateStorage storage.Storage[int64, mux.UserState]) *TelegramAdapter {
	return &TelegramAdapter{
		Telegram:         tg,
		userStateStorage: userStateStorage,
	}
}

func (*TelegramAdapter) Name() string {
	return "telegram"
}

func parseUserState(str string) (mux.UserState, error) {
	us, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return mux.Idle, err
	}
	return mux.UserState(us), nil
}

func encodeUserState(us mux.UserState) string {
	return strconv.FormatInt(int64(us), 10)
}

func (ta *TelegramAdapter) SendMessage(msg mux.Message) error {
	Buttons := make([]telegram.Button, 0)
	URLButtons := make([]telegram.Button, 0)
	for _, button := range msg.Content.Buttons {
		Buttons = append(Buttons, telegram.Button{
			Text: button,
			Data: encodeUserState(msg.UserState),
		})
	}
	for _, url := range msg.Content.URLs {
		URLButtons = append(URLButtons, telegram.Button{
			Text: url.Text,
			Data: url.Link,
		})
	}

	return ta.Telegram.SendMessage(msg.UserID, msg.Content.Text, Buttons, URLButtons)
}

func (ta *TelegramAdapter) GetMessage() (msg mux.Message, err error) {
	telegramMessage, err := ta.Telegram.GetMessage()
	if err != nil {
		return mux.Message{}, err
	}
	msg.Content.Text = strings.ToLower(strings.Trim(telegramMessage.Text, " \n\r\t"))
	msg.UserID = telegramMessage.ChatID
	if telegramMessage.Data == "" {
		msg.UserState, err = ta.userStateStorage.Get(msg.UserID)
		if err != nil {
			return mux.Message{}, err
		}
	} else {
		us, err := parseUserState(telegramMessage.Data)
		if err != nil {
			us, err = ta.userStateStorage.Get(msg.UserID)
			if err != nil {
				return mux.Message{}, err
			}
		}
		msg.UserState = us
	}
	return msg, err
}
