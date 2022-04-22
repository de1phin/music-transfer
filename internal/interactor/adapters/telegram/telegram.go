package telegram

import (
	"encoding/xml"
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
	content := mux.Content{}
	xml.Unmarshal([]byte(msg.Content), &content)
	Buttons := make([]telegram.Button, 0)
	URLButtons := make([]telegram.Button, 0)
	for _, i := range content.Button {
		Buttons = append(Buttons, telegram.Button{
			Text: i,
			Data: encodeUserState(msg.UserState),
		})
	}
	for _, i := range content.URL {
		URLButtons = append(URLButtons, telegram.Button{
			Text: i.Text,
			Data: i.Link,
		})
	}
	for _, i := range content.Either {
		if i.Button != "" {
			Buttons = append(Buttons, telegram.Button{
				Text: i.Button,
				Data: encodeUserState(msg.UserState),
			})
			continue
		}
		if i.URL.Link != "" {
			URLButtons = append(URLButtons, telegram.Button{
				Text: i.URL.Text,
				Data: i.URL.Link,
			})
		}
	}
	return ta.Telegram.SendMessage(msg.UserID, strings.Join(content.Text, "\n"), Buttons, URLButtons)
}

func (ta *TelegramAdapter) GetMessage() (mux.Message, error) {
	telegramMessage, err := ta.Telegram.GetMessage()
	if err != nil {
		return mux.Message{}, err
	}
	msg := mux.Message{
		Content: strings.ToLower(strings.Trim(telegramMessage.Text, " \n\r\t")),
		UserID:  telegramMessage.ChatID,
	}
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
	msg.Content = strings.ToLower(strings.Trim(msg.Content, " \n\r\t"))
	return msg, nil
}
