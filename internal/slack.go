package internal

import (
	"fmt"
	"github.com/slack-go/slack"
	"log"
	"reflect"
	"strings"
)

// 에러가 났을 때 slack으로 에러메세지를 전송해주는 코드

// sendSlackMessage 함수는 slack으로 message를 보내준다.
func sendSlackMessage(message string) {
	api := slack.New(CONFIG.Slack.Token)

	_, timestamp, err := api.PostMessage(CONFIG.Slack.Channel, slack.MsgOptionText(message, false))

	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Message successfully sent to channel at %s", timestamp)
}

// makeSlackMessage 함수는 interface를 받아 slack에 보낼 message 형태로 이쁘게 출력해준다.
func makeSlackMessage(info interface{}) string {
	values := reflect.ValueOf(info)
	keys := values.Type()
	var message []string

	for i := 0; i < values.NumField(); i++ {
		key := keys.Field(i).Name
		value := values.Field(i).Interface()
		keyValue := fmt.Sprintf("%s: %v", key, value)
		message = append(message, keyValue)
	}

	finalString := fmt.Sprintf("```%s```", strings.Join(message, "\n"))

	return finalString
}
