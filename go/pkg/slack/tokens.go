package slack

const OauthSecretsEnvKey = "SLACK_OAUTH"

type SlackToken interface {
	BotOauthKey() string
	VerificationToken() string
}

type BasicSlackOauth struct {
	OauthKey string `json:"oauth_key,omitempty"`
	VToken   string `json:"verification_token,omitempty"`
}

func (bso BasicSlackOauth) BotOauthKey() string {
	return bso.OauthKey
}

func (bso BasicSlackOauth) VerificationToken() string {
	return bso.VToken
}
