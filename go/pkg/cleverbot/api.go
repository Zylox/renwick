package cleverbot

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
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

type PostHandler struct {
	topicArn  string
	snsClient *sns.SNS
}

func init() {
	messagePostParameters = nslack.NewPostMessageParameters()
	messagePostParameters.EscapeText = false
}

func NewCalbackHandler(awsSession *session.Session, topicArn string) slack.SlackAppMessageEventHandler {
	return PostHandler{
		snsClient: sns.New(awsSession),
		topicArn:  topicArn,
	}
}

func (chatter *BotChatter) initIfNeeded() {
	if chatter.session == nil {
		log.InfoF("Cleverbot.initIfNeeded - Initing cleverbot session")
		chatter.session = cleverbot.New(chatter.secret.MustGetSecret())
	}
}

func (chatter *BotChatter) Chat(clientContainer slack.ClientContainer, event slack.SlackAppMessageEvent) error {
	chatter.initIfNeeded()
	client := clientContainer.GetClient()
	log.InfoF("Entering Chat for event %s", event.TimeStamp)
	answer, err := chatter.session.Ask(event.Text)
	if err != nil {
		log.ErrorF("cleverbot.Chat - Error when asking. Err: %s", err.Error())
		client.PostMessage(event.Channel, "It is about time you leave.", nslack.NewPostMessageParameters())
		return err
	}
	userID := slack.UserID{ID: event.User}

	log.InfoF("cleverbot.Chat - Sending mesasge: %s", slack.UserResponse(userID, slackutilsx.EscapeMessage(answer)))
	client.PostMessage(event.Channel, slack.UserResponse(userID, slackutilsx.EscapeMessage(answer)), messagePostParameters)
	return nil
}

func (_ PostHandler) Name() string {
	return "CleverBot"
}

func (_ PostHandler) Is(_ slack.ClientContainer, event slack.SlackAppMessageEvent) bool {
	return true
}

func (poster PostHandler) Act(clientContainer slack.ClientContainer, event slack.SlackAppMessageEvent) error {
	msg, err := json.Marshal(event)
	if err != nil {
		log.ErrorF("cleverbot.Act - Failed to marshal message. Err: %s", err.Error())
		return err
	}
	params := &sns.PublishInput{
		Message:  aws.String(string(msg)),
		TopicArn: aws.String(poster.topicArn),
	}
	_, err = poster.snsClient.Publish(params)
	if err != nil {
		log.ErrorF("cleverbot.Act - Failed to public message. Err: %s", err.Error())
	}
	return err
}
