package secrets

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/zylox/renwick/go/pkg/log"
)

// No invalidation for now
var secretsCache map[string]*string = make(map[string]*string)

type LazySecret struct {
	SecretID string
	Session  *session.Session
}

func NewLazySecret(session *session.Session, secretID string) LazySecret {
	return LazySecret{Session: session, SecretID: secretID}
}

func (ls LazySecret) GetSecret() *string {
	return GetSecret(ls.Session, ls.SecretID)
}

func (ls LazySecret) MustGetSecret() string {
	return MustGetSecret(ls.Session, ls.SecretID)
}

func GetSecret(sess *session.Session, secretID string) *string {
	if secret, ok := secretsCache[secretID]; ok && secret != nil {
		return secret
	}
	sm := secretsmanager.New(sess)
	input := &secretsmanager.GetSecretValueInput{SecretId: aws.String(secretID)}
	value, err := sm.GetSecretValue(input)
	if err != nil {
		log.FatalF("GetSecret - Can't get secret %s. Err: %+v", secretID, err)
	}
	log.InfoF("GetSecret - Caching secret for secretID: %s", secretID)
	secretsCache[secretID] = value.SecretString
	return value.SecretString
}

func MustGetSecret(sess *session.Session, secretID string) string {
	secret := GetSecret(sess, secretID)
	if secret == nil {
		log.FatalF("MustGetSecret - Secret %s was not populated. Failing out", secretID)
	}
	return *secret
}
