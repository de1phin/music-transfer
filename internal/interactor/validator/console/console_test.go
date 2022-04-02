package console_test

import (
	"testing"

	"github.com/de1phin/music-transfer/internal/interactor"
	"github.com/de1phin/music-transfer/internal/interactor/validator/console"
)

func TestConsoleValidator(t *testing.T) {
	cv := console.Validator{}
	text := interactor.Message{UserID: 0, Text: " AbObUs\n"}
	cv.Validate(&text)
	if text.Text != "abobus" {
		t.Fatal("Expected abobus, got", text.Text)
	}
	if text.UserID != 0 {
		t.Fatal("Corrupted UserID")
	}
}
