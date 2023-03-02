package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/869413421/wechatbot/config"
	"github.com/gin-gonic/gin"
	sdkginext "github.com/larksuite/oapi-sdk-gin"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type FeishuHandler struct {
	baseHandler func(c *gin.Context)
}
type FeishuValidate struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
}

func (f *FeishuHandler) Init() {
	f.baseHandler = sdkginext.NewEventHandlerFunc(f.GenFeiHandler())
}

func (f *FeishuHandler) GenValidateHandler(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	if body != nil {
		fmt.Printf("请求body内容为:%s", body)
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	rJson := &FeishuValidate{}
	c.BindJSON(rJson)
	if rJson.Type == "url_verification" {
		fmt.Printf("%+v", rJson)
		if rJson.Token == config.LoadConfig().FeiToken {
			c.JSON(http.StatusOK, gin.H{
				"challenge": rJson.Challenge,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"challenge": "error",
			})
		}
	} else {
		f.baseHandler(c)
	}
}

func (f *FeishuHandler) GenFeiHandler() *dispatcher.EventDispatcher {
	handler := dispatcher.NewEventDispatcher(config.LoadConfig().FeiToken, config.LoadConfig().FeiEncrpy).
		OnP2MessageReceiveV1(f.onP2MessageReceiveV1).
		OnP2MessageReadV1(f.onP2MessageReadV1)
		// OnP2UserCreatedV3(func(ctx context.Context, event *larkcontact.P2UserCreatedV3) error {
		// 	fmt.Println(larkcore.Prettify(event))
		// 	fmt.Println(event.RequestId())
		// 	return nil
		// })
	return handler
}

func (f *FeishuHandler) onP2MessageReceiveV1(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	fmt.Println(larkcore.Prettify(event))
	fmt.Println(event.RequestId())
	return nil
}

func (f *FeishuHandler) onP2MessageReadV1(ctx context.Context, event *larkim.P2MessageReadV1) error {
	fmt.Println(larkcore.Prettify(event))
	fmt.Println(event.RequestId())
	return nil
}
