package isit

import (
	"github.com/zylox/renwick/go/pkg/slack"
	"strings"
)


const CommandString = "is it"
type Isit struct {}


func (_ Isit) Name() string {
	return "Isit"
}

func ()

func (_ Isit) Is(_ slack.ClientContainer, event slack.SlackAppMessageEvent) bool {
	msg := slack.ParseSimpleCommand(event.BotID, event.Text)
	return slack.QualifiesForCommand(msg, CommandString)
}

func (poster Isit) Act(clientContainer slack.ClientContainer, event slack.SlackAppMessageEvent) error {
	// msg, err := json.Marshal(event)
	// if err != nil {
	// 	log.ErrorF("cleverbot.Act - Failed to marshal message. Err: %s", err.Error())
	// 	return err
	// }
	// params := &sns.PublishInput{
	// 	Message:  aws.String(string(msg)),
	// 	TopicArn: aws.String(poster.topicArn),
	// }
	// _, err = poster.snsClient.Publish(params)
	// if err != nil {
	// 	log.ErrorF("cleverbot.Act - Failed to public message. Err: %s", err.Error())
	// }
	// return err
	
	

	return nil
}
