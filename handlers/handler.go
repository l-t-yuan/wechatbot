package handlers

import (
	"log"
	"strings"

	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/gtp"
	"github.com/eatmoreapple/openwechat"
)

// MessageHandlerInterface 消息处理接口
type MessageHandlerInterface interface {
	handle(*openwechat.Message) error
	ReplyText(*openwechat.Message) error
}

type HandlerType string

const (
	GroupHandler = "group"
	UserHandler  = "user"
)

// handlers 所有消息类型类型的处理器
var handlers map[HandlerType]MessageHandlerInterface

func init() {
	handlers = make(map[HandlerType]MessageHandlerInterface)
	handlers[GroupHandler] = NewGroupMessageHandler()
	handlers[UserHandler] = NewUserMessageHandler()
}

// Handler 全局处理入口
func Handler(msg *openwechat.Message) {
	log.Printf("hadler Received msg : %v", msg.Content)
	// 处理群消息
	if msg.IsSendByGroup() {
		handlers[GroupHandler].handle(msg)
		return
	}

	// 好友申请
	if msg.IsFriendAdd() {
		if config.LoadConfig().AutoPass {
			_, err := msg.Agree("你好我是基于chatGPT引擎开发的微信机器人，你可以向我提问任何问题。")
			if err != nil {
				log.Fatalf("add friend agree error : %v", err)
				return
			}
		}
	}

	// 私聊
	handlers[UserHandler].handle(msg)
}

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

func TeleHandler(msg, unitKey string) string {
	// 向GPT发起请求
	requestText := strings.TrimSpace(msg)
	requestText = strings.Trim(msg, "\n")
	log.Printf(requestText)
	// return "dd"

	reply, err := gtp.CompletionsMore(requestText, unitKey)
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
