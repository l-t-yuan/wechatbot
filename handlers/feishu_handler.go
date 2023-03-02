package handlers

import (
	"context"
	"fmt"

	"github.com/869413421/wechatbot/config"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type FeishuHandler struct {
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
