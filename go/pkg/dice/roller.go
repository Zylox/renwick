package dice

import (
	"github.com/nlopes/slack/slackevents"
	"github.com/zylox/renwick/go/pkg/utils"
)

func NewCallbackHandler() utils.CallbackHandler {
	return DiceRoller{}
}

type DiceRoller struct{}

func (d DiceRoller) Is(event slackevents.EventsAPIEvent) bool {
	return false
}

func (d DiceRoller) Act(event slackevents.EventsAPIEvent) error {
	return nil
}
