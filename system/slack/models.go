package slack

type SlackMessage struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}
