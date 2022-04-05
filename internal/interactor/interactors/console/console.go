package console

import (
	"bufio"
	"os"

	"github.com/de1phin/music-transfer/internal/mux"
)

type ConsoleInteractor struct {
	readWriter    *bufio.ReadWriter
	defaultUserID int64
}

func NewConsoleInteractor(defaultUserID int64) *ConsoleInteractor {
	console := ConsoleInteractor{
		bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)),
		defaultUserID,
	}
	return &console
}

func (console *ConsoleInteractor) GetMessage() mux.Message {
	msg, _ := console.readWriter.ReadString('\n')
	return mux.Message{UserID: console.defaultUserID, Content: msg}
}

func (console *ConsoleInteractor) SendMessage(msg mux.Message) {
	console.readWriter.WriteString(msg.Content)
	console.readWriter.Flush()
}
