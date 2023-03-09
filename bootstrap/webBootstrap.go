package bootstrap

import (
	"net/http"

	"github.com/869413421/wechatbot/handlers"
	"github.com/gin-gonic/gin"
)

func RunWeb() {
	r := gin.Default()
	// 在已有 Gin 实例上注册消息处理路由
	// r.POST("/feishu/event", sdkginext.NewEventHandlerFunc(feishuHandler.GenFeiHandler()))
	// feishuGinRoute := sdkginext.NewEventHandlerFunc(feishuHandler.GenFeiHandler())
	r.POST("/feishu/:appId/event", func(c *gin.Context) {
		id := c.Param("appId")
		h := handlers.GetFeishuHandler(id)
		if h == nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "appId not found",
			})
		}
		h.GenValidateHandler(c)
	})

	// r.POST("/feishu/validate", feishuHandler.GenValidateHandler)
	r.POST("/web/send", func(c *gin.Context) {
		sendInfo := c.PostForm("send")
		c.JSON(http.StatusOK, gin.H{
			"message": handlers.WebHandler(sendInfo),
		})
	})

	r.POST("/web/image", func(c *gin.Context) {
		sendInfo := c.PostForm("send")
		c.JSON(http.StatusOK, gin.H{
			"message": handlers.WebImageHandler(sendInfo),
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
