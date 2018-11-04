package slack

import (
	nslack "github.com/nlopes/slack"
	"github.com/zylox/renwick/go/pkg/log"
)

type ClientContainer interface {
	GetClient() *nslack.Client
}

type LazyClient struct {
	slackToken SlackToken
	client     *nslack.Client
}

func NewClientContainer(slackToken SlackToken) ClientContainer {
	return &LazyClient{slackToken: slackToken}
}

func (lc *LazyClient) GetClient() *nslack.Client {
	if lc.client == nil {
		log.InfoF("slack.LazyClient.GetClient - Creating slack client")
		lc.client = nslack.New(lc.slackToken.BotOauthKey())
	}
	return lc.client
}
