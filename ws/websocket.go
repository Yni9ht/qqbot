package ws

import (
	"encoding/json"
	"errors"
	"github.com/Yni9ht/qqbot/log"
	"github.com/Yni9ht/qqbot/model"
	"github.com/gorilla/websocket"
	"time"
)

type Client struct {
	URL         string
	Token       string
	Intent      model.Intent
	WssCon      *websocket.Conn
	MessageChan chan *model.WSPayload
	CloseChan   chan error
	LastSeq     uint32
	Heart       *time.Ticker
}

func (c *Client) Start() error {
	defer c.close()

	// 监听消息
	go c.readMessage()

	// 处理消息
	go c.handleMessage()

	for {
		select {
		case err := <-c.CloseChan:
			log.Errorf("ws client close: %v", err)
			c.close()
			return err
		case <-c.Heart.C:
			// 发送心跳信息
			heartBeatEvent := &model.WSPayload{
				WSPayloadBase: model.WSPayloadBase{
					OPCode: model.WSHeartbeat,
				},
				Data: c.LastSeq,
			}
			c.Write(heartBeatEvent)
		}
	}
}

func (c *Client) close() {
	c.WssCon.Close()
	c.Heart.Stop()
}

func (c *Client) readMessage() {
	for {
		_, message, err := c.WssCon.ReadMessage()
		if err != nil {
			return
		}
		payload := &model.WSPayload{}
		err = json.Unmarshal(message, payload)
		if err != nil {
			return
		}
		payload.RawMessage = message

		// 更新最后接收到的消息序列号
		c.updateSeq(payload.Seq)

		// 处理一些默认信息
		if c.handleInsideMessage(payload) {
			continue
		}
		c.MessageChan <- payload
	}
}

func (c *Client) handleMessage() {
	for msg := range c.MessageChan {
		switch msg.Type {
		case "READY":
			// 暂时忽略 READY
			log.Info("ws client receive ready msg")
		case model.EventAtMessageCreate:
			handleAtMessageCreate(msg)
		default:
			defaultHandle(msg)
		}
	}
	log.Infof("ws client handleMessage exit")
}

func (c *Client) updateSeq(seq uint32) {
	if seq > 0 && seq > c.LastSeq {
		c.LastSeq = seq
	}
}

// handleInsideMessage 处理一些默认信息
func (c *Client) handleInsideMessage(payload *model.WSPayload) bool {
	switch payload.OPCode {
	case model.WSHello:
		c.startHeart(payload)
	case model.WSHeartbeatAck:
		c.heartAck(payload)
	case model.WSReconnect:
		c.CloseChan <- errors.New("ws client reconnect")
	case model.WSInvalidSession:
		// 无效的 sessionLog，需要重新鉴权
		c.CloseChan <- errors.New("ws client invalid session")
	default:
		return false
	}
	return true
}

// startHeart 更新心跳时间
func (c *Client) startHeart(payload *model.WSPayload) {
	heartMsg := model.WSHelloData{}

	if err := parseData(payload, &heartMsg); err != nil {
		log.Errorf("ws client startHeart parseData error: %s", err.Error())
		return
	}
	c.Heart.Reset(time.Duration(heartMsg.HeartbeatInterval) * time.Millisecond)
}

func (c *Client) heartAck(payload *model.WSPayload) {
	log.Infof("ws client receive heartAck")
}

// Identify 发送鉴权信息
func (c *Client) Identify() error {
	event := &model.WSPayload{
		Data: &model.WSIdentityData{
			Token:   c.Token,
			Intents: c.Intent,
		},
	}
	event.OPCode = model.WSIdentity

	_ = c.Write(event)
	return nil
}

func (c *Client) Write(event *model.WSPayload) error {
	m, _ := json.Marshal(event)

	if err := c.WssCon.WriteMessage(websocket.TextMessage, m); err != nil {
		log.Errorf("发送消息失败: %s", err.Error())
		c.CloseChan <- err
		return err
	}
	return nil
}

// parseData 解析数据
func parseData(payload *model.WSPayload, w interface{}) error {
	b, err := json.Marshal(payload.Data)
	if err != nil {
		log.Errorf("ws client startHeart json.Marshal error: %s", err.Error())
		return err
	}
	err = json.Unmarshal(b, &w)
	if err != nil {
		log.Errorf("startHeart json.Unmarshal error: %v", err)
		return err
	}
	return nil

}
