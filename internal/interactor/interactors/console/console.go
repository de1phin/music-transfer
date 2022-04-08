package console

import (
	"bufio"
	"os"
)

type ConsoleInteractor struct {
	readWriter *bufio.ReadWriter
}

func NewConsoleInteractor() *ConsoleInteractor {
	console := ConsoleInteractor{
		bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)),
	}
	return &console
}

func (console *ConsoleInteractor) GetMessage() string {
	msg, _ := console.readWriter.ReadString('\n')
	return msg
}

func (console *ConsoleInteractor) SendMessage(msg string) {
	console.readWriter.WriteString(msg)
	console.readWriter.Flush()
}
