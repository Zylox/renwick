package slack

import (
	"regexp"
)

const all = -1

type UserID struct {
	string
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

// func PrefixedWithBot(slackevents.AppMentionEvent) {
// 	event.
// }
