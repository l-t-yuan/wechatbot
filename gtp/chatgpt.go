package gtp

import (
	"context"
	"fmt"

	"github.com/869413421/wechatbot/config"
	"github.com/otiai10/openaigo"

	"log"
)

type ChatGptBot struct {
	botChannels map[string]*[]openaigo.ChatMessage
	client      *openaigo.Client
}

var chatGptBotIns *ChatGptBot

func GetChatGptBot() *ChatGptBot {
	if chatGptBotIns == nil {
		chatGptBotIns = &ChatGptBot{}
		chatGptBotIns.init()
	}
	return chatGptBotIns
}

func (c *ChatGptBot) init() {
	c.client = openaigo.NewClient(config.LoadConfig().ApiKey)
	c.botChannels = make(map[string]*[]openaigo.ChatMessage)
}

func (c *ChatGptBot) getChannel(unitKey string) *[]openaigo.ChatMessage {
	chat, isOk := c.botChannels[unitKey]
	if !isOk {
		chat = &[]openaigo.ChatMessage{}
		c.botChannels[unitKey] = chat
	}
	return chat
}

func (c *ChatGptBot) Chat(msg, unitKey string) (string, error) {

	chat := c.getChannel(unitKey)
	tryChat := *chat

	tryChat = append(tryChat, openaigo.ChatMessage{Role: "user", Content: msg})

	request := openaigo.ChatCompletionRequestBody{
		Model:    "gpt-3.5-turbo",
		Messages: tryChat,
	}
	ctx := context.Background()
	response, err := c.client.Chat(ctx, request)
	fmt.Println(response, err)
	if err != nil {
		return "机器人出错了", err
	}
	reply := response.Choices[0].Message.Content
	*chat = append(*chat, openaigo.ChatMessage{Role: "user", Content: msg})
	*chat = append(*chat, openaigo.ChatMessage{Role: "assistant", Content: reply})
	log.Printf("gpt response text: %s \n", reply)
	return reply, nil
}
