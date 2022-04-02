package console

import (
	"bufio"
	"os"

	"github.com/de1phin/music-transfer/internal/interactor"
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

func (console *consoleInteractor) GetMessage() interactor.Message {
	msg, _ := console.readWriter.ReadString('\n')
	return interactor.Message{UserID: console.defaultUserID, Text: msg}
}

func (console *consoleInteractor) SendMessage(msg interactor.Message) {
	console.readWriter.WriteString(msg.Text + "\n")
	console.readWriter.Flush()
}
