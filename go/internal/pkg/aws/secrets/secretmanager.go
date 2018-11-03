package secrets

import (
	"github.com/zylox/renwick/go/internal/pkg/log"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)


func GetSecret(sess *session.Session, secretID string) *string {
	sm := secretsmanager.New(sess)
	input := &secretsmanager.GetSecretValueInput{SecretId: aws.String(secretID)}
	value, err := sm.GetSecretValue(input)
	if err != nil {
		log.FatalF("GetSecret - Can't get secret %s. Err: %+v", secretID, err)
	}
	return value.SecretString
}

func MustGetSecret(sess *session.Session, secretID string) string {
	secret := GetSecret(sess, secretID)
	if secret == nil {
		log.FatalF("MustGetSecret - Secret %s was not populated. Failing out", secretID)
	}
	return *secret
}