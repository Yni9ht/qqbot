package ws

import (
	"context"
	"fmt"
	"github.com/Yni9ht/qqbot/log"
	"github.com/Yni9ht/qqbot/model"
	"github.com/Yni9ht/qqbot/redis"
	"github.com/Yni9ht/qqbot/vars"
	"regexp"
	"strings"
	"time"
)

// 用于过滤 at 结构的正则
var atRE = regexp.MustCompile(`<@!\d+>`)

// 用于过滤用户发送消息中的空格符号，\u00A0 是 &nbsp; 的 unicode 编码，某些 mac/pc 版本，连续多个空格的时候会转换成这个符号发送到后台
const spaceCharSet = " \u00A0"

// 默认签到积分
const defaultSignScore = 1

const (
	// signSuccessMsgTemplate 签到成功消息模板
	signSuccessMsgTemplate = "恭喜 %s 打卡成功，获得 %d 积分，当前总积分: %d"
)

func handleAtMessageCreate(msg *model.WSPayload) {
	data := &model.WSATMessageData{}

	if err := parseData(msg, data); err != nil {
		log.Errorf("handleAtMessageCreate: %v", err)
		return
	}

	cmd, err := parseCommand(data.Content)
	if err != nil {
		log.Errorf("handleAtMessageCreate: %v", err)
		return
	}
	log.Debugf("handleAtMessageCreate, cmd:%v", cmd)

	switch cmd {
	case "/打卡":
		handleClockIn(data)
	default:
		handleDefaultAtMessage(data)
	}
}

// handleClockIn 处理打卡命令
func handleClockIn(data *model.WSATMessageData) {
	// 获取 RedisSvc
	ctx := context.TODO()
	redisSvc := redis.NewRedisSvc()
	defer redisSvc.Close()

	// 判断用户今天是否已经打卡
	now := time.Now()
	signFlag, err := redisSvc.CheckUserClockIn(ctx, data.Author.ID, now)
	if err != nil {
		log.Errorf("检查用户是否打卡失败 msgID: %v, userId: %v, time: %v", data.ID, data.Author.ID, now.Format("2006-01-02 15:04:05"))
	}
	if signFlag {
		// 已经打卡
		log.Warnf("用户已经打卡 msgID: %v, userId: %v, time: %v", data.ID, data.Author.ID, now.Format("2006-01-02 15:04:05"))
		return
	}

	// 创建用户今日打卡记录
	err = redisSvc.CreateUserClockIn(ctx, data.Author.ID, now)
	if err != nil {
		log.Errorf("用户打卡失败 CreateUserClockIn err: %v", err)
		return
	}
	log.Infof("用户打卡成功 userID: %v, time: %v", data.Author.ID, now.Format("2006-01-02 15:04:05"))

	// 累加积分记录
	signScore := defaultSignScore
	currentScoreTotal, err := redisSvc.IncrUserSignScore(ctx, data.Author.ID, signScore)
	if err != nil {
		log.Errorf("给用户更新积分失败 userID: %v, score: %v, err: %v", data.Author.ID, signScore, err)
		_ = redisSvc.CancelUserClockIn(ctx, data.Author.ID, now)
	}

	// 发送消息
	content := fmt.Sprintf(signSuccessMsgTemplate, data.Author.Username, signScore, currentScoreTotal)
	msg := model.OpenAPIMessageReq{
		Content: content,
		MsgID:   data.ID,
		MessageReference: &model.MessageReference{
			// 引用这条消息
			MessageID:             data.ID,
			IgnoreGetMessageError: true,
		},
	}
	err = vars.OpenAPI.SendMessage(data.ChannelID, msg)
	if err != nil {
		log.Errorf("发送用户打卡成功信息失败 handleClockIn: %v", err)
		_ = redisSvc.CancelUserClockIn(ctx, data.Author.ID, now)
		_ = redisSvc.CancelUserSignScore(ctx, data.Author.ID, signScore)
	}
}

// handleDefaultAtMessage 默认 AT 回复
func handleDefaultAtMessage(data *model.WSATMessageData) {
	msg := model.OpenAPIMessageReq{
		Content: "您好，我暂时无法处理该信息哦。",
		MsgID:   data.ID,
		MessageReference: &model.MessageReference{
			// 引用这条消息
			MessageID:             data.ID,
			IgnoreGetMessageError: true,
		},
	}
	err := vars.OpenAPI.SendMessage(data.ChannelID, msg)
	if err != nil {
		log.Errorf("handleDefaultAtMessage: %v", err)
	}
}

func ETLInput(input string) string {
	etlData := string(atRE.ReplaceAll([]byte(input), []byte("")))
	etlData = strings.Trim(etlData, spaceCharSet)
	return etlData
}
func parseCommand(content string) (string, error) {
	input := ETLInput(content)
	cmds := strings.Split(input, " ")

	return cmds[0], nil
}

func defaultHandle(event *model.WSPayload) {
	log.Infof("ws client defaultHandle %+v", event)
}
