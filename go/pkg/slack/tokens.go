package slack

import (
	"encoding/json"

	"github.com/zylox/renwick/go/pkg/aws/secrets"
	"github.com/zylox/renwick/go/pkg/log"
)

const OauthSecretsEnvKey = "SLACK_OAUTH"

type SlackToken interface {
	BotOauthKey() string
	VerificationToken() string
}

type BasicSlackOauth struct {
	OauthKey string `json:"oauth_key,omitempty"`
	VToken   string `json:"verification_token,omitempty"`
}

type LazySlackOauth struct {
	container *BasicSlackOauth
	Secret    secrets.LazySecret
}

func NewLazySlackOauth(secret secrets.LazySecret) *LazySlackOauth {
	return &LazySlackOauth{Secret: secret}
}

func (lso *LazySlackOauth) BotOauthKey() string {
	if lso.container == nil {
		lso.container = &BasicSlackOauth{}
		secretValue := lso.Secret.MustGetSecret()
		err := json.Unmarshal([]byte(secretValue), lso.container)
		if err != nil {
			log.FatalF("LazySlackOauth.BotOauthKey - Failed to unmarshal secret, fix secret. Err: %s", err.Error())
		}
	}
	return lso.container.OauthKey
}

func (lso *LazySlackOauth) VerificationToken() string {
	if lso.container == nil {
		lso.container = &BasicSlackOauth{}
		secretValue := lso.Secret.MustGetSecret()
		err := json.Unmarshal([]byte(secretValue), lso.container)
		if err != nil {
			log.FatalF("LazySlackOauth.VerificationToken - Failed to unmarshal secret, fix secret. Err: %s", err.Error())
		}
	}
	return lso.container.VToken
}

func (bso BasicSlackOauth) BotOauthKey() string {
	return bso.OauthKey
}

func (bso BasicSlackOauth) VerificationToken() string {
	return bso.VToken
}
