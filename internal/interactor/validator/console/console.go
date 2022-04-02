package console

import (
	"strings"

	"github.com/de1phin/music-transfer/internal/interactor"
)

type Validator struct{}

func (Validator) Validate(msg *interactor.Message) {
	msg.Text = strings.ToLower(
		strings.Trim(msg.Text, " \n\t\r"),
	)
}
