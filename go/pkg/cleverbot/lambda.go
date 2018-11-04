package cleverbot

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/zylox/renwick/go/pkg/aws/secrets"
	"github.com/zylox/renwick/go/pkg/log"
	"github.com/zylox/renwick/go/pkg/slack"
	"github.com/zylox/renwick/go/pkg/utils"
)

const CleverbotArnKey = "CLEVERBOT_ARN"

func Create() {
	awsSession, err := session.NewSession()
	if err != nil {
		log.FatalF("dispatch.Main - Failed to init aws session. Err: %+v", err)
	}

	slackOauthSecret := secrets.NewLazySecret(
		awsSession,
		utils.MustGetEnv(slack.OauthSecretsEnvKey),
	)

	cbsek := utils.MustGetEnv(CleverbotSecretEnvKey)
	cleverBotSecret := secrets.NewLazySecret(awsSession, cbsek)
	chatter := NewBotChatter(cleverBotSecret)

	oauthKey := slack.NewLazySlackOauth(slackOauthSecret)
	clientContainer := slack.NewClientContainer(oauthKey)
	lambda.Start(bootStrapHandler(clientContainer, chatter))
}

type SNSEventHandler func(ctx context.Context, snsEvent events.SNSEvent) error

func bootStrapHandler(clientContainer slack.ClientContainer, chatter *BotChatter) SNSEventHandler {
	return func(ctx context.Context, snsEvent events.SNSEvent) error {
		return handler(ctx, snsEvent, clientContainer, chatter)
	}
}

func handler(ctx context.Context, snsEvent events.SNSEvent, clientContainer slack.ClientContainer, chatter *BotChatter) error {
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		slackEvent := slack.SlackAppMessageEvent{}
		err := json.Unmarshal([]byte(snsRecord.Message), &slackEvent)
		if err != nil {
			log.ErrorF("cleverbot.Handler - Could not unmarshal message. Err: %s", err.Error())
			return err
		}
		log.InfoF("cleverbot.handler - Initiating chat")
		chatter.Chat(clientContainer, slackEvent)
	}
	return nil
}
