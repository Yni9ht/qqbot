package main

import (
	"github.com/Yni9ht/qqbot/log"
	"github.com/Yni9ht/qqbot/model"
	"github.com/Yni9ht/qqbot/startup"
	"github.com/Yni9ht/qqbot/vars"
	"github.com/Yni9ht/qqbot/ws"
	"github.com/gorilla/websocket"
	"time"
)

func main() {
	// 加载配置文件
	err := startup.LoadConfig()
	if err != nil {
		log.Errorf("加载配置文件失败：%s", err.Error())
		panic(err)
	}
	vars.Config.Sandbox = true

	// 初始化 openapi client
	err = startup.LoadAPIClient()
	if err != nil {
		log.Errorf("初始化 openapi client 失败：%s", err.Error())
		panic(err)
	}

	// 获取 ws 网关地址
	url, err := vars.OpenAPI.GetGatewayURL()
	if err != nil {
		log.Errorf("get gateway url error: %s", err.Error())
		panic(err)
	}
	log.Debugf("ws gateway url: %s", url)

	// 初始化 ws 连接
	con, res, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Errorf("dial ws error: %v", err)
		panic(err)
	}
	log.Debugf("ws gateway response: %v", res.StatusCode)

	wsClient := &ws.Client{
		URL:         url,
		Token:       vars.Config.GetToken(),
		Intent:      model.IntentGuildAtMessage,
		WssCon:      con,
		MessageChan: make(chan *model.WSPayload, 10),
		CloseChan:   make(chan error, 10),
		Heart:       time.NewTicker(time.Second * 60),
	}

	// 进行鉴权
	if err := wsClient.Identify(); err != nil {
		log.Errorf("identify error: %s", err.Error())
		panic(err)
	}
	log.Infof("identify success")

	if err := wsClient.Start(); err != nil {
		log.Errorf("ws client start error: %s", err.Error())
		return
	}
}
