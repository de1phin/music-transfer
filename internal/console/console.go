package console

import (
	"bufio"
	"os"
)

type consoleInteractor struct {
	readWriter    *bufio.ReadWriter
	defaultUserID int64
}

func NewConsoleInteractor(defaultUserID int64) *consoleInteractor {
	console := consoleInteractor{
		bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)),
		defaultUserID,
	}
	return &console
}

func (console *consoleInteractor) GetMessageFrom() (int64, string) {
	msg, _ := console.readWriter.ReadString('\n')
	return console.defaultUserID, msg
}

func (console *consoleInteractor) SendMessageTo(id int64, msg string) {
	console.readWriter.WriteString(msg + "\n")
	console.readWriter.Flush()
}

func (console *consoleInteractor) ChooseFrom(id int64, msg string, options []string) {
	console.readWriter.WriteString(msg)

	for _, option := range options {
		console.readWriter.WriteString(option + "\n")
	}
	console.readWriter.Flush()
}

func (console *consoleInteractor) SendURL(id int64, msg string, url string) {
	console.readWriter.WriteString(msg)
	console.readWriter.WriteString(url + "\n")
	console.readWriter.Flush()
}
