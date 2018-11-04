package cleverbot

// import (
// 	"github.com/zylox/renwick/go/pkg/aws/secrets"
// 	"encoding/json"
// 	"github.com/zylox/renwick/go/pkg/log"
// 	"github.com/zylox/renwick/go/pkg/slack"
// 	"github.com/zylox/renwick/go/pkg/utils"
// 	"github.com/aws/aws-sdk-go/aws/session"

// )

// var callBackHandler

// func Create() {
// 	awsSession, err := session.NewSession()
// 	if err != nil {
// 		log.FatalF("dispatch.Main - Failed to init aws session. Err: %+v", err)
// 	}

// 	slackOauthSecretKey := utils.MustGetEnv(slack.OauthSecretsEnvKey)
// 	oauthKey := slack.BasicSlackOauth{}
// 	json.Unmarshal([]byte(secrets.MustGetSecret(awsSession, slackOauthSecretKey)), &oauthKey)

// 	if cbsek := utils.GetEnv(CleverbotSecretEnvKey); cbsek != "" {
// 		cleverBotKey := secrets.MustGetSecret(awsSession, cbsek)
// 		fallbackHandler = cleverbot.NewCalbackHandler(cleverBotKey)
// 	}

// 	lambda.Start(bootStrapHandler(oauthKey))
// }

// if cbsek := utils.GetEnv(cleverbot.CleverbotSecretEnvKey); cbsek != "" {
// 	cleverBotKey := secrets.MustGetSecret(awsSession, cbsek)
// 	fallbackHandler = cleverbot.NewCalbackHandler(cleverBotKey)
// }
