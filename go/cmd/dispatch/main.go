package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	// "github.com/nlopes/slack"
)

func GetSecret(sess *session.Session, secretID string) *string {
	sm := secretsmanager.New(sess)
	input := &secretsmanager.GetSecretValueInput{SecretId: aws.String(secretID)}
	// output, err := sm.GetSecret
	secretsmanager.GetSecretValue(input)
}

func HandleRequest(ctx context.Context, name interface{}) (string, error) {
	return fmt.Sprintf("Hello %s!", ""), nil
}

func main() {
	// slackOauth :=

	lambda.Start(HandleRequest)
}
