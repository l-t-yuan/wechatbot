package gtp

import (
	chatGptPkg "github.com/golang-infrastructure/go-ChatGPT"
	"github.com/869413421/wechatbot/config"
	"strings"
	"log"
)

var ChatGptChannel = make(map[string]*chatGptPkg.ChatGPT)

func getChannel(unitKey string) *chatGptPkg.ChatGPT {
	chat, isOk := ChatGptChannel[unitKey]
	if !isOk {
		jwt := config.LoadConfig().JwtToken
		chat = chatGptPkg.NewChatGPT(jwt)
		ChatGptChannel[unitKey] = chat
	}
	return chat
}

func ChatGptCompletions(msg, unitKey string) (string, error) {

	chat := getChannel(unitKey)

	talk, err := chat.Talk(msg)

	if err != nil {
		return "", err
	}
	reply := strings.Join(talk.Message.Content.Parts, ",")
	log.Printf("gpt response text: %s \n", reply)
	return reply, nil
}
