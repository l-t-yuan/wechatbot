package handlers

import (
	"log"
	"strings"

	"github.com/869413421/wechatbot/gtp"
)

func WebHandler(msg string) string {
	// 向GPT发起请求
	requestText := strings.TrimSpace(msg)
	requestText = strings.Trim(msg, "\n")
	log.Printf(requestText)
	// return "dd"
	client := gtp.GetChatGptBot()
	reply, err := client.Chat(requestText, "base")
	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		return "机器人神了，我一会发现了就去修。"
	}
	if reply == "" {
		return "error"
	}

	// 回复用户
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")
	if err != nil {
		log.Printf("response user error: %v \n", err)
	}
	return reply
}

func WebImageHandler(msg string) string {
	// 向GPT发起请求
	requestText := strings.TrimSpace(msg)
	requestText = strings.Trim(msg, "\n")
	log.Printf(requestText)
	// return "dd"
	client := gtp.GetChatGptBot()
	reply, err := client.DrawImg(requestText)
	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		return "机器人神了，我一会发现了就去修。"
	}
	if reply == "" {
		return "error"
	}

	// 回复用户
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")
	if err != nil {
		log.Printf("response user error: %v \n", err)
	}
	return reply
}
