package cleverbot

import (
	nslack "github.com/nlopes/slack"
	"github.com/nlopes/slack/slackutilsx"
	"github.com/ugjka/cleverbot-go"
	"github.com/zylox/renwick/go/pkg/log"
	"github.com/zylox/renwick/go/pkg/slack"
	"github.com/zylox/renwick/go/pkg/utils"
)

const CleverbotSecretEnvKey = "CLEVERBOT_KEY"

var messagePostParameters nslack.PostMessageParameters

type BotChatter struct {
	secret  utils.Secret
	session *cleverbot.Session
}

func init() {
	messagePostParameters = nslack.NewPostMessageParameters()
	messagePostParameters.EscapeText = false
}

func NewCalbackHandler(secret utils.Secret) slack.SlackAppMessageEventHandler {
	return &BotChatter{
		secret: secret,
	}
}

func (chatter *BotChatter) initIfNeeded() {
	if chatter.session == nil {
		chatter.session = cleverbot.New(chatter.secret.MustGetSecret())
	}
}

// func Chat(client *nslack.Client, event slack.SlackAppMessageEvent) {
// 	chatter.initIfNeeded()
// 	log.InfoF("Entering Act for event %s", event.TimeStamp)
// 	answer, err := chatter.session.Ask(event.Text)
// 	if err != nil {
// 		log.ErrorF("cleverbot.Act - Error when asking. Err: %s", err.Error())
// 		client.PostMessage(event.Channel, "It is about time you leave.", nslack.NewPostMessageParameters())
// 		return err
// 	}
// 	userID := slack.UserID{ID: event.User}

// 	log.InfoF("cleverbot.Act - Sending mesasge: %s", slack.UserResponse(userID, slackutilsx.EscapeMessage(answer)))
// 	client.PostMessage(event.Channel, slack.UserResponse(userID, slackutilsx.EscapeMessage(answer)), messagePostParameters)
// 	return nil
// }

func (_ *BotChatter) Name() string {
	return "CleverBot"
}

func (chatter *BotChatter) Is(_ slack.ClientContainer, event slack.SlackAppMessageEvent) bool {
	return true
}

func (chatter *BotChatter) Act(clientContainer slack.ClientContainer, event slack.SlackAppMessageEvent) error {
	chatter.initIfNeeded()
	client := clientContainer.GetClient()
	log.InfoF("Entering Act for event %s", event.TimeStamp)
	answer, err := chatter.session.Ask(event.Text)
	if err != nil {
		log.ErrorF("cleverbot.Act - Error when asking. Err: %s", err.Error())
		client.PostMessage(event.Channel, "It is about time you leave.", nslack.NewPostMessageParameters())
		return err
	}
	userID := slack.UserID{ID: event.User}

	log.InfoF("cleverbot.Act - Sending mesasge: %s", slack.UserResponse(userID, slackutilsx.EscapeMessage(answer)))
	client.PostMessage(event.Channel, slack.UserResponse(userID, slackutilsx.EscapeMessage(answer)), messagePostParameters)
	return nil
}
