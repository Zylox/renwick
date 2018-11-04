package slack

import (
	nslack "github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

type SlackAppMessageEvent struct {
	slackevents.AppMentionEvent
	BotID  UserID
	Token  string `json:"token"`
	TeamID string `json:"team_id"`
}

type SlackAppMessageEventHandler interface {
	Is(*nslack.Client, SlackAppMessageEvent) bool
	Act(*nslack.Client, SlackAppMessageEvent) error
}
