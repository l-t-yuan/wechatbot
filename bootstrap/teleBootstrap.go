package bootstrap

import (
	"github.com/869413421/wechatbot/handlers"
	"github.com/869413421/wechatbot/config"
	"log"

	"github.com/yanzay/tbot/v2"
)

func RunTele() {
	// log.Printf("tele token: %s \n", config.LoadConfig().TeleToken)
	bot := tbot.New(config.LoadConfig().TeleToken)
	c := bot.Client()
	log.Printf("tele start: %s \n", "try")
	bot.HandleMessage(".*", func(m *tbot.Message) {
		// c.SendChatAction(m.Chat.ID, tbot.ActionTyping)
		// time.Sleep(1 * time.Second)
		// c.SendMessage(m.Chat.ID, "hello!")
		log.Printf("tele send text: %#v \n", m.Text)
		c.SendMessage(m.Chat.ID, handlers.TeleHandler(m.Text, m.Chat.ID))
		// c.SendMessage(m.Chat.ID, "ok")
	})
	err := bot.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tele response text: %s \n", "成功")
}	
