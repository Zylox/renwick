package dispatch

import (
	"github.com/aws/aws-sdk-go/aws/session"
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
	json.Unmarshal([]byte(secrets.MustGetSecret(awsSession, slackOauthSecretKey)), oauthKey)
	log.InfoF("dont look: %s", oauthKey)
	lambda.Start(bootStrapHandler(oauthKey))
}

func bootStrapHandler(oauthKey slack.SlackToken) GatewayProxyFn {
	return func(ctx context.Context, gatewayEvent events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return HandleRequest(ctx, gatewayEvent, oauthKey)
	}
}

func HandleRequest(ctx context.Context, gatewayEvent events.APIGatewayProxyRequest, oauthKey slack.SlackToken) (events.APIGatewayProxyResponse, error) {
	// api := nslack.New(oauthKey.BotKey())

	eventsAPIEvent, err := slackevents.ParseEvent(
		json.RawMessage(gatewayEvent.Body),
		slackevents.OptionVerifyToken(
			&slackevents.TokenComparator{VerificationToken: oauthKey.VerificationToken()},
		),
	)

	if err != nil {
		log.ErrorF("Failed to parse slack event: %+v", err)
	}

	log.InfoF("WTF: %+v", eventsAPIEvent)

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(gatewayEvent.Body), &r)
		if err != nil {
			log.ErrorF("Failed to unmarshal slack event: %+v", err)
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       r.Challenge,
		}, nil
	}

	log.InfoF("Message: %+v", gatewayEvent)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Pong",
	}, nil
}
