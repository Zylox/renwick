package cleverbot

import (
	"encoding/json"
	"strconv"

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

type CleverbotConfig struct {
	Attentiveness uint8
	Talkativeness uint8
	Wackiness     uint8
}

type BotChatter struct {
	secret   utils.Secret
	config   CleverbotConfig
	sessions map[string]*cleverbot.Session
}

type PostHandler struct {
	topicArn  string
	snsClient *sns.SNS
}

func init() {
	messagePostParameters = nslack.NewPostMessageParameters()
	messagePostParameters.EscapeText = false
}

func getConfigForKey(key string) uint8 {
	s := utils.GetEnv(key)
	if s == "" {
		return 255
	} else {
		v, _ := strconv.Atoi(s)
		return uint8(v)
	}
}

func NewBotChatter(secret utils.Secret, config *CleverbotConfig) *BotChatter {
	if config == nil {
		config = &CleverbotConfig{
			Attentiveness: getConfigForKey(AttentivenessKey),
			Talkativeness: getConfigForKey(TalkativenessKey),
			Wackiness:     getConfigForKey(WackinessKey),
		}
	}
	return &BotChatter{
		secret:   secret,
		config:   *config,
		sessions: make(map[string]*cleverbot.Session),
	}
}

func NewCalbackHandler(awsSession *session.Session, topicArn string) slack.SlackAppMessageEventHandler {
	return PostHandler{
		snsClient: sns.New(awsSession),
		topicArn:  topicArn,
	}
}

func (chatter *BotChatter) GetUserSession(userID slack.UserID) *cleverbot.Session {
	if session, ok := chatter.sessions[userID.ID]; ok {
		return session
	}
	session := cleverbot.New(chatter.secret.MustGetSecret())
	if chatter.config.Talkativeness != 255 {
		session.Talkativeness(chatter.config.Talkativeness)
	}
	if chatter.config.Attentiveness != 255 {
		session.Attentiveness(chatter.config.Attentiveness)
	}
	if chatter.config.Wackiness != 255 {
		session.Wackiness(chatter.config.Wackiness)
	}
	chatter.sessions[userID.ID] = cleverbot.New(chatter.secret.MustGetSecret())
	return chatter.sessions[userID.ID]
}

func (chatter *BotChatter) Chat(clientContainer slack.ClientContainer, event slack.SlackAppMessageEvent) error {
	session := chatter.GetUserSession(slack.UserID{ID: event.User})
	client := clientContainer.GetClient()
	log.InfoF("Entering Chat for event %s", event.TimeStamp)
	answer, err := session.Ask(event.Text)
	if err != nil {
		log.ErrorF("cleverbot.Chat - Error when asking. Err: %s", err.Error())
		client.PostMessage(event.Channel, "It is about time you leave.", nslack.NewPostMessageParameters())
		return err
	}
	userID := slack.UserID{ID: event.User}

	log.InfoF("cleverbot.Chat - Sending mesasge: %s", slack.UserResponse(userID, slackutilsx.EscapeMessage(answer)))
	client.PostMessage(
		event.Channel,
		slack.UserResponse(userID, slackutilsx.EscapeMessage(answer)),
		slack.MaybeAddThread(messagePostParameters, event.ThreadTimeStamp))
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
