package bootstrap

import (
	// "github.com/869413421/wechatbot/handlers"
	"github.com/869413421/wechatbot/config"
	"time"
	"log"

	"github.com/yanzay/tbot/v2"
)

func RunTele() {
	// log.Printf("tele token: %s \n", config.LoadConfig().TeleToken)
	bot := tbot.New(config.LoadConfig().TeleToken)
	c := bot.Client()

	bot.HandleMessage("*", func(m *tbot.Message) {
		// c.SendChatAction(m.Chat.ID, tbot.ActionTyping)
		// time.Sleep(1 * time.Second)
		// c.SendMessage(m.Chat.ID, "hello!")
			log.Printf("tele response text: %#v \n", m)
		// c.SendMessage(m.Chat.ID, hadler.TeleHandler("", m.Chat.ID.String()))
		c.SendMessage(m.Chat.ID, "ok")
	})
	err := bot.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tele response text: %s \n", "成功")
}	
