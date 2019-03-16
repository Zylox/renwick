package slack

import (
	"regexp"
	"strings"
)

const all = -1

type UserID struct {
	ID string
}

func (u UserID) ToLiteral() string {
	return `<@` + u.ID + `>`
}

func DetectUsers(msg string) []UserID {
	re := regexp.MustCompile(`\<\@([^>]*)\>`)
	users := re.FindAllStringSubmatch(msg, all)
	userIDs := []UserID{}
	for _, user := range users {
		userIDs = append(userIDs, UserID{user[1]})
	}
	return userIDs
}

func GetPositiveKarmaRecipients(msg string) []string {
	re := regexp.MustCompile(`([^\s-+]*)\+\+`)
	recipientMatches := re.FindAllStringSubmatch(msg, all)
	recipients := []string{}
	for _, recipient := range recipientMatches {
		if recipient[1] != "" {
			recipients = append(recipients, recipient[1])
		}
	}
	return recipients
}

func GetNegativeKarmaRecipients(msg string) []string {
	re := regexp.MustCompile(`([^\s-+]*)\-\-`)
	recipientMatches := re.FindAllStringSubmatch(msg, all)
	recipients := []string{}
	for _, recipient := range recipientMatches {
		if recipient[1] != "" {
			recipients = append(recipients, recipient[1])
		}
	}
	return recipients
}

func UserResponse(userID UserID, msg string) string {
	return userID.ToLiteral() + ` ` + msg
}

func StripUserOnce(userID UserID, msg string) string {
	return strings.Replace(msg, userID.ToLiteral(), "", 1)
}

func PrefixedWithBot(botID UserID, msg string) bool {
	return strings.HasPrefix(msg, botID.ToLiteral())
}

func ParseSimpleCommand(botID UserID, msg string) string {
	if c := PrefixedWithBot(botID, msg); c {
		return strings.TrimSpace(StripUserOnce(botID, msg))
	}
	return ""
}

func QualifiesForCommand(msg string, commandString string) bool {
	return msg != "" && strings.HasPrefix(strings.ToLower(msg), commandString)
}
