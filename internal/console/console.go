package console

import "fmt"

type consoleInteractor struct {
	defaultUserID int64
}

func NewConsoleInteractor(defaultUserID int64) *consoleInteractor {
	console := consoleInteractor{defaultUserID}
	return &console
}

func (console *consoleInteractor) GetMessageFrom() (int64, string) {
	var msg string
	fmt.Scanln(&msg)
	return console.defaultUserID, msg
}

func (console *consoleInteractor) SendMessageTo(id int64, msg string) {
	fmt.Println(msg)
}
