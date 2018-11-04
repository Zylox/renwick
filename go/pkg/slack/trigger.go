package slack

import (
	"github.com/nlopes/slack/slackevents"
)

type SlackAppMessageEvent struct {
	slackevents.AppMentionEvent
	BotID  UserID
	Token  string `json:"token"`
	TeamID string `json:"team_id"`
}

type SlackAppMessageEventHandler interface {
	Is(ClientContainer, SlackAppMessageEvent) bool
	Act(ClientContainer, SlackAppMessageEvent) error
	Name() string
}
