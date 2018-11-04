package dispatch

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/zylox/renwick/go/pkg/aws/secrets"
	"github.com/zylox/renwick/go/pkg/cleverbot"
	"github.com/zylox/renwick/go/pkg/dice"
	"github.com/zylox/renwick/go/pkg/log"
	"github.com/zylox/renwick/go/pkg/slack"
	"github.com/zylox/renwick/go/pkg/utils"

	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack/slackevents"
)

var appMentionCallBackHandlers []slack.SlackAppMessageEventHandler
var fallbackHandler slack.SlackAppMessageEventHandler

func init() {
	appMentionCallBackHandlers = append(
		appMentionCallBackHandlers,
		dice.NewCallbackHandler(),
	)
}

type GatewayProxyFn func(ctx context.Context, gatewayEvent events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func Create() {
	awsSession, err := session.NewSession()
	if err != nil {
		log.FatalF("dispatch.Main - Failed to init aws session. Err: %+v", err)
	}

	slackOauthSecret := secrets.NewLazySecret(
		awsSession,
		utils.MustGetEnv(slack.OauthSecretsEnvKey),
	)

	if cbsek := utils.GetEnv(cleverbot.CleverbotSecretEnvKey); cbsek != "" {
		// cleverBotKey := secrets.MustGetSecret(awsSession, cbsek)
		cleverBotKey := secrets.NewLazySecret(awsSession, cbsek)
		fallbackHandler = cleverbot.NewCalbackHandler(cleverBotKey)
	}
	oauthKey := slack.NewLazySlackOauth(slackOauthSecret)
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

func HandleSlackCallback(clientContainer slack.ClientContainer, event slackevents.EventsAPIEvent, authedUsers []string) {
	botID := slack.UserID{ID: authedUsers[0]}
	innerEvent := event.InnerEvent
	switch ev := innerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		same := slack.SlackAppMessageEvent{
			AppMentionEvent: *ev,
			BotID:           botID,
			Token:           event.Token,
			TeamID:          event.TeamID,
		}

		handledAtleastOnce := false
		for _, handler := range appMentionCallBackHandlers {
			if handler.Is(clientContainer, same) {
				handledAtleastOnce = true
				handler.Act(clientContainer, same)
			}
		}
		if !handledAtleastOnce && fallbackHandler != nil {
			log.InfoF("HandleSlackCallback - entering fallback callback")
			fallbackHandler.Act(clientContainer, same)
		}

		// triggeringUser, err := client.GetUserInfo(ev.User)
		// var response string
		// if err != nil {
		// 	log.ErrorF("dispatch.HandleSlackCallback.AppMention - Could not lookup user name. Err: %+v", err)
		// 	response = "Strange...i can't figure out who you are"
		// } else {
		// 	response = fmt.Sprintf("Go away %s", triggeringUser.Profile.DisplayName)
		// }

	}
}

func HandleRequest(ctx context.Context, gatewayEvent events.APIGatewayProxyRequest, oauthKey slack.SlackToken) (events.APIGatewayProxyResponse, error) {
	//client := nslack.New(oauthKey.BotOauthKey())
	clientContainer := slack.NewClientContainer(oauthKey)
	rawE := json.RawMessage(gatewayEvent.Body)
	eventsAPIEvent, err := slackevents.ParseEvent(
		rawE,
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
		// Don't retry, likely caused by lambda warmup
		if reason, ok := gatewayEvent.Headers["X-Slack-Retry-Reason"]; ok && reason == "http_timeout" {
			log.InfoF("HandleRequest - Sending no retry response")
			headers := make(map[string]string)
			headers["X-Slack-No-Retry"] = "1"
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNoContent,
				Body:       "Stop it",
				Headers:    headers,
			}, nil
		}
		cbEvent := &slackevents.EventsAPICallbackEvent{}
		err = json.Unmarshal(rawE, cbEvent)
		if err != nil {
			log.ErrorF("Failed to parse slack Outer event to get users: %+v", err)
		}
		HandleSlackCallback(clientContainer, eventsAPIEvent, cbEvent.AuthedUsers)
	}

	log.InfoF("Message: %+v", gatewayEvent)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Pong",
	}, nil
}
