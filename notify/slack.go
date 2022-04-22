package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const SlackTimeout = 5 * time.Second

type SlackNotify struct {
	Username   string
	WebHookURL string
}

func NewSlackNotify(username, webhook string) *SlackNotify {
	return &SlackNotify{
		Username:   username,
		WebHookURL: webhook,
	}
}

type SlackMessage struct {
	Username string `json:"username,omitempty"`
	Text     string `json:"text,omitempty"`
	Icon     string `json:"icon_emoji"`
}

func NewMessage(text string) *SlackMessage {
	return &SlackMessage{
		Text: text,
	}
}

func (me SlackNotify) SendMessage(slackRequest *SlackMessage) error {
	if me.WebHookURL == "" {
		return nil
	}

	slackRequest.Username = me.Username
	slackRequest.Icon = ":robot_face:"
	slackMsg, err := json.Marshal(slackRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, me.WebHookURL, bytes.NewBuffer(slackMsg))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: SlackTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack api is unavailable %d", resp.StatusCode)
	}

	return nil
}
