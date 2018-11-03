package main

import (
	"github.com/renwick/go/internal/pkg/log"
	"github.com/renwick/go/internal/pkg/utils"
	"github.com/renwick/go/internal/pkg/slack"
	
	"github.com/renwick/go/internal/pkg/aws/secrets"
	"github.com/aws/aws-sdk-go/aws/session"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	
	// "github.com/nlopes/slack"
)

type GatewayProxyFn func(ctx context.Context, gatewayEvent events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func bootStrapHandler() GatewayProxyFn {
	return func(ctx context.Context, gatewayEvent events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return HandleRequest(ctx, gatewayEvent)
	}
}


func HandleRequest(ctx context.Context, gatewayEvent events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.InfoF("Message: %+v", gatewayEvent)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body: "Pong",
	}, nil
}

func main() {

	awsSession, err := session.NewSession()
	if err != nil {
		log.FatalF("dispatch.Main - Failed to init aws session. Err: %+v", err)
	}

	slackOauthSecretKey := utils.MustGetEnv(slack.OauthSecretsEnvKey)
	oauthKey := secrets.MustGetSecret(awsSession, slackOauthSecretKey)
	log.InfoF("dont look: %s", oauthKey)
	lambda.Start(bootStrapHandler())
}
