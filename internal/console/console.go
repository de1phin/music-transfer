package console

import "fmt"

type consoleInteractor struct {
	defaultUserID int64
}

func NewConsoleInteractor(defaultUserID int64) *consoleInteractor {
	console := consoleInteractor{defaultUserID}
	return &console
}

func (console *consoleInteractor) GetMessageFrom() (string, int64) {
	var msg string
	fmt.Scanln(&msg)
	return msg, console.defaultUserID
}

func (console *consoleInteractor) SendMessageTo(msg string, id int64) {
	fmt.Println(msg)
}
