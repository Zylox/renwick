package slack

const OauthSecretsEnvKey = "SLACK_OAUTH"

type SlackOauth interface {
	BotKey() string
}

type BasicSlackOauth struct {
	Key string
}

func (bso BasicSlackOauth) BotKey() string {
	return bso.Key
}
