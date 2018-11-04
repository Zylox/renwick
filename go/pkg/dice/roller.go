package dice

import (
	"strings"

	"github.com/justinian/dice"
	nslack "github.com/nlopes/slack"
	"github.com/zylox/renwick/go/pkg/log"
	"github.com/zylox/renwick/go/pkg/slack"
)

func NewCallbackHandler() slack.SlackAppMessageEventHandler {
	return DiceRoller{}
}

const DiceRollCommandString = "roll me"

type DiceRoller struct{}

func (d DiceRoller) Is(_ *nslack.Client, event slack.SlackAppMessageEvent) bool {
	msg := slack.ParseSimpleCommand(event.BotID, event.Text)
	return msg != "" && strings.HasPrefix(strings.ToLower(msg), DiceRollCommandString)
}

func (d DiceRoller) Act(client *nslack.Client, event slack.SlackAppMessageEvent) error {
	msg := slack.ParseSimpleCommand(event.BotID, event.Text)
	trimmedMsg := strings.TrimSpace(strings.TrimPrefix(msg, DiceRollCommandString))
	result, reason, err := dice.Roll(trimmedMsg)
	if err != nil {
		log.ErrorF("dice.Act - Error when rolling. Err: %s", err.Error())
		client.PostMessage(event.Channel, "You managed to break dice. Impressive.", nslack.NewPostMessageParameters())
		return err
	}
	log.InfoF("dice.Act - User: %s,Roll: %s, Dice Reason: %s", event.User, result, reason)
	client.PostMessage(event.Channel, result.String(), nslack.NewPostMessageParameters())
	return nil
}