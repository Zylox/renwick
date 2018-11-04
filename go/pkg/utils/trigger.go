package utils

import (
	"github.com/nlopes/slack/slackevents"
)

type CallbackHandler interface {
	Is(slackevents.EventsAPIEvent) bool
	Act(slackevents.EventsAPIEvent) error
}
