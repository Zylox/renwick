package cleverbot

import (
	nslack "github.com/nlopes/slack"
	"github.com/ugjka/cleverbot-go"
	"github.com/zylox/renwick/go/pkg/log"
	"github.com/zylox/renwick/go/pkg/slack"
)

const CleverbotSecretEnvKey = "CLEVERBOT_KEY"

type BotChatter struct {
	session *cleverbot.Session
}

func NewCalbackHandler(apiKey string) slack.SlackAppMessageEventHandler {
	return BotChatter{
		session: cleverbot.New(apiKey),
	}
}

func (chatter BotChatter) Is(_ *nslack.Client, event slack.SlackAppMessageEvent) bool {
	return true
}

func (chatter BotChatter) Act(client *nslack.Client, event slack.SlackAppMessageEvent) error {
	answer, err := chatter.session.Ask(event.Text)
	if err != nil {
		log.ErrorF("cleverbot.Act - Error when asking. Err: %s", err.Error())
		client.PostMessage(event.Channel, "It is about time you leave.", nslack.NewPostMessageParameters())
		return err
	}
	userID := slack.UserID{ID: event.User}
	client.PostMessage(event.Channel, slack.UserResponse(userID, answer), nslack.NewPostMessageParameters())
	return nil
}