package bootstrap

import (
	"github.com/869413421/wechatbot/handlers"
	"net/http"
	"github.com/gin-gonic/gin"
)

func RunWeb() {
	r := gin.Default()
	r.POST("/send", func(c *gin.Context) {
		sendInfo := c.PostForm("send")
		c.JSON(http.StatusOK, gin.H{
			"message": handlers.WebHandler(sendInfo),
		})
	})
  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}	
