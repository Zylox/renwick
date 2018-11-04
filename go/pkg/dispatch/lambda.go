package dispatch

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	nslack "github.com/nlopes/slack"
	"github.com/zylox/renwick/go/internal/pkg/aws/secrets"
	"github.com/zylox/renwick/go/internal/pkg/log"
	"github.com/zylox/renwick/go/internal/pkg/slack"
	"github.com/zylox/renwick/go/internal/pkg/utils"

	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack/slackevents"
)

type GatewayProxyFn func(ctx context.Context, gatewayEvent events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func Create() {

	awsSession, err := session.NewSession()
	if err != nil {
		log.FatalF("dispatch.Main - Failed to init aws session. Err: %+v", err)
	}

	slackOauthSecretKey := utils.MustGetEnv(slack.OauthSecretsEnvKey)
	oauthKey := slack.BasicSlackOauth{}
	json.Unmarshal([]byte(secrets.MustGetSecret(awsSession, slackOauthSecretKey)), &oauthKey)
	log.InfoF("dont look: %s %s", oauthKey, secrets.MustGetSecret(awsSession, slackOauthSecretKey))
	lambda.Start(bootStrapHandler(oauthKey))
}

func Challenge(gatewayEvent events.APIGatewayProxyRequest) string {
	var r *slackevents.ChallengeResponse
	err := json.Unmarshal([]byte(gatewayEvent.Body), &r)
	if err != nil {
		log.ErrorF("Failed to unmarshal slack event: %+v", err)
	}
	return r.Challenge
}

func bootStrapHandler(oauthKey slack.SlackToken) GatewayProxyFn {
	return func(ctx context.Context, gatewayEvent events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return HandleRequest(ctx, gatewayEvent, oauthKey)
	}
}

func HandleSlackCallback(client *nslack.Client, event slackevents.EventsAPIEvent) {
	innerEvent := event.InnerEvent
	switch ev := innerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		triggeringUser, err := client.GetUserProfile(ev.User, false)
		var response string
		if err != nil {
			log.ErrorF("dispatch.HandleSlackCallback.AppMention - Could not lookup user name. Err: %+v", err)
			response = "Strange...i can't figure out who you are"
		} else {
			response = fmt.Sprintf("Go away %s", triggeringUser.DisplayName)
		}
		client.PostMessage(ev.Channel, response, nslack.NewPostMessageParameters())
	}
}

func HandleRequest(ctx context.Context, gatewayEvent events.APIGatewayProxyRequest, oauthKey slack.SlackToken) (events.APIGatewayProxyResponse, error) {
	client := nslack.New(oauthKey.BotOauthKey())
	eventsAPIEvent, err := slackevents.ParseEvent(
		json.RawMessage(gatewayEvent.Body),
		slackevents.OptionVerifyToken(
			&slackevents.TokenComparator{VerificationToken: oauthKey.VerificationToken()},
		),
	)
	if err != nil {
		log.ErrorF("Failed to parse slack event: %+v", err)
	}

	log.InfoF("Event: %+v", eventsAPIEvent)

	if eventsAPIEvent.Type == slackevents.URLVerification {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       Challenge(gatewayEvent),
		}, nil
	} else if eventsAPIEvent.Type == slackevents.CallbackEvent {
		HandleSlackCallback(client, eventsAPIEvent)
	}

	log.InfoF("Message: %+v", gatewayEvent)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Pong",
	}, nil
}
