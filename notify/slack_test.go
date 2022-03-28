package notify_test

import (
	"testing"

	"github.com/rafaelbmateus/binance-bot/notify"
	"github.com/stretchr/testify/assert"
)

var webhook = ""

func TestSendMessage(t *testing.T) {
	msg := notify.NewMessage("Hello!")
	assert.Equal(t, "Hello!", msg.Text)

	err := notify.NewSlackNotify("Test", webhook).SendMessage(msg)
	assert.NoError(t, err)
}
