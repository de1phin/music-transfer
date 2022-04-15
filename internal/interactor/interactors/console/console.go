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

func (console *ConsoleInteractor) GetMessage() (string, error) {
	return console.readWriter.ReadString('\n')
}

func (console *ConsoleInteractor) SendMessage(msg string) error {
	_, err := console.readWriter.WriteString(msg)
	if err != nil {
		return err
	}
	return console.readWriter.Flush()
}
