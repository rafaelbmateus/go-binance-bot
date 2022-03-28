package notify_test

import (
	"os"
	"testing"

	"github.com/rafaelbmateus/binance-bot/notify"
	"github.com/stretchr/testify/assert"
)

var webhook = os.Getenv("SLACK_WEBHOOK_URL")

func TestSendMessage(t *testing.T) {
	msg := notify.NewMessage("Hello!")
	assert.Equal(t, "Hello!", msg.Text)

	err := notify.NewSlackNotify("Test", webhook).SendMessage(msg)
	assert.NoError(t, err)
}
