package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/gtp"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	sdkginext "github.com/larksuite/oapi-sdk-gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type FeishuHandler struct {
	baseHandler func(c *gin.Context)
	cli         *lark.Client
	eventIdList sync.Map
}
type FeishuValidate struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
}

type FMessageText struct {
	Text string `json:"text"`
}

type FMessageImg struct {
	ImageKey string `json:"image_key"`
}

func (f *FeishuHandler) Init() {
	f.baseHandler = sdkginext.NewEventHandlerFunc(f.GenFeiHandler())
	f.cli = lark.NewClient(config.LoadConfig().FeiAppId, config.LoadConfig().FeiAppSecret, lark.WithLogReqAtDebug(true), lark.WithLogLevel(larkcore.LogLevelDebug))
}
func (f *FeishuHandler) SetCache(key string, value bool, exp time.Duration) {
	f.eventIdList.Store(key, value)
	time.AfterFunc(exp, func() {
		f.eventIdList.Delete(key)
	})
}
func (f *FeishuHandler) GenValidateHandler(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	if body != nil {
		fmt.Printf("请求body内容为:%s", body)
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	rJson := &FeishuValidate{}
	c.ShouldBindBodyWith(rJson, binding.JSON)
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
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
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
	go f.sloveEvent(ctx, event)

	return nil
}

func (f *FeishuHandler) sloveEvent(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	cacheKey := *event.Event.Message.MessageId
	if _, ok := f.eventIdList.Load(cacheKey); ok {
		return nil
	}
	fmt.Println(cacheKey)

	receiveMessageBody := &FMessageText{}
	eventErr := json.Unmarshal([]byte(*event.Event.Message.Content), &receiveMessageBody)
	if eventErr != nil {
		fmt.Println("================json str 转struct==")
		fmt.Println(*event.Event.Message.Content)
		return nil
	}
	receiveMessageText := receiveMessageBody.Text
	unitKey := "feishu_" + *event.Event.Sender.SenderId.OpenId
	client := gtp.GetChatGptBot()
	if strings.HasPrefix(receiveMessageText, "/清理") {
		client.CleanChat(unitKey)
		eventErr = f.responseChat(ctx, "数据已清理", event)
	} else if strings.HasPrefix(receiveMessageText, "/genImg") {
		reply, err := client.DrawImg(receiveMessageText[len("/genImg")+1:])
		if err != nil {
			return nil
		}
		eventErr = f.responseImage(ctx, reply, event)
	} else {
		reply, err := client.Chat(receiveMessageText, unitKey)
		if err != nil {
			return nil
		}
		eventErr = f.responseChat(ctx, reply, event)
	}

	if eventErr == nil {
		f.SetCache(cacheKey, true, time.Second*60)
	}
	return nil
}

func (f *FeishuHandler) responseChat(ctx context.Context, reply string, event *larkim.P2MessageReceiveV1) error {

	tenantKey := event.TenantKey()
	openId := *event.Event.Sender.SenderId.OpenId

	replayStruct := &FMessageText{
		Text: reply,
	}
	replayStructString, err := json.Marshal(replayStruct)
	if err != nil {
		fmt.Printf("\n replayStructString error %#v\n", replayStruct) //{"cost":123.33,"name":"天马星空"}
	}
	fmt.Printf("\n replayStructString %s\n", string(replayStructString))
	// ISV 给指定租户发送消息
	_, err = f.cli.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			ReceiveId(openId).
			Content(string(replayStructString)).
			Build()).
		Build(), larkcore.WithTenantKey(tenantKey))

	// 发送结果处理，resp,err
	// fmt.Println(resp, err)

	return err
}

func (f *FeishuHandler) responseImage(ctx context.Context, url string, event *larkim.P2MessageReceiveV1) error {
	// fmt.Printf("\n responseImage prop %#v\n", msg)
	tenantKey := event.TenantKey()
	openId := *event.Event.Sender.SenderId.OpenId

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	req := larkim.NewCreateImageReqBuilder().
		Body(larkim.NewCreateImageReqBodyBuilder().
			ImageType("message").
			Image(resp.Body).
			Build()).
		Build()
	fresp, err := f.cli.Im.Image.Create(context.Background(), req)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	// 服务端错误处理
	if !fresp.Success() {
		fmt.Println(fresp.Code, fresp.Msg, fresp.RequestId())
		return nil
	}

	replayStruct := &FMessageImg{
		ImageKey: *fresp.Data.ImageKey,
	}
	replayStructString, err := json.Marshal(replayStruct)
	if err != nil {
		fmt.Printf("\n replayStructString error %#v\n", replayStruct)
	}
	fmt.Printf("\n replayStructString %s\n", string(replayStructString))
	// ISV 给指定租户发送消息
	_, err = f.cli.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeImage).
			ReceiveId(openId).
			Content(string(replayStructString)).
			Build()).
		Build(), larkcore.WithTenantKey(tenantKey))

	// 发送结果处理，resp,err
	// fmt.Println(resp, err)

	return err
}

func (f *FeishuHandler) onP2MessageReadV1(ctx context.Context, event *larkim.P2MessageReadV1) error {
	fmt.Println(larkcore.Prettify(event))
	fmt.Println(event.RequestId())
	return nil
}
